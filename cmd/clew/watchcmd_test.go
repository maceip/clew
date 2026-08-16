package main

import (
	"testing"
	"unicode/utf8"
)

func TestClipStrHonorsCharacterLimit(t *testing.T) {
	got := clipStr("one\ntwo three", 8)
	if got != "one two…" {
		t.Fatalf("clipStr() = %q, want %q", got, "one two…")
	}
	if n := utf8.RuneCountInString(got); n != 8 {
		t.Fatalf("clipStr() returned %d characters, want 8", n)
	}

	got = clipStr("ééé", 2)
	if got != "é…" {
		t.Fatalf("clipStr() split or miscounted Unicode: %q", got)
	}
}
