package seed

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
)

var testNow = time.Date(2026, 8, 18, 15, 0, 0, 0, time.UTC)

func addEntry(t *testing.T, j *journal.Journal, typ model.EntryType, title string, at time.Time) *model.Entry {
	t.Helper()
	e := &model.Entry{
		ID: ids.NewEntry(at), Type: typ, Title: title, Body: "body: " + title,
		Quote: "quote: " + title, UtteranceBy: model.ByUser,
		Source: model.Source{
			Kind: model.SrcSession, Ref: "codex:/sessions/source.jsonl#L42",
			Agent: "codex", Surface: "laptop", At: at,
		},
		Confidence: .9,
	}
	if typ == model.Question {
		e.Asks = "any"
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	return e
}

func addEvent(t *testing.T, j *journal.Journal, entry string, kind model.EventKind, at time.Time, payload map[string]any) *model.Event {
	t.Helper()
	e := &model.Event{
		ID: ids.NewEvent(at), Kind: kind, Entry: entry, Payload: payload,
		By: model.By{Who: "human", Surface: "laptop"}, At: at,
	}
	if err := j.AddEvent(e); err != nil {
		t.Fatal(err)
	}
	return e
}

func fixture(t *testing.T) (*Snapshot, *model.Entry, *model.Entry, *model.Entry, *model.Event) {
	t.Helper()
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	decision := addEntry(t, j, model.Decision, "keep the narrow substrate", testNow.Add(-4*time.Hour))
	finding := addEntry(t, j, model.Finding, "warm sessions are the latency gate", testNow.Add(-3*time.Hour))
	open := addEntry(t, j, model.Question, "which vendor remains", testNow.Add(-2*time.Hour))
	grave := addEntry(t, j, model.Intent, "discard the relay", testNow.Add(-time.Hour))
	reject := addEvent(t, j, grave.ID, model.EvReject, testNow, map[string]any{"reason": "wrong substrate"})
	exhibit := addEvent(t, j, finding.ID, model.EvEvidence, testNow.Add(time.Minute), map[string]any{"kind": "benchmark", "ref": "commit:abc"})
	st := journal.Compute(j, testNow.Add(2*time.Minute))
	s, err := Build(j, st, BuildInput{
		Repository: Repository{ID: "repo:substrate", Name: "substrate", Path: "/work/substrate", Remote: "git@example/substrate.git"},
		Lifecycle:  Lifecycle{State: "tombstoned", At: time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)},
		Topics:     []string{"Agent substrate", "warm-session", "agent"},
		Ancestors:  []string{"repo:ancestor", "repo:ancestor"},
		OrganBank:  &OrganBankPin{Remote: "git@example/substrate.git", Commit: "0123456789abcdef0123456789abcdef01234567", Dirty: true, At: testNow},
	})
	if err != nil {
		t.Fatal(err)
	}
	if open == nil || decision == nil || reject == nil {
		t.Fatal("fixture setup failed")
	}
	return s, decision, finding, grave, exhibit
}

func TestBuildSelectsLessonsGraveyardExhibitsAndPin(t *testing.T) {
	s, decision, finding, grave, exhibit := fixture(t)
	if len(s.Decisions) != 1 || s.Decisions[0].Entry.ID != decision.ID {
		t.Fatalf("decisions = %#v", s.Decisions)
	}
	if len(s.Findings) != 1 || s.Findings[0].Entry.ID != finding.ID {
		t.Fatalf("findings = %#v", s.Findings)
	}
	if len(s.Graveyard) != 1 || s.Graveyard[0].Entry.ID != grave.ID || len(s.Graveyard[0].Events) != 1 {
		t.Fatalf("graveyard = %#v", s.Graveyard)
	}
	if len(s.Exhibits) != 1 || s.Exhibits[0].ID != exhibit.ID {
		t.Fatalf("exhibits = %#v", s.Exhibits)
	}
	if s.OrganBank == nil || s.OrganBank.Commit != "0123456789abcdef0123456789abcdef01234567" {
		t.Fatalf("organ pin = %#v", s.OrganBank)
	}
	if got := strings.Join(s.Topics, ","); got != "agent,session,substrate,warm" {
		t.Fatalf("normalized topics = %q", got)
	}
	if len(s.Ancestors) != 1 || s.Ancestors[0] != "repo:ancestor" {
		t.Fatalf("ancestors = %#v", s.Ancestors)
	}
	// Active questions/intents are deliberately not part of the ambient lore
	// seed. A manifest may carry them during a deliberate big restart.
	for _, section := range [][]Lesson{s.Decisions, s.Findings} {
		for _, lesson := range section {
			if lesson.Entry.Type == model.Question || lesson.Entry.Type == model.Intent {
				t.Fatalf("active project work leaked into ambient seed: %#v", lesson.Entry)
			}
		}
	}
}

func TestRenderParseIsDeterministicAndPreservesExactProvenance(t *testing.T) {
	s, decision, _, _, _ := fixture(t)
	one, err := Render(s)
	if err != nil {
		t.Fatal(err)
	}
	two, err := Render(s)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(one, two) {
		t.Fatal("same snapshot rendered different bytes")
	}
	for _, want := range []string{"# Project seed — substrate", "## Decisions", "## Findings", "## Graveyard", "## Exhibits", "## Organ-bank pin", "tombstoned 2026-07-14"} {
		if !bytes.Contains(one, []byte(want)) {
			t.Errorf("human seed missing %q", want)
		}
	}
	got, err := Parse(one)
	if err != nil {
		t.Fatal(err)
	}
	entry := got.Decisions[0].Entry
	if entry.ID != decision.ID || entry.Source.Kind != decision.Source.Kind ||
		entry.Source.Ref != decision.Source.Ref || entry.Source.Agent != decision.Source.Agent ||
		entry.Source.Surface != decision.Source.Surface || !entry.Source.At.Equal(decision.Source.At) {
		t.Fatalf("provenance changed across seed round-trip:\n got %#v\nwant %#v", entry.Source, decision.Source)
	}
}

func TestParseRejectsTamperedMachinePayload(t *testing.T) {
	s, _, _, _, _ := fixture(t)
	b, err := Render(s)
	if err != nil {
		t.Fatal(err)
	}
	b = bytes.Replace(b, []byte("keep the narrow substrate"), []byte("poison the new project"), 1)
	if _, err := Parse(b); err == nil || !strings.Contains(err.Error(), "digest mismatch") {
		t.Fatalf("tampered seed error = %v", err)
	}
}

func TestWriteSkipsUnchangedBytesAndAtomicallyReplacesChange(t *testing.T) {
	s, _, _, _, _ := fixture(t)
	path := filepath.Join(t.TempDir(), ".clew", "SEED.md")
	changed, err := Write(path, s)
	if err != nil || !changed {
		t.Fatalf("first write = %v, %v", changed, err)
	}
	old := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(path, old, old); err != nil {
		t.Fatal(err)
	}
	changed, err = Write(path, s)
	if err != nil || changed {
		t.Fatalf("unchanged write = %v, %v", changed, err)
	}
	info, _ := os.Stat(path)
	if !info.ModTime().Equal(old) {
		t.Fatalf("unchanged seed was rewritten: modtime=%s", info.ModTime())
	}
	s.Findings[0].Entry.Body = "new journal-derived finding body"
	changed, err = Write(path, s)
	if err != nil || !changed {
		t.Fatalf("changed write = %v, %v", changed, err)
	}
	read, err := Read(path)
	if err != nil {
		t.Fatal(err)
	}
	if read.Findings[0].Entry.Body != "new journal-derived finding body" {
		t.Fatal("replacement did not contain new seed")
	}
}

func TestValidationRejectsDanglingExhibitBeforeCarry(t *testing.T) {
	s, _, _, _, _ := fixture(t)
	s.Exhibits[0].Entry = "e01M00000000000000000000000"
	if err := Validate(s); err == nil || !strings.Contains(err.Error(), "outside seed") {
		t.Fatalf("dangling exhibit error = %v", err)
	}
}

func TestDroppedDispositionMovesOtherwiseLiveLessonToGraveyard(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	decision := addEntry(t, j, model.Decision, "deliberately left behind", testNow)
	addEvent(t, j, decision.ID, model.EvDisposition, testNow.Add(time.Minute), map[string]any{"disposition": "dropped"})
	st := journal.Compute(j, testNow.Add(2*time.Minute))
	if st[decision.ID].Status != journal.StActive {
		t.Fatalf("fixture relies on disposition not changing algebra status, got %s", st[decision.ID].Status)
	}
	s, err := Build(j, st, BuildInput{
		Repository: Repository{ID: "repo:dropped", Name: "dropped"},
		Lifecycle:  Lifecycle{State: "active"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(s.Decisions) != 0 || len(s.Graveyard) != 1 || s.Graveyard[0].Entry.ID != decision.ID {
		t.Fatalf("dropped disposition classification: decisions=%#v graveyard=%#v", s.Decisions, s.Graveyard)
	}
	if s.Graveyard[0].Status != journal.StSuperseded {
		t.Fatalf("decision disposition was not normalized to a type-compatible terminal status: %s", s.Graveyard[0].Status)
	}
}

func TestJournalRevisionIncludesExactLineageLinkIdentity(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	addEntry(t, j, model.Finding, "same project journal", testNow)
	st := journal.Compute(j, testNow)
	base := BuildInput{
		Repository: Repository{ID: "repo:successor", Name: "successor"},
		Lifecycle:  Lifecycle{State: "active"}, Ancestors: []string{"repo:predecessor"},
		LineageRevision: []string{"link-a\x1frepo:predecessor\x1fsha256:aaa"},
	}
	first, err := Build(j, st, base)
	if err != nil {
		t.Fatal(err)
	}
	base.LineageRevision = []string{"link-b\x1frepo:predecessor\x1fsha256:bbb"}
	second, err := Build(j, st, base)
	if err != nil {
		t.Fatal(err)
	}
	if first.JournalRevision == second.JournalRevision {
		t.Fatal("different explicit lineage link identity/digest collapsed behind the same ancestor id")
	}
}

func TestJournalGateIgnoresPollOnlyMetadataUntilJournalChanges(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	finding := addEntry(t, j, model.Finding, "stable journal fact", testNow)
	st := journal.Compute(j, testNow)
	base, err := Build(j, st, BuildInput{
		Repository: Repository{ID: "repo:gate", Name: "gate"},
		Lifecycle:  Lifecycle{State: "active"},
		Topics:     []string{"old README topic"},
		OrganBank:  &OrganBankPin{Commit: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"},
	})
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "SEED.md")
	if changed, err := WriteOnJournalChange(path, base); err != nil || !changed {
		t.Fatalf("initial gated write = %v, %v", changed, err)
	}
	polled, err := Build(j, st, BuildInput{
		Repository: Repository{ID: "repo:gate", Name: "gate"},
		Lifecycle:  Lifecycle{State: "active"},
		Topics:     []string{"new README-only topic"},
		OrganBank:  &OrganBankPin{Commit: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if polled.JournalRevision != base.JournalRevision {
		t.Fatal("code/poll-only metadata changed journal revision")
	}
	if changed, err := WriteOnJournalChange(path, polled); err != nil || changed {
		t.Fatalf("poll-only gated write = %v, %v", changed, err)
	}
	kept, _ := Read(path)
	if kept.OrganBank.Commit != base.OrganBank.Commit {
		t.Fatal("poll-only organ pin rewrote ambient seed")
	}

	addEvent(t, j, finding.ID, model.EvEvidence, testNow.Add(time.Minute), map[string]any{"kind": "commit", "ref": "bbbb"})
	st = journal.Compute(j, testNow.Add(time.Minute))
	journalChanged, err := Build(j, st, BuildInput{
		Repository: Repository{ID: "repo:gate", Name: "gate"},
		Lifecycle:  Lifecycle{State: "active"},
		Topics:     []string{"new README-only topic"},
		OrganBank:  &OrganBankPin{Commit: "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"},
	})
	if err != nil {
		t.Fatal(err)
	}
	if journalChanged.JournalRevision == base.JournalRevision {
		t.Fatal("journal event did not change revision")
	}
	if changed, err := WriteOnJournalChange(path, journalChanged); err != nil || !changed {
		t.Fatalf("journal-change gated write = %v, %v", changed, err)
	}
	updated, _ := Read(path)
	if updated.OrganBank.Commit != journalChanged.OrganBank.Commit {
		t.Fatal("journal-triggered write did not capture current organ pin")
	}
}
