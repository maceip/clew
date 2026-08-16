// Package llm shells out to a configured headless agent CLI or an
// OpenAI-compatible endpoint (JOURNAL_SPEC §6.3, bring-your-own-agent).
// Default "auto" picks the first available of claude → codex → openai.
package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"clew/internal/config"
)

type Result struct {
	Text   string
	Tokens int // provider-reported when available, else bytes/4 estimate
}

type Provider interface {
	Name() string
	Call(prompt string) (*Result, error)
}

// Pick resolves the configured provider. Returns nil (with a reason) when
// extraction is off/unavailable — a loud status line, never silent (I2).
func Pick(cfg *config.Config) (Provider, string) {
	ex := cfg.Extractor
	switch ex.Provider {
	case "off":
		return nil, "extractor off (config)"
	case "claude":
		return claudeIfPresent(ex)
	case "codex":
		return codexIfPresent(ex)
	case "openai":
		return openaiIfPresent(ex)
	case "command":
		if len(ex.Command) == 0 {
			return nil, "extractor provider=command but no command configured"
		}
		if err := validateNeutralCommand(ex.Command); err != nil {
			return nil, err.Error()
		}
		return &cmdProvider{argv: ex.Command}, ""
	default: // auto
		if p, _ := claudeIfPresent(ex); p != nil {
			return p, ""
		}
		if p, _ := codexIfPresent(ex); p != nil {
			return p, ""
		}
		if p, _ := openaiIfPresent(ex); p != nil {
			return p, ""
		}
		return nil, "no extraction provider available (claude/codex not in PATH, no " + ex.OpenAIKeyEnv + ")"
	}
}

func validateNeutralCommand(argv []string) error {
	if strings.ContainsRune(argv[0], os.PathSeparator) && !filepath.IsAbs(argv[0]) {
		return fmt.Errorf("extractor command executable %q must be absolute or PATH-resolved (provider calls run from a neutral directory)", argv[0])
	}
	for _, arg := range argv[1:] {
		if strings.HasPrefix(arg, "./") || strings.HasPrefix(arg, "../") {
			return fmt.Errorf("extractor command argument %q must be absolute (provider calls run from a neutral directory)", arg)
		}
	}
	return nil
}

func claudeIfPresent(ex config.Extractor) (Provider, string) {
	if _, err := exec.LookPath("claude"); err != nil {
		return nil, "claude not in PATH"
	}
	return &claudeCLI{model: ex.Model}, ""
}

func codexIfPresent(ex config.Extractor) (Provider, string) {
	if _, err := exec.LookPath("codex"); err != nil {
		return nil, "codex not in PATH"
	}
	return &codexCLI{model: ex.Model}, ""
}

func openaiIfPresent(ex config.Extractor) (Provider, string) {
	key := os.Getenv(ex.OpenAIKeyEnv)
	if key == "" {
		return nil, ex.OpenAIKeyEnv + " not set"
	}
	base := ex.OpenAIBase
	if base == "" {
		base = "https://api.openai.com/v1"
	}
	model := ex.Model
	if model == "" {
		model = "gpt-4o-mini"
	}
	return &openaiHTTP{base: base, key: key, model: model}, ""
}

const callTimeout = 5 * time.Minute

func runArgv(argv []string, stdin string) (string, error) {
	cmd := exec.Command(argv[0], argv[1:]...)
	// Provider CLIs often persist their own transcript keyed by cwd. Keep
	// extraction calls outside every registered project so the watcher cannot
	// observe and recursively extract its own prompts.
	cmd.Dir = os.TempDir()
	cmd.Stdin = strings.NewReader(stdin)
	var out, errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Start(); err != nil {
		return "", err
	}
	done := make(chan error, 1)
	go func() { done <- cmd.Wait() }()
	select {
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("%s: %w: %s", argv[0], err, firstKB(errb.String()))
		}
		return out.String(), nil
	case <-time.After(callTimeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("%s: timeout after %s", argv[0], callTimeout)
	}
}

func firstKB(s string) string {
	if len(s) > 1024 {
		return s[:1024]
	}
	return s
}

// ---- claude -p ----

type claudeCLI struct{ model string }

func (c *claudeCLI) Name() string { return "claude" }

func (c *claudeCLI) Call(prompt string) (*Result, error) {
	argv := []string{"claude", "-p", "--output-format", "json"}
	if c.model != "" {
		argv = append(argv, "--model", c.model)
	}
	out, err := runArgv(argv, prompt)
	if err != nil {
		return nil, err
	}
	return decodeClaudeResult(prompt, out), nil
}

func decodeClaudeResult(prompt, out string) *Result {
	var env struct {
		Result string `json:"result"`
		Usage  struct {
			Input         int `json:"input_tokens"`
			Output        int `json:"output_tokens"`
			CacheCreation int `json:"cache_creation_input_tokens"`
			CacheRead     int `json:"cache_read_input_tokens"`
		} `json:"usage"`
	}
	if err := json.Unmarshal([]byte(out), &env); err != nil || env.Result == "" {
		// Envelope drift: fall back to raw output, estimated tokens.
		return &Result{Text: out, Tokens: (len(prompt) + len(out)) / 4}
	}
	tok := env.Usage.Input + env.Usage.Output + env.Usage.CacheCreation + env.Usage.CacheRead
	if tok == 0 {
		tok = (len(prompt) + len(env.Result)) / 4
	}
	return &Result{Text: env.Result, Tokens: tok}
}

// ---- codex exec ----

type codexCLI struct{ model string }

func (c *codexCLI) Name() string { return "codex" }

func (c *codexCLI) Call(prompt string) (*Result, error) {
	tmp, err := os.CreateTemp("", "clew-codex-*.txt")
	if err != nil {
		return nil, err
	}
	tmp.Close()
	defer os.Remove(tmp.Name())
	argv := []string{"codex", "exec", "--skip-git-repo-check", "-s", "read-only",
		"--output-last-message", tmp.Name()}
	if c.model != "" {
		argv = append(argv, "-m", c.model)
	}
	argv = append(argv, "-") // prompt on stdin
	if _, err := runArgv(argv, prompt); err != nil {
		return nil, err
	}
	b, err := os.ReadFile(tmp.Name())
	if err != nil {
		return nil, err
	}
	return &Result{Text: string(b), Tokens: (len(prompt) + len(b)) / 4}, nil
}

// ---- OpenAI-compatible HTTP ----

type openaiHTTP struct{ base, key, model string }

func (o *openaiHTTP) Name() string { return "openai" }

func (o *openaiHTTP) Call(prompt string) (*Result, error) {
	body, _ := json.Marshal(map[string]any{
		"model":           o.model,
		"messages":        []map[string]string{{"role": "user", "content": prompt}},
		"response_format": map[string]string{"type": "json_object"},
	})
	req, err := http.NewRequest("POST", strings.TrimSuffix(o.base, "/")+"/chat/completions", bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+o.key)
	req.Header.Set("Content-Type", "application/json")
	cl := &http.Client{Timeout: callTimeout}
	resp, err := cl.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var env struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
		Usage struct {
			Total int `json:"total_tokens"`
		} `json:"usage"`
		Error *struct {
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}
	if env.Error != nil {
		return nil, fmt.Errorf("openai: %s", env.Error.Message)
	}
	if len(env.Choices) == 0 {
		return nil, fmt.Errorf("openai: empty response")
	}
	return &Result{Text: env.Choices[0].Message.Content, Tokens: env.Usage.Total}, nil
}

// ---- custom command ----

type cmdProvider struct{ argv []string }

func (c *cmdProvider) Name() string { return "command:" + c.argv[0] }

func (c *cmdProvider) Call(prompt string) (*Result, error) {
	out, err := runArgv(c.argv, prompt)
	if err != nil {
		return nil, err
	}
	return &Result{Text: out, Tokens: (len(prompt) + len(out)) / 4}, nil
}

// ExtractJSON pulls the first top-level JSON object out of provider text
// (providers wrap JSON in prose/fences at times; validation stays strict).
func ExtractJSON(s string) (string, bool) {
	start := strings.IndexByte(s, '{')
	if start < 0 {
		return "", false
	}
	depth, inStr, esc := 0, false, false
	for i := start; i < len(s); i++ {
		ch := s[i]
		if inStr {
			switch {
			case esc:
				esc = false
			case ch == '\\':
				esc = true
			case ch == '"':
				inStr = false
			}
			continue
		}
		switch ch {
		case '"':
			inStr = true
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return s[start : i+1], true
			}
		}
	}
	return "", false
}
