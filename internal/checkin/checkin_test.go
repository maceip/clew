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
	if got := view.Items[0].Line; !strings.Contains(got, "Human-facing words stay calm everywhere") || strings.Contains(got, "tree uncommitted") {
		t.Fatalf("wording item did not pass the amnesia translation: %q", got)
	}
	var out bytes.Buffer
	view.Repo = "restart"
	if err := Render(&out, view); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	for _, want := range []string{"KNOWLEDGE MERGE — restart remembers what we decide", "apply · explain · defer", "apply all"} {
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
	if len(view.Items) != 1 {
		t.Fatalf("gap items = %#v, want only the unfinished hook wiring", view.Items)
	}
	lines := view.Items[0].Line
	if !strings.Contains(lines, "Every agent needs new decisions") {
		t.Fatalf("gap missing hook wiring:\n%s", lines)
	}
	if strings.Contains(lines, "spend floor") {
		t.Fatalf("finished listener work remained in the gap:\n%s", lines)
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
		{name: "broken", issues: 1, want: "The attending agent must repair saved knowledge", not: "Nothing new."},
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

func TestVerifiedWorkSettlesAndOnlyShowsOnce(t *testing.T) {
	base := time.Date(2026, 8, 18, 17, 0, 0, 0, time.UTC)
	done := testEntry(ids.NewEntry(base), model.Decision, "Name the system restart", base)
	pending := testEntry(ids.NewEntry(base.Add(time.Minute)), model.Decision, "Choose the next color", base.Add(time.Minute))
	j := &journal.Journal{Entries: map[string]*model.Entry{done.ID: done, pending.ID: pending}}
	evidence := &model.Event{
		ID: ids.NewEvent(base.Add(2 * time.Minute)), Kind: model.EvEvidence, Entry: done.ID,
		Payload: map[string]any{"kind": "commit", "ref": "abc123"}, By: model.By{Who: "differ"}, At: base.Add(2 * time.Minute),
	}
	j.Events = []*model.Event{evidence}

	first := BuildMerge(j, nil)
	if len(first.Settled) != 1 || len(first.Items) != 1 || EntryIDs(first.Settled[0])[0] != done.ID {
		t.Fatalf("first merge = settled %#v actionable %#v", first.Settled, first.Items)
	}
	var out bytes.Buffer
	first.Repo = "restart"
	if err := Render(&out, first); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); !strings.Contains(got, "Settled while you were away") || strings.Contains(got, "The system is named restart — apply") {
		t.Fatalf("settled work was actionable or silent:\n%s", got)
	}
	second := BuildMerge(j, map[string]string{done.ID: "settled"})
	if len(second.Settled) != 0 || len(second.Items) != 1 {
		t.Fatalf("settled receipt repeated: %#v", second)
	}
}

func TestJournalCommitDoesNotSettleWork(t *testing.T) {
	base := time.Date(2026, 8, 18, 17, 0, 0, 0, time.UTC)
	entry := testEntry(ids.NewEntry(base), model.Decision, "Choose the next color", base)
	j := &journal.Journal{Entries: map[string]*model.Entry{entry.ID: entry}, Events: []*model.Event{{
		ID: ids.NewEvent(base.Add(time.Minute)), Kind: model.EvEvidence, Entry: entry.ID,
		Payload: map[string]any{"kind": "commit", "ref": "journal-sha", "note": "journal: decision entry"},
		By:      model.By{Who: "differ"}, At: base.Add(time.Minute),
	}}}
	view := BuildMerge(j, nil)
	if len(view.Settled) != 0 || len(view.Items) != 1 {
		t.Fatalf("journal-only commit settled work: %#v", view)
	}
}

func TestMemoryLagAppearsOnBothScreens(t *testing.T) {
	for _, screen := range []Screen{KnowledgeMerge, IntentGap} {
		var out bytes.Buffer
		if err := Render(&out, View{Screen: screen, Repo: "restart", MemoryLag: 91 * time.Second}); err != nil {
			t.Fatal(err)
		}
		if got := out.String(); !strings.Contains(got, "memory is 2 minutes behind") || strings.Contains(got, "paused") {
			t.Fatalf("lag wording for %s:\n%s", screen, got)
		}
	}
}

func TestPresentationFoldsCloudWorkAndHidesHeldWork(t *testing.T) {
	base := time.Date(2026, 8, 18, 17, 0, 0, 0, time.UTC)
	pr := testEntry(ids.NewEntry(base), model.Intent, "Surface coverage: PR-only cloud agents (Codex-app-class) contribute knowledge", base)
	write := testEntry(ids.NewEntry(base.Add(time.Minute)), model.Intent, "Surface coverage: repo-write cloud agents (Cursor-class) are full journal nodes", base.Add(time.Minute))
	held := testEntry(ids.NewEntry(base.Add(2*time.Minute)), model.Intent, "Held: a restart tab — stage selected drift into the next generation", base.Add(2*time.Minute))
	held.Body = "Owner direction, held for more thinking; not buildable spec yet."
	j := &journal.Journal{Entries: map[string]*model.Entry{pr.ID: pr, write.ID: write, held.ID: held}}

	view := BuildGap(j, nil)
	if len(view.Items) != 1 {
		t.Fatalf("gap items = %#v, want one folded cloud line and no held work", view.Items)
	}
	if got := len(EntryIDs(view.Items[0])); got != 2 {
		t.Fatalf("folded cloud sources = %d, want 2", got)
	}
	var out bytes.Buffer
	view.Repo = "restart"
	if err := Render(&out, view); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if strings.Count(got, "Cloud agents return what they learned") != 1 || strings.Contains(got, "restart tab") {
		t.Fatalf("fold or held filtering failed:\n%s", got)
	}
	if strings.Contains(got, pr.ID) || strings.Contains(got, write.ID) {
		t.Fatalf("human line leaked machine identity:\n%s", got)
	}
}

func TestSpokenWordsResolveWithoutRelayingIdentity(t *testing.T) {
	base := time.Date(2026, 8, 18, 17, 0, 0, 0, time.UTC)
	rename := testEntry(ids.NewEntry(base), model.Finding, "Codex finished I13 stale: tree uncommitted, law wording on human surfaces", base)
	rename.Body = "README and other human surfaces need the wording rename."
	other := testEntry(ids.NewEntry(base.Add(time.Minute)), model.Intent, "Build the freshness ladder: one delta payload, five delivery layers", base.Add(time.Minute))
	view := BuildMerge(&journal.Journal{Entries: map[string]*model.Entry{rename.ID: rename, other.ID: other}}, nil)

	item, err := Resolve(view, "apply the readme rename")
	if err != nil {
		t.Fatal(err)
	}
	if ids := EntryIDs(item); len(ids) != 1 || ids[0] != rename.ID {
		t.Fatalf("spoken words resolved to %#v", ids)
	}
}

func TestBrokenLiveCheckCarriesItsFix(t *testing.T) {
	var out bytes.Buffer
	view := View{Screen: KnowledgeMerge, Repo: "restart", Repairs: []string{"Listening is paused until the spend floor is built — build"}}
	if err := Render(&out, view); err != nil {
		t.Fatal(err)
	}
	got := out.String()
	if !strings.Contains(got, "spend floor is built — build") || strings.Contains(got, "needs attention") {
		t.Fatalf("live check is not actionable:\n%s", got)
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
