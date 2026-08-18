// Package extract implements the extractor (JOURNAL_SPEC §6): incremental,
// schema-validated LLM distillation of transcript slices into entries and
// events. Validation failure → one retry → park the slice + red status (I2).
// Every entry must carry a verbatim quote or it is rejected (I7).
package extract

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"clew/internal/adapters"
	"clew/internal/config"
	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/model"
	"clew/internal/scrub"
	"clew/internal/state"
)

//go:embed assets/instruction.md
var Instruction string

const (
	SliceCap     = 100 * 1024 // §6.2: cap 100 KB, middle-elision beyond
	IdleTrigger  = 2 * time.Minute
	BytesTrigger = 50 * 1024 // §6.1
)

// Outcome reports one extraction call's results.
type Outcome struct {
	Entries []*model.Entry
	// PromotionCandidates are accepted findings the extractor judged useful
	// across unrelated projects. They are proposals only: the watcher may put
	// them in the owner's docket, but only `clew journal promote` can copy one
	// into the owner-scope law journal.
	PromotionCandidates []*model.Entry
	Events              []*model.Event
	Rejected            int // entries rejected by validation (fabricated quote etc.)
	Redactions          int
	Tokens              int
	ObservedTokens      int // transcript tokens consumed (byte estimate), for backfill metering
	Parked              bool
	ParkReason          string
	NewOffset           int64
}

// wire mirrors the strict output schema of the instruction.
type wire struct {
	Entries []struct {
		Type               string     `json:"type"`
		Title              string     `json:"title"`
		Body               string     `json:"body"`
		Quote              string     `json:"quote"`
		Line               int        `json:"line"`
		UtteranceBy        string     `json:"utterance_by"`
		Confidence         float64    `json:"confidence"`
		Tags               []string   `json:"tags"`
		Env                *model.Env `json:"env"`
		Affects            []string   `json:"affects"`
		Asks               string     `json:"asks"`
		PromotionCandidate bool       `json:"promotion_candidate"`
	} `json:"entries"`
	Supersedes []struct{ Old, By string }      `json:"supersedes"`
	Answers    []struct{ Question, By string } `json:"answers"`
	Links      []struct {
		Entry string `json:"entry"`
		Line  int    `json:"line"`
	} `json:"links"`
}

// Gate enforces the unattended-live I9 budget: extraction tokens ≤
// session_pct% of observed session tokens today AND under the absolute daily
// cap. Cap hit = pause loudly; sensors keep recording.
func Gate(db *state.DB, cfg *config.Config, estimate int) (bool, string) {
	spent := db.TokensToday("spent")
	observed := db.TokensToday("observed")
	cap := cfg.Extractor.DailyCapTokens
	if spent+estimate > cap {
		return false, fmt.Sprintf("daily cap: spent %d + est %d > %d", spent, estimate, cap)
	}
	sessionSpent := db.TokensToday("extraction-spent")
	pctBudget := int(cfg.Extractor.SessionPct / 100 * float64(observed))
	if sessionSpent+estimate > pctBudget {
		return false, fmt.Sprintf("2%%-rule: spent %d + est %d > %d (%.0f%% of %d observed)",
			sessionSpent, estimate, pctBudget, cfg.Extractor.SessionPct, observed)
	}
	return true, ""
}

// GateDaily applies the absolute machine-wide LLM cap to non-session calls,
// which have no honest session-token denominator (archaeology and differ).
func GateDaily(db *state.DB, cfg *config.Config, estimate int) (bool, string) {
	spent := db.TokensToday("spent")
	if spent+estimate > cfg.Extractor.DailyCapTokens {
		return false, fmt.Sprintf("daily cap: spent %d + est %d > %d", spent, estimate, cfg.Extractor.DailyCapTokens)
	}
	return true, ""
}

// Run extracts from one session file starting at the extraction watermark.
// The caller persists outcome entries/events (already applied to j) and the
// new watermark.
func Run(j *journal.Journal, p llm.Provider, a adapters.Adapter, file string, offset int64, surface string, now time.Time) (*Outcome, error) {
	return run(j, p, a, file, offset, -1, surface, now)
}

// RunUntil is the explicit-backfill path. end is the immutable enrollment
// boundary, so history extraction can never consume a later live suffix.
func RunUntil(j *journal.Journal, p llm.Provider, a adapters.Adapter, file string, offset, end int64, surface string, now time.Time) (*Outcome, error) {
	return run(j, p, a, file, offset, end, surface, now)
}

func run(j *journal.Journal, p llm.Provider, a adapters.Adapter, file string, offset, end int64, surface string, now time.Time) (*Outcome, error) {
	var d *adapters.Delta
	var err error
	if end >= 0 {
		d, err = adapters.ParseRange(a, file, offset, end)
	} else {
		d, err = a.Parse(file, offset)
	}
	if err != nil {
		return nil, err
	}
	out := &Outcome{NewOffset: d.NewOffset, ObservedTokens: d.Bytes / 4}
	if len(d.Messages) == 0 {
		return out, nil
	}
	transcript := renderSlice(d.Messages)
	digest := readDigest(j)
	prompt := Instruction +
		"\n\n## Current journal digest\n\n" + digest +
		"\n\n## New transcript slice\n\nSession file: " + file + " (agent: " + d.Agent + ")\n\n" + transcript

	var w *wire
	var tokens int
	for attempt := 0; attempt < 2; attempt++ { // §6.2: one retry
		pr := prompt
		if attempt == 1 {
			pr += "\n\nREMINDER: output STRICT JSON only, exactly matching the schema. No prose."
		}
		res, err := p.Call(pr)
		if err != nil {
			return nil, err
		}
		tokens += res.Tokens
		raw, ok := llm.ExtractJSON(res.Text)
		if !ok {
			continue
		}
		cand := &wire{}
		if validWireEnvelope([]byte(raw), cand) {
			w = cand
			break
		}
	}
	out.Tokens = tokens
	if w == nil {
		out.Parked = true
		out.ParkReason = "schema validation failed after retry"
		out.NewOffset = offset // do not consume; slice is parked raw
		return out, nil
	}

	// Validate + apply entries.
	newIDs := make([]string, len(w.Entries))
	for i, we := range w.Entries {
		e, ok := validateEntry(we.Type, we.Title, we.Body, we.Quote, we.Line, we.UtteranceBy,
			we.Confidence, we.Tags, we.Env, we.Affects, we.Asks, d, surface)
		if !ok {
			out.Rejected++
			continue
		}
		var n int
		e.Quote, n = scrub.Scrub(e.Quote)
		out.Redactions += n
		e.Body, n = scrub.Scrub(e.Body)
		out.Redactions += n
		e.Title, n = scrub.Scrub(e.Title)
		out.Redactions += n
		// Persist the proposal bit on the immutable project finding. It is not
		// authority: it only lets any watcher rebuild the owner's ruling card if
		// local state is lost. Keep the mechanical boundary deliberately narrow.
		e.PromotionCandidate = we.PromotionCandidate && e.Type == model.Finding &&
			e.UtteranceBy != model.ByToolResult && e.Confidence >= 0.8 &&
			len(e.Tags) == 0 && e.Env == nil && len(e.Affects) == 0 &&
			!journal.Imperative(e)
		if err := j.AddEntry(e); err != nil {
			out.Rejected++
			continue
		}
		newIDs[i] = e.ID
		out.Entries = append(out.Entries, e)
		// A provider flag is never authority to create an ambient law. Keep the
		// proposal boundary narrow and reject obviously project/environment-bound
		// candidates mechanically before the watcher surfaces them to the owner.
		if e.PromotionCandidate {
			out.PromotionCandidates = append(out.PromotionCandidates, e)
		}
	}

	resolve := func(ref string) string {
		if s, ok := strings.CutPrefix(ref, "new:"); ok {
			var idx int
			if _, err := fmt.Sscanf(s, "%d", &idx); err == nil && idx >= 0 && idx < len(newIDs) {
				return newIDs[idx]
			}
			return ""
		}
		if _, exists := j.Entries[ref]; exists {
			return ref
		}
		return ""
	}
	by := model.By{Who: "extractor", Surface: surface}
	addEvent := func(kind model.EventKind, entry string, payload map[string]any) {
		v := &model.Event{ID: ids.NewEvent(now), Kind: kind, Entry: entry, Payload: payload, By: by, At: now.UTC()}
		if err := j.AddEvent(v); err == nil {
			out.Events = append(out.Events, v)
		}
	}
	for _, s := range w.Supersedes {
		old, nw := resolve(s.Old), resolve(s.By)
		if old != "" && nw != "" && !j.HasEvent(model.EvSupersede, old, "by", nw) {
			addEvent(model.EvSupersede, old, map[string]any{"by": nw})
		}
	}
	for _, an := range w.Answers {
		q, nw := resolve(an.Question), resolve(an.By)
		if q != "" && nw != "" && !j.HasEvent(model.EvAnswer, q, "by", nw) {
			addEvent(model.EvAnswer, q, map[string]any{"by": nw})
		}
	}
	for _, l := range w.Links {
		en := resolve(l.Entry)
		if en == "" {
			continue
		}
		ref := fmt.Sprintf("%s#L%d", file, l.Line)
		if !j.HasEvent(model.EvEvidence, en, "ref", ref) {
			addEvent(model.EvEvidence, en, map[string]any{"kind": "session", "ref": ref})
		}
	}
	return out, nil
}

func validWireEnvelope(raw []byte, dst *wire) bool {
	var object map[string]json.RawMessage
	if json.Unmarshal(raw, &object) != nil {
		return false
	}
	entries, ok := object["entries"]
	if !ok || string(entries) == "null" {
		return false // {} or a drifted provider envelope is not success (I2)
	}
	var shape []json.RawMessage
	if json.Unmarshal(entries, &shape) != nil {
		return false
	}
	return json.Unmarshal(raw, dst) == nil
}

// validateEntry enforces the schema (§6.2) and I7: the quote must be found
// verbatim (whitespace-normalized) in a transcript message; utterance_by is
// overridden by the actual role of the message the quote was found in, so a
// web-page quote can never dodge taint by claiming to be the assistant.
func validateEntry(typ, title, body, quote string, line int, uttBy string,
	conf float64, tags []string, env *model.Env, affects []string, asks string,
	d *adapters.Delta, surface string) (*model.Entry, bool) {

	var et model.EntryType
	switch typ {
	case "decision":
		et = model.Decision
	case "finding":
		et = model.Finding
	case "question":
		et = model.Question
	case "intent":
		et = model.Intent
	default:
		return nil, false
	}
	if strings.TrimSpace(quote) == "" || strings.TrimSpace(title) == "" {
		return nil, false
	}
	msg, foundLine := findQuote(d.Messages, quote, line)
	if msg == nil {
		return nil, false // fabricated or stitched quote → rejected (I7)
	}
	at := msg.At
	if at.IsZero() {
		return nil, false // source time is evidence; ingest time is not a substitute
	}
	e := &model.Entry{
		ID:          ids.NewEntry(at),
		Type:        et,
		Title:       truncate(title, model.MaxTitle),
		Body:        truncate(body, model.MaxBody),
		Quote:       quote,
		UtteranceBy: model.UtteranceBy(msg.Role),
		Source: model.Source{
			Kind:    model.SrcSession,
			Ref:     fmt.Sprintf("%s:%s#L%d", d.Adapter, d.File, foundLine),
			Agent:   d.Agent,
			Surface: surface,
			At:      at.UTC(),
		},
		Confidence: clamp(conf),
	}
	if msg.Role == "tool_result" {
		e.UtteranceBy = model.ByToolResult
	}
	for _, t := range tags {
		if t != "" && len(t) < 128 {
			e.Tags = append(e.Tags, t)
		}
	}
	if et == model.Finding {
		if env != nil && (env.Host != "" || env.HW != "" || env.Dataset != "") {
			e.Env = env
		}
		e.Affects = affects
	}
	if et == model.Question {
		if asks != "human" {
			asks = "any"
		}
		e.Asks = asks
	}
	return e, e.Validate() == nil
}

// findQuote locates the quote in the messages, preferring the claimed line.
func findQuote(msgs []adapters.Message, quote string, line int) (*adapters.Message, int) {
	nq := normWS(quote)
	if nq == "" {
		return nil, 0
	}
	for i := range msgs {
		if msgs[i].Line == line && strings.Contains(normWS(msgs[i].Text), nq) {
			return &msgs[i], msgs[i].Line
		}
	}
	for i := range msgs {
		if strings.Contains(normWS(msgs[i].Text), nq) {
			return &msgs[i], msgs[i].Line
		}
	}
	return nil, 0
}

func normWS(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func truncate(s string, n int) string {
	if utf8.RuneCountInString(s) <= n {
		return s
	}
	r := []rune(s)
	return string(r[:n-1]) + "…"
}

func clamp(f float64) float64 {
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}

// renderSlice formats messages with [Lnn] role markers, middle-eliding to
// SliceCap (§6.2).
func renderSlice(msgs []adapters.Message) string {
	lines := make([]string, len(msgs))
	total := 0
	for i, m := range msgs {
		lines[i] = fmt.Sprintf("[L%d] %s: %s", m.Line, m.Role, m.Text)
		total += len(lines[i]) + 1
	}
	if total <= SliceCap {
		return strings.Join(lines, "\n")
	}
	// Keep head and tail halves by bytes; elide the middle whole-message.
	half := SliceCap / 2
	var head, tail []string
	n := 0
	for _, l := range lines {
		if n+len(l)+1 > half {
			break
		}
		head = append(head, l)
		n += len(l) + 1
	}
	n = 0
	for i := len(lines) - 1; i >= len(head); i-- {
		if n+len(lines[i])+1 > half {
			break
		}
		tail = append([]string{lines[i]}, tail...)
		n += len(lines[i]) + 1
	}
	elided := len(lines) - len(head) - len(tail)
	return strings.Join(head, "\n") +
		fmt.Sprintf("\n…[%d messages elided for size]…\n", elided) +
		strings.Join(tail, "\n")
}

func readDigest(j *journal.Journal) string {
	b, err := os.ReadFile(filepath.Join(j.Dir, "digest.md"))
	if err != nil || len(b) == 0 {
		return "(journal is empty)"
	}
	return string(b)
}

// ParkSlice saves the raw slice for later inspection (I2: parked, not lost).
func ParkSlice(db *state.DB, a adapters.Adapter, file string, offset int64, reason string) error {
	d, err := a.Parse(file, offset)
	if err != nil {
		return err
	}
	return parkDelta(db, a, file, offset, reason, d)
}

func ParkSliceUntil(db *state.DB, a adapters.Adapter, file string, offset, end int64, reason string) error {
	d, err := adapters.ParseRange(a, file, offset, end)
	if err != nil {
		return err
	}
	return parkDelta(db, a, file, offset, reason, d)
}

// ParkRawRange preserves the exact adapter bytes for an unknown or broken
// record. Rendered messages are insufficient here because the unrecognized
// envelope is precisely the evidence needed to update the pinned parser.
func ParkRawRange(db *state.DB, a adapters.Adapter, file string, offset, end int64, reason string) error {
	if end <= offset {
		return fmt.Errorf("empty raw range [%d,%d)", offset, end)
	}
	src, err := os.Open(file)
	if err != nil {
		return err
	}
	defer src.Close()
	dir := filepath.Join(gitx.Home(), "parked")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	p := filepath.Join(dir, fmt.Sprintf("%d-%s.jsonl", time.Now().UnixNano(), a.ID()))
	dst, err := os.OpenFile(p, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	_, copyErr := io.CopyN(dst, io.NewSectionReader(src, offset, end-offset), end-offset)
	closeErr := dst.Close()
	if copyErr != nil {
		os.Remove(p)
		return copyErr
	}
	if closeErr != nil {
		os.Remove(p)
		return closeErr
	}
	if err := db.Park(a.ID(), file, offset, reason, p); err != nil {
		os.Remove(p)
		return err
	}
	return nil
}

func parkDelta(db *state.DB, a adapters.Adapter, file string, offset int64, reason string, d *adapters.Delta) error {
	dir := filepath.Join(gitx.Home(), "parked")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	raw := renderSlice(d.Messages)
	name := fmt.Sprintf("%d-%s.txt", time.Now().Unix(), a.ID())
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(raw), 0o644); err != nil {
		return err
	}
	return db.Park(a.ID(), file, offset, reason, p)
}
