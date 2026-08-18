package docket

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/state"
)

var testNow = time.Date(2026, 8, 16, 20, 0, 0, 0, time.UTC)

func TestRenderRejectsSyntheticFYIWithoutWriting(t *testing.T) {
	card := validCard("fyi")
	card.Shape = FYIShape
	var out bytes.Buffer
	err := Render(&out, View{Repo: "repo", Now: testNow, Cards: []Card{card}})
	if !errors.Is(err, ErrNotDecisionShaped) {
		t.Fatalf("error = %v, want ErrNotDecisionShaped", err)
	}
	if out.Len() != 0 {
		t.Fatalf("renderer wrote rejected FYI: %q", out.String())
	}
}

func TestRenderRequiresOneToThreeDiscreteAnswers(t *testing.T) {
	for _, n := range []int{0, 4} {
		t.Run(fmt.Sprintf("answers-%d", n), func(t *testing.T) {
			card := validCard("verbs")
			card.Answers = nil
			for i := 0; i < n; i++ {
				card.Answers = append(card.Answers, Verb{Name: fmt.Sprintf("answer-%d", i)})
			}
			err := card.Validate()
			if !errors.Is(err, ErrNotDecisionShaped) {
				t.Fatalf("error = %v, want ErrNotDecisionShaped", err)
			}
		})
	}
}

func TestBuildPromotionCandidateIsDecisionShaped(t *testing.T) {
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	e := &model.Entry{
		ID: ids.NewEntry(testNow), Type: model.Finding,
		Title: "Verify completion directly", Body: "A durable cross-project rule.",
		Quote:       "verify the affected state before declaring completion",
		UtteranceBy: model.ByUser, Confidence: .9,
		Source: model.Source{Kind: model.SrcSession, Ref: "session#L7", At: testNow},
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	cards := Build(Input{Journal: j, Now: testNow.Add(time.Minute), Alerts: []state.Alert{{
		Key: "promotion:" + e.ID, Kind: "promotion", EntryIDs: e.ID, Blocking: true, CreatedAt: testNow,
		WithdrawWhen: "promotion:" + e.ID + ":ruled",
	}}})
	if len(cards) != 1 {
		t.Fatalf("promotion cards = %d, want 1", len(cards))
	}
	if err := cards[0].Validate(); err != nil {
		t.Fatalf("promotion card invalid: %v", err)
	}
	if cards[0].Answers[0].Name != "promote" || cards[0].Answers[1].Name != "keep-local" {
		t.Fatalf("promotion verbs = %#v", cards[0].Answers)
	}
}

func TestRenderEighthCardBecomesOneOverflowFailure(t *testing.T) {
	cards := make([]Card, 8)
	for i := range cards {
		cards[i] = validCard(fmt.Sprintf("card-%d", i))
	}
	precision := &PushPrecision{Needed: 3, Total: 4}
	got := renderString(t, View{Repo: "repo", Now: testNow, Cards: cards, PushPrecision: precision})

	if n := strings.Count(got, "┌─"); n != 1 {
		t.Fatalf("rendered %d cards, want one overflow card:\n%s", n, got)
	}
	if !strings.Contains(got, "8 more items — the system is misconfigured; push-precision report attached") {
		t.Fatalf("missing exact overflow failure: %s", got)
	}
	if !strings.Contains(got, "push precision: 75.0% (3 needed / 4 total; 1 unneeded failure)") {
		t.Fatalf("missing attached precision: %s", got)
	}
	if !strings.Contains(got, "leaves when the decision-card count returns to seven or fewer") || strings.Contains(got, "when=") {
		t.Fatalf("overflow withdrawal leaked its machine condition: %s", got)
	}
	for i := range cards {
		if strings.Contains(got, cards[i].Headline) {
			t.Fatalf("normal card %d leaked through overflow: %s", i, got)
		}
	}
}

func TestRenderSevenCardsWithoutScrollTail(t *testing.T) {
	cards := make([]Card, MaxCards)
	for i := range cards {
		cards[i] = validCard(fmt.Sprintf("card-%d", i))
	}
	got := renderString(t, View{Repo: "repo", Now: testNow, Cards: cards})
	if n := strings.Count(got, "┌─ DECIDE"); n != MaxCards {
		t.Fatalf("rendered %d decision cards, want %d", n, MaxCards)
	}
	if strings.Contains(got, "more items") || strings.Contains(got, "scroll") {
		t.Fatalf("seven-card view grew a tail: %s", got)
	}
}

func TestOverflowStillRejectsFYIInHiddenInput(t *testing.T) {
	cards := make([]Card, 8)
	for i := range cards {
		cards[i] = validCard(fmt.Sprintf("card-%d", i))
	}
	cards[7].Shape = FYIShape
	var out bytes.Buffer
	err := Render(&out, View{Cards: cards, Now: testNow})
	if !errors.Is(err, ErrNotDecisionShaped) || out.Len() != 0 {
		t.Fatalf("FYI was hidden by overflow: err=%v out=%q", err, out.String())
	}
}

func TestOverflowPushPrecisionPlaceholderAndValidation(t *testing.T) {
	cards := make([]Card, 8)
	for i := range cards {
		cards[i] = validCard(fmt.Sprintf("card-%d", i))
	}
	got := renderString(t, View{Cards: cards, Now: testNow})
	if !strings.Contains(got, "push precision: unavailable (metric not supplied)") {
		t.Fatalf("missing precision placeholder: %s", got)
	}

	var out bytes.Buffer
	err := Render(&out, View{Cards: cards, Now: testNow, PushPrecision: &PushPrecision{Needed: 2, Total: 1}})
	if err == nil || out.Len() != 0 {
		t.Fatalf("invalid precision rendered: err=%v out=%q", err, out.String())
	}
}

func TestBuildAndRenderNeverReadJournalOrAlertParaphrases(t *testing.T) {
	question := entry(model.Question, testNow.Add(-3*time.Hour), "FORBIDDEN TITLE PARAPHRASE", "exact owner words without a question mark")
	question.Body = "FORBIDDEN JOURNAL BODY PARAPHRASE"
	question.Asks = "human"
	j := &journal.Journal{Entries: map[string]*model.Entry{question.ID: question}}
	alerts := []state.Alert{
		{
			Key: "adapter:/tmp/session.jsonl", Kind: "adapter", Blocking: true,
			Body: "FORBIDDEN ALERT BODY PARAPHRASE", CreatedAt: testNow.Add(-15 * time.Minute),
			WithdrawWhen: "adapter:session:resumed",
		},
		{Key: "suspect:fyi", Kind: "suspect", Blocking: false, Body: "FORBIDDEN FYI"},
		{Key: "unknown:blocking", Kind: "unknown", Blocking: true, Body: "FORBIDDEN UNKNOWN"},
	}
	cards := Build(Input{Journal: j, Alerts: alerts, Now: testNow})
	if len(cards) != 2 {
		t.Fatalf("Build returned %d cards, want question + adapter: %#v", len(cards), cards)
	}
	got := renderString(t, View{Repo: "repo", Cards: cards, Now: testNow})

	for _, forbidden := range []string{
		"FORBIDDEN TITLE PARAPHRASE", "FORBIDDEN JOURNAL BODY PARAPHRASE",
		"FORBIDDEN ALERT BODY PARAPHRASE", "FORBIDDEN FYI", "FORBIDDEN UNKNOWN",
	} {
		if strings.Contains(got, forbidden) {
			t.Fatalf("paraphrase %q leaked:\n%s", forbidden, got)
		}
	}
	if !strings.Contains(got, strconv.Quote(question.Quote)) {
		t.Fatalf("exact quote missing: %s", got)
	}
	if strings.Contains(got, question.Source.Ref) || strings.Contains(got, "adapter:/tmp/session.jsonl") {
		t.Fatalf("machine provenance leaked: %s", got)
	}
	if !strings.Contains(got, strconv.Quote("The watcher found this while checking current work")) {
		t.Fatalf("calm watcher evidence missing: %s", got)
	}
}

func TestBuildUsesShortVerbatimQuestionAsHeadline(t *testing.T) {
	question := entry(model.Question, testNow, "extractor title must not appear", "Rotate attestation tokens per launch?")
	question.Asks = "human"
	j := &journal.Journal{Entries: map[string]*model.Entry{question.ID: question}}
	cards := Build(Input{Journal: j, Now: testNow})
	if len(cards) != 1 {
		t.Fatalf("cards = %d, want 1", len(cards))
	}
	if cards[0].Headline != question.Quote {
		t.Fatalf("headline = %q, want exact quote %q", cards[0].Headline, question.Quote)
	}
	if strings.Contains(cards[0].Headline, question.Title) {
		t.Fatal("extractor title was used as headline")
	}
}

func TestQuestionHeadlineFallsBackRatherThanTruncatingOrParaphrasing(t *testing.T) {
	long := strings.Repeat("界", MaxHeadline) + "?"
	question := entry(model.Question, testNow, "forbidden title", long)
	question.Asks = "human"
	j := &journal.Journal{Entries: map[string]*model.Entry{question.ID: question}}
	card := Build(Input{Journal: j, Now: testNow})[0]
	if card.Headline != "How should this human-only question be answered?" {
		t.Fatalf("long quote was altered into headline: %q", card.Headline)
	}
	if utf8.RuneCountInString(card.Headline) > MaxHeadline {
		t.Fatalf("fallback headline too long: %d", utf8.RuneCountInString(card.Headline))
	}
	if card.Evidence[0].Text != long {
		t.Fatal("exact long evidence was modified")
	}
}

func TestValidateHeadlineIsQuestionAndAtMostEightyRunes(t *testing.T) {
	card := validCard("headline")
	card.Headline = strings.Repeat("x", MaxHeadline) + "?"
	if err := card.Validate(); err == nil {
		t.Fatal("accepted 81-rune headline")
	}
	card.Headline = "This is not a question"
	if err := card.Validate(); err == nil {
		t.Fatal("accepted statement headline")
	}
}

func TestEvidenceMustBeVerbatimAndCarryProvenance(t *testing.T) {
	card := validCard("evidence")
	card.Evidence[0].Verbatim = false
	if err := card.Validate(); err == nil {
		t.Fatal("accepted non-verbatim evidence")
	}
	card = validCard("evidence")
	card.Evidence[0].Provenance.Ref = ""
	if err := card.Validate(); err == nil {
		t.Fatal("accepted evidence without provenance")
	}

	card = validCard("evidence")
	card.Evidence[0].Text = "exact line one\n\"exact line two\""
	got := renderString(t, View{Cards: []Card{card}, Now: testNow})
	if !strings.Contains(got, strconv.Quote(card.Evidence[0].Text)) {
		t.Fatalf("multiline evidence was not reversibly exact: %s", got)
	}
}

func TestAssumptionAppearsExactlyOnHighMagnitudeOnly(t *testing.T) {
	normal := validCard("normal")
	normal.Assumption = "forbidden on normal"
	if err := normal.Validate(); err == nil {
		t.Fatal("accepted assumption on normal-magnitude card")
	}

	high := validCard("high")
	high.Magnitude = HighMagnitude
	if err := high.Validate(); err == nil {
		t.Fatal("accepted high-magnitude card without assumption")
	}
	high.Assumption = "production latency matches the cited measurement"
	got := renderString(t, View{Cards: []Card{high}, Now: testNow})
	line := "accepting this assumes: " + high.Assumption
	if strings.Count(got, line) != 1 {
		t.Fatalf("assumption line count != 1: %s", got)
	}
}

func TestForeignImportIsHighMagnitudeAndNeedsCallerAssumption(t *testing.T) {
	alert := state.Alert{
		Key: "import:batch-1", Kind: "import", Blocking: true, EntryIDs: "batch-1",
		CreatedAt: testNow, WithdrawWhen: "import:batch-1:closed",
	}
	cards := Build(Input{Alerts: []state.Alert{alert}, Now: testNow})
	if len(cards) != 1 || cards[0].Magnitude != HighMagnitude {
		t.Fatalf("foreign import card = %#v", cards)
	}
	if err := cards[0].Validate(); err == nil {
		t.Fatal("foreign import rendered without high-magnitude assumption")
	}

	key := cards[0].Key
	cards = Build(Input{
		Alerts: []state.Alert{alert}, Now: testNow,
		Assumptions: map[string]string{key: "the foreign provenance is authentic"},
		Evidence: map[string][]Evidence{key: {{
			Text: "exact proposed words", Verbatim: true,
			Provenance: Provenance{Kind: "foreign", Ref: "bundle.yaml#entry-1"},
		}}},
	})
	if err := cards[0].Validate(); err != nil {
		t.Fatalf("foreign import with assumption invalid: %v", err)
	}
	got := renderString(t, View{Cards: cards, Now: testNow})
	if !strings.Contains(got, `"exact proposed words"`) || !strings.Contains(got, "[enter] open batch diff") {
		t.Fatalf("proposal evidence/open missing: %s", got)
	}
}

func TestDeferralIsEventBoundAndWithdrawalIsPrintedMachineReadable(t *testing.T) {
	card := validCard("withdraw")
	card.Defer.Until = ""
	if err := card.Validate(); err == nil {
		t.Fatal("accepted defer without an event")
	}
	card = validCard("withdraw")
	card.Withdrawal.When = "human prose with spaces"
	if err := card.Validate(); err == nil {
		t.Fatal("accepted non-machine withdrawal")
	}
	card = validCard("withdraw")
	got := renderString(t, View{Cards: []Card{card}, Now: testNow})
	if !strings.Contains(got, "[defer]") || strings.Contains(got, "entry:e1:next-event") {
		t.Fatalf("human defer verb missing or machine condition leaked: %s", got)
	}
	if !strings.Contains(got, "leaves when") || strings.Contains(got, "when=") {
		t.Fatalf("human withdrawal line leaked its machine condition: %s", got)
	}
}

func TestOrderingIsRunningAgentsThenBlockingCostThenAge(t *testing.T) {
	waiting := validCard("waiting")
	waiting.Headline = "Waiting card?"
	waiting.Why.Class = CostWaiting
	waiting.Why.Since = testNow.Add(-10 * time.Hour)

	irreversible := validCard("irreversible")
	irreversible.Headline = "Irreversible card?"
	irreversible.Why.Class = CostIrreversible
	irreversible.Why.Since = testNow.Add(-time.Hour)

	running := validCard("running")
	running.Headline = "Running agents card?"
	running.Why.Class = CostWaiting
	running.Why.RunningAgents = 2
	running.Why.Since = testNow.Add(-5 * time.Minute)

	got := renderString(t, View{Cards: []Card{waiting, irreversible, running}, Now: testNow})
	positions := []int{
		strings.Index(got, running.Headline),
		strings.Index(got, irreversible.Headline),
		strings.Index(got, waiting.Headline),
	}
	if positions[0] < 0 || !(positions[0] < positions[1] && positions[1] < positions[2]) {
		t.Fatalf("wrong order %v:\n%s", positions, got)
	}
}

func TestStallTimerTicksAtRenderTime(t *testing.T) {
	card := validCard("timer")
	card.Why.Since = testNow.Add(-5 * time.Minute)
	card.Why.RunningAgents = 2
	first := renderString(t, View{Cards: []Card{card}, Now: testNow})
	second := renderString(t, View{Cards: []Card{card}, Now: testNow.Add(time.Minute)})
	if !strings.Contains(first, "· 5m ─") || !strings.Contains(second, "· 6m ─") {
		t.Fatalf("timer did not tick:\nfirst=%s\nsecond=%s", first, second)
	}
}

func TestExactParameterizedEmptyState(t *testing.T) {
	got := renderString(t, View{
		Repo: "clew", Now: testNow,
		Empty: EmptyMetrics{DaysSinceLastRuling: 2, EntriesLearnedSince: 14},
	})
	want := "DOCKET — clew\n\nNothing needs you · last ruling 2d ago · 14 entries learned since.\n"
	if got != want {
		t.Fatalf("empty state:\n got %q\nwant %q", got, want)
	}
}

func TestEmptyMetricsAtUsesWholeDays(t *testing.T) {
	got := EmptyMetricsAt(testNow, testNow.Add(-49*time.Hour), 6)
	if got.DaysSinceLastRuling != 2 || got.EntriesLearnedSince != 6 {
		t.Fatalf("metrics = %#v", got)
	}
}

func TestBuildBatchesAgingAlertWithItsQuestion(t *testing.T) {
	question := entry(model.Question, testNow.Add(-8*24*time.Hour), "forbidden title", "Which release should ship?")
	question.Asks = "human"
	j := &journal.Journal{Entries: map[string]*model.Entry{question.ID: question}}
	alert := state.Alert{
		Key: "aging:" + question.ID, Kind: "aging", EntryIDs: question.ID, Blocking: true,
		CreatedAt: testNow.Add(-time.Hour), WithdrawWhen: "question:" + question.ID + ":closed",
	}
	cards := Build(Input{Journal: j, Alerts: []state.Alert{alert}, Now: testNow})
	if len(cards) != 1 || cards[0].Kind != "question" {
		t.Fatalf("question identity was not batched: %#v", cards)
	}
}

func TestBuildFiltersStaleJournalAlerts(t *testing.T) {
	intent := entry(model.Intent, testNow.Add(-time.Hour), "intent title", "exact intent words")
	decisionA := entry(model.Decision, testNow.Add(-time.Hour), "decision a", "exact decision A")
	decisionB := entry(model.Decision, testNow.Add(-time.Hour), "decision b", "exact decision B")
	// No shared tags means the decisions are active, not contradictory; an
	// old alert must not resurrect them into the docket.
	j := &journal.Journal{Entries: map[string]*model.Entry{
		intent.ID: intent, decisionA.ID: decisionA, decisionB.ID: decisionB,
	}}
	alerts := []state.Alert{
		{Key: "absence:old", Kind: "absence", EntryIDs: intent.ID, Blocking: true, WithdrawWhen: "intent:active"},
		{Key: "contradiction:old", Kind: "contradiction", EntryIDs: decisionA.ID + "+" + decisionB.ID, Blocking: true, WithdrawWhen: "pair:resolved"},
	}
	if cards := Build(Input{Journal: j, Alerts: alerts, Now: testNow}); len(cards) != 0 {
		t.Fatalf("stale journal alerts rendered: %#v", cards)
	}
}

func TestBuildDropsResolvedOrDismissedStomp(t *testing.T) {
	active := state.Alert{
		Key: "stomp:file.go", Kind: "stomp", EntryIDs: "file.go", Blocking: true,
		CreatedAt: testNow.Add(-5 * time.Minute), WithdrawWhen: "stomp:file.go:clean",
	}
	if cards := Build(Input{Alerts: []state.Alert{active}, Now: testNow}); len(cards) != 1 {
		t.Fatalf("active stomp cards = %d, want 1", len(cards))
	}
	active.DroppedAt = testNow.Format(time.RFC3339)
	if cards := Build(Input{Alerts: []state.Alert{active}, Now: testNow}); len(cards) != 0 {
		t.Fatalf("withdrawn stomp rendered: %#v", cards)
	}
}

func TestAlertWithdrawalTokenIsPreservedExactly(t *testing.T) {
	alert := state.Alert{
		Key: "stomp:file.go", Kind: "stomp", EntryIDs: "file.go", Blocking: true,
		CreatedAt: testNow, WithdrawWhen: "stomp:repo%2Ffile.go:merge-clean",
	}
	card := Build(Input{
		Alerts: []state.Alert{alert}, Now: testNow,
		Assumptions: map[string]string{"alert:" + alert.Key: "the selected owner has the work that should survive"},
	})[0]
	if card.Withdrawal.When != alert.WithdrawWhen {
		t.Fatalf("withdrawal = %q, want %q", card.Withdrawal.When, alert.WithdrawWhen)
	}
	got := renderString(t, View{Cards: []Card{card}, Now: testNow})
	if !strings.Contains(got, "leaves when") || strings.Contains(got, alert.WithdrawWhen) {
		t.Fatalf("human withdrawal line leaked its machine condition: %s", got)
	}
}

func TestQuestionWithdrawsAfterAnswerEvent(t *testing.T) {
	j := newJournal(t)
	question := entry(model.Question, testNow.Add(-time.Hour), "title not rendered", "Ship the guarded path?")
	question.Asks = "human"
	if err := j.AddEntry(question); err != nil {
		t.Fatal(err)
	}
	if got := Build(Input{Journal: j, Now: testNow}); len(got) != 1 {
		t.Fatalf("open question cards = %d, want 1", len(got))
	}
	event := &model.Event{
		ID: ids.NewEvent(testNow), Kind: model.EvAnswer, Entry: question.ID,
		Payload: map[string]any{"by": "human"}, By: model.By{Who: "human"}, At: testNow,
	}
	if err := j.AddEvent(event); err != nil {
		t.Fatal(err)
	}
	if got := Build(Input{Journal: j, Now: testNow.Add(time.Minute)}); len(got) != 0 {
		t.Fatalf("answered question remained: %#v", got)
	}
}

func TestKnownNonBlockingAndUnknownBlockingAlertsNeverBecomeCards(t *testing.T) {
	alerts := []state.Alert{
		{Key: "overlap:x", Kind: "overlap", Blocking: false},
		{Key: "suspect:x", Kind: "suspect", Blocking: false},
		{Key: "budget:x", Kind: "budget", Blocking: false},
		{Key: "new-kind:x", Kind: "new-kind", Blocking: true},
	}
	if cards := Build(Input{Alerts: alerts, Now: testNow}); len(cards) != 0 {
		t.Fatalf("FYI alerts became cards: %#v", cards)
	}
}

func validCard(key string) Card {
	return Card{
		Key:      key,
		Kind:     "question",
		Shape:    DecisionShape,
		Headline: "What should happen next?",
		Why: WhyYou{
			RuleCode: "human-only", Rule: "only the human can choose",
			CostOfDelay: "work remains blocked", Since: testNow.Add(-time.Hour), Class: CostBlockedWork,
		},
		Evidence: []Evidence{{
			Text: "exact source words", Verbatim: true,
			Provenance: Provenance{Kind: "session", Ref: "codex:/tmp/session.jsonl#L9", EntryID: "e1"},
		}},
		Magnitude: NormalMagnitude,
		Answers:   []Verb{{Name: "answer", Target: "e1"}},
		Defer: &Deferral{
			Verb: Verb{Name: "defer", Target: "e1"}, Until: "entry:e1:next-event",
		},
		Withdrawal: Withdrawal{When: "entry:e1:status!=open", Text: "the question closes"},
	}
}

func entry(typ model.EntryType, at time.Time, title, quote string) *model.Entry {
	return &model.Entry{
		ID: ids.NewEntry(at), Type: typ, Title: title, Body: "extractor body", Quote: quote,
		UtteranceBy: model.ByUser,
		Source:      model.Source{Kind: model.SrcSession, Ref: "codex:/tmp/session.jsonl#L9", At: at},
		Confidence:  0.95,
	}
}

func newJournal(t *testing.T) *journal.Journal {
	t.Helper()
	dir := t.TempDir()
	for _, child := range []string{"entries", "events"} {
		if err := os.MkdirAll(filepath.Join(dir, child), 0o755); err != nil {
			t.Fatal(err)
		}
	}
	j, err := journal.Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func renderString(t *testing.T, view View) string {
	t.Helper()
	var out bytes.Buffer
	if err := Render(&out, view); err != nil {
		t.Fatalf("Render: %v", err)
	}
	return out.String()
}
