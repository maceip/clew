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

	"restart/internal/globx"
	"restart/internal/ids"
	"restart/internal/journal"
	"restart/internal/llm"
	"restart/internal/model"
	"restart/internal/poller"
	"restart/internal/state"
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
	commits := db.RecentCommits(in.Repo, now.Add(-14*24*time.Hour), true)
	for _, c := range commits {
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
	emitAlerts(db, in, st, res, now)
	return res, nil
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

func emitAlerts(db *state.DB, in *Input, st map[string]*journal.Computed, res *Result, now time.Time) {
	j := in.Journal
	upsert := func(kind, entryIDs, body string, blocking bool) {
		key := kind + ":" + shortHash(in.Repo+entryIDs+body)
		a := state.Alert{Key: key, RepoPath: in.Repo, Kind: kind, Body: body, EntryIDs: entryIDs, Blocking: blocking}
		if isNew, _ := db.UpsertAlert(a); isNew {
			res.NewAlerts = append(res.NewAlerts, a)
		}
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
				upsert("contradiction", pk, body, true)
			}
		case journal.StAbsent:
			body := fmt.Sprintf("absence: intent %q has zero evidence while %d+ sibling intents progressed — was it forgotten? (blocks: only a human can decide drop vs revive)",
				e.Title, journal.AbsenceSiblings)
			upsert("absence", id, body, true)
		case journal.StOpen:
			if e.Asks == "human" && now.Sub(e.Created()) > journal.QuestionAging {
				body := fmt.Sprintf("question aging: %q open %dd, addressed to you (blocks: agents cannot answer it)",
					e.Title, int(now.Sub(e.Created()).Hours()/24))
				upsert("aging", id, body, true)
			}
		case journal.StSuspect:
			body := fmt.Sprintf("suspect finding: %q — affected paths churned since it was measured; re-measure or supersede", e.Title)
			upsert("suspect", id, body, false)
		}
	}

	// Overlap radar (§5.2): footprint intersection between live sessions on
	// this repo; same file while dirty → stomp (inbox), else map annotation.
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
				if dirty[rel] {
					body := fmt.Sprintf("stomp risk: sessions %s(%s) and %s(%s) both touched %s and it is dirty — lost-work scenario",
						sessions[i].Agent, sessions[i].Surface, sessions[k].Agent, sessions[k].Surface, rel)
					upsert("stomp", rel, body, true)
				} else {
					upsert("overlap", rel, fmt.Sprintf("overlap: two sessions touched %s on this repo+branch", rel), false)
				}
			}
		}
	}
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
