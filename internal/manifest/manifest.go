// Package manifest implements the clew manifest (JOURNAL_SPEC §9): the
// tool for the reset moment. Disposition pass → MANIFEST.md, human pass
// (edit carry/drop marks), then SEED.md + genesis/ outputs. The choice to
// lose knowledge is itself journaled as disposition events.
package manifest

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/model"
)

const SeedCap = 4 * 1024 // §9.3: SEED.md ≤ 4 KB

type Disposition struct {
	Entry    string
	Verdict  string // covered | missing | contradicts | "" (no spec given)
	SpecLine string
}

// Generate writes MANIFEST.md into the journal dir. With a spec, every live
// entry is dispositioned (LLM when available, else a mechanical keyword
// pass, labeled as such). Without a spec, all entries render as the carry
// candidate list.
func Generate(j *journal.Journal, st map[string]*journal.Computed, specPath string, provider llm.Provider, now time.Time) (string, error) {
	var live []*model.Entry
	for id, e := range j.Entries {
		if c := st[id]; c != nil && journal.Live(c.Status) {
			live = append(live, e)
		}
	}
	sort.Slice(live, func(a, b int) bool { return live[a].ID < live[b].ID })

	disp := map[string]Disposition{}
	method := ""
	if specPath != "" {
		spec, err := os.ReadFile(specPath)
		if err != nil {
			return "", fmt.Errorf("read spec: %w", err)
		}
		if provider != nil {
			disp, method = llmDispositions(live, string(spec), provider)
		}
		if method == "" {
			disp = keywordDispositions(live, string(spec))
			method = "mechanical keyword pass (no LLM provider available)"
		}
	}

	var b strings.Builder
	fmt.Fprintf(&b, "# Restart manifest — %s\n\n", now.UTC().Format("2006-01-02"))
	b.WriteString("Mark each entry `[carry]` or `[drop]`, then run `clew manifest --out <dir>`.\n")
	if specPath != "" {
		fmt.Fprintf(&b, "Disposition vs spec %s (%s): covered / missing / contradicts.\n", filepath.Base(specPath), method)
	}
	b.WriteString("\n")
	for _, e := range live {
		c := st[e.ID]
		fmt.Fprintf(&b, "## %s — %s (%s, %s)\n", e.ID, e.Title, e.Type, c.Status)
		d := disp[e.ID]
		mark := "[carry]"
		if d.Verdict == "covered" {
			mark = "[drop]" // covered by the new spec → default drop (already there)
		}
		fmt.Fprintf(&b, "disposition: %s", mark)
		if d.Verdict != "" {
			fmt.Fprintf(&b, "   spec: %s", d.Verdict)
			if d.SpecLine != "" {
				fmt.Fprintf(&b, " → %q", clip(d.SpecLine, 100))
			}
		}
		b.WriteString("\n")
		fmt.Fprintf(&b, "> %s\n", clip(e.Quote, 200))
		if e.Body != "" {
			fmt.Fprintf(&b, "%s\n", clip(e.Body, 400))
		}
		b.WriteString("\n")
	}
	p := filepath.Join(j.Dir, "MANIFEST.md")
	return p, os.WriteFile(p, []byte(b.String()), 0o644)
}

var manifestEntryRe = regexp.MustCompile(`(?m)^## (e[0-9A-HJKMNP-TV-Z]{26}) `)
var manifestMarkRe = regexp.MustCompile(`disposition:\s*\[(carry|drop)\]`)

type ApplyResult struct {
	Carried, Dropped []string
	SeedPath         string
	GenesisDir       string
}

// Apply reads the marked MANIFEST.md, records disposition events, and emits
// SEED.md + genesis/ into outDir.
func Apply(j *journal.Journal, st map[string]*journal.Computed, manifestPath, outDir, surface string, now time.Time) (*ApplyResult, error) {
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		return nil, err
	}
	text := string(b)
	locs := manifestEntryRe.FindAllStringSubmatchIndex(text, -1)
	marks := map[string]string{}
	for i, loc := range locs {
		id := text[loc[2]:loc[3]]
		end := len(text)
		if i+1 < len(locs) {
			end = locs[i+1][0]
		}
		if m := manifestMarkRe.FindStringSubmatch(text[loc[0]:end]); m != nil {
			marks[id] = m[1]
		}
	}
	if len(marks) == 0 {
		return nil, fmt.Errorf("no [carry]/[drop] marks found in %s", manifestPath)
	}
	res := &ApplyResult{GenesisDir: filepath.Join(outDir, "genesis")}
	if err := os.MkdirAll(res.GenesisDir, 0o755); err != nil {
		return nil, err
	}

	var carried []*model.Entry
	for id, mark := range marks {
		e, ok := j.Entries[id]
		if !ok {
			continue
		}
		// The choice to lose (or keep) knowledge is itself journaled (§9.2).
		if !j.HasEvent(model.EvDisposition, id, "disposition", mark) {
			j.AddEvent(&model.Event{
				ID: ids.NewEvent(now), Kind: model.EvDisposition, Entry: id,
				Payload: map[string]any{"disposition": mapMark(mark)},
				By:      model.By{Who: "human", Surface: surface}, At: now.UTC(),
			})
		}
		if mark == "carry" {
			carried = append(carried, e)
			res.Carried = append(res.Carried, id)
		} else {
			res.Dropped = append(res.Dropped, id)
		}
	}
	sort.Slice(carried, func(a, b int) bool { return carried[a].ID < carried[b].ID })
	sort.Strings(res.Carried)
	sort.Strings(res.Dropped)

	// genesis/: carried entries as files, provenance intact (§9.3).
	for _, e := range carried {
		ge := *e
		ge.Source.Kind = model.SrcCarried // original ref/at/agent preserved
		gb, err := yaml.Marshal(&ge)
		if err != nil {
			return nil, err
		}
		if err := os.WriteFile(filepath.Join(res.GenesisDir, ge.ID+".yaml"), gb, 0o644); err != nil {
			return nil, err
		}
	}
	res.SeedPath = filepath.Join(outDir, "SEED.md")
	return res, os.WriteFile(res.SeedPath, []byte(seed(carried, st, now)), 0o644)
}

// seed renders SEED.md (≤4KB): carried decisions WITH reasons, findings as
// constraints/warnings, open questions. Paste-ready first prompt (§9.3).
func seed(carried []*model.Entry, st map[string]*journal.Computed, now time.Time) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Project seed — carried knowledge (%s)\n\nThis is the distilled, human-ratified knowledge of the predecessor project.\nTreat decisions as standing constraints; treat findings as measured reality.\n", now.UTC().Format("2006-01-02"))
	section := func(title string, t model.EntryType, render func(*model.Entry) string) {
		var lines []string
		for _, e := range carried {
			if e.Type == t {
				lines = append(lines, render(e))
			}
		}
		if len(lines) == 0 {
			return
		}
		fmt.Fprintf(&b, "\n## %s\n", title)
		for _, l := range lines {
			if b.Len()+len(l) > SeedCap-100 {
				b.WriteString("- … (truncated at 4KB — see genesis/ for the rest)\n")
				return
			}
			b.WriteString(l)
		}
	}
	section("Decisions (with reasons)", model.Decision, func(e *model.Entry) string {
		return fmt.Sprintf("- %s — %s\n", e.Title, oneLine(e.Body))
	})
	section("Findings (constraints & warnings)", model.Finding, func(e *model.Entry) string {
		env := ""
		if e.Env != nil && e.Env.Host != "" {
			env = fmt.Sprintf(" [measured on %s — budget accordingly]", e.Env.Host)
		}
		return fmt.Sprintf("- %s — %s%s\n", e.Title, oneLine(e.Body), env)
	})
	section("Open questions (still unanswered)", model.Question, func(e *model.Entry) string {
		return fmt.Sprintf("- %s (asks: %s)\n", e.Title, e.Asks)
	})
	section("Unfinished intents", model.Intent, func(e *model.Entry) string {
		return fmt.Sprintf("- %s\n", e.Title)
	})
	return b.String()
}

// Import seeds a journal from a genesis directory (init --carry, §9.4).
func Import(j *journal.Journal, genesisDir string) (int, error) {
	files, err := os.ReadDir(genesisDir)
	if err != nil {
		return 0, err
	}
	n := 0
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		b, err := os.ReadFile(filepath.Join(genesisDir, f.Name()))
		if err != nil {
			return n, err
		}
		e := &model.Entry{}
		if err := yaml.Unmarshal(b, e); err != nil {
			return n, fmt.Errorf("%s: %w", f.Name(), err)
		}
		if e.Source.Kind != model.SrcCarried {
			e.Source.Kind = model.SrcCarried
		}
		if _, dup := j.Entries[e.ID]; dup {
			continue // idempotent re-import
		}
		if err := j.AddEntry(e); err != nil {
			return n, fmt.Errorf("%s: %w", f.Name(), err)
		}
		n++
	}
	return n, nil
}

// llmDispositions maps entries to the new spec via the provider (§7.1
// mapping machinery re-used).
func llmDispositions(live []*model.Entry, spec string, p llm.Provider) (map[string]Disposition, string) {
	var b strings.Builder
	b.WriteString("For each journal entry decide against the NEW SPEC below: covered (the spec already includes it — cite the spec line), missing (the spec lost it), or contradicts (the spec says otherwise — cite the line). STRICT JSON only: {\"dispositions\":[{\"entry\":\"id\",\"verdict\":\"covered|missing|contradicts\",\"spec_line\":\"…\"}]}\n\nENTRIES:\n")
	for _, e := range live {
		fmt.Fprintf(&b, "- %s [%s] %s — %s\n", e.ID, e.Type, e.Title, clip(e.Body, 150))
	}
	b.WriteString("\nNEW SPEC:\n" + clip(spec, 24*1024))
	res, err := p.Call(b.String())
	if err != nil {
		return nil, ""
	}
	raw, ok := llm.ExtractJSON(res.Text)
	if !ok {
		return nil, ""
	}
	var loose struct {
		Dispositions []map[string]string `json:"dispositions"`
	}
	if err := yamlUnmarshalJSON(raw, &loose); err != nil {
		return nil, ""
	}
	disp := map[string]Disposition{}
	for _, d := range loose.Dispositions {
		v := d["verdict"]
		if v != "covered" && v != "missing" && v != "contradicts" {
			continue
		}
		disp[d["entry"]] = Disposition{Entry: d["entry"], Verdict: v, SpecLine: d["spec_line"]}
	}
	return disp, "LLM disposition pass (" + p.Name() + ")"
}

func yamlUnmarshalJSON(raw string, out any) error {
	return yaml.Unmarshal([]byte(raw), out) // YAML is a JSON superset
}

// keywordDispositions is the provider-less fallback: title keyword overlap.
func keywordDispositions(live []*model.Entry, spec string) map[string]Disposition {
	specLower := strings.ToLower(spec)
	specLines := strings.Split(spec, "\n")
	disp := map[string]Disposition{}
	for _, e := range live {
		words := significantWords(e.Title)
		hits := 0
		for _, w := range words {
			if strings.Contains(specLower, w) {
				hits++
			}
		}
		if len(words) > 0 && hits*2 >= len(words) { // half the title words appear
			line := ""
			for _, l := range specLines {
				ll := strings.ToLower(l)
				for _, w := range words {
					if strings.Contains(ll, w) {
						line = strings.TrimSpace(l)
						break
					}
				}
				if line != "" {
					break
				}
			}
			disp[e.ID] = Disposition{Entry: e.ID, Verdict: "covered", SpecLine: line}
		} else {
			disp[e.ID] = Disposition{Entry: e.ID, Verdict: "missing"}
		}
	}
	return disp
}

func significantWords(s string) []string {
	stop := map[string]bool{"the": true, "a": true, "an": true, "of": true, "to": true,
		"and": true, "or": true, "in": true, "on": true, "for": true, "with": true,
		"is": true, "are": true, "be": true, "it": true, "this": true, "that": true}
	var out []string
	for _, w := range strings.Fields(strings.ToLower(s)) {
		w = strings.Trim(w, ".,;:!?\"'()[]{}")
		if len(w) >= 3 && !stop[w] {
			out = append(out, w)
		}
	}
	return out
}

func mapMark(m string) string {
	if m == "carry" {
		return "carried"
	}
	return "dropped"
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func clip(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}
