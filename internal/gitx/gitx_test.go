package gitx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"restart/internal/ids"
	"restart/internal/journal"
	"restart/internal/model"
)

// mkRemote creates a bare "remote" and two clones ("machines" A and B).
func mkRemote(t *testing.T) (bare, repoA, repoB string) {
	t.Helper()
	t.Setenv("RESTART_HOME", t.TempDir())
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

	// Redact on A: scrub in place, rewrite root, force-push.
	b2, _ := os.ReadFile(p)
	os.WriteFile(p, []byte(strings.Replace(string(b2), "AKIAIOSFODNN7EXAMPLE", "‹redacted›", 1)), 0o644)
	if err := RewriteRoot(repoA, "redact "+secretID); err != nil {
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

func TestLocalOnlyNoRemote(t *testing.T) {
	t.Setenv("RESTART_HOME", t.TempDir())
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
}
