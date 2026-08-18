package main

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/lineage"
	"clew/internal/model"
	"clew/internal/seed"
)

func captureFromStdout(t *testing.T, args []string) (string, error) {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	old := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = old }()
	callErr := cmdFrom(args)
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	b, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	_ = r.Close()
	return string(b), callErr
}

func chdirForTest(t *testing.T, dir string) {
	t.Helper()
	old, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func TestFromWithoutArgsDoesNotBirthOrMutateTheProject(t *testing.T) {
	home := t.TempDir()
	clewHome := filepath.Join(home, "clew-home")
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", clewHome)
	repo := filepath.Join(t.TempDir(), "new-topic")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	canonicalRepo, err := gitx.Root(repo)
	if err != nil {
		t.Fatal(err)
	}
	repo = canonicalRepo
	before, err := gitx.Run(repo, "status", "--porcelain")
	if err != nil {
		t.Fatal(err)
	}
	chdirForTest(t, repo)
	out, err := captureFromStdout(t, nil)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "nothing is carried") {
		t.Fatalf("candidate listing did not make its read-only behavior clear: %q", out)
	}
	after, err := gitx.Run(repo, "status", "--porcelain")
	if err != nil {
		t.Fatal(err)
	}
	if before != after {
		t.Fatalf("no-argument listing mutated project status: %q -> %q", before, after)
	}
	if _, err := os.Stat(filepath.Join(repo, ".clew")); !os.IsNotExist(err) {
		t.Fatalf("no-argument listing birthed project .clew: %v", err)
	}
	if _, err := os.Stat(gitx.WorktreeDir(repo)); !os.IsNotExist(err) {
		t.Fatalf("no-argument listing created a project journal worktree: %v", err)
	}
	if _, err := os.Stat(clewHome); !os.IsNotExist(err) {
		t.Fatalf("empty-home candidate listing created machine config/state: %v", err)
	}
}

func TestMaintainedSeedPrefersReadyCanonicalJournalBranch(t *testing.T) {
	t.Setenv("CLEW_HOME", filepath.Join(t.TempDir(), "clew-home"))
	repo := t.TempDir()
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	wt, err := gitx.EnsureJournal(repo)
	if err != nil {
		t.Fatal(err)
	}
	base := seed.Snapshot{
		Repository: seed.Repository{ID: "repo:source", Name: "source"},
		ChangedAt:  time.Date(2026, 8, 18, 19, 0, 0, 0, time.UTC),
		Lifecycle:  seed.Lifecycle{State: "active"},
	}
	stale := base
	stale.JournalRevision = "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	canonical := base
	canonical.JournalRevision = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	canonical.Ancestors = []string{"repo:latest-ancestor"}
	if _, err := seed.Write(filepath.Join(repo, ".clew", "SEED.md"), &stale); err != nil {
		t.Fatal(err)
	}
	if _, err := seed.Write(filepath.Join(wt, "SEED.md"), &canonical); err != nil {
		t.Fatal(err)
	}
	got, err := maintainedSeed(repo)
	if err != nil {
		t.Fatal(err)
	}
	if got.JournalRevision != canonical.JournalRevision || strings.Join(got.Ancestors, ",") != "repo:latest-ancestor" {
		t.Fatalf("selected stale workspace seed over canonical journal branch: %#v", got)
	}
}

func TestExplicitFromReadsMaintainedSourceWithoutMutatingIt(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	base := t.TempDir()
	sourceRepo := filepath.Join(base, "predecessor")
	successorRepo := filepath.Join(base, "successor")
	for _, repo := range []string{sourceRepo, successorRepo} {
		if err := os.MkdirAll(repo, 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := gitx.Run(repo, "init", "-q"); err != nil {
			t.Fatal(err)
		}
	}
	var err error
	sourceRepo, err = gitx.Root(sourceRepo)
	if err != nil {
		t.Fatal(err)
	}
	successorRepo, err = gitx.Root(successorRepo)
	if err != nil {
		t.Fatal(err)
	}
	at := time.Date(2026, 8, 18, 20, 0, 0, 0, time.UTC)
	lesson := model.Entry{
		ID: ids.NewEntry(at), Type: model.Decision, Title: "keep predecessor lesson",
		Body:  "The explicit lineage ceremony selected this decision.",
		Quote: "carry this decision when I explicitly choose the predecessor", UtteranceBy: model.ByUser,
		Source: model.Source{Kind: model.SrcHuman, Ref: "test:predecessor", At: at}, Confidence: 1,
	}
	snapshot := &seed.Snapshot{
		Repository:      seed.Repository{ID: "repo:predecessor", Name: "predecessor"},
		JournalRevision: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		ChangedAt:       at, Lifecycle: seed.Lifecycle{State: "active"},
		Decisions: []seed.Lesson{{Entry: lesson, Status: journal.StActive}},
	}
	sourceSeed := filepath.Join(sourceRepo, ".clew", "SEED.md")
	if _, err := seed.Write(sourceSeed, snapshot); err != nil {
		t.Fatal(err)
	}
	beforeBytes, err := os.ReadFile(sourceSeed)
	if err != nil {
		t.Fatal(err)
	}
	beforeInfo, err := os.Stat(sourceSeed)
	if err != nil {
		t.Fatal(err)
	}
	beforeStatus, _ := gitx.Run(sourceRepo, "status", "--porcelain")
	chdirForTest(t, successorRepo)
	out, err := captureFromStdout(t, []string{sourceRepo})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "declared lineage successor") {
		t.Fatalf("explicit ceremony output = %q", out)
	}
	afterBytes, err := os.ReadFile(sourceSeed)
	if err != nil {
		t.Fatal(err)
	}
	afterInfo, err := os.Stat(sourceSeed)
	if err != nil {
		t.Fatal(err)
	}
	afterStatus, _ := gitx.Run(sourceRepo, "status", "--porcelain")
	if !bytes.Equal(beforeBytes, afterBytes) || !beforeInfo.ModTime().Equal(afterInfo.ModTime()) || beforeStatus != afterStatus {
		t.Fatalf("explicit import mutated predecessor: bytes=%v mtime=%v status=%q -> %q",
			bytes.Equal(beforeBytes, afterBytes), beforeInfo.ModTime().Equal(afterInfo.ModTime()), beforeStatus, afterStatus)
	}
	if _, err := os.Stat(gitx.WorktreeDir(sourceRepo)); !os.IsNotExist(err) {
		t.Fatalf("explicit import generated a predecessor journal on demand: %v", err)
	}
	successorJournal, err := journal.Load(gitx.WorktreeDir(successorRepo))
	if err != nil {
		t.Fatal(err)
	}
	if successorJournal.Entries[lesson.ID] == nil {
		t.Fatalf("explicit import did not carry predecessor decision %s", lesson.ID)
	}
	links, err := lineage.LoadLinks(successorJournal.Dir)
	if err != nil || len(links) != 1 || links[0].From.ID != snapshot.Repository.ID {
		t.Fatalf("durable successor lineage = %#v, err=%v", links, err)
	}
}
