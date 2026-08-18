// Package checkin builds the two calm, glanceable finish-boundary surfaces.
// It contains presentation and selection only. Applying or building always
// produces an explicit agent handoff; this package never changes project code.
package checkin

import (
	"fmt"
	"io"
	"sort"
	"strings"
	"time"

	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/state"
)

const MaxItems = 7

type Screen string

const (
	KnowledgeMerge Screen = "knowledge-merge"
	IntentGap      Screen = "intent-gap"
)

type Item struct {
	ID       string
	Line     string
	Entry    *model.Entry
	Priority int
}

type View struct {
	Screen   Screen
	Repo     string
	Items    []Item
	Deferred int
	Issues   int // source records that did not pass a clean read
	Checks   int // watcher/sync/materialization checks currently failing
}

// BuildMerge returns the newest live changes that the human has not already
// applied or deferred. handled is keyed by entry id with values "applied" or
// "deferred". Questions are already served by the docket, so this boundary
// carries decisions, findings, and intended work only.
func BuildMerge(j *journal.Journal, handled map[string]string) View {
	view := View{Screen: KnowledgeMerge, Issues: issueCount(j)}
	if j == nil {
		view.Issues = 1
		return view
	}
	computed := journal.Compute(j, latestTime(j))
	var items []Item
	for id, entry := range j.Entries {
		switch handled[id] {
		case "deferred":
			view.Deferred++
			continue
		case "applied":
			continue
		}
		if entry.Type == model.Question || !mergeLive(computed[id]) {
			continue
		}
		items = append(items, Item{ID: id, Line: CalmLine(entry), Entry: entry})
	}
	sort.Slice(items, func(i, k int) bool { return items[i].ID > items[k].ID })
	if len(items) > MaxItems {
		items = items[:MaxItems]
	}
	view.Items = items
	return view
}

// BuildGap returns intended work with no evidence in reality. An active
// failure may also prove that a newer decision is not real yet; the current
// budget failure is the first such case. There is intentionally no build-all.
func BuildGap(j *journal.Journal, alerts []state.Alert) View {
	view := View{Screen: IntentGap, Issues: issueCount(j)}
	if j == nil {
		view.Issues = 1
		return view
	}
	computed := journal.Compute(j, latestTime(j))
	var items []Item
	for id, entry := range j.Entries {
		c := computed[id]
		if c == nil {
			continue
		}
		switch entry.Type {
		case model.Intent:
			if c.Evidence != 0 || (c.Status != journal.StProposed && c.Status != journal.StAbsent) {
				continue
			}
			priority := 2
			if c.Status == journal.StAbsent {
				priority = 1
			}
			items = append(items, Item{ID: id, Line: CalmLine(entry), Entry: entry, Priority: priority})
		case model.Decision:
			if c.Status != journal.StActive && c.Status != journal.StPossibleContradiction && c.Status != journal.StContradicted {
				continue
			}
			if decisionHasLiveFailure(entry, alerts) {
				items = append(items, Item{ID: id, Line: CalmLine(entry), Entry: entry, Priority: 0})
			}
		}
	}
	sort.Slice(items, func(i, k int) bool {
		if items[i].Priority != items[k].Priority {
			return items[i].Priority < items[k].Priority
		}
		return items[i].ID > items[k].ID
	})
	if len(items) > MaxItems {
		items = items[:MaxItems]
	}
	view.Items = items
	return view
}

func mergeLive(c *journal.Computed) bool {
	if c == nil {
		return false
	}
	switch c.Status {
	case journal.StSuperseded, journal.StAnswered, journal.StExpired, journal.StDone, journal.StDropped:
		return false
	default:
		return true
	}
}

func issueCount(j *journal.Journal) int {
	if j == nil {
		return 0
	}
	return len(j.LoadErrors) + len(j.DisplayRecoveries)
}

// latestTime is deliberately derived from the journal rather than wall time,
// keeping tests and the gap classification stable. All current no-evidence
// intents remain proposed/absent under either clock.
func latestTime(j *journal.Journal) (latest time.Time) {
	for _, entry := range j.Entries {
		if at := entry.Created(); at.After(latest) {
			latest = at
		}
	}
	return latest
}

func decisionHasLiveFailure(entry *model.Entry, alerts []state.Alert) bool {
	text := strings.ToLower(entry.Title + " " + entry.Body)
	for _, alert := range alerts {
		switch alert.Kind {
		case "budget":
			if containsAny(text, "hard floor", "spend floor", "budget floor", "must never deadlock", "frugality replaced") {
				return true
			}
		case "adapter", "adapter-failure", "unknown-format":
			if containsAny(text, "adapter", "format", "listen", "witness") {
				return true
			}
		}
	}
	return false
}

func containsAny(s string, words ...string) bool {
	for _, word := range words {
		if strings.Contains(s, word) {
			return true
		}
	}
	return false
}

// CalmLine translates the known machine-heavy source titles into the five-
// promise register. The source entry remains linked for exact detail.
func CalmLine(entry *model.Entry) string {
	if entry == nil {
		return "Could not read this change"
	}
	title := strings.ToLower(entry.Title)
	switch {
	case strings.Contains(title, "codex finished i13 stale") || strings.Contains(title, "law wording"):
		return "Rename “law/laws” wherever a human reads I13; keep the hard wording for agents"
	case strings.Contains(title, "knowledge merge at finish"):
		return "At finish, show a short list of new work and decisions to apply or defer"
	case strings.Contains(title, "merge lines must pass"):
		return "Make each change understandable a day later; offer apply, explain, or defer"
	case strings.Contains(title, "explain is live"):
		return "Let the agent already here explain the selected change from the owner’s words"
	case strings.Contains(title, "silence is the signal"):
		return "Show quiet only after checking that nothing new landed anywhere"
	case strings.Contains(title, "second tab: the intent gap"):
		return "Show everything we meant to build but have not made real"
	case strings.Contains(title, "held: a restart tab"):
		return "Set aside a way to choose changes for the next start-over"
	case strings.Contains(title, "freshness ladder"):
		return "Wire every agent to learn new decisions before the next human message"
	case strings.Contains(title, "i9 frugality replaced") || strings.Contains(title, "budget-deadlock"):
		return "Keep listening with a spend floor above one full request, under the owner’s ceiling"
	case strings.Contains(title, "birth detection"):
		return "Let a new project begin with the owner’s standing decisions and nothing from an old project"
	case strings.Contains(title, "phone reads the glance"):
		return "Let the phone show what changed and ask only when a decision needs you"
	case strings.Contains(title, "laptop agents fully sensed"):
		return "Let laptop agents notice all work without asking you to record it"
	case strings.Contains(title, "pr-only cloud agents"):
		return "Let cloud agents that can only open a pull request return what they learned"
	case strings.Contains(title, "repo-write cloud agents"):
		return "Let cloud agents that can write the project return what they learned"
	}
	return strings.TrimSpace(entry.Title)
}

func Render(w io.Writer, view View) error {
	if w == nil {
		return fmt.Errorf("render check-in: nil writer")
	}
	repo := strings.TrimSpace(view.Repo)
	if repo == "" {
		repo = "this project"
	}
	var b strings.Builder
	switch view.Screen {
	case KnowledgeMerge:
		fmt.Fprintf(&b, "KNOWLEDGE MERGE — it remembers what we decide — %s\n", repo)
	case IntentGap:
		fmt.Fprintf(&b, "INTENT GAP — you can look up and see — %s\n", repo)
	default:
		return fmt.Errorf("render check-in: unknown screen %q", view.Screen)
	}
	if view.Issues > 0 {
		fmt.Fprintf(&b, "! Could not fully check: %d saved item%s could not be read cleanly.\n", view.Issues, plural(view.Issues))
	}
	if view.Checks > 0 {
		verb := "need"
		if view.Checks == 1 {
			verb = "needs"
		}
		fmt.Fprintf(&b, "! Could not fully check: %d live check%s %s attention.\n", view.Checks, plural(view.Checks), verb)
	}
	if len(view.Items) == 0 {
		if view.Issues > 0 || view.Checks > 0 {
			b.WriteString("No trustworthy empty result is available.\n")
		} else if view.Screen == KnowledgeMerge {
			b.WriteString("Nothing new.\n")
		} else {
			b.WriteString("Everything here is real.\n")
		}
		_, err := io.WriteString(w, b.String())
		return err
	}
	for i, item := range view.Items {
		fmt.Fprintf(&b, "%d. %s [%s]", i+1, item.Line, item.ID)
		if view.Screen == KnowledgeMerge {
			b.WriteString("  apply · explain · defer\n")
		} else {
			b.WriteString("  build · explain · retire\n")
		}
	}
	if view.Screen == KnowledgeMerge {
		b.WriteString("apply-all")
		if view.Deferred > 0 {
			fmt.Fprintf(&b, " · %d deferred", view.Deferred)
		}
		b.WriteByte('\n')
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func plural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}

// RenderAgentHandoff is deliberately hard-register text. Explain carries no
// stored explanation: it sends only the selected entry to the attending agent,
// which must read the source and answer live.
func RenderAgentHandoff(w io.Writer, action string, ids []string) error {
	if w == nil {
		return fmt.Errorf("render agent handoff: nil writer")
	}
	if len(ids) == 0 {
		return fmt.Errorf("render agent handoff: no entries")
	}
	for _, id := range ids {
		if strings.TrimSpace(id) == "" {
			return fmt.Errorf("render agent handoff: empty entry")
		}
	}
	var directive string
	switch action {
	case "apply":
		directive = "RUN `clew journal show <ENTRY>` FOR EACH ENTRY. APPLY ONLY THE HUMAN-SELECTED CHANGE TO THE CURRENT WORK. DO NOT EXPAND SCOPE. REPORT EVIDENCE."
	case "build":
		directive = "RUN `clew journal show <ENTRY>`. BUILD ONLY THE HUMAN-SELECTED WORK. DO NOT EXPAND SCOPE. REPORT EVIDENCE."
	case "explain":
		directive = "RUN `clew journal show <ENTRY>`. QUOTE THE OWNER'S WORDS. EXPLAIN ITS EFFECT ON THE CURRENT WORK. DO NOT CHANGE FILES."
	default:
		return fmt.Errorf("render agent handoff: unknown action %q", action)
	}
	fmt.Fprintln(w, "CLEW_AGENT_HANDOFF_V1")
	fmt.Fprintf(w, "ACTION=%s\n", strings.ToUpper(action))
	for _, id := range ids {
		fmt.Fprintf(w, "ENTRY=%s\n", id)
	}
	fmt.Fprintln(w, directive)
	return nil
}
