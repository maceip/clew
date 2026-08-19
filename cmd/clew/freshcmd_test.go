package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
)

func TestFreshHookInjectsNewDecisionsOnce(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	repo := filepath.Join(t.TempDir(), "fresh")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	if err := runBirthHook(strings.NewReader(fmt.Sprintf(
		`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo)), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	j, err := journal.Load(gitx.WorktreeDir(repo))
	if err != nil {
		t.Fatal(err)
	}
	base := time.Date(2026, 8, 19, 4, 0, 0, 0, time.UTC)
	old := testFreshDecision(base, "old decision")
	if err := j.AddEntry(old); err != nil {
		t.Fatal(err)
	}
	start := fmt.Sprintf(`{"hook_event_name":"SessionStart","session_id":"s1","cwd":%q}`, repo)
	if err := runFreshStartHook(strings.NewReader(start)); err != nil {
		t.Fatal(err)
	}
	newDecision := testFreshDecision(base.Add(time.Minute), "new decision")
	if err := j.AddEntry(newDecision); err != nil {
		t.Fatal(err)
	}
	if err := j.AddEntry(&model.Entry{
		ID: ids.NewEntry(base.Add(2 * time.Minute)), Type: model.Finding, Title: "not injected",
		Body: "finding body", Quote: "finding quote", UtteranceBy: model.ByUser, Confidence: 1,
		Source: model.Source{Kind: model.SrcHuman, Ref: "test:finding", At: base.Add(2 * time.Minute)},
	}); err != nil {
		t.Fatal(err)
	}
	prompt := fmt.Sprintf(`{"hook_event_name":"UserPromptSubmit","session_id":"s1","cwd":%q}`, repo)
	var out bytes.Buffer
	if err := runFreshHook(strings.NewReader(prompt), &out); err != nil {
		t.Fatal(err)
	}
	var envelope freshHookOutput
	if err := json.Unmarshal(out.Bytes(), &envelope); err != nil {
		t.Fatalf("decode hook output: %v\n%s", err, out.String())
	}
	context := envelope.HookSpecificOutput.AdditionalContext
	for _, want := range []string{"CLEW_DECISION_DELTA_V1", "REGISTER=HARD", newDecision.ID, "new decision"} {
		if !strings.Contains(context, want) {
			t.Errorf("fresh context missing %q:\n%s", want, context)
		}
	}
	for _, forbidden := range []string{"TITLE=\"old decision\"", "not injected"} {
		if strings.Contains(context, forbidden) {
			t.Errorf("fresh context contains %q:\n%s", forbidden, context)
		}
	}
	out.Reset()
	if err := runFreshHook(strings.NewReader(prompt), &out); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("decision delivered twice: %s", out.String())
	}
}

func TestFreshHookBaselinesAnAlreadyRunningSession(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	t.Setenv("CLEW_HOME", filepath.Join(home, "clew-home"))
	repo := filepath.Join(t.TempDir(), "fresh-existing")
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	repo, _ = gitx.Root(repo)
	if err := runBirthHook(strings.NewReader(fmt.Sprintf(
		`{"hook_event_name":"SessionStart","source":"startup","cwd":%q}`, repo)), &bytes.Buffer{}); err != nil {
		t.Fatal(err)
	}
	j, _ := journal.Load(gitx.WorktreeDir(repo))
	if err := j.AddEntry(testFreshDecision(time.Now().UTC(), "existing")); err != nil {
		t.Fatal(err)
	}
	prompt := fmt.Sprintf(`{"hook_event_name":"UserPromptSubmit","session_id":"late","cwd":%q}`, repo)
	var out bytes.Buffer
	if err := runFreshHook(strings.NewReader(prompt), &out); err != nil {
		t.Fatal(err)
	}
	if out.Len() != 0 {
		t.Fatalf("late hook flooded history: %s", out.String())
	}
}

func testFreshDecision(at time.Time, title string) *model.Entry {
	return &model.Entry{
		ID: ids.NewEntry(at), Type: model.Decision, Title: title,
		Body: "decision body", Quote: "owner said it", UtteranceBy: model.ByUser, Confidence: 1,
		Source: model.Source{Kind: model.SrcHuman, Ref: "test:" + title, At: at},
	}
}
