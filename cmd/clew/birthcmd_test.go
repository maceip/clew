package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
	"time"

	"github.com/maceip/clew/internal/adapters"
	"github.com/maceip/clew/internal/config"
	"github.com/maceip/clew/internal/docket"
	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/lineage"
	"github.com/maceip/clew/internal/model"
	"github.com/maceip/clew/internal/owner"
	"github.com/maceip/clew/internal/seed"
	"github.com/maceip/clew/internal/state"
)

func testProcessEnv(overrides ...string) []string {
	keys := make(map[string]bool, len(overrides))
	for _, override := range overrides {
		if key, _, ok := strings.Cut(override, "="); ok {
			keys[key] = true
		}
	}
	env := make([]string, 0, len(os.Environ())+len(overrides))
	for _, item := range os.Environ() {
		key, _, ok := strings.Cut(item, "=")
		if ok && keys[key] {
			continue
		}
		env = append(env, item)
	}
	return append(env, overrides...)
}

func buildClewTestBinary(t *testing.T) string {
	t.Helper()
	bin := filepath.Join(t.TempDir(), "clew")
	_, thisFile, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("resolve birth test source directory")
	}
	build := exec.Command("go", "build", "-o", bin, ".")
	build.Dir = filepath.Dir(thisFile)
	// The package under test has already resolved its modules. Refuse network
	// access here so this process-boundary fixture stays hermetic.
	build.Env = testProcessEnv("GOPROXY=off", "GOSUMDB=off")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build compiled hook command: %v\n%s", err, output)
	}
	return bin
}

func TestSessionStartAutoBirthInjectsOnlyOwnerLaws(t *testing.T) {
	home := t.TempDir()
	clewHome := filepath.Join(home, "clew-home")
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", clewHome)

	// Certify one project-agnostic finding before the successor exists.
	now := time.Date(2026, 8, 18, 18, 0, 0, 0, time.UTC)
	source, err := journal.Load(filepath.Join(t.TempDir(), "source-journal"))
	if err != nil {
		t.Fatal(err)
	}
	law := &model.Entry{
		ID: ids.NewEntry(now), Type: model.Finding,
		Title:       "Verify affected states directly",
		Body:        "Completion claims require evidence from the affected state.",
		Quote:       "verify the affected state before declaring completion",
		UtteranceBy: model.ByUser, Confidence: 1,
		Source: model.Source{Kind: model.SrcHuman, Ref: "test:owner-law", Surface: "test", At: now},
	}
	if err := source.AddEntry(law); err != nil {
		t.Fatal(err)
	}
	if _, err := owner.Default("").Promote(source, law.ID, "source", "test", now); err != nil {
		t.Fatal(err)
	}

	// This is the literal post-`mkdir x && git init` state: no clew command,
	// no carry card, no archaeology, and no predecessor selection.
	repo := filepath.Join(t.TempDir(), "x")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, err = gitx.Root(repo)
	if err != nil {
		t.Fatal(err)
	}
	input := strings.NewReader(fmt.Sprintf(`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
	var injected bytes.Buffer
	if err := runBirthHook(input, &injected); err != nil {
		t.Fatal(err)
	}

	contextPath := filepath.Join(repo, ".clew", "context.md")
	context, err := os.ReadFile(contextPath)
	if err != nil {
		t.Fatal(err)
	}
	if injected.String() != string(context) {
		t.Fatalf("SessionStart stdout differs from context.md\nstdout: %q\nfile: %q", injected.String(), context)
	}
	for _, want := range []string{"## Owner laws", law.ID, law.Quote} {
		if !strings.Contains(string(context), want) {
			t.Errorf("birth context missing %q:\n%s", want, context)
		}
	}
	for _, forbidden := range []string{"## Active decisions", "## Current findings", "## Open questions", "promotion candidate", "clew from"} {
		if strings.Contains(string(context), forbidden) {
			t.Errorf("laws-only birth context contains %q:\n%s", forbidden, context)
		}
	}

	db, err := state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if !db.RepoRegistered(repo) {
		t.Fatal("new repository was not registered with the watcher")
	}
	projectJournal, err := journal.Load(gitx.WorktreeDir(repo))
	if err != nil {
		t.Fatal(err)
	}
	if len(projectJournal.Entries) != 0 || len(projectJournal.Events) != 0 {
		t.Fatalf("birth imported lore: %d entries, %d events", len(projectJournal.Entries), len(projectJournal.Events))
	}
	snapshot, err := seed.Read(filepath.Join(repo, ".clew", "SEED.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(snapshot.Decisions)+len(snapshot.Findings)+len(snapshot.Graveyard) != 0 {
		t.Fatalf("newborn seed contains lore: %#v", snapshot)
	}
	for _, path := range []string{
		filepath.Join(repo, "CLAUDE.md"),
		filepath.Join(repo, "AGENTS.md"),
		filepath.Join(repo, ".claude", "settings.json"),
		filepath.Join(gitx.WorktreeDir(repo), "SEED.md"),
	} {
		if _, err := os.Stat(path); err != nil {
			t.Errorf("birth did not create %s: %v", path, err)
		}
	}
}

func TestFallbackBirthKeepsTriggeringSessionFromOffsetZero(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	repo := filepath.Join(t.TempDir(), "fallback")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	now := time.Now().UTC().Truncate(time.Second)
	day := now.Format("2006/01/02")
	file := filepath.Join(home, ".codex", "sessions", day, "rollout-birth.jsonl")
	if err := os.MkdirAll(filepath.Dir(file), 0o755); err != nil {
		t.Fatal(err)
	}
	content := fmt.Sprintf(
		"{\"timestamp\":%q,\"type\":\"session_meta\",\"payload\":{\"id\":\"fallback-1\",\"cwd\":%q}}\n"+
			"{\"timestamp\":%q,\"type\":\"response_item\",\"payload\":{\"type\":\"message\",\"role\":\"user\",\"content\":[{\"type\":\"input_text\",\"text\":\"FIRST PROMPT MUST BE CAPTURED\"}]}}\n",
		now.Format(time.RFC3339), repo, now.Add(time.Second).Format(time.RFC3339))
	if err := os.WriteFile(file, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := config.WriteSkeleton(); err != nil {
		t.Fatal(err)
	}
	a, err := load()
	if err != nil {
		t.Fatal(err)
	}
	defer a.close()
	w := newWatcher(a)
	w.birthSince = now.Add(-time.Minute)
	w.discoverBirths()
	if !a.db.RepoRegistered(repo) {
		t.Fatal("fallback watcher did not register the newborn repo")
	}
	for _, key := range []string{"tail:", "extract:", "history-end:"} {
		if got := a.db.Watermark(key + file); got != 0 {
			t.Fatalf("%s triggering-session watermark = %d, want 0", key, got)
		}
	}
	w.tail(repo, &adapters.Codex{}, file)
	if got := a.db.Watermark("tail:" + file); got == 0 {
		t.Fatal("fallback triggering transcript was not tailed")
	}
	sessions := a.db.LiveSessions(repo, time.Hour)
	if len(sessions) != 1 || !strings.Contains(sessions[0].Title, "FIRST PROMPT MUST BE CAPTURED") {
		t.Fatalf("first fallback prompt was not parsed: %+v", sessions)
	}
	if got := a.db.Watermark("extract:" + file); got != 0 {
		t.Fatalf("extraction start moved past first prompt: %d", got)
	}
}

func TestBirthAtReusedPathCreatesNewJournalIncarnation(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	parent := t.TempDir()
	repo := filepath.Join(parent, "x")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	envelope := func() *strings.Reader {
		return strings.NewReader(fmt.Sprintf(`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
	}
	if err := runBirthHook(envelope(), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	oldWorktree := gitx.WorktreeDir(repo)
	oldSeed, err := os.ReadFile(filepath.Join(oldWorktree, "SEED.md"))
	if err != nil {
		t.Fatal(err)
	}
	oldSnapshot, err := seed.Parse(oldSeed)
	if err != nil {
		t.Fatal(err)
	}
	oldAlertBody := "PREDECESSOR ALERT MUST NOT CROSS THE BIRTH BOUNDARY"
	oldSessionID := "predecessor-session"
	db, err := state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := db.UpsertAlert(state.Alert{
		Key: "stomp:predecessor", RepoPath: repo, Kind: "stomp",
		Body: oldAlertBody, EntryIDs: "old/file.go", Blocking: true,
		WithdrawWhen: "stomp:predecessor-resolved",
	}); err != nil {
		db.Close()
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if err := db.UpsertSession(state.Session{
		ID: oldSessionID, Adapter: "test", Agent: "claude", File: "old-session.jsonl",
		RepoPath: repo, Surface: "test", Title: "predecessor session",
		StartedAt: now.Add(-time.Minute), LastActivity: now,
	}); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.AddFootprints(oldSessionID, []string{"old/file.go"}); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.Set("sync-error:"+repo, "predecessor sync failure"); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}
	predecessor := filepath.Join(parent, "predecessor-moved-aside")
	if err := os.Rename(repo, predecessor); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	var injected bytes.Buffer
	if err := runBirthHook(envelope(), &injected); err != nil {
		t.Fatal(err)
	}
	if injected.Len() == 0 {
		t.Fatal("reused-path newborn received no context")
	}
	if strings.Contains(injected.String(), oldAlertBody) {
		t.Fatalf("predecessor alert entered newborn context:\n%s", injected.String())
	}
	newWorktree := gitx.WorktreeDir(repo)
	if newWorktree == oldWorktree {
		t.Fatalf("reused path rebound predecessor worktree %s", oldWorktree)
	}
	if !gitx.JournalReady(repo) {
		t.Fatal("new journal worktree is not bound to current repo incarnation")
	}
	if got, err := os.ReadFile(filepath.Join(oldWorktree, "SEED.md")); err != nil || !bytes.Equal(got, oldSeed) {
		t.Fatalf("predecessor seed was mutated: err=%v", err)
	}
	newSeed, err := seed.Read(filepath.Join(newWorktree, "SEED.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(newSeed.Decisions)+len(newSeed.Findings)+len(newSeed.Graveyard) != 0 {
		t.Fatalf("reused-path newborn inherited lore: %#v", newSeed)
	}
	if newSeed.Repository.ID == oldSnapshot.Repository.ID {
		t.Fatalf("reused-path seed identity did not advance: predecessor and newborn are both %s", newSeed.Repository.ID)
	}
	db, err = state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if alerts := db.OpenAlerts(repo, false); len(alerts) != 0 {
		t.Fatalf("predecessor alerts survived path incarnation reset: %#v", alerts)
	}
	if sessions := db.LiveSessions(repo, time.Hour); len(sessions) != 0 {
		t.Fatalf("predecessor sessions survived path incarnation reset: %#v", sessions)
	}
	if footprints := db.Footprints(oldSessionID); len(footprints) != 0 {
		t.Fatalf("predecessor footprints survived path incarnation reset: %#v", footprints)
	}
	if got := db.Get("sync-error:" + repo); got != "" {
		t.Fatalf("predecessor repo status survived path incarnation reset: %q", got)
	}
	newJournal, err := journal.Load(newWorktree)
	if err != nil {
		t.Fatal(err)
	}
	if cards := docket.Build(docket.Input{
		Journal: newJournal, Alerts: db.OpenAlerts(repo, true), Now: time.Now(),
	}); len(cards) != 0 {
		t.Fatalf("predecessor state entered newborn docket: %#v", cards)
	}

	// The distinct incarnation ids make the moved predecessor an eligible
	// explicit source instead of a false self-lineage cycle.
	priorCWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(priorCWD)
	if err := cmdFrom([]string{predecessor}); err != nil {
		t.Fatalf("explicit lineage from moved predecessor: %v", err)
	}
	links, err := lineage.LoadLinks(newWorktree)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 || links[0].From.ID != oldSnapshot.Repository.ID || links[0].TargetRepository != newSeed.Repository.ID {
		t.Fatalf("reused-path lineage link = %#v; predecessor=%s newborn=%s", links, oldSnapshot.Repository.ID, newSeed.Repository.ID)
	}
}

func TestBirthRepairsKnownIncarnationWithoutResettingState(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	repo := filepath.Join(t.TempDir(), "repair")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	var err error
	repo, err = gitx.Root(repo)
	if err != nil {
		t.Fatal(err)
	}
	envelope := func() *strings.Reader {
		return strings.NewReader(fmt.Sprintf(`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
	}
	if err := runBirthHook(envelope(), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}

	db, err := state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	const alertBody = "KNOWN INCARNATION STATE MUST SURVIVE REPAIR"
	if _, err := db.UpsertAlert(state.Alert{
		Key: "stomp:repair", RepoPath: repo, Kind: "stomp", Body: alertBody,
		EntryIDs: "repair/file.go", Blocking: true, WithdrawWhen: "repair:complete",
	}); err != nil {
		db.Close()
		t.Fatal(err)
	}
	now := time.Now().UTC()
	if err := db.UpsertSession(state.Session{
		ID: "repair-session", RepoPath: repo, StartedAt: now.Add(-time.Minute), LastActivity: now,
	}); err != nil {
		db.Close()
		t.Fatal(err)
	}
	if err := db.Close(); err != nil {
		t.Fatal(err)
	}

	worktree := gitx.WorktreeDir(repo)
	if _, err := gitx.Run(repo, "worktree", "remove", "--force", worktree); err != nil {
		t.Fatal(err)
	}
	if gitx.JournalReady(repo) {
		t.Fatal("removed journal worktree still reported ready")
	}
	if !gitx.HasBirthIncarnation(repo) {
		t.Fatal("known repository lost local incarnation evidence")
	}
	var injected bytes.Buffer
	if err := runBirthHook(envelope(), &injected); err != nil {
		t.Fatalf("repair known incarnation: %v", err)
	}
	if !gitx.JournalReady(repo) {
		t.Fatal("birth did not repair the journal worktree")
	}
	if !strings.Contains(injected.String(), alertBody) {
		t.Fatalf("known-incarnation alert was reset during repair:\n%s", injected.String())
	}
	db, err = state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if alerts := db.OpenAlerts(repo, false); len(alerts) != 1 || alerts[0].Body != alertBody {
		t.Fatalf("known-incarnation alerts changed during repair: %#v", alerts)
	}
	if sessions := db.LiveSessions(repo, time.Hour); len(sessions) != 1 || sessions[0].ID != "repair-session" {
		t.Fatalf("known-incarnation sessions changed during repair: %#v", sessions)
	}
}

func TestReusedPathGetsFreshIdentityWhenPredecessorWorktreeWasRemoved(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	parent := t.TempDir()
	repo := filepath.Join(parent, "reused")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	var err error
	repo, err = gitx.Root(repo)
	if err != nil {
		t.Fatal(err)
	}
	envelope := func() *strings.Reader {
		return strings.NewReader(fmt.Sprintf(`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
	}
	if err := runBirthHook(envelope(), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	oldWorktree := gitx.WorktreeDir(repo)
	oldSeed, err := seed.Read(filepath.Join(repo, ".clew", "SEED.md"))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "worktree", "remove", "--force", oldWorktree); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(oldWorktree); !os.IsNotExist(err) {
		t.Fatalf("predecessor worktree was not removed: %v", err)
	}

	predecessor := filepath.Join(parent, "predecessor")
	if err := os.Rename(repo, predecessor); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	if err := runBirthHook(envelope(), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	incarnation := gitx.ConfiguredJournalID(repo)
	if incarnation == "" || incarnation == gitx.RepoID(repo) {
		t.Fatalf("reused-path newborn did not persist a fresh journal incarnation: %q", incarnation)
	}
	newWorktree := gitx.WorktreeDir(repo)
	if newWorktree == oldWorktree {
		t.Fatalf("newborn reused deleted predecessor namespace %s", oldWorktree)
	}
	newSeed, err := seed.Read(filepath.Join(repo, ".clew", "SEED.md"))
	if err != nil {
		t.Fatal(err)
	}
	if newSeed.Repository.ID == oldSeed.Repository.ID {
		t.Fatalf("deleted predecessor worktree caused self identity %s", newSeed.Repository.ID)
	}

	priorCWD, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(repo); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(priorCWD)
	if err := cmdFrom([]string{predecessor}); err != nil {
		t.Fatalf("explicit lineage after predecessor worktree removal: %v", err)
	}
	links, err := lineage.LoadLinks(newWorktree)
	if err != nil {
		t.Fatal(err)
	}
	if len(links) != 1 || links[0].From.ID != oldSeed.Repository.ID || links[0].TargetRepository != newSeed.Repository.ID {
		t.Fatalf("lineage link after predecessor worktree removal = %#v", links)
	}
}

func TestSteadySessionStartReadsContextWithoutOpeningMachineState(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	repo := filepath.Join(t.TempDir(), "steady")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	envelope := func() *strings.Reader {
		return strings.NewReader(fmt.Sprintf(`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
	}
	if err := runBirthHook(envelope(), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	contextPath := filepath.Join(repo, ".clew", "context.md")
	context, err := os.ReadFile(contextPath)
	if err != nil {
		t.Fatal(err)
	}
	context = append(context, []byte("\nFAST-PATH-SENTINEL\n")...)
	if err := os.WriteFile(contextPath, context, 0o644); err != nil {
		t.Fatal(err)
	}
	statePath := state.DefaultPath()
	if err := os.Rename(statePath, statePath+".saved"); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(statePath, 0o755); err != nil {
		t.Fatal(err)
	}
	var output bytes.Buffer
	if err := runBirthHook(envelope(), &output); err != nil {
		t.Fatalf("steady SessionStart touched unavailable state DB: %v", err)
	}
	if output.String() != string(context) {
		t.Fatalf("steady SessionStart did not return maintained context\ngot: %q\nwant: %q", output.String(), context)
	}
}

func TestInstalledBirthHandlerExecutesCompiledCommand(t *testing.T) {
	bin := buildClewTestBinary(t)
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	now := time.Now().UTC()
	source, err := journal.Load(filepath.Join(t.TempDir(), "source-journal"))
	if err != nil {
		t.Fatal(err)
	}
	law := &model.Entry{
		ID: ids.NewEntry(now), Type: model.Finding, Title: "Generated title is not ambient",
		Body: "Generated body is not ambient", Quote: "inspect the affected state before claiming completion",
		UtteranceBy: model.ByUser, Confidence: 1,
		Source: model.Source{Kind: model.SrcHuman, Ref: "test:compiled-hook", Surface: "test", At: now},
	}
	if err := source.AddEntry(law); err != nil {
		t.Fatal(err)
	}
	if _, err := owner.Default("").Promote(source, law.ID, "source", "test", now); err != nil {
		t.Fatal(err)
	}

	if err := installClaudeBirthHook(bin); err != nil {
		t.Fatal(err)
	}
	settingsPath, err := claudeUserSettingsPath()
	if err != nil {
		t.Fatal(err)
	}
	settings := readJSONMap(t, settingsPath)
	hooks := settings["hooks"].(map[string]any)["SessionStart"].([]any)
	var command string
	var args []string
	var matcher string
	var timeout float64
	for _, rawGroup := range hooks {
		group, _ := rawGroup.(map[string]any)
		for _, rawHandler := range group["hooks"].([]any) {
			handler, _ := rawHandler.(map[string]any)
			if !isBirthHandler(handler, bin) {
				continue
			}
			matcher, _ = group["matcher"].(string)
			command, _ = handler["command"].(string)
			timeout, _ = handler["timeout"].(float64)
			for _, rawArg := range handler["args"].([]any) {
				args = append(args, rawArg.(string))
			}
		}
	}
	if matcher != "startup" || command != bin || len(args) != 1 || args[0] != "_birth" || timeout != 30 {
		t.Fatalf("installed handler shape: matcher=%q command=%q args=%q timeout=%v", matcher, command, args, timeout)
	}

	repo := filepath.Join(t.TempDir(), "compiled-birth")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	run := exec.Command(command, args...)
	run.Stdin = strings.NewReader(fmt.Sprintf(`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
	run.Env = os.Environ()
	output, err := run.CombinedOutput()
	if err != nil {
		t.Fatalf("installed handler failed: %v\n%s", err, output)
	}
	if !strings.Contains(string(output), law.Quote) || strings.Contains(string(output), law.Title) {
		t.Fatalf("compiled handler did not inject laws-only context:\n%s", output)
	}
	db, err := state.Open(state.DefaultPath())
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if !db.RepoRegistered(repo) || !gitx.BirthReady(repo) {
		t.Fatal("compiled installed handler did not complete watcher enrollment")
	}
}

func TestConcurrentColdBirthProcessesShareMachineBootstrap(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("the watcher supports launchd and systemd-user; gate uses POSIX sh")
	}
	bin := buildClewTestBinary(t)
	home := t.TempDir()
	clewHome := filepath.Join(home, "clew-home")
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", clewHome)
	t.Setenv("CLAUDE_CONFIG_DIR", "")
	t.Setenv("GIT_CONFIG_NOSYSTEM", "1")

	repos := []string{
		filepath.Join(t.TempDir(), "cold-a"),
		filepath.Join(t.TempDir(), "cold-b"),
	}
	for i, repo := range repos {
		if err := os.MkdirAll(repo, 0o755); err != nil {
			t.Fatal(err)
		}
		if _, err := gitx.Run(repo, "init", "-q"); err != nil {
			t.Fatal(err)
		}
		root, err := gitx.Root(repo)
		if err != nil {
			t.Fatal(err)
		}
		repos[i] = root
	}

	type birthProcess struct {
		repo           string
		cmd            *exec.Cmd
		stdout, stderr bytes.Buffer
	}
	children := make([]birthProcess, len(repos))
	gate := filepath.Join(t.TempDir(), "release")
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	const gatedExec = `while [ ! -e "$1" ]; do sleep 0.01; done; exec "$2" _birth`
	for i, repo := range repos {
		children[i].repo = repo
		children[i].cmd = exec.CommandContext(ctx, "sh", "-c", gatedExec, "clew-birth-gate", gate, bin)
		children[i].cmd.Stdin = strings.NewReader(fmt.Sprintf(
			`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo))
		children[i].cmd.Stdout = &children[i].stdout
		children[i].cmd.Stderr = &children[i].stderr
		children[i].cmd.Env = testProcessEnv(
			"HOME="+home,
			"CLEW_HOME="+clewHome,
			"CLAUDE_CONFIG_DIR=",
			"GIT_CONFIG_NOSYSTEM=1",
		)
		if err := children[i].cmd.Start(); err != nil {
			cancel()
			for k := 0; k < i; k++ {
				_ = children[k].cmd.Wait()
			}
			t.Fatalf("start cold birth %d: %v", i, err)
		}
	}
	if err := os.WriteFile(gate, []byte("go\n"), 0o600); err != nil {
		cancel()
		for i := range children {
			_ = children[i].cmd.Wait()
		}
		t.Fatal(err)
	}

	var failures []string
	for i := range children {
		err := children[i].cmd.Wait()
		if err != nil {
			failures = append(failures, fmt.Sprintf(
				"%s: %v; stdout=%q stderr=%q",
				children[i].repo, err, children[i].stdout.String(), children[i].stderr.String()))
		}
	}
	if ctx.Err() != nil {
		t.Fatalf("simultaneous cold births exceeded 20s bound: %v", ctx.Err())
	}
	if len(failures) != 0 {
		t.Fatalf("simultaneous cold births failed:\n%s", strings.Join(failures, "\n"))
	}

	db, err := state.Open(filepath.Join(clewHome, "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	for i := range children {
		if !db.RepoRegistered(children[i].repo) {
			t.Errorf("cold birth did not register %s", children[i].repo)
		}
		if !gitx.BirthReady(children[i].repo) {
			t.Errorf("cold birth did not publish readiness for %s", children[i].repo)
		}
		contextPath := filepath.Join(children[i].repo, ".clew", "context.md")
		contextBytes, err := os.ReadFile(contextPath)
		if err != nil {
			t.Errorf("read cold-birth context for %s: %v", children[i].repo, err)
			continue
		}
		if children[i].stdout.String() != string(contextBytes) {
			t.Errorf("cold-birth stdout differs from %s", contextPath)
		}
		if !strings.Contains(string(contextBytes), "clew context: project memory") {
			t.Errorf("cold-birth context for %s is not a clew context", children[i].repo)
		}
	}
}
