AME = clew (owner's decision; supersedes the builder's unilateral restart from §12.1)

Rename the product end-to-end, then journal the decision. Mechanical sweep, verify, then record:

Sweep: go.mod module path · binary + cmd/ dir · every CLI string and help text · README (including install instructions) · the journal branch ref (restart/journal → NAME/journal, including the branch already created locally) · workspace dir (.restart/ → .NAME/) · state dir (~/.restart/ → ~/.NAME/) · env var prefixes · MCP server name · ntfy/webhook topic defaults · the CLAUDE.md/AGENTS.md snippet text that init installs · hook names · JOURNAL_SPEC.md working-name header and §12.1 (mark the decision CLOSED with the chosen name) · fixture references · test strings.
Verify: grep -ri restart . returns only deliberate historical mentions (spec history, this decision's record); full test suite + hermetic acceptance gates green after rename; fresh go install produces a working NAME on PATH.
Journal the decision as a decision entry with the reason recorded: chosen name, alternatives considered (restart — verb collision, names the crisis not the daily loop; lore — rejected for binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper considered), source: owner. This entry supersedes the builder's naming resolution in the README table.
Two pending one-liners while you're in there: add the header to fixtures/strategy-session/labels.yaml — "test fixture: humans edit this once, ever" (or move under testdata/); and add a real install path to the README quickstart (go install … — the owner hit command not found within a minute of first contact).
Commit the rename as its own commit, separate from anything else.
