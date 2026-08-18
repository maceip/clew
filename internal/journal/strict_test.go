package journal

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadKeepsMalformedEntryLoud(t *testing.T) {
	dir := t.TempDir()
	entries := filepath.Join(dir, "entries")
	if err := os.MkdirAll(entries, 0o755); err != nil {
		t.Fatal(err)
	}
	id := "e01M0AXCM561K3C5QXAVGVGT46T"
	raw := "id: " + id + "\n" +
		"type: intent\n" +
		"title: Build the freshness ladder: one delta payload, five delivery layers\n" +
		"quote: owner words\n" +
		"utterance_by: user\n" +
		"source:\n  kind: session\n  ref: session:test\n  at: 2026-08-18T17:00:00Z\n" +
		"confidence: 1\ntags: []\nenv: null\naffects: []\n"
	if err := os.WriteFile(filepath.Join(entries, id+".yaml"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	got, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(got.LoadErrors) != 1 || got.Entries[id] != nil {
		t.Fatalf("malformed source was guessed or hidden: errors=%v entry=%#v", got.LoadErrors, got.Entries[id])
	}
}
