package llm

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/maceip/clew/internal/config"
)

func TestRunArgvUsesNeutralWorkingDirectory(t *testing.T) {
	out, err := runArgv([]string{"pwd"}, "")
	if err != nil {
		t.Fatal(err)
	}
	got, err := filepath.EvalSymlinks(strings.TrimSpace(out))
	if err != nil {
		t.Fatal(err)
	}
	want, err := filepath.EvalSymlinks(os.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("provider cwd = %q, want neutral %q", got, want)
	}
}

func TestCustomCommandRejectsRelativePathsBeforeNeutralCWD(t *testing.T) {
	cfg := config.Default()
	cfg.Extractor.Provider = "command"
	cfg.Extractor.Command = []string{"./bin/extractor", "../schema.json"}
	p, note := Pick(cfg)
	if p != nil || !strings.Contains(note, "must be absolute") {
		t.Fatalf("Pick() = %v, %q; want loud absolute-path requirement", p, note)
	}
}

func TestClaudeUsageIncludesCacheTokens(t *testing.T) {
	out := `{"result":"{\"entries\":[]}","usage":{"input_tokens":10,"output_tokens":5,"cache_creation_input_tokens":20,"cache_read_input_tokens":30}}`
	got := decodeClaudeResult("prompt", out)
	if got.Tokens != 65 {
		t.Fatalf("Claude usage = %d, want input+output+cache = 65", got.Tokens)
	}
}
