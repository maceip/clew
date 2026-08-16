package main

import (
	"bytes"
	"testing"
)

func TestJournalNoteHelpCanNeverBecomeEntryText(t *testing.T) {
	for _, args := range [][]string{{"-h"}, {"--help"}, {"real text", "-h"}, {"real text", "--help"}} {
		var out bytes.Buffer
		note, help, err := parseJournalNoteArgs(args, &out)
		if err != nil || !help || note.Text != "" {
			t.Fatalf("args=%v note=%#v help=%v err=%v", args, note, help, err)
		}
	}
}

func TestJournalNoteParsesFlagsOnlyAfterText(t *testing.T) {
	var out bytes.Buffer
	note, help, err := parseJournalNoteArgs([]string{"exact words", "--type", "decision", "--title", "Chosen path", "--tags", "a/**,b"}, &out)
	if err != nil || help {
		t.Fatalf("help=%v err=%v", help, err)
	}
	if note.Text != "exact words" || note.Type != "decision" || note.Title != "Chosen path" || note.Tags != "a/**,b" {
		t.Fatalf("note=%#v", note)
	}
}
