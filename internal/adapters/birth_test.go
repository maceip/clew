package adapters

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func writeBirthFixture(t *testing.T, path, content string, at time.Time) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.Chtimes(path, at, at); err != nil {
		t.Fatal(err)
	}
}

func TestBirthCandidatesFindRecentClaudeCodexAndWrapSessions(t *testing.T) {
	home := t.TempDir()
	clewHome := filepath.Join(home, "clew-state")
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", clewHome)
	now := time.Now().UTC().Truncate(time.Second)

	claudeOld := filepath.Join(home, ".claude", "projects", "-repo-claude", "old.jsonl")
	writeBirthFixture(t, claudeOld,
		`{"type":"user","cwd":"/repo/stale","sessionId":"old"}`+"\n", now.Add(-2*time.Hour))
	claude := filepath.Join(home, ".claude", "projects", "-repo-claude", "new.jsonl")
	writeBirthFixture(t, claude,
		`{"type":"mode","mode":"normal","sessionId":"c1"}`+"\n"+
			`{"type":"user","cwd":"/repo/claude","sessionId":"c1"}`+"\n", now.Add(-time.Minute))

	day := time.Now().Format("2006/01/02")
	codex := filepath.Join(home, ".codex", "sessions", day, "rollout-birth.jsonl")
	writeBirthFixture(t, codex,
		`{"timestamp":"2026-08-18T12:00:00Z","type":"session_meta","payload":{"id":"cx","cwd":"/repo/codex"}}`+"\n",
		now.Add(-2*time.Minute))

	wrap := filepath.Join(clewHome, "raw", "wrap-birth.jsonl")
	writeBirthFixture(t, wrap,
		`{"kind":"meta","argv":["gemini"],"cwd":"/repo/wrap","session":"w1","at":"2026-08-18T12:00:00Z"}`+"\n",
		now.Add(-3*time.Minute))

	got := BirthCandidates(now.Add(-time.Hour))
	if len(got) != 3 {
		t.Fatalf("BirthCandidates() returned %d, want 3: %+v", len(got), got)
	}
	want := []struct {
		adapter string
		cwd     string
	}{
		{"claude-code", "/repo/claude"},
		{"codex", "/repo/codex"},
		{"wrap", "/repo/wrap"},
	}
	for i := range want {
		if got[i].Adapter != want[i].adapter || got[i].CWD != want[i].cwd {
			t.Errorf("candidate %d = %+v, want adapter=%s cwd=%s", i, got[i], want[i].adapter, want[i].cwd)
		}
		if got[i].At.Location() != time.UTC {
			t.Errorf("candidate %d time is not UTC: %s", i, got[i].At)
		}
	}
}

func TestBirthCandidatesKeepNewestSessionPerCWD(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-state"))
	now := time.Now().UTC().Truncate(time.Second)
	dir := filepath.Join(home, ".claude", "projects", "-same")
	older := filepath.Join(dir, "older.jsonl")
	newer := filepath.Join(dir, "newer.jsonl")
	content := `{"type":"user","cwd":"/repo/same","sessionId":"s"}` + "\n"
	writeBirthFixture(t, older, content, now.Add(-2*time.Minute))
	writeBirthFixture(t, newer, content, now.Add(-time.Minute))

	got := BirthCandidates(now.Add(-time.Hour))
	if len(got) != 1 {
		t.Fatalf("deduplicated candidates = %d, want 1: %+v", len(got), got)
	}
	if got[0].File != newer {
		t.Fatalf("kept %s, want newest session %s", got[0].File, newer)
	}
}

func TestBirthCandidatesSkipMalformedAndRelativeCWD(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-state"))
	now := time.Now().UTC().Truncate(time.Second)
	dir := filepath.Join(home, ".claude", "projects", "-bad")
	writeBirthFixture(t, filepath.Join(dir, "malformed.jsonl"), "{bad}\n", now)
	writeBirthFixture(t, filepath.Join(dir, "relative.jsonl"),
		`{"type":"user","cwd":"relative/repo","sessionId":"s"}`+"\n", now)
	if got := BirthCandidates(now.Add(-time.Minute)); len(got) != 0 {
		t.Fatalf("invalid candidates were returned: %+v", got)
	}
}

func TestBirthCandidatesHonorClaudeConfigDir(t *testing.T) {
	home := t.TempDir()
	custom := filepath.Join(t.TempDir(), "claude-custom")
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-state"))
	t.Setenv("CLAUDE_CONFIG_DIR", custom)
	now := time.Now().UTC().Truncate(time.Second)
	file := filepath.Join(custom, "projects", "-custom", "session.jsonl")
	writeBirthFixture(t, file, `{"type":"user","cwd":"/repo/custom","sessionId":"s"}`+"\n", now)
	got := BirthCandidates(now.Add(-time.Minute))
	if len(got) != 1 || got[0].File != file || got[0].CWD != "/repo/custom" {
		t.Fatalf("custom Claude config candidate = %+v", got)
	}
}
