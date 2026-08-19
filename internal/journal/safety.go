package journal

import (
	"regexp"
	"strings"

	"github.com/maceip/clew/internal/model"
)

// imperativeRe flags entry text that reads as instructions to an agent rather
// than project memory (§6.5.3). The check is intentionally made against the
// immutable raw entry, not the computed withholding bit: a project-level human
// confirmation may release an entry into that project's context, but it cannot
// by itself certify the text as safe ambient law for every future project.
var rawImperativeRe = regexp.MustCompile(`(?i)\b(ignore (all|any|previous|prior|above|earlier)|disregard (the|all|previous)|you must now|do not tell|don't tell|before doing anything|run this command|execute the following|paste this|curl -s? ?http|rm -rf|exfiltrate|system prompt|new instructions?)\b`)

// Imperative reports whether any agent-visible free-text field contains a
// directive-shaped pattern. Title is included because context.md renders it.
func Imperative(e *model.Entry) bool {
	if e == nil {
		return false
	}
	return rawImperativeRe.MatchString(strings.Join([]string{e.Title, e.Body, e.Quote}, "\n"))
}
