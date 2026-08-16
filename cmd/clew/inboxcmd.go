package main

import (
	"fmt"
	"os"
	"time"

	"clew/internal/docket"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/state"
)

// cmdDocket renders the bounded decision-only surface (§8.2). The old inbox
// command is a hidden alias at dispatch; it does not have separate semantics.
func cmdDocket(args []string) error {
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
		return docketList(a, repo)
	}
	switch args[0] {
	case "answer":
		if len(args) < 3 {
			return fmt.Errorf("usage: clew docket answer <question-id> \"text\"")
		}
		return journalAnswer(a, repo, args[1], args[2], "")
	case "ack", "drop":
		if len(args) < 2 {
			return fmt.Errorf("usage: clew docket %s <alert-key>", args[0])
		}
		return docketDismiss(a, repo, args[1], args[0])
	default:
		return fmt.Errorf("unknown docket verb %q (answer|ack|drop)", args[0])
	}
}

// cmdInbox remains for source compatibility with older internal callers.
func cmdInbox(args []string) error { return cmdDocket(args) }

func docketList(a *app, repo string) error {
	now := time.Now()
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	lastRuling, learned := docketMetrics(j, now)
	alerts := a.db.OpenAlerts(repo, true)
	assumptions := docketAssumptions(alerts)
	cards := docket.Build(docket.Input{Journal: j, Alerts: alerts, Now: now, Assumptions: assumptions})

	failureKey := "system-failure:docket:" + repoBase(repo)
	if len(cards) > docket.MaxCards {
		if err := a.db.Set(failureKey, fmt.Sprintf("%d active decision cards; cap=%d", len(cards), docket.MaxCards)); err != nil {
			return err
		}
	} else if err := a.db.Set(failureKey, ""); err != nil {
		return err
	}
	precision := currentPushPrecision(alerts)
	return docket.Render(os.Stdout, docket.View{
		Repo: repoBase(repo), Cards: cards, Now: now,
		Empty: docket.EmptyMetricsAt(now, lastRuling, learned), PushPrecision: precision,
	})
}

func docketAssumptions(alerts []state.Alert) map[string]string {
	assumptions := make(map[string]string)
	for _, alert := range alerts {
		switch alert.Kind {
		case "contradiction":
			assumptions["alert:"+alert.Key] = "the retained decision matches the owner's current intent"
		case "stomp":
			assumptions["alert:"+alert.Key] = "the selected owner has the work that should survive"
		case "import":
			assumptions["alert:"+alert.Key] = "the foreign provenance and batch diff are trustworthy"
		}
	}
	return assumptions
}

func currentPushPrecision(alerts []state.Alert) *docket.PushPrecision {
	for _, alert := range alerts {
		if alert.PushedAt != "" {
			return nil // delivery exists but no needed/unneeded adjudication is stored yet
		}
	}
	return &docket.PushPrecision{} // zero delivered implies exactly 0/0
}

func docketMetrics(j *journal.Journal, now time.Time) (time.Time, int) {
	var last time.Time
	for _, event := range j.Events {
		if event.By.Who == "human" && event.At.After(last) {
			last = event.At
		}
	}
	learned := 0
	for _, entry := range j.Entries {
		created := entry.Created()
		if last.IsZero() || (created.After(last) && !created.After(now)) {
			learned++
		}
	}
	return last, learned
}

// docketDismiss acks/drops locally AND records a disposition event on the
// journal so the ruling propagates to every other machine's docket.
func docketDismiss(a *app, repo, key, verb string) error {
	entry := "alert:" + key
	for _, al := range a.db.OpenAlerts(repo, false) {
		if al.Key == key && al.EntryIDs != "" {
			entry = al.EntryIDs
			break
		}
	}
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
