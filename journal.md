# Journal

_generated 2026-08-19 01:20 UTC · 152 live entries (74 decisions · 47 findings · 9 questions · 22 intents) · 174 total in history_

## DECIDED

- `e01M0BFY164YEV6DAEFVGGH18VT` The limiter gates distillation timing, never sensing; failure is lag — 2h · active
- `e01M0BFRYY2WVHJZ3R3TDV4CFTS` The wording sweep covers every fear-attached word; docket stays by name — 2h · active
- `e01M0BFRYY1EW9CFMJNDJ5QH3M2` Evidence settles merge lines; apply is never asked for finished work — 2h · active
- `e01M0BF0WMX264RA9D0VTM9R24K` Entry ids are machine plumbing: never shown to or relayed through humans — 3h · active
- `e01M0BEYP65CE70G0VVSX3PV01B` Finished means shared: work ends pushed or PR'd; local-only is an alarm — 3h · active
- `e01M0BETRRAZ54063PRJK1JSQS7` The finish message is a surface: what exists, where it lives, my next move — 3h · active
- `e01M0BER1Q412DRFQYESPCN0Q30` Lines are plain speech, no ids; near-duplicates fold; held items rest — 3h · active

## LEARNED

- `e01M0AYE066QK08QK8MTPX4XNFX` Codex finished I13 stale: tree uncommitted, law wording on human surfaces — 7h · current
- `e01M0AXCM55N0QM9RCRYF48TQ6C` Universal injection point: every model API call rebuilds the mind — 8h · current
- `e01M0AXCM54CDSPV84DBH1PWGWD` Spec nudge matrix is stale: codex and gemini now ship injection hooks — 8h · current
- `e01M0AV5721XB0ZFQW42241YJ4A` Daemon fallback birth discarded the triggering session's first prompt — 8h · current
- `e01M0ATNY3H0WHX91ASD0AMR8TE` Cold CLEW_HOME loses concurrent births; warm machine state is safe — 9h · current
- `e01M0ATNY3H0WHX91ASCWE1MM7R` Repo identity is the absolute path, so a rebuilt repo at a reused path is not a… — 9h · current
- `e01M0ATARAN8NWCBGNTM58QRVPX` Promotion candidates enter project context before the human rules on them — 9h · current

## OPEN

- `e01M0BER1Q5W312NDD7GZ283RCM` Cap the two screens at 7+7; overflow policy is an open question — 3h · open ★
- `e01M0AY97K0AZHNJW0D87KTD7YH` Two-grade telling: block dependent tasks on unresolved drift? — 8h · open ★
- `e01M0ASQGM2P07Z36KVZX6P8EH4` Adopt the complexity law: additions must be a verb, label, rendering, or config — 9h · open ★
- `e01M0ASQGKZ9ZY4XRF6K7KRCT75` Adopt clew witness <transcript> as the cloud-session gap fix? — 9h · open ★
- `e01M0ASQGKX2F9THW496N79Z83F` Approve selfwatch + journal add + owner-laws relocation to a git-reachable repo — 9h · open ★
- `e01M0ASQGKVW3RPZ2M51J3WQ9G9` Approve cloud env recipes: install clew + wire MCP in cursor/codex/claude envs — 9h · open ★
- `e01M0ASQGKRQBVSGHP3RCYNWYEG` Approve extractor over-firing tune: rollup >32KB and docket hit 8 cards — 9h · open ★

## ALERTS

- **ABSENT** `e01M0AQHSRF4DVDYZ989PVHGA7R` Birth detection: auto-init a new repo with owner laws only — 9h
- **ABSENT** `e01M0ARJE1XNN8Q45DJ36FP47YT` Surface coverage: repo-write cloud agents (Cursor-class) are full journal nodes — 9h
- **ABSENT** `e01M0ARJFTWEWY5H6JFJ17656W4` Surface coverage: PR-only cloud agents (Codex-app-class) contribute knowledge — 9h
- **ABSENT** `e01M0ARJHKDP6Z6R1FKZSJ8AN4S` Surface coverage: laptop agents fully sensed with zero human effort — 9h
- **ABSENT** `e01M0ARJKGSQFSH8WSVZEG520DH` Surface coverage: phone reads the glance and receives decision cards — 9h
- **ABSENT** `e01M0ASDMN8S3TNV9502VTFK27E` Wire the seed/lineage libraries into watcher, materialization, and clew from — 9h
- possible-contradiction `e01M05SA72DPRGNTY7GCN1P7CED` Rename the inbox surface to "docket"; keep inbox as hidden alias — 2d

## Intent × reality

| Intent | Age | Reality | State |
|---|---:|---:|---|
| `e01M0AZN6HJETV241AK5RSBDHNR` Held: a restart tab — stage selected drift into the next generation | 7h | 0 evidence | proposed |
| `e01M0AXCM561K3C5QXAVGVGT46T` Build the freshness ladder: one delta payload, five delivery layers | 8h | 0 evidence | proposed |
| `e01M0AST7PZRJP8XBWDAHWK6QNV` Build invariant I13: ambient seed, birth detection, owner laws, clew from | 9h | 1 evidence | in_flight |
| `e01M0ASDMN8S3TNV9502VTFK27E` Wire the seed/lineage libraries into watcher, materialization, and clew from | 9h | 0 evidence | **ABSENT** |
| `e01M0ARJKGSQFSH8WSVZEG520DH` Surface coverage: phone reads the glance and receives decision cards | 9h | 0 evidence | **ABSENT** |
| `e01M0ARJHKDP6Z6R1FKZSJ8AN4S` Surface coverage: laptop agents fully sensed with zero human effort | 9h | 0 evidence | **ABSENT** |
| `e01M0ARJFTWEWY5H6JFJ17656W4` Surface coverage: PR-only cloud agents (Codex-app-class) contribute knowledge | 9h | 0 evidence | **ABSENT** |

## Decisions

### e01M0BFY164YEV6DAEFVGGH18VT — The limiter gates distillation timing, never sensing; failure is lag  `active`
> a deaf agent is useless and will lead to more problems for me to deal with later?

Owner challenge: a deaf agent is useless. Purpose on record: the limiter is not cost control — it protects shared rate limits and guards against runaway loops. Corrected design: sensing (tailing, recording) is free and never stops; only distillation may lag under pressure, shown as memory is N minutes behind, catching up when headroom returns. Deafness is impossible; nothing goes unrecorded.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BFRYY2WVHJZ3R3TDV4CFTS — The wording sweep covers every fear-attached word; docket stays by name  `active`
> why is 7 indexing on "law" when we discussed _all_ legal sounding words, with docket being ok

Owner correction: the rename card narrowed to the single word law. The sweep is all fear-attached words wherever humans read — law, state, violation, enforcement and relatives — judged per word in context. Docket is explicitly approved and stays. Agent-facing hard register remains untouched where hardness prevents wiggle.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BFRYY1EW9CFMJNDJ5QH3M2 — Evidence settles merge lines; apply is never asked for finished work  `active`
> why are there still 7 merges that need to take place ?  number 1 i thought you already merged?

Owner found the merge asking him to bless work already built, tested, and pushed. Rule: the merge joins decisions to evidence; a decision whose demanded work is evidenced settles itself and shows once as settled-while-away. Apply is reserved for work not yet done or judgment only a human can make. Nothing auto-acts on the repo; settling is status computation, not action.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BF0WMX264RA9D0VTM9R24K — Entry ids are machine plumbing: never shown to or relayed through humans  `active`
> you should never print "e01M0BER1Q412DRFQYESPCN0Q30" i dont know what it means and its super long -- what is the intent of showing that to the human?

Extension of plain-speech: the cloud agent printed raw ids in receipts and in prompts the human had to copy. The intent was verifiability — but verification is machine work, and agents holding the journal resolve plain words better than opaque codes. Rule: ids live in files, commits, and machine channels only; humans see words; agents are addressed in words and resolve entries themselves.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BEYP65CE70G0VVSX3PV01B — Finished means shared: work ends pushed or PR'd; local-only is an alarm  `active`
> why do i need to tell it to push? with all the journals and dockets and intents and knowledge that should never be the case?

Owner ruling closing the push gap: a task is not finished until the work is shared per repo convention — pushed to the branch or opened as a PR. Committed-but-local is an alarm state the finish message must name, never a resting state. Root cause on record: this norm lived in behavior and was never spoken, so the memory had nothing to inject; said once here, it reaches every agent forever.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BETRRAZ54063PRJK1JSQS7 — The finish message is a surface: what exists, where it lives, my next move  `active`
> when i read that my first thought is: "ok it didnt push anything so thats not good", then i dont know what the rest of it means

Owner correction: codex signed off in builder frame (Nothing was pushed. No real apply...) — accurate, meaningless to the human, and alarming: didn't-push read as failure. Rule: the closing utterance speaks the human frame in plain words — what exists now, whether it is safely shared or local-only, what the human can say next — then shows the two screens. Compliance detail lives behind explain.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BER1Q412DRFQYESPCN0Q30 — Lines are plain speech, no ids; near-duplicates fold; held items rest  `active`
> language in each <knowledge merge/intent gap> thats confusing, e.g., "Let cloud agents that can only"

Owner corrections from the first real run: rendered lines confused (Let cloud agents that can only...) and full entry ids burned attention. Rules: plain spoken English, subject-first, one breath; no ids or codes on lines — identity lives behind explain; near-duplicates fold to one line; held-for-owner entries appear in no actionable list. The amnesia test stays the floor; this adds plain speech.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0BEM04PAZ85R5YNRM2Y31Z6 — Broken states carry their verb: no unactionable warnings for humans  `active`
> that is state the human cant immediately fix/help with

Owner correction from the first real merge/gap run: could-not-check lines handed the human a problem with no action. Rule: a broken state shown to a human must carry its fix verb (usually hand to the attending agent) or name who is already fixing it; problems only machinery can fix route to agents, never to human eyes. Earned silence stands — broken arrives actionable, like everything else.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0AZDDYC49Y1741H7W74Y1QY — Second tab: the intent gap — everything intended, not yet real  `active`
> the intent gap is a similar thing that lists simply all the crazy shit thats not yet implemented

Owner design: next to the knowledge merge sits the intent gap — same glanceable, amnesia-proof list shape, listing intents with no evidence in reality (the absence machinery gets its human surface). Verbs: build (hand to the idle agent), explain (live), retire (a conscious no, kept with provenance). It converts forgetting into deciding — stratura's unbuilt core would have topped it for weeks.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0AZ7JBPJRNHFJQC6WEEQB9Y — Silence is the signal: an absent merge means truly nothing new  `active`
> if there is no knowledge diff shown, the human knows nothing new has been added somewhere else

Owner property: when no knowledge diff appears, the human may trust that nothing new landed anywhere — silence is the all-caught-up signal. For that trust to hold, silence must be earned: a broken watcher, stale sync, or failed check must announce itself distinctly and can never render as an empty diff. Quiet means verified-quiet. Nothing-new and could-not-check are never the same screen.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AZ4BGF1VC0VSXA05VYVEQ3 — Explain is live: the attending agent reads the entry and explains  `active`
> so you dont need to encode the "explain in more detail in the kowledge diff" itself

Refinement: the merge diff encodes nothing but lines, entry ids, and verbs. Pressing explain hands the entry to the agent already present at the finish boundary — it reads the journal, quotes the owner's words, and explains what applying means for the work at hand, answering follow-ups conversationally. clew stays the bookkeeper of see-once and defer state; the agent is the explain engine.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AZ0K17AA5D9P7KZDSPJQSY — Merge lines must pass the amnesia test; verbs are apply/explain/defer  `active`
> these diff entries need to be something the human can read after maybe a day where he totally forgot the convo he had with you

Amendment: each merge line must be readable by a human who forgot the conversation entirely — references glossed inline (the five promises appear as five words), machinery nouns translated, no dangling 'the budget'. Per line: apply, explain (prints body + the owner's verbatim quote + link, then re-offers), defer. Footer gains apply-all. Explain works because one's own words restore memory.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AYQ3VJM71SMM1YQFYE1X61 — Knowledge Merge at finish: glanceable apply/defer list, external memory  `active`
> this list cannot be verbose, it needs to be glancable, and it serves the "humans forget" thing as well

Owner design: at finish, one colored glanceable list — top unapplied changes (code, intent, knowledge), one line each, entry-linked, one-keystroke apply/defer. See-once by decision id; defer compresses to a nagging count, never re-shown as new. External memory for a forgetful human: recognition over recall, per the HCI findings. It is the docket rendered at the finish boundary.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0AY6VWJV133F811JYAANJPE — Stale finish: know and tell, never act — the reconcile is read-only  `active`
> finishes, then syncs, should _NOT_ then remove the work it just did, or do anything for that matter

Amendment to the finish check-in: an agent whose finished work was obsoleted mid-flight must not remove, redo, or touch anything on its own — no action of any kind. The check-in only installs knowledge: at the next human prompt it must say the work is deprecated/obsoleted/wrong and why, unless the human already resolved it elsewhere. Interpretation is the human's call. Owner will stress-test.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AXZTTHG5FKXETX0X8PR6EX — On finish, check in first: reconcile against the delta before next steps  `active`
> but what i do expect is that after it finishes, to check in first and then figure out what to do based on the new synced state

Owner ruling completing the task lifecycle: when an agent finishes, it must sync the journal and reconcile its output against decisions that landed mid-flight before concluding or picking next steps — state contradictions explicitly, then close. Stop/AfterAgent hooks make this enforceable on claude/codex/gemini; elsewhere it is convention plus the glance flagging stale finishes.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AXZ5DEXRWRGD3GT1HYEYSD — Redaction and secret handling are deferred: no monkey-patching without team inp…  `possible-contradiction`
> that requires more input and other subsystems than we want you to try to monkey patch without inputs from our team

The owner instructed stopping the proposed redaction/secret-handling follow-up. That subsystem needs inputs and other subsystems beyond what should be patched unilaterally. Current journal PATs are 24-hour ephemeral credentials, intentionally pushable, and must not block I13 or trigger redaction machinery.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-22-35-01a01577-9dd5-7e63-811f-0b52ba3e3cab.jsonl#L2925 · confidence: 0.94 · tags: internal/gitx/**, internal/owner/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR_

### e01M0AXXKMNNKKY721HJ3REN3KH — Freshness is owed at contact points; a task runs on its snapshot  `active`
> i dont expect an agent on task to stop mid task and change based on a cloud agent decision i made at the same time

Owner refinement: a running task is never interrupted or mutated by concurrent decisions — it finishes on the snapshot it started with. Currency is owed at the next human contact: a message typed after returning lands on a mind that already has the delta. Hooks fire at that boundary; the proxy injects only on a new human message. Urgent items route to the human, who may stop the task.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AXRKQ8C7FZNKARW83CMBMX — The five promises are the foundation (owner ratified)  `active`
> yeah those 5 things are the foundation, agreed

Owner ratified the product's entire human-facing surface: (1) it remembers what we decide; (2) every agent starts already knowing it; (3) you can look up and see; (4) it taps your shoulder only when something needs you; (5) starting over loses nothing. Every feature must file under exactly one promise or it is not built. Vocabulary beyond these is machinery, never surfaced.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AXMXK3SAATFDFAYZ932TPC — Two registers, one memory: calm words for humans, hard words for agents  `active`
> however -- if using softer language means agents would "wiggle" out of those constructs, then we keep the harsher language

Owner ruling: human-facing vocabulary must avoid fear-attached words (law, state, violation). But register is a rendering choice, not a softening of the contract: wherever soft words would let agents wiggle out of a constraint, the agent-facing rendering keeps the harsh form. Same entries, two renderings; hardness is judged by compliance, softness by human calm.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · evidence: 1_

### e01M0AXF5FM9RDFP817NF7QW5BN — Human-facing surface must reduce to the desires it satisfies  `active`
> its already too much to understand given how it should reduce down to the simple set of human desires it satisfies

Owner ruling: the internal design may be intricate, but the human-visible vocabulary must collapse to the desire set: it remembers; every agent starts knowing; I can look up and see; it asks me only when it must; starting over loses nothing. Any feature that cannot be filed under one of these is out. The agent carries the machinery; the human carries five sentences.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AXCM53D36PDN0H85S0WE3W — Mind-plane freshness is vendor-neutral; hooks accelerate, never carry it  `active`
> this should work if i decide to spin up ollama locally with deepseek4-flash on my laptop

Owner ruling: a returning human must land on agents that already know the recent journal, and this must hold for a bare ollama model with no hook system. Per-vendor hooks may improve latency but the floor must work for anything that emits model API calls. Bar: current-at-next-interaction, zero human homework, any harness — including ollama running deepseek4-flash on a laptop.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AWM6188V135JTTP4S3MQAA — Incarnation reset needs all three facts; any Clew marker protects the repo  `possible-contradiction`
> stale machine state is erased only when the path is still registered, its journal binding is invalid, and the current `.git` contains no Clew incarnation evidence. A journal branch or either local Clew marker protects a damaged-but-existing repository from cleanup.

Disposable machine state for a reused path is wiped only when the path is still registered, the journal binding is invalid, AND the current .git holds no Clew incarnation evidence (clew.birth-ready, clew.journal-id, or the local journal branch). Any one marker means a damaged-but-existing repository gets repaired, not erased. The reset runs in a single SQLite transaction and leaves the moved pred…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-30-01a0159d-11bd-7b13-81c5-26b7f910b998.jsonl#L549 · confidence: 0.90 · tags: cmd/clew/**, internal/state/** · pairs-with: e01M05V5HWA6TFT0A0KDZY8S45K, e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0ATNY3H0WHX91ASD478G9E7_

### e01M0AVXA7EH2KD1BPZ4GJNKN67 — Witness-node role adopted: always-on ear with owner API creds, sequenced  `active`
> that sgtm with one wrinkle still in my brain, how does this system work for 2,10,100 projects

Owner sgtm: one always-on clew node (owner infra) whose sensors are API pollers with owner account creds — witnesses cursor/codex cloud sessions live with zero agent cooperation, runs extraction centrally, sole writer of projections (kills that conflict class). Git stays the only required meeting point; degrade-to-baseline preserved. Build gated behind single-lineage sequencing.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AVWK3HH2R9M55FQSAWZFF1 — I9 frugality replaced: listening completeness is the invariant, cost is a dial  `active`
> you assume i'm price sensitive and are using token cost as being prohibitive but ive never mentioned we need to make this work cheaply

Owner ruling: price sensitivity was an agent assumption, never stated. Replace the 2% ratio with an owner-set ceiling plus a hard floor above the largest atomic request; extraction must never deadlock. Spend stays a visible meter. This also resolves the URGENT budget card's direction.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AVAR0N8ZCCN1VJW3GQ5PF4 — Restart-with-mutation is the flagship advertised workflow, not a failure mode  `active`
> i allow myself to clew from clew "but with cloud agent witnesses" ... something we _actually_ advertise as a path

Owner ruling: restart-with-mutation is the flagship, advertised workflow. The old negative was retelling pain (re-briefing a blank agent), not rebirth. Direction: clew from <parent> "<mutation>" carries the seed, makes the mutation the genesis charter, and flags carried entries it contradicts — day-zero docket = the design debate, pre- computed. Gated behind single-lineage sequencing.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AV0H7T9P69CNPB56MRAG8V — Sequencing: single-lineage from + one-project glance FIRST; fleet/multi later  `active`
> we need to get clew from working well with just one project lineage, and the glance infra ui working with just one project well first, else we risk scope creep and me clew from clew

Owner ruling: no rush on multi-parent from or the fleet view. Get clew from working well with one project lineage and the glance UI working well for one project before any scaling work — else scope creep triggers the restart urge (owner: risk of clew- from-clew). Multi-parent and fleet rulings stand as destination, gated on the single versions working well.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0AV0H7RY7DG79VCMHPEMPJP — Glance direction ruling: graphic, two zooms — deferred behind single-project  `active`
> the glance view _graphic_ not text, needs to support a one project view, and global view

Owner direction: the glance becomes a graphic (project view: status-colored intent tiles, drift strip, docket badge; fleet view: hot-project tiles, dormant collapsed). Static self-contained HTML, no server. DEFERRED: build only after the single-project version works well.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0ATYJG615JE6BV5MG5RAF9Z — clew from must accept multiple parent projects, with strand selection  `active`
> clew from needs to support _multiple_ projects as inputs, the human may want parts from multiple

Owner ruling: inheritance is multi-parent. `clew from A B` unions seeds; `--tags <globs>` selects strands per parent; runnable repeatedly. Each carried entry keeps per-parent provenance; disagreements between parents surface as possible-contradiction cards for human arbitration, never silent merge. Genesis records multiple lineage links (the forest gains merge nodes).

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0ATNY3H0WHX91ASD478G9E7 — SessionStart consumes the carry-kit; the watcher prepares it  `possible-contradiction`
> the watcher should have prepared the carry-kit before the session, while `SessionStart` should primarily consume it

An already-registered session must print the atomic .clew/context.md immediately and return, instead of taking the birth lock and running git sync, journal reload, owner load and full materialization on every `claude` startup. Reason: startup would otherwise block on network/worktree work, a transient error would suppress a perfectly valid existing context, and it blurs the boundary where ambient…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-30-01a0159d-11bd-7b13-81c5-26b7f910b998.jsonl#L344 · confidence: 0.85 · tags: cmd/clew/** · pairs-with: e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0AST7PZRJP8XBWDAHVVX8NR — Owner-law layer: promoted findings, <=1KB injected into every project forever  `possible-contradiction`
> findings promoted via clew journal promote <id>; extractor proposes promotion when a finding's content is project-agnostic; ≤1KB injection budget into every project's context, forever.

Owner laws live in an owner-scope journal synced like any other. A finding becomes a law via clew journal promote <id>; the extractor proposes promotion when a finding's content is project-agnostic. The rendered law block is capped at a 1KB injection budget and enters every project's context indefinitely.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L9 · confidence: 0.93 · tags: internal/owner/**, internal/extract/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0AST7PZRJP8XBWDAGG7AD44 — SEED.md is maintained continuously by the watcher, never built on demand  `possible-contradiction`
> the watcher maintains SEED.md continuously alongside context.md — regenerated on journal change, never on demand. The carry-kit always predates the urge to restart.

The watcher regenerates SEED.md alongside context.md on every journal change rather than when someone asks for it, so the carry-kit always exists before the urge to restart appears. Restart never waits on seed generation.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L9 · confidence: 0.94 · tags: cmd/clew/**, internal/materialize/** · pairs-with: e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AR79QQ9PPSZKQ2EVTFNQMV, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0AST7PZRJP8XBWDAF0AMC15 — Birth detection auto-inits a new repo with owner laws only, no lore  `possible-contradiction`
> new git init + agent session on a watched machine → auto-init (watch, journal branch, context.md) with owner laws only. No lore, no card required.

A fresh git init plus an agent session on a watched machine triggers auto-init (watch, journal branch, context.md) carrying owner laws only — no lore, no card required. Acceptance test: mkdir x && git init && claude yields a context containing the owner's laws with zero clew commands typed.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L9 · confidence: 0.94 · tags: cmd/clew/** · pairs-with: e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0AST7PZRJP8XBWDABXXWWQH — Lineage is explicit, laws are ambient: a wrong guess poisons a fresh project  `possible-contradiction`
> lore inheritance was made explicit because a wrong lineage guess poisons a fresh project worse than no inheritance at all — laws are safe to auto-carry precisely because promotion certified them project-agnostic.

clew from <repo> is the only way lore crosses projects — runnable at birth or any time later, never automatic; a birth card may suggest it on blatant name/topic overlap but only suggests, never acts. Owner laws auto-carry instead, because promotion already certified them project-agnostic. Reason: a wrong lineage guess poisons a fresh project worse than no inheritance at all.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L9 · confidence: 0.95 · tags: cmd/clew/**, internal/owner/** · evidence: 1 · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0AST4FW4WFMMX225XEXGBD1 — Seed is ambient: watcher regenerates SEED.md on journal change, never on demand  `possible-contradiction`
> the watcher maintains SEED.md continuously alongside context.md — regenerated on journal change, never on demand. The carry-kit always predates the urge to restart.

The carry-kit is maintained continuously by the watcher next to context.md, rebuilt whenever the journal changes rather than assembled when someone asks to restart. Reason: the seed must already exist before the urge to restart arrives, so no deliberate ceremony is needed to preserve a dying project's lessons.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-30-01a0159d-11bd-7b13-81c5-26b7f910b998.jsonl#L9 · confidence: 0.90 · tags: cmd/clew/**, .clew/** · pairs-with: e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AQHSRF4DVDYZ989M4DVHYX, e01M0AQHSRF4DVDYZ989MTZQV7E, e01M0AQHSRF4DVDYZ989W6K185N, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0AST4FW4WFMMX225TQ56V5E — I13: laws are ambient and auto-carried; lineage must be declared explicitly  `possible-contradiction`
> lore inheritance was made explicit because a wrong lineage guess poisons a fresh project worse than no inheritance at all — laws are safe to auto-carry precisely because promotion certified them project-agnostic.

Under invariant I13 birth costs nothing: owner laws (promoted, project-agnostic findings) are injected into every new project automatically, but project lore is never inherited automatically. Lore travels only through the explicit `clew from <repo>` command. Reason: a wrong lineage guess poisons a fresh project worse than no inheritance at all, while promotion has already certified laws as safe a…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-30-01a0159d-11bd-7b13-81c5-26b7f910b998.jsonl#L9 · confidence: 0.95 · tags: cmd/clew/**, .clew/** · evidence: 1 · pairs-with: e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AQHSRF4DVDYZ989M4DVHYX, e01M0AQHSRF4DVDYZ989MTZQV7E, e01M0AQHSRF4DVDYZ989W6K185N, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0ASQGK4FMKRNCNR91KJ06JD — Restart machinery must be zero human effort: ambient inheritance, opt-out  `active`
> the solution we create out of clew needs to make the restart acceleration zero effort from the human or zero cognitive load

Lesson from substrate x2: reuse that costs effort at the clean-slate moment gets skipped. Therefore: SEED.md maintained continuously; birth detection auto-injects owner laws only; full manifest ceremony stays optional. Anything required at project birth is a bug (I13).

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0ASQGK2ZK03KRP4WPAQG0MJ — Bet: restart-accelerated development plus drift guardrails, one shared substrate  `active`
> lets bet the farm on strong coordination via our journal and intent < - > current reality drift

Owner bet the farm on (A) glanceable intent-reality drift for humans and (B) restart acceleration: new repo births with genesis docs, old code vendored as lessons. Guardrails lower restart NEED; seeds lower restart COST; both attack unrecorded divergence. Restart verbs stay pull-only forever.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0ASQGJZPFJD3FSRT7P4HX44 — Owner-configured cloud environments are full clew nodes  `active`
> we have installable skills, MCP, and I can configure the environments cursor, codex, and claude agents run in the cloud which can absolutely include our golang services

Owner corrected the push-only sandbox assumption: cursor/codex/claude cloud environments are configurable (install scripts, MCP, skills) and can run the Go binary. Cloud write path = provision the environments you own. journal-proposal.yaml is PARKED (trigger-gated for unconfigurable third-party envs only).

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0ASDMN8S3TNV9502TGKHRPD — Tombstone language requires explicit lifecycle metadata; inactivity is not death  `possible-contradiction`
> Tombstone language only when lifecycle metadata explicitly says `tombstoned`; inactivity is never interpreted as death.

Candidate summaries say "died"/"tombstoned" only when the predecessor's lifecycle metadata explicitly reads tombstoned. A quiet or inactive repository is rendered as "changed <date> · active" instead, so silence is never misreported as a dead project.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.90 · tags: internal/lineage/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0ASDMN8S3TNV9502NGMZVQC, e01M0ASDMN8S3TNV9502RT7746E, e01M0ASDMN8S3TNV9502T3QKJ52_

### e01M0ASDMN8S3TNV9502T3QKJ52 — Candidate ranking formula: 0.65 topic overlap + 0.35 recency decay  `possible-contradiction`
> score = 0.65 × binary-cosine topic overlap       + 0.35 × 1 / (1 + ageDays / 30)

clew from with no args ranks predecessor candidates by score = 0.65 x binary-cosine topic overlap + 0.35 x 1/(1 + ageDays/30), with deterministic tie-breaks on score, overlap, recency, repo name, then repo ID. A separate Blatant signal exists only to let the birth card suggest a lineage; it never triggers an import.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.88 · tags: internal/lineage/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0ASDMN8S3TNV9502NGMZVQC, e01M0ASDMN8S3TNV9502RT7746E, e01M0ASDMN8S3TNV9502TGKHRPD_

### e01M0ASDMN8S3TNV9502RT7746E — Lineage import must be human-invoked; carry provenance lives in a separate link  `possible-contradiction`
> `By.Who` must be `human`; automatic lineage imports are rejected.

lineage.Import rejects any request whose By.Who is not "human", so automatic lineage imports are impossible. Original entry/event provenance is kept verbatim rather than rewriting source.kind; the carry fact is recorded in an append-only lineage/<ULID>.yaml link written last, making interrupted runs safely resumable.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.91 · tags: internal/lineage/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0ASDMN8S3TNV9502NGMZVQC, e01M0ASDMN8S3TNV9502T3QKJ52, e01M0ASDMN8S3TNV9502TGKHRPD_

### e01M0ASDMN8S3TNV9502PGAAE99 — Seed regenerates only on journal change, never on metadata polling  `possible-contradiction`
> It compares repository identity and `JournalRevision`, so repeated sync polls cannot rewrite `SEED.md` merely because README topics, `HEAD`, dirty state, or other sampled repository metadata changed.

seed.WriteOnJournalChange is the watcher-facing gate: it compares repository identity and JournalRevision so repeated sync polls cannot rewrite SEED.md when only README topics, HEAD, or dirty state changed. Recurring sync paths must call it instead of seed.Write, and it refuses to overwrite a corrupt existing seed.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.92 · tags: internal/seed/** · pairs-with: e01M0ASDMN8S3TNV9502NGMZVQC_

### e01M0ASDMN8S3TNV9502NGMZVQC — Ambient seed carries lore only; questions and intents stay in the manifest path  `possible-contradiction`
> Active questions and intents are deliberately excluded from ambient lore; they remain available through the separate deliberate manifest path.

SEED.md deliberately carries decisions, findings, graveyard, exhibits, and an optional organ-bank pin. Active questions and intents are excluded from ambient lore and remain available only through the separate deliberate manifest path, which stays pull-only and is never a birth gate.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.93 · tags: internal/seed/**, internal/lineage/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0ASDMN8S3TNV9502PGAAE99, e01M0ASDMN8S3TNV9502RT7746E, e01M0ASDMN8S3TNV9502T3QKJ52, e01M0ASDMN8S3TNV9502TGKHRPD_

### e01M0ARSXN8PR0Y3PVAG8JVYZFY — Ambient budget: refuse over-budget promotion, keep oldest-certified prefix  `possible-contradiction`
> handles a concurrent remote overflow deterministically by retaining the oldest-certified prefix, never letting a newer promotion evict an older ambient law;

Promotion is refused before either the entry or its certification is written when the 1,024-byte ambient budget would be exceeded, so there is never a partial law write. If a concurrent remote promotion causes overflow anyway, the oldest-certified prefix is retained deterministically — a newer promotion can never evict an older ambient law. Overflow state is reported for loud surfacing.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L379 · confidence: 0.90 · tags: internal/owner/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0ARSXN8PR0Y3PVAG5XRP4W8 — Owner journal is its own git repo, never registered as a project  `possible-contradiction`
> Does not register the owner repository as a project, so adapters, archaeology, poller, and project-session discovery cannot scan it.

The owner law store lives in a dedicated normal git repository at $CLEW_HOME/owner, reusing the clew/journal append-only branch and gitx.Sync. It is deliberately not registered as a project so adapters, archaeology, the poller, and project-session discovery cannot scan it. A configured remote is optional; the default empty owner.remote yields a fully functional local-only owner journal.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L379 · confidence: 0.92 · tags: internal/owner/**, internal/config/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0ARA4Y7Z4Q65SP0T7A8EFEY — One canonical renderer: admission and injection measure the same bytes  `possible-contradiction`
> The package will own one canonical renderer so admission and injection measure the exact same bytes; that prevents a law from being admitted under one format and silently omitted under another.

The owner package owns a single renderer used both when admitting a promoted law against the size budget and when injecting laws later. Reason: two formatters would let a law be admitted under one format and then silently omitted under another.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L207 · confidence: 0.93 · tags: internal/owner/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0AR79QQ9PPSZKQ2EVTFNQMV — Ambient seed carries project lineage only; owner laws join at materialize.Conte…  `possible-contradiction`
> The ambient seed should be project lineage data, not ambient owner law data. Otherwise promoted laws would be duplicated into every predecessor seed and then carried as lore despite I13’s explicit separation.

Promoted owner laws are not written into project journals or into SEED.md. They are passed into materialize.Context as a separate capped section (≤1KiB, after the safety preamble, never dropped under truncation pressure). Reason: laws inside every seed would be duplicated and then re-carried as lore, collapsing I13's law/lineage separation; a law change should rematerialize context, not rewrite s…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-23-54-01a01578-d123-7200-af99-a2105dfb139e.jsonl#L200 · confidence: 0.85 · tags: internal/materialize/**, internal/journal/** · evidence: 1 · pairs-with: e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AST7PZRJP8XBWDAGG7AD44_

### e01M0AR79QQ9PPSZKQ2ETN2K1Y0 — Birth boundary is the user-scope SessionStart hook, installed by watch install  `possible-contradiction`
> Claude Code’s user-scope `SessionStart` hook is the correct first-session boundary. Its input includes `cwd`; successful stdout or `additionalContext` is inserted before the first prompt.

Claude's machine-level SessionStart hook is the seam for first-session laws: it receives cwd and its stdout reaches the model before the first prompt. It is installed by clew watch install (not project init) so the machine, not the repo, is what gets watched. A daemon-side global session scanner stays only as fallback for Codex/wrap — a two-second transcript poll is inherently later than startup.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-23-54-01a01578-d123-7200-af99-a2105dfb139e.jsonl#L200 · confidence: 0.85 · tags: cmd/clew/birthhook.go, cmd/clew/watchcmd.go, cmd/clew/initcmd.go · pairs-with: e01M0AR2HXKR05KE2VRTH1ETY28_

### e01M0AR71TW3S282XMYX6HV601V — Owner-law admission rejects overflow; it never evicts an older law  `possible-contradiction`
> owner laws remain ordinary findings in a separate append-only journal, but only a human promotion disposition makes one injectable; promotion that would exceed the full 1KB law block is rejected instead of silently evicting an older law

Owner laws live in a separate append-only journal and become injectable only via a human promotion disposition. The law block is capped at exactly 1,024 bytes; a promotion that would exceed the cap is rejected rather than silently evicting an older law. The extractor's project-agnostic signal is a proposal only, never ambient before certification.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-22-35-01a01577-9dd5-7e63-811f-0b52ba3e3cab.jsonl#L225 · confidence: 0.92 · tags: internal/owner/**, cmd/clew/ownercmd.go, internal/extract/assets/instruction.md · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0AR2HXKR05KE2VRTH1ETY28 — Birth runs on a synchronous Claude SessionStart hook, not polling  `possible-contradiction`
> A key birth-path constraint surfaced: polling session files cannot satisfy the first Claude turn because the project hook would be installed after Claude has already started.

Polling session files cannot serve the first Claude turn, because a project-scope hook would only be installed after Claude has already started. Birth therefore installs a machine-scope Claude SessionStart hook that auto-initializes the repo and emits laws-only context before the first prompt; the daemon's session scan stays as fallback for other agent surfaces.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-22-35-01a01577-9dd5-7e63-811f-0b52ba3e3cab.jsonl#L147 · confidence: 0.92 · tags: cmd/clew/birthhook.go, cmd/clew/birthcmd.go, cmd/clew/watchcmd.go · pairs-with: e01M0AR79QQ9PPSZKQ2ETN2K1Y0_

### e01M0AR2BRV4T3ASCNNB7F5QA10 — Law is a journal scope, not a fifth entry type  `possible-contradiction`
> keep “law” as a journal scope, not a fifth entry type. A promotion copies the original finding unchanged into a dedicated owner journal and records a human promotion disposition; only entries with that disposition are eligible for ambient injection.

Owner-level "laws" stay findings: promotion copies the original finding unchanged into a dedicated owner journal and records a human promotion disposition. Only entries carrying that disposition are eligible for ambient injection. Chosen because it preserves the original evidence and avoids teaching every status, differ, and render path a new epistemic type.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L162 · confidence: 0.94 · tags: internal/owner/** · pairs-with: e01M0AQGPHNFVPKVFT67DKX5923, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M0AQHSRF4DVDYZ989W6K185N — Owner laws live in an owner-scope journal with a ≤1KB injection budget  `possible-contradiction`
> an owner-scope journal synced like any other; findings promoted via clew journal promote <id>; extractor proposes promotion when a finding's content is project-agnostic; ≤1KB injection budget into every project's context, forever.

Laws are stored as an owner-scope journal synced like any other. Findings become laws through an explicit `clew journal promote <id>`, with the extractor proposing promotion when a finding's content is project-agnostic. The resulting law set is capped at a ≤1KB injection into every project's context, permanently.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 · confidence: 0.92 · tags: .clew/** · pairs-with: e01M0AQHSRF4DVDYZ989M4DVHYX, e01M0AQHSRF4DVDYZ989MTZQV7E, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1_

### e01M0AQHSRF4DVDYZ989RNZVC46 — `clew from` is the one explicit lineage command; never automatic  `active`
> Runnable at birth or any time later; un-carrying is a reject (carried entries keep provenance like everything else). Never automatic. The birth card may suggest clew from X on blatant name/topic overlap — suggest only, never act.

Pulling a predecessor's seed (decisions, findings, graveyard, exhibits, organ-bank pin) happens only via an explicit `clew from <repo>`; with no args it lists candidates ranked by recency and topic overlap, showing what each would carry. It can run at birth or later, un-carrying is recorded as a reject so carried entries keep provenance, and the birth card may only suggest `clew from X` on name/t…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 · confidence: 0.93 · evidence: 1_

### e01M0AQHSRF4DVDYZ989MTZQV7E — SEED.md is watcher-maintained continuously, never generated on demand  `possible-contradiction`
> the watcher maintains SEED.md continuously alongside context.md — regenerated on journal change, never on demand. The carry-kit always predates the urge to restart.

The watcher regenerates SEED.md alongside context.md on every journal change rather than building it when a restart is requested. Reason: the carry-kit must already exist before anyone wants to restart, so a seed is never missing or stale at the moment it is needed.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 · confidence: 0.94 · tags: .clew/** · pairs-with: e01M0AQHSRF4DVDYZ989M4DVHYX, e01M0AQHSRF4DVDYZ989W6K185N, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1_

### e01M0AQHSRF4DVDYZ989M4DVHYX — Lineage inheritance is explicit; only promoted laws auto-carry  `possible-contradiction`
> lore inheritance was made explicit because a wrong lineage guess poisons a fresh project worse than no inheritance at all — laws are safe to auto-carry precisely because promotion certified them project-agnostic.

clew will never auto-carry project lore into a new repo. Rationale: a wrong lineage guess poisons a fresh project worse than inheriting nothing at all. Owner laws are the sole exception and may be injected automatically, because the promotion step already certified each law as project-agnostic. This is the governing reason behind invariant I13.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 · confidence: 0.95 · tags: .clew/** · evidence: 1 · pairs-with: e01M0AQHSRF4DVDYZ989MTZQV7E, e01M0AQHSRF4DVDYZ989W6K185N, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1_

### e01M0AQHMB1WWTYWTM79JHACMNP — Ambient seed: SEED.md is maintained continuously, never generated on demand  `possible-contradiction`
> the watcher maintains SEED.md continuously alongside context.md — regenerated on journal change, never on demand. The carry-kit always predates the urge to restart.

The watcher keeps SEED.md up to date beside context.md, regenerated on journal change rather than when someone asks for it, so the carry-kit already exists before anyone wants to restart. This separates ambient seed from the deliberate manifest ceremony, which stays pull-only for big restarts and is never a gate.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-23-54-01a01578-d123-7200-af99-a2105dfb139e.jsonl#L9 · confidence: 0.90 · tags: internal/materialize/**, internal/journal/**, cmd/clew/** · pairs-with: e01M0AQHMB1WWTYWTM79HGM5C34, e01M0AR79QQ9PPSZKQ2EVTFNQMV, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0AQHMB1WWTYWTM79HGM5C34 — Lineage is explicit, laws are ambient: wrong inheritance poisons a fresh project  `possible-contradiction`
> lore inheritance was made explicit because a wrong lineage guess poisons a fresh project worse than no inheritance at all — laws are safe to auto-carry precisely because promotion certified them project-agnostic.

Predecessor lore is never auto-carried. Lineage is pulled by one explicit command (clew from <repo>); the birth card may suggest a match but never acts. Owner laws, by contrast, are injected into every new project automatically. Reason: a wrong lineage guess poisons a fresh project worse than no inheritance, while promotion has already certified laws as project-agnostic.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-23-54-01a01578-d123-7200-af99-a2105dfb139e.jsonl#L9 · confidence: 0.95 · tags: cmd/clew/**, internal/manifest/** · evidence: 1 · pairs-with: e01M0AQHMB1WWTYWTM79JHACMNP, e01M0AST4FW4WFMMX225TQ56V5E, e01M0AST4FW4WFMMX225XEXGBD1, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAF0AMC15, e01M0AST7PZRJP8XBWDAGG7AD44, e01M0ATNY3H0WHX91ASD478G9E7, e01M0AWM6188V135JTTP4S3MQAA_

### e01M0AQGPHNFVPKVFT67DKX5923 — Lore inheritance is explicit; only certified laws auto-carry  `possible-contradiction`
> lore inheritance was made explicit because a wrong lineage guess poisons a fresh project worse than no inheritance at all — laws are safe to auto-carry precisely because promotion certified them project-agnostic.

Project lore never crosses repositories automatically — only the explicit `clew from` command carries it. Owner laws are safe to inject ambiently because human promotion certified them project-agnostic. Reason: a wrong lineage guess poisons a fresh project worse than having no inheritance at all.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-22-35-01a01577-9dd5-7e63-811f-0b52ba3e3cab.jsonl#L9 · confidence: 0.95 · tags: internal/lineage/**, internal/owner/**, cmd/clew/fromcmd.go · evidence: 1 · pairs-with: e01M0AR2BRV4T3ASCNNB7F5QA10, e01M0AR71TW3S282XMYX6HV601V, e01M0ARA4Y7Z4Q65SP0T7A8EFEY, e01M0ARSXN8PR0Y3PVAG5XRP4W8, e01M0ARSXN8PR0Y3PVAG8JVYZFY, e01M0ASDMN8S3TNV9502NGMZVQC, e01M0ASDMN8S3TNV9502RT7746E, e01M0ASDMN8S3TNV9502T3QKJ52, e01M0ASDMN8S3TNV9502TGKHRPD, e01M0AST7PZRJP8XBWDABXXWWQH, e01M0AST7PZRJP8XBWDAHVVX8NR, e01M0AXZ5DEXRWRGD3GT1HYEYSD_

### e01M068ECYE067WF6BH7F26VC3D — Cursor v1 stays CLI-only: desktop 0 vs CLI 44 in 7d  `active`
> window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.

window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.

_source: human cli:note · confidence: 1.00_

### e01M065T92NXY1ER6R73YCQNH84 — Recover docket empty-state and withdrawal wording from Task 3 journal source  `possible-contradiction`
> The generated journal contains the Task 3 source pointer and the fixed card decisions. I’m using that exact source to recover the missing empty-state wording and withdrawal semantics before locking the package API.

Rather than re-deriving the card semantics, the generated journal's Task 3 source pointer and fixed card decisions are used as the authoritative source to recover the missing empty-state wording and withdrawal semantics before the package API is locked.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L75 · confidence: 0.80 · tags: docket/** · pairs-with: e01M05SA72DPRGNTY7GCN1P7CED, e01M05SA72DPRGNTY7GCPEX9W2N, e01M05SA72DPRGNTY7GCQ2RASTP_

### e01M065SK8W1ZT32KZF8YGGTP7W — Alerts self-clean via one scoped reconcile with an explicit withdrawal condition  `active`
> I’m shaping this as one state-level reconcile operation: upsert the active alerts for a repo and auto-drop previously open differ-owned kinds that are absent from that poll, with an explicit withdrawal condition stored on each alert.

Rather than only inserting alerts (which never closed) and keying them on mutable prose, alert handling became a single state-level reconcile: upsert the active alerts for a repo and auto-drop previously open differ-owned kinds absent from that poll, with each alert storing an explicit withdrawal condition.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L45 · confidence: 0.93 · tags: state/**, differ/**_

### e01M064DQTWYDVGGAE3M5QRTGME — Re-evaluate the current tree instead of carrying the prior gate verdict forward  `active`
> The earlier gate’s three blockers are the right pressure points, but the checkout has moved: reservation callers and the neutral-workdir behavior now have new code and tests. I’m re-evaluating the present tree, including untracked test files, instead of carrying that verdict forward.

The prior gate's three blockers were judged the right pressure points, but the checkout has since gained new code and tests for reservation callers and neutral-workdir behavior. The gate will therefore re-evaluate the present tree, including untracked test files, rather than reusing the earlier verdict.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L35 · confidence: 0.90_

### e01M05V5HWA6TFT0A0KDZY8S45K — Confine the cap/ratio admission fix to internal/state; no caller or spec changes  `possible-contradiction`
> I’m implementing this entirely inside `internal/state`

The reservation/settlement work is scoped entirely to internal/state rather than changing call sites or the specification. This keeps the enforcement change contained and closes off the alternative of reshaping the caller-facing API or amending the spec to fix over-admission.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L188 · confidence: 0.82 · tags: internal/state/** · evidence: 5 · pairs-with: e01M0AWM6188V135JTTP4S3MQAA_

### e01M05V0G3Q9F41V62P6T44TC04 — Cursor migration must be monotonic — never rewind an existing cursor  `active`
> I’m correcting the migration to be monotonic

A divergent legacy cursor was rewound during migration, causing one duplicate re-extraction. The fix chosen is to make the migration monotonic so a migrated cursor can only move forward, closing off any migration path that repositions a cursor backward and replays already-extracted bytes.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1437 · confidence: 0.78_

### e01M05TG628KAGZ9HCBQSAG562P — Separate session-extraction budget from one-time archaeology budget  `active`
> I’m also separating session-extraction budget from one-time archaeology, because applying a session-token ratio to cold-start docs makes archaeology mathematically impossible at zero observed sessions.

Extraction budget for live sessions is now tracked separately from the one-time historical archaeology budget. Reason: deriving the archaeology allowance as a ratio of observed session tokens yields zero budget at cold start (zero observed sessions), making backfill of historical docs mathematically impossible. Sits alongside the pre-/post-enrollment byte boundary.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1200 · confidence: 0.90_

### e01M05TBYJXEW5N5FE397XYMEHY — Treat first dogfood run's cursor/push/adapter failures as required failure sign…  `active`
> the first dogfood run exposed real cursor, push, and adapter failures, so I’m treating that run as the required failure signal—not as acceptance

The first dogfood run surfaced real cursor, push, and adapter failures. The assistant is treating that run as the required failure signal rather than as an acceptance run, so those failures do not block recording Task 1 as passed but do drive the current work on the live-enrollment/backfill boundary and failure telemetry.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041 · confidence: 0.82_

### e01M05SVGK1Q2MR34Y3CMHR7DXM — No cursor translation: keep `extract:` for live, add bounded `backfill:` for hi…  `active`
> the live watcher keeps the existing `extract:` cursor, explicit history gets a new bounded `backfill:` cursor, and `history-end` freezes their boundary

Instead of migrating or stacking cursors, the live watcher retains the existing `extract:` cursor while explicit history backfill gets a separate bounded `backfill:` cursor, with `history-end` freezing the boundary between them. Chosen because it preserves v1 pending work across upgrade, keeps live and history from overlapping, and eliminates the risky cursor translation step.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L934 · confidence: 0.93 · evidence: 1_

### e01M05SG3NTP5W2JX7Y6MG00HQJ — Hold the Task 2 commit; make backfill and live watch disjoint by construction  `active`
> I’m holding the Task 2 commit and adding a one-time migration, complete-record offsets, and a fixed historical upper bound so backfill and live watch are disjoint by construction.

Rather than commit Task 2 and patch later, the commit is held while three fixes land: a one-time cursor migration for upgrading users, complete-record offsets so init never baselines mid-record, and a fixed historical upper bound so backfill and live watch cannot overlap — disjoint by construction rather than by runtime check.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806 · confidence: 0.90 · evidence: 2_

### e01M05SA72DPRGNTY7GD0TGHEX7 — Scope freeze: relay, TUI, team mode, adapters need a §11 trigger measurement  `possible-contradiction`
> Frozen (build nothing here without citing a §11 trigger measurement in the journal and stopping for owner review): relay server, TUI/native apps, team mode, semantic code graphs, treemaps, new adapters, orchestration.

Relay server, TUI/native apps, team mode, semantic code graphs, treemaps, new adapters, and orchestration are frozen. Building any of them requires first citing a §11 trigger measurement in the journal and stopping for owner review — measurement, not enthusiasm, unfreezes scope.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.93 · tags: JOURNAL_SPEC.md · evidence: 7 · pairs-with: e01M05SA72DPRGNTY7GCPEX9W2N_

### e01M05SA72DPRGNTY7GCQ2RASTP — Cards show verbatim quotes + clickable provenance, never extractor paraphrase  `possible-contradiction`
> Design consequences: cards show verbatim quotes + clickable provenance, never the extractor's paraphrase or reasoning; high-magnitude cards carry one assumptions line; no other friction, ever.

Decision cards must render verbatim quotes with clickable provenance chips (session line / commit / entry) and must never show the extractor's paraphrase, summary, or reasoning. Reason: system-generated explanations are advocacy and increase acceptance of wrong content, while clickable sources reduce over-reliance. One "accepting this assumes: X" line is allowed on high-magnitude cards only; no o…

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.92 · tags: docket/** · pairs-with: e01M05SA72DPRGNTY7GCN1P7CED, e01M05SA72DPRGNTY7GCPEX9W2N, e01M065T92NXY1ER6R73YCQNH84_

### e01M05SA72DPRGNTY7GCPEX9W2N — I10–I12 added as hard invariants enforced in code and tests  `possible-contradiction`
> Add these three as I10–I12, hard law, enforced in code and tests, not convention:

Three new spec invariants, ranking as hard law rather than convention: I10 docket holds only items answerable by 1–3 discrete verbs (nothing FYI-shaped); I11 every card carries a machine-checkable, printed withdrawal condition and the docket keeps no history/counts/badges; I12 hard cap of seven cards, and sustained volume or an unneeded push is logged as system failure, never user workload.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.93 · tags: JOURNAL_SPEC.md, docket/** · evidence: 7 · pairs-with: e01M05SA72DPRGNTY7GCN1P7CED, e01M05SA72DPRGNTY7GCQ2RASTP, e01M05SA72DPRGNTY7GD0TGHEX7, e01M065T92NXY1ER6R73YCQNH84_

### e01M05SA72DPRGNTY7GCN1P7CED — Rename the inbox surface to "docket"; keep inbox as hidden alias  `possible-contradiction`
> Rename the surface — vocabulary is a forcing function against email-drift. It's a docket of decision cards (clew docket, with inbox as a hidden alias for muscle memory).

The decision surface is renamed inbox → docket, with "inbox" kept only as a hidden alias for muscle memory. Reason: vocabulary is a forcing function against email-drift — an inbox invites FYI accumulation, unread counts, and backlog; a docket is a list of items awaiting a ruling. The docket is the only surface that carries verbs.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.94 · tags: docket/**, inbox/** · pairs-with: e01M05SA72DPRGNTY7GCPEX9W2N, e01M05SA72DPRGNTY7GCQ2RASTP, e01M065T92NXY1ER6R73YCQNH84_

### e01M05RHSWXDNR10P1PY8ERYA9S — Dogfood metrics predeclared; D0 snapshot recorded  `active`
> Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.

Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.

_source: human cli:note · confidence: 1.00_

### e01M04WCGJS9FS7FQB0YFX9DTYG — Name the system clew (owner decision)  `active`
> Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

_source: human cli:note · confidence: 1.00_

## Findings

### e01M0AYE066QK08QK8MTPX4XNFX — Codex finished I13 stale: tree uncommitted, law wording on human surfaces  `current`
> The main-branch implementation remains uncommitted in the working tree; I did not commit or push the code changes.

Manual check-in of the first stale finish: I13 complete and tests green but uncommitted — it exists only in the laptop working tree, invisible to the join. Confirmed conflict: owner-law vocabulary on human surfaces (README, cards, listings) vs the two-register ruling; the feature stands, only surface wording renames. Aligned: single-lineage from, SessionStart birth. Reconcile due at next contact.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95 · evidence: 1_

### e01M0AXCM55N0QM9RCRYF48TQ6C — Universal injection point: every model API call rebuilds the mind  `current`
> It intercepts every chat request, injects relevant memories into the system prompt, and forwards the enriched request to Ollama — all without any changes to your client code.

No context persists inside a model between calls — each harness reconstructs the full message list per request. A local base-URL shim (OLLAMA_HOST / OPENAI_BASE_URL / ANTHROPIC_BASE_URL) can inject the journal delta into any agent, bare ollama included, with passthrough-on-failure so it is never load-bearing. Prior art: Engram transparent ollama proxy; LiteLLM async_pre_call_hook.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.90 · taint: tool_result_

### e01M0AXCM54CDSPV84DBH1PWGWD — Spec nudge matrix is stale: codex and gemini now ship injection hooks  `current`
> Plain text on `stdout` is added as extra developer context.

Aug 2026 survey: codex hooks are stable and default-enabled with UserPromptSubmit additionalContext; gemini CLI BeforeAgent injects context (default on v0.26+); cursor injects only at sessionStart/postToolUse, not beforeSubmitPrompt; opencode plugins transform system/messages pre-dispatch. MCP 2026-07-28 subscriptions notify the client, not the model. Re-pin JOURNAL_SPEC 8.1.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.90 · taint: tool_result_

### e01M0AV5721XB0ZFQW42241YJ4A — Daemon fallback birth discarded the triggering session's first prompt  `current`
> I reproduced this end to end with the compiled daemon. A fresh Codex transcript contained `session_meta` plus `FIRST PROMPT MUST BE JOURNALED`. After five seconds:

Reproduced with the compiled daemon: on fallback discovery, bootstrap baselined every discovered transcript to current EOF, so a fresh Codex transcript's first prompt was never journaled (tail/extract/history-end all set to the 320-byte EOF, zero occurrences in the journal). Safe only for synchronous Claude SessionStart, which runs before the first prompt.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-22-35-01a01577-9dd5-7e63-811f-0b52ba3e3cab.jsonl#L1510 · confidence: 0.88 · tags: cmd/clew/watchcmd.go, cmd/clew/birthcmd.go, cmd/clew/initcmd.go · taint: tool_result_

### e01M0ATNY3H0WHX91ASD0AMR8TE — Cold CLEW_HOME loses concurrent births; warm machine state is safe  `current`
> On an empty `CLEW_HOME`, two simultaneous `_birth` processes in different new repos caused one to exit with `open state.db: database is locked`.

With an empty CLEW_HOME, two simultaneous `_birth` processes in different repos left one dead on `open state.db: database is locked`; with the DB precreated but the owner store uninitialized, five of six concurrent births failed on owner git init/config lock/template copy/worktree add. Once the owner store was fully initialized, six concurrent births passed. The danger window is first-run bootstr…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-30-01a0159d-11bd-7b13-81c5-26b7f910b998.jsonl#L344 · confidence: 0.90 · tags: cmd/clew/**, internal/state/**_

### e01M0ATNY3H0WHX91ASCWE1MM7R — Repo identity is the absolute path, so a rebuilt repo at a reused path is not a…  `current`
> repository identity is only an absolute path, so a fresh repo at a reused path is not a newborn

Registration keys on the checkout path (gitx.RepoID hashes only the absolute path), so moving a checkout aside and running `mkdir x && git init` at the same path is treated as already registered. Reproduced: the second birth exited 1 with a fatal missing-worktree error and emitted zero context. It fails safely today, but path-as-identity is exactly the stale machine state that can rebind an old j…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-30-01a0159d-11bd-7b13-81c5-26b7f910b998.jsonl#L344 · confidence: 0.92 · tags: internal/gitx/**, cmd/clew/**_

### e01M0ATARAN8NWCBGNTM58QRVPX — Promotion candidates enter project context before the human rules on them  `current`
> Extraction adds a promotion candidate as an ordinary live finding

Extraction stores a promotion candidate as an ordinary live finding, and materialization includes all live findings — only the promotion alert is filtered out of the alerts and nudge sections. The candidate's title and body therefore land in context.md immediately, contradicting the documented boundary that a candidate cannot enter context until promote or keep-local resolves the card.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L229 · confidence: 0.86 · tags: internal/extract/**, internal/materialize/**_

### e01M0ATARAN8NWCBGNTM2VNR895 — Title-only imperatives bypass injection withholding into every project  `current`
> Injection withholding scans only `Body` and `Quote`

The imperative-withholding scan checks only Body and Quote (internal/journal/algebra.go:85), but ambient owner laws render Title and Body (internal/owner/owner.go:325). A benign body and quote carrying an imperative title passes the safety gate and becomes ambient in every project's context. Proposed fix: include Title in the scan, add a title-only regression test, recheck raw content at owner ad…

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L229 · confidence: 0.88 · tags: internal/journal/**, internal/owner/**_

### e01M0ASSNH1HP68M1QERV9AKG5A — Attachments bypass the secret scrub; GitHub push protection caught PATs  `current`
> GITHUB PUSH PROTECTION - Push cannot contain secrets

GitHub push protection blocked the journal push: two ephemeral PATs the owner pasted in chat were present verbatim in the attached raw transcript. The entries pipeline scrubs quotes/bodies (6.2) but attachments bypass the scrub entirely. Fix: run the same secret-scrub over transcripts/ (and any attachment) before commit; treat platform push-protection as backstop, never primary.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · tags: internal/scrub/** · taint: tool_result_

### e01M0ASQGKMHGY4F7JBEQX1ZT3T — This chat became the supernode: 788 messages of unjournaled load-bearing context  `current`
> i am afraid to leave this session with you because you have all the context (the exact thing clew is supposed to help with, funnily enough)

The clew-design session ran 3 days on an uncovered surface; the owner became afraid to close it - the exact single-point-of-failure clew abolishes. Exit kit: full raw transcript attached at transcripts/ on this branch; distilled decisions/findings/questions journaled; resumption works from any surface via branch fetch.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00_

### e01M0ASQGKJN308G4HVJTE4A4AB — note-then-edit has a limbo: placeholders auto-commit and leaked into a seed  `current`
> phone surface intent placeholder

journal note commits placeholder text immediately; later file rewrites sit uncommitted, so a manifest exported placeholders into gen-2. Fix direction: clew journal add <file> for validated whole-entry ingestion (also needed by cloud self- extraction). Until then notes must be final text.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · taint: tool_result_

### e01M0ASQGKF09RFDGJD90VNADCQ — Spawn test: 63/63 entries carried into scratch gen-2 with guardrails intact  `current`
> carried 63, dropped 0 (dispositions journaled - the loss is deliberate and dated)

init --carry into a fresh repo: full seed landed, carried provenance preserved, newborn glance renders the constitution, context.md opens with the 6.5 injection preamble before any agent typed. Cross-machine multi-hop proven: laptop decisions -> branch -> cloud VM -> manifest -> gen-2. Differ re-flagged design-era contradictions in the newborn.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · taint: tool_result_

### e01M0ASQGKDZ8J6SQQK3BXKYAV5 — module clew blocks go install by URL; release binaries or rename needed  `current`
> go 1.26.3

go.mod declares 'module clew', not the repo path, so go install github.com/maceip/clew/... fails. Env recipes need git clone && go build until the module is renamed or release binaries ship. Blocks the one-line cloud env install.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · tags: go.mod · taint: tool_result_

### e01M0ASQGKB9VP38VBTQ7YC9GBR — Regime detector: composition and earned-state separate corpses from control  `current`
> my GitHub probably has a hundred or more examples

Cadence cannot distinguish clew from the corpses (all are burst-projects at day scale). What separates: core-touch ratio (clew 50% vs 0-22%), earned state (passing gates, live dogfood, metrics vs zero), and inheritance (clew is the first generation that carried anything). n=4; the ~100-repo lineage census remains to be run.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0ASQGK9TDRH9RHAJG9YDPME — Census: security-substrate human-sealed in 1 day; stratura inherited nothing  `current`
> preserved as future MOSS input

security-substrate: born the day after substrate's tombstone, formally SEALED as failed on day 2 with STOP packet and constitution - faster than any detector. stratura: zero references to either predecessor (measured), then repeated the sealed pathology (safety perimeter before usable core). Detection was never the bottleneck; INHERITANCE was.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · taint: tool_result_

### e01M0ASQGK6W0Y1NEDQW45KDECS — Corpse census: substrate died in a 6-day burst; tombstone came 5 weeks late  `current`
> document adoption failure and harden project creation

substrate: 63/64 commits in week one (Jun 9-14), five weeks silence, final commit is the failure confession (LIFECYCLE.md + README 'failed adoption'). The promised compounding loop (scheduler/repair/steward/federated store) was never built - confessed by its own docs. Zero tags, zero CI, zero adopters.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · taint: tool_result_

### e01M0ASDMN8S3TNV9502WFBA70W — Losing the lineage/ directory in sync destroys cycle protection  `current`
> Entries and events alone are insufficient; losing the links loses durable lineage declarations and transitive-cycle protection.

lineage.AncestorIDs reads the append-only lineage/ links to compute transitive ancestry, which is what makes A→B→C then C→A rejectable without contacting either predecessor. Entries and events alone are insufficient: if journal-branch sync or remote adoption drops lineage/, durable declarations and cycle protection are lost.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.87 · tags: internal/lineage/**_

### e01M0ARPN4E0VHAP2DEJSGP7GZB — macOS /var vs /private/var alias made an initialized owner repo look foreign  `current`
> I also tightened repository identity so macOS’s `/var` versus `/private/var` alias cannot make an initialized owner repo look foreign.

On macOS, /var is a symlink to /private/var, so path-based repository identity could treat an already-initialized owner repository as a different, foreign repo. Repository identity was tightened to resolve the alias. Surfaced while the core owner package tests passed, including a two-clone git round trip and an over-budget refusal with no partial law write.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L315 · confidence: 0.90 · tags: internal/owner/**_

### e01M069MQYJX6QVW3YCWWTAWV34 — --help  `current`
> --help

--help

_source: human cli:note · confidence: 1.00_

### e01M068FQH1ND9MM1WH851AF45M — Task 6 gate: flags 0 writes; algebra, poller, manifest pass  `current`
> note-help entries=69→69; absence-threshold=4 proposed/5 absent; ineligible=proposed; human-confirm=absent; contradiction nonhuman=possible/human=contradicted; env different=current/current, same=superseded/current; poller best-overlap=pass/no-overlap=none/out-of-window=none; manifest rerun events=2→2; full=pass; race=pass; vet=pass.

note-help entries=69→69; absence-threshold=4 proposed/5 absent; ineligible=proposed; human-confirm=absent; contradiction nonhuman=possible/human=contradicted; env different=current/current, same=superseded/current; poller best-overlap=pass/no-overlap=none/out-of-window=none; manifest rerun events=2→2; full=pass; race=pass; vet=pass.

_source: human cli:note · confidence: 1.00_

### e01M06821D53QYHBJS1FEC2CK7G — Task 5 gate: 3 formats, 1 card, confirm boundary pass  `current`
> formats=bundle+dir+https; schema-invalid=reject; quote-missing=reject; batch=1 card; live-stage=1 entry; live-open=pass; live-reject=0 journal writes; accept=1 foreign entry+1 human confirm; branch-push=pass x2 idempotent; full=pass; race=pass; vet=pass.

formats=bundle+dir+https; schema-invalid=reject; quote-missing=reject; batch=1 card; live-stage=1 entry; live-open=pass; live-reject=0 journal writes; accept=1 foreign entry+1 human confirm; branch-push=pass x2 idempotent; full=pass; race=pass; vet=pass.

_source: human cli:note · confidence: 1.00_

### e01M067G2QY8BZNZXBE46QBEXEW — Task 4 gate: 10ms, HTML 30s, ntfy 5/5  `current`
> bare-clew=10ms x5; dashboard-sections=5; html-refresh=30s; title-light=nonempty-only; ntfy-delivered=5; payload-valid=5/5; full=pass; race=pass; vet=pass.

bare-clew=10ms x5; dashboard-sections=5; html-refresh=30s; title-light=nonempty-only; ntfy-delivered=5; payload-valid=5/5; full=pass; race=pass; vet=pass.

_source: human cli:note · confidence: 1.00_

### e01M067AZ3VJM73XP82REZEZ1QC — Task 4 channel: ntfy 2383c10c…, card creation only  `current`
> ntfy-topic=https://ntfy.sh/2383c10ce6438813da9969532f2df2f7; push-trigger=docket-card-creation-only; payload=headline+why-you; html-refresh=30s.

ntfy-topic=https://ntfy.sh/2383c10ce6438813da9969532f2df2f7; push-trigger=docket-card-creation-only; payload=headline+why-you; html-refresh=30s.

_source: human cli:note · confidence: 1.00_

### e01M066XWD8F1SAWQ8HVXGW4J4Z — Task 3 docket gate: 8→1, FYI 0, withdrawal 1 poll  `current`
> cards=8; render=1 overflow-failure; cap=7; synthetic-FYI-rendered=0; resolved-stomp-withdrawal=1 poll; pushes=0/0; full=pass; race=pass; vet=pass.

cards=8; render=1 overflow-failure; cap=7; synthetic-FYI-rendered=0; resolved-stomp-withdrawal=1 poll; pushes=0/0; full=pass; race=pass; vet=pass.

_source: human cli:note · confidence: 1.00_

### e01M066CV9W2WX7FZA44QYME204 — Alert self-cleaning shipped: WithdrawWhen, ReconcileAlerts, six kinds withdrawn  `current`
> - Differ now withdraws stale contradiction, absence, aging, suspect, stomp, and overlap alerts each poll.

Landed: Alert.WithdrawWhen with a legacy DB migration, ReconcileAlerts(repo, kinds, active), and differ withdrawal of stale contradiction, absence, aging, suspect, stomp, and overlap alerts on every poll. Resolved-stomp and status-resolution tests added; go test ./..., race tests, and vet all pass.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L263 · confidence: 0.88 · tags: state/**, differ/**_

### e01M0660RW03A1VSC6ENMH9M2J7 — Parallel agent task killed by gpt-5.6-sol TPM rate limit  `current`
> stream disconnected before completion: Rate limit reached for gpt-5.6-sol in organization ‹redacted› on tokens per min (TPM): Limit 500000, Used 397298, Requested 204717.

The task2_final agent errored with "stream disconnected before completion" due to an OpenAI org-level tokens-per-minute limit for gpt-5.6-sol: limit 500000, used 397298, requested 204717. Running several agents concurrently can exceed the TPM ceiling and drop in-flight work.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L160 · confidence: 0.85 · taint: tool_result_

### e01M065G8RTKVH6466KE7GQREGX — Task 2 passes its live gate on all five acceptance checks  `current`
> Task 2 now passes its live gate: 52 automatic session entries, 0 delivered/unneeded pushes, monotonic cursors, exact installed binary, and no active adapter/LLM errors.

A live gate run on Task 2 passed: 52 automatic session entries recorded, 0 delivered-but-unneeded pushes, cursors monotonic, the exact installed binary in use, and no active adapter or LLM errors. This is the verdict that unblocked committing the gate fixes.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988 · confidence: 0.93_

### e01M064YRS4S9NK7KW9NQ1JMDV2 — Malformed or missing pinned timestamps silently fall back to ingest time  `suspect`
> Source time: malformed/missing pinned timestamps silently become ingest `now`.

When a source record's pinned timestamp is missing or malformed, the adapter/extract path substitutes the ingest-time `now` without signalling, so entries get fabricated source times. Located at adapters.go:151 and extract.go:264; flagged as a gate blocker.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252 · confidence: 0.92 · tags: internal/adapters/**, internal/extract/** · evidence: 3_

### e01M064YRS4S9NK7KW9NN351FQF — I9: Claude settlement ignores cache token fields, letting spend exceed caps  `suspect`
> I9: Claude settlement ignores `cache_creation_input_tokens` and `cache_read_input_tokens`, permitting cumulative spend beyond caps.

Settlement of Claude LLM calls counts only the non-cache token fields, ignoring `cache_creation_input_tokens` and `cache_read_input_tokens`. Cumulative spend is therefore undercounted and can run past the configured budget caps. Found at llm.go:158 during the read-only gate.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252 · confidence: 0.93 · tags: internal/llm/**, internal/state/** · evidence: 1_

### e01M064K88F7SDMHA1SPAB51HK7 — D2 final: 52 automatic entries; live extraction 0.631%  `current`
> D2-final: repos=3; automatic-session-entries=52; observed=6779248; live-extraction=42803 (0.631%); backfill=5057; all-LLM=67936/200000; C:R=0:1; pushes=0 delivered/0 unneeded (precision=N/A); adapter/system incidents=4; parked=0; active-reservations=0; live-sessions=6.

D2-final: repos=3; automatic-session-entries=52; observed=6779248; live-extraction=42803 (0.631%); backfill=5057; all-LLM=67936/200000; C:R=0:1; pushes=0 delivered/0 unneeded (precision=N/A); adapter/system incidents=4; parked=0; active-reservations=0; live-sessions=6.

_source: human cli:note · confidence: 1.00_

### e01M0640H4ZAR0BQMS73R8QW9E7 — Repaired watcher installed as launchd agent dev.clew.watch  `current`
> installed launchd agent dev.clew.watch (log: /Users/mac/.clew/logs/watch.log)

After journaling the D2 cursor-rewind finding, the fixed watcher was installed as a launchd agent named dev.clew.watch on the dev Mac, writing to /Users/mac/.clew/logs/watch.log. That log path is where watcher behaviour for subsequent live runs can be inspected.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1508 · confidence: 0.86 · tags: cmd/clew/**_

### e01M05V9MQWYX3BAX0VXZ70SHTD — D2: cursor rewind replayed 58,754 bytes once  `current`
> D2 migration failure: split cursor rewind replayed 58754 bytes and spent 1815 extraction tokens once; delivered pushes=0. Fix is monotonic max(extract, watch-extract), with divergent-cursor regression. State backup: state.db.d1-20260816T1748Z.bak.

D2 migration failure: split cursor rewind replayed 58754 bytes and spent 1815 extraction tokens once; delivered pushes=0. Fix is monotonic max(extract, watch-extract), with divergent-cursor regression. State backup: state.db.d1-20260816T1748Z.bak.

_source: human cli:note · confidence: 1.00_

### e01M05TX51R4VZHJMNM804YPMRT — D1: 30 automatic entries; Codex metadata incident pinned  `current`
> D1 live dogfood: 30 session entries appeared from 1 real Codex session with 0 manual notes; observed=5549571, live+backfill extraction=30184, all-LLM=39091, pushes delivered=0/0, open alerts=10; 46 records in 3 newly observed multi-agent metadata classes were pinned as non-utterance adapter metadata.

D1 live dogfood: 30 session entries appeared from 1 real Codex session with 0 manual notes; observed=5549571, live+backfill extraction=30184, all-LLM=39091, pushes delivered=0/0, open alerts=10; 46 records in 3 newly observed multi-agent metadata classes were pinned as non-utterance adapter metadata.

_source: human cli:note · confidence: 1.00_

### e01M05T4XG0RJWQTP25T1K58B61 — Confirm/reject only in event YAML; adapter unknowns undated, absent from status  `suspect`
> Confirm/reject exists only in journal event YAML.

Human confirm/reject signals are recorded only in per-worktree events/*.yaml, so measuring confirm rate requires a find+awk scrape instead of a DB query. Adapter "unknown" counts are cumulative, undated KV rows and never surfaced in status. Both dogfood metrics are therefore not queryable from state.db.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441 · confidence: 0.87 · tags: internal/**, cmd/clew/** · evidence: 15 · taint: tool_result_

### e01M05T4XG0RJWQTP25SYT4FH0B — Task 2 not passable: `spent` conflates extraction, differ, and archaeology  `suspect`
> `spent` combines extraction, differ, and archaeology; it is not extraction-only.

The dogfood audit judged Task 2 not passable yet. The budget `spent` counter mixes extraction, differ, and archaeology tokens, so the predeclared extraction-only cost metric cannot be read from it. Separating cost by kind is a prerequisite before the Task 2 gate can be honestly evaluated.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441 · confidence: 0.90 · tags: internal/**, cmd/clew/** · evidence: 15 · taint: tool_result_

### e01M05SPB2EMMC4F4PR0BDQA8S5 — First watch treated historical sessions as live: 342 overlaps, 27 stomps, 12.9M…  `suspect`
> First watch misclassified historical sessions as live, producing 342 overlaps, 27 stomps, and 12,895,847 observed tokens.

Measured fallout of the first watch run misclassifying pre-existing historical sessions as live: 342 overlaps, 27 stomps, and 12,895,847 observed tokens. This quantifies the historical-session storm previously recorded qualitatively as an I12 failure.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335 · confidence: 0.91 · tags: cmd/clew/** · evidence: 4 · taint: tool_result_

### e01M05SG3NTP5W2JX7Y6KZD5N6P — Pre-commit review found three Task 2 blockers: cursor migration, backfill overl…  `current`
> Pre-commit review found three real blockers: upgrade users lacked cursor migration, backfill could overlap live suffixes, and init could baseline inside a partial JSONL record.

A pre-commit review of the Task 2 work surfaced three real blockers: users upgrading had no cursor migration path; backfill could overlap with live suffixes; and init could baseline in the middle of a partial JSONL record, producing a corrupt offset.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806 · confidence: 0.93_

### e01M05SA72DPRGNTY7GCV0CBDK7 — Assumptions prompt cut over-reliance 42%→22%; stacked/delay friction backfired  `current`
> the assumptions prompt ("accepting this assumes X") reduced over-reliance from 42%→22% without added cognitive load, while stacked/delay-based friction backfired

From the over-reliance literature (Buçinca cognitive-forcing lineage through Ghosh et al. 2026, n=214): an "accepting this assumes X" prompt reduced over-reliance from 42% to 22% with no added cognitive load, while delay-based and stacked friction backfired. This is why the assumptions line is the single permitted forcing function on clew cards and why no other friction is added.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.85 · tags: docket/**_

### e01M05S9SFKAAM813AR1B8DXEYW — Codex format now detected; watcher tracks only post-baseline bytes  `current`
> The current Codex format is now detected: re-init found 3 large-metadata sessions, and the watcher is tracking only post-baseline bytes with source time.

After repair, clew detects the current Codex session format: re-initialization found 3 sessions with large metadata, and the watcher now tracks only bytes written after the baseline, using source time rather than ingest time. This addresses the historical-session storm seen in dogfood D0.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729 · confidence: 0.85_

### e01M05S9SFKAAM813AR1AY96QWH — Final Task 2 dogfood snapshot: 0.113% extraction, 0:1 confirm:reject, 0 pushes  `current`
> Final Task 2 snapshot is: 6 automatic entries; extraction 5,057 / 4,491,713 observed = 0.113%; confirm:reject 0:1; actual pushes 0/0; adapter incidents 1 (the journaled D0 storm); parked 0.

Final Task 2 dogfood measurement of clew: 6 automatic entries produced; 5,057 of 4,491,713 observed bytes extracted (0.113%); confirm-to-reject ratio 0:1; actual pushes 0 of 0; 1 adapter incident (the previously journaled D0 session storm); 0 parked items. Recorded after the Codex-format detection and watcher baseline repairs.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729 · confidence: 0.93_

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

### e01M0BER1Q5W312NDD7GZ283RCM — Cap the two screens at 7+7; overflow policy is an open question  `open`
> 7 merge, 7 diff decisions => 14 decisions is probably the max we should allow, although not sure what to do if there are more than that

Owner bound: seven merge plus seven gap — fourteen decisions — is probably the most a human should ever face at once. Unresolved: what happens when more exist. Candidates on record: rank by stakes and fold the rest into one waiting line; treat sustained overflow as restart pressure feeding the held restart tab. Owner has not ruled.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0AY97K0AZHNJW0D87KTD7YH — Two-grade telling: block dependent tasks on unresolved drift?  `open`
> go to agent Y that has drifted and ask it to work on task B that has dependency on the "unmerged decison A", and i forget to "tell it to undo its own drift"

Owner stress-test found the hole: human forgets decision A, asks drifted agent Y for task B that depends on the contested ground; know-and-tell alone lets B build on it. Proposed rule: telling has two grades — courtesy mention when unrelated; blocking question when the new task depends on unresolved drift. Y proceeds only after a ruling. The docket card persists regardless; forgetting stays safe.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0ASQGM2P07Z36KVZX6P8EH4 — Adopt the complexity law: additions must be a verb, label, rendering, or config  `open`
> i dont want to increase the complexity from here

Proposed governance line to stop subsystem sprawl: anything wanting its own process or protocol is frozen by default. Needs explicit ratification.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0ASQGKZ9ZY4XRF6K7KRCT75 — Adopt clew witness <transcript> as the cloud-session gap fix?  `open`
> lets think hard about a simple fix for it

One verb: ingest an exported transcript as a session (adapter-parse, extract in budget, cloud provenance). Ritual line in each env/skill: at checkpoints and end, export + witness + push. No daemon/relay/protocol; reuses the whole pipeline; makes the volume-based 11 trigger moot for owned surfaces. Also rewrite that trigger density-based.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0ASQGKX2F9THW496N79Z83F — Approve selfwatch + journal add + owner-laws relocation to a git-reachable repo  `open`
> i dont want the laptop to be the supernode here

Cloud nodes need: selfwatch (poller+sync+context, no tailing), clew journal add <file> (validated ingestion, fixes note-limbo), self_extracted provenance label, and owner laws moved from ~/.clew scope to a git-reachable repo so cloud births inherit laws.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0ASQGKVW3RPZ2M51J3WQ9G9 — Approve cloud env recipes: install clew + wire MCP in cursor/codex/claude envs  `open`
> I can configure the environments cursor, codex, and claude agents run in the cloud

Row 4 revised: (a) fix module path or ship release binaries; (b) document three env recipes (cursor environment.json, codex cloud setup, claude cloud) installing clew + clew mcp; (c) package the node ritual as an installable skill; (d) one-time test that codex cloud creds can push clew/journal.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0ASQGKRQBVSGHP3RCYNWYEG — Approve extractor over-firing tune: rollup >32KB and docket hit 8 cards  `open`
> !! rollup exceeds 32KB - the extractor is over-firing (tune, don't scroll)

Both 3.3 tripwires fired. Proposed: raise confidence bar / tighten firing triggers; never raise caps. Needs ruling.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · taint: tool_result_

### e01M0ASQGKPPZ4Y3H8P1RV88DZN — URGENT: approve budget-deadlock fix - floor rule or smaller extraction slices  `open`
> extraction paused: LLM budget reservation denied

Extraction is fully starved: single request (29-34k tokens) exceeds the 2% daily budget (22.9k at 1.14M observed), so today's sessions are not being journaled. Options: budget floor = max(2% observed, 1.5x max request); shrink slices; both. Needs ruling before knowledge piles up unprocessed.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 1.00 · taint: tool_result_

### e01M04XVQEEZ38J9TC5NZNKC16B — Run the live fidelity gate (RealProvider) on a machine with provider keys  `answered`
> If you cannot reach P>=0.9 / R>=0.75, stop all other work and report that the kill criterion fired — that outcome is a valid and useful result, and it is written down on purpose

Hermetic pipeline passes, but the go/no-go gate — P≥0.9/R≥0.75 within 5 instruction iterations vs ratified labels — has never run against a live provider. Until it does, extraction quality (the stated existential risk) is unmeasured; the kill criterion is theoretical. Needs claude/codex CLI or OpenAI-compatible key; env flag was RESTART_FIDELITY=1 pre-rename — verify.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 1.00 · tags: acceptance/**_

## Intents

### e01M0AZN6HJETV241AK5RSBDHNR — Held: a restart tab — stage selected drift into the next generation  `proposed`
> users select items from each and put them in "restart tab", and that tab also shows the same output "clew from" would show

Owner direction, held for more thinking: a third tab beside merge and gap. The human selects items from both and stages them into a restart; the tab previews exactly what clew from would emit — live seed curation from the drift you are already reading. Gives merge/gap overflow a relief valve: too heavy to absorb here becomes carry it forward. Not buildable spec yet; owner explains later.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0AXCM561K3C5QXAVGVGT46T — Build the freshness ladder: one delta payload, five delivery layers  `proposed`
> the human (me) would think that our mind plane would inject "knowledge" into all running agents on a specific project at some interval

Deliver one idempotent journal-delta digest via: (1) base-URL proxy shim — the floor, covers ollama; (2) MCP tool-result banner on every clew tool reply; (3) harness hooks where present (claude/codex/gemini prompt boundary, cursor postToolUse, opencode transform); (4) PTY wrap; (5) human relay. Dedupe by entry-ULID watermark so no agent sees a delta twice. Payload is data, not instructions.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 · confidence: 0.95_

### e01M0AST7PZRJP8XBWDAHWK6QNV — Build invariant I13: ambient seed, birth detection, owner laws, clew from  `in_flight`
> invariant I13 — birth costs nothing; laws are ambient, lineage is declared. Build:

Commitment to implement I13 'birth costs nothing': continuous SEED.md, auto-init at repo birth with laws-only injection, an owner-scope law journal with promote, and clew from <repo> for declared lineage (no-arg form lists candidates ranked by recency and topic overlap, each line showing what would be carried). The manifest ceremony stays pull-only and never a gate.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T12-03-33-01a0159d-1ea8-7f70-bfac-3fecc3a17c09.jsonl#L9 · confidence: 0.92 · tags: cmd/clew/**, internal/owner/**, internal/extract/** · evidence: 1_

### e01M0ASDMN8S3TNV9502VTFK27E — Wire the seed/lineage libraries into watcher, materialization, and clew from  `absent`
> I did not edit command wiring, watcher/init behavior, materialization, documentation, manifest, extractor, or parent-owned repository metadata code.

The I13 data layer is done but unwired. Parent integration must: build seeds with lineage.AncestorIDs before writing, swap recurring seed.Write calls for WriteOnJournalChange, preserve the journal's lineage/ directory across journal-branch sync and unrelated-root remote adoption, and implement clew from as read → import → materialize.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-04-01a01578-f76e-74f0-8e19-f7387068f66a.jsonl#L573 · confidence: 0.86 · tags: internal/seed/**, internal/lineage/**_

### e01M0ARJKGSQFSH8WSVZEG520DH — Surface coverage: phone reads the glance and receives decision cards  `absent`
> phone is typically either a variant of cloud or remote to laptop

Owner requirement: phone is read + interrupt only — journal.md via GitHub mobile as the away-glance; ntfy delivering decision cards (headline + why-you) as the only push. No write path in v1 (known gap, trigger on file). Evidence classes: ntfy deliveries with measured push precision; first-screenful glance fix landed. Blocked on owner's 3-minute pairing: rotate leaked topic, redact, subscribe.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-18 · confidence: 0.95_

### e01M0ARJHKDP6Z6R1FKZSJ8AN4S — Surface coverage: laptop agents fully sensed with zero human effort  `absent`
> agents on laptop

Owner requirement: local agents (claude/codex/cursor CLI) on watched machines are tailed, extracted, and journaled with no manual notes. Evidence classes: automatic session entries; extraction spend within I9; adapter incidents loud. Largely met by current dogfood; kept as an intent so regressions surface as absence rather than silence.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-18 · confidence: 0.95_

### e01M0ARJFTWEWY5H6JFJ17656W4 — Surface coverage: PR-only cloud agents (Codex-app-class) contribute knowledge  `absent`
> when I open the android codex/chatgpt app to do more work in the cloud what happens

Owner requirement: agents in sandboxes that can only open PRs (Android Codex/ChatGPT app work) must still contribute session knowledge, not just code. Mechanism blocked on the write-path ruling (see the open question entry). Evidence classes: journal entries whose provenance is a PR-only sandbox; until then, their sessions remain a visible gap in status, never a silent one.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-18 · confidence: 0.95_

### e01M0ARJE1XNN8Q45DJ36FP47YT — Surface coverage: repo-write cloud agents (Cursor-class) are full journal nodes  `absent`
> my originally stated (yet somehow lost, even using the journal) surfaces: agents on web (typically cloud), agents in cloud, agents on laptop, phone is typically either a variant of cloud or remote to laptop

Owner requirement: cloud/web agents working on watched repos read the journal at session start (digest from branch) and write their session knowledge back. Evidence classes: journal-branch pushes originating from cloud runs; digest fetches in cloud AGENTS.md startup. This entry itself was pushed by a credentialed cloud agent as the first proof.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-18 · confidence: 0.95_

### e01M0AQHSRF4DVDYZ989PVHGA7R — Birth detection: auto-init a new repo with owner laws only  `absent`
> new git init + agent session on a watched machine → auto-init (watch, journal branch, context.md) with owner laws only. No lore, no card required. Acceptance: mkdir x && git init && claude yields a context containing the owner's laws with zero clew commands typed.

Build auto-init so that a fresh git init plus an agent session on a watched machine sets up watch, journal branch, and context.md carrying only the owner's laws — no lore and no birth card required. Acceptance test: `mkdir x && git init && claude` yields a context containing the owner's laws with zero clew commands typed.

_source: session codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 · confidence: 0.93 · tags: .clew/**_

### e01M0664DX5WAZ01J6KJBCZP4QC — Add contract tests for docket invariants and failure modes  `in_flight`
> The core package now compiles. I’m adding contract tests around the three hard invariants plus the less visible failure modes: no paraphrase leakage, stale withdrawal, timer ordering, high-magnitude assumptions, exact provenance, event-bound defer, and empty/overflow rendering.

With the core package compiling, the plan is to write contract tests covering three hard invariants plus less visible failure modes: no paraphrase leakage, stale withdrawal, timer ordering, high-magnitude assumptions, exact provenance, event-bound defer, and empty/overflow rendering.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L218 · confidence: 0.88 · tags: docket/** · evidence: 1_

### e01M065G8RTKVH6466KEB0FRJ8N — Commit the gate fixes as one spec-amended change, then start the docket  `in_flight`
> I’m committing the gate fixes as one spec-amended change, then moving to the docket.

After the Task 2 live gate passed, the plan is to land all gate fixes as a single commit that also amends the spec, and then move on to the docket surface work.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988 · confidence: 0.90 · evidence: 1_

### e01M064DHNTXFM29PNWCW0VBBF2 — Run a strict read-only gate proving each acceptance point with tests or state q…  `in_flight`
> I’m doing a strict read-only gate: first loading Clew’s generated context and relevant prior memory, then I’ll inspect the current diff and prove each acceptance point with tests or direct state queries. I’ll report only a blocker or a PASS with exact evidence.

Plan for this session: load Clew's generated context and prior memory, inspect the current diff, and prove each acceptance point via tests or direct state queries — reporting only a blocker or a PASS with exact evidence, making no writes.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L16 · confidence: 0.88 · evidence: 2_

### e01M0642VRXV9PCGA4NDJF92E2Y — Wire atomic budget reservations into every LLM call  `in_flight`
> I’m wiring the new atomic budget reservations into every LLM call next; this closes the remaining race where live extraction and backfill could both spend against the same allowance.

Next step: route every LLM call through the new atomic budget reservations, so live extraction and backfill can no longer both spend against the same allowance. Closing this race is the stated purpose of the change.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1577 · confidence: 0.88 · tags: internal/state/** · evidence: 7_

### e01M05VFAW9A783PMZZEER0G6FX — Second pass on rollover, double-settlement, migration; then run wider suite  `in_flight`
> I’m doing a second pass for rollover, double-settlement, and migration behavior before running the wider suite.

Before treating the internal/state reservation/settlement work as done, do a second review pass covering rollover, double-settlement, and migration behavior, then run the wider test suite beyond the state package.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L268 · confidence: 0.85 · tags: internal/state/** · evidence: 6_

### e01M05V5HWA6TFT0A0KDY1QV6C0 — Add transactional reservation + settlement accounting in internal/state with co…  `in_flight`
> a transactional reservation record plus settlement accounting, with contention tests that prove the cap/ratio cannot be over-admitted

Commitment to implement a transactional reservation record plus settlement accounting inside internal/state, accompanied by contention tests that demonstrate the cap/ratio cannot be over-admitted under concurrent access.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L188 · confidence: 0.88 · tags: internal/state/** · evidence: 11_

### e01M05V0G3Q9F41V62P6R5QV53G — Fix migration monotonicity, restore from D1 boundary, rerun cycle before passin…  `in_flight`
> restoring from the D1 boundary, and will rerun the cycle before calling Task 2 passed.

After the first upgraded live cycle exposed a cursor-rewind defect, the watcher was stopped. Committed follow-up work: correct the cursor migration, restore state from the D1 boundary, and rerun the live cycle. Task 2 will not be declared passed until that rerun is clean.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1437 · confidence: 0.92 · evidence: 2_

### e01M05TPA9ZK4C8RGC3KBPV2MRF — Re-run the live gate: install binary, normalize spend category, restart watcher  `in_flight`
> I’m moving back to the live gate now: install this exact binary, normalize the one dogfood spend category, restart the watcher, then verify one full tail/poll cycle has no historical replay, false sessions, or false pushes.

Next step after the unit path went green is a live-gate run: install the exact built binary, normalize the single dogfood spend category, restart the watcher, then observe one complete tail/poll cycle. Acceptance is negative evidence — no historical replay, no false sessions, no false pushes in that cycle.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1323 · confidence: 0.88 · evidence: 2_

### e01M05TBYJXEW5N5FE398NJ4RD7 — Tighten live-enrollment/backfill boundary and add failure telemetry  `in_flight`
> I’m tightening the live-enrollment/backfill boundary and failure telemetry now

Work in progress to harden the boundary between live enrollment and backfill, and to add telemetry for failures. This follows the recorded pass of Task 1 and the cursor/push/adapter failures observed in the first dogfood run.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041 · confidence: 0.85 · evidence: 2_

### e01M05SPB2EMMC4F4PR0BSPFZKZ — Narrow fix for the watch storm: transactional baselines, source-time, bounded c…  `in_flight`
> Root cause and narrow fix sent to parent: transactional live baselines, source-time sessions, and bounded separate backfill cursor.

The semantics investigation reported a root cause and a narrow fix to the parent: make live baselines transactional, use source-time (not observation-time) session timestamps, and give backfill its own bounded cursor separate from live watch. This is the proposed work to make backfill and live watch disjoint.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335 · confidence: 0.62 · tags: cmd/clew/**, internal/** · evidence: 20 · taint: tool_result_

### e01M05SG3NTP5W2JX7Y6MP1F1K6 — Add cursor migration, complete-record offsets, and fixed historical upper bound  `in_flight`
> adding a one-time migration, complete-record offsets, and a fixed historical upper bound

Committed follow-up work before the Task 2 commit can land: a one-time cursor migration for existing installs, offsets that always fall on complete JSONL record boundaries, and a fixed upper bound on the historical backfill range.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806 · confidence: 0.88 · evidence: 4_

### e01M05SA72DPRGNTY7GCXARNKYH — Implement the decision card and enforce I10–I12 in renderer and tests  `in_flight`
> Enforce I10–I12 in the renderer and in tests: a synthetic FYI item must be unrenderable; an 8th card must collapse to the overflow-failure card; a resolved stomp must withdraw within one poll cycle.

Commitment (Task 3): build the docket card to the fixed anatomy — headline-as-question ≤80 chars, why-you strip with rule fired and ticking stall timers, verbatim-quote evidence rows with provenance chips, assumptions line on high-magnitude only, 1–3 verbs plus defer-until-event, printed withdrawal condition, ordering by blocking cost, designed empty state. Tests must prove the three invariants.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.90 · tags: docket/** · evidence: 1_

### e01M05S9SFKAAM813AR1EG3X3WR — Land the dogfood fixes after recording the Task 2 snapshot  `in_flight`
> I’m recording that snapshot and then landing the dogfood fixes.

Commitment to record the final Task 2 dogfood snapshot and then land the outstanding dogfood fixes in the codebase.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729 · confidence: 0.82 · evidence: 2_
