package main

import (
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/maceip/clew/internal/config"
	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
	"github.com/maceip/clew/internal/state"
)

func TestKeepLocalTargetsFindingWithoutDisposableAlertAndIsIdempotent(t *testing.T) {
	home := t.TempDir()
	t.Setenv("CLEW_HOME", home)
	repo := filepath.Join(t.TempDir(), "project")
	if out, err := exec.Command("git", "init", "-q", repo).CombinedOutput(); err != nil {
		t.Fatalf("git init: %v: %s", err, out)
	}
	db, err := state.Open(filepath.Join(home, "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	a := &app{cfg: config.Default(), db: db}
	wt, err := gitx.EnsureJournal(repo)
	if err != nil {
		t.Fatal(err)
	}
	j, err := journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	finding := &model.Entry{
		ID: ids.NewEntry(now), Type: model.Finding,
		Title: "Portable candidate", Body: "generated summary",
		Quote:       "verify the affected state before declaring completion",
		UtteranceBy: model.ByUser,
		Source:      model.Source{Kind: model.SrcSession, Ref: "session#L1", At: now},
		Confidence:  .9, PromotionCandidate: true,
	}
	if err := j.AddEntry(finding); err != nil {
		t.Fatal(err)
	}
	key := "promotion:" + finding.ID
	if err := docketDismiss(a, repo, key, "drop"); err != nil {
		t.Fatal(err)
	}
	loaded, err := journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}
	if !promotionRuled(loaded, finding.ID) {
		t.Fatal("keep-local without a state alert did not durably rule the finding")
	}
	count := dispositionCount(loaded, finding.ID, key, "drop")
	if count != 1 {
		t.Fatalf("keep-local dispositions = %d, want 1", count)
	}
	if got := len(loaded.EventsFor("alert:" + key)); got != 0 {
		t.Fatalf("keep-local targeted disposable alert identity in %d event(s)", got)
	}
	if err := docketDismiss(a, repo, key, "drop"); err != nil {
		t.Fatal(err)
	}
	loaded, err = journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}
	if got := dispositionCount(loaded, finding.ID, key, "drop"); got != count {
		t.Fatalf("idempotent keep-local wrote %d dispositions, want %d", got, count)
	}

	// If the disposable local index cannot be updated, the journal ruling has
	// already been written and synced. A later retry can repair local state.
	second := *finding
	second.ID = ids.NewEntry(now.Add(time.Minute))
	second.Source.At = now.Add(time.Minute)
	if err := loaded.AddEntry(&second); err != nil {
		t.Fatal(err)
	}
	secondKey := "promotion:" + second.ID
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	if err := docketDismiss(a, repo, secondKey, "drop"); err == nil {
		t.Fatal("keep-local unexpectedly updated a closed state database")
	}
	loaded, err = journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}
	if got := dispositionCount(loaded, second.ID, secondKey, "drop"); got != 1 {
		t.Fatalf("journal-before-local failure ordering wrote %d dispositions, want 1", got)
	}
}

func dispositionCount(j *journal.Journal, id, key, verb string) int {
	count := 0
	for _, event := range j.EventsFor(id) {
		if event.Kind == model.EvDisposition && event.By.Who == "human" &&
			event.PStr("ack") == key && event.PStr("verb") == verb {
			count++
		}
	}
	return count
}
