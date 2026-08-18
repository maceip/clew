package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/docket"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/state"
)

func TestGlanceHTMLIsSelfContainedRefreshingAndUsesTitleLight(t *testing.T) {
	t.Setenv("CLEW_HOME", t.TempDir())
	view := glanceView{
		Repo: "project", Generated: "12:00:00", Title: "clew ●", Docket: 2,
		Sections: []glanceSection{
			{Name: "DECIDED"}, {Name: "LEARNED"},
			{Name: "OPEN", Items: []glanceItem{{ID: "e1", Title: "Owner choice?", Age: "2d", Status: "open", Star: true}}},
		},
		Intents: 3, InFlight: 1, Absent: 1, Proposed: 1,
	}
	path, err := writeGlanceHTML(view)
	if err != nil {
		t.Fatal(err)
	}
	if path != filepath.Join(os.Getenv("CLEW_HOME"), "glance.html") {
		t.Fatalf("path = %s", path)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	got := string(b)
	for _, want := range []string{`http-equiv="refresh" content="30"`, `<title>clew ●</title>`, "DECIDED", "DOCKET"} {
		if !strings.Contains(got, want) {
			t.Errorf("HTML missing %q", want)
		}
	}
	for _, forbidden := range []string{"<script", "http://", "https://"} {
		if strings.Contains(got, forbidden) {
			t.Errorf("HTML is not self-contained: found %q", forbidden)
		}
	}
}

func TestBuildGlanceShowsCalmSectionsMapAndDocketCount(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	now := time.Date(2026, 8, 16, 18, 0, 0, 0, time.UTC)
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	add := func(typ model.EntryType, title string) *model.Entry {
		e := &model.Entry{
			ID: ids.NewEntry(now), Type: typ, Title: title, Quote: title,
			UtteranceBy: model.ByUser, Confidence: .9,
			Source: model.Source{Kind: model.SrcSession, Ref: "s#L1", At: now},
		}
		if err := j.AddEntry(e); err != nil {
			t.Fatal(err)
		}
		return e
	}
	add(model.Decision, "ship")
	add(model.Finding, "measured")
	q := add(model.Question, "Owner choice?")
	q.Asks = "human"
	add(model.Intent, "build")
	app := &app{db: db}
	view := buildGlance(app, "/tmp/project", j, now)
	if len(view.Sections) != 3 || view.Intents != 1 || view.Docket != 1 || view.Title != "clew ●" {
		t.Fatalf("view = %#v", view)
	}
}

func TestCardCreatedByAlertNeverSelectsUnrelatedOrFYICard(t *testing.T) {
	alert := state.Alert{Key: "stomp:x", Kind: "stomp"}
	cards := []docket.Card{{Key: "question:e1", Kind: "question"}, {Key: "alert:stomp:x", Kind: "stomp"}}
	got, ok := cardCreatedByAlert(cards, alert)
	if !ok || got.Key != "alert:stomp:x" {
		t.Fatalf("selected = %#v, %v", got, ok)
	}
	if _, ok := cardCreatedByAlert(cards, state.Alert{Key: "suspect:x", Kind: "suspect"}); ok {
		t.Fatal("FYI alert selected a push card")
	}
}

func TestDocketPushMessageIsHeadlinePlusWhyStrip(t *testing.T) {
	card := docket.Card{
		Headline: "Which owner should keep the file?",
		Why: docket.WhyYou{
			RuleCode: "dirty-same-file-overlap", Rule: "running sessions overlap on a dirty file",
			CostOfDelay: "continued edits risk lost work",
		},
	}
	title, body := docketPushMessage(card)
	if title != card.Headline || body != "running sessions overlap on a dirty file · cost of delay: continued edits risk lost work" {
		t.Fatalf("push = %q / %q", title, body)
	}
}
