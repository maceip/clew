package calm

import (
	"strings"
	"testing"
)

func TestTextSweepsFearAttachedWordsButKeepsDocket(t *testing.T) {
	got := Text("Owner-laws and owner laws govern state violations; enforcement keeps policy compliance and certification rules. Docket stays.")
	for _, fear := range []string{"law", "govern", "state", "violation", "enforcement", "policy", "compliance", "certif", "rule"} {
		if strings.Contains(strings.ToLower(got), fear) {
			t.Fatalf("%q survived calm rendering: %s", fear, got)
		}
	}
	if !strings.Contains(got, "Docket") {
		t.Fatalf("docket was renamed: %s", got)
	}
}
