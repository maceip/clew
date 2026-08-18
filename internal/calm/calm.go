// Package calm translates machine and source vocabulary at human rendering
// boundaries. Stored journal bytes and agent instructions stay exact and hard.
package calm

import (
	"regexp"
	"strings"
)

type replacement struct {
	re *regexp.Regexp
	to string
}

var humanWords = []replacement{
	{regexp.MustCompile(`(?i)\bowner-laws?\b`), "owner guidance"},
	{regexp.MustCompile(`(?i)\bowner laws\b`), "owner guidance"},
	{regexp.MustCompile(`(?i)\bowner law\b`), "owner guidance"},
	{regexp.MustCompile(`(?i)\billegal\b`), "not allowed"},
	{regexp.MustCompile(`(?i)\blaws\b`), "standing decisions"},
	{regexp.MustCompile(`(?i)\blaw\b`), "standing decision"},
	{regexp.MustCompile(`(?i)\bstates\b`), "statuses"},
	{regexp.MustCompile(`(?i)\bstate\b`), "status"},
	{regexp.MustCompile(`(?i)\bviolations\b`), "problems"},
	{regexp.MustCompile(`(?i)\bviolation\b`), "problem"},
	{regexp.MustCompile(`(?i)\benforcement\b`), "follow-through"},
	{regexp.MustCompile(`(?i)\benforced\b`), "kept"},
	{regexp.MustCompile(`(?i)\benforces\b`), "keeps"},
	{regexp.MustCompile(`(?i)\benforce\b`), "keep"},
	{regexp.MustCompile(`(?i)\blegal\b`), "allowed"},
	{regexp.MustCompile(`(?i)\bpolicies\b`), "guidance"},
	{regexp.MustCompile(`(?i)\bpolicy\b`), "guidance"},
	{regexp.MustCompile(`(?i)\brules\b`), "agreements"},
	{regexp.MustCompile(`(?i)\brule\b`), "agreement"},
	{regexp.MustCompile(`(?i)\bgovernance\b`), "stewardship"},
	{regexp.MustCompile(`(?i)\bgoverned\b`), "guided"},
	{regexp.MustCompile(`(?i)\bgoverns\b`), "guides"},
	{regexp.MustCompile(`(?i)\bgovern\b`), "guide"},
	{regexp.MustCompile(`(?i)\bcompliance\b`), "follow-through"},
	{regexp.MustCompile(`(?i)\bconstitution\b`), "foundation"},
	{regexp.MustCompile(`(?i)\bcertification\b`), "approval"},
	{regexp.MustCompile(`(?i)\bcertified\b`), "approved"},
	{regexp.MustCompile(`(?i)\bcertify\b`), "approve"},
	{regexp.MustCompile(`(?i)\bgraveyard\b`), "retired work"},
	{regexp.MustCompile(`(?i)\bcontradicted\b`), "in conflict"},
	{regexp.MustCompile(`(?i)\bcontradictions\b`), "conflicts"},
	{regexp.MustCompile(`(?i)\bcontradiction\b`), "conflict"},
	{regexp.MustCompile(`(?i)\bsuperseded\b`), "replaced"},
	{regexp.MustCompile(`(?i)\btainted\b`), "from an untrusted source"},
	{regexp.MustCompile(`(?i)\btaint\b`), "source warning"},
	{regexp.MustCompile(`(?i)\bwithheld\b`), "kept out"},
	{regexp.MustCompile(`(?i)\bfailures\b`), "problems"},
	{regexp.MustCompile(`(?i)\bfailure\b`), "problem"},
	{regexp.MustCompile(`(?i)\balarms\b`), "attention items"},
	{regexp.MustCompile(`(?i)\balarm\b`), "attention item"},
	{regexp.MustCompile(`(?i)\bratified\b`), "agreed"},
	{regexp.MustCompile(`(?i)\bratify\b`), "agree"},
	{regexp.MustCompile(`(?i)\brulings\b`), "decisions"},
	{regexp.MustCompile(`(?i)\bruling\b`), "decision"},
	{regexp.MustCompile(`(?i)\bdispositions\b`), "choices"},
	{regexp.MustCompile(`(?i)\bdisposition\b`), "choice"},
	{regexp.MustCompile(`(?i)\badmission\b`), "acceptance"},
	{regexp.MustCompile(`(?i)\bpromotions\b`), "sharing"},
	{regexp.MustCompile(`(?i)\bpromotion\b`), "sharing"},
	{regexp.MustCompile(`(?i)\bpromoted\b`), "shared"},
}

// Text returns calm display prose without changing the stored source. Docket
// is intentionally absent from the replacement list: its name is approved.
func Text(s string) string {
	for _, word := range humanWords {
		s = word.re.ReplaceAllStringFunc(s, func(found string) string {
			to := word.to
			if found != strings.ToLower(found) && len(to) > 0 {
				to = strings.ToUpper(to[:1]) + to[1:]
			}
			return to
		})
	}
	return s
}
