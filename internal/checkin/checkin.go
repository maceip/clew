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

	"clew/internal/calm"
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
	IDs      []string
	Line     string
	Entry    *model.Entry
	Priority int
}

type View struct {
	Screen    Screen
	Repo      string
	Items     []Item
	Settled   []Item
	Deferred  int
	Issues    int // source records that did not pass a clean read
	Repairs   []string
	MemoryLag time.Duration
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
		if heldForOwner(entry) {
			continue
		}
		switch handled[id] {
		case "deferred":
			view.Deferred++
			continue
		case "applied":
			continue
		case "settled":
			continue
		}
		if entry.Type == model.Question || !mergeLive(computed[id]) {
			continue
		}
		item := Item{ID: id, IDs: []string{id}, Line: CalmLine(entry), Entry: entry}
		if verifiedWork(j, entry) {
			view.Settled = appendFolded(view.Settled, item)
			continue
		}
		items = appendFolded(items, item)
	}
	sort.Slice(items, func(i, k int) bool { return items[i].ID > items[k].ID })
	sort.Slice(view.Settled, func(i, k int) bool { return view.Settled[i].ID > view.Settled[k].ID })
	if len(items) > MaxItems {
		items = items[:MaxItems]
	}
	if len(view.Settled) > MaxItems {
		view.Settled = view.Settled[:MaxItems]
	}
	view.Items = items
	return view
}

func verifiedWork(j *journal.Journal, entry *model.Entry) bool {
	if j == nil || entry == nil || (entry.Type != model.Decision && entry.Type != model.Intent) {
		return false
	}
	for _, event := range j.Events {
		if event.Entry != entry.ID || event.Kind != model.EvEvidence {
			continue
		}
		switch event.PStr("kind") {
		case "completion":
			return true
		case "commit":
			if event.PStr("ref") != "" {
				return true
			}
		}
	}
	return false
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
		if heldForOwner(entry) {
			continue
		}
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
			items = appendFolded(items, Item{ID: id, IDs: []string{id}, Line: CalmLine(entry), Entry: entry, Priority: priority})
		case model.Decision:
			if c.Status != journal.StActive && c.Status != journal.StPossibleContradiction && c.Status != journal.StContradicted {
				continue
			}
			if decisionHasLiveFailure(entry, alerts) {
				items = appendFolded(items, Item{ID: id, IDs: []string{id}, Line: CalmLine(entry), Entry: entry, Priority: 0})
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

func heldForOwner(entry *model.Entry) bool {
	if entry == nil {
		return false
	}
	text := strings.ToLower(entry.Title + " " + entry.Body)
	return strings.Contains(text, "held for") || strings.Contains(text, "not buildable spec")
}

func appendFolded(items []Item, item Item) []Item {
	key := foldKey(item.Line)
	for i := range items {
		if foldKey(items[i].Line) != key {
			continue
		}
		items[i].IDs = appendUnique(items[i].IDs, item.IDs...)
		if item.ID > items[i].ID {
			items[i].ID = item.ID
			items[i].Entry = item.Entry
		}
		if item.Priority < items[i].Priority {
			items[i].Priority = item.Priority
		}
		return items
	}
	return append(items, item)
}

func appendUnique(ids []string, more ...string) []string {
	seen := make(map[string]bool, len(ids)+len(more))
	for _, id := range ids {
		seen[id] = true
	}
	for _, id := range more {
		if id != "" && !seen[id] {
			ids = append(ids, id)
			seen[id] = true
		}
	}
	return ids
}

func foldKey(line string) string {
	return strings.Join(words(line), " ")
}

// EntryIDs returns the machine identities behind one human line. Folded lines
// deliberately keep every source so apply/build/retire records each one.
func EntryIDs(item Item) []string {
	ids := append([]string(nil), item.IDs...)
	if len(ids) == 0 && item.ID != "" {
		ids = append(ids, item.ID)
	}
	sort.Strings(ids)
	return ids
}

// Resolve maps plain spoken words back to one displayed line. Exact identity
// remains accepted for agent plumbing, but never has to pass through a human.
func Resolve(view View, spoken string) (Item, error) {
	query := words(spoken)
	if len(query) == 0 {
		return Item{}, fmt.Errorf("name the change in your own words")
	}
	best, bestScore, tied := -1, 0, false
	for i, item := range view.Items {
		for _, id := range EntryIDs(item) {
			if strings.TrimSpace(spoken) == id {
				return item, nil
			}
		}
		text := item.Line
		if item.Entry != nil {
			text += " " + item.Entry.Title + " " + item.Entry.Body
		}
		hay := make(map[string]bool)
		for _, word := range words(text) {
			hay[word] = true
		}
		score := 0
		for _, word := range query {
			if hay[word] {
				score++
			}
		}
		if score > bestScore {
			best, bestScore, tied = i, score, false
		} else if score > 0 && score == bestScore {
			tied = true
		}
	}
	if best < 0 || bestScore < len(query) || tied {
		return Item{}, fmt.Errorf("that wording does not name one change on this screen")
	}
	return view.Items[best], nil
}

func words(s string) []string {
	var b strings.Builder
	for _, r := range strings.ToLower(s) {
		switch {
		case r >= 'a' && r <= 'z', r >= '0' && r <= '9':
			b.WriteRune(r)
		default:
			b.WriteByte(' ')
		}
	}
	stop := map[string]bool{"a": true, "an": true, "and": true, "apply": true, "build": true, "defer": true, "explain": true, "retire": true, "the": true, "to": true}
	var out []string
	for _, word := range strings.Fields(b.String()) {
		if !stop[word] {
			out = append(out, word)
		}
	}
	return out
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
	return len(j.LoadErrors)
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
		return "Human-facing words stay calm everywhere, and docket keeps its name"
	case strings.Contains(title, "wording sweep covers every fear-attached word"):
		return "Human-facing words stay calm everywhere, and docket keeps its name"
	case strings.Contains(title, "limiter gates distillation timing"):
		return "Recording never stops, and memory catches up under your ceiling"
	case strings.Contains(title, "evidence settles merge lines"):
		return "Finished and verified work settles itself and tells you once"
	case strings.Contains(title, "name the system restart"):
		return "The system is named restart"
	case strings.Contains(title, "entry ids are machine plumbing"):
		return "Screen lines use plain speech, fold repeats, and leave held work alone"
	case strings.Contains(title, "finished means shared"), strings.Contains(title, "finish message is a surface"):
		return "Finished work is shared, and its message says what exists, where it lives, and what you can say next"
	case strings.Contains(title, "lines are plain speech"):
		return "Screen lines use plain speech, fold repeats, and leave held work alone"
	case strings.Contains(title, "broken states carry their verb"):
		return "Broken lines name the fix or the agent already fixing it"
	case strings.Contains(title, "knowledge merge at finish"), strings.Contains(title, "merge lines must pass"), strings.Contains(title, "explain is live"):
		return "Each remembered change stays clear enough to apply, explain, or defer tomorrow"
	case strings.Contains(title, "silence is the signal"):
		return "Quiet appears only after every place has been checked for something new"
	case strings.Contains(title, "second tab: the intent gap"):
		return "The intent gap shows everything we meant to build but have not made real"
	case strings.Contains(title, "freshness ladder"):
		return "Every agent needs new decisions before the next human message"
	case strings.Contains(title, "i9 frugality replaced") || strings.Contains(title, "budget-deadlock"):
		return "Listening needs a spend floor above one full request and under the owner’s ceiling"
	case strings.Contains(title, "birth detection"):
		return "A new project starts with the owner’s standing decisions and nothing from an old project"
	case strings.Contains(title, "phone reads the glance"):
		return "The phone shows what changed and asks only when a decision needs you"
	case strings.Contains(title, "laptop agents fully sensed"):
		return "Laptop agents notice all work without asking you to record it"
	case strings.Contains(title, "pr-only cloud agents"), strings.Contains(title, "repo-write cloud agents"):
		return "Cloud agents return what they learned whether they write or open a pull request"
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
		fmt.Fprintf(&b, "KNOWLEDGE MERGE — %s remembers what we decide\n", repo)
	case IntentGap:
		fmt.Fprintf(&b, "INTENT GAP — %s shows what is not real yet\n", repo)
	default:
		return fmt.Errorf("render check-in: unknown screen %q", view.Screen)
	}
	if view.MemoryLag > 0 {
		fmt.Fprintf(&b, "memory is %d minutes behind\n", lagMinutes(view.MemoryLag))
	}
	for _, repair := range view.Repairs {
		if !lineAlreadyShown(repair, view.Items) {
			b.WriteString(calm.Text(repair))
			b.WriteByte('\n')
		}
	}
	if view.Screen == KnowledgeMerge && len(view.Settled) > 0 {
		b.WriteString("Settled while you were away\n")
		for _, item := range view.Settled {
			fmt.Fprintf(&b, "✓ %s\n", calm.Text(item.Line))
		}
	}
	if len(view.Items) == 0 {
		if view.Issues > 0 {
			b.WriteString("The attending agent must repair saved knowledge before this can be called quiet.\n")
		} else if len(view.Repairs) > 0 {
			b.WriteString("Quiet waits until that work is fixed.\n")
		} else if view.Screen == KnowledgeMerge && len(view.Settled) > 0 {
			b.WriteString("Nothing needs applying.\n")
		} else if view.Screen == KnowledgeMerge {
			b.WriteString("Nothing new.\n")
		} else {
			b.WriteString("Everything here is real.\n")
		}
		_, err := io.WriteString(w, b.String())
		return err
	}
	for i, item := range view.Items {
		fmt.Fprintf(&b, "%d. %s", i+1, calm.Text(item.Line))
		if view.Screen == KnowledgeMerge {
			b.WriteString(" — apply · explain · defer\n")
		} else {
			b.WriteString(" — build · explain · retire\n")
		}
	}
	if view.Screen == KnowledgeMerge {
		b.WriteString("apply all")
		if view.Deferred > 0 {
			fmt.Fprintf(&b, " · %d deferred", view.Deferred)
		}
		b.WriteByte('\n')
	}
	_, err := io.WriteString(w, b.String())
	return err
}

func lagMinutes(lag time.Duration) int {
	minutes := int((lag + time.Minute - 1) / time.Minute)
	if minutes < 1 {
		return 1
	}
	return minutes
}

func lineAlreadyShown(repair string, items []Item) bool {
	repairWords := foldKey(repair)
	for _, item := range items {
		lineWords := foldKey(item.Line)
		if strings.Contains(repairWords, "spend floor") && strings.Contains(lineWords, "spend floor") {
			return true
		}
	}
	return false
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
