package main

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/scrub"
)

func TestProjectRedactionLeaseRetryPreservesConcurrentJournalAppend(t *testing.T) {
	t.Setenv("CLEW_HOME", t.TempDir())
	base := t.TempDir()
	remote := filepath.Join(base, "remote.git")
	if _, err := gitx.Run(base, "init", "-q", "--bare", remote); err != nil {
		t.Fatal(err)
	}
	repoA := filepath.Join(base, "a")
	repoB := filepath.Join(base, "b")
	for _, repo := range []string{repoA, repoB} {
		if _, err := gitx.Run(base, "clone", "-q", remote, repo); err != nil {
			t.Fatal(err)
		}
	}
	wtA, err := gitx.EnsureJournal(repoA)
	if err != nil {
		t.Fatal(err)
	}
	wtB, err := gitx.EnsureJournal(repoB)
	if err != nil {
		t.Fatal(err)
	}

	stamp := time.Date(2026, 8, 18, 18, 0, 0, 0, time.UTC)
	secret := "project-redaction-lease-secret"
	target := addRedactionTestEntry(t, wtA, "Compromised project fact", "exact "+secret, stamp)
	if _, err := gitx.Sync(repoA, regenFor(repoA)); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Sync(repoB, regenFor(repoB)); err != nil {
		t.Fatal(err)
	}
	observed, err := gitx.Sync(repoA, regenFor(repoA))
	if err != nil {
		t.Fatal(err)
	}
	lease, err := gitx.RewriteLeaseFromSync(repoA, observed)
	if err != nil {
		t.Fatal(err)
	}
	journalA, err := journal.Load(wtA)
	if err != nil {
		t.Fatal(err)
	}

	concurrent := addRedactionTestEntry(t, wtB, "Concurrent project fact", "independent append survives", stamp.Add(time.Second))
	if _, err := gitx.Sync(repoB, regenFor(repoB)); err != nil {
		t.Fatal(err)
	}
	if err := rewriteProjectRedactedRoot(repoA, journalA, target.ID, "desk", stamp.Add(time.Minute), lease); err != nil {
		t.Fatal(err)
	}

	if _, err := gitx.Sync(repoB, regenFor(repoB)); err != nil {
		t.Fatal(err)
	}
	peer, err := journal.Load(wtB)
	if err != nil {
		t.Fatal(err)
	}
	gotTarget := peer.Entries[target.ID]
	if gotTarget == nil || gotTarget.Title != scrub.Mark || gotTarget.Body != scrub.Mark || gotTarget.Quote != scrub.Mark {
		t.Fatalf("target was not scrubbed after project lease retry: %+v", gotTarget)
	}
	if got := peer.Entries[concurrent.ID]; got == nil || got.Quote != concurrent.Quote {
		t.Fatalf("concurrent project append was erased: got=%+v want=%+v", got, concurrent)
	}
	out, err := exec.Command("git", "--git-dir", remote, "log", "--all", "-S"+secret, "--format=%H").CombinedOutput()
	if err != nil {
		t.Fatalf("inspect rewritten project history: %v: %s", err, out)
	}
	if strings.TrimSpace(string(out)) != "" {
		t.Fatalf("project secret remains in reachable history after lease retry: %s", out)
	}
}

func addRedactionTestEntry(t *testing.T, wt, title, quote string, at time.Time) *model.Entry {
	t.Helper()
	j, err := journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}
	e := &model.Entry{
		ID: ids.NewEntry(at), Type: model.Finding, Title: title, Body: "generated summary", Quote: quote,
		UtteranceBy: model.ByUser,
		Source:      model.Source{Kind: model.SrcSession, Ref: "test-session#L1", At: at},
		Confidence:  0.9,
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	return e
}
