# Journal

_generated 2026-08-17 05:14 UTC · 53 live entries (14 decisions · 26 findings · 1 questions · 12 intents) · 71 total in history_

## DECIDED

- `e01M068ECYE067WF6BH7F26VC3D` Cursor v1 stays CLI-only: desktop 0 vs CLI 44 in 7d — 7h · active
- `e01M064DQTWYDVGGAE3M5QRTGME` Re-evaluate the current tree instead of carrying the prior gate verdict forward — 8h · active
- `e01M05V5HWA6TFT0A0KDZY8S45K` Confine the cap/ratio admission fix to internal/state; no caller or spec changes — 11h · active
- `e01M05V0G3Q9F41V62P6T44TC04` Cursor migration must be monotonic — never rewind an existing cursor — 11h · active
- `e01M05TG628KAGZ9HCBQSAG562P` Separate session-extraction budget from one-time archaeology budget — 11h · active
- `e01M05TBYJXEW5N5FE397XYMEHY` Treat first dogfood run's cursor/push/adapter failures as required failure sign… — 11h · active
- `e01M05SVGK1Q2MR34Y3CMHR7DXM` No cursor translation: keep `extract:` for live, add bounded `backfill:` for hi… — 11h · active

## LEARNED

- `e01M068FQH1ND9MM1WH851AF45M` Task 6 gate: flags 0 writes; algebra, poller, manifest pass — 7h · current
- `e01M06821D53QYHBJS1FEC2CK7G` Task 5 gate: 3 formats, 1 card, confirm boundary pass — 7h · current
- `e01M067G2QY8BZNZXBE46QBEXEW` Task 4 gate: 10ms, HTML 30s, ntfy 5/5 — 7h · current
- `e01M067AZ3VJM73XP82REZEZ1QC` Task 4 channel: ntfy 2383c10c…, card creation only — 7h · current
- `e01M066XWD8F1SAWQ8HVXGW4J4Z` Task 3 docket gate: 8→1, FYI 0, withdrawal 1 poll — 7h · current
- `e01M065G8RTKVH6466KE7GQREGX` Task 2 passes its live gate on all five acceptance checks — 8h · current
- `e01M064YRS4S9NK7KW9NQ1JMDV2` Malformed or missing pinned timestamps silently fall back to ingest time — 8h · suspect

## OPEN

_None._

## ALERTS

- possible-contradiction `e01M05SA72DPRGNTY7GCN1P7CED` Rename the inbox surface to "docket"; keep inbox as hidden alias — 11h
- possible-contradiction `e01M05SA72DPRGNTY7GCPEX9W2N` I10–I12 added as hard invariants enforced in code and tests — 11h
- possible-contradiction `e01M05SA72DPRGNTY7GCQ2RASTP` Cards show verbatim quotes + clickable provenance, never extractor paraphrase — 11h
- possible-contradiction `e01M05SA72DPRGNTY7GD0TGHEX7` Scope freeze: relay, TUI, team mode, adapters need a §11 trigger measurement — 11h
- suspect `e01M05SPB2EMMC4F4PR0BDQA8S5` First watch treated historical sessions as live: 342 overlaps, 27 stomps, 12.9M… — 11h
- suspect `e01M05T4XG0RJWQTP25SYT4FH0B` Task 2 not passable: `spent` conflates extraction, differ, and archaeology — 11h
- suspect `e01M05T4XG0RJWQTP25T1K58B61` Confirm/reject only in event YAML; adapter unknowns undated, absent from status — 11h

## Intent × reality

| Intent | Age | Reality | State |
|---|---:|---:|---|
| `e01M065G8RTKVH6466KEB0FRJ8N` Commit the gate fixes as one spec-amended change, then start the docket | 8h | 1 evidence | in_flight |
| `e01M064DHNTXFM29PNWCW0VBBF2` Run a strict read-only gate proving each acceptance point with tests or state q… | 8h | 3 evidence | in_flight |
| `e01M0642VRXV9PCGA4NDJF92E2Y` Wire atomic budget reservations into every LLM call | 8h | 5 evidence | in_flight |
| `e01M05VFAW9A783PMZZEER0G6FX` Second pass on rollover, double-settlement, migration; then run wider suite | 11h | 4 evidence | in_flight |
| `e01M05V5HWA6TFT0A0KDY1QV6C0` Add transactional reservation + settlement accounting in internal/state with co… | 11h | 9 evidence | in_flight |
| `e01M05V0G3Q9F41V62P6R5QV53G` Fix migration monotonicity, restore from D1 boundary, rerun cycle before passin… | 11h | 3 evidence | in_flight |
| `e01M05TPA9ZK4C8RGC3KBPV2MRF` Re-run the live gate: install binary, normalize spend category, restart watcher | 11h | 3 evidence | in_flight |

## Decisions

### e01M068ECYE067WF6BH7F26VC3D — Cursor v1 stays CLI-only: desktop 0 vs CLI 44 in 7d  `active`
> window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.

window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.

_source: human cli:note · confidence: 1.00_

### e01M064DQTWYDVGGAE3M5QRTGME — Re-evaluate the current tree instead of carrying the prior gate verdict forward  `active`
> The earlier gate’s three blockers are the right pressure points, but the checkout has moved: reservation callers and the neutral-workdir behavior now have new code and tests. I’m re-evaluating the present tree, including untracked test files, instead of carrying that verdict forward.

The prior gate's three blockers were judged the right pressure points, but the checkout has since gained new code and tests for reservation callers and neutral-workdir behavior. The gate will therefore re-evaluate the present tree, including untracked test files, rather than reusing the earlier verdict.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L35 · confidence: 0.90_

### e01M05V5HWA6TFT0A0KDZY8S45K — Confine the cap/ratio admission fix to internal/state; no caller or spec changes  `active`
> I’m implementing this entirely inside `internal/state`

The reservation/settlement work is scoped entirely to internal/state rather than changing call sites or the specification. This keeps the enforcement change contained and closes off the alternative of reshaping the caller-facing API or amending the spec to fix over-admission.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L188 · confidence: 0.82 · tags: internal/state/** · evidence: 2_

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

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.93 · tags: JOURNAL_SPEC.md · evidence: 5 · pairs-with: e01M05SA72DPRGNTY7GCPEX9W2N_

### e01M05SA72DPRGNTY7GCQ2RASTP — Cards show verbatim quotes + clickable provenance, never extractor paraphrase  `possible-contradiction`
> Design consequences: cards show verbatim quotes + clickable provenance, never the extractor's paraphrase or reasoning; high-magnitude cards carry one assumptions line; no other friction, ever.

Decision cards must render verbatim quotes with clickable provenance chips (session line / commit / entry) and must never show the extractor's paraphrase, summary, or reasoning. Reason: system-generated explanations are advocacy and increase acceptance of wrong content, while clickable sources reduce over-reliance. One "accepting this assumes: X" line is allowed on high-magnitude cards only; no o…

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.92 · tags: docket/** · pairs-with: e01M05SA72DPRGNTY7GCN1P7CED, e01M05SA72DPRGNTY7GCPEX9W2N_

### e01M05SA72DPRGNTY7GCPEX9W2N — I10–I12 added as hard invariants enforced in code and tests  `possible-contradiction`
> Add these three as I10–I12, hard law, enforced in code and tests, not convention:

Three new spec invariants, ranking as hard law rather than convention: I10 docket holds only items answerable by 1–3 discrete verbs (nothing FYI-shaped); I11 every card carries a machine-checkable, printed withdrawal condition and the docket keeps no history/counts/badges; I12 hard cap of seven cards, and sustained volume or an unneeded push is logged as system failure, never user workload.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.93 · tags: JOURNAL_SPEC.md, docket/** · evidence: 5 · pairs-with: e01M05SA72DPRGNTY7GCN1P7CED, e01M05SA72DPRGNTY7GCQ2RASTP, e01M05SA72DPRGNTY7GD0TGHEX7_

### e01M05SA72DPRGNTY7GCN1P7CED — Rename the inbox surface to "docket"; keep inbox as hidden alias  `possible-contradiction`
> Rename the surface — vocabulary is a forcing function against email-drift. It's a docket of decision cards (clew docket, with inbox as a hidden alias for muscle memory).

The decision surface is renamed inbox → docket, with "inbox" kept only as a hidden alias for muscle memory. Reason: vocabulary is a forcing function against email-drift — an inbox invites FYI accumulation, unread counts, and backlog; a docket is a list of items awaiting a ruling. The docket is the only surface that carries verbs.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.94 · tags: docket/**, inbox/** · pairs-with: e01M05SA72DPRGNTY7GCPEX9W2N, e01M05SA72DPRGNTY7GCQ2RASTP_

### e01M05RHSWXDNR10P1PY8ERYA9S — Dogfood metrics predeclared; D0 snapshot recorded  `active`
> Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.

Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.

_source: human cli:note · confidence: 1.00_

### e01M04WCGJS9FS7FQB0YFX9DTYG — Name the system clew (owner decision)  `active`
> Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.

_source: human cli:note · confidence: 1.00_

## Findings

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

### e01M065G8RTKVH6466KE7GQREGX — Task 2 passes its live gate on all five acceptance checks  `current`
> Task 2 now passes its live gate: 52 automatic session entries, 0 delivered/unneeded pushes, monotonic cursors, exact installed binary, and no active adapter/LLM errors.

A live gate run on Task 2 passed: 52 automatic session entries recorded, 0 delivered-but-unneeded pushes, cursors monotonic, the exact installed binary in use, and no active adapter or LLM errors. This is the verdict that unblocked committing the gate fixes.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988 · confidence: 0.93_

### e01M064YRS4S9NK7KW9NQ1JMDV2 — Malformed or missing pinned timestamps silently fall back to ingest time  `suspect`
> Source time: malformed/missing pinned timestamps silently become ingest `now`.

When a source record's pinned timestamp is missing or malformed, the adapter/extract path substitutes the ingest-time `now` without signalling, so entries get fabricated source times. Located at adapters.go:151 and extract.go:264; flagged as a gate blocker.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252 · confidence: 0.92 · tags: internal/adapters/**, internal/extract/** · evidence: 1_

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

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441 · confidence: 0.87 · tags: internal/**, cmd/clew/** · evidence: 5 · taint: tool_result_

### e01M05T4XG0RJWQTP25SYT4FH0B — Task 2 not passable: `spent` conflates extraction, differ, and archaeology  `suspect`
> `spent` combines extraction, differ, and archaeology; it is not extraction-only.

The dogfood audit judged Task 2 not passable yet. The budget `spent` counter mixes extraction, differ, and archaeology tokens, so the predeclared extraction-only cost metric cannot be read from it. Separating cost by kind is a prerequisite before the Task 2 gate can be honestly evaluated.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441 · confidence: 0.90 · tags: internal/**, cmd/clew/** · evidence: 5 · taint: tool_result_

### e01M05SPB2EMMC4F4PR0BDQA8S5 — First watch treated historical sessions as live: 342 overlaps, 27 stomps, 12.9M…  `suspect`
> First watch misclassified historical sessions as live, producing 342 overlaps, 27 stomps, and 12,895,847 observed tokens.

Measured fallout of the first watch run misclassifying pre-existing historical sessions as live: 342 overlaps, 27 stomps, and 12,895,847 observed tokens. This quantifies the historical-session storm previously recorded qualitatively as an I12 failure.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335 · confidence: 0.91 · tags: cmd/clew/** · evidence: 2 · taint: tool_result_

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

### e01M04XVQEEZ38J9TC5NZNKC16B — Run the live fidelity gate (RealProvider) on a machine with provider keys  `answered`
> If you cannot reach P>=0.9 / R>=0.75, stop all other work and report that the kill criterion fired — that outcome is a valid and useful result, and it is written down on purpose

Hermetic pipeline passes, but the go/no-go gate — P≥0.9/R≥0.75 within 5 instruction iterations vs ratified labels — has never run against a live provider. Until it does, extraction quality (the stated existential risk) is unmeasured; the kill criterion is theoretical. Needs claude/codex CLI or OpenAI-compatible key; env flag was RESTART_FIDELITY=1 pre-rename — verify.

_source: session chat:cursor-cloud-agent/stratura-strategy-2026-08-15 · confidence: 1.00 · tags: acceptance/**_

## Intents

### e01M065G8RTKVH6466KEB0FRJ8N — Commit the gate fixes as one spec-amended change, then start the docket  `in_flight`
> I’m committing the gate fixes as one spec-amended change, then moving to the docket.

After the Task 2 live gate passed, the plan is to land all gate fixes as a single commit that also amends the spec, and then move on to the docket surface work.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988 · confidence: 0.90 · evidence: 1_

### e01M064DHNTXFM29PNWCW0VBBF2 — Run a strict read-only gate proving each acceptance point with tests or state q…  `in_flight`
> I’m doing a strict read-only gate: first loading Clew’s generated context and relevant prior memory, then I’ll inspect the current diff and prove each acceptance point with tests or direct state queries. I’ll report only a blocker or a PASS with exact evidence.

Plan for this session: load Clew's generated context and prior memory, inspect the current diff, and prove each acceptance point via tests or direct state queries — reporting only a blocker or a PASS with exact evidence, making no writes.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L16 · confidence: 0.88 · evidence: 3_

### e01M0642VRXV9PCGA4NDJF92E2Y — Wire atomic budget reservations into every LLM call  `in_flight`
> I’m wiring the new atomic budget reservations into every LLM call next; this closes the remaining race where live extraction and backfill could both spend against the same allowance.

Next step: route every LLM call through the new atomic budget reservations, so live extraction and backfill can no longer both spend against the same allowance. Closing this race is the stated purpose of the change.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1577 · confidence: 0.88 · tags: internal/state/** · evidence: 5_

### e01M05VFAW9A783PMZZEER0G6FX — Second pass on rollover, double-settlement, migration; then run wider suite  `in_flight`
> I’m doing a second pass for rollover, double-settlement, and migration behavior before running the wider suite.

Before treating the internal/state reservation/settlement work as done, do a second review pass covering rollover, double-settlement, and migration behavior, then run the wider test suite beyond the state package.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L268 · confidence: 0.85 · tags: internal/state/** · evidence: 4_

### e01M05V5HWA6TFT0A0KDY1QV6C0 — Add transactional reservation + settlement accounting in internal/state with co…  `in_flight`
> a transactional reservation record plus settlement accounting, with contention tests that prove the cap/ratio cannot be over-admitted

Commitment to implement a transactional reservation record plus settlement accounting inside internal/state, accompanied by contention tests that demonstrate the cap/ratio cannot be over-admitted under concurrent access.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L188 · confidence: 0.88 · tags: internal/state/** · evidence: 9_

### e01M05V0G3Q9F41V62P6R5QV53G — Fix migration monotonicity, restore from D1 boundary, rerun cycle before passin…  `in_flight`
> restoring from the D1 boundary, and will rerun the cycle before calling Task 2 passed.

After the first upgraded live cycle exposed a cursor-rewind defect, the watcher was stopped. Committed follow-up work: correct the cursor migration, restore state from the D1 boundary, and rerun the live cycle. Task 2 will not be declared passed until that rerun is clean.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1437 · confidence: 0.92 · evidence: 3_

### e01M05TPA9ZK4C8RGC3KBPV2MRF — Re-run the live gate: install binary, normalize spend category, restart watcher  `in_flight`
> I’m moving back to the live gate now: install this exact binary, normalize the one dogfood spend category, restart the watcher, then verify one full tail/poll cycle has no historical replay, false sessions, or false pushes.

Next step after the unit path went green is a live-gate run: install the exact built binary, normalize the single dogfood spend category, restart the watcher, then observe one complete tail/poll cycle. Acceptance is negative evidence — no historical replay, no false sessions, no false pushes in that cycle.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1323 · confidence: 0.88 · evidence: 3_

### e01M05TBYJXEW5N5FE398NJ4RD7 — Tighten live-enrollment/backfill boundary and add failure telemetry  `in_flight`
> I’m tightening the live-enrollment/backfill boundary and failure telemetry now

Work in progress to harden the boundary between live enrollment and backfill, and to add telemetry for failures. This follows the recorded pass of Task 1 and the cursor/push/adapter failures observed in the first dogfood run.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041 · confidence: 0.85 · evidence: 3_

### e01M05SPB2EMMC4F4PR0BSPFZKZ — Narrow fix for the watch storm: transactional baselines, source-time, bounded c…  `in_flight`
> Root cause and narrow fix sent to parent: transactional live baselines, source-time sessions, and bounded separate backfill cursor.

The semantics investigation reported a root cause and a narrow fix to the parent: make live baselines transactional, use source-time (not observation-time) session timestamps, and give backfill its own bounded cursor separate from live watch. This is the proposed work to make backfill and live watch disjoint.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335 · confidence: 0.62 · tags: cmd/clew/**, internal/** · evidence: 11 · taint: tool_result_

### e01M05SG3NTP5W2JX7Y6MP1F1K6 — Add cursor migration, complete-record offsets, and fixed historical upper bound  `in_flight`
> adding a one-time migration, complete-record offsets, and a fixed historical upper bound

Committed follow-up work before the Task 2 commit can land: a one-time cursor migration for existing installs, offsets that always fall on complete JSONL record boundaries, and a fixed upper bound on the historical backfill range.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806 · confidence: 0.88 · evidence: 5_

### e01M05SA72DPRGNTY7GCXARNKYH — Implement the decision card and enforce I10–I12 in renderer and tests  `in_flight`
> Enforce I10–I12 in the renderer and in tests: a synthetic FYI item must be unrenderable; an 8th card must collapse to the overflow-failure card; a resolved stomp must withdraw within one poll cycle.

Commitment (Task 3): build the docket card to the fixed anatomy — headline-as-question ≤80 chars, why-you strip with rule fired and ticking stall timers, verbatim-quote evidence rows with provenance chips, assumptions line on high-magnitude only, 1–3 verbs plus defer-until-event, printed withdrawal condition, ordering by blocking cost, designed empty state. Tests must prove the three invariants.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 · confidence: 0.90 · tags: docket/** · evidence: 1_

### e01M05S9SFKAAM813AR1EG3X3WR — Land the dogfood fixes after recording the Task 2 snapshot  `in_flight`
> I’m recording that snapshot and then landing the dogfood fixes.

Commitment to record the final Task 2 dogfood snapshot and then land the outstanding dogfood fixes in the codebase.

_source: session codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729 · confidence: 0.82 · evidence: 3_
