package main

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"time"

	"clew/internal/docket"
	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/model"
	"clew/internal/proposal"
	"clew/internal/push"
	"clew/internal/state"
)

func cmdImport(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: clew import <bundle.yaml|dir|https-url>")
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	manager := proposal.Default()
	batch, err := manager.Stage(context.Background(), repo, args[0])
	if err != nil {
		return err
	}
	if batch.Status != "pending" {
		return fmt.Errorf("proposal %s is already %s", batch.ID, batch.Status)
	}
	if branch := os.Getenv("CLEW_PROPOSAL_BRANCH"); branch != "" {
		if err := manager.PushBranch(repo, batch, branch); err != nil {
			return err
		}
		fmt.Printf("proposal %s pushed to %s; open a PR with base %s (merge = confirm)\n", batch.ID, branch, gitx.Branch)
		return nil
	}
	alert := proposalAlert(repo, batch)
	created, err := a.db.UpsertAlert(alert)
	if err != nil {
		return err
	}
	if created {
		if err := pushProposalCard(a, repo, alert, batch); err != nil {
			a.db.Set("push-error:"+repo, err.Error())
		} else {
			a.db.Set("push-error:"+repo, "")
		}
	}
	path, _ := manager.DiffPath(repo, batch.ID)
	fmt.Printf("proposal %s · entries=%d · diff=%s\n", batch.ID, len(batch.Entries), path)
	fmt.Printf("review: clew docket open %s · accept|reject %s\n", batch.ID, batch.ID)
	return nil
}

func proposalAlert(repo string, batch *proposal.Batch) state.Alert {
	return state.Alert{
		Key: "import:" + batch.ID, RepoPath: repo, Kind: "import",
		Body:     fmt.Sprintf("foreign proposal %s · %d entries", batch.ID, len(batch.Entries)),
		EntryIDs: batch.ID, Blocking: true,
		WithdrawWhen: "proposal:" + batch.ID + ":status!=pending",
	}
}

func proposalEvidence(batch *proposal.Batch) []docket.Evidence {
	n := len(batch.Entries)
	if n > 3 {
		n = 3
	}
	out := make([]docket.Evidence, 0, n)
	for _, entry := range batch.Entries[:n] {
		out = append(out, docket.Evidence{
			Text: entry.Quote, Verbatim: true,
			Provenance: docket.Provenance{Kind: string(entry.Source.Kind), Ref: entry.Source.Ref, EntryID: entry.ID},
		})
	}
	return out
}

func proposalDocketEvidence(repo string, alerts []state.Alert) (map[string][]docket.Evidence, error) {
	manager := proposal.Default()
	out := make(map[string][]docket.Evidence)
	for _, alert := range alerts {
		if alert.Kind != "import" {
			continue
		}
		batch, err := manager.Load(repo, alert.EntryIDs)
		if err != nil {
			return nil, fmt.Errorf("proposal card %s: %w", alert.Key, err)
		}
		if batch.Status != "pending" {
			return nil, fmt.Errorf("proposal card %s points to %s batch", alert.Key, batch.Status)
		}
		out["alert:"+alert.Key] = proposalEvidence(batch)
	}
	return out, nil
}

func pushProposalCard(a *app, repo string, alert state.Alert, batch *proposal.Batch) error {
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	key := "alert:" + alert.Key
	cards := docket.Build(docket.Input{
		Journal: j, Alerts: []state.Alert{alert}, Now: time.Now(),
		Assumptions: docketAssumptions([]state.Alert{alert}),
		Evidence:    map[string][]docket.Evidence{key: proposalEvidence(batch)},
	})
	card, ok := cardCreatedByAlert(cards, alert)
	if !ok {
		return fmt.Errorf("proposal %s did not produce a docket card", batch.ID)
	}
	title, body := docketPushMessage(card)
	sent, err := push.Send(a.cfg.Push, title, body)
	if err != nil {
		return err
	}
	if sent {
		return a.db.MarkAlert(alert.Key, "pushed_at")
	}
	return nil
}

func docketProposalAction(a *app, repo, verb, batchID string) error {
	manager := proposal.Default()
	batch, err := manager.Load(repo, batchID)
	if err != nil {
		return err
	}
	switch verb {
	case "open":
		path, err := manager.DiffPath(repo, batchID)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		fmt.Print(string(b))
		return nil
	case "reject":
		if err := manager.SetStatus(repo, batchID, "rejected"); err != nil {
			return err
		}
		if err := a.db.MarkAlert("import:"+batchID, "dropped_at"); err != nil {
			return err
		}
		fmt.Printf("rejected %s\n", batchID)
		return nil
	case "accept":
		return acceptProposal(a, repo, manager, batch)
	default:
		return fmt.Errorf("unknown proposal verb %q", verb)
	}
}

func acceptProposal(a *app, repo string, manager *proposal.Manager, batch *proposal.Batch) error {
	if batch.Status != "pending" && batch.Status != "accepted" {
		return fmt.Errorf("proposal %s is %s", batch.ID, batch.Status)
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	for _, entry := range batch.Entries {
		if existing := j.Entries[entry.ID]; existing != nil && !reflect.DeepEqual(existing, entry) {
			return fmt.Errorf("entry %s already exists with different content", entry.ID)
		}
	}
	now := time.Now().UTC()
	added := 0
	for _, entry := range batch.Entries {
		if j.Entries[entry.ID] == nil {
			copy := *entry
			if err := j.AddEntry(&copy); err != nil {
				return err
			}
			added++
		}
		if !j.HasEvent(model.EvConfirm, entry.ID, "proposal", batch.ID) {
			if err := j.AddEvent(&model.Event{
				ID: ids.NewEvent(now), Kind: model.EvConfirm, Entry: entry.ID,
				Payload: map[string]any{"proposal": batch.ID, "foreign": true},
				By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: now,
			}); err != nil {
				return err
			}
		}
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	if err := manager.SetStatus(repo, batch.ID, "accepted"); err != nil {
		return err
	}
	if err := a.db.MarkAlert("import:"+batch.ID, "acked_at"); err != nil {
		return err
	}
	fmt.Printf("accepted %s · entries=%d added=%d · human-confirmed\n", batch.ID, len(batch.Entries), added)
	return nil
}

func proposalBatchID(arg string) string { return strings.TrimPrefix(arg, "import:") }
