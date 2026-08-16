package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"restart/internal/model"
)

const (
	DigestCap       = 4 * 1024  // §4: digest.md ≤ 4 KB
	RollupOverfire  = 32 * 1024 // §3.3: rollup > 32 KB → extractor over-firing
	rollupTimestamp = "2006-01-02 15:04 MST"
)

// sortedLive returns live entries newest-first, optionally filtered by type.
func sortedLive(j *Journal, st map[string]*Computed, t model.EntryType) []*model.Entry {
	var out []*model.Entry
	for id, e := range j.Entries {
		if t != "" && e.Type != t {
			continue
		}
		if c, ok := st[id]; ok && Live(c.Status) {
			out = append(out, e)
		}
	}
	sort.Slice(out, func(a, b int) bool { return out[a].ID > out[b].ID })
	return out
}

// Rollup renders journal.md: only non-superseded, non-expired entries (§3.3).
func Rollup(j *Journal, st map[string]*Computed, now time.Time) string {
	var b strings.Builder
	counts := map[model.EntryType]int{}
	live := 0
	for id, e := range j.Entries {
		if c, ok := st[id]; ok && Live(c.Status) {
			counts[e.Type]++
			live++
		}
	}
	fmt.Fprintf(&b, "# Journal\n\n_generated %s · %d live entries (%d decisions · %d findings · %d questions · %d intents) · %d total in history_\n",
		now.UTC().Format(rollupTimestamp), live,
		counts[model.Decision], counts[model.Finding], counts[model.Question], counts[model.Intent], len(j.Entries))

	section := func(title string, t model.EntryType) {
		entries := sortedLive(j, st, t)
		if len(entries) == 0 {
			return
		}
		fmt.Fprintf(&b, "\n## %s\n", title)
		for _, e := range entries {
			c := st[e.ID]
			fmt.Fprintf(&b, "\n### %s — %s  `%s`\n", e.ID, e.Title, c.Status)
			if e.Quote != "" {
				fmt.Fprintf(&b, "> %s\n", strings.ReplaceAll(e.Quote, "\n", " "))
			}
			if e.Body != "" {
				fmt.Fprintf(&b, "\n%s\n", strings.TrimSpace(e.Body))
			}
			meta := []string{fmt.Sprintf("source: %s %s", e.Source.Kind, e.Source.Ref)}
			meta = append(meta, fmt.Sprintf("confidence: %.2f", c.Confidence))
			if len(e.Tags) > 0 {
				meta = append(meta, "tags: "+strings.Join(e.Tags, ", "))
			}
			if c.Evidence > 0 {
				meta = append(meta, fmt.Sprintf("evidence: %d", c.Evidence))
			}
			if len(c.Contradicts) > 0 {
				meta = append(meta, "pairs-with: "+strings.Join(c.Contradicts, ", "))
			}
			if c.Tainted {
				meta = append(meta, "taint: tool_result")
			}
			fmt.Fprintf(&b, "\n_%s_\n", strings.Join(meta, " · "))
		}
	}
	section("Decisions", model.Decision)
	section("Findings", model.Finding)
	section("Open questions", model.Question)
	section("Intents", model.Intent)
	return b.String()
}

// Digest renders digest.md, the ≤4KB projection fed to extraction calls and
// agents (§4, §6.2): ids/titles/statuses of live entries.
func Digest(j *Journal, st map[string]*Computed) string {
	var b strings.Builder
	b.WriteString("# Journal digest (live entries: id · type/status · title)\n")
	var all []*model.Entry
	for id, e := range j.Entries {
		if c, ok := st[id]; ok && Live(c.Status) {
			all = append(all, e)
		}
	}
	sort.Slice(all, func(a, b int) bool { return all[a].ID > all[b].ID })
	for _, e := range all {
		line := fmt.Sprintf("- %s %s/%s %s\n", e.ID, e.Type, st[e.ID].Status, e.Title)
		if b.Len()+len(line) > DigestCap-64 {
			fmt.Fprintf(&b, "… (%d more omitted for size)\n", len(all))
			break
		}
		b.WriteString(line)
	}
	return b.String()
}

// WriteProjections regenerates journal.md and digest.md in the journal dir.
// Deterministic projections of entries+events; races are harmless (§4).
func WriteProjections(j *Journal, st map[string]*Computed, now time.Time) error {
	if err := os.WriteFile(filepath.Join(j.Dir, "journal.md"), []byte(Rollup(j, st, now)), 0o644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(j.Dir, "digest.md"), []byte(Digest(j, st)), 0o644)
}

// Overfiring reports the §3.3 hygiene flag (rollup > 32 KB).
func Overfiring(j *Journal, st map[string]*Computed, now time.Time) bool {
	return len(Rollup(j, st, now)) > RollupOverfire
}
