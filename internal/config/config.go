// Package config loads ~/.clew/config.yaml. Everything has a working
// default; the file is optional.
package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/maceip/clew/internal/gitx"
)

type Push struct {
	Kind string `yaml:"kind"` // ntfy | webhook (§12.2: one config line either way)
	URL  string `yaml:"url"`
}

// Owner configures the owner's project-agnostic law journal. The local store
// always works without a remote; when Remote is set it uses the same git sync
// protocol and credentials as project journals.
type Owner struct {
	Remote string `yaml:"remote"`
}

type Extractor struct {
	// Provider: auto | claude | codex | openai | command | off.
	// auto picks the first available of claude → codex → openai (§6.3:
	// "cheapest configured"; order documented in README).
	Provider string `yaml:"provider"`
	// Command: explicit argv; prompt on stdin, strict JSON on stdout.
	Command []string `yaml:"command,omitempty"`
	Model   string   `yaml:"model,omitempty"`

	OpenAIBase   string `yaml:"openai_base,omitempty"`
	OpenAIKeyEnv string `yaml:"openai_key_env,omitempty"`

	// Budget (I9, §6.4).
	DailyCapTokens int `yaml:"daily_cap_tokens"`
}

type Config struct {
	Surface   string    `yaml:"surface"` // machine label
	Push      Push      `yaml:"push"`
	Owner     Owner     `yaml:"owner"`
	Extractor Extractor `yaml:"extractor"`
	// LinkPass enables the batched LLM link pass (§7.1.2).
	LinkPass bool `yaml:"link_pass"`
}

func Path() string { return filepath.Join(gitx.Home(), "config.yaml") }

func Default() *Config {
	host, _ := os.Hostname()
	if host == "" {
		host = "local"
	}
	return &Config{
		Surface: shortHost(host),
		Extractor: Extractor{
			Provider:       "auto",
			OpenAIKeyEnv:   "OPENAI_API_KEY",
			DailyCapTokens: 200_000, // owner ceiling; recording itself is unmetered
		},
		LinkPass: true,
	}
}

func Load() *Config {
	c := Default()
	b, err := os.ReadFile(Path())
	if err != nil {
		return c
	}
	_ = yaml.Unmarshal(b, c)
	if c.Extractor.DailyCapTokens <= 0 {
		c.Extractor.DailyCapTokens = 200_000
	}
	if c.Extractor.OpenAIKeyEnv == "" {
		c.Extractor.OpenAIKeyEnv = "OPENAI_API_KEY"
	}
	if c.Surface == "" {
		c.Surface = Default().Surface
	}
	return c
}

// WriteSkeleton writes a commented default config if none exists.
func WriteSkeleton() error {
	p := Path()
	if _, err := os.Stat(p); err == nil {
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		return err
	}
	c := Default()
	b, _ := yaml.Marshal(c)
	head := "# clew config. Defaults shown; push.url enables phone push (ntfy/webhook).\n"
	return os.WriteFile(p, append([]byte(head), b...), 0o644)
}

func shortHost(h string) string {
	for i := 0; i < len(h); i++ {
		if h[i] == '.' {
			return h[:i]
		}
	}
	return h
}
