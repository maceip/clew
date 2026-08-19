package gitx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
)

// mkRemote creates a bare "remote" and two clones ("machines" A and B).
func mkRemote(t *testing.T) (bare, repoA, repoB string) {
	t.Helper()
	t.Setenv("CLEW_HOME", t.TempDir())
	base := t.TempDir()
	bare = filepath.Join(base, "remote.git")
	if _, err := Run(base, "init", "-q", "--bare", bare); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"a", "b"} {
		p := filepath.Join(base, name)
		if _, err := Run(base, "clone", "-q", bare, p); err != nil {
			t.Fatal(err)
		}
		if name == "a" {
			repoA = p
		} else {
			repoB = p
		}
	}
	return
}

func addEntry(t *testing.T, wt, title string) string {
	t.Helper()
	j, err := journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}
	e := &model.Entry{
		ID: ids.NewEntry(time.Now()), Type: model.Intent, Title: title,
		Quote: "q: " + title, UtteranceBy: model.ByUser,
		Source:     model.Source{Kind: model.SrcSession, Ref: "test#L1", At: time.Now().UTC()},
		Confidence: 0.9,
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	return e.ID
}

func regen(wt string) error {
	j, err := journal.Load(wt)
	if err != nil {
		return err
	}
	st := journal.Compute(j, time.Now())
	return journal.WriteProjections(j, st, time.Now())
}

func TestTwoMachineUnionSync(t *testing.T) {
	_, repoA, repoB := mkRemote(t)

	wtA, err := EnsureJournal(repoA)
	if err != nil {
		t.Fatal(err)
	}
	// B bootstraps after A pushed: must adopt A's root, not race a new one.
	wtB, err := EnsureJournal(repoB)
	if err != nil {
		t.Fatal(err)
	}

	idA := addEntry(t, wtA, "from machine A")
	idB := addEntry(t, wtB, "from machine B")
	if _, err := Sync(repoA, regen); err != nil {
		t.Fatal(err)
	}
	// B syncs: rebase over A's push; projections conflict is auto-resolved.
	if _, err := Sync(repoB, regen); err != nil {
		t.Fatal(err)
	}
	if _, err := Sync(repoA, regen); err != nil {
		t.Fatal(err)
	}
	for _, wt := range []string{wtA, wtB} {
		for _, id := range []string{idA, idB} {
			if _, err := os.Stat(filepath.Join(wt, "entries", id+".yaml")); err != nil {
				t.Errorf("union violated: %s missing in %s", id, wt)
			}
		}
	}
}

func TestUnrelatedRootAdoption(t *testing.T) {
	// Two machines race init with NO remote branch check possible in between:
	// simulate by creating journals before either pushes.
	_, repoA, repoB := mkRemote(t)
	wtA, _ := EnsureJournal(repoA)
	// Sever B's view: create its journal while remote is still empty.
	wtB, _ := EnsureJournal(repoB)
	addEntry(t, wtA, "A wins the race")
	addEntry(t, wtB, "B loses the race")
	if _, err := Sync(repoA, regen); err != nil {
		t.Fatal(err)
	}
	res, err := Sync(repoB, regen)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Adopted && !res.Pulled {
		t.Fatalf("B should have adopted or rebased; got %+v", res)
	}
	// Both entries must survive on the remote.
	if _, err := Sync(repoA, regen); err != nil {
		t.Fatal(err)
	}
	entries, _ := os.ReadDir(filepath.Join(wtA, "entries"))
	if len(entries) != 2 {
		t.Fatalf("want 2 entries after race, got %d", len(entries))
	}
}

func TestRedactRewriteDoesNotResurrect(t *testing.T) {
	_, repoA, repoB := mkRemote(t)
	wtA, _ := EnsureJournal(repoA)
	wtB, _ := EnsureJournal(repoB)

	secretID := addEntry(t, wtA, "leaked key entry")
	// Poison the quote with a secret.
	p := filepath.Join(wtA, "entries", secretID+".yaml")
	b, _ := os.ReadFile(p)
	os.Remove(p)
	os.WriteFile(p, []byte(strings.Replace(string(b), "q: leaked key entry", "AKIAIOSFODNN7EXAMPLE", 1)), 0o644)
	if _, err := Sync(repoA, regen); err != nil {
		t.Fatal(err)
	}
	if _, err := Sync(repoB, regen); err != nil { // B now has the secret locally
		t.Fatal(err)
	}
	lastSync, err := Sync(repoA, regen)
	if err != nil {
		t.Fatal(err)
	}
	lease, err := RewriteLeaseFromSync(repoA, lastSync)
	if err != nil {
		t.Fatal(err)
	}

	// Redact on A: scrub in place, rewrite root, force-push.
	b2, _ := os.ReadFile(p)
	os.WriteFile(p, []byte(strings.Replace(string(b2), "AKIAIOSFODNN7EXAMPLE", "‹redacted›", 1)), 0o644)
	if err := RewriteRoot(repoA, "redact "+secretID, lease); err != nil {
		t.Fatal(err)
	}

	// B syncs: unrelated root → adopt; scrubbed file must NOT be re-added
	// from B's poisoned local copy.
	res, err := Sync(repoB, regen)
	if err != nil {
		t.Fatal(err)
	}
	if !res.Adopted {
		t.Fatalf("B must adopt the rewritten root, got %+v", res)
	}
	got, _ := os.ReadFile(filepath.Join(wtB, "entries", secretID+".yaml"))
	if strings.Contains(string(got), "AKIA") {
		t.Fatal("secret resurrected on machine B after redaction")
	}
	// And the remote history must be a single fresh root.
	if _, err := Sync(repoA, regen); err != nil {
		t.Fatal(err)
	}
	log, _ := Run(wtA, "log", "--format=%s", Branch)
	if strings.Contains(log, "leaked") {
		t.Fatal("old history survived the rewrite")
	}
}

func TestRewriteRootLeaseRejectsConcurrentAppendWithoutResettingLocalBranch(t *testing.T) {
	_, repoA, repoB := mkRemote(t)
	wtA, err := EnsureJournal(repoA)
	if err != nil {
		t.Fatal(err)
	}
	wtB, err := EnsureJournal(repoB)
	if err != nil {
		t.Fatal(err)
	}
	targetID := addEntry(t, wtA, "secret to scrub")
	if _, err := Sync(repoA, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := Sync(repoB, nil); err != nil {
		t.Fatal(err)
	}
	observed, err := Sync(repoA, nil)
	if err != nil {
		t.Fatal(err)
	}
	lease, err := RewriteLeaseFromSync(repoA, observed)
	if err != nil {
		t.Fatal(err)
	}
	headBefore, err := Run(wtA, "rev-parse", "HEAD")
	if err != nil {
		t.Fatal(err)
	}

	concurrentID := addEntry(t, wtB, "concurrent law")
	if _, err := Sync(repoB, nil); err != nil {
		t.Fatal(err)
	}
	p := filepath.Join(wtA, "entries", targetID+".yaml")
	b, err := os.ReadFile(p)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(strings.Replace(string(b), "secret to scrub", "‹redacted›", -1)), 0o644); err != nil {
		t.Fatal(err)
	}
	err = RewriteRoot(repoA, "redact "+targetID, lease)
	if !IsRewriteLeaseError(err) {
		t.Fatalf("stale rewrite lease error = %v, want RewriteLeaseError", err)
	}
	headAfter, headErr := Run(wtA, "rev-parse", "HEAD")
	if headErr != nil {
		t.Fatal(headErr)
	}
	if headAfter != headBefore {
		t.Fatalf("lease rejection reset local branch: before=%s after=%s", headBefore, headAfter)
	}
	if _, err := Run(repoB, "fetch", "-q", "origin", Branch); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(repoB, "show", "FETCH_HEAD:entries/"+concurrentID+".yaml"); err != nil {
		t.Fatalf("concurrent append was erased by stale rewrite: %v", err)
	}
}

func TestLocalOnlyNoRemote(t *testing.T) {
	t.Setenv("CLEW_HOME", t.TempDir())
	repo := t.TempDir()
	if _, err := Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	wt, err := EnsureJournal(repo)
	if err != nil {
		t.Fatal(err)
	}
	addEntry(t, wt, "works offline")
	res, err := Sync(repo, regen)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, n := range res.Notes {
		if strings.Contains(n, "local-only") {
			found = true
		}
	}
	if !found {
		t.Errorf("local-only mode must be noted (I2), got %+v", res.Notes)
	}
	for _, projection := range []string{"journal.md", "digest.md"} {
		if _, err := os.Stat(filepath.Join(wt, projection)); err != nil {
			t.Errorf("local-only sync did not regenerate %s: %v", projection, err)
		}
	}
}

func TestOfflineFetchStillRegeneratesAndCommitsAmbientSeed(t *testing.T) {
	t.Setenv("CLEW_HOME", t.TempDir())
	repo := t.TempDir()
	if _, err := Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	missingRemote := filepath.Join(t.TempDir(), "does-not-exist.git")
	if _, err := Run(repo, "remote", "add", "origin", missingRemote); err != nil {
		t.Fatal(err)
	}
	wt, err := EnsureJournal(repo)
	if err != nil {
		t.Fatal(err)
	}
	addEntry(t, wt, "journal change before offline fetch")
	content := []byte("ambient seed at local journal revision\n")
	res, err := Sync(repo, func(branch string) error {
		return os.WriteFile(filepath.Join(branch, "SEED.md"), content, 0o644)
	})
	if err != nil {
		t.Fatal(err)
	}
	fetchFailed := false
	for _, note := range res.Notes {
		if strings.Contains(note, "fetch failed") {
			fetchFailed = true
		}
	}
	if !fetchFailed {
		t.Fatalf("missing loud offline fetch note: %#v", res.Notes)
	}
	got, err := Show(repo, "SEED.md")
	if err != nil {
		t.Fatal(err)
	}
	if got != strings.TrimSpace(string(content)) {
		t.Fatalf("journal branch did not commit regenerated seed before returning: %q", got)
	}
	if status, err := Run(wt, "status", "--porcelain"); err != nil || status != "" {
		t.Fatalf("regenerated seed was not committed: status=%q err=%v", status, err)
	}
}

func TestBirthIncarnationEvidenceSurvivesDamagedExternalWorktree(t *testing.T) {
	t.Setenv("CLEW_HOME", t.TempDir())
	repo := t.TempDir()
	if _, err := Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	if HasBirthIncarnation(repo) {
		t.Fatal("fresh git repository unexpectedly has clew incarnation evidence")
	}
	wt, err := EnsureJournal(repo)
	if err != nil {
		t.Fatal(err)
	}
	if !JournalReady(repo) || !HasBirthIncarnation(repo) {
		t.Fatal("journal enrollment did not establish local incarnation evidence")
	}
	gitFile := filepath.Join(wt, ".git")
	if err := os.Rename(gitFile, gitFile+".damaged"); err != nil {
		t.Fatal(err)
	}
	if JournalReady(repo) {
		t.Fatal("damaged external worktree still reported ready")
	}
	if !HasBirthIncarnation(repo) {
		t.Fatal("external worktree damage erased local incarnation evidence")
	}
}
