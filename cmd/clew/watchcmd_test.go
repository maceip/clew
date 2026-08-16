package main

import (
	"os"
	"path/filepath"
	"testing"
	"time"
	"unicode/utf8"

	"clew/internal/adapters"
	"clew/internal/state"
)

func TestClipStrHonorsCharacterLimit(t *testing.T) {
	got := clipStr("one\ntwo three", 8)
	if got != "one two…" {
		t.Fatalf("clipStr() = %q, want %q", got, "one two…")
	}
	if n := utf8.RuneCountInString(got); n != 8 {
		t.Fatalf("clipStr() returned %d characters, want 8", n)
	}

	got = clipStr("ééé", 2)
	if got != "é…" {
		t.Fatalf("clipStr() split or miscounted Unicode: %q", got)
	}
}

func TestSupervisorPATHIncludesInstalledBinaryDirectory(t *testing.T) {
	t.Setenv("PATH", "/usr/bin"+string(os.PathListSeparator)+"/bin")
	got := supervisorPATH("/Users/test/.local/bin/clew")
	want := "/Users/test/.local/bin" + string(os.PathListSeparator) +
		"/usr/bin" + string(os.PathListSeparator) + "/bin"
	if got != want {
		t.Fatalf("supervisorPATH() = %q, want %q", got, want)
	}
}

func TestXMLTextEscapesPlistValues(t *testing.T) {
	if got := xmlText("a&b<c>d"); got != "a&amp;b&lt;c&gt;d" {
		t.Fatalf("xmlText() = %q", got)
	}
}

type discoverOnlyAdapter struct {
	files []string
}

func (a *discoverOnlyAdapter) ID() string               { return "fixture" }
func (a *discoverOnlyAdapter) Discover(string) []string { return a.files }
func (a *discoverOnlyAdapter) Parse(string, int64) (*adapters.Delta, error) {
	return nil, nil
}

func TestBaselineSessionFilesIsForwardOnlyAndIdempotent(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	file := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(file, []byte("old\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	adapter := &discoverOnlyAdapter{files: []string{file}}
	if n, err := baselineSessionFiles(db, "/repo", []adapters.Adapter{adapter}); err != nil || n != 1 {
		t.Fatalf("first baseline = %d, %v; want 1, nil", n, err)
	}
	if got := db.Watermark("tail:" + file); got != 4 {
		t.Fatalf("tail baseline = %d, want 4", got)
	}
	if got := db.Watermark("extract:" + file); got != 4 {
		t.Fatalf("watch extraction baseline = %d, want 4", got)
	}
	if got := db.Watermark("history-end:" + file); got != 4 {
		t.Fatalf("history boundary = %d, want 4", got)
	}
	if err := os.WriteFile(file, []byte("old\nnew\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if n, err := baselineSessionFiles(db, "/repo", []adapters.Adapter{adapter}); err != nil || n != 0 {
		t.Fatalf("repeat baseline = %d, %v; want 0, nil", n, err)
	}
	if got := db.Watermark("extract:" + file); got != 4 {
		t.Fatalf("repeat baseline advanced to %d, want 4", got)
	}
}

func TestMigrateLiveCursorsRunsOnceForRegisteredRepos(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.RegisterRepo("/repo", "origin"); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(file, []byte("old\npartial"), 0o600); err != nil {
		t.Fatal(err)
	}
	adapter := &discoverOnlyAdapter{files: []string{file}}
	if err := migrateLiveCursorsDB(db, []adapters.Adapter{adapter}); err != nil {
		t.Fatal(err)
	}
	for _, prefix := range []string{"tail:", "extract:", "history-end:"} {
		if got := db.Watermark(prefix + file); got != 4 {
			t.Fatalf("%s migration offset = %d, want safe baseline 4", prefix, got)
		}
	}
	if got := db.Watermark("backfill:" + file); got != 0 {
		t.Fatalf("backfill migration offset = %d, want explicit history from 0", got)
	}
	if err := os.WriteFile(file, []byte("old\npartial\nnew\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := migrateLiveCursorsDB(db, []adapters.Adapter{adapter}); err != nil {
		t.Fatal(err)
	}
	if got := db.Watermark("history-end:" + file); got != 4 {
		t.Fatalf("completed migration reran and advanced to %d", got)
	}
}

func TestMigrateLiveCursorsPreservesExistingPendingOffset(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.RegisterRepo("/repo", "origin"); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(file, []byte("a\nb\nc\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("tail:"+file, "fixture", "/repo", 6); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("extract:"+file, "fixture", "/repo", 2); err != nil {
		t.Fatal(err)
	}
	adapter := &discoverOnlyAdapter{files: []string{file}}
	if err := migrateLiveCursorsDB(db, []adapters.Adapter{adapter}); err != nil {
		t.Fatal(err)
	}
	for _, prefix := range []string{"extract:", "history-end:", "backfill:"} {
		if got := db.Watermark(prefix + file); got != 2 {
			t.Fatalf("%s offset = %d, want preserved legacy offset 2", prefix, got)
		}
	}
}

func TestMigrateLiveCursorsSeedsMissingTailFromLegacyExtract(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.RegisterRepo("/repo", "origin"); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(file, []byte("old\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("extract:"+file, "fixture", "/repo", 4); err != nil {
		t.Fatal(err)
	}
	adapter := &discoverOnlyAdapter{files: []string{file}}
	if err := migrateLiveCursorsDB(db, []adapters.Adapter{adapter}); err != nil {
		t.Fatal(err)
	}
	if got := db.Watermark("tail:" + file); got != 4 {
		t.Fatalf("missing legacy tail seeded to %d, want extract offset 4", got)
	}
}

func TestMigrateSplitDogfoodCursorWithoutReplayingLiveTail(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.RegisterRepo("/repo", "origin"); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(file, []byte("history\nlive\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("watch-extract:"+file, "fixture", "/repo", 13); err != nil {
		t.Fatal(err)
	}
	adapter := &discoverOnlyAdapter{files: []string{file}}
	if err := migrateLiveCursorsDB(db, []adapters.Adapter{adapter}); err != nil {
		t.Fatal(err)
	}
	for _, prefix := range []string{"extract:", "history-end:", "backfill:"} {
		if got := db.Watermark(prefix + file); got != 13 {
			t.Fatalf("%s offset = %d, want consumed legacy boundary 13", prefix, got)
		}
	}
	if db.Get("migration-note") == "" {
		t.Fatal("history tradeoff was silent")
	}
}

func TestMigrateSplitCursorNeverRewindsConsumedExtract(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.RegisterRepo("/repo", "origin"); err != nil {
		t.Fatal(err)
	}
	file := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(file, []byte("history\nlive\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("watch-extract:"+file, "fixture", "/repo", 8); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("extract:"+file, "fixture", "/repo", 13); err != nil {
		t.Fatal(err)
	}
	adapter := &discoverOnlyAdapter{files: []string{file}}
	if err := migrateLiveCursorsDB(db, []adapters.Adapter{adapter}); err != nil {
		t.Fatal(err)
	}
	if got := db.Watermark("extract:" + file); got != 13 {
		t.Fatalf("migration rewound consumed extract cursor to %d, want 13", got)
	}
	if got := db.Watermark("history-end:" + file); got != 13 {
		t.Fatalf("history boundary = %d, want consumed historical offset 13", got)
	}
}

func TestExtractionFailureBacksOffAndIsVisible(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	w := &watcher{a: &app{db: db}}
	p := &pendState{}
	w.failExtraction("/tmp/session.jsonl", p, os.ErrPermission)
	firstRetry := p.retryAt
	if p.failures != 1 || !firstRetry.After(time.Now()) {
		t.Fatalf("first failure state = %+v", p)
	}
	if got := db.Get("extract-error:/tmp/session.jsonl"); got == "" {
		t.Fatal("provider failure was not persisted for status")
	}
	w.failExtraction("/tmp/session.jsonl", p, os.ErrPermission)
	if !p.retryAt.After(firstRetry) {
		t.Fatalf("retry did not back off: first=%s second=%s", firstRetry, p.retryAt)
	}
}

func TestLiveExtractionOffsetsResetOnRotation(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	file := "/tmp/rotated-session.jsonl"
	if err := db.SetWatermark("extract:"+file, "fixture", "/repo", 100); err != nil {
		t.Fatal(err)
	}
	if err := db.SetWatermark("tail:"+file, "fixture", "/repo", 12); err != nil {
		t.Fatal(err)
	}
	extractOff, tailOff, rotated, err := liveExtractionOffsets(db, &discoverOnlyAdapter{}, "/repo", file)
	if err != nil {
		t.Fatal(err)
	}
	if !rotated || extractOff != 0 || tailOff != 12 {
		t.Fatalf("rotation offsets = %d,%d rotated=%v", extractOff, tailOff, rotated)
	}
	if got := db.Watermark("extract:" + file); got != 0 {
		t.Fatalf("persisted extract cursor = %d, want 0", got)
	}
}

func TestDeltaTimesUsesEvidenceTime(t *testing.T) {
	fallback := time.Date(2026, 8, 16, 17, 0, 0, 0, time.UTC)
	early := fallback.Add(-2 * time.Hour)
	late := fallback.Add(-time.Hour)
	delta := &adapters.Delta{Messages: []adapters.Message{{At: late}, {At: early}}}
	started, last := deltaTimes(delta, fallback)
	if !started.Equal(early) || !last.Equal(late) {
		t.Fatalf("deltaTimes() = %s, %s; want %s, %s", started, last, early, late)
	}
}
