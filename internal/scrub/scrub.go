// Package scrub redacts secrets from quotes/bodies before anything is
// written to the journal (JOURNAL_SPEC §6.2): a fixed key/token regex set
// plus an entropy check. Redactions render as ‹redacted› and are counted.
package scrub

import (
	"math"
	"regexp"
	"strings"
)

const Mark = "‹redacted›"

var patterns = []*regexp.Regexp{
	// Cloud/API key shapes.
	regexp.MustCompile(`\bAKIA[0-9A-Z]{16}\b`),                                              // AWS access key
	regexp.MustCompile(`\bASIA[0-9A-Z]{16}\b`),                                              // AWS STS
	regexp.MustCompile(`\bgh[pousr]_[A-Za-z0-9]{20,}\b`),                                    // GitHub tokens
	regexp.MustCompile(`\bgithub_pat_[A-Za-z0-9_]{20,}\b`),                                  // GitHub fine-grained
	regexp.MustCompile(`\bsk-[A-Za-z0-9_-]{20,}\b`),                                         // OpenAI-style
	regexp.MustCompile(`\bxox[baprs]-[A-Za-z0-9-]{10,}\b`),                                  // Slack
	regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{30,}\b`),                                        // Google API
	regexp.MustCompile(`\bglpat-[A-Za-z0-9_-]{20,}\b`),                                      // GitLab
	regexp.MustCompile(`\bnpm_[A-Za-z0-9]{30,}\b`),                                          // npm
	regexp.MustCompile(`\beyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\b`), // JWT
	// Assignments: password=..., api_key: "...", secret=...
	regexp.MustCompile(`(?i)\b(password|passwd|secret|api[_-]?key|access[_-]?token|auth[_-]?token|private[_-]?key)\b\s*[:=]\s*['"]?[^\s'"]{6,}`),
	// PEM blocks.
	regexp.MustCompile(`-----BEGIN [A-Z ]*PRIVATE KEY-----[\s\S]*?-----END [A-Z ]*PRIVATE KEY-----`),
	// Basic-auth URLs.
	regexp.MustCompile(`\b[a-z][a-z0-9+.-]*://[^/\s:@]+:[^@\s]+@`),
}

// Scrub replaces secret-shaped substrings and returns the count replaced.
func Scrub(s string) (string, int) {
	n := 0
	for _, re := range patterns {
		s = re.ReplaceAllStringFunc(s, func(string) string { n++; return Mark })
	}
	// Entropy pass: long unbroken token-ish runs with high Shannon entropy.
	words := tokenRe.FindAllStringIndex(s, -1)
	if len(words) > 0 {
		var b strings.Builder
		last := 0
		for _, w := range words {
			tok := s[w[0]:w[1]]
			if len(tok) >= 28 && entropy(tok) > 4.2 && !looksStructural(tok) {
				b.WriteString(s[last:w[0]])
				b.WriteString(Mark)
				n++
			} else {
				b.WriteString(s[last:w[1]])
			}
			last = w[1]
		}
		b.WriteString(s[last:])
		s = b.String()
	}
	return s, n
}

var tokenRe = regexp.MustCompile(`[A-Za-z0-9+/_=-]{28,}`)

// looksStructural spares hashes users legitimately discuss (git SHAs, ULIDs,
// paths) — lowercase-hex and crockford-upper runs are common in engineering
// conversation and carry no exploitable entropy pattern worth hiding.
func looksStructural(tok string) bool {
	if regexp.MustCompile(`^[0-9a-f]{28,64}$`).MatchString(tok) {
		return true // git sha / hex digest
	}
	if regexp.MustCompile(`^[0-9A-HJKMNP-TV-Z]{26,27}$`).MatchString(tok) {
		return true // ULID
	}
	return false
}

func entropy(s string) float64 {
	freq := map[rune]float64{}
	for _, r := range s {
		freq[r]++
	}
	e := 0.0
	n := float64(len(s))
	for _, c := range freq {
		p := c / n
		e -= p * math.Log2(p)
	}
	return e
}
