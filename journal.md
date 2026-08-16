# Journal

_generated 2026-08-16 09:27 UTC · 7 live entries (1 decisions · 5 findings · 1 questions · 0 intents) · 9 total in history_

## Decisions

### e01M04WCGJS9FS7FQB0YFX9DTYG — Name the system clew (owner decision)  `active`
> Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

_source: human cli:note · confidence: 1.00_

## Findings

### e01M04YAQN85TF6YQDP33VB4JQ0 — placeholder for write-path finding  `current`
> placeholder for write-path finding

placeholder for write-path finding

_source: human cli:note · confidence: 1.00_

### e01M04XVPW5H5JCPS92CFQ4EBY3 — Independent verify: clean clone green; differ/poller/manifest lack unit tests  `current`
> --- PASS: TestAcceptance1_AbsenceDetection / TestAcceptance2_ExtractionFidelityPipeline / TestAcceptance3_RestartRoundTrip; SKIP: TestAcceptance2_RealProvider

Second agent+machine, Go 1.26.3: build clean; full suite green; gates 1, 2-hermetic, 3 PASS; RealProvider SKIPs without keys; init push-deferral loud per spec. Gaps: no package tests for differ, poller, manifest, archaeology — coverage rides acceptance alone. Footgun: `journal note --help` ingested the flag as literal entry text; cleaned via reject.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 1.00 · tags: internal/**, acceptance/**_

### e01M04XVP9QBE041XCK43C78BVP — Decision-dense sessions live on uncovered surfaces; manual notes = homework  `current`
> ok not sure im liking this me needing to add to the journal each time a file in a repo is edited let me think about this

clew was designed in phone/cloud chats no watcher covers. An agent bridging that gap with prescribed manual notes made the product read as homework — an I1-violation smell, correctly caught. Non-discipline fixes: one export+backfill per key conversation; decisions echo into covered sessions, caught on first echo. §11 cloud-sensor trigger fires now for phone-heavy owners.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 0.95_

### e01M04XVNN32W4FTFKKSKJSN01V — varalys/lore owns session-provenance + git-sync; clew's edge is diff + absence  `current`
> lore is retrospective provenance — whole-session storage answering "why does this line exist." Restart is prospective state — distilled typed entries answering "what's true about this project right now"

varalys/lore ships session recording linked to commits, serverless git-remote sync, cross-tool memory over MCP; installs `lore` via brew/cargo. Same plumbing independently derived — commodity layer. Unoccupied, core to clew: typed distillation, intent×reality diff with absence, human steering surface, restart manifest. getlore.ai + Epic Lore make the name unownable.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 0.95_

### e01M04TDKJ79MEGBWTCSD8HTE5M — Lock-and-notary pitches fail for any team whose job isn't security — not just …  `current`
> Lock-and-notary pitches fail for any team whose job isn't security — not just solo devs. Evidence framing must lead with anti-forgetting/throughput, never audit.

Lock-and-notary pitches fail for any team whose job isn't security — not just solo devs. Evidence framing must lead with anti-forgetting/throughput, never audit.

_source: human cli:note · confidence: 1.00_

## Open questions

### e01M04XVQEEZ38J9TC5NZNKC16B — Run the live fidelity gate (RealProvider) on a machine with provider keys  `open`
> If you cannot reach P>=0.9 / R>=0.75, stop all other work and report that the kill criterion fired — that outcome is a valid and useful result, and it is written down on purpose

Hermetic pipeline passes, but the go/no-go gate — P≥0.9/R≥0.75 within 5 instruction iterations vs ratified labels — has never run against a live provider. Until it does, extraction quality (the stated existential risk) is unmeasured; the kill criterion is theoretical. Needs claude/codex CLI or OpenAI-compatible key; env flag was RESTART_FIDELITY=1 pre-rename — verify.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 1.00 · tags: acceptance/**_
