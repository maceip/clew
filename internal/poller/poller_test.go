package poller

import (
	"path/filepath"
	"testing"
	"time"

	"github.com/maceip/clew/internal/state"
)

func TestAttributionRequiresTimeAndFootprintAndChoosesBestOverlap(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	repo := filepath.Join(t.TempDir(), "repo")
	now := time.Now().UTC()
	for _, session := range []state.Session{
		{ID: "one", Agent: "codex", RepoPath: repo, StartedAt: now.Add(-time.Hour), LastActivity: now},
		{ID: "two", Agent: "claude", RepoPath: repo, StartedAt: now.Add(-time.Hour), LastActivity: now},
		{ID: "stale", Agent: "cursor", RepoPath: repo, StartedAt: now.Add(-4 * time.Hour), LastActivity: now.Add(-3 * time.Hour)},
	} {
		if err := db.UpsertSession(session); err != nil {
			t.Fatal(err)
		}
	}
	if err := db.AddFootprints("one", []string{filepath.Join(repo, "a.go")}); err != nil {
		t.Fatal(err)
	}
	if err := db.AddFootprints("two", []string{filepath.Join(repo, "a.go"), filepath.Join(repo, "b.go")}); err != nil {
		t.Fatal(err)
	}
	if err := db.AddFootprints("stale", []string{filepath.Join(repo, "a.go"), filepath.Join(repo, "b.go"), filepath.Join(repo, "c.go")}); err != nil {
		t.Fatal(err)
	}
	commit := &state.Commit{At: now.Add(-5 * time.Minute), Files: []string{"a.go", "b.go"}}
	if got := attribute(db, repo, commit); got != "two" {
		t.Fatalf("best attribution = %q, want two", got)
	}
	commit.Files = []string{"unseen.go"}
	if got := attribute(db, repo, commit); got != "" {
		t.Fatalf("no-overlap attribution = %q", got)
	}
	commit.Files = []string{"a.go"}
	commit.At = now.Add(-2 * time.Hour)
	if got := attribute(db, repo, commit); got != "" {
		t.Fatalf("outside-window attribution = %q", got)
	}
}
