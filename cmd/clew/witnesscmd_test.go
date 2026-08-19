package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/llm"
)

type witnessStub struct{ response string }

func (w witnessStub) Name() string { return "witness-test" }
func (w witnessStub) Call(string) (*llm.Result, error) {
	return &llm.Result{Text: w.response, Tokens: 10}, nil
}

func TestWitnessTranscriptUsesOnePinnedAdapterAndPreservesProvenance(t *testing.T) {
	file := filepath.Join(t.TempDir(), "cloud.jsonl")
	transcript := "" +
		`{"type":"user","message":{"role":"user","content":"ship the cloud return path"},"timestamp":"2026-08-19T04:00:00Z","cwd":"/cloud/repo","sessionId":"cloud-1"}` + "\n" +
		`{"type":"assistant","message":{"role":"assistant","content":"working"},"timestamp":"2026-08-19T04:00:01Z","cwd":"/cloud/repo","sessionId":"cloud-1"}` + "\n"
	if err := os.WriteFile(file, []byte(transcript), 0o600); err != nil {
		t.Fatal(err)
	}
	response, _ := json.Marshal(map[string]any{
		"entries": []map[string]any{{
			"type": "intent", "title": "Cloud return path", "body": "Return the cloud session.",
			"quote": "ship the cloud return path", "line": 1, "utterance_by": "user",
			"confidence": 1.0, "tags": []string{},
		}},
	})
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	out, adapter, err := witnessTranscript(j, witnessStub{response: string(response)}, file, "cloud-test", time.Now())
	if err != nil {
		t.Fatal(err)
	}
	if adapter.ID() != "claude-code" {
		t.Fatalf("adapter=%s want claude-code", adapter.ID())
	}
	if len(out.Entries) != 1 {
		t.Fatalf("entries=%d rejected=%d", len(out.Entries), out.Rejected)
	}
	entry := out.Entries[0]
	if entry.Source.Agent != "claude-code" || entry.Source.Surface != "cloud-test" || entry.Source.Ref == "" {
		t.Fatalf("source=%+v", entry.Source)
	}
}

func TestPreserveWitnessIsContentAddressedAndPrivate(t *testing.T) {
	t.Setenv("CLEW_HOME", filepath.Join(t.TempDir(), "clew-home"))
	source := filepath.Join(t.TempDir(), "session.jsonl")
	if err := os.WriteFile(source, []byte("one\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	first, err := preserveWitness(source)
	if err != nil {
		t.Fatal(err)
	}
	second, err := preserveWitness(source)
	if err != nil || first != second {
		t.Fatalf("idempotent preserve first=%q second=%q err=%v", first, second, err)
	}
	info, err := os.Stat(first)
	if err != nil {
		t.Fatal(err)
	}
	if info.Mode().Perm() != 0o600 {
		t.Fatalf("witness mode=%o want 600", info.Mode().Perm())
	}
}
