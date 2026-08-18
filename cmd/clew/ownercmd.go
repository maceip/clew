package main

import (
	"fmt"
	"strings"
	"time"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/materialize"
	"clew/internal/model"
	"clew/internal/owner"
)

// journalPromote is the sole project-to-owner crossing. The extractor may
// create a proposal card, but only this human command writes the owner-scope
// certification that makes a finding ambient.
func journalPromote(a *app, repo, id string) error {
	source, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	now := time.Now().UTC()
	store := owner.Default(a.cfg.Owner.Remote)
	result, err := store.Promote(source, id, repo, a.cfg.Surface, now)
	if err != nil {
		return err
	}

	// Record the scope crossing in the source journal as well. This disposition
	// carries the alert acknowledgement to other machines and makes retries
	// idempotent if either journal sync was interrupted.
	if !sourcePromotionRecorded(source, id) {
		if err := source.AddEvent(&model.Event{
			ID: ids.NewEvent(now), Kind: model.EvDisposition, Entry: id,
			Payload: map[string]any{
				"action":      owner.ActionPromote,
				"scope":       owner.ScopeOwner,
				"owner_entry": id,
				"ack":         "promotion:" + id,
			},
			By: model.By{Who: "human", Surface: a.cfg.Surface}, At: now,
		}); err != nil {
			return err
		}
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	if err := a.db.MarkAlert("promotion:"+id, "acked_at"); err != nil {
		return err
	}

	// Existing projects receive the newly certified ambient layer now; future
	// projects receive the same block during their laws-only birth.
	laws, err := a.ownerLawBlock(false, now)
	if err != nil {
		return err
	}
	if err := fanoutOwnerLaws(a, laws, now); err != nil {
		return err
	}
	action := "certified"
	if !result.CertificationAdded {
		action = "already certified"
	}
	fmt.Printf("%s: %s as an ambient owner law (%d/%d bytes)\n", id, action, result.Render.RequiredBytes, owner.LawCap)
	return nil
}

func sourcePromotionRecorded(j *journal.Journal, id string) bool {
	if j == nil {
		return false
	}
	for _, event := range j.EventsFor(id) {
		if event.Kind == model.EvDisposition && event.By.Who == "human" &&
			event.PStr("scope") == owner.ScopeOwner && event.PStr("action") == owner.ActionPromote {
			return true
		}
	}
	return false
}

func fanoutOwnerLaws(a *app, laws string, now time.Time) error {
	repos, err := a.db.Repos()
	if err != nil {
		return err
	}
	var failures []string
	for _, registered := range repos {
		j, err := a.openJournal(registered.Path)
		if err == nil {
			err = j.Reload()
		}
		if err == nil {
			err = materialize.WriteWithOwner(registered.Path, j, journal.Compute(j, now), a.db, laws, now)
		}
		if err != nil {
			_ = a.db.Set("materialize-error:"+registered.Path, err.Error())
			failures = append(failures, repoBase(registered.Path)+": "+err.Error())
			continue
		}
		_ = a.db.Set("materialize-error:"+registered.Path, "")
	}
	if len(failures) > 0 {
		return fmt.Errorf("owner law was certified, but context fanout failed: %s", strings.Join(failures, "; "))
	}
	return nil
}
