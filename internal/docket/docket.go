// Package docket builds and renders the bounded, decision-only human ruling
// surface. It deliberately has no interaction loop: callers execute the
// printed verbs and rebuild the view from current journal and state data.
package docket

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/state"
)

const (
	// MaxCards is the I12 hard cap. Render replaces any larger set with one
	// overflow-failure card; it never emits a scrollable tail.
	MaxCards = 7
	// MaxHeadline is measured in runes, matching the journal title limit.
	MaxHeadline = 80
)

var (
	// ErrNotDecisionShaped means an informational/FYI item reached the docket.
	// Such an item is rejected rather than silently made to look actionable.
	ErrNotDecisionShaped = errors.New("docket item is not decision-shaped")
	// ErrInvalidCard means a nominal decision card violates its fixed anatomy.
	ErrInvalidCard = errors.New("invalid docket card")

	machineTokenRE = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._~:/@%+=,!<>-]*$`)
)

// Shape makes decision-ness explicit and non-default. A zero-value Card, or a
// synthetic FYI card, therefore cannot accidentally render.
type Shape string

const (
	DecisionShape Shape = "decision"
	FYIShape      Shape = "fyi"
)

// Magnitude controls the sole permitted cognitive forcing function.
type Magnitude string

const (
	NormalMagnitude Magnitude = "normal"
	HighMagnitude   Magnitude = "high"
)

// CostClass is the stable coarse ordering before elapsed delay. RunningAgents
// is considered first regardless of class.
type CostClass int

const (
	CostWaiting CostClass = iota + 1
	CostBlockedWork
	CostIrreversible
)

// WhyYou states the rule that fired and the cost of waiting. Since remains a
// timestamp so Render can make the delay timer tick on every refresh.
type WhyYou struct {
	RuleCode      string
	Rule          string
	CostOfDelay   string
	Since         time.Time
	Class         CostClass
	RunningAgents int
}

// Provenance is an exact source pointer. Ref is intentionally not rewritten:
// session line, commit, and entry references remain openable by callers that
// know their native scheme.
type Provenance struct {
	Kind    string
	Ref     string
	EntryID string
}

// Evidence is rendered only when Verbatim is true. Text is quoted with a
// reversible escaping, never clipped, summarized, or paraphrased.
type Evidence struct {
	Text       string
	Provenance Provenance
	Verbatim   bool
}

// Verb is one discrete ruling. Name is the machine verb; Label is optional UI
// copy; Target is an opaque entry/alert identity for command integration.
type Verb struct {
	Name   string
	Label  string
	Target string
}

// Deferral is event-bound, never duration-bound. Until must name the machine
// event which causes the card to be reconsidered or withdrawn.
type Deferral struct {
	Verb  Verb
	Until string
}

// Withdrawal is printed in both human and machine-checkable form. When is an
// opaque token evaluated by the subsystem that owns the source item.
type Withdrawal struct {
	When string
	Text string
}

// Card has no Body, Summary, or Explanation field by design. The only source
// text it can carry is marked verbatim Evidence; a high-magnitude assumption
// is the sole exception.
type Card struct {
	Key        string
	Kind       string
	Shape      Shape
	Headline   string
	Why        WhyYou
	Evidence   []Evidence
	Magnitude  Magnitude
	Assumption string
	Answers    []Verb // one to three discrete ruling verbs
	Defer      *Deferral
	Redirect   *Verb
	Open       *Verb // opens the full source/diff; not a ruling verb
	Withdrawal Withdrawal
}

// Input is the narrow journal+state boundary used by CLI integration. The two
// maps are keyed by Card.Key. They allow a caller to mark a genuinely
// high-magnitude item without putting generated explanation text on the card.
type Input struct {
	Journal     *journal.Journal
	Alerts      []state.Alert
	Now         time.Time
	Magnitudes  map[string]Magnitude
	Assumptions map[string]string
	Evidence    map[string][]Evidence
}

// EmptyMetrics parameterize the designed empty state. They are not docket
// backlog or unread counts.
type EmptyMetrics struct {
	DaysSinceLastRuling int
	EntriesLearnedSince int
}

// PushPrecision is the dogfood metric attached to an overflow failure.
// Needed is the number of delivered pushes that truly needed the human.
type PushPrecision struct {
	Needed int
	Total  int
}

// View carries presentation-only context. Cards may come directly from Build
// or from another typed producer; Render enforces the same laws either way.
type View struct {
	Repo          string
	Cards         []Card
	Now           time.Time
	Empty         EmptyMetrics
	PushPrecision *PushPrecision
}

// EmptyMetricsAt computes the integer-day N used by the exact empty template.
func EmptyMetricsAt(now, lastRuling time.Time, learned int) EmptyMetrics {
	days := 0
	if !lastRuling.IsZero() && now.After(lastRuling) {
		days = int(now.Sub(lastRuling) / (24 * time.Hour))
	}
	return EmptyMetrics{DaysSinceLastRuling: days, EntriesLearnedSince: learned}
}

// Build derives current decision cards from journal entries and open state
// alerts. Findings, non-blocking alerts, stale journal statuses, and unknown
// alert kinds produce no card. Entry Body and Title are never read.
func Build(in Input) []Card {
	now := in.Now
	if now.IsZero() {
		now = time.Now()
	}

	var computed map[string]*journal.Computed
	if in.Journal != nil {
		computed = journal.Compute(in.Journal, now)
	}

	cards := make([]Card, 0)
	seen := make(map[string]bool)
	if in.Journal != nil {
		ids := make([]string, 0, len(in.Journal.Entries))
		for id := range in.Journal.Entries {
			ids = append(ids, id)
		}
		sort.Strings(ids)
		for _, id := range ids {
			e := in.Journal.Entries[id]
			c := computed[id]
			if e == nil || c == nil || e.Type != model.Question || e.Asks != "human" || c.Status != journal.StOpen {
				continue
			}
			card := questionCard(e)
			applyOverrides(&card, in)
			cards = append(cards, card)
			seen[card.Key] = true
		}
	}

	alerts := append([]state.Alert(nil), in.Alerts...)
	sort.SliceStable(alerts, func(i, k int) bool {
		if alerts[i].CreatedAt.Equal(alerts[k].CreatedAt) {
			return alerts[i].Key < alerts[k].Key
		}
		return alerts[i].CreatedAt.Before(alerts[k].CreatedAt)
	})
	for _, alert := range alerts {
		if !alert.Blocking || alert.AckedAt != "" || alert.DroppedAt != "" {
			continue
		}
		if alert.Kind == "aging" {
			if id := firstEntryID(alert.EntryIDs, in.Journal); id != "" && seen["question:"+id] {
				continue // one identity, one card
			}
		}
		card, ok := alertCard(alert, in.Journal, computed, now)
		if !ok || seen[card.Key] {
			continue
		}
		applyOverrides(&card, in)
		cards = append(cards, card)
		seen[card.Key] = true
	}

	sortCards(cards)
	return cards
}

func applyOverrides(card *Card, in Input) {
	if magnitude, ok := in.Magnitudes[card.Key]; ok {
		card.Magnitude = magnitude
	}
	if assumption, ok := in.Assumptions[card.Key]; ok {
		card.Assumption = assumption
	}
	if evidence, ok := in.Evidence[card.Key]; ok {
		card.Evidence = append([]Evidence(nil), evidence...)
	}
}

func questionCard(e *model.Entry) Card {
	id := e.ID
	return Card{
		Key:      "question:" + id,
		Kind:     "question",
		Shape:    DecisionShape,
		Headline: exactQuestionOr("How should this human-only question be answered?", e.Quote),
		Why: WhyYou{
			RuleCode: "human-only-question", Rule: "only the human can answer",
			CostOfDelay: "the question remains unresolved", Since: e.Created(), Class: CostBlockedWork,
		},
		Evidence:   []Evidence{entryEvidence(e)},
		Magnitude:  NormalMagnitude,
		Answers:    []Verb{{Name: "answer", Target: id}},
		Defer:      eventDeferral("entry:"+atom(id)+":next-event", id),
		Withdrawal: Withdrawal{When: "entry:" + atom(id) + ":status!=open", Text: "the question is answered or expires"},
	}
}

func alertCard(alert state.Alert, j *journal.Journal, computed map[string]*journal.Computed, now time.Time) (Card, bool) {
	ids := entryIDs(alert.EntryIDs, j)
	since := alert.CreatedAt
	if since.IsZero() {
		since = now
	}
	key := "alert:" + alert.Key
	card := Card{Key: key, Kind: alert.Kind, Shape: DecisionShape, Magnitude: NormalMagnitude}
	card.Evidence = alertEvidence(alert, ids, j)
	card.Withdrawal = Withdrawal{When: alert.WithdrawWhen}

	switch alert.Kind {
	case "contradiction":
		if !statusAmong(ids, computed, journal.StPossibleContradiction, journal.StContradicted) {
			return Card{}, false
		}
		card.Headline = "Which conflicting decision should govern?"
		card.Magnitude = HighMagnitude
		card.Why = WhyYou{RuleCode: "decision-contradiction", Rule: "active decisions conflict", CostOfDelay: "work may follow incompatible constraints", Since: since, Class: CostIrreversible}
		for i, id := range ids {
			if i == 2 {
				break
			}
			card.Answers = append(card.Answers, Verb{Name: fmt.Sprintf("keep-%d", i+1), Target: id})
		}
		if len(card.Answers) == 0 {
			card.Answers = []Verb{{Name: "resolve", Target: alert.Key}}
		}
		card.Redirect = &Verb{Name: "redirect", Target: alert.Key}
		card.Defer = eventDeferral("alert:"+atom(alert.Key)+":source-event", alert.Key)
		card.Withdrawal.Text = "the decisions no longer conflict or one is superseded"
	case "absence":
		if !statusAmong(ids, computed, journal.StAbsent) {
			return Card{}, false
		}
		card.Headline = "Should this absent intent be revived or dropped?"
		card.Why = WhyYou{RuleCode: "intent-absence", Rule: "sibling intents progressed without this one", CostOfDelay: "the intended work has no evidence", Since: since, Class: CostBlockedWork}
		card.Answers = []Verb{{Name: "revive", Target: alert.EntryIDs}, {Name: "drop", Target: alert.EntryIDs}}
		card.Redirect = &Verb{Name: "redirect", Target: alert.EntryIDs}
		card.Defer = eventDeferral("alert:"+atom(alert.Key)+":source-event", alert.Key)
		card.Withdrawal.Text = "the intent gains evidence, is revived, or is dropped"
	case "aging":
		if !openHumanQuestion(ids, j, computed) {
			return Card{}, false
		}
		card.Headline = "How should this aging human-only question be answered?"
		card.Why = WhyYou{RuleCode: "question-aging", Rule: "a human-only question is aging", CostOfDelay: "agents cannot close the question", Since: since, Class: CostBlockedWork}
		card.Answers = []Verb{{Name: "answer", Target: alert.EntryIDs}}
		card.Defer = eventDeferral("alert:"+atom(alert.Key)+":source-event", alert.Key)
		card.Withdrawal.Text = "the question is answered or expires"
	case "stomp":
		card.Headline = "How should the dirty same-file overlap be resolved?"
		card.Magnitude = HighMagnitude
		card.Why = WhyYou{RuleCode: "dirty-same-file-overlap", Rule: "running sessions overlap on a dirty file", CostOfDelay: "continued edits risk lost work", Since: since, Class: CostIrreversible, RunningAgents: 2}
		card.Answers = []Verb{{Name: "choose-owner", Target: alert.EntryIDs}}
		card.Redirect = &Verb{Name: "redirect", Target: alert.EntryIDs}
		card.Defer = eventDeferral("alert:"+atom(alert.Key)+":next-poll", alert.Key)
		card.Withdrawal.Text = "the branches merge cleanly or the same-file overlap ends"
	case "adapter":
		card.Headline = "How should the paused session adapter be recovered?"
		card.Why = WhyYou{RuleCode: "adapter-paused", Rule: "a session format cannot be parsed safely", CostOfDelay: "new session evidence remains parked", Since: since, Class: CostBlockedWork}
		card.Answers = []Verb{{Name: "update", Target: alert.Key}, {Name: "report", Target: alert.Key}}
		card.Defer = eventDeferral("alert:"+atom(alert.Key)+":adapter-resumed", alert.Key)
		card.Withdrawal.Text = "the adapter resumes parsing the source"
	case "import":
		card.Headline = "Should this foreign proposal batch be accepted?"
		card.Why = WhyYou{RuleCode: "foreign-proposal", Rule: "a foreign proposal needs owner confirmation", CostOfDelay: "the proposal cannot enter project memory", Since: since, Class: CostIrreversible}
		card.Magnitude = HighMagnitude
		card.Answers = []Verb{{Name: "accept", Target: alert.EntryIDs}, {Name: "reject", Target: alert.EntryIDs}}
		card.Redirect = &Verb{Name: "redirect", Target: alert.EntryIDs}
		card.Open = &Verb{Name: "open", Label: "open batch diff", Target: alert.EntryIDs}
		card.Defer = eventDeferral("alert:"+atom(alert.Key)+":proposal-change", alert.Key)
		card.Withdrawal.Text = "the proposal is accepted, rejected, or replaced"
	default:
		return Card{}, false // unknown and FYI-shaped alerts belong to the glance
	}

	if card.Withdrawal.When == "" {
		card.Withdrawal.When = "alert:" + atom(alert.Key) + ":inactive"
	}
	return card, true
}

func entryEvidence(e *model.Entry) Evidence {
	return Evidence{
		Text: e.Quote, Verbatim: true,
		Provenance: Provenance{Kind: string(e.Source.Kind), Ref: e.Source.Ref, EntryID: e.ID},
	}
}

func alertEvidence(alert state.Alert, ids []string, j *journal.Journal) []Evidence {
	out := make([]Evidence, 0, len(ids)+1)
	for _, id := range ids {
		if e := j.Entries[id]; e != nil {
			out = append(out, entryEvidence(e))
		}
	}
	if len(out) == 0 {
		text := alert.EntryIDs
		if text == "" {
			text = alert.Key
		}
		out = append(out, Evidence{
			Text: text, Verbatim: true,
			Provenance: Provenance{Kind: "state", Ref: "alert:" + alert.Key},
		})
	}
	return out
}

func entryIDs(raw string, j *journal.Journal) []string {
	if j == nil || raw == "" {
		return nil
	}
	parts := strings.FieldsFunc(raw, func(r rune) bool {
		switch r {
		case '+', ',', ' ', '\t', '\r', '\n':
			return true
		}
		return false
	})
	seen := map[string]bool{}
	ids := make([]string, 0, len(parts))
	for _, part := range parts {
		if _, ok := j.Entries[part]; ok && !seen[part] {
			ids = append(ids, part)
			seen[part] = true
		}
	}
	sort.Strings(ids)
	return ids
}

func firstEntryID(raw string, j *journal.Journal) string {
	ids := entryIDs(raw, j)
	if len(ids) == 0 {
		return ""
	}
	return ids[0]
}

func statusAmong(ids []string, computed map[string]*journal.Computed, statuses ...journal.Status) bool {
	if len(ids) == 0 {
		return true // state-only producer; its reconciliation owns freshness
	}
	want := map[journal.Status]bool{}
	for _, status := range statuses {
		want[status] = true
	}
	for _, id := range ids {
		if c := computed[id]; c != nil && want[c.Status] {
			return true
		}
	}
	return false
}

func openHumanQuestion(ids []string, j *journal.Journal, computed map[string]*journal.Computed) bool {
	if len(ids) == 0 {
		return true
	}
	for _, id := range ids {
		e := j.Entries[id]
		if e != nil && e.Type == model.Question && e.Asks == "human" && computed[id] != nil && computed[id].Status == journal.StOpen {
			return true
		}
	}
	return false
}

func eventDeferral(until, target string) *Deferral {
	return &Deferral{Verb: Verb{Name: "defer", Target: target}, Until: until}
}

func exactQuestionOr(fallback, quote string) string {
	q := strings.TrimSpace(quote)
	if !strings.ContainsAny(q, "\r\n") && utf8.RuneCountInString(q) <= MaxHeadline && strings.HasSuffix(q, "?") {
		return q
	}
	return fallback
}

func atom(s string) string { return url.QueryEscape(s) }

// Validate enforces I10 and the complete fixed card anatomy.
func (c Card) Validate() error {
	invalid := func(format string, args ...any) error {
		return fmt.Errorf("%w %q: %s", ErrInvalidCard, c.Key, fmt.Sprintf(format, args...))
	}
	if c.Shape != DecisionShape {
		return fmt.Errorf("%w %q (shape %q)", ErrNotDecisionShaped, c.Key, c.Shape)
	}
	if c.Key == "" || c.Kind == "" {
		return invalid("key and kind are required")
	}
	if strings.TrimSpace(c.Headline) == "" || !strings.HasSuffix(strings.TrimSpace(c.Headline), "?") {
		return invalid("headline must be an answerable question")
	}
	if utf8.RuneCountInString(c.Headline) > MaxHeadline {
		return invalid("headline exceeds %d runes", MaxHeadline)
	}
	if !token(c.Why.RuleCode) || strings.TrimSpace(c.Why.Rule) == "" || strings.TrimSpace(c.Why.CostOfDelay) == "" || c.Why.Since.IsZero() {
		return invalid("why-you requires a machine rule, rule text, cost of delay, and start time")
	}
	if c.Why.Class < CostWaiting || c.Why.Class > CostIrreversible || c.Why.RunningAgents < 0 {
		return invalid("invalid blocking cost")
	}
	if len(c.Evidence) == 0 {
		return invalid("at least one exact evidence row is required")
	}
	for i, row := range c.Evidence {
		if !row.Verbatim || strings.TrimSpace(row.Text) == "" {
			return invalid("evidence %d is not marked verbatim", i)
		}
		if strings.TrimSpace(row.Provenance.Kind) == "" || strings.TrimSpace(row.Provenance.Ref) == "" {
			return invalid("evidence %d lacks provenance", i)
		}
	}
	if c.Magnitude == "" {
		return invalid("magnitude is required")
	}
	switch c.Magnitude {
	case HighMagnitude:
		if strings.TrimSpace(c.Assumption) == "" {
			return invalid("high-magnitude card requires one assumption")
		}
	case NormalMagnitude:
		if c.Assumption != "" {
			return invalid("assumptions are forbidden below high magnitude")
		}
	default:
		return invalid("unknown magnitude %q", c.Magnitude)
	}
	if len(c.Answers) < 1 || len(c.Answers) > 3 {
		return fmt.Errorf("%w %q: requires 1-3 discrete verbs", ErrNotDecisionShaped, c.Key)
	}
	seenVerbs := map[string]bool{}
	for _, answer := range c.Answers {
		if err := validateVerb(answer); err != nil {
			return invalid("answer: %v", err)
		}
		if seenVerbs[answer.Name] {
			return invalid("duplicate verb %q", answer.Name)
		}
		seenVerbs[answer.Name] = true
	}
	if c.Defer != nil {
		if err := validateVerb(c.Defer.Verb); err != nil || c.Defer.Verb.Name != "defer" || !token(c.Defer.Until) {
			return invalid("defer must name a machine-checkable event")
		}
		if seenVerbs[c.Defer.Verb.Name] {
			return invalid("duplicate verb %q", c.Defer.Verb.Name)
		}
		seenVerbs[c.Defer.Verb.Name] = true
	}
	if c.Redirect != nil {
		if err := validateVerb(*c.Redirect); err != nil || c.Redirect.Name != "redirect" {
			return invalid("redirect must be a discrete redirect verb")
		}
		if seenVerbs[c.Redirect.Name] {
			return invalid("duplicate verb %q", c.Redirect.Name)
		}
	}
	if c.Open != nil {
		if err := validateVerb(*c.Open); err != nil || c.Open.Name != "open" {
			return invalid("open must be a discrete source-opening verb")
		}
	}
	if !token(c.Withdrawal.When) || strings.TrimSpace(c.Withdrawal.Text) == "" {
		return invalid("printed machine withdrawal condition is required")
	}
	return nil
}

func validateVerb(v Verb) error {
	if !token(v.Name) {
		return fmt.Errorf("invalid machine verb %q", v.Name)
	}
	return nil
}

func token(s string) bool { return machineTokenRE.MatchString(s) }

func sortCards(cards []Card) {
	sort.SliceStable(cards, func(i, k int) bool {
		a, b := cards[i], cards[k]
		if (a.Why.RunningAgents > 0) != (b.Why.RunningAgents > 0) {
			return a.Why.RunningAgents > 0
		}
		if a.Why.RunningAgents != b.Why.RunningAgents {
			return a.Why.RunningAgents > b.Why.RunningAgents
		}
		if a.Why.Class != b.Why.Class {
			return a.Why.Class > b.Why.Class
		}
		if !a.Why.Since.Equal(b.Why.Since) {
			return a.Why.Since.Before(b.Why.Since) // older delay first
		}
		return a.Key < b.Key
	})
}

// Render writes a complete snapshot. It validates before writing, sorts by
// blocking cost, and replaces eight or more valid cards with exactly one
// overflow-failure card.
func Render(w io.Writer, view View) error {
	if w == nil {
		return fmt.Errorf("render docket: nil writer")
	}
	now := view.Now
	if now.IsZero() {
		now = time.Now()
	}
	if view.Empty.DaysSinceLastRuling < 0 || view.Empty.EntriesLearnedSince < 0 {
		return fmt.Errorf("render docket: empty-state metrics cannot be negative")
	}
	if err := validatePushPrecision(view.PushPrecision); err != nil {
		return err
	}
	cards := append([]Card(nil), view.Cards...)
	for _, card := range cards {
		if err := card.Validate(); err != nil {
			return err
		}
	}
	sortCards(cards)

	var out strings.Builder
	repo := strings.TrimSpace(view.Repo)
	if repo == "" {
		repo = "current repo"
	}
	fmt.Fprintf(&out, "DOCKET — %s\n\n", repo)
	switch {
	case len(cards) == 0:
		fmt.Fprintf(&out, "Nothing needs you · last ruling %dd ago · %d entries learned since.\n",
			view.Empty.DaysSinceLastRuling, view.Empty.EntriesLearnedSince)
	case len(cards) > MaxCards:
		renderOverflow(&out, len(cards), view.PushPrecision)
	default:
		for i, card := range cards {
			if i > 0 {
				out.WriteByte('\n')
			}
			renderCard(&out, card, now)
		}
	}
	_, err := io.WriteString(w, out.String())
	return err
}

func renderCard(out *strings.Builder, card Card, now time.Time) {
	fmt.Fprintf(out, "┌─ DECIDE ─ blocking %s · %s ─\n", card.Why.RuleCode, elapsed(now, card.Why.Since))
	fmt.Fprintf(out, "│ %s\n", card.Headline)
	stall := ""
	if card.Why.RunningAgents > 0 {
		stall = fmt.Sprintf(" · %d running agents stalled", card.Why.RunningAgents)
	}
	fmt.Fprintf(out, "│ why you: %s fired · cost of delay: %s%s\n", card.Why.Rule, card.Why.CostOfDelay, stall)
	out.WriteString("│\n")
	for _, row := range card.Evidence {
		fmt.Fprintf(out, "│ ▸ %s   [%s · %s", strconv.Quote(row.Text), row.Provenance.Kind, row.Provenance.Ref)
		if row.Provenance.EntryID != "" {
			fmt.Fprintf(out, " · entry %s", row.Provenance.EntryID)
		}
		out.WriteString(" ↗]\n")
	}
	if card.Magnitude == HighMagnitude {
		out.WriteString("│\n")
		fmt.Fprintf(out, "│ accepting this assumes: %s\n", card.Assumption)
	}
	out.WriteString("│\n│ ")
	for i, answer := range card.Answers {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprintf(out, "[%s]", verbLabel(answer))
	}
	if card.Redirect != nil {
		fmt.Fprintf(out, " [%s]", verbLabel(*card.Redirect))
	}
	if card.Defer != nil {
		fmt.Fprintf(out, " [%s → %s]", verbLabel(card.Defer.Verb), card.Defer.Until)
	}
	if card.Open != nil {
		fmt.Fprintf(out, " [enter] %s", verbLabel(*card.Open))
	}
	out.WriteByte('\n')
	fmt.Fprintf(out, "│ withdraws itself if %s [when=%s]\n", card.Withdrawal.Text, card.Withdrawal.When)
	out.WriteString("└─\n")
}

func renderOverflow(out *strings.Builder, n int, precision *PushPrecision) {
	out.WriteString("┌─ SYSTEM FAILURE ─ I12 docket overflow ─\n")
	fmt.Fprintf(out, "│ %d more items — the system is misconfigured; push-precision report attached\n", n)
	fmt.Fprintf(out, "│ %s\n", pushPrecisionText(precision))
	out.WriteString("│ [inspect-config]\n")
	out.WriteString("│ withdraws itself if the decision-card count returns to seven or fewer [when=docket:card-count<=7]\n")
	out.WriteString("└─\n")
}

func verbLabel(v Verb) string {
	if v.Label != "" {
		return v.Label
	}
	return v.Name
}

func elapsed(now, since time.Time) string {
	d := now.Sub(since)
	if d < 0 {
		d = 0
	}
	d = d.Truncate(time.Minute)
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d/time.Minute))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh%02dm", int(d/time.Hour), int((d%time.Hour)/time.Minute))
	}
	return fmt.Sprintf("%dd%02dh", int(d/(24*time.Hour)), int((d%(24*time.Hour))/time.Hour))
}

func validatePushPrecision(p *PushPrecision) error {
	if p == nil {
		return nil
	}
	if p.Needed < 0 || p.Total < 0 || p.Needed > p.Total {
		return fmt.Errorf("render docket: invalid push precision %d/%d", p.Needed, p.Total)
	}
	return nil
}

func pushPrecisionText(p *PushPrecision) string {
	if p == nil {
		return "push precision: unavailable (metric not supplied)"
	}
	if p.Total == 0 {
		return "push precision: N/A (0 needed / 0 total)"
	}
	unneeded := p.Total - p.Needed
	return fmt.Sprintf("push precision: %.1f%% (%d needed / %d total; %d unneeded failure%s)",
		100*float64(p.Needed)/float64(p.Total), p.Needed, p.Total, unneeded, plural(unneeded))
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
