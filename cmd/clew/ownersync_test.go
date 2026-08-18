package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/config"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/owner"
	"clew/internal/state"
)

func TestOwnerLawSyncNotesDistinguishColdFromPreviouslyVerifiedCache(t *testing.T) {
	t.Run("cold cache refuses offline birth data", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("CLEW_HOME", home)
		remote := filepath.Join(home, "owner-remote.git")
		if out, err := exec.Command("git", "init", "--bare", "-q", remote).CombinedOutput(); err != nil {
			t.Fatalf("git init --bare: %v: %s", err, out)
		}
		cfg := config.Default()
		cfg.Owner.Remote = remote
		db, err := state.Open(filepath.Join(home, "state.db"))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		a := &app{cfg: cfg, db: db}
		// Ensure bootstraps the remote branch but deliberately does not certify
		// that this cache has completed a successful owner-law sync.
		if _, err := owner.Default(remote).Ensure(); err != nil {
			t.Fatal(err)
		}
		if err := os.Rename(remote, remote+".offline"); err != nil {
			t.Fatal(err)
		}
		if _, err := a.ownerLawBlock(true, time.Now()); err == nil {
			t.Fatal("cold owner sync unexpectedly accepted a fetch-failure note")
		}
		if _, err := a.ownerLawBlock(false, time.Now()); err == nil || !strings.Contains(err.Error(), "has not been verified") {
			t.Fatalf("cold cached owner load error = %v", err)
		}
		if db.Get("owner-sync-error") == "" {
			t.Fatal("cold owner sync failure was not made loud")
		}
	})

	t.Run("verified cache survives offline sync loudly", func(t *testing.T) {
		home := t.TempDir()
		t.Setenv("CLEW_HOME", home)
		remote := filepath.Join(home, "owner-remote.git")
		if out, err := exec.Command("git", "init", "--bare", "-q", remote).CombinedOutput(); err != nil {
			t.Fatalf("git init --bare: %v: %s", err, out)
		}
		cfg := config.Default()
		cfg.Owner.Remote = remote
		db, err := state.Open(filepath.Join(home, "state.db"))
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		a := &app{cfg: cfg, db: db}

		source, err := journal.Load(t.TempDir())
		if err != nil {
			t.Fatal(err)
		}
		now := time.Now().UTC()
		finding := &model.Entry{
			ID: ids.NewEntry(now), Type: model.Finding,
			Title: "Verify affected state", Body: "generated summary",
			Quote:       "verify the affected state before declaring completion",
			UtteranceBy: model.ByUser,
			Source:      model.Source{Kind: model.SrcHuman, Ref: "human:owner-law", At: now},
			Confidence:  1,
		}
		if err := source.AddEntry(finding); err != nil {
			t.Fatal(err)
		}
		if _, err := owner.Default(remote).Promote(source, finding.ID, "/project", "desk", now.Add(time.Minute)); err != nil {
			t.Fatal(err)
		}
		laws, err := a.ownerLawBlock(true, now.Add(2*time.Minute))
		if err != nil || !strings.Contains(laws, finding.Quote) {
			t.Fatalf("verified owner laws = %q, %v", laws, err)
		}
		if err := os.Rename(remote, remote+".offline"); err != nil {
			t.Fatal(err)
		}
		if _, err := a.ownerLawBlock(true, now.Add(3*time.Minute)); err == nil {
			t.Fatal("offline refresh unexpectedly reported current owner laws")
		}
		loud := db.Get("owner-sync-error")
		if loud == "" {
			t.Fatal("offline refresh was not visible in owner-sync-error")
		}
		cached, err := a.ownerLawBlock(false, now.Add(3*time.Minute))
		if err != nil || cached != laws {
			t.Fatalf("verified cached laws = %q, %v; want %q", cached, err, laws)
		}
		if got := db.Get("owner-sync-error"); got != loud {
			t.Fatalf("cached fallback cleared loud sync failure: before=%q after=%q", loud, got)
		}
	})
}
