package scrub

import (
	"strings"
	"testing"
)

func TestScrub(t *testing.T) {
	cases := []struct {
		in       string
		wantHits int
		keep     string
	}{
		{"my key is AKIAIOSFODNN7EXAMPLE ok", 1, "my key is"},
		{"token ghp_abcdefghijklmnopqrstuvwxyz123456", 1, "token"},
		{"set password = hunter2secret in env", 1, "in env"},
		{"the p95 was 340ms on the emulator", 0, "340ms"},
		{"commit 9c2e41f7aa31b2244a1c5f4de1e2ab3c4d5e6f70 fixed it", 0, "commit"},
		{"OPENAI sk-proj-Ab3dEfGh1jK2lMnOpQrStUvWx", 1, "OPENAI"},
	}
	for _, c := range cases {
		got, n := Scrub(c.in)
		if n != c.wantHits {
			t.Errorf("%q: hits %d want %d (→ %q)", c.in, n, c.wantHits, got)
		}
		if c.wantHits > 0 && !strings.Contains(got, Mark) {
			t.Errorf("%q: no redaction mark in %q", c.in, got)
		}
		if !strings.Contains(got, c.keep) {
			t.Errorf("%q: over-scrubbed to %q", c.in, got)
		}
	}
}

func TestScrubPEM(t *testing.T) {
	in := "here -----BEGIN RSA PRIVATE KEY-----\nMII...\n-----END RSA PRIVATE KEY----- done"
	got, n := Scrub(in)
	if n != 1 || strings.Contains(got, "MII") {
		t.Errorf("PEM not scrubbed: %q (%d)", got, n)
	}
}
