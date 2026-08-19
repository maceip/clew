// Package archaeology implements the cold-start pass (JOURNAL_SPEC §5.3):
// one-time distillation of README/docs/ADRs, TODO/FIXME scan, recent commit
// messages, and open PR titles into seeded entries with source.kind:
// archaeology and confidence ≤ 0.6 — the differ's day-one intent baseline.
package archaeology

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/llm"
	"github.com/maceip/clew/internal/model"
	"github.com/maceip/clew/internal/scrub"
)

const (
	maxTODOs    = 50
	maxConf     = 0.6 // §5.3 hard cap
	maxDocBytes = 16 * 1024
	maxCommits  = 200
)

type Result struct {
	Added   int
	Skipped []string // sources skipped, with reasons (I2: displayed)
	Tokens  int
}

// Run performs archaeology over an existing repo. Mechanical passes always
// run; README/docs/commit-message distillation runs only when a provider is
// configured (skipped loudly otherwise).
func Run(j *journal.Journal, repoPath, surface string, provider llm.Provider, now time.Time) (*Result, error) {
	res := &Result{}
	add := func(e *model.Entry) {
		if e.Confidence > maxConf {
			e.Confidence = maxConf
		}
		var n int
		e.Quote, n = scrub.Scrub(e.Quote)
		_ = n
		e.Body, _ = scrub.Scrub(e.Body)
		if dupTitle(j, e.Title) {
			return
		}
		if err := j.AddEntry(e); err == nil {
			res.Added++
		}
	}

	todoScan(j, repoPath, surface, now, add)
	adrScan(j, repoPath, surface, now, add)
	prScan(repoPath, surface, now, add, res)

	if provider == nil {
		res.Skipped = append(res.Skipped, "README/docs/commit-message distillation (no LLM provider; mechanical passes only)")
		return res, nil
	}
	tok, err := distill(j, repoPath, surface, provider, now, add)
	res.Tokens = tok
	if err != nil {
		res.Skipped = append(res.Skipped, "LLM distillation failed: "+err.Error())
	}
	return res, nil
}

var todoRe = regexp.MustCompile(`(?i)\b(TODO|FIXME)\b[:( ]\s*(.{4,160})`)

var skipDirs = map[string]bool{
	".git": true, "node_modules": true, "vendor": true, "dist": true,
	"build": true, ".clew": true, "target": true, ".next": true,
	"__pycache__": true, ".venv": true,
}

func todoScan(j *journal.Journal, repo, surface string, now time.Time, add func(*model.Entry)) {
	count := 0
	filepath.WalkDir(repo, func(path string, d os.DirEntry, err error) error {
		if err != nil || count >= maxTODOs {
			return filepath.SkipAll
		}
		if d.IsDir() {
			if skipDirs[d.Name()] {
				return filepath.SkipDir
			}
			return nil
		}
		if fi, err := d.Info(); err != nil || fi.Size() > 1<<20 {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()
		head := make([]byte, 512)
		n, _ := f.Read(head)
		if bytes.IndexByte(head[:n], 0) >= 0 {
			return nil // binary
		}
		f.Seek(0, 0)
		rel, _ := filepath.Rel(repo, path)
		sc := bufio.NewScanner(f)
		sc.Buffer(make([]byte, 1024*1024), 1024*1024)
		lineNo := 0
		for sc.Scan() && count < maxTODOs {
			lineNo++
			line := sc.Text()
			m := todoRe.FindStringSubmatch(line)
			if m == nil {
				continue
			}
			count++
			title := strings.TrimSpace(m[2])
			if len(title) > model.MaxTitle {
				title = title[:model.MaxTitle-1] + "…"
			}
			add(&model.Entry{
				ID: ids.NewEntry(now), Type: model.Intent,
				Title: title, Body: fmt.Sprintf("%s comment in %s.", strings.ToUpper(m[1]), rel),
				Quote: strings.TrimSpace(line), UtteranceBy: model.ByUser,
				Source: model.Source{Kind: model.SrcArchaeology,
					Ref: fmt.Sprintf("%s#L%d", rel, lineNo), Surface: surface, At: now.UTC()},
				Confidence: 0.5,
				Tags:       []string{dirGlob(rel)},
			})
		}
		return nil
	})
}

var adrTitleRe = regexp.MustCompile(`(?m)^#\s+(.+)$`)

func adrScan(j *journal.Journal, repo, surface string, now time.Time, add func(*model.Entry)) {
	var candidates []string
	for _, pat := range []string{"docs/adr/*.md", "docs/adrs/*.md", "docs/decisions/*.md", "adr/*.md", "ADR-*.md", "docs/ADR-*.md"} {
		m, _ := filepath.Glob(filepath.Join(repo, pat))
		candidates = append(candidates, m...)
	}
	for _, p := range candidates {
		b, err := os.ReadFile(p)
		if err != nil || len(b) == 0 {
			continue
		}
		rel, _ := filepath.Rel(repo, p)
		title := filepath.Base(p)
		var quote string
		if m := adrTitleRe.FindSubmatch(b); m != nil {
			title = string(m[1])
			quote = string(m[1])
		}
		// Prefer the first line of a "Decision" section as the quote.
		if i := strings.Index(string(b), "## Decision"); i >= 0 {
			rest := string(b)[i:]
			for _, l := range strings.Split(rest, "\n")[1:] {
				if t := strings.TrimSpace(l); t != "" && !strings.HasPrefix(t, "#") {
					quote = t
					break
				}
			}
		}
		if quote == "" {
			quote = strings.TrimSpace(strings.Split(string(b), "\n")[0])
		}
		if len(title) > model.MaxTitle {
			title = title[:model.MaxTitle-1] + "…"
		}
		add(&model.Entry{
			ID: ids.NewEntry(now), Type: model.Decision,
			Title: title, Body: fmt.Sprintf("Recorded ADR at %s.", rel),
			Quote: quote, UtteranceBy: model.ByUser,
			Source: model.Source{Kind: model.SrcArchaeology, Ref: rel + "#L1",
				Surface: surface, At: now.UTC()},
			Confidence: 0.6,
			Tags:       []string{dirGlob(rel)},
		})
	}
}

func prScan(repo, surface string, now time.Time, add func(*model.Entry), res *Result) {
	if _, err := exec.LookPath("gh"); err != nil {
		res.Skipped = append(res.Skipped, "open PR titles (gh not installed)")
		return
	}
	cmd := exec.Command("gh", "pr", "list", "--state", "open", "--json", "number,title", "--limit", "50")
	cmd.Dir = repo
	out, err := cmd.Output()
	if err != nil {
		res.Skipped = append(res.Skipped, "open PR titles (gh failed: not a github repo or not authed)")
		return
	}
	var prs []struct {
		Number int    `json:"number"`
		Title  string `json:"title"`
	}
	if json.Unmarshal(out, &prs) != nil {
		return
	}
	for _, pr := range prs {
		title := pr.Title
		if len(title) > model.MaxTitle {
			title = title[:model.MaxTitle-1] + "…"
		}
		add(&model.Entry{
			ID: ids.NewEntry(now), Type: model.Intent,
			Title: title, Body: fmt.Sprintf("Open PR #%d.", pr.Number),
			Quote: pr.Title, UtteranceBy: model.ByUser,
			Source: model.Source{Kind: model.SrcArchaeology,
				Ref: fmt.Sprintf("pr#%d", pr.Number), Surface: surface, At: now.UTC()},
			Confidence: 0.5,
		})
	}
}

// distill runs the LLM pass over README + key docs + 90d of commit subjects.
// Quotes are verified verbatim against the numbered material (I7 applies to
// archaeology too).
func distill(j *journal.Journal, repo, surface string, p llm.Provider, now time.Time, add func(*model.Entry)) (int, error) {
	type srcLine struct {
		ref  string
		text string
	}
	var lines []srcLine
	appendDoc := func(rel string) {
		b, err := os.ReadFile(filepath.Join(repo, rel))
		if err != nil {
			return
		}
		if len(b) > maxDocBytes {
			b = b[:maxDocBytes]
		}
		for i, l := range strings.Split(string(b), "\n") {
			if t := strings.TrimSpace(l); t != "" {
				lines = append(lines, srcLine{fmt.Sprintf("%s#L%d", rel, i+1), t})
			}
		}
	}
	appendDoc("README.md")
	docs, _ := filepath.Glob(filepath.Join(repo, "docs", "*.md"))
	for i, d := range docs {
		if i >= 5 {
			break
		}
		rel, _ := filepath.Rel(repo, d)
		appendDoc(rel)
	}
	if out, err := gitx.Run(repo, "log", "--since=90.days", "--no-merges", "--pretty=%h %s"); err == nil && out != "" {
		for i, l := range strings.Split(out, "\n") {
			if i >= maxCommits {
				break
			}
			lines = append(lines, srcLine{"commit:" + strings.Fields(l)[0], l})
		}
	}
	if len(lines) == 0 {
		return 0, nil
	}
	var b strings.Builder
	b.WriteString(`Distill durable project knowledge from this material (README/docs/commit subjects of an existing repo). Extract decisions (choices with reasons), findings (measured/learned facts), questions (open unknowns), intents (committed future work incl. unfinished TODO themes). Every entry needs a quote copied VERBATIM from exactly one numbered line, plus that line number. STRICT JSON only:
{"entries":[{"type":"decision|finding|question|intent","title":"≤80","body":"≤400","quote":"verbatim","line":3,"confidence":0.5,"tags":["path/**"]}]}
Max 20 entries; prefer fewer, stronger ones. Material:` + "\n\n")
	for i, l := range lines {
		fmt.Fprintf(&b, "[L%d] (%s) %s\n", i+1, l.ref, l.text)
	}
	resp, err := p.Call(b.String())
	if err != nil {
		return 0, err
	}
	raw, ok := llm.ExtractJSON(resp.Text)
	if !ok {
		return resp.Tokens, fmt.Errorf("no JSON in provider output")
	}
	var out struct {
		Entries []struct {
			Type, Title, Body, Quote string
			Line                     int
			Confidence               float64
			Tags                     []string
		} `json:"entries"`
	}
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return resp.Tokens, err
	}
	for _, we := range out.Entries {
		if we.Line < 1 || we.Line > len(lines) {
			continue
		}
		src := lines[we.Line-1]
		if !strings.Contains(normWS(src.text), normWS(we.Quote)) || strings.TrimSpace(we.Quote) == "" {
			continue // fabricated → rejected (I7)
		}
		var et model.EntryType
		switch we.Type {
		case "decision":
			et = model.Decision
		case "finding":
			et = model.Finding
		case "question":
			et = model.Question
		case "intent":
			et = model.Intent
		default:
			continue
		}
		e := &model.Entry{
			ID: ids.NewEntry(now), Type: et,
			Title: we.Title, Body: we.Body, Quote: we.Quote, UtteranceBy: model.ByUser,
			Source: model.Source{Kind: model.SrcArchaeology, Ref: src.ref,
				Surface: surface, At: now.UTC()},
			Confidence: we.Confidence, Tags: we.Tags,
		}
		if et == model.Question {
			e.Asks = "any"
		}
		if len(e.Title) > model.MaxTitle {
			e.Title = e.Title[:model.MaxTitle-1] + "…"
		}
		add(e)
	}
	return resp.Tokens, nil
}

func dupTitle(j *journal.Journal, title string) bool {
	for _, e := range j.Entries {
		if strings.EqualFold(e.Title, title) {
			return true
		}
	}
	return false
}

func dirGlob(rel string) string {
	d := filepath.Dir(rel)
	if d == "." {
		return filepath.Base(rel)
	}
	return strings.Split(d, string(os.PathSeparator))[0] + "/**"
}

func normWS(s string) string { return strings.Join(strings.Fields(s), " ") }
