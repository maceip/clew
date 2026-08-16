// Package acceptance implements JOURNAL_SPEC §10: the three load-bearing
// acceptance tests, no ceremony.
package acceptance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"restart/internal/adapters"
	"restart/internal/config"
	"restart/internal/differ"
	"restart/internal/extract"
	"restart/internal/gitx"
	"restart/internal/ids"
	"restart/internal/journal"
	"restart/internal/llm"
	"restart/internal/manifest"
	"restart/internal/materialize"
	"restart/internal/model"
	"restart/internal/poller"
	"restart/internal/state"
)

// ---------- shared plumbing ----------

func mkRepo(t *testing.T, base, name string) string {
	t.Helper()
	bare := filepath.Join(base, name+".git")
	if _, err := gitx.Run(base, "init", "-q", "--bare", bare); err != nil {
		t.Fatal(err)
	}
	repo := filepath.Join(base, name)
	if _, err := gitx.Run(base, "clone", "-q", bare, repo); err != nil {
		t.Fatal(err)
	}
	return repo
}

func commitFile(t *testing.T, repo, rel, content, msg string, when time.Time) {
	t.Helper()
	p := filepath.Join(repo, rel)
	os.MkdirAll(filepath.Dir(p), 0o755)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "add", "-A"); err != nil {
		t.Fatal(err)
	}
	env := when.UTC().Format(time.RFC3339)
	if _, err := gitx.Run(repo, "-c", "user.name=dev", "-c", "user.email=dev@x",
		"commit", "-q", "-m", msg, "--date", env); err != nil {
		t.Fatal(err)
	}
}

func addIntent(t *testing.T, j *journal.Journal, title string, tags []string, at time.Time) *model.Entry {
	t.Helper()
	e := &model.Entry{
		ID: ids.NewEntry(at), Type: model.Intent, Title: title,
		Quote: "we should " + title, UtteranceBy: model.ByUser,
		Source:     model.Source{Kind: model.SrcSession, Ref: "claude:s.jsonl#L1", At: at},
		Confidence: 0.9, Tags: tags,
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	return e
}

// ---------- Acceptance 1: absence detection (§10.1) ----------
// Fixture shape: the agentdesk ground truth — a workload-runner intent that
// never gains evidence while surface intents do. `map` must mark it absent
// no later than the point where 5 eligible sibling intents gained evidence.
// The sibling rule is the criterion, not commit counts.

func TestAcceptance1_AbsenceDetection(t *testing.T) {
	t.Setenv("RESTART_HOME", t.TempDir())
	base := t.TempDir()
	repo := mkRepo(t, base, "agentdesk")
	commitFile(t, repo, "README.md", "agentdesk", "init", time.Now().Add(-40*24*time.Hour))

	db, err := state.Open(filepath.Join(t.TempDir(), "s.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	db.RegisterRepo(repo, "origin")

	wt, err := gitx.EnsureJournal(repo)
	if err != nil {
		t.Fatal(err)
	}
	j, err := journal.Load(wt)
	if err != nil {
		t.Fatal(err)
	}

	start := time.Now().Add(-30 * 24 * time.Hour)
	target := addIntent(t, j, "workload runner: launch an agent process against a tip",
		[]string{"supervisor/**"}, start)
	var siblings []*model.Entry
	for i := 0; i < 11; i++ {
		siblings = append(siblings, addIntent(t, j,
			fmt.Sprintf("surface %d: panel work", i),
			[]string{fmt.Sprintf("surfaces/s%d/**", i)}, start.Add(time.Hour)))
	}

	run := func(now time.Time) journal.Status {
		snap, err := poller.Poll(db, repo)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := differ.Run(db, &differ.Input{
			Repo: repo, Journal: j, Snapshot: snap, Surface: "test",
		}, now); err != nil {
			t.Fatal(err)
		}
		return journal.Compute(j, now)[target.ID].Status
	}

	// 4 surface intents gain real commit evidence → NOT absent yet.
	// (Dates sit inside the poller's 14-day commit window.)
	for i := 0; i < 4; i++ {
		commitFile(t, repo, fmt.Sprintf("surfaces/s%d/panel.go", i), "x",
			fmt.Sprintf("surface %d work", i), time.Now().Add(-time.Duration(13-i)*24*time.Hour))
	}
	if got := run(time.Now()); got == journal.StAbsent {
		t.Fatalf("absent after only 4 evidenced siblings — K=5 rule violated")
	}

	// 5th sibling gains evidence → absent, and the alert is in the inbox.
	commitFile(t, repo, "surfaces/s4/panel.go", "x", "surface 4 work", time.Now().Add(-9*24*time.Hour))
	if got := run(time.Now()); got != journal.StAbsent {
		t.Fatalf("want absent after 5 evidenced siblings, got %s", got)
	}
	found := false
	for _, al := range db.OpenAlerts(repo, true) {
		if al.Kind == "absence" && strings.Contains(al.Body, "workload runner") {
			found = true
		}
	}
	if !found {
		t.Fatal("absence must reach the inbox as a human-blocking alert")
	}
	// Sanity: siblings are in_flight/proposed, never absent (they have evidence).
	st := journal.Compute(j, time.Now())
	for i := 0; i < 5; i++ {
		if st[siblings[i].ID].Status == journal.StAbsent {
			t.Fatalf("evidenced sibling %d wrongly absent", i)
		}
	}
}

// ---------- Acceptance 2: extraction fidelity (§10.2) ----------
// Ground truth: fixtures/strategy-session/labels.yaml (Appendix A), pending
// one-time human ratification. This hermetic run exercises the full
// pipeline + harness with a faithful stub provider and proves the
// zero-fabrication gate; the real-provider gate runs with RESTART_FIDELITY=1.

type labelFile struct {
	Ratified bool `yaml:"ratified"`
	Targets  struct {
		DecisionsMin          int     `yaml:"decisions_min"`
		FindingsMin           int     `yaml:"findings_min"`
		Precision             float64 `yaml:"precision"`
		Recall                float64 `yaml:"recall"`
		InstructionIterations int     `yaml:"instruction_iterations"`
	} `yaml:"targets"`
	Decisions []label `yaml:"decisions"`
	Findings  []label `yaml:"findings"`
}

type label struct {
	ID    string   `yaml:"id"`
	Title string   `yaml:"title"`
	Match []string `yaml:"match"`
}

func loadLabels(t *testing.T) *labelFile {
	t.Helper()
	b, err := os.ReadFile(fixture(t, "labels.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	lf := &labelFile{}
	if err := yaml.Unmarshal(b, lf); err != nil {
		t.Fatal(err)
	}
	return lf
}

func fixture(t *testing.T, name string) string {
	t.Helper()
	// Tests run from the package dir; fixtures live at the repo root.
	p := filepath.Join("..", "fixtures", "strategy-session", name)
	if _, err := os.Stat(p); err != nil {
		t.Fatalf("fixture missing: %s", p)
	}
	return p
}

// evalFidelity scores extracted entries against the ratifiable label set:
// an entry matches a label when all match-keywords appear in its text.
func evalFidelity(entries []*model.Entry, labels []label, typ model.EntryType) (matched int, precise int, extracted int) {
	matchedSet := map[string]bool{}
	for _, e := range entries {
		if e.Type != typ {
			continue
		}
		extracted++
		hay := strings.ToLower(e.Title + " " + e.Body + " " + e.Quote)
		for _, l := range labels {
			all := true
			for _, kw := range l.Match {
				if !strings.Contains(hay, strings.ToLower(kw)) {
					all = false
					break
				}
			}
			if all {
				matchedSet[l.ID] = true
				precise++
				break
			}
		}
	}
	return len(matchedSet), precise, extracted
}

type fidelityStub struct{ out string }

func (s *fidelityStub) Name() string { return "fidelity-stub" }
func (s *fidelityStub) Call(string) (*llm.Result, error) {
	return &llm.Result{Text: s.out, Tokens: 100}, nil
}

func TestAcceptance2_ExtractionFidelityPipeline(t *testing.T) {
	lf := loadLabels(t)
	if !lf.Ratified {
		t.Log("NOTE: labels not yet human-ratified (§10.2) — this hermetic run validates the pipeline+harness; ratify before treating scores as truth")
	}
	transcript := fixture(t, "transcript.jsonl")

	// The stub answers with entries whose quotes are copied verbatim from the
	// fixture (as a competent extractor would), plus one fabricated entry
	// that the I7 gate must reject.
	entries := []map[string]any{
		mk("decision", "Knowledge plane first, control plane opt-in", "the knowledge plane ships first and the control plane stays opt-in for a possible team phase", 2, "assistant"),
		mk("decision", "Wedge is observe-only, zero workflow change", "the wedge is observe-only — zero workflow change for me or the agents", 3, "user"),
		mk("decision", "Journal is typed entries with provenance, distillate never transcript", "the journal is typed entries with provenance: distillate never transcript, and the human can edit it", 4, "assistant"),
		mk("decision", "Divergence by default, radar never locks", "divergence is the default; reconcile is a run you spawn when needed; the radar never locks anyone", 5, "user"),
		mk("decision", "Git is transport for the journal branch", "git is the transport — the journal rides an orphan branch in the project's own remote; sensors are session streams and fs events", 6, "assistant"),
		mk("decision", "Glance + reflexion map; push only when human-blocking", "the human surface is a glance plus a reflexion map; we push only when the human is blocking", 7, "assistant"),
		mk("decision", "Intent–reality diff incl. absence is the core product claim", "the intent–reality diff, including absence, is the core product claim", 8, "user"),
		mk("finding", "agentdesk is a clean drift specimen", "the agentdesk repo is a clean drift specimen — eleven PRs of surfaces landed while the workload core never went in-flight", 9, "assistant"),
		mk("finding", "Console category owned; console features non-differentiating", "OpenHands' Agent Canvas already owns the parallel-agent console category, so console features are non-differentiating", 10, "assistant"),
		mk("finding", "Lock-and-notary fails solo persona; anti-forgetting resonates", "the lock-and-notary pitch failed with the solo persona; anti-forgetting and laptop-close-and-it-keeps-working resonated", 11, "user"),
		mk("finding", "Merge conflicts are not the pain; forgotten work is", "merge conflicts are not the pain — LLMs resolve them fine; forgotten work, lost findings and unanswered questions are the pain", 12, "user"),
		mk("finding", "Agent Teams autopsy: advisory locks, shared cwd, silent fallback", "Claude Agent Teams' published failure modes — advisory locks, shared cwd, silent fallback to isolated subagents — show publish-by-discipline coordination failing", 13, "assistant"),
		mk("finding", "FABRICATED: users demand blockchain sync", "we must add blockchain sync immediately", 3, "user"),
	}
	resp, _ := json.Marshal(map[string]any{"entries": entries})

	j, _ := journal.Load(t.TempDir())
	out, err := extract.Run(j, &fidelityStub{out: string(resp)}, &adapters.Claude{}, transcript, 0, "fixture", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if out.Rejected != 1 {
		t.Fatalf("the fabricated entry must be rejected by the I7 gate (rejected=%d)", out.Rejected)
	}

	decMatched, decPrecise, decN := evalFidelity(out.Entries, lf.Decisions, model.Decision)
	finMatched, finPrecise, finN := evalFidelity(out.Entries, lf.Findings, model.Finding)
	t.Logf("decisions: %d/%d recovered (extracted %d), findings: %d/%d (extracted %d)",
		decMatched, len(lf.Decisions), decN, finMatched, len(lf.Findings), finN)

	if decMatched < lf.Targets.DecisionsMin {
		t.Errorf("decisions recovered %d < %d target", decMatched, lf.Targets.DecisionsMin)
	}
	if finMatched < lf.Targets.FindingsMin {
		t.Errorf("findings recovered %d < %d target", finMatched, lf.Targets.FindingsMin)
	}
	precision := float64(decPrecise+finPrecise) / float64(decN+finN)
	recall := float64(decMatched+finMatched) / float64(len(lf.Decisions)+len(lf.Findings))
	if precision < lf.Targets.Precision {
		t.Errorf("precision %.2f < %.2f — §10.2 go/no-go gate", precision, lf.Targets.Precision)
	}
	if recall < lf.Targets.Recall {
		t.Errorf("recall %.2f < %.2f — §10.2 go/no-go gate", recall, lf.Targets.Recall)
	}
	// Correct quote + source span, zero fabricated (I7 already enforced).
	for _, e := range out.Entries {
		if !strings.Contains(e.Source.Ref, "transcript.jsonl#L") {
			t.Errorf("entry %s missing source span: %s", e.ID, e.Source.Ref)
		}
	}
}

// TestAcceptance2_RealProvider runs the same gate against the configured
// live provider. Opt-in: RESTART_FIDELITY=1 (costs tokens; §10.2 kill
// criterion applies after human ratification of the labels).
func TestAcceptance2_RealProvider(t *testing.T) {
	if os.Getenv("RESTART_FIDELITY") != "1" {
		t.Skip("set RESTART_FIDELITY=1 to run the live extraction fidelity gate")
	}
	lf := loadLabels(t)
	p, note := llm.Pick(config.Load())
	if p == nil {
		t.Skip("no provider: " + note)
	}
	j, _ := journal.Load(t.TempDir())
	out, err := extract.Run(j, p, &adapters.Claude{}, fixture(t, "transcript.jsonl"), 0, "fixture", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	decMatched, decPrecise, decN := evalFidelity(out.Entries, lf.Decisions, model.Decision)
	finMatched, finPrecise, finN := evalFidelity(out.Entries, lf.Findings, model.Finding)
	precision := 1.0
	if decN+finN > 0 {
		precision = float64(decPrecise+finPrecise) / float64(decN+finN)
	}
	recall := float64(decMatched+finMatched) / float64(len(lf.Decisions)+len(lf.Findings))
	t.Logf("LIVE %s: decisions %d/7, findings %d/5, precision %.2f, recall %.2f, rejected(fabrication gate) %d",
		p.Name(), decMatched, finMatched, precision, recall, out.Rejected)
	if decMatched < lf.Targets.DecisionsMin || finMatched < lf.Targets.FindingsMin ||
		precision < lf.Targets.Precision || recall < lf.Targets.Recall {
		t.Errorf("§10.2 gate FAILED — if this persists across %d instruction iterations, the product premise fails: stop and rethink (the kill criterion is written down on purpose)",
			lf.Targets.InstructionIterations)
	}
}

func mk(typ, title, quote string, line int, by string) map[string]any {
	return map[string]any{
		"type": typ, "title": title, "body": title + ".", "quote": quote,
		"line": line, "utterance_by": by, "confidence": 0.9, "tags": []string{},
	}
}

// ---------- Acceptance 3: restart round-trip (§10.3) ----------
// journal ≥10 live entries → manifest → init --carry into an empty repo →
// 100% of carried findings/decisions in the new context.md with provenance
// intact; dropped entries recorded with disposition: dropped.

func TestAcceptance3_RestartRoundTrip(t *testing.T) {
	t.Setenv("RESTART_HOME", t.TempDir())
	base := t.TempDir()
	oldRepo := mkRepo(t, base, "old-project")
	commitFile(t, oldRepo, "README.md", "old", "init", time.Now().Add(-time.Hour))

	wtOld, err := gitx.EnsureJournal(oldRepo)
	if err != nil {
		t.Fatal(err)
	}
	jOld, err := journal.Load(wtOld)
	if err != nil {
		t.Fatal(err)
	}
	now := time.Now().UTC()
	mkEntry := func(typ model.EntryType, title string) *model.Entry {
		e := &model.Entry{
			ID: ids.NewEntry(now), Type: typ, Title: title,
			Body: "because of reasons: " + title, Quote: "verbatim: " + title,
			UtteranceBy: model.ByUser,
			Source: model.Source{Kind: model.SrcSession,
				Ref: "claude:~/.claude/projects/-w/9f2c.jsonl#L42", Agent: "claude-code",
				Surface: "laptop", At: now},
			Confidence: 0.9,
		}
		if typ == model.Question {
			e.Asks = "human"
		}
		if err := jOld.AddEntry(e); err != nil {
			t.Fatal(err)
		}
		return e
	}
	var all []*model.Entry
	for i := 0; i < 4; i++ {
		all = append(all, mkEntry(model.Decision, fmt.Sprintf("decision %d", i)))
	}
	for i := 0; i < 4; i++ {
		all = append(all, mkEntry(model.Finding, fmt.Sprintf("finding %d", i)))
	}
	all = append(all, mkEntry(model.Question, "question 0"), mkEntry(model.Question, "question 1"))
	all = append(all, mkEntry(model.Intent, "intent 0"), mkEntry(model.Intent, "intent 1"))
	if len(all) < 10 {
		t.Fatal("need ≥10 live entries")
	}

	st := journal.Compute(jOld, now)
	if _, err := manifest.Generate(jOld, st, "", nil, now); err != nil {
		t.Fatal(err)
	}
	// Human pass: drop one intent and one question, carry the rest.
	mpath := filepath.Join(wtOld, "MANIFEST.md")
	mb, _ := os.ReadFile(mpath)
	text := string(mb)
	text = dropMark(t, text, all[10].ID) // intent 0
	text = dropMark(t, text, all[8].ID)  // question 0
	os.WriteFile(mpath, []byte(text), 0o644)

	outDir := filepath.Join(base, "kit")
	res, err := manifest.Apply(jOld, st, mpath, outDir, "laptop", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Carried) != 10 || len(res.Dropped) != 2 {
		t.Fatalf("carried %d dropped %d, want 10/2", len(res.Carried), len(res.Dropped))
	}
	// Dropped recorded as disposition events (§9.2: loss is deliberate+dated).
	for _, id := range res.Dropped {
		if !jOld.HasEvent(model.EvDisposition, id, "disposition", "dropped") {
			t.Errorf("dropped %s not journaled", id)
		}
	}
	if fi, _ := os.Stat(res.SeedPath); fi.Size() > manifest.SeedCap {
		t.Errorf("SEED.md exceeds 4KB: %d", fi.Size())
	}

	// Import into an EMPTY repo (no commits at all).
	newRepo := mkRepo(t, base, "new-project")
	wtNew, err := gitx.EnsureJournal(newRepo)
	if err != nil {
		t.Fatal(err)
	}
	jNew, err := journal.Load(wtNew)
	if err != nil {
		t.Fatal(err)
	}
	n, err := manifest.Import(jNew, res.GenesisDir)
	if err != nil {
		t.Fatal(err)
	}
	if n != 10 {
		t.Fatalf("imported %d, want 10", n)
	}

	db, _ := state.Open(filepath.Join(t.TempDir(), "s.db"))
	defer db.Close()
	stNew := journal.Compute(jNew, now)
	if err := materialize.Write(newRepo, jNew, stNew, db, now); err != nil {
		t.Fatal(err)
	}
	ctx, err := os.ReadFile(filepath.Join(newRepo, ".restart", "context.md"))
	if err != nil {
		t.Fatal(err)
	}
	// 100% of carried decisions+findings present in the new context.md.
	for _, e := range all[:8] {
		if !strings.Contains(string(ctx), e.Title) {
			t.Errorf("carried %s %q missing from new context.md", e.Type, e.Title)
		}
	}
	// Provenance intact on the imported entries.
	for _, id := range res.Carried {
		ne, ok := jNew.Entries[id]
		if !ok {
			t.Errorf("carried %s missing in new journal", id)
			continue
		}
		if ne.Source.Kind != model.SrcCarried {
			t.Errorf("%s source.kind = %s, want carried", id, ne.Source.Kind)
		}
		if !strings.Contains(ne.Source.Ref, "9f2c.jsonl#L42") {
			t.Errorf("%s original ref lost: %s", id, ne.Source.Ref)
		}
		if !ne.Source.At.Equal(jOld.Entries[id].Source.At) {
			t.Errorf("%s original timestamp lost", id)
		}
	}
	// Dropped entries must NOT be in the new journal.
	for _, id := range res.Dropped {
		if _, ok := jNew.Entries[id]; ok {
			t.Errorf("dropped %s leaked into the new journal", id)
		}
	}
}

func dropMark(t *testing.T, text, id string) string {
	t.Helper()
	i := strings.Index(text, "## "+id)
	if i < 0 {
		t.Fatalf("entry %s not in MANIFEST.md", id)
	}
	j := strings.Index(text[i:], "[carry]")
	if j < 0 {
		t.Fatalf("no [carry] mark for %s", id)
	}
	return text[:i+j] + "[drop]" + text[i+j+len("[carry]"):]
}
