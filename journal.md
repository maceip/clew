# Journal

_generated 2026-08-16 17:16 UTC · 11 live entries (2 decisions · 8 findings · 1 questions · 0 intents) · 13 total in history_

## Decisions

### e01M05RHSWXDNR10P1PY8ERYA9S — Dogfood metrics predeclared; D0 snapshot recorded  `active`
> Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.

Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.

_source: human cli:note · confidence: 1.00_

### e01M04WCGJS9FS7FQB0YFX9DTYG — Name the system clew (owner decision)  `active`
> Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

_source: human cli:note · confidence: 1.00_

## Findings

### e01M05RPB2N6QDXN1MP31SB64B2 — Dogfood D0 historical-session storm is an I12 failure  `current`
> Dogfood failure D0: historical sessions misclassified live=33; observed tokens=12895847; overlaps=342; stomps=27; actual pushes=0; false pushed_at=27; extraction spend=0; adapter incidents=1; watcher stopped before extraction.

Dogfood failure D0: historical sessions misclassified live=33; observed tokens=12895847; overlaps=342; stomps=27; actual pushes=0; false pushed_at=27; extraction spend=0; adapter incidents=1; watcher stopped before extraction.

_source: human cli:note · confidence: 1.00_

### e01M05RF74PBCZVAPRY49GCVM7H — answer: Run the live fidelity gate (RealProvider) on a machine with provider ke…  `current`
> P=.91 R=.83; D=6/7 F=4/5; reject=0; claude; iter=1; PASS

P=.91 R=.83; D=6/7 F=4/5; reject=0; claude; iter=1; PASS

_source: human inbox:answer:e01M04XVQEEZ38J9TC5NZNKC16B · confidence: 1.00 · tags: acceptance/**_

### e01M05REDJZAERAWRQG7349E7Y1 — Live fidelity gate passed on iteration 1  `current`
> Live fidelity gate iteration 1: P=0.91; R=0.83; decisions=6/7; findings=4/5; rejected=0; provider=claude; PASS.

Live fidelity gate iteration 1: P=0.91; R=0.83; decisions=6/7; findings=4/5; rejected=0; provider=claude; PASS.

_source: human cli:note · confidence: 1.00_

### e01M04YAQN85TF6YQDP33VB4JQ0 — Foreign agents can read the journal but not write; contribution path unbuilt  `current`
> so if you cant even push to the public clew directory does that mean our first real test failed?

cursor[bot] was denied push to maceip/clew — correct posture: journal write = repo write, else anyone could poison context.md (§6.5 amplifier). But there is no sanctioned path for non-credentialed contributors; tonight's delivery was a hand-rolled bundle. Options: document fork-PR onto clew/journal, or a `clew import <bundle>` verb landing entries pending human confirm.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 1.00_

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

### e01M04XVQEEZ38J9TC5NZNKC16B — Run the live fidelity gate (RealProvider) on a machine with provider keys  `answered`
> If you cannot reach P>=0.9 / R>=0.75, stop all other work and report that the kill criterion fired — that outcome is a valid and useful result, and it is written down on purpose

Hermetic pipeline passes, but the go/no-go gate — P≥0.9/R≥0.75 within 5 instruction iterations vs ratified labels — has never run against a live provider. Until it does, extraction quality (the stated existential risk) is unmeasured; the kill criterion is theoretical. Needs claude/codex CLI or OpenAI-compatible key; env flag was RESTART_FIDELITY=1 pre-rename — verify.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 1.00 · tags: acceptance/**_
