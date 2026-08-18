package journal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestLoadForDisplayRecoversExactUnquotedColonAndStaysLoud(t *testing.T) {
	dir := t.TempDir()
	entries := filepath.Join(dir, "entries")
	if err := os.MkdirAll(entries, 0o755); err != nil {
		t.Fatal(err)
	}
	id := "e01M0AXCM561K3C5QXAVGVGT46T"
	body := strings.Repeat("exact source words ", 30)
	raw := "id: " + id + "\n" +
		"type: intent\n" +
		"title: Build the freshness ladder: one delta payload, five delivery layers\n" +
		"body: >-\n  " + body + "\n" +
		"quote: owner words\n" +
		"utterance_by: user\n" +
		"source:\n  kind: session\n  ref: session:test\n  at: 2026-08-18T17:00:00Z\n" +
		"confidence: 1\ntags: []\nenv: null\naffects: []\n"
	if err := os.WriteFile(filepath.Join(entries, id+".yaml"), []byte(raw), 0o644); err != nil {
		t.Fatal(err)
	}

	strict, err := Load(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(strict.LoadErrors) != 1 || strict.Entries[id] != nil {
		t.Fatalf("strict load guessed malformed source: errors=%v entry=%#v", strict.LoadErrors, strict.Entries[id])
	}

	display, err := LoadForDisplay(dir)
	if err != nil {
		t.Fatal(err)
	}
	entry := display.Entries[id]
	if entry == nil {
		t.Fatalf("display load did not recover exact source: errors=%v recoveries=%v", display.LoadErrors, display.DisplayRecoveries)
	}
	if entry.Title != "Build the freshness ladder: one delta payload, five delivery layers" || !strings.HasPrefix(entry.Body, "exact source words") {
		t.Fatalf("display recovery changed source: %#v", entry)
	}
	if len(display.DisplayRecoveries) != 1 {
		t.Fatalf("recovery was not loud: %#v", display.DisplayRecoveries)
	}
}
