// Package differ implements the intent×reality join (JOURNAL_SPEC §7):
// glob mapping, optional LLM link pass, finding auto-supersession, churn
// marking, and the alert set (contradiction, absence, stomp, question-aging).
// The differ writes journal events (by.who: differ) so its facts propagate;
// alerts are machine-local state. Radar, never locks (I3): nothing blocks.
package differ

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"time"

	"clew/internal/globx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/model"
	"clew/internal/poller"
	"clew/internal/state"
)

type Input struct {
	Repo     string
	Journal  *journal.Journal
	Snapshot *poller.Snapshot // may be nil
	Surface  string
	Provider llm.Provider // nil → glob-only mapping (link pass skipped)
	LinkPass bool
}

type Result struct {
	EventsAdded int
	NewAlerts   []state.Alert
	Unmapped    []string // commit shas with no entry linkage (map view)
	Tokens      int
}

// Run executes the mapping pass and status-algebra side effects, then emits
// alerts. Idempotent: events are deduped, alerts keyed.
func Run(db *state.DB, in *Input, now time.Time) (*Result, error) {
	res := &Result{}
	j := in.Journal
	by := model.By{Who: "differ", Surface: in.Surface}

	addEvent := func(kind model.EventKind, entry string, payload map[string]any, at time.Time) {
		v := &model.Event{ID: ids.NewEvent(at), Kind: kind, Entry: entry, Payload: payload, By: by, At: at.UTC()}
		if err := j.AddEvent(v); err == nil {
			res.EventsAdded++
		}
	}

	// ---- 7.1(1) glob match: new commits × entry tags → evidence ----
	commits := db.RecentCommits(in.Repo, now.Add(-14*24*time.Hour), false)
	for _, c := range commits {
		if c.Mapped {
			continue
		}
		matched := false
		for id, e := range j.Entries {
			if e.Type != model.Intent && e.Type != model.Decision {
				continue
			}
			if len(e.Tags) == 0 || !globx.AnyMatch(e.Tags, c.Files) {
				continue
			}
			if !j.HasEvent(model.EvEvidence, id, "ref", c.SHA) {
				payload := map[string]any{"kind": "commit", "ref": c.SHA, "note": c.Subject}
				if c.SessionID != "" {
					payload["session"] = c.SessionID
				}
				addEvent(model.EvEvidence, id, payload, c.At)
			}
			matched = true
		}
		// Churn: commits touching a finding's affects → churn event (suspect).
		for id, e := range j.Entries {
			if e.Type != model.Finding || len(e.Affects) == 0 {
				continue
			}
			if globx.AnyMatch(e.Affects, c.Files) && !j.HasEvent(model.EvEvidence, id, "ref", "churn:"+c.SHA) {
				addEvent(model.EvEvidence, id, map[string]any{"kind": "churn", "ref": "churn:" + c.SHA, "note": c.Subject}, c.At)
			}
		}
		if matched {
			db.MarkCommitMapped(in.Repo, c.SHA)
		} else {
			res.Unmapped = append(res.Unmapped, c.SHA)
		}
	}
	// Plain commit subjects can prove recently requested work without spending
	// a model call. Require two distinctive shared words and only consider
	// decisions/intents that existed before the commit; everything ambiguous
	// remains for the optional link pass.
	linkedBySubject := subjectLinkPass(db, in, commits, addEvent)
	linkedBySubject = append(linkedBySubject, repairCoordinationLinks(db, in, commits, addEvent)...)
	res.Unmapped = subtract(res.Unmapped, linkedBySubject)

	// ---- 7.1(2) LLM link pass: unmatched commits × unevidenced intents ----
	if in.LinkPass && in.Provider != nil && len(res.Unmapped) > 0 {
		linked, tokens := linkPass(db, in, res.Unmapped, addEvent, now)
		res.Tokens = tokens
		res.Unmapped = subtract(res.Unmapped, linked)
	}
	// Persist unmapped list for `map` (displayed, not guessed — §7.1.2).
	um, _ := json.Marshal(res.Unmapped)
	db.Set("unmapped:"+in.Repo, string(um))

	// ---- finding auto-supersession: newer measurement, same env+topic ----
	st := journal.Compute(j, now)
	autoSupersede(j, st, addEvent, now)

	// ---- alerts ----
	st = journal.Compute(j, now) // recompute after event writes
	if err := emitAlerts(db, in, st, res, now); err != nil {
		return nil, err
	}
	return res, nil
}

func subjectLinkPass(db *state.DB, in *Input, commits []state.Commit, addEvent func(model.EventKind, string, map[string]any, time.Time)) []string {
	var linked []string
	for _, commit := range commits {
		if strings.HasPrefix(strings.ToLower(strings.TrimSpace(commit.Subject)), "journal:") {
			continue
		}
		clauses := evidenceClauses(commit.Subject)
		matched := false
		for id, entry := range in.Journal.Entries {
			if !workEvidenceEntry(entry) || entry.Created().After(commit.At) || hasWorkEvidence(in.Journal, id) {
				continue
			}
			entryWords := evidenceWords(entry.Title + " " + entry.Body)
			if !clauseMatches(clauses, entryWords) {
				continue
			}
			addEvent(model.EvEvidence, id, map[string]any{
				"kind": "commit", "ref": commit.SHA, "note": commit.Subject,
				"via": "subject-match",
			}, commit.At)
			matched = true
		}
		if matched {
			db.MarkCommitMapped(in.Repo, commit.SHA)
			linked = append(linked, commit.SHA)
		}
	}
	return linked
}

// repairCoordinationLinks replaces no evidence and rewrites no history. It
// handles the narrow legacy case where a high-confidence link-pass selected a
// clew/journal commit: exactly one real code commit must fall between the
// entry's source time and that false event. Ambiguity leaves the work open.
func repairCoordinationLinks(db *state.DB, in *Input, commits []state.Commit, addEvent func(model.EventKind, string, map[string]any, time.Time)) []string {
	var linked []string
	for id, entry := range in.Journal.Entries {
		if !workEvidenceEntry(entry) || hasWorkEvidence(in.Journal, id) {
			continue
		}
		for _, event := range in.Journal.EventsFor(id) {
			if journal.CountsAsRealityEvidence(event) || event.Kind != model.EvEvidence || event.PStr("via") != "link-pass" || eventConfidence(event) < 0.8 {
				continue
			}
			var candidates []state.Commit
			for _, commit := range commits {
				if strings.HasPrefix(strings.ToLower(strings.TrimSpace(commit.Subject)), "journal:") ||
					commit.At.Before(entry.Created()) || commit.At.After(event.At) {
					continue
				}
				candidates = append(candidates, commit)
			}
			if len(candidates) != 1 {
				continue
			}
			commit := candidates[0]
			addEvent(model.EvEvidence, id, map[string]any{
				"kind": "commit", "ref": commit.SHA, "note": commit.Subject,
				"via": "coordination-repair",
			}, commit.At)
			db.MarkCommitMapped(in.Repo, commit.SHA)
			linked = append(linked, commit.SHA)
			break
		}
	}
	return linked
}

func workEvidenceEntry(entry *model.Entry) bool {
	if entry == nil {
		return false
	}
	if entry.Type == model.Decision || entry.Type == model.Intent {
		return true
	}
	if entry.Type != model.Finding {
		return false
	}
	text := strings.ToLower(entry.Title + " " + entry.Body)
	return strings.Contains(text, "finished") && (strings.Contains(text, "stale") || strings.Contains(text, "uncommitted"))
}

func eventConfidence(event *model.Event) float64 {
	if event == nil {
		return 0
	}
	switch value := event.Payload["confidence"].(type) {
	case float64:
		return value
	case float32:
		return float64(value)
	case int:
		return float64(value)
	case int64:
		return float64(value)
	}
	return 0
}

func evidenceClauses(subject string) []map[string]bool {
	parts := strings.Split(subject, ";")
	out := make([]map[string]bool, 0, len(parts))
	for _, part := range parts {
		if words := evidenceWords(part); len(words) > 0 {
			out = append(out, words)
		}
	}
	return out
}

func clauseMatches(clauses []map[string]bool, entryWords map[string]bool) bool {
	for _, clause := range clauses {
		if sharedWordCount(clause, entryWords) >= 2 {
			return true
		}
	}
	return false
}

func hasWorkEvidence(j *journal.Journal, id string) bool {
	for _, event := range j.EventsFor(id) {
		if journal.CountsAsRealityEvidence(event) && (event.PStr("kind") == "commit" || event.PStr("kind") == "completion") {
			return true
		}
	}
	return false
}

func evidenceWords(text string) map[string]bool {
	var b strings.Builder
	for _, r := range strings.ToLower(text) {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' {
			b.WriteRune(r)
		} else {
			b.WriteByte(' ')
		}
	}
	stop := map[string]bool{
		"a": true, "an": true, "and": true, "as": true, "at": true, "by": true,
		"for": true, "from": true, "in": true, "is": true, "it": true, "never": true,
		"of": true, "on": true, "or": true, "the": true, "to": true, "under": true,
		"while": true, "with": true,
		// Common commit scaffolding is not evidence. Matches must come from the
		// subject matter inside one clause, not generic project vocabulary.
		"add": true, "agent": true, "entry": true, "fix": true, "human": true,
		"journal": true, "keep": true, "make": true, "memory": true, "screen": true,
		"surface": true, "system": true, "update": true, "work": true,
	}
	out := make(map[string]bool)
	for _, word := range strings.Fields(b.String()) {
		if stop[word] {
			continue
		}
		for _, suffix := range []string{"ing", "ied", "ed", "s"} {
			if strings.HasSuffix(word, suffix) && len(word) > len(suffix)+3 {
				word = strings.TrimSuffix(word, suffix)
				if suffix == "ied" {
					word += "y"
				}
				break
			}
		}
		out[word] = true
	}
	return out
}

func sharedWordCount(a, b map[string]bool) int {
	n := 0
	for word := range a {
		if b[word] {
			n++
		}
	}
	return n
}

// autoSupersede pairs findings sharing a tag with equal env; the newer
// supersedes the older (§7.1.3). Changed env → both stay current.
func autoSupersede(j *journal.Journal, st map[string]*journal.Computed, addEvent func(model.EventKind, string, map[string]any, time.Time), now time.Time) {
	var findings []*model.Entry
	for id, e := range j.Entries {
		if e.Type == model.Finding && st[id] != nil && st[id].Status == journal.StCurrent {
			findings = append(findings, e)
		}
	}
	sort.Slice(findings, func(a, b int) bool { return findings[a].ID < findings[b].ID })
	for i := 0; i < len(findings); i++ {
		for k := i + 1; k < len(findings); k++ {
			old, nw := findings[i], findings[k]
			// Same topic (shared tag) and same env → newer supersedes older.
			// Changed env → both stay current (§7.1.3).
			if !tagsShared(old.Tags, nw.Tags) || !old.Env.Equal(nw.Env) {
				continue
			}
			if !j.HasEvent(model.EvSupersede, old.ID, "by", nw.ID) {
				addEvent(model.EvSupersede, old.ID, map[string]any{"by": nw.ID, "why": "newer measurement, same env"}, now)
			}
		}
	}
}

var differAlertKinds = []string{"contradiction", "absence", "aging", "suspect", "stomp", "overlap"}

func emitAlerts(db *state.DB, in *Input, st map[string]*journal.Computed, res *Result, now time.Time) error {
	j := in.Journal
	active := []state.Alert{}
	emit := func(kind, identity, entryIDs, body, resolvedWhen string, blocking bool) {
		fingerprint := shortHash(in.Repo + "\x00" + kind + "\x00" + identity)
		active = append(active, state.Alert{
			Key:          kind + ":" + fingerprint,
			RepoPath:     in.Repo,
			Kind:         kind,
			Body:         body,
			EntryIDs:     entryIDs,
			WithdrawWhen: kind + ":" + resolvedWhen + ":" + fingerprint,
			Blocking:     blocking,
		})
	}

	seenPair := map[string]bool{}
	for id, c := range st {
		e := j.Entries[id]
		switch c.Status {
		case journal.StPossibleContradiction:
			for _, other := range c.Contradicts {
				pk := pairID(id, other)
				if seenPair[pk] {
					continue
				}
				seenPair[pk] = true
				oe := j.Entries[other]
				body := fmt.Sprintf("possible contradiction: %q vs %q — quotes: “%s” / “%s” (needs human: only you can rule it a real contradiction)",
					e.Title, oe.Title, clip(e.Quote, 80), clip(oe.Quote, 80))
				emit("contradiction", pk, pk, body, "status-resolved", true)
			}
		case journal.StAbsent:
			body := fmt.Sprintf("absence: intent %q has zero evidence while %d+ sibling intents progressed — was it forgotten? (blocks: only a human can decide drop vs revive)",
				e.Title, journal.AbsenceSiblings)
			emit("absence", id, id, body, "status-resolved", true)
		case journal.StOpen:
			if e.Asks == "human" && now.Sub(e.Created()) > journal.QuestionAging {
				body := fmt.Sprintf("question aging: %q open %dd, addressed to you (blocks: agents cannot answer it)",
					e.Title, int(now.Sub(e.Created()).Hours()/24))
				emit("aging", id, id, body, "status-resolved", true)
			}
		case journal.StSuspect:
			body := fmt.Sprintf("suspect finding: %q — affected paths churned since it was measured; re-measure or supersede", e.Title)
			emit("suspect", id, id, body, "status-resolved", false)
		}
	}

	// Overlap radar (§5.2): footprint intersection between live sessions on
	// this repo; same file while dirty → stomp (docket), else map annotation.
	sessions := db.LiveSessions(in.Repo, 30*time.Minute)
	dirty := map[string]bool{}
	if in.Snapshot != nil {
		for _, f := range in.Snapshot.Dirty {
			dirty[f] = true
		}
	}
	for i := 0; i < len(sessions); i++ {
		for k := i + 1; k < len(sessions); k++ {
			fa, fb := db.Footprints(sessions[i].ID), db.Footprints(sessions[k].ID)
			for _, f := range intersect(fa, fb) {
				rel := relPath(in.Repo, f)
				identity := pairID(sessions[i].ID, sessions[k].ID) + "\x00" + rel
				if dirty[rel] {
					body := fmt.Sprintf("stomp risk: sessions %s(%s) and %s(%s) both touched %s and it is dirty — lost-work scenario",
						sessions[i].Agent, sessions[i].Surface, sessions[k].Agent, sessions[k].Surface, rel)
					emit("stomp", identity, rel, body, "path-clean-or-sessions-diverged", true)
				} else {
					emit("overlap", identity, rel, fmt.Sprintf("overlap: two sessions touched %s on this repo+branch", rel), "path-dirty-or-sessions-diverged", false)
				}
			}
		}
	}
	created, err := db.ReconcileAlerts(in.Repo, differAlertKinds, active)
	if err != nil {
		return fmt.Errorf("reconcile differ alerts: %w", err)
	}
	res.NewAlerts = append(res.NewAlerts, created...)
	return nil
}

// linkPass batches unmapped commits × unevidenced intents through the
// provider; ≥0.8 auto-attach (§7.1.2).
func linkPass(db *state.DB, in *Input, unmapped []string, addEvent func(model.EventKind, string, map[string]any, time.Time), now time.Time) ([]string, int) {
	j := in.Journal
	st := journal.Compute(j, now)
	var intents []*model.Entry
	for id, e := range j.Entries {
		if e.Type == model.Intent && st[id].Evidence == 0 && journal.Live(st[id].Status) {
			intents = append(intents, e)
		}
	}
	if len(intents) == 0 {
		return nil, 0
	}
	commits := map[string]state.Commit{}
	for _, c := range db.RecentCommits(in.Repo, now.Add(-14*24*time.Hour), true) {
		commits[c.SHA] = c
	}
	var b strings.Builder
	b.WriteString("Map commits to intents. Output STRICT JSON {\"links\":[{\"commit\":\"sha\",\"intent\":\"id\",\"confidence\":0.0}]} — only pairs you are confident about; no prose.\n\nINTENTS:\n")
	for _, e := range intents {
		fmt.Fprintf(&b, "- %s: %s — %s\n", e.ID, e.Title, clip(e.Body, 120))
	}
	b.WriteString("\nCOMMITS:\n")
	for _, sha := range unmapped {
		c, ok := commits[sha]
		if !ok {
			continue
		}
		fmt.Fprintf(&b, "- %s: %s (files: %s)\n", sha[:min(10, len(sha))], c.Subject, clip(strings.Join(c.Files, ", "), 200))
	}
	res, err := in.Provider.Call(b.String())
	if err != nil {
		return nil, 0
	}
	raw, ok := llm.ExtractJSON(res.Text)
	if !ok {
		return nil, res.Tokens
	}
	var out struct {
		Links []struct {
			Commit     string  `json:"commit"`
			Intent     string  `json:"intent"`
			Confidence float64 `json:"confidence"`
		} `json:"links"`
	}
	if json.Unmarshal([]byte(raw), &out) != nil {
		return nil, res.Tokens
	}
	var linked []string
	for _, l := range out.Links {
		if l.Confidence < 0.8 {
			continue // below threshold → stays unmapped (displayed, not guessed)
		}
		full := l.Commit
		for sha := range commits {
			if strings.HasPrefix(sha, l.Commit) {
				full = sha
			}
		}
		c, ok := commits[full]
		if !ok {
			continue
		}
		if _, exists := j.Entries[l.Intent]; !exists {
			continue
		}
		if !j.HasEvent(model.EvEvidence, l.Intent, "ref", full) {
			addEvent(model.EvEvidence, l.Intent, map[string]any{
				"kind": "commit", "ref": full, "note": c.Subject,
				"via": "link-pass", "confidence": l.Confidence,
			}, c.At)
		}
		db.MarkCommitMapped(in.Repo, full)
		linked = append(linked, full)
	}
	return linked, res.Tokens
}

func tagsShared(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x != "" && x == y {
				return true
			}
		}
	}
	return false
}

func intersect(a, b []string) []string {
	set := map[string]bool{}
	for _, x := range a {
		set[x] = true
	}
	var out []string
	for _, y := range b {
		if set[y] {
			out = append(out, y)
		}
	}
	return out
}

func subtract(a, remove []string) []string {
	rm := map[string]bool{}
	for _, r := range remove {
		rm[r] = true
	}
	var out []string
	for _, x := range a {
		if !rm[x] {
			out = append(out, x)
		}
	}
	return out
}

func relPath(repo, p string) string {
	if strings.HasPrefix(p, repo+"/") {
		return p[len(repo)+1:]
	}
	return p
}

func pairID(a, b string) string {
	if a > b {
		a, b = b, a
	}
	return a + "+" + b
}

func clip(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func shortHash(s string) string {
	h := sha1.Sum([]byte(s))
	return hex.EncodeToString(h[:])[:10]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
