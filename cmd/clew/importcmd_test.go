package main

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/config"
	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/proposal"
	"clew/internal/state"
)

func TestAcceptProposalIsHumanConfirmBoundary(t *testing.T) {
	t.Setenv("CLEW_HOME", filepath.Join(t.TempDir(), "clew-home"))
	repo := filepath.Join(t.TempDir(), "repo")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	db, err := state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.RegisterRepo(repo, ""); err != nil {
		t.Fatal(err)
	}
	a := &app{cfg: config.Default(), db: db}
	at := time.Now().UTC()
	entry := &model.Entry{
		ID: ids.NewEntry(at), Type: model.Decision, Title: "foreign choice",
		Quote: "choose the guarded path", UtteranceBy: model.ByUser, Confidence: .8,
		Source: model.Source{Kind: model.SrcHuman, Ref: "proposal.md:3", At: at},
	}
	bundle, _ := yaml.Marshal(proposal.Bundle{Version: 1, Entries: []*model.Entry{entry}})
	source := filepath.Join(t.TempDir(), "bundle.yaml")
	if err := os.WriteFile(source, bundle, 0o600); err != nil {
		t.Fatal(err)
	}
	manager := proposal.Default()
	batch, err := manager.Stage(context.Background(), repo, source)
	if err != nil {
		t.Fatal(err)
	}
	alert := proposalAlert(repo, batch)
	if _, err := db.UpsertAlert(alert); err != nil {
		t.Fatal(err)
	}
	if err := acceptProposal(a, repo, manager, batch); err != nil {
		t.Fatal(err)
	}
	j, err := journal.Load(gitx.WorktreeDir(repo))
	if err != nil {
		t.Fatal(err)
	}
	imported := j.Entries[entry.ID]
	if imported == nil || imported.Source.Kind != model.SrcForeign || imported.Quote != entry.Quote {
		t.Fatalf("imported = %#v", imported)
	}
	if !j.HasEvent(model.EvConfirm, entry.ID, "proposal", batch.ID) {
		t.Fatal("accepted foreign entry lacks human confirm event")
	}
	loaded, err := manager.Load(repo, batch.ID)
	if err != nil || loaded.Status != "accepted" {
		t.Fatalf("batch = %#v, %v", loaded, err)
	}
	if got := db.OpenAlerts(repo, true); len(got) != 0 {
		t.Fatalf("accepted card remained open: %#v", got)
	}
	// Idempotent retry cannot duplicate the immutable entry or confirm event.
	if err := acceptProposal(a, repo, manager, loaded); err != nil {
		t.Fatal(err)
	}
	j, _ = journal.Load(gitx.WorktreeDir(repo))
	if got := len(j.EventsFor(entry.ID)); got != 1 {
		t.Fatalf("events after retry = %d", got)
	}
}
