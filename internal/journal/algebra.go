package journal

import (
	"sort"
	"strings"
	"time"

	"clew/internal/model"
)

// Status values are computed by each watcher from entries + events and never
// persisted to the branch (§3.2). See §3.1 for the per-type status sets and
// §7.1(3) for the algebra.
type Status string

const (
	// decision
	StActive                Status = "active"
	StSuperseded            Status = "superseded"
	StPossibleContradiction Status = "possible-contradiction"
	StContradicted          Status = "contradicted" // human-confirmed only
	// finding
	StCurrent Status = "current"
	StSuspect Status = "suspect"
	// question
	StOpen     Status = "open"
	StAnswered Status = "answered"
	StExpired  Status = "expired"
	// intent
	StProposed Status = "proposed"
	StInFlight Status = "in_flight"
	StDone     Status = "done"
	StAbsent   Status = "absent"
	StDropped  Status = "dropped"
)

const (
	InFlightWindow  = 7 * 24 * time.Hour  // evidence in last 7d → in_flight
	QuestionAging   = 7 * 24 * time.Hour  // human-addressed question alert (§7.1.4)
	QuestionExpiry  = 45 * 24 * time.Hour // §3.3
	AbsenceSiblings = 5                   // K (§7.1.3)
	EligibleMinConf = 0.8                 // session-sourced absence eligibility
)

// Computed is everything the algebra derives for one entry.
type Computed struct {
	Status             Status
	Confidence         float64 // effective: human confirm raises to 1.0
	Evidence           int     // count of evidence events
	LastActivity       time.Time
	SupersededBy       string
	AnsweredBy         string
	LineageStatus      Status    // explicit human carry marker for a seed grave
	LineageStatusAt    time.Time // new successor evidence may revive carried absence
	LineageStatusEvent string
	Contradicts        []string // ids of pair partners (possible or confirmed)
	Eligible           bool     // absence-eligibility (§7.1.3 guard)
	Tainted            bool     // quote originated in tool_result (§6.5)
	Withheld           bool     // imperative-to-agent pattern, pending human confirm (§6.5.3)
}

// Compute runs the status algebra over the whole journal at instant now.
func Compute(j *Journal, now time.Time) map[string]*Computed {
	out := make(map[string]*Computed, len(j.Entries))

	// supersededBy index: supersede/reject events + entry.Supersedes chains.
	supBy := map[string]string{}
	for _, e := range j.Entries {
		if e.Supersedes != "" {
			supBy[e.Supersedes] = e.ID
		}
	}
	for _, v := range j.Events {
		switch v.Kind {
		case model.EvSupersede:
			supBy[v.Entry] = v.PStr("by")
		case model.EvReject:
			if _, ok := supBy[v.Entry]; !ok {
				supBy[v.Entry] = "human:" + v.ID
			}
		}
	}

	for id, e := range j.Entries {
		c := &Computed{Confidence: e.Confidence, LastActivity: e.Created()}
		c.Tainted = e.UtteranceBy == model.ByToolResult
		c.Withheld = Imperative(e)
		for _, v := range j.EventsFor(id) {
			if v.At.After(c.LastActivity) {
				c.LastActivity = v.At
			}
			switch v.Kind {
			case model.EvEvidence:
				if IsRealityEvidence(j, id, v) {
					c.Evidence++
				}
			case model.EvConfirm:
				if v.By.Who == "human" {
					c.Confidence = 1.0
					c.Withheld = false
				}
			case model.EvAnswer:
				c.AnsweredBy = v.PStr("by")
			case model.EvDisposition:
				status := Status(v.PStr("lineage_status"))
				if v.By.Who == "human" && v.PStr("lineage_from_repository") != "" &&
					v.PStr("lineage_seed_revision") != "" && lineageStatusCompatible(e.Type, status) &&
					(c.LineageStatusAt.IsZero() || !v.At.Before(c.LineageStatusAt)) {
					c.LineageStatus = status
					c.LineageStatusAt = v.At
					c.LineageStatusEvent = v.ID
				}
			}
		}
		c.SupersededBy = supBy[id]
		if c.SupersededBy == "" && c.LineageStatus == StSuperseded {
			c.SupersededBy = "lineage:" + c.LineageStatusEvent
		}
		out[id] = c
	}

	// Absence eligibility (§7.1.3 guard): human-confirmed, or session-sourced
	// with confidence ≥ 0.8, or human/carried-sourced (carry is a human act).
	for id, e := range j.Entries {
		c := out[id]
		humanConfirmed := c.Confidence >= 1.0 && hasHumanConfirm(j, id)
		switch {
		case humanConfirmed:
			c.Eligible = true
		case e.Source.Kind == model.SrcSession && e.Confidence >= EligibleMinConf:
			c.Eligible = true
		case e.Source.Kind == model.SrcHuman || e.Source.Kind == model.SrcCarried:
			c.Eligible = true
		}
	}

	// Per-type status resolution.
	for id, e := range j.Entries {
		c := out[id]
		switch e.Type {
		case model.Question:
			switch {
			case c.AnsweredBy != "":
				c.Status = StAnswered
			case c.LineageStatus == StExpired || c.SupersededBy != "" || now.Sub(c.LastActivity) > QuestionExpiry:
				c.Status = StExpired
			default:
				c.Status = StOpen
			}
		case model.Finding:
			switch {
			case c.LineageStatus == StSuperseded || c.SupersededBy != "":
				c.Status = StSuperseded
			case len(e.Affects) > 0 && hasChurnAfter(j, id, e.Created()):
				c.Status = StSuspect
			default:
				c.Status = StCurrent
			}
		case model.Intent:
			c.Status = intentStatus(j, out, e, c, now)
		case model.Decision:
			if c.LineageStatus == StSuperseded || c.SupersededBy != "" {
				c.Status = StSuperseded
			} else {
				c.Status = StActive // contradiction pass below
			}
		}
	}

	// Decision contradiction pass (§7.1.3): a decision sharing a tag with
	// another active decision → both surfaced as a possible-contradiction
	// pair. Only a human confirm event sets contradicted.
	decisions := []string{}
	for id, e := range j.Entries {
		if e.Type == model.Decision && out[id].Status == StActive {
			decisions = append(decisions, id)
		}
	}
	sort.Strings(decisions)
	confirmedPairs := map[[2]string]bool{}
	for _, v := range j.Events {
		if v.Kind == model.EvConfirm && v.By.Who == "human" {
			if o := v.PStr("contradicts"); o != "" {
				confirmedPairs[pairKey(v.Entry, o)] = true
			}
		}
	}
	for i := 0; i < len(decisions); i++ {
		for k := i + 1; k < len(decisions); k++ {
			a, b := decisions[i], decisions[k]
			if !tagsOverlap(j.Entries[a].Tags, j.Entries[b].Tags) {
				continue
			}
			ca, cb := out[a], out[b]
			ca.Contradicts = append(ca.Contradicts, b)
			cb.Contradicts = append(cb.Contradicts, a)
			st := StPossibleContradiction
			if confirmedPairs[pairKey(a, b)] {
				st = StContradicted
			}
			// contradicted (confirmed) outranks possible.
			if ca.Status != StContradicted {
				ca.Status = st
			}
			if cb.Status != StContradicted {
				cb.Status = st
			}
		}
	}
	return out
}

func intentStatus(j *Journal, all map[string]*Computed, e *model.Entry, c *Computed, now time.Time) Status {
	if c.SupersededBy != "" {
		return StDropped
	}
	if c.LineageStatus == StDropped {
		return StDropped
	}
	hasSuccessorEvidence := false
	for _, v := range j.EventsFor(e.ID) {
		if !c.LineageStatusAt.IsZero() && !v.At.After(c.LineageStatusAt) {
			continue
		}
		if v.Kind == model.EvConfirm && v.PBool("done") {
			return StDone
		}
		if IsRealityEvidence(j, e.ID, v) {
			hasSuccessorEvidence = true
			if v.PStr("kind") == "completion" {
				return StDone
			}
		}
	}
	recent := false
	for _, v := range j.EventsFor(e.ID) {
		if IsRealityEvidence(j, e.ID, v) &&
			(c.LineageStatusAt.IsZero() || v.At.After(c.LineageStatusAt)) &&
			now.Sub(v.At) <= InFlightWindow {
			recent = true
		}
	}
	if recent {
		return StInFlight
	}
	if c.LineageStatus == StAbsent && !hasSuccessorEvidence {
		return StAbsent
	}
	// Absence rule (§7.1.3): zero evidence ever AND ≥K eligible sibling
	// intents gained evidence since this entry was created. Relative to
	// project activity, not wall-clock.
	if c.Evidence == 0 && c.Eligible {
		created := e.Created()
		siblings := 0
		for oid, oe := range j.Entries {
			if oid == e.ID || oe.Type != model.Intent || !all[oid].Eligible {
				continue
			}
			for _, v := range j.EventsFor(oid) {
				if IsRealityEvidence(j, oid, v) && v.At.After(created) {
					siblings++
					break
				}
			}
		}
		if siblings >= AbsenceSiblings {
			return StAbsent
		}
	}
	return StProposed
}

// CountsAsRealityEvidence rejects coordination-only commits. Older watcher
// builds could observe cached clew/journal commits as if they were code
// reality and attach subject-match events to the decision the commit merely
// recorded. Those immutable events stay in history, but never move an intent
// or settle a merge line.
func CountsAsRealityEvidence(event *model.Event) bool {
	if event == nil || event.Kind != model.EvEvidence {
		return false
	}
	return event.PStr("kind") != "commit" || !strings.HasPrefix(strings.ToLower(strings.TrimSpace(event.PStr("note"))), "journal:")
}

// IsRealityEvidence also applies append-only corrections written when a
// previously accepted subject match no longer passes the current conservative
// matcher. The original event remains available as provenance.
func IsRealityEvidence(j *Journal, entryID string, event *model.Event) bool {
	if !CountsAsRealityEvidence(event) || j == nil {
		return false
	}
	ref := event.PStr("ref")
	return ref == "" || !j.HasEvent(model.EvEvidenceWithdrawn, entryID, "ref", ref)
}

func lineageStatusCompatible(typ model.EntryType, status Status) bool {
	switch typ {
	case model.Decision, model.Finding:
		return status == StSuperseded
	case model.Question:
		return status == StExpired
	case model.Intent:
		return status == StAbsent || status == StDropped
	}
	return false
}

func hasHumanConfirm(j *Journal, id string) bool {
	for _, v := range j.EventsFor(id) {
		if v.Kind == model.EvConfirm && v.By.Who == "human" {
			return true
		}
	}
	return false
}

// hasChurnAfter reports churn evidence (differ-written, payload kind=churn)
// newer than the finding's measurement (§7.1.3 suspect rule).
func hasChurnAfter(j *Journal, id string, t time.Time) bool {
	for _, v := range j.EventsFor(id) {
		if v.Kind == model.EvEvidence && v.PStr("kind") == "churn" && v.At.After(t) {
			return true
		}
	}
	return false
}

func tagsOverlap(a, b []string) bool {
	for _, x := range a {
		for _, y := range b {
			if x != "" && x == y {
				return true
			}
		}
	}
	return false
}

func pairKey(a, b string) [2]string {
	if a > b {
		a, b = b, a
	}
	return [2]string{a, b}
}

// Live reports whether an entry should render in rollup/context (§3.3:
// non-superseded, non-expired; dropped intents are superseded-by-human).
func Live(st Status) bool {
	switch st {
	case StSuperseded, StExpired, StDropped:
		return false
	}
	return true
}
