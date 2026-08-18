package lineage

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/seed"
)

var now = time.Date(2026, 8, 18, 16, 0, 0, 0, time.UTC)

func entry(typ model.EntryType, title string, at time.Time) model.Entry {
	e := model.Entry{
		ID: ids.NewEntry(at), Type: typ, Title: title, Body: "body of " + title,
		Quote: "quote of " + title, UtteranceBy: model.ByUser,
		Source: model.Source{
			Kind: model.SrcSession, Ref: "claude:/sessions/original.jsonl#L17",
			Agent: "claude-code", Surface: "laptop", At: at,
		},
		Confidence: .91,
	}
	if typ == model.Question {
		e.Asks = "any"
	}
	return e
}

func eventFor(entryID string, kind model.EventKind, at time.Time, payload map[string]any) model.Event {
	return model.Event{
		ID: ids.NewEvent(at), Kind: kind, Entry: entryID, Payload: payload,
		By: model.By{Who: "human", Surface: "laptop"}, At: at,
	}
}

func snapshot(id, name string, changed time.Time, topics ...string) *seed.Snapshot {
	decision := entry(model.Decision, "choose the narrow substrate", changed.Add(-3*time.Minute))
	finding := entry(model.Finding, "warm latency is the real gate", changed.Add(-2*time.Minute))
	grave := entry(model.Intent, "abandon the relay", changed.Add(-time.Minute))
	reject := eventFor(grave.ID, model.EvReject, changed, map[string]any{"reason": "wrong transport"})
	exhibit := eventFor(finding.ID, model.EvEvidence, changed.Add(time.Second), map[string]any{"kind": "benchmark", "ref": "commit:abc"})
	return &seed.Snapshot{
		Repository:      seed.Repository{ID: id, Name: name, Path: "/work/" + name, Remote: "git@example/" + name + ".git"},
		JournalRevision: "sha256:aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
		ChangedAt:       changed.Add(time.Second), Lifecycle: seed.Lifecycle{State: "active"}, Topics: topics,
		Decisions: []seed.Lesson{{Entry: decision, Status: journal.StActive}},
		Findings:  []seed.Lesson{{Entry: finding, Status: journal.StCurrent}},
		Graveyard: []seed.Grave{{Entry: grave, Status: journal.StDropped, Events: []model.Event{reject}}},
		Exhibits:  []model.Event{exhibit},
		OrganBank: &seed.OrganBankPin{Remote: "git@example/" + name + ".git", Commit: "0123456789abcdef0123456789abcdef01234567", At: changed},
	}
}

func TestRankUsesTopicAndRecencyAndRendersRequestedSummary(t *testing.T) {
	topical := snapshot("repo:substrate", "substrate", time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC), "security substrate")
	topical.Lifecycle = seed.Lifecycle{State: "tombstoned", At: time.Date(2026, 7, 14, 0, 0, 0, 0, time.UTC)}
	recent := snapshot("repo:calendar", "calendar", now.Add(-time.Hour), "calendar scheduling")
	self := snapshot("repo:target", "target", now, "security substrate")
	ranked, err := Rank(Target{RepositoryID: "repo:target", Name: "security-substrate", Topics: []string{"security substrate"}}, []Candidate{
		{Snapshot: recent, Location: "/work/calendar"},
		{Snapshot: self, Location: "/work/target"},
		{Snapshot: topical, Location: "/work/substrate"},
	}, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(ranked) != 2 {
		t.Fatalf("ranked candidates = %d, want current repo excluded", len(ranked))
	}
	if ranked[0].Snapshot.Repository.ID != "repo:substrate" {
		t.Fatalf("topic match did not outrank unrelated recency: %#v", ranked)
	}
	if !ranked[0].Blatant {
		t.Fatal("strong two-token topic overlap was not exposed as suggestion-only signal")
	}
	if got, want := ranked[0].Summary(), "substrate · 2 lessons · died Jul 14 · tombstoned"; got != want {
		t.Fatalf("summary = %q, want %q", got, want)
	}
}

func TestRankHasStableNameAndIDTieBreaks(t *testing.T) {
	a := snapshot("repo:a", "alpha", now.Add(-time.Hour), "unrelated")
	b := snapshot("repo:b", "beta", now.Add(-time.Hour), "unrelated")
	ranked, err := Rank(Target{RepositoryID: "repo:target", Name: "target"}, []Candidate{{Snapshot: b}, {Snapshot: a}}, now)
	if err != nil {
		t.Fatal(err)
	}
	if ranked[0].Snapshot.Repository.ID != "repo:a" {
		t.Fatalf("tie order = %s, want repo:a", ranked[0].Snapshot.Repository.ID)
	}
}

func TestImportPreservesProvenanceAndWritesLinkLast(t *testing.T) {
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "old topic")
	result, err := Import(ImportRequest{
		Journal: j, TargetRepository: "repo:new", Snapshot: s,
		By: model.By{Who: "human", Surface: "laptop"}, At: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.NewEntries != 3 || result.NewEvents != 3 || result.Link == nil {
		t.Fatalf("import result = %#v", result)
	}
	want := s.Decisions[0].Entry
	got := j.Entries[want.ID]
	if got == nil || got.Source.Kind != want.Source.Kind || got.Source.Ref != want.Source.Ref ||
		got.Source.Agent != want.Source.Agent || got.Source.Surface != want.Source.Surface || !got.Source.At.Equal(want.Source.At) {
		t.Fatalf("original provenance changed:\n got %#v\nwant %#v", got, want)
	}
	links, err := LoadLinks(j.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 || links[0].From.ID != "repo:old" || links[0].SeedJournalRevision != s.JournalRevision ||
		len(links[0].ImportedEntries) != 3 || len(links[0].ImportedEvents) != 3 {
		t.Fatalf("durable links = %#v", links)
	}
	if got, want := links[0].OrganBank, s.OrganBank; got == nil || want == nil ||
		got.Remote != want.Remote || got.Commit != want.Commit || got.Dirty != want.Dirty || !got.At.Equal(want.At) {
		t.Fatalf("durable organ-bank pin = %#v, want %#v", got, want)
	}
	if _, err := os.Stat(filepath.Join(j.Dir, linkDir, result.Link.ID+".yaml")); err != nil {
		t.Fatal(err)
	}
	st := journal.Compute(j, now.Add(time.Minute))
	if grave := s.Graveyard[0].Entry.ID; st[grave].Status != journal.StDropped {
		t.Fatalf("grave status = %s, want dropped", st[grave].Status)
	}
}

func TestImportSameSeedIsIdempotentAndRejectIsNeverResurrected(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
	req := ImportRequest{Journal: j, TargetRepository: "repo:new", Snapshot: s, By: model.By{Who: "human"}, At: now}
	if _, err := Import(req); err != nil {
		t.Fatal(err)
	}
	carried := s.Decisions[0].Entry.ID
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(now.Add(time.Minute)), Kind: model.EvReject, Entry: carried,
		Payload: map[string]any{"reason": "not for this successor"}, By: model.By{Who: "human"}, At: now.Add(time.Minute),
	}); err != nil {
		t.Fatal(err)
	}
	result, err := Import(req)
	if err != nil {
		t.Fatal(err)
	}
	if !result.AlreadyImported || result.NewEntries != 0 || result.NewEvents != 0 {
		t.Fatalf("repeat result = %#v", result)
	}
	if got := journal.Compute(j, now.Add(2*time.Minute))[carried].Status; got != journal.StSuperseded {
		t.Fatalf("explicit un-carry reject was resurrected: status=%s", got)
	}
	links, _ := LoadLinks(j.Dir)
	if len(links) != 1 {
		t.Fatalf("repeat created %d lineage links", len(links))
	}
}

func TestImportUpdatedSeedIsDeltaOnlyAndStillExplicit(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
	base := ImportRequest{Journal: j, TargetRepository: "repo:new", Snapshot: s, By: model.By{Who: "human"}, At: now}
	if _, err := Import(base); err != nil {
		t.Fatal(err)
	}
	newFinding := entry(model.Finding, "new lesson", now.Add(time.Minute))
	s.Findings = append(s.Findings, seed.Lesson{Entry: newFinding, Status: journal.StCurrent})
	s.ChangedAt = now.Add(time.Minute)
	s.JournalRevision = "sha256:bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	// Merely changing the predecessor seed does nothing to the successor.
	if j.Entries[newFinding.ID] != nil {
		t.Fatal("successor changed without an explicit second import")
	}
	base.At = now.Add(2 * time.Minute)
	result, err := Import(base)
	if err != nil {
		t.Fatal(err)
	}
	if result.NewEntries != 1 || result.NewEvents != 1 || j.Entries[newFinding.ID] == nil {
		t.Fatalf("delta result = %#v", result)
	}
	links, _ := LoadLinks(j.Dir)
	if len(links) != 2 {
		t.Fatalf("updated explicit import links = %d, want 2 seed revisions", len(links))
	}
}

func TestImportIdempotencyUsesRepositoryAndJournalRevisionNotMachineMetadata(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
	req := ImportRequest{Journal: j, TargetRepository: "repo:new", Snapshot: s, By: model.By{Who: "human"}, At: now}
	first, err := Import(req)
	if err != nil {
		t.Fatal(err)
	}
	variant := *s
	variant.Repository = s.Repository
	variant.Repository.Path = "/a/different/clone/path"
	variant.Repository.Remote = "ssh://git@example/old.git"
	variant.Topics = []string{"different", "machine", "sampling"}
	variant.OrganBank = &seed.OrganBankPin{Commit: "ffffffffffffffffffffffffffffffffffffffff", Dirty: true}
	firstDigest, _ := seed.Digest(s)
	variantDigest, _ := seed.Digest(&variant)
	if firstDigest == variantDigest {
		t.Fatal("fixture did not vary complete seed digest")
	}
	req.Snapshot = &variant
	req.At = now.Add(time.Minute)
	second, err := Import(req)
	if err != nil {
		t.Fatal(err)
	}
	if !second.AlreadyImported || second.Link.ID != first.Link.ID {
		t.Fatalf("same predecessor journal revision was treated as new lineage: %#v", second)
	}
	links, _ := LoadLinks(j.Dir)
	if len(links) != 1 {
		t.Fatalf("machine-local seed metadata created %d links, want 1", len(links))
	}
}

func TestImportDurablyPreservesEveryDerivedGraveStatus(t *testing.T) {
	decision := entry(model.Decision, "superseded by an omitted successor", now.Add(-5*time.Hour))
	question := entry(model.Question, "expired before the restart", now.Add(-4*time.Hour))
	absent := entry(model.Intent, "lost amid sibling progress", now.Add(-3*time.Hour))
	dropped := entry(model.Intent, "deliberately dispositioned away", now.Add(-2*time.Hour))
	s := &seed.Snapshot{
		Repository:      seed.Repository{ID: "repo:old", Name: "old"},
		JournalRevision: "sha256:cccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccccc",
		ChangedAt:       now.Add(-time.Hour), Lifecycle: seed.Lifecycle{State: "active"},
		Graveyard: []seed.Grave{
			{Entry: decision, Status: journal.StSuperseded},
			{Entry: question, Status: journal.StExpired},
			{Entry: absent, Status: journal.StAbsent},
			{Entry: dropped, Status: journal.StDropped},
		},
	}
	j, _ := journal.Load(t.TempDir())
	result, err := Import(ImportRequest{
		Journal: j, TargetRepository: "repo:new", Snapshot: s,
		By: model.By{Who: "human", Surface: "test"}, At: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.NewEntries != 4 || result.NewEvents != 4 {
		t.Fatalf("grave carry = %#v, want one durable marker per grave", result)
	}
	computed := journal.Compute(j, now.Add(time.Minute))
	for id, want := range map[string]journal.Status{
		decision.ID: journal.StSuperseded,
		question.ID: journal.StExpired,
		absent.ID:   journal.StAbsent,
		dropped.ID:  journal.StDropped,
	} {
		if got := computed[id].Status; got != want {
			t.Errorf("carried grave %s resurrected as %s, want %s", id, got, want)
		}
		if computed[id].LineageStatusEvent == "" {
			t.Errorf("carried grave %s has no durable lineage marker", id)
		}
	}

	// Absence describes the predecessor's observed history; actual work in the
	// successor may revive it, while the other terminal graves stay terminal.
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(now.Add(2 * time.Minute)), Kind: model.EvEvidence, Entry: absent.ID,
		Payload: map[string]any{"kind": "commit", "ref": "successor:abc"},
		By:      model.By{Who: "differ"}, At: now.Add(2 * time.Minute),
	}); err != nil {
		t.Fatal(err)
	}
	if got := journal.Compute(j, now.Add(3*time.Minute))[absent.ID].Status; got != journal.StInFlight {
		t.Fatalf("new successor evidence did not revive carried absence: %s", got)
	}
}

func TestImportConflictFailsBeforeAnyWrite(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
	conflict := s.Decisions[0].Entry
	conflict.Title = "same id but poisoned content"
	if err := j.AddEntry(&conflict); err != nil {
		t.Fatal(err)
	}
	beforeEntries, beforeEvents := len(j.Entries), len(j.Events)
	_, err := Import(ImportRequest{
		Journal: j, TargetRepository: "repo:new", Snapshot: s,
		By: model.By{Who: "human"}, At: now,
	})
	if err == nil || !strings.Contains(err.Error(), "wrote nothing") {
		t.Fatalf("conflict error = %v", err)
	}
	if len(j.Entries) != beforeEntries || len(j.Events) != beforeEvents {
		t.Fatalf("preflight conflict partially wrote journal: entries %d→%d events %d→%d", beforeEntries, len(j.Entries), beforeEvents, len(j.Events))
	}
	links, loadErr := LoadLinks(j.Dir)
	if loadErr != nil || len(links) != 0 {
		t.Fatalf("preflight conflict wrote link: links=%#v err=%v", links, loadErr)
	}
}

func TestImportAcceptsExistingManifestCarriedKind(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
	existing := s.Decisions[0].Entry
	existing.Source.Kind = model.SrcCarried
	if err := j.AddEntry(&existing); err != nil {
		t.Fatal(err)
	}
	result, err := Import(ImportRequest{
		Journal: j, TargetRepository: "repo:new", Snapshot: s,
		By: model.By{Who: "human"}, At: now,
	})
	if err != nil {
		t.Fatal(err)
	}
	if result.NewEntries != 2 {
		t.Fatalf("manifest-compatible import wrote %d entries, want remaining 2", result.NewEntries)
	}
}

func TestImportRejectsSelfCycleAndNonHumanInvocation(t *testing.T) {
	for _, tc := range []struct {
		name   string
		target string
		who    string
		mutate func(*seed.Snapshot)
		want   string
	}{
		{name: "self", target: "repo:old", who: "human", want: "cannot inherit from itself"},
		{name: "cycle", target: "repo:new", who: "human", mutate: func(s *seed.Snapshot) { s.Ancestors = []string{"repo:new"} }, want: "lineage cycle"},
		{name: "automatic", target: "repo:new", who: "watcher", want: "explicitly declared by a human"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			j, _ := journal.Load(t.TempDir())
			s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
			if tc.mutate != nil {
				tc.mutate(s)
			}
			_, err := Import(ImportRequest{Journal: j, TargetRepository: tc.target, Snapshot: s, By: model.By{Who: tc.who}, At: now})
			if err == nil || !strings.Contains(err.Error(), tc.want) {
				t.Fatalf("error = %v, want %q", err, tc.want)
			}
			if len(j.Entries) != 0 {
				t.Fatal("invalid lineage import mutated journal")
			}
		})
	}
}

func TestImportedAncestorsRemainTransitiveForSuccessorSeed(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	s := snapshot("repo:b", "b", now.Add(-time.Hour), "topic")
	s.Ancestors = []string{"repo:a"}
	if _, err := Import(ImportRequest{
		Journal: j, TargetRepository: "repo:c", Snapshot: s,
		By: model.By{Who: "human"}, At: now,
	}); err != nil {
		t.Fatal(err)
	}
	ancestors, err := AncestorIDs(j.Dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := strings.Join(ancestors, ","); got != "repo:a,repo:b" {
		t.Fatalf("transitive ancestors = %q, want repo:a,repo:b", got)
	}

	// A generated C seed carries both IDs, so A can reject C as a cycle using
	// only the already-generated artifact.
	cSeed := snapshot("repo:c", "c", now, "topic")
	cSeed.Ancestors = ancestors
	target, _ := journal.Load(t.TempDir())
	if _, err := Import(ImportRequest{
		Journal: target, TargetRepository: "repo:a", Snapshot: cSeed,
		By: model.By{Who: "human"}, At: now.Add(time.Minute),
	}); err == nil || !strings.Contains(err.Error(), "lineage cycle") {
		t.Fatalf("multi-hop cycle error = %v", err)
	}
}

func TestLoadLinksIsAppendOnlyAndMalformedLinksAreLoud(t *testing.T) {
	dir := t.TempDir()
	s := snapshot("repo:old", "old", now.Add(-time.Hour), "topic")
	digest, _ := seed.Digest(s)
	link := &Link{
		ID: NewLinkID(now), TargetRepository: "repo:new", From: s.Repository,
		SeedDigest: digest, By: model.By{Who: "human"}, At: now,
	}
	if err := AddLink(dir, link); err != nil {
		t.Fatal(err)
	}
	if err := AddLink(dir, link); err == nil {
		t.Fatal("append-only link was overwritten")
	}
	bad := filepath.Join(dir, linkDir, "broken.yaml")
	if err := os.WriteFile(bad, []byte("not: [valid"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := LoadLinks(dir); err == nil || !strings.Contains(err.Error(), "broken.yaml") {
		t.Fatalf("malformed link was silent: %v", err)
	}
}
