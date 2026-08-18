package checkin

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/state"
)

func TestKnowledgeMergeIsBoundedAmnesiaProofAndHasExactVerbs(t *testing.T) {
	base := time.Date(2026, 8, 18, 17, 0, 0, 0, time.UTC)
	j := &journal.Journal{Entries: map[string]*model.Entry{}}
	for i := 0; i < 8; i++ {
		title := fmt.Sprintf("ordinary remembered change %d", i)
		if i == 7 {
			title = "Codex finished I13 stale: tree uncommitted, law wording on human surfaces"
		}
		e := testEntry(ids.NewEntry(base.Add(time.Duration(i)*time.Minute)), model.Finding, title, base.Add(time.Duration(i)*time.Minute))
		j.Entries[e.ID] = e
	}

	view := BuildMerge(j, nil)
	if len(view.Items) != MaxItems {
		t.Fatalf("merge items = %d, want %d", len(view.Items), MaxItems)
	}
	if got := view.Items[0].Line; !strings.Contains(got, "Rename “law/laws”") || strings.Contains(got, "tree uncommitted") {
		t.Fatalf("law item did not pass the amnesia translation: %q", got)
	}
	var out bytes.Buffer
	view.Repo = "restart"
	if err := Render(&out, view); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"KNOWLEDGE MERGE — it remembers what we decide", "apply · explain · defer", "apply-all"} {
		if !strings.Contains(got, want) {
			t.Fatalf("render missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "build ·") || strings.Contains(got, "retire") {
		t.Fatalf("merge leaked sibling verbs:\n%s", got)
	}
}

func TestIntentGapIncludesHookWiringAndLiveSpendFailureWithoutBuildAll(t *testing.T) {
	base := time.Date(2026, 8, 18, 16, 0, 0, 0, time.UTC)
	hook := testEntry(ids.NewEntry(base), model.Intent, "Build the freshness ladder: one delta payload, five delivery layers", base)
	spend := testEntry(ids.NewEntry(base.Add(time.Minute)), model.Decision, "I9 frugality replaced: listening completeness is the invariant, cost is a dial", base.Add(time.Minute))
	unrelated := testEntry(ids.NewEntry(base.Add(2*time.Minute)), model.Decision, "Owner laws live in an owner-scope journal with a ≤1KB injection budget", base.Add(2*time.Minute))
	j := &journal.Journal{Entries: map[string]*model.Entry{
		hook.ID: hook, spend.ID: spend, unrelated.ID: unrelated,
	}}

	view := BuildGap(j, []state.Alert{{Kind: "budget", Body: "extraction paused"}})
	if len(view.Items) != 2 {
		t.Fatalf("gap items = %#v, want hook and spend only", view.Items)
	}
	lines := view.Items[0].Line + "\n" + view.Items[1].Line
	for _, want := range []string{"spend floor above one full request", "Wire every agent"} {
		if !strings.Contains(lines, want) {
			t.Fatalf("gap missing %q:\n%s", want, lines)
		}
	}
	var out bytes.Buffer
	view.Repo = "restart"
	if err := Render(&out, view); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "build · explain · retire") {
		t.Fatalf("gap verbs missing:\n%s", got)
	}
	if strings.Contains(got, "build-all") || strings.Contains(got, "apply-all") {
		t.Fatalf("gap rendered a forbidden all-action:\n%s", got)
	}
}

func TestEmptyAndBrokenAreNeverTheSameScreen(t *testing.T) {
	for _, tc := range []struct {
		name   string
		issues int
		want   string
		not    string
	}{
		{name: "empty", want: "Nothing new.", not: "No trustworthy empty result"},
		{name: "broken", issues: 1, want: "No trustworthy empty result is available.", not: "Nothing new."},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer
			if err := Render(&out, View{Screen: KnowledgeMerge, Issues: tc.issues}); err != nil {
				t.Fatal(err)
			}
			if got := out.String(); !strings.Contains(got, tc.want) || strings.Contains(got, tc.not) {
				t.Fatalf("wrong state:\n%s", got)
			}
		})
	}
}

func TestExplainHandoffCarriesOnlyEntryIdentityAndHardDirective(t *testing.T) {
	var out bytes.Buffer
	id := "e01M0AZ4BGF1VC0VSXA05VYVEQ3"
	if err := RenderAgentHandoff(&out, "explain", []string{id}); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"CLEW_AGENT_HANDOFF_V1", "ACTION=EXPLAIN", "ENTRY=" + id, "DO NOT CHANGE FILES"} {
		if !strings.Contains(got, want) {
			t.Fatalf("handoff missing %q:\n%s", want, got)
		}
	}
	if strings.Contains(got, "because") || strings.Contains(got, "selected change") {
		t.Fatalf("handoff encoded an explanation instead of handing over the entry:\n%s", got)
	}
}

func testEntry(id string, typ model.EntryType, title string, at time.Time) *model.Entry {
	return &model.Entry{
		ID: id, Type: typ, Title: title, Body: title, Quote: title,
		UtteranceBy: model.ByUser,
		Source:      model.Source{Kind: model.SrcSession, Ref: "session:test", At: at},
		Confidence:  1,
	}
}
