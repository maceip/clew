package main

import (
	"fmt"
	"time"

	"restart/internal/ids"
	"restart/internal/journal"
	"restart/internal/model"
)

// cmdInbox lists human-blocking items only (§8.2): ★questions, possible-
// contradiction pairs, absences, stomps, adapter breaks. Verbs: answer <id>
// "text" · ack <key> · drop <key>. Radar, never locks (I3).
func cmdInbox(args []string) error {
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return inboxList(a, repo)
	}
	switch args[0] {
	case "answer":
		if len(args) < 3 {
			return fmt.Errorf("usage: restart inbox answer <question-id> \"text\"")
		}
		return journalAnswer(a, repo, args[1], args[2], "")
	case "ack", "drop":
		if len(args) < 2 {
			return fmt.Errorf("usage: restart inbox %s <alert-key>", args[0])
		}
		return inboxDismiss(a, repo, args[1], args[0])
	default:
		return fmt.Errorf("unknown inbox verb %q (answer|ack|drop)", args[0])
	}
}

func inboxList(a *app, repo string) error {
	now := time.Now()
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	st := journal.Compute(j, now)

	n := 0
	fmt.Printf("INBOX — %s (human-blocking only; every item names why it blocks)\n\n", repoBase(repo))
	// ★ questions addressed to the human.
	for id, e := range j.Entries {
		if e.Type != model.Question || e.Asks != "human" || st[id].Status != journal.StOpen {
			continue
		}
		age := int(now.Sub(e.Created()).Hours() / 24)
		fmt.Printf("★ %s %q — open %dd; blocks: agents cannot answer it\n    ↳ restart inbox answer %s \"…\"\n", id, e.Title, age, id)
		n++
	}
	for _, al := range a.db.OpenAlerts(repo, true) {
		fmt.Printf("● [%s] %s\n    ↳ restart inbox ack %s   (or: drop)\n", al.Kind, al.Body, al.Key)
		n++
	}
	if n == 0 {
		fmt.Println("(empty — nothing needs you)")
	}
	return nil
}

// inboxDismiss acks/drops locally AND records a disposition event on the
// journal so the ack propagates to every other machine's inbox.
func inboxDismiss(a *app, repo, key, verb string) error {
	col := "acked_at"
	if verb == "drop" {
		col = "dropped_at"
	}
	if err := a.db.MarkAlert(key, col); err != nil {
		return err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	entry := "alert:" + key
	for _, al := range a.db.OpenAlerts(repo, false) {
		if al.Key == key && al.EntryIDs != "" {
			entry = al.EntryIDs
			break
		}
	}
	j.AddEvent(&model.Event{
		ID: ids.NewEvent(time.Now()), Kind: model.EvDisposition, Entry: entry,
		Payload: map[string]any{"ack": key, "verb": verb},
		By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: time.Now().UTC(),
	})
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("%sed %s\n", verb, key)
	return nil
}
