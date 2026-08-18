---
format: clew.seed/v1
digest: sha256:e4e3cfae2180aaf4f035dfa49436d86dec265616aeb97cfc6b6a73a4d499ee7a
snapshot:
    repository:
        id: r8631b465b9e40c83e4d3e137
        name: restart
        remote: https://github.com/maceip/clew.git
    journal_revision: sha256:791c744975067a137d8a36173ef09d714f6435aacd2f1d28dc3b17d4ccb6f080
    changed_at: 2026-08-18T23:10:35Z
    lifecycle:
        state: active
    topics:
        - 10ms
        - "113"
        - 1kb
        - 2383c10c
        - 30s
        - 32kb
        - "342"
        - "631"
        - "754"
        - "788"
        - absence
        - absent
        - accelerate
        - accelerated
        - accept
        - acceptance
        - accounting
        - act
        - adapter
        - adapters
        - add
        - added
        - additions
        - admission
        - adopt
        - adopted
        - advances
        - advertised
        - after
        - against
        - agent
        - agents
        - alarm
        - alert
        - alerts
        - algebra
        - alias
        - all
        - always
        - ambient
        - amended
        - amnesia
        - ancient
        - answer
        - any
        - api
        - app
        - apply
        - approve
        - archaeology
        - asked
        - assumptions
        - atomic
        - attached
        - attachments
        - attending
        - audit
        - auto
        - automatic
        - back
        - backfill
        - backfillcmd
        - backfired
        - baked
        - base
        - baseline
        - baselines
        - bearing
        - became
        - before
        - begin
        - behind
        - bet
        - big
        - bin
        - binaries
        - binary
        - birth
        - block
        - blockers
        - blocks
        - bootstrap
        - born
        - bound
        - boundary
        - bounded
        - branch
        - breaks
        - broken
        - budget
        - build
        - burst
        - but
        - bypass
        - bytes
        - cache
        - call
        - caller
        - callers
        - calm
        - came
        - can
        - cap
        - caps
        - card
        - cards
        - carried
        - carry
        - carrying
        - category
        - caught
        - census
        - ceremony
        - change
        - changes
        - channel
        - chat
        - check
        - checks
        - class
        - claude
        - clean
        - cleaning
        - clew
        - cli
        - clickable
        - clone
        - closed
    decisions:
        - entry:
            id: e01M04WCGJS9FS7FQB0YFX9DTYG
            type: decision
            title: Name the system clew (owner decision)
            body: 'Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder''s unilateral restart from §12.1.'
            quote: 'Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder''s unilateral restart from §12.1.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T08:53:03.193378Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05RHSWXDNR10P1PY8ERYA9S
            type: decision
            title: Dogfood metrics predeclared; D0 snapshot recorded
            body: 'Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.'
            quote: 'Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:05:16.70191Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05SA72DPRGNTY7GCN1P7CED
            type: decision
            title: Rename the inbox surface to "docket"; keep inbox as hidden alias
            body: 'The decision surface is renamed inbox → docket, with "inbox" kept only as a hidden alias for muscle memory. Reason: vocabulary is a forcing function against email-drift — an inbox invites FYI accumulation, unread counts, and backlog; a docket is a list of items awaiting a ruling. The docket is the only surface that carries verbs.'
            quote: Rename the surface — vocabulary is a forcing function against email-drift. It's a docket of decision cards (clew docket, with inbox as a hidden alias for muscle memory).
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:36.621Z
            confidence: 0.94
            tags:
                - docket/**
                - inbox/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M05SA72DPRGNTY7GCPEX9W2N
            type: decision
            title: I10–I12 added as hard invariants enforced in code and tests
            body: 'Three new spec invariants, ranking as hard law rather than convention: I10 docket holds only items answerable by 1–3 discrete verbs (nothing FYI-shaped); I11 every card carries a machine-checkable, printed withdrawal condition and the docket keeps no history/counts/badges; I12 hard cap of seven cards, and sustained volume or an unneeded push is logged as system failure, never user workload.'
            quote: 'Add these three as I10–I12, hard law, enforced in code and tests, not convention:'
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:36.621Z
            confidence: 0.93
            tags:
                - JOURNAL_SPEC.md
                - docket/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M05SA72DPRGNTY7GCQ2RASTP
            type: decision
            title: Cards show verbatim quotes + clickable provenance, never extractor paraphrase
            body: 'Decision cards must render verbatim quotes with clickable provenance chips (session line / commit / entry) and must never show the extractor''s paraphrase, summary, or reasoning. Reason: system-generated explanations are advocacy and increase acceptance of wrong content, while clickable sources reduce over-reliance. One "accepting this assumes: X" line is allowed on high-magnitude cards only; no o…'
            quote: 'Design consequences: cards show verbatim quotes + clickable provenance, never the extractor''s paraphrase or reasoning; high-magnitude cards carry one assumptions line; no other friction, ever.'
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:36.621Z
            confidence: 0.92
            tags:
                - docket/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M05SA72DPRGNTY7GD0TGHEX7
            type: decision
            title: 'Scope freeze: relay, TUI, team mode, adapters need a §11 trigger measurement'
            body: Relay server, TUI/native apps, team mode, semantic code graphs, treemaps, new adapters, and orchestration are frozen. Building any of them requires first citing a §11 trigger measurement in the journal and stopping for owner review — measurement, not enthusiasm, unfreezes scope.
            quote: 'Frozen (build nothing here without citing a §11 trigger measurement in the journal and stopping for owner review): relay server, TUI/native apps, team mode, semantic code graphs, treemaps, new adapters, orchestration.'
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:36.621Z
            confidence: 0.93
            tags:
                - JOURNAL_SPEC.md
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M05SG3NTP5W2JX7Y6MG00HQJ
            type: decision
            title: Hold the Task 2 commit; make backfill and live watch disjoint by construction
            body: 'Rather than commit Task 2 and patch later, the commit is held while three fixes land: a one-time cursor migration for upgrading users, complete-record offsets so init never baselines mid-record, and a fixed historical upper bound so backfill and live watch cannot overlap — disjoint by construction rather than by runtime check.'
            quote: I’m holding the Task 2 commit and adding a one-time migration, complete-record offsets, and a fixed historical upper bound so backfill and live watch are disjoint by construction.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:21:49.754Z
            confidence: 0.9
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05SVGK1Q2MR34Y3CMHR7DXM
            type: decision
            title: 'No cursor translation: keep `extract:` for live, add bounded `backfill:` for hi…'
            body: Instead of migrating or stacking cursors, the live watcher retains the existing `extract:` cursor while explicit history backfill gets a separate bounded `backfill:` cursor, with `history-end` freezing the boundary between them. Chosen because it preserves v1 pending work across upgrade, keeps live and history from overlapping, and eliminates the risky cursor translation step.
            quote: the live watcher keeps the existing `extract:` cursor, explicit history gets a new bounded `backfill:` cursor, and `history-end` freezes their boundary
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L934
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:28:03.425Z
            confidence: 0.93
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05TBYJXEW5N5FE397XYMEHY
            type: decision
            title: Treat first dogfood run's cursor/push/adapter failures as required failure sign…
            body: The first dogfood run surfaced real cursor, push, and adapter failures. The assistant is treating that run as the required failure signal rather than as an acceptance run, so those failures do not block recording Task 1 as passed but do drive the current work on the live-enrollment/backfill boundary and failure telemetry.
            quote: the first dogfood run exposed real cursor, push, and adapter failures, so I’m treating that run as the required failure signal—not as acceptance
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:37:02.045Z
            confidence: 0.82
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05TG628KAGZ9HCBQSAG562P
            type: decision
            title: Separate session-extraction budget from one-time archaeology budget
            body: 'Extraction budget for live sessions is now tracked separately from the one-time historical archaeology budget. Reason: deriving the archaeology allowance as a ratio of observed session tokens yields zero budget at cold start (zero observed sessions), making backfill of historical docs mathematically impossible. Sits alongside the pre-/post-enrollment byte boundary.'
            quote: I’m also separating session-extraction budget from one-time archaeology, because applying a session-token ratio to cold-start docs makes archaeology mathematically impossible at zero observed sessions.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1200
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:39:20.776Z
            confidence: 0.9
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05V0G3Q9F41V62P6T44TC04
            type: decision
            title: Cursor migration must be monotonic — never rewind an existing cursor
            body: A divergent legacy cursor was rewound during migration, causing one duplicate re-extraction. The fix chosen is to make the migration monotonic so a migrated cursor can only move forward, closing off any migration path that repositions a cursor backward and replays already-extracted bytes.
            quote: I’m correcting the migration to be monotonic
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1437
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:48:15.351Z
            confidence: 0.78
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M05V5HWA6TFT0A0KDZY8S45K
            type: decision
            title: Confine the cap/ratio admission fix to internal/state; no caller or spec changes
            body: The reservation/settlement work is scoped entirely to internal/state rather than changing call sites or the specification. This keeps the enforcement change contained and closes off the alternative of reshaping the caller-facing API or amending the spec to fix over-admission.
            quote: I’m implementing this entirely inside `internal/state`
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L188
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:51:01.002Z
            confidence: 0.82
            tags:
                - internal/state/**
            env: null
            affects: []
          status: active
        - entry:
            id: e01M064DQTWYDVGGAE3M5QRTGME
            type: decision
            title: Re-evaluate the current tree instead of carrying the prior gate verdict forward
            body: The prior gate's three blockers were judged the right pressure points, but the checkout has since gained new code and tests for reservation callers and neutral-workdir behavior. The gate will therefore re-evaluate the present tree, including untracked test files, rather than reusing the earlier verdict.
            quote: 'The earlier gate’s three blockers are the right pressure points, but the checkout has moved: reservation callers and the neutral-workdir behavior now have new code and tests. I’m re-evaluating the present tree, including untracked test files, instead of carrying that verdict forward.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L35
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:32:46.428Z
            confidence: 0.9
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M065SK8W1ZT32KZF8YGGTP7W
            type: decision
            title: Alerts self-clean via one scoped reconcile with an explicit withdrawal condition
            body: 'Rather than only inserting alerts (which never closed) and keying them on mutable prose, alert handling became a single state-level reconcile: upsert the active alerts for a repo and auto-drop previously open differ-owned kinds absent from that poll, with each alert storing an explicit withdrawal condition.'
            quote: 'I’m shaping this as one state-level reconcile operation: upsert the active alerts for a repo and auto-drop previously open differ-owned kinds that are absent from that poll, with an explicit withdrawal condition stored on each alert.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L45
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:56:43.548Z
            confidence: 0.93
            tags:
                - state/**
                - differ/**
            env: null
            affects: []
          status: active
        - entry:
            id: e01M065T92NXY1ER6R73YCQNH84
            type: decision
            title: Recover docket empty-state and withdrawal wording from Task 3 journal source
            body: Rather than re-deriving the card semantics, the generated journal's Task 3 source pointer and fixed card decisions are used as the authoritative source to recover the missing empty-state wording and withdrawal semantics before the package API is locked.
            quote: The generated journal contains the Task 3 source pointer and the fixed card decisions. I’m using that exact source to recover the missing empty-state wording and withdrawal semantics before locking the package API.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L75
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:57:05.877Z
            confidence: 0.8
            tags:
                - docket/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M068ECYE067WF6BH7F26VC3D
            type: decision
            title: 'Cursor v1 stays CLI-only: desktop 0 vs CLI 44 in 7d'
            body: 'window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.'
            quote: 'window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:43:02.350185Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AQHSRF4DVDYZ989M4DVHYX
            type: decision
            title: Lineage inheritance is explicit; only promoted laws auto-carry
            body: 'clew will never auto-carry project lore into a new repo. Rationale: a wrong lineage guess poisons a fresh project worse than inheriting nothing at all. Owner laws are the sole exception and may be injected automatically, because the promotion step already certified each law as project-agnostic. This is the governing reason behind invariant I13.'
            quote: lore inheritance was made explicit because a wrong lineage guess poisons a fresh project worse than no inheritance at all — laws are safe to auto-carry precisely because promotion certified them project-agnostic.
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-18T15:24:00.143Z
            confidence: 0.95
            tags:
                - .clew/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M0AQHSRF4DVDYZ989MTZQV7E
            type: decision
            title: SEED.md is watcher-maintained continuously, never generated on demand
            body: 'The watcher regenerates SEED.md alongside context.md on every journal change rather than building it when a restart is requested. Reason: the carry-kit must already exist before anyone wants to restart, so a seed is never missing or stale at the moment it is needed.'
            quote: the watcher maintains SEED.md continuously alongside context.md — regenerated on journal change, never on demand. The carry-kit always predates the urge to restart.
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-18T15:24:00.143Z
            confidence: 0.94
            tags:
                - .clew/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M0AQHSRF4DVDYZ989RNZVC46
            type: decision
            title: '`clew from` is the one explicit lineage command; never automatic'
            body: Pulling a predecessor's seed (decisions, findings, graveyard, exhibits, organ-bank pin) happens only via an explicit `clew from <repo>`; with no args it lists candidates ranked by recency and topic overlap, showing what each would carry. It can run at birth or later, un-carrying is recorded as a reject so carried entries keep provenance, and the birth card may only suggest `clew from X` on name/t…
            quote: Runnable at birth or any time later; un-carrying is a reject (carried entries keep provenance like everything else). Never automatic. The birth card may suggest clew from X on blatant name/topic overlap — suggest only, never act.
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-18T15:24:00.143Z
            confidence: 0.93
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AQHSRF4DVDYZ989W6K185N
            type: decision
            title: Owner laws live in an owner-scope journal with a ≤1KB injection budget
            body: Laws are stored as an owner-scope journal synced like any other. Findings become laws through an explicit `clew journal promote <id>`, with the extractor proposing promotion when a finding's content is project-agnostic. The resulting law set is capped at a ≤1KB injection into every project's context, permanently.
            quote: an owner-scope journal synced like any other; findings promoted via clew journal promote <id>; extractor proposes promotion when a finding's content is project-agnostic; ≤1KB injection budget into every project's context, forever.
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-18T15:24:00.143Z
            confidence: 0.92
            tags:
                - .clew/**
            env: null
            affects: []
          status: possible-contradiction
        - entry:
            id: e01M0ASQGJZPFJD3FSRT7P4HX44
            type: decision
            title: Owner-configured cloud environments are full clew nodes
            body: 'Owner corrected the push-only sandbox assumption: cursor/codex/claude cloud environments are configurable (install scripts, MCP, skills) and can run the Go binary. Cloud write path = provision the environments you own. journal-proposal.yaml is PARKED (trigger-gated for unconfigurable third-party envs only).'
            quote: we have installable skills, MCP, and I can configure the environments cursor, codex, and claude agents run in the cloud which can absolutely include our golang services
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0ASQGK2ZK03KRP4WPAQG0MJ
            type: decision
            title: 'Bet: restart-accelerated development plus drift guardrails, one shared substrate'
            body: 'Owner bet the farm on (A) glanceable intent-reality drift for humans and (B) restart acceleration: new repo births with genesis docs, old code vendored as lessons. Guardrails lower restart NEED; seeds lower restart COST; both attack unrecorded divergence. Restart verbs stay pull-only forever.'
            quote: lets bet the farm on strong coordination via our journal and intent < - > current reality drift
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0ASQGK4FMKRNCNR91KJ06JD
            type: decision
            title: 'Restart machinery must be zero human effort: ambient inheritance, opt-out'
            body: 'Lesson from substrate x2: reuse that costs effort at the clean-slate moment gets skipped. Therefore: SEED.md maintained continuously; birth detection auto-injects owner laws only; full manifest ceremony stays optional. Anything required at project birth is a bug (I13).'
            quote: the solution we create out of clew needs to make the restart acceleration zero effort from the human or zero cognitive load
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0ATYJG615JE6BV5MG5RAF9Z
            type: decision
            title: clew from must accept multiple parent projects, with strand selection
            body: 'Owner ruling: inheritance is multi-parent. `clew from A B` unions seeds; `--tags <globs>` selects strands per parent; runnable repeatedly. Each carried entry keeps per-parent provenance; disagreements between parents surface as possible-contradiction cards for human arbitration, never silent merge. Genesis records multiple lineage links (the forest gains merge nodes).'
            quote: clew from needs to support _multiple_ projects as inputs, the human may want parts from multiple
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:22:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AV0H7RY7DG79VCMHPEMPJP
            type: decision
            title: 'Glance direction ruling: graphic, two zooms — deferred behind single-project'
            body: 'Owner direction: the glance becomes a graphic (project view: status-colored intent tiles, drift strip, docket badge; fleet view: hot-project tiles, dormant collapsed). Static self-contained HTML, no server. DEFERRED: build only after the single-project version works well.'
            quote: the glance view _graphic_ not text, needs to support a one project view, and global view
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:23:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AV0H7T9P69CNPB56MRAG8V
            type: decision
            title: 'Sequencing: single-lineage from + one-project glance FIRST; fleet/multi later'
            body: 'Owner ruling: no rush on multi-parent from or the fleet view. Get clew from working well with one project lineage and the glance UI working well for one project before any scaling work — else scope creep triggers the restart urge (owner: risk of clew- from-clew). Multi-parent and fleet rulings stand as destination, gated on the single versions working well.'
            quote: we need to get clew from working well with just one project lineage, and the glance infra ui working with just one project well first, else we risk scope creep and me clew from clew
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:23:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AVAR0N8ZCCN1VJW3GQ5PF4
            type: decision
            title: Restart-with-mutation is the flagship advertised workflow, not a failure mode
            body: 'Owner ruling: restart-with-mutation is the flagship, advertised workflow. The old negative was retelling pain (re-briefing a blank agent), not rebirth. Direction: clew from <parent> "<mutation>" carries the seed, makes the mutation the genesis charter, and flags carried entries it contradicts — day-zero docket = the design debate, pre- computed. Gated behind single-lineage sequencing.'
            quote: i allow myself to clew from clew "but with cloud agent witnesses" ... something we _actually_ advertise as a path
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:28:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AVWK3HH2R9M55FQSAWZFF1
            type: decision
            title: 'I9 frugality replaced: listening completeness is the invariant, cost is a dial'
            body: 'Owner ruling: price sensitivity was an agent assumption, never stated. Replace the 2% ratio with an owner-set ceiling plus a hard floor above the largest atomic request; extraction must never deadlock. Spend stays a visible meter. This also resolves the URGENT budget card''s direction.'
            quote: you assume i'm price sensitive and are using token cost as being prohibitive but ive never mentioned we need to make this work cheaply
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:38:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AVXA7EH2KD1BPZ4GJNKN67
            type: decision
            title: 'Witness-node role adopted: always-on ear with owner API creds, sequenced'
            body: 'Owner sgtm: one always-on clew node (owner infra) whose sensors are API pollers with owner account creds — witnesses cursor/codex cloud sessions live with zero agent cooperation, runs extraction centrally, sole writer of projections (kills that conflict class). Git stays the only required meeting point; degrade-to-baseline preserved. Build gated behind single-lineage sequencing.'
            quote: that sgtm with one wrinkle still in my brain, how does this system work for 2,10,100 projects
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:38:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AXCM53D36PDN0H85S0WE3W
            type: decision
            title: Mind-plane freshness is vendor-neutral; hooks accelerate, never carry it
            body: 'Owner ruling: a returning human must land on agents that already know the recent journal, and this must hold for a bare ollama model with no hook system. Per-vendor hooks may improve latency but the floor must work for anything that emits model API calls. Bar: current-at-next-interaction, zero human homework, any harness — including ollama running deepseek4-flash on a laptop.'
            quote: this should work if i decide to spin up ollama locally with deepseek4-flash on my laptop
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:00:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AXF5FM9RDFP817NF7QW5BN
            type: decision
            title: Human-facing surface must reduce to the desires it satisfies
            body: 'Owner ruling: the internal design may be intricate, but the human-visible vocabulary must collapse to the desire set: it remembers; every agent starts knowing; I can look up and see; it asks me only when it must; starting over loses nothing. Any feature that cannot be filed under one of these is out. The agent carries the machinery; the human carries five sentences.'
            quote: its already too much to understand given how it should reduce down to the simple set of human desires it satisfies
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:06:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AXMXK3SAATFDFAYZ932TPC
            type: decision
            title: 'Two registers, one memory: calm words for humans, hard words for agents'
            body: 'Owner ruling: human-facing vocabulary must avoid fear-attached words (law, state, violation). But register is a rendering choice, not a softening of the contract: wherever soft words would let agents wiggle out of a constraint, the agent-facing rendering keeps the harsh form. Same entries, two renderings; hardness is judged by compliance, softness by human calm.'
            quote: however -- if using softer language means agents would "wiggle" out of those constructs, then we keep the harsher language
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:09:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AXRKQ8C7FZNKARW83CMBMX
            type: decision
            title: The five promises are the foundation (owner ratified)
            body: 'Owner ratified the product''s entire human-facing surface: (1) it remembers what we decide; (2) every agent starts already knowing it; (3) you can look up and see; (4) it taps your shoulder only when something needs you; (5) starting over loses nothing. Every feature must file under exactly one promise or it is not built. Vocabulary beyond these is machinery, never surfaced.'
            quote: yeah those 5 things are the foundation, agreed
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:12:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AXXKMNNKKY721HJ3REN3KH
            type: decision
            title: Freshness is owed at contact points; a task runs on its snapshot
            body: 'Owner refinement: a running task is never interrupted or mutated by concurrent decisions — it finishes on the snapshot it started with. Currency is owed at the next human contact: a message typed after returning lands on a mind that already has the delta. Hooks fire at that boundary; the proxy injects only on a new human message. Urgent items route to the human, who may stop the task.'
            quote: i dont expect an agent on task to stop mid task and change based on a cloud agent decision i made at the same time
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:14:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AXZTTHG5FKXETX0X8PR6EX
            type: decision
            title: 'On finish, check in first: reconcile against the delta before next steps'
            body: 'Owner ruling completing the task lifecycle: when an agent finishes, it must sync the journal and reconcile its output against decisions that landed mid-flight before concluding or picking next steps — state contradictions explicitly, then close. Stop/AfterAgent hooks make this enforceable on claude/codex/gemini; elsewhere it is convention plus the glance flagging stale finishes.'
            quote: but what i do expect is that after it finishes, to check in first and then figure out what to do based on the new synced state
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:15:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AY6VWJV133F811JYAANJPE
            type: decision
            title: 'Stale finish: know and tell, never act — the reconcile is read-only'
            body: 'Amendment to the finish check-in: an agent whose finished work was obsoleted mid-flight must not remove, redo, or touch anything on its own — no action of any kind. The check-in only installs knowledge: at the next human prompt it must say the work is deprecated/obsoleted/wrong and why, unless the human already resolved it elsewhere. Interpretation is the human''s call. Owner will stress-test.'
            quote: finishes, then syncs, should _NOT_ then remove the work it just did, or do anything for that matter
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:19:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AYQ3VJM71SMM1YQFYE1X61
            type: decision
            title: 'Knowledge Merge at finish: glanceable apply/defer list, external memory'
            body: 'Owner design: at finish, one colored glanceable list — top unapplied changes (code, intent, knowledge), one line each, entry-linked, one-keystroke apply/defer. See-once by decision id; defer compresses to a nagging count, never re-shown as new. External memory for a forgetful human: recognition over recall, per the HCI findings. It is the docket rendered at the finish boundary.'
            quote: this list cannot be verbose, it needs to be glancable, and it serves the "humans forget" thing as well
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:28:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AZ0K17AA5D9P7KZDSPJQSY
            type: decision
            title: Merge lines must pass the amnesia test; verbs are apply/explain/defer
            body: 'Amendment: each merge line must be readable by a human who forgot the conversation entirely — references glossed inline (the five promises appear as five words), machinery nouns translated, no dangling ''the budget''. Per line: apply, explain (prints body + the owner''s verbatim quote + link, then re-offers), defer. Footer gains apply-all. Explain works because one''s own words restore memory.'
            quote: these diff entries need to be something the human can read after maybe a day where he totally forgot the convo he had with you
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:33:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AZ4BGF1VC0VSXA05VYVEQ3
            type: decision
            title: 'Explain is live: the attending agent reads the entry and explains'
            body: 'Refinement: the merge diff encodes nothing but lines, entry ids, and verbs. Pressing explain hands the entry to the agent already present at the finish boundary — it reads the journal, quotes the owner''s words, and explains what applying means for the work at hand, answering follow-ups conversationally. clew stays the bookkeeper of see-once and defer state; the agent is the explain engine.'
            quote: so you dont need to encode the "explain in more detail in the kowledge diff" itself
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:35:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AZ7JBPJRNHFJQC6WEEQB9Y
            type: decision
            title: 'Silence is the signal: an absent merge means truly nothing new'
            body: 'Owner property: when no knowledge diff appears, the human may trust that nothing new landed anywhere — silence is the all-caught-up signal. For that trust to hold, silence must be earned: a broken watcher, stale sync, or failed check must announce itself distinctly and can never render as an empty diff. Quiet means verified-quiet. Nothing-new and could-not-check are never the same screen.'
            quote: if there is no knowledge diff shown, the human knows nothing new has been added somewhere else
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:37:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0AZDDYC49Y1741H7W74Y1QY
            type: decision
            title: 'Second tab: the intent gap — everything intended, not yet real'
            body: 'Owner design: next to the knowledge merge sits the intent gap — same glanceable, amnesia-proof list shape, listing intents with no evidence in reality (the absence machinery gets its human surface). Verbs: build (hand to the idle agent), explain (live), retire (a conscious no, kept with provenance). It converts forgetting into deciding — stratura''s unbuilt core would have topped it for weeks.'
            quote: the intent gap is a similar thing that lists simply all the crazy shit thats not yet implemented
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:40:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BEM04PAZ85R5YNRM2Y31Z6
            type: decision
            title: 'Broken states carry their verb: no unactionable warnings for humans'
            body: 'Owner correction from the first real merge/gap run: could-not-check lines handed the human a problem with no action. Rule: a broken state shown to a human must carry its fix verb (usually hand to the attending agent) or name who is already fixing it; problems only machinery can fix route to agents, never to human eyes. Earned silence stands — broken arrives actionable, like everything else.'
            quote: that is state the human cant immediately fix/help with
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:04:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BER1Q412DRFQYESPCN0Q30
            type: decision
            title: Lines are plain speech, no ids; near-duplicates fold; held items rest
            body: 'Owner corrections from the first real run: rendered lines confused (Let cloud agents that can only...) and full entry ids burned attention. Rules: plain spoken English, subject-first, one breath; no ids or codes on lines — identity lives behind explain; near-duplicates fold to one line; held-for-owner entries appear in no actionable list. The amnesia test stays the floor; this adds plain speech.'
            quote: language in each <knowledge merge/intent gap> thats confusing, e.g., "Let cloud agents that can only"
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:08:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BETRRAZ54063PRJK1JSQS7
            type: decision
            title: 'The finish message is a surface: what exists, where it lives, my next move'
            body: 'Owner correction: codex signed off in builder frame (Nothing was pushed. No real apply...) — accurate, meaningless to the human, and alarming: didn''t-push read as failure. Rule: the closing utterance speaks the human frame in plain words — what exists now, whether it is safely shared or local-only, what the human can say next — then shows the two screens. Compliance detail lives behind explain.'
            quote: 'when i read that my first thought is: "ok it didnt push anything so thats not good", then i dont know what the rest of it means'
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:10:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BEYP65CE70G0VVSX3PV01B
            type: decision
            title: 'Finished means shared: work ends pushed or PR''d; local-only is an alarm'
            body: 'Owner ruling closing the push gap: a task is not finished until the work is shared per repo convention — pushed to the branch or opened as a PR. Committed-but-local is an alarm state the finish message must name, never a resting state. Root cause on record: this norm lived in behavior and was never spoken, so the memory had nothing to inject; said once here, it reaches every agent forever.'
            quote: why do i need to tell it to push? with all the journals and dockets and intents and knowledge that should never be the case?
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:11:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BF0WMX264RA9D0VTM9R24K
            type: decision
            title: 'Entry ids are machine plumbing: never shown to or relayed through humans'
            body: 'Extension of plain-speech: the cloud agent printed raw ids in receipts and in prompts the human had to copy. The intent was verifiability — but verification is machine work, and agents holding the journal resolve plain words better than opaque codes. Rule: ids live in files, commits, and machine channels only; humans see words; agents are addressed in words and resolve entries themselves.'
            quote: you should never print "e01M0BER1Q412DRFQYESPCN0Q30" i dont know what it means and its super long -- what is the intent of showing that to the human?
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:13:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BFRYY1EW9CFMJNDJ5QH3M2
            type: decision
            title: Evidence settles merge lines; apply is never asked for finished work
            body: 'Owner found the merge asking him to bless work already built, tested, and pushed. Rule: the merge joins decisions to evidence; a decision whose demanded work is evidenced settles itself and shows once as settled-while-away. Apply is reserved for work not yet done or judgment only a human can make. Nothing auto-acts on the repo; settling is status computation, not action.'
            quote: why are there still 7 merges that need to take place ?  number 1 i thought you already merged?
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:25:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BFRYY2WVHJZ3R3TDV4CFTS
            type: decision
            title: The wording sweep covers every fear-attached word; docket stays by name
            body: 'Owner correction: the rename card narrowed to the single word law. The sweep is all fear-attached words wherever humans read — law, state, violation, enforcement and relatives — judged per word in context. Docket is explicitly approved and stays. Agent-facing hard register remains untouched where hardness prevents wiggle.'
            quote: why is 7 indexing on "law" when we discussed _all_ legal sounding words, with docket being ok
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:25:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
        - entry:
            id: e01M0BFY164YEV6DAEFVGGH18VT
            type: decision
            title: The limiter gates distillation timing, never sensing; failure is lag
            body: 'Owner challenge: a deaf agent is useless. Purpose on record: the limiter is not cost control — it protects shared rate limits and guards against runaway loops. Corrected design: sensing (tailing, recording) is free and never stops; only distillation may lag under pressure, shown as memory is N minutes behind, catching up when headroom returns. Deafness is impossible; nothing goes unrecorded.'
            quote: a deaf agent is useless and will lead to more problems for me to deal with later?
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T22:28:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: active
    findings:
        - entry:
            id: e01M04TDKJ79MEGBWTCSD8HTE5M
            type: finding
            title: Lock-and-notary pitches fail for any team whose job isn't security — not just …
            body: Lock-and-notary pitches fail for any team whose job isn't security — not just solo devs. Evidence framing must lead with anti-forgetting/throughput, never audit.
            quote: Lock-and-notary pitches fail for any team whose job isn't security — not just solo devs. Evidence framing must lead with anti-forgetting/throughput, never audit.
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T08:18:41.863449Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M04XVNN32W4FTFKKSKJSN01V
            type: finding
            title: varalys/lore owns session-provenance + git-sync; clew's edge is diff + absence
            body: 'varalys/lore ships session recording linked to commits, serverless git-remote sync, cross-tool memory over MCP; installs `lore` via brew/cargo. Same plumbing independently derived — commodity layer. Unoccupied, core to clew: typed distillation, intent×reality diff with absence, human steering surface, restart manifest. getlore.ai + Epic Lore make the name unownable.'
            quote: lore is retrospective provenance — whole-session storage answering "why does this line exist." Restart is prospective state — distilled typed entries answering "what's true about this project right now"
            utterance_by: assistant
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-16T08:19:00Z
            confidence: 0.95
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M04XVP9QBE041XCK43C78BVP
            type: finding
            title: Decision-dense sessions live on uncovered surfaces; manual notes = homework
            body: 'clew was designed in phone/cloud chats no watcher covers. An agent bridging that gap with prescribed manual notes made the product read as homework — an I1-violation smell, correctly caught. Non-discipline fixes: one export+backfill per key conversation; decisions echo into covered sessions, caught on first echo. §11 cloud-sensor trigger fires now for phone-heavy owners.'
            quote: ok not sure im liking this me needing to add to the journal each time a file in a repo is edited let me think about this
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-16T08:19:00Z
            confidence: 0.95
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M04XVPW5H5JCPS92CFQ4EBY3
            type: finding
            title: 'Independent verify: clean clone green; differ/poller/manifest lack unit tests'
            body: 'Second agent+machine, Go 1.26.3: build clean; full suite green; gates 1, 2-hermetic, 3 PASS; RealProvider SKIPs without keys; init push-deferral loud per spec. Gaps: no package tests for differ, poller, manifest, archaeology — coverage rides acceptance alone. Footgun: `journal note --help` ingested the flag as literal entry text; cleaned via reject.'
            quote: '--- PASS: TestAcceptance1_AbsenceDetection / TestAcceptance2_ExtractionFidelityPipeline / TestAcceptance3_RestartRoundTrip; SKIP: TestAcceptance2_RealProvider'
            utterance_by: assistant
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-16T09:20:00Z
            confidence: 1
            tags:
                - internal/**
                - acceptance/**
            env:
                host: cursor-cloud-vm
                hw: linux/amd64 go1.26.3
            affects: []
          status: current
        - entry:
            id: e01M04YAQN85TF6YQDP33VB4JQ0
            type: finding
            title: Foreign agents can read the journal but not write; contribution path unbuilt
            body: 'cursor[bot] was denied push to maceip/clew — correct posture: journal write = repo write, else anyone could poison context.md (§6.5 amplifier). But there is no sanctioned path for non-credentialed contributors; tonight''s delivery was a hand-rolled bundle. Options: document fork-PR onto clew/journal, or a `clew import <bundle>` verb landing entries pending human confirm.'
            quote: so if you cant even push to the public clew directory does that mean our first real test failed?
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-16T09:25:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M05REDJZAERAWRQG7349E7Y1
            type: finding
            title: Live fidelity gate passed on iteration 1
            body: 'Live fidelity gate iteration 1: P=0.91; R=0.83; decisions=6/7; findings=4/5; rejected=0; provider=claude; PASS.'
            quote: 'Live fidelity gate iteration 1: P=0.91; R=0.83; decisions=6/7; findings=4/5; rejected=0; provider=claude; PASS.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:03:25.791199Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M05RF74PBCZVAPRY49GCVM7H
            type: finding
            title: 'answer: Run the live fidelity gate (RealProvider) on a machine with provider ke…'
            body: P=.91 R=.83; D=6/7 F=4/5; reject=0; claude; iter=1; PASS
            quote: P=.91 R=.83; D=6/7 F=4/5; reject=0; claude; iter=1; PASS
            utterance_by: user
            source:
                kind: human
                ref: inbox:answer:e01M04XVQEEZ38J9TC5NZNKC16B
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:03:51.958265Z
            confidence: 1
            tags:
                - acceptance/**
            env: null
            affects: []
          status: current
        - entry:
            id: e01M05S9SFKAAM813AR1AY96QWH
            type: finding
            title: 'Final Task 2 dogfood snapshot: 0.113% extraction, 0:1 confirm:reject, 0 pushes'
            body: 'Final Task 2 dogfood measurement of clew: 6 automatic entries produced; 5,057 of 4,491,713 observed bytes extracted (0.113%); confirm-to-reject ratio 0:1; actual pushes 0 of 0; 1 adapter incident (the previously journaled D0 session storm); 0 parked items. Recorded after the Codex-format detection and watcher baseline repairs.'
            quote: 'Final Task 2 snapshot is: 6 automatic entries; extraction 5,057 / 4,491,713 observed = 0.113%; confirm:reject 0:1; actual pushes 0/0; adapter incidents 1 (the journaled D0 storm); parked 0.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:22.707Z
            confidence: 0.93
            tags: []
            env:
                host: local dogfood run
                dataset: clew dogfood sessions, 4,491,713 bytes observed
            affects: []
          status: current
        - entry:
            id: e01M05S9SFKAAM813AR1B8DXEYW
            type: finding
            title: Codex format now detected; watcher tracks only post-baseline bytes
            body: 'After repair, clew detects the current Codex session format: re-initialization found 3 sessions with large metadata, and the watcher now tracks only bytes written after the baseline, using source time rather than ingest time. This addresses the historical-session storm seen in dogfood D0.'
            quote: 'The current Codex format is now detected: re-init found 3 large-metadata sessions, and the watcher is tracking only post-baseline bytes with source time.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:22.707Z
            confidence: 0.85
            tags: []
            env: null
            affects:
                - e01M05RPB2N6QDXN1MP31SB64B2
          status: current
        - entry:
            id: e01M05SA72DPRGNTY7GCV0CBDK7
            type: finding
            title: Assumptions prompt cut over-reliance 42%→22%; stacked/delay friction backfired
            body: 'From the over-reliance literature (Buçinca cognitive-forcing lineage through Ghosh et al. 2026, n=214): an "accepting this assumes X" prompt reduced over-reliance from 42% to 22% with no added cognitive load, while delay-based and stacked friction backfired. This is why the assumptions line is the single permitted forcing function on clew cards and why no other friction is added.'
            quote: the assumptions prompt ("accepting this assumes X") reduced over-reliance from 42%→22% without added cognitive load, while stacked/delay-based friction backfired
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:36.621Z
            confidence: 0.85
            tags:
                - docket/**
            env:
                dataset: Ghosh et al. 2026 over-reliance study, n=214
            affects:
                - docket/**
          status: current
        - entry:
            id: e01M05SG3NTP5W2JX7Y6KZD5N6P
            type: finding
            title: 'Pre-commit review found three Task 2 blockers: cursor migration, backfill overl…'
            body: 'A pre-commit review of the Task 2 work surfaced three real blockers: users upgrading had no cursor migration path; backfill could overlap with live suffixes; and init could baseline in the middle of a partial JSONL record, producing a corrupt offset.'
            quote: 'Pre-commit review found three real blockers: upgrade users lacked cursor migration, backfill could overlap live suffixes, and init could baseline inside a partial JSONL record.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:21:49.754Z
            confidence: 0.93
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M05SPB2EMMC4F4PR0BDQA8S5
            type: finding
            title: 'First watch treated historical sessions as live: 342 overlaps, 27 stomps, 12.9M…'
            body: 'Measured fallout of the first watch run misclassifying pre-existing historical sessions as live: 342 overlaps, 27 stomps, and 12,895,847 observed tokens. This quantifies the historical-session storm previously recorded qualitatively as an I12 failure.'
            quote: First watch misclassified historical sessions as live, producing 342 overlaps, 27 stomps, and 12,895,847 observed tokens.
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:25:13.934Z
            confidence: 0.91
            tags:
                - cmd/clew/**
            env:
                dataset: clew dogfood D0 first watch
            affects:
                - cmd/clew/watchcmd.go
          status: suspect
        - entry:
            id: e01M05T4XG0RJWQTP25SYT4FH0B
            type: finding
            title: 'Task 2 not passable: `spent` conflates extraction, differ, and archaeology'
            body: The dogfood audit judged Task 2 not passable yet. The budget `spent` counter mixes extraction, differ, and archaeology tokens, so the predeclared extraction-only cost metric cannot be read from it. Separating cost by kind is a prerequisite before the Task 2 gate can be honestly evaluated.
            quote: '`spent` combines extraction, differ, and archaeology; it is not extraction-only.'
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:33:11.552Z
            confidence: 0.9
            tags:
                - internal/**
                - cmd/clew/**
            env:
                host: local dogfood machine
                dataset: $CLEW_HOME/state.db budget table
            affects:
                - internal/**
          status: suspect
        - entry:
            id: e01M05T4XG0RJWQTP25T1K58B61
            type: finding
            title: Confirm/reject only in event YAML; adapter unknowns undated, absent from status
            body: Human confirm/reject signals are recorded only in per-worktree events/*.yaml, so measuring confirm rate requires a find+awk scrape instead of a DB query. Adapter "unknown" counts are cumulative, undated KV rows and never surfaced in status. Both dogfood metrics are therefore not queryable from state.db.
            quote: Confirm/reject exists only in journal event YAML.
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:33:11.552Z
            confidence: 0.87
            tags:
                - internal/**
                - cmd/clew/**
            env:
                host: local dogfood machine
                dataset: $CLEW_HOME/state.db + worktrees events/*.yaml
            affects:
                - internal/**
                - cmd/clew/**
          status: suspect
        - entry:
            id: e01M05TX51R4VZHJMNM804YPMRT
            type: finding
            title: 'D1: 30 automatic entries; Codex metadata incident pinned'
            body: 'D1 live dogfood: 30 session entries appeared from 1 real Codex session with 0 manual notes; observed=5549571, live+backfill extraction=30184, all-LLM=39091, pushes delivered=0/0, open alerts=10; 46 records in 3 newly observed multi-agent metadata classes were pinned as non-utterance adapter metadata.'
            quote: 'D1 live dogfood: 30 session entries appeared from 1 real Codex session with 0 manual notes; observed=5549571, live+backfill extraction=30184, all-LLM=39091, pushes delivered=0/0, open alerts=10; 46 records in 3 newly observed multi-agent metadata classes were pinned as non-utterance adapter metadata.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:46:25.720035Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M05V9MQWYX3BAX0VXZ70SHTD
            type: finding
            title: 'D2: cursor rewind replayed 58,754 bytes once'
            body: 'D2 migration failure: split cursor rewind replayed 58754 bytes and spent 1815 extraction tokens once; delivered pushes=0. Fix is monotonic max(extract, watch-extract), with divergent-cursor regression. State backup: state.db.d1-20260816T1748Z.bak.'
            quote: 'D2 migration failure: split cursor rewind replayed 58754 bytes and spent 1815 extraction tokens once; delivered pushes=0. Fix is monotonic max(extract, watch-extract), with divergent-cursor regression. State backup: state.db.d1-20260816T1748Z.bak.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:53:15.004943Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0640H4ZAR0BQMS73R8QW9E7
            type: finding
            title: Repaired watcher installed as launchd agent dev.clew.watch
            body: After journaling the D2 cursor-rewind finding, the fixed watcher was installed as a launchd agent named dev.clew.watch on the dev Mac, writing to /Users/mac/.clew/logs/watch.log. That log path is where watcher behaviour for subsequent live runs can be inspected.
            quote: 'installed launchd agent dev.clew.watch (log: /Users/mac/.clew/logs/watch.log)'
            utterance_by: user
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1508
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:25:33.599Z
            confidence: 0.86
            tags:
                - cmd/clew/**
            env:
                host: local macOS dev machine (launchd)
            affects:
                - /Users/mac/.clew/logs/watch.log
          status: current
        - entry:
            id: e01M064K88F7SDMHA1SPAB51HK7
            type: finding
            title: 'D2 final: 52 automatic entries; live extraction 0.631%'
            body: 'D2-final: repos=3; automatic-session-entries=52; observed=6779248; live-extraction=42803 (0.631%); backfill=5057; all-LLM=67936/200000; C:R=0:1; pushes=0 delivered/0 unneeded (precision=N/A); adapter/system incidents=4; parked=0; active-reservations=0; live-sessions=6.'
            quote: 'D2-final: repos=3; automatic-session-entries=52; observed=6779248; live-extraction=42803 (0.631%); backfill=5057; all-LLM=67936/200000; C:R=0:1; pushes=0 delivered/0 unneeded (precision=N/A); adapter/system incidents=4; parked=0; active-reservations=0; live-sessions=6.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:35:47.087963Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M064YRS4S9NK7KW9NN351FQF
            type: finding
            title: 'I9: Claude settlement ignores cache token fields, letting spend exceed caps'
            body: Settlement of Claude LLM calls counts only the non-cache token fields, ignoring `cache_creation_input_tokens` and `cache_read_input_tokens`. Cumulative spend is therefore undercounted and can run past the configured budget caps. Found at llm.go:158 during the read-only gate.
            quote: 'I9: Claude settlement ignores `cache_creation_input_tokens` and `cache_read_input_tokens`, permitting cumulative spend beyond caps.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:42:04.452Z
            confidence: 0.93
            tags:
                - internal/llm/**
                - internal/state/**
            env: null
            affects:
                - internal/llm/llm.go
          status: suspect
        - entry:
            id: e01M064YRS4S9NK7KW9NQ1JMDV2
            type: finding
            title: Malformed or missing pinned timestamps silently fall back to ingest time
            body: When a source record's pinned timestamp is missing or malformed, the adapter/extract path substitutes the ingest-time `now` without signalling, so entries get fabricated source times. Located at adapters.go:151 and extract.go:264; flagged as a gate blocker.
            quote: 'Source time: malformed/missing pinned timestamps silently become ingest `now`.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:42:04.452Z
            confidence: 0.92
            tags:
                - internal/adapters/**
                - internal/extract/**
            env: null
            affects:
                - internal/adapters/adapters.go
                - internal/extract/extract.go
          status: suspect
        - entry:
            id: e01M065G8RTKVH6466KE7GQREGX
            type: finding
            title: Task 2 passes its live gate on all five acceptance checks
            body: 'A live gate run on Task 2 passed: 52 automatic session entries recorded, 0 delivered-but-unneeded pushes, cursors monotonic, the exact installed binary in use, and no active adapter or LLM errors. This is the verdict that unblocked committing the gate fixes.'
            quote: 'Task 2 now passes its live gate: 52 automatic session entries, 0 delivered/unneeded pushes, monotonic cursors, exact installed binary, and no active adapter/LLM errors.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:51:37.882Z
            confidence: 0.93
            tags: []
            env:
                dataset: live session watch
            affects: []
          status: current
        - entry:
            id: e01M0660RW03A1VSC6ENMH9M2J7
            type: finding
            title: Parallel agent task killed by gpt-5.6-sol TPM rate limit
            body: 'The task2_final agent errored with "stream disconnected before completion" due to an OpenAI org-level tokens-per-minute limit for gpt-5.6-sol: limit 500000, used 397298, requested 204717. Running several agents concurrently can exceed the TPM ceiling and drop in-flight work.'
            quote: 'stream disconnected before completion: Rate limit reached for gpt-5.6-sol in organization ‹redacted› on tokens per min (TPM): Limit 500000, Used 397298, Requested 204717.'
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L160
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:00:38.656Z
            confidence: 0.85
            tags: []
            env:
                host: openai api
                dataset: gpt-5.6-sol
            affects:
                - task2_final
          status: current
        - entry:
            id: e01M066CV9W2WX7FZA44QYME204
            type: finding
            title: 'Alert self-cleaning shipped: WithdrawWhen, ReconcileAlerts, six kinds withdrawn'
            body: 'Landed: Alert.WithdrawWhen with a legacy DB migration, ReconcileAlerts(repo, kinds, active), and differ withdrawal of stale contradiction, absence, aging, suspect, stomp, and overlap alerts on every poll. Resolved-stomp and status-resolution tests added; go test ./..., race tests, and vet all pass.'
            quote: '- Differ now withdraws stale contradiction, absence, aging, suspect, stomp, and overlap alerts each poll.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L263
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:07:14.364Z
            confidence: 0.88
            tags:
                - state/**
                - differ/**
            env: null
            affects:
                - state/**
                - differ/**
          status: current
        - entry:
            id: e01M066XWD8F1SAWQ8HVXGW4J4Z
            type: finding
            title: 'Task 3 docket gate: 8→1, FYI 0, withdrawal 1 poll'
            body: cards=8; render=1 overflow-failure; cap=7; synthetic-FYI-rendered=0; resolved-stomp-withdrawal=1 poll; pushes=0/0; full=pass; race=pass; vet=pass.
            quote: cards=8; render=1 overflow-failure; cap=7; synthetic-FYI-rendered=0; resolved-stomp-withdrawal=1 poll; pushes=0/0; full=pass; race=pass; vet=pass.
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:16:32.552542Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M067AZ3VJM73XP82REZEZ1QC
            type: finding
            title: 'Task 4 channel: ntfy 2383c10c…, card creation only'
            body: ntfy-topic=https://ntfy.sh/2383c10ce6438813da9969532f2df2f7; push-trigger=docket-card-creation-only; payload=headline+why-you; html-refresh=30s.
            quote: ntfy-topic=https://ntfy.sh/2383c10ce6438813da9969532f2df2f7; push-trigger=docket-card-creation-only; payload=headline+why-you; html-refresh=30s.
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:23:41.307799Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M067G2QY8BZNZXBE46QBEXEW
            type: finding
            title: 'Task 4 gate: 10ms, HTML 30s, ntfy 5/5'
            body: bare-clew=10ms x5; dashboard-sections=5; html-refresh=30s; title-light=nonempty-only; ntfy-delivered=5; payload-valid=5/5; full=pass; race=pass; vet=pass.
            quote: bare-clew=10ms x5; dashboard-sections=5; html-refresh=30s; title-light=nonempty-only; ntfy-delivered=5; payload-valid=5/5; full=pass; race=pass; vet=pass.
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:26:28.862743Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M06821D53QYHBJS1FEC2CK7G
            type: finding
            title: 'Task 5 gate: 3 formats, 1 card, confirm boundary pass'
            body: formats=bundle+dir+https; schema-invalid=reject; quote-missing=reject; batch=1 card; live-stage=1 entry; live-open=pass; live-reject=0 journal writes; accept=1 foreign entry+1 human confirm; branch-push=pass x2 idempotent; full=pass; race=pass; vet=pass.
            quote: formats=bundle+dir+https; schema-invalid=reject; quote-missing=reject; batch=1 card; live-stage=1 entry; live-open=pass; live-reject=0 journal writes; accept=1 foreign entry+1 human confirm; branch-push=pass x2 idempotent; full=pass; race=pass; vet=pass.
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:36:17.317827Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M068FQH1ND9MM1WH851AF45M
            type: finding
            title: 'Task 6 gate: flags 0 writes; algebra, poller, manifest pass'
            body: note-help entries=69→69; absence-threshold=4 proposed/5 absent; ineligible=proposed; human-confirm=absent; contradiction nonhuman=possible/human=contradicted; env different=current/current, same=superseded/current; poller best-overlap=pass/no-overlap=none/out-of-window=none; manifest rerun events=2→2; full=pass; race=pass; vet=pass.
            quote: note-help entries=69→69; absence-threshold=4 proposed/5 absent; ineligible=proposed; human-confirm=absent; contradiction nonhuman=possible/human=contradicted; env different=current/current, same=superseded/current; poller best-overlap=pass/no-overlap=none/out-of-window=none; manifest rerun events=2→2; full=pass; race=pass; vet=pass.
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:43:45.953155Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M069MQYJX6QVW3YCWWTAWV34
            type: finding
            title: --help
            body: --help
            quote: --help
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: cursor
                at: 2026-08-16T22:03:58.802665362Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGK6W0Y1NEDQW45KDECS
            type: finding
            title: 'Corpse census: substrate died in a 6-day burst; tombstone came 5 weeks late'
            body: 'substrate: 63/64 commits in week one (Jun 9-14), five weeks silence, final commit is the failure confession (LIFECYCLE.md + README ''failed adoption''). The promised compounding loop (scheduler/repair/steward/federated store) was never built - confessed by its own docs. Zero tags, zero CI, zero adopters.'
            quote: document adoption failure and harden project creation
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGK9TDRH9RHAJG9YDPME
            type: finding
            title: 'Census: security-substrate human-sealed in 1 day; stratura inherited nothing'
            body: 'security-substrate: born the day after substrate''s tombstone, formally SEALED as failed on day 2 with STOP packet and constitution - faster than any detector. stratura: zero references to either predecessor (measured), then repeated the sealed pathology (safety perimeter before usable core). Detection was never the bottleneck; INHERITANCE was.'
            quote: preserved as future MOSS input
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGKB9VP38VBTQ7YC9GBR
            type: finding
            title: 'Regime detector: composition and earned-state separate corpses from control'
            body: 'Cadence cannot distinguish clew from the corpses (all are burst-projects at day scale). What separates: core-touch ratio (clew 50% vs 0-22%), earned state (passing gates, live dogfood, metrics vs zero), and inheritance (clew is the first generation that carried anything). n=4; the ~100-repo lineage census remains to be run.'
            quote: my GitHub probably has a hundred or more examples
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 0.95
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGKDZ8J6SQQK3BXKYAV5
            type: finding
            title: module clew blocks go install by URL; release binaries or rename needed
            body: go.mod declares 'module clew', not the repo path, so go install github.com/maceip/clew/... fails. Env recipes need git clone && go build until the module is renamed or release binaries ship. Blocks the one-line cloud env install.
            quote: go 1.26.3
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags:
                - go.mod
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGKF09RFDGJD90VNADCQ
            type: finding
            title: 'Spawn test: 63/63 entries carried into scratch gen-2 with guardrails intact'
            body: 'init --carry into a fresh repo: full seed landed, carried provenance preserved, newborn glance renders the constitution, context.md opens with the 6.5 injection preamble before any agent typed. Cross-machine multi-hop proven: laptop decisions -> branch -> cloud VM -> manifest -> gen-2. Differ re-flagged design-era contradictions in the newborn.'
            quote: carried 63, dropped 0 (dispositions journaled - the loss is deliberate and dated)
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGKJN308G4HVJTE4A4AB
            type: finding
            title: 'note-then-edit has a limbo: placeholders auto-commit and leaked into a seed'
            body: 'journal note commits placeholder text immediately; later file rewrites sit uncommitted, so a manifest exported placeholders into gen-2. Fix direction: clew journal add <file> for validated whole-entry ingestion (also needed by cloud self- extraction). Until then notes must be final text.'
            quote: phone surface intent placeholder
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASQGKMHGY4F7JBEQX1ZT3T
            type: finding
            title: 'This chat became the supernode: 788 messages of unjournaled load-bearing context'
            body: 'The clew-design session ran 3 days on an uncovered surface; the owner became afraid to close it - the exact single-point-of-failure clew abolishes. Exit kit: full raw transcript attached at transcripts/ on this branch; distilled decisions/findings/questions journaled; resumption works from any surface via branch fetch.'
            quote: i am afraid to leave this session with you because you have all the context (the exact thing clew is supposed to help with, funnily enough)
            utterance_by: user
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T15:58:00Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0ASSNH1HP68M1QERV9AKG5A
            type: finding
            title: Attachments bypass the secret scrub; GitHub push protection caught PATs
            body: 'GitHub push protection blocked the journal push: two ephemeral PATs the owner pasted in chat were present verbatim in the attached raw transcript. The entries pipeline scrubs quotes/bodies (6.2) but attachments bypass the scrub entirely. Fix: run the same secret-scrub over transcripts/ (and any attachment) before commit; treat platform push-protection as backstop, never primary.'
            quote: GITHUB PUSH PROTECTION - Push cannot contain secrets
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T16:02:00Z
            confidence: 1
            tags:
                - internal/scrub/**
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0AXCM54CDSPV84DBH1PWGWD
            type: finding
            title: 'Spec nudge matrix is stale: codex and gemini now ship injection hooks'
            body: 'Aug 2026 survey: codex hooks are stable and default-enabled with UserPromptSubmit additionalContext; gemini CLI BeforeAgent injects context (default on v0.26+); cursor injects only at sessionStart/postToolUse, not beforeSubmitPrompt; opencode plugins transform system/messages pre-dispatch. MCP 2026-07-28 subscriptions notify the client, not the model. Re-pin JOURNAL_SPEC 8.1.'
            quote: Plain text on `stdout` is added as extra developer context.
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:03:00Z
            confidence: 0.9
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0AXCM55N0QM9RCRYF48TQ6C
            type: finding
            title: 'Universal injection point: every model API call rebuilds the mind'
            body: 'No context persists inside a model between calls — each harness reconstructs the full message list per request. A local base-URL shim (OLLAMA_HOST / OPENAI_BASE_URL / ANTHROPIC_BASE_URL) can inject the journal delta into any agent, bare ollama included, with passthrough-on-failure so it is never load-bearing. Prior art: Engram transparent ollama proxy; LiteLLM async_pre_call_hook.'
            quote: It intercepts every chat request, injects relevant memories into the system prompt, and forwards the enriched request to Ollama — all without any changes to your client code.
            utterance_by: tool_result
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: cursor-cloud-agent
                surface: cursor-cloud-vm
                at: 2026-08-18T17:04:00Z
            confidence: 0.9
            tags: []
            env: null
            affects: []
          status: current
        - entry:
            id: e01M0AYE066QK08QK8MTPX4XNFX
            type: finding
            title: 'Codex finished I13 stale: tree uncommitted, law wording on human surfaces'
            body: 'Manual check-in of the first stale finish: I13 complete and tests green but uncommitted — it exists only in the laptop working tree, invisible to the join. Confirmed conflict: owner-law vocabulary on human surfaces (README, cards, listings) vs the two-register ruling; the feature stands, only surface wording renames. Aligned: single-lineage from, SessionStart birth. Reconcile due at next contact.'
            quote: The main-branch implementation remains uncommitted in the working tree; I did not commit or push the code changes.
            utterance_by: assistant
            source:
                kind: session
                ref: chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-18T17:22:00Z
            confidence: 0.95
            tags: []
            env: null
            affects: []
          status: current
    graveyard:
        - entry:
            id: e01M04H0KBY8QQWXPE8DP99N012
            type: decision
            title: Name the system restart
            body: the app name is restart; the spec placeholder stratura was renamed per spec section 12.1
            quote: the app name is restart; the spec placeholder stratura was renamed per spec section 12.1
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T05:34:18.4941Z
            confidence: 1
            tags:
                - cmd/**
            env: null
            affects: []
          status: superseded
          events:
            - id: v01M04WCPA7K7CJTN42BAP273CH
              kind: supersede
              entry: e01M04H0KBY8QQWXPE8DP99N012
              payload:
                by: e01M04WCGJS9FS7FQB0YFX9DTYG
              by:
                who: human
                surface: macs-MacBook-Pro
              at: 2026-08-16T08:53:09.063873Z
        - entry:
            id: e01M04XS4RN9Z26XEE0CNHX8AT5
            type: finding
            title: --help
            body: --help
            quote: --help
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: cursor
                at: 2026-08-16T09:17:25.653480808Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: superseded
          events:
            - id: v01M04XT9DXQESZQ8YM4ZA4FNQE
              kind: reject
              entry: e01M04XS4RN9Z26XEE0CNHX8AT5
              payload:
                reason: ""
              by:
                who: human
                surface: cursor
              at: 2026-08-16T09:18:03.197113285Z
        - entry:
            id: e01M05RPB2N6QDXN1MP31SB64B2
            type: finding
            title: Dogfood D0 historical-session storm is an I12 failure
            body: 'Dogfood failure D0: historical sessions misclassified live=33; observed tokens=12895847; overlaps=342; stomps=27; actual pushes=0; false pushed_at=27; extraction spend=0; adapter incidents=1; watcher stopped before extraction.'
            quote: 'Dogfood failure D0: historical sessions misclassified live=33; observed tokens=12895847; overlaps=342; stomps=27; actual pushes=0; false pushed_at=27; extraction spend=0; adapter incidents=1; watcher stopped before extraction.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:07:45.365367Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: superseded
          events:
            - id: v01M05SQFQX6EE5DJTG141CRERF
              kind: supersede
              entry: e01M05RPB2N6QDXN1MP31SB64B2
              payload:
                by: e01M05SPB2EMMC4F4PR0BDQA8S5
              by:
                who: extractor
                surface: macs-MacBook-Pro
              at: 2026-08-16T17:25:51.48576Z
        - entry:
            id: e01M05S9TK3X86C59JCJDWZ7Z6R
            type: finding
            title: Dogfood D0 metrics after repair
            body: 'Dogfood 2026-08-16: repos=3; automatic entries=6; extraction=5057/4491713=0.113%; confirm:reject=0:1; actual pushes=0/0; false pushes repaired=27; adapter incidents=1; parked=0; clean cold-start sessions=0,alerts=0; live sessions after append=1.'
            quote: 'Dogfood 2026-08-16: repos=3; automatic entries=6; extraction=5057/4491713=0.113%; confirm:reject=0:1; actual pushes=0/0; false pushes repaired=27; adapter incidents=1; parked=0; clean cold-start sessions=0,alerts=0; live sessions after append=1.'
            utterance_by: user
            source:
                kind: human
                ref: cli:note
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:18:23.843973Z
            confidence: 1
            tags: []
            env: null
            affects: []
          status: superseded
          events:
            - id: v01M05SD9PYMN6XEVRG3T1S4WE5
              kind: supersede
              entry: e01M05S9TK3X86C59JCJDWZ7Z6R
              payload:
                by: e01M05S9SFKAAM813AR1AY96QWH
              by:
                who: extractor
                surface: macs-MacBook-Pro
              at: 2026-08-16T17:20:17.630392Z
        - entry:
            id: e01M05SPB2EMMC4F4PR0928NA31
            type: finding
            title: 'Task 2 dogfood audit verdict: not passable; four metrics gaps'
            body: 'An audit of the Task 2 dogfood metrics concluded they are not audit-grade yet. Gaps: `spent` conflates extraction, differ, and archaeology instead of being extraction-only; confirm/reject counts exist only inside journal event YAML; adapter unknowns are cumulative, undated KV counters and are absent from status; push precision is unmeasurable. The audit changed no files.'
            quote: 'Task 2 audit: not passable yet.'
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:25:13.934Z
            confidence: 0.88
            tags:
                - cmd/clew/**
                - internal/**
            env:
                dataset: clew dogfood Task 2 snapshot
            affects:
                - cmd/clew/**
          status: superseded
          events:
            - id: v01M05SSAA5K5QG6KBBW3PNMPJF
              kind: supersede
              entry: e01M05SPB2EMMC4F4PR0928NA31
              payload:
                by: e01M05SPB2EMMC4F4PR09QWMXG8
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T17:26:51.461075Z
        - entry:
            id: e01M05SPB2EMMC4F4PR09QWMXG8
            type: finding
            title: 'Push precision unmeasurable: unset push returns success, faking pushed_at'
            body: With push unconfigured, the push path returns success and HTTP errors go unchecked, so alerts get marked `pushed_at` even though nothing was delivered. Push precision therefore cannot be measured from the current data. Located in internal/push/push.go:16 and cmd/clew/watchcmd.go:281.
            quote: 'Push precision is unmeasurable: unset push returns success, causing false `pushed_at` marks'
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:25:13.934Z
            confidence: 0.92
            tags:
                - internal/push/**
                - cmd/clew/**
            env:
                dataset: clew dogfood Task 2 snapshot
            affects:
                - internal/push/push.go
                - cmd/clew/watchcmd.go
          status: superseded
          events:
            - id: v01M05TDFT5TVYCMX5XSB47NKWX
              kind: supersede
              entry: e01M05SPB2EMMC4F4PR09QWMXG8
              payload:
                by: e01M05TCM2N1P9P2EQSE8YKVQX3
              by:
                who: extractor
                surface: macs-MacBook-Pro
              at: 2026-08-16T17:37:52.453011Z
        - entry:
            id: e01M05T52RRS9EXQ92M3PA9WBVF
            type: finding
            title: Baseline/upgrade paths have state-transition races unit tests don't exercise
            body: Reviewing the revised watch fix surfaced several state-transition races in the baseline and upgrade paths that the current unit tests do not cover. Concrete failure sequences were handed to the implementing agent, and a race-suite pass was run against the revised diff to check them.
            quote: The baseline/upgrade paths exposed several state-transition races that unit tests don’t currently exercise.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L446
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:33:16.952Z
            confidence: 0.9
            tags:
                - cmd/clew/**
                - internal/**
            env: null
            affects:
                - cmd/clew/watchcmd.go
          status: superseded
          events:
            - id: v01M05TE8VA89CTM9KCE8X4HZBT
              kind: supersede
              entry: e01M05T52RRS9EXQ92M3PA9WBVF
              payload:
                by: e01M05TCM2N1P9P2EQSE8YKVQX3
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T17:38:18.090903Z
        - entry:
            id: e01M05TCM2N1P9P2EQSE8YKVQX3
            type: finding
            title: 'Push failure paths unchecked: unset endpoint returns success, HTTP errors ignor…'
            body: 'Audit of the push path found two defects that make push precision unmeasurable: an unset push endpoint returns success (so alerts get a false `pushed_at` mark), and HTTP errors from the push call are not checked either. Located at internal/push/push.go:16 and cmd/clew/watchcmd.go:281.'
            quote: HTTP errors are also unchecked.
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1075
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:37:24.053Z
            confidence: 0.87
            tags:
                - internal/push/**
                - cmd/clew/**
            env: null
            affects:
                - internal/push/push.go
                - cmd/clew/watchcmd.go
          status: superseded
          events:
            - id: v01M05VZ9GSVSBHHQGNTZ5974HQ
              kind: supersede
              entry: e01M05TCM2N1P9P2EQSE8YKVQX3
              payload:
                by: e01M05VTCM3AR0WFY9TZH6DJYDA
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
            - id: v01M05VZ9GSVSBHHQGNTZVGJCV0
              kind: supersede
              entry: e01M05TCM2N1P9P2EQSE8YKVQX3
              payload:
                by: e01M05VTCM3AR0WFY9TZPG9W1J8
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
        - entry:
            id: e01M05VFAW9A783PMZZEBHJ1TZH
            type: finding
            title: State-package tests pass, including 20-way concurrent reservation races
            body: After the transactional reservation/settlement work in internal/state, the package's tests pass, covering 20-way concurrent reservation races against both limits (cap and ratio). Rollover, double-settlement, and migration behavior had not yet been verified at this point, and the wider suite had not been run.
            quote: The state-package tests now pass, including 20-way concurrent reservation races against both limits.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L268
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:56:21.513Z
            confidence: 0.88
            tags:
                - internal/state/**
            env: null
            affects:
                - internal/state/**
          status: superseded
          events:
            - id: v01M05VT6X9V07HN2WWXJ60H2D3
              kind: supersede
              entry: e01M05VFAW9A783PMZZEBHJ1TZH
              payload:
                by: e01M05VNN1T5TKJ15Q47KSS3FMJ
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:02:17.897377Z
        - entry:
            id: e01M05VNN1T5TKJ15Q47KSS3FMJ
            type: finding
            title: Cap/ratio admission serialized with SQLite BEGIN IMMEDIATE; typed reservation e…
            body: 'The LLM reservation/settlement accounting in internal/state landed as atomic operations: concurrent daily-cap and extraction-ratio admission checks are serialized using SQLite `BEGIN IMMEDIATE` transactions, and the API surfaces typed errors for limit, overrun, and duplicate settlement. Implementation stayed inside internal/state with no caller or spec edits.'
            quote: Concurrent daily-cap and extraction-ratio admission is serialized via SQLite `BEGIN IMMEDIATE`.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L319
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T17:59:48.538Z
            confidence: 0.91
            tags:
                - internal/state/**
            env: null
            affects:
                - internal/state/**
          status: superseded
          events:
            - id: v01M05VZ9GSVSBHHQGNV3PS9STB
              kind: supersede
              entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
              payload:
                by: e01M05VTCM3AR0WFY9TZH6DJYDA
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
            - id: v01M05VZ9GSVSBHHQGNV3WRTWNG
              kind: supersede
              entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
              payload:
                by: e01M05VTCM3AR0WFY9TZJXAN7SE
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
            - id: v01M05VZ9GSVSBHHQGNV4TKZWJA
              kind: supersede
              entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
              payload:
                by: e01M05VTCM3AR0WFY9TZPG9W1J8
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
        - entry:
            id: e01M05VTCM3AR0WFY9TZH6DJYDA
            type: finding
            title: I9 reservation API is dead code; callers still use racy estimate gates
            body: 'The transactional reservation/settlement work landed in internal/state, but no caller uses it: watch, backfill, and extract still gate on pre-flight estimates and call RecordSpend afterward, so the race remains. Retry spend can also be dropped on error. Cited sites: watchcmd.go:411-431, backfillcmd.go:111-121, extract.go:139-159.'
            quote: I9 reservation code is unused; all callers still use racy estimate gates and `RecordSpend` afterward.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L961
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T18:02:23.747Z
            confidence: 0.93
            tags:
                - cmd/clew/**
                - internal/state/**
            env: null
            affects:
                - cmd/clew/watchcmd.go
                - cmd/clew/backfillcmd.go
                - extract.go
                - internal/state/**
          status: superseded
          events:
            - id: v01M05VZ9GSVSBHHQGNV55DXPH4
              kind: supersede
              entry: e01M05VTCM3AR0WFY9TZH6DJYDA
              payload:
                by: e01M05VTCM3AR0WFY9TZJXAN7SE
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
            - id: v01M05VZ9GSVSBHHQGNV59XHZV1
              kind: supersede
              entry: e01M05VTCM3AR0WFY9TZH6DJYDA
              payload:
                by: e01M05VTCM3AR0WFY9TZPG9W1J8
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
        - entry:
            id: e01M05VTCM3AR0WFY9TZJXAN7SE
            type: finding
            title: Settlement records overruns before reporting them, so the hard cap is unenforce…
            body: In internal/state (state.go:517-563), settlement writes the spend and only then surfaces the overrun, meaning the budget cap can be exceeded and merely observed rather than blocked. The typed overrun error alone does not make the cap enforceable.
            quote: 'Settlement only reports overruns after recording them, so it still cannot enforce the hard cap: `state.go:517-563`.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L961
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T18:02:23.747Z
            confidence: 0.92
            tags:
                - internal/state/**
            env: null
            affects:
                - internal/state/state.go
          status: superseded
          events:
            - id: v01M05VZ9GSVSBHHQGNV7AYYDV4
              kind: supersede
              entry: e01M05VTCM3AR0WFY9TZJXAN7SE
              payload:
                by: e01M05VTCM3AR0WFY9TZPG9W1J8
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T18:05:04.409743Z
        - entry:
            id: e01M05VTCM3AR0WFY9TZKZMBMA4
            type: finding
            title: Neutral cwd breaks relative custom extractor commands like ./bin/extractor
            body: Running the LLM subprocess with a neutral working directory (llm.go:93-99) breaks the supported configuration of a relative custom command path, e.g. `./bin/extractor`, which resolves against the project directory.
            quote: 'Neutral cwd breaks supported relative custom commands such as `./bin/extractor`: `llm.go:93-99`.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L961
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T18:02:23.747Z
            confidence: 0.9
            tags:
                - internal/llm/**
            env: null
            affects:
                - llm.go
          status: superseded
          events:
            - id: v01M0653VVE690K37G2ZTSESQ12
              kind: supersede
              entry: e01M05VTCM3AR0WFY9TZKZMBMA4
              payload:
                by: e01M064YRS4S9NK7KW9NN3114KH
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
            - id: v01M0653VVE690K37G2ZV5GK70P
              kind: supersede
              entry: e01M05VTCM3AR0WFY9TZKZMBMA4
              payload:
                by: e01M064YRS4S9NK7KW9NN351FQF
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
        - entry:
            id: e01M05VTCM3AR0WFY9TZPG9W1J8
            type: finding
            title: Cursor migration and init/bootstrap races fixed; test, race, vet, diff checks p…
            body: 'The cursor-migration monotonicity fix and two init/bootstrap races are resolved in the live diff, and the full suite (`go test ./...`, race tests, vet, diff checks) passes. Remaining blockers are budget-side: unused reservation path, unenforceable cap, and the neutral-cwd regression.'
            quote: Cursor migration and init/bootstrap races are fixed. `go test ./...`, race tests, vet, and diff checks pass.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L961
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T18:02:23.747Z
            confidence: 0.9
            tags:
                - cmd/clew/**
                - internal/state/**
            env: null
            affects:
                - internal/state/**
          status: superseded
          events:
            - id: v01M0644MDBVNY4B9EFATFQXGXT
              kind: supersede
              entry: e01M05VTCM3AR0WFY9TZPG9W1J8
              payload:
                by: e01M0642VRXV9PCGA4NDHV479C9
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:27:48.011004Z
        - entry:
            id: e01M0642VRXV9PCGA4NDHV479C9
            type: finding
            title: Cursor path now monotonic and watcher stable
            body: After the migration-monotonicity work, the cursor path no longer rewinds and the watcher is reported stable. This is the status following the earlier fixes for cursor migration and the watch storm; remaining work moved on to budget accounting.
            quote: The cursor path is now monotonic and the watcher is stable.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1577
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:26:50.013Z
            confidence: 0.78
            tags:
                - internal/state/**
            env: null
            affects:
                - e01M05V0G3Q9F41V62P6T44TC04
                - e01M05V0G3Q9F41V62P6R5QV53G
          status: superseded
          events:
            - id: v01M0653VVE690K37G2ZYEWQ2RV
              kind: supersede
              entry: e01M0642VRXV9PCGA4NDHV479C9
              payload:
                by: e01M064XJ84ZTY2HNFWWV8RSRQR
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
            - id: v01M0653VVE690K37G2ZZPGM9V6
              kind: supersede
              entry: e01M0642VRXV9PCGA4NDHV479C9
              payload:
                by: e01M064YRS4S9NK7KW9NN351FQF
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
        - entry:
            id: e01M064XJ84ZTY2HNFWWQZATFA9
            type: finding
            title: 'Reservation API is dead code: callers still use racy estimate gates + RecordSpe…'
            body: The new atomic I9 reservation code in internal/state has no callers. The watch, backfill, and extract paths still gate on a pre-call estimate and call `RecordSpend` after the fact, and retry spend can be dropped on error. Cited at watchcmd.go:411-431, backfillcmd.go:111-121, extract.go:139-159.
            quote: I9 reservation code is unused; all callers still use racy estimate gates and `RecordSpend` afterward.
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L245
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:41:24.996Z
            confidence: 0.91
            tags:
                - internal/extract/**
                - '**/watchcmd.go'
                - '**/backfillcmd.go'
            env: null
            affects:
                - internal/extract/extract.go
                - e01M0642VRXV9PCGA4NDJF92E2Y
          status: superseded
          events:
            - id: v01M0653VVE690K37G3019M4A1T
              kind: supersede
              entry: e01M064XJ84ZTY2HNFWWQZATFA9
              payload:
                by: e01M064YRS4S9NK7KW9NN3114KH
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
            - id: v01M0653VVE690K37G303870E38
              kind: supersede
              entry: e01M064XJ84ZTY2HNFWWQZATFA9
              payload:
                by: e01M064YRS4S9NK7KW9NQ1JMDV2
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
        - entry:
            id: e01M064XJ84ZTY2HNFWWV8RSRQR
            type: finding
            title: Settlement records overruns before reporting them, so the hard cap is unenforce…
            body: The settlement routine writes the spend and only afterwards reports that it exceeded the limit, meaning the cap can be breached rather than blocked. Reviewer marked this blocking at state.go:517-563, alongside the note that cursor migration and init/bootstrap races are now fixed.
            quote: 'Settlement only reports overruns after recording them, so it still cannot enforce the hard cap: `state.go:517-563`.'
            utterance_by: tool_result
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L245
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:41:24.996Z
            confidence: 0.9
            tags:
                - internal/state/**
            env: null
            affects:
                - internal/state/state.go
          status: superseded
          events:
            - id: v01M0653VVE690K37G304QMZR7P
              kind: supersede
              entry: e01M064XJ84ZTY2HNFWWV8RSRQR
              payload:
                by: e01M064YRS4S9NK7KW9NN351FQF
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
        - entry:
            id: e01M064YRS4S9NK7KW9NN3114KH
            type: finding
            title: 'I2: any JSON object counts as extraction success and advances the cursor'
            body: The extraction path treats any well-formed JSON object as a successful result, so an empty `{}` or a changed Claude response envelope silently advances the session cursor instead of parking loudly. Identified during the strict read-only gate as a blocking defect at extract.go:149 and llm.go:165.
            quote: 'I2: any JSON object is accepted as extraction success; `{}` or a changed Claude envelope advances the cursor instead of parking loudly.'
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:42:04.452Z
            confidence: 0.93
            tags:
                - internal/extract/**
                - internal/llm/**
            env: null
            affects:
                - internal/extract/extract.go
                - internal/llm/llm.go
          status: superseded
          events:
            - id: v01M0653VVE690K37G305GVYZV2
              kind: supersede
              entry: e01M064YRS4S9NK7KW9NN3114KH
              payload:
                by: e01M064YRS4S9NK7KW9NN351FQF
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
            - id: v01M0653VVE690K37G309DCEANA
              kind: supersede
              entry: e01M064YRS4S9NK7KW9NN3114KH
              payload:
                by: e01M064YRS4S9NK7KW9NQ1JMDV2
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-16T20:44:51.438665Z
        - entry:
            id: e01M065G8RTKVH6466KEB0FRJ8N
            type: intent
            title: Commit the gate fixes as one spec-amended change, then start the docket
            body: After the Task 2 live gate passed, the plan is to land all gate fixes as a single commit that also amends the spec, and then move on to the docket surface work.
            quote: I’m committing the gate fixes as one spec-amended change, then moving to the docket.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:51:37.882Z
            confidence: 0.9
            tags: []
            env: null
            affects: []
          status: absent
        - entry:
            id: e01M065SK8W1ZT32KZF92CP8KRT
            type: finding
            title: Alerts only inserted; nothing closed them and keys used mutable prose
            body: Before the reconcile work, the implementation had no poll path that closed alerts, so open alerts accumulated indefinitely, and alert keys were built from mutable prose — making identity unstable across polls.
            quote: The current implementation only inserts alerts; no poll ever closes them, and keys include mutable prose.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L45
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T20:56:43.548Z
            confidence: 0.9
            tags:
                - state/**
                - differ/**
            env: null
            affects:
                - state/**
                - differ/**
          status: superseded
          events:
            - id: v01M0ASVJ3AZSG7D8MX1D4SPJAR
              kind: supersede
              entry: e01M065SK8W1ZT32KZF92CP8KRT
              payload:
                by: e01M06656MKNEEY96GBNNDYYR36
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-18T16:04:17.130585Z
            - id: v01M0ASVJ3AZSG7D8MX1G8TPTRZ
              kind: supersede
              entry: e01M065SK8W1ZT32KZF92CP8KRT
              payload:
                by: e01M066CV9W2WX7FZA44QYME204
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-18T16:04:17.130585Z
        - entry:
            id: e01M06656MKNEEY96GBNNDYYR36
            type: finding
            title: Stomp withdrawal verified on dirty-path and session-expiry in the next Run
            body: 'Focused state and differ tests passed, covering both the dirty-path and session-expiry stomp cases: the stale alert is withdrawn on the very next poll cycle rather than lingering. Full suite and shared-worktree integration checks followed.'
            quote: Focused state/differ tests pass, including both dirty-path and session-expiry stomp withdrawal in the very next `Run`.
            utterance_by: assistant
            source:
                kind: session
                ref: codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L194
                agent: codex
                surface: macs-MacBook-Pro
                at: 2026-08-16T21:03:03.827Z
            confidence: 0.9
            tags:
                - state/**
                - differ/**
            env: null
            affects:
                - state/**
                - differ/**
          status: superseded
          events:
            - id: v01M0ASVJ3AZSG7D8MX1GDY68A2
              kind: supersede
              entry: e01M06656MKNEEY96GBNNDYYR36
              payload:
                by: e01M066CV9W2WX7FZA44QYME204
                why: newer measurement, same env
              by:
                who: differ
                surface: macs-MacBook-Pro
              at: 2026-08-18T16:04:17.130585Z
    exhibits:
        - id: v01M04W6B48BX7A75HCZX5KDVEX
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Implement the journal system per JOURNAL_SPEC v1
            ref: 3ec6cb91176446ac4f216e5a64c2daa72821717e
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T08:49:41Z
        - id: v01M04WEHTGWW1X69CNKG9TKS0F
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: 'Rename product: restart → clew (owner decision; JOURNAL_SPEC §12.1 closed)'
            ref: 26fe9c230acc7056e1a88346129870dd1cba82d2
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T08:54:10Z
        - id: v01M05RFD1GC8N5VXNY49DY4VP8
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Fix generated journal title clipping
            ref: de4282ea478934e6435dceea690c927cafaa47be
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T17:03:58Z
        - id: v01M05SQFQX6EE5DJTG14C12QJG
          kind: evidence
          entry: e01M05SG3NTP5W2JX7Y6MG00HQJ
          payload:
            kind: session
            ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335
          by:
            who: extractor
            surface: macs-MacBook-Pro
          at: 2026-08-16T17:25:51.48576Z
        - id: v01M05SXM6B43DFB7JFTJMHBAQW
          kind: evidence
          entry: e01M05SG3NTP5W2JX7Y6MG00HQJ
          payload:
            kind: session
            ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L934
          by:
            who: extractor
            surface: macs-MacBook-Pro
          at: 2026-08-16T17:29:12.651541Z
        - id: v01M05TC7V3N8F0E55N3P20HECD
          kind: evidence
          entry: e01M05SVGK1Q2MR34Y3CMHR7DXM
          payload:
            kind: session
            ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041
          by:
            who: extractor
            surface: macs-MacBook-Pro
          at: 2026-08-16T17:37:11.523232Z
        - id: v01M064DVAVCVBMK76Z5QXTC778
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZKZMBMA4
          payload:
            kind: session
            ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L24
          by:
            who: extractor
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:32:50.011227Z
        - id: v01M065G9VR22YS55XF05Y7NS7C
          kind: evidence
          entry: e01M05V5HWA6TFT0A0KDZY8S45K
          payload:
            kind: commit
            note: Make zero-homework watch safe and measurable
            ref: b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF06ASCXHV
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Make zero-homework watch safe and measurable
            ref: b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF08JR8J8J
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Make zero-homework watch safe and measurable
            ref: b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0C89YR0F
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Make zero-homework watch safe and measurable
            ref: b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0KRWW4WN
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWQZATFA9
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0QJ8JSJ5
          kind: evidence
          entry: e01M05T52RRS9EXQ92M3PA9WBVF
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0RZ7JZNY
          kind: evidence
          entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0SXDMH1B
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NN3114KH
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0VNM97P8
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZKZMBMA4
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0XDGXPQV
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NN351FQF
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0XWP02N4
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZH6DJYDA
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF0ZH82WM9
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZPG9W1J8
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF126BTMET
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NQ1JMDV2
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF165TS28J
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF19ESRP71
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1CY24HTF
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0BDQA8S5
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1G706X3B
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWV8RSRQR
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1H7W3PDP
          kind: evidence
          entry: e01M05VFAW9A783PMZZEBHJ1TZH
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1N08WC3B
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZJXAN7SE
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1R53Y5CJ
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1T2DPDNG
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR09QWMXG8
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065G9VR22YS55XF1VABNACQ
          kind: evidence
          entry: e01M05TCM2N1P9P2EQSE8YKVQX3
          payload:
            kind: churn
            note: Make zero-homework watch safe and measurable
            ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:51:39Z
        - id: v01M065GZB8C2SS9H9S5BP8EC30
          kind: evidence
          entry: e01M065G8RTKVH6466KEB0FRJ8N
          payload:
            confidence: 0.85
            kind: commit
            note: 'journal: 3 file(s)'
            ref: 1be8c5020483730f4effd68b53f599bdf3db33c5
            via: link-pass
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T20:52:01Z
        - id: v01M066XZS0BVBJQSSZ17NZTTSR
          kind: evidence
          entry: e01M05V5HWA6TFT0A0KDZY8S45K
          payload:
            kind: commit
            note: Replace inbox with a bounded self-cleaning docket
            ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ1AGPD2A0
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Replace inbox with a bounded self-cleaning docket
            ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ1HCM3885
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Replace inbox with a bounded self-cleaning docket
            ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ1K3WA4W9
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Replace inbox with a bounded self-cleaning docket
            ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ1S581QX6
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR09QWMXG8
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ1V6DVERB
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWV8RSRQR
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ1YF4SXXX
          kind: evidence
          entry: e01M05TCM2N1P9P2EQSE8YKVQX3
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ20MBGE1Z
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ21B27KX0
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZH6DJYDA
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ22CNBTJD
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZPG9W1J8
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ22DBXSXV
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ2381J7VK
          kind: evidence
          entry: e01M05VFAW9A783PMZZEBHJ1TZH
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ276QFMFV
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZJXAN7SE
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ286Y2N4Q
          kind: evidence
          entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M066XZS0BVBJQSSZ2BTB26YB
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Replace inbox with a bounded self-cleaning docket
            ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:16:36Z
        - id: v01M067GFJGKHWY2CPHZ5S24TM0
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Add calm glance tiers and card-only push
            ref: a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZBR0CAD7
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Add calm glance tiers and card-only push
            ref: a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZFE3K07M
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Add calm glance tiers and card-only push
            ref: a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZJ44YF8N
          kind: evidence
          entry: e01M05T52RRS9EXQ92M3PA9WBVF
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZJFCQ8JR
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR09QWMXG8
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZMYS6KZ9
          kind: evidence
          entry: e01M05TCM2N1P9P2EQSE8YKVQX3
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZPRSWTHN
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZH6DJYDA
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZSSTFBVY
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZTY3TBCS
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0BDQA8S5
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZX049R86
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M067GFJGKHWY2CPHZXV43ZEQ
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Add calm glance tiers and card-only push
            ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:26:42Z
        - id: v01M068329RM9V5WNBWXP03ZJW5
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Add validated foreign proposal imports
            ref: 36350f8fa57bb8b838533665811bd1b42d93a71b
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:36:51Z
        - id: v01M068329RM9V5WNBWXPRJDB45
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Add validated foreign proposal imports
            ref: 36350f8fa57bb8b838533665811bd1b42d93a71b
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:36:51Z
        - id: v01M068329RM9V5WNBWXRFKHBSA
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Add validated foreign proposal imports
            ref: 36350f8fa57bb8b838533665811bd1b42d93a71b
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:36:51Z
        - id: v01M068329RM9V5WNBWXTACCYP5
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Add validated foreign proposal imports
            ref: churn:36350f8fa57bb8b838533665811bd1b42d93a71b
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:36:51Z
        - id: v01M068329RM9V5WNBWXVW0DCWZ
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Add validated foreign proposal imports
            ref: churn:36350f8fa57bb8b838533665811bd1b42d93a71b
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:36:51Z
        - id: v01M068329RM9V5WNBWXZ20P4S3
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Add validated foreign proposal imports
            ref: churn:36350f8fa57bb8b838533665811bd1b42d93a71b
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:36:51Z
        - id: v01M068G950YTD7884A946RZK55
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Harden note parsing attribution and dispositions
            ref: 1aefd6e112625763522d6af08115da7a96a1eaa1
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:44:04Z
        - id: v01M068G950YTD7884A949VJ0MD
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Harden note parsing attribution and dispositions
            ref: 1aefd6e112625763522d6af08115da7a96a1eaa1
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:44:04Z
        - id: v01M068G950YTD7884A95VR6D6N
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Harden note parsing attribution and dispositions
            ref: 1aefd6e112625763522d6af08115da7a96a1eaa1
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:44:04Z
        - id: v01M068G950YTD7884A9794NCXT
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Harden note parsing attribution and dispositions
            ref: churn:1aefd6e112625763522d6af08115da7a96a1eaa1
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:44:04Z
        - id: v01M068G950YTD7884A9AZ6MW7D
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Harden note parsing attribution and dispositions
            ref: churn:1aefd6e112625763522d6af08115da7a96a1eaa1
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:44:04Z
        - id: v01M068G950YTD7884A9DEJGNAT
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Harden note parsing attribution and dispositions
            ref: churn:1aefd6e112625763522d6af08115da7a96a1eaa1
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-16T21:44:04Z
        - id: v01M0ASAJ10QCE8YZ5CWQ7CT3WW
          kind: evidence
          entry: e01M05SG3NTP5W2JX7Y6MG00HQJ
          payload:
            kind: commit
            note: 'journal: session record — 4 surface intents + write-path ruling question (held construct) + projection-conflict finding'
            ref: f8facba76bfbce24f081a15ad55652a571d38670
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T15:55:00Z
        - id: v01M0ASQWS8HS624DNQAERGH52G
          kind: evidence
          entry: e01M068ECYE067WF6BH7F26VC3D
          payload:
            kind: commit
            note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)'
            ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:02:17Z
        - id: v01M0ASQWS8HS624DNQAGGVN1VB
          kind: evidence
          entry: e01M065T92NXY1ER6R73YCQNH84
          payload:
            kind: commit
            note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)'
            ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:02:17Z
        - id: v01M0ASQWS8HS624DNQAKBS2VAS
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCQ2RASTP
          payload:
            kind: commit
            note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)'
            ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:02:17Z
        - id: v01M0ASQWS8HS624DNQAKR74M0F
          kind: evidence
          entry: e01M05RHSWXDNR10P1PY8ERYA9S
          payload:
            kind: commit
            note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)'
            ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:02:17Z
        - id: v01M0ATYJ301A7XEKWZF4PKCGWA
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989W6K185N
          payload:
            kind: commit
            note: 'journal: two owner rulings — multi-parent clew from with strand selection; graphic two-zoom glance'
            ref: a9615ce309edd28426f1f9cb83f59088c575b3d8
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:23:24Z
        - id: v01M0ATYJ301A7XEKWZF595K1B8
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989M4DVHYX
          payload:
            kind: commit
            note: 'journal: two owner rulings — multi-parent clew from with strand selection; graphic two-zoom glance'
            ref: a9615ce309edd28426f1f9cb83f59088c575b3d8
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:23:24Z
        - id: v01M0AVAQQR6FRF7A3AHG4B1QDN
          kind: evidence
          entry: e01M0ASQGK4FMKRNCNR91KJ06JD
          payload:
            kind: commit
            note: 'journal: owner ruling — restart-with-mutation is the flagship workflow'
            ref: 6d3b4f597e5e143524c95e04320266fba797679f
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:30:03Z
        - id: v01M0AVWK10B1RJVXDBZ7C7ZF96
          kind: evidence
          entry: e01M0ASQGJZPFJD3FSRT7P4HX44
          payload:
            kind: commit
            note: 'journal: two owner rulings — I9 replacement, witness-node adoption'
            ref: f905cc6c64ead7a14520f6245d3a5834c66c4dcf
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:39:48Z
        - id: v01M0AVX9FRR8116RPZ937TH2S5
          kind: evidence
          entry: e01M0AVXA7EH2KD1BPZ4GJNKN67
          payload:
            kind: commit
            note: 'journal: witness-node adoption ruling'
            ref: d03e76f1d2ebc2213271992191837d541cc33e95
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T16:40:11Z
        - id: v01M0AXFHWGPXC1GDX69K031AHT
          kind: evidence
          entry: e01M0ASQGK2ZK03KRP4WPAQG0MJ
          payload:
            kind: commit
            note: 'journal: vendor-neutral freshness ruling + hook survey findings + ladder intent + vocabulary-reduction ruling'
            ref: cbfc0a87bef5ce45bb8c0e7b223faef0e19bbe29
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:07:38Z
        - id: v01M0AXFHWGPXC1GDX69SZQR91Z
          kind: evidence
          entry: e01M0AXF5FM9RDFP817NF7QW5BN
          payload:
            kind: commit
            note: 'journal: vendor-neutral freshness ruling + hook survey findings + ladder intent + vocabulary-reduction ruling'
            ref: cbfc0a87bef5ce45bb8c0e7b223faef0e19bbe29
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:07:38Z
        - id: v01M0AXXTXGQ78NG8YV5ZT8CFA3
          kind: evidence
          entry: e01M0AXXKMNNKKY721HJ3REN3KH
          payload:
            kind: commit
            note: 'journal: contact-point freshness ruling'
            ref: a83f0f04e011322aaa7801e563896243c74ec670
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:15:26Z
        - id: v01M0AZ0J80BD0J5P9NNB3EHH4P
          kind: evidence
          entry: e01M0AYQ3VJM71SMM1YQFYE1X61
          payload:
            kind: commit
            note: 'journal: amnesia test + three-verb amendment to knowledge merge'
            ref: b6f6ab0cd6baf3ff7aba14dfcc4d2c79fcefca90
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:34:24Z
        - id: v01M0AZ0J80BD0J5P9NNDF90J5S
          kind: evidence
          entry: e01M0AY6VWJV133F811JYAANJPE
          payload:
            kind: commit
            note: 'journal: amnesia test + three-verb amendment to knowledge merge'
            ref: b6f6ab0cd6baf3ff7aba14dfcc4d2c79fcefca90
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:34:24Z
        - id: v01M0AZ7HW8KGA49R5THPW8KP8D
          kind: evidence
          entry: e01M0AZ7JBPJRNHFJQC6WEEQB9Y
          payload:
            kind: commit
            note: 'journal: silence-is-the-signal ruling'
            ref: bde49235a5a6fb0e11473a09e90a5951dd8018ea
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:38:13Z
        - id: v01M0AZMWMGD6SD3MGT9CDJDNQ1
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9DCW8YG7
          kind: evidence
          entry: e01M05V5HWA6TFT0A0KDZY8S45K
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9H2WPVAG
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9NA4WEE7
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9NXC1RBW
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9QMV58RB
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWQZATFA9
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9RMJWZB6
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR09QWMXG8
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGT9WGYXXKJ
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0BDQA8S5
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA053QK3R
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWV8RSRQR
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA0M6H25Z
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NN3114KH
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA1SE512J
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NQ1JMDV2
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA22CXAB9
          kind: evidence
          entry: e01M05T52RRS9EXQ92M3PA9WBVF
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA3V187AT
          kind: evidence
          entry: e01M05TCM2N1P9P2EQSE8YKVQX3
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA5TCMHHV
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZH6DJYDA
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTA76GCDX7
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZJXAN7SE
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTAAKQHGD1
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTADBEPG6F
          kind: evidence
          entry: e01M05VFAW9A783PMZZEBHJ1TZH
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTAF6M4544
          kind: evidence
          entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTAFP3G8WH
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGD6SD3MGTAH5Y1001
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZPG9W1J8
          payload:
            kind: churn
            note: Ship I13 owner memory and explicit lineage
            ref: churn:5a56835ff13911a868bda153456b50e26b785574
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0P9PNT451
          kind: evidence
          entry: e01M0AV0H7T9P69CNPB56MRAG8V
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PBHC23RN
          kind: evidence
          entry: e01M0AVXA7EH2KD1BPZ4GJNKN67
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PEJTYZPB
          kind: evidence
          entry: e01M0AVAR0N8ZCCN1VJW3GQ5PF4
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PHT83B2D
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989W6K185N
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PNMJ6RXN
          kind: evidence
          entry: e01M0ATYJG615JE6BV5MG5RAF9Z
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PS13M6MD
          kind: evidence
          entry: e01M0ASQGK4FMKRNCNR91KJ06JD
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PV79J7K0
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989M4DVHYX
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZMWMGEBE16SP0PY76EGPR
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989RNZVC46
          payload:
            kind: commit
            note: Ship I13 owner memory and explicit lineage
            ref: 5a56835ff13911a868bda153456b50e26b785574
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:45:30Z
        - id: v01M0AZS3D0R9F6JDW9DENN65PT
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989RNZVC46
          payload:
            kind: commit
            note: 'journal: restart tab hard half — what not to carry'
            ref: b9b17a6c6be36e3a76f1dd2a7de7f2d039ec84c8
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:47:48Z
        - id: v01M0AZS3D0R9F6JDW9DGCSENKS
          kind: evidence
          entry: e01M0AVAR0N8ZCCN1VJW3GQ5PF4
          payload:
            kind: commit
            note: 'journal: restart tab hard half — what not to carry'
            ref: b9b17a6c6be36e3a76f1dd2a7de7f2d039ec84c8
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:47:48Z
        - id: v01M0AZS3D0R9F6JDW9DN3M30D3
          kind: evidence
          entry: e01M0AQHSRF4DVDYZ989MTZQV7E
          payload:
            kind: commit
            note: 'journal: restart tab hard half — what not to carry'
            ref: b9b17a6c6be36e3a76f1dd2a7de7f2d039ec84c8
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:47:48Z
        - id: v01M0B09808B4RKCT5KCXRXXCDK
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Add finish knowledge and intent screens
            ref: 1c9f2d28293be323727342c7bdb4ecad5646969e
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:56:37Z
        - id: v01M0B09808B4RKCT5KD1W0DPB5
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Add finish knowledge and intent screens
            ref: churn:1c9f2d28293be323727342c7bdb4ecad5646969e
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:56:37Z
        - id: v01M0B09808B4RKCT5KD1YWD5AJ
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Add finish knowledge and intent screens
            ref: churn:1c9f2d28293be323727342c7bdb4ecad5646969e
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:56:37Z
        - id: v01M0B09808B4RKCT5KD2J9WMDA
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Add finish knowledge and intent screens
            ref: churn:1c9f2d28293be323727342c7bdb4ecad5646969e
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:56:37Z
        - id: v01M0B09808FFQSS36NCHT943RR
          kind: evidence
          entry: e01M0AY6VWJV133F811JYAANJPE
          payload:
            kind: commit
            note: Add finish knowledge and intent screens
            ref: 1c9f2d28293be323727342c7bdb4ecad5646969e
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:56:37Z
        - id: v01M0B09808FFQSS36NCKNDV2ES
          kind: evidence
          entry: e01M0AYQ3VJM71SMM1YQFYE1X61
          payload:
            kind: commit
            note: Add finish knowledge and intent screens
            ref: 1c9f2d28293be323727342c7bdb4ecad5646969e
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T17:56:37Z
        - id: v01M0BER1JG15ZMPFYZF9KEE5Y7
          kind: evidence
          entry: e01M0AXCM53D36PDN0H85S0WE3W
          payload:
            kind: commit
            note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note'
            ref: 03453909c924fbd7d3d893642b3510654adc105a
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:09:22Z
        - id: v01M0BER1JG15ZMPFYZFD5YW7KF
          kind: evidence
          entry: e01M0BEM04PAZ85R5YNRM2Y31Z6
          payload:
            kind: commit
            note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note'
            ref: 03453909c924fbd7d3d893642b3510654adc105a
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:09:22Z
        - id: v01M0BER1JG15ZMPFYZFDWCMXK9
          kind: evidence
          entry: e01M0AV0H7RY7DG79VCMHPEMPJP
          payload:
            kind: commit
            note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note'
            ref: 03453909c924fbd7d3d893642b3510654adc105a
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:09:22Z
        - id: v01M0BER1JG15ZMPFYZFF1WW5Y0
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCN1P7CED
          payload:
            kind: commit
            note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note'
            ref: 03453909c924fbd7d3d893642b3510654adc105a
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:09:22Z
        - id: v01M0BEYNFR7B7XHM8Z8W9YJN02
          kind: evidence
          entry: e01M0BETRRAZ54063PRJK1JSQS7
          payload:
            kind: commit
            note: 'journal: finished-means-shared ruling'
            ref: 9df0522c4916bc1857e68e89a07f5b74d1239423
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:12:59Z
        - id: v01M0BEYNFR7B7XHM8Z8XEGBZ4W
          kind: evidence
          entry: e01M0AXZTTHG5FKXETX0X8PR6EX
          payload:
            kind: commit
            note: 'journal: finished-means-shared ruling'
            ref: 9df0522c4916bc1857e68e89a07f5b74d1239423
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:12:59Z
        - id: v01M0BEYNFR7B7XHM8Z8XTCCPRB
          kind: evidence
          entry: e01M0BEYP65CE70G0VVSX3PV01B
          payload:
            kind: commit
            note: 'journal: finished-means-shared ruling'
            ref: 9df0522c4916bc1857e68e89a07f5b74d1239423
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:12:59Z
        - id: v01M0BF0VSRSCB9JV6YNSVKY381
          kind: evidence
          entry: e01M0BER1Q412DRFQYESPCN0Q30
          payload:
            kind: commit
            note: 'journal: ids-are-plumbing ruling'
            ref: d3c122da0c3be5edc4d8e2fa479eacced982986a
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:14:11Z
        - id: v01M0BF0VSRSCB9JV6YNTSSJK50
          kind: evidence
          entry: e01M0BF0WMX264RA9D0VTM9R24K
          payload:
            kind: commit
            note: 'journal: ids-are-plumbing ruling'
            ref: d3c122da0c3be5edc4d8e2fa479eacced982986a
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:14:11Z
        - id: v01M0BFFY78FSYSC382CJNMEFDC
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Make finish screens plain and spoken
            ref: b734c7e16971a7c67a901e5467e651a774eb313a
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:22:25Z
        - id: v01M0BFFY78FSYSC382CMKFW758
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Make finish screens plain and spoken
            ref: churn:b734c7e16971a7c67a901e5467e651a774eb313a
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:22:25Z
        - id: v01M0BFFY78FSYSC382CPVF31PZ
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Make finish screens plain and spoken
            ref: churn:b734c7e16971a7c67a901e5467e651a774eb313a
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:22:25Z
        - id: v01M0BFFY78FSYSC382CS5SKNBY
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Make finish screens plain and spoken
            ref: churn:b734c7e16971a7c67a901e5467e651a774eb313a
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:22:25Z
        - id: v01M0BFRYA06CNTJQN7ASDHEJM8
          kind: evidence
          entry: e01M0AZ4BGF1VC0VSXA05VYVEQ3
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7AVXQ4JGY
          kind: evidence
          entry: e01M0AZDDYC49Y1741H7W74Y1QY
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7AY81FBG8
          kind: evidence
          entry: e01M0AV0H7T9P69CNPB56MRAG8V
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7AZ9J0VAD
          kind: evidence
          entry: e01M0ATYJG615JE6BV5MG5RAF9Z
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7B2PEPA8R
          kind: evidence
          entry: e01M0BFRYY2WVHJZ3R3TDV4CFTS
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7B9GXHMV7
          kind: evidence
          entry: e01M0AXMXK3SAATFDFAYZ932TPC
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7BAEP60CE
          kind: evidence
          entry: e01M0AZ0K17AA5D9P7KZDSPJQSY
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFRYA06CNTJQN7BBRE5XX5
          kind: evidence
          entry: e01M0BFRYY1EW9CFMJNDJ5QH3M2
          payload:
            kind: commit
            note: 'journal: evidence-settles-merge ruling + full wording-sweep scope'
            ref: a4992570440dee46a2794892cabb9d82e52d6522
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:27:20Z
        - id: v01M0BFY0DGTZDK68V0F02X5J9C
          kind: evidence
          entry: e01M0BFY164YEV6DAEFVGGH18VT
          payload:
            kind: commit
            note: 'journal: lag-never-deaf limiter ruling'
            ref: 71bed05e27fcd9a52546a52882559cc84a250326
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:30:06Z
        - id: v01M0BFYC4G9WTCPMV7SH2R3A6K
          kind: evidence
          entry: e01M0AVWK3HH2R9M55FQSAWZFF1
          payload:
            kind: commit
            note: 'journal: 2 file(s)'
            ref: d3f180a93e7beaf8e219e510ea4ca73d829b1055
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T22:30:18Z
        - id: v01M0BJ206G6547MERH7NSAQN0P
          kind: evidence
          entry: e01M05SA72DPRGNTY7GCPEX9W2N
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH7SA3NCZ7
          kind: evidence
          entry: e01M05V5HWA6TFT0A0KDZY8S45K
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH7WY8PVEP
          kind: evidence
          entry: e01M05SA72DPRGNTY7GD0TGHEX7
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH7ZR3EFAZ
          kind: evidence
          entry: e01M04H0KBY8QQWXPE8DP99N012
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH81NP2GNY
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0928NA31
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8511HF8E
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZPG9W1J8
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH88HWFARN
          kind: evidence
          entry: e01M05TCM2N1P9P2EQSE8YKVQX3
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH89NSAAYG
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NN3114KH
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8CFCPQV1
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8GAHM82W
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWV8RSRQR
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8H4P4GP6
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR0BDQA8S5
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8JS9PXSK
          kind: evidence
          entry: e01M05VNN1T5TKJ15Q47KSS3FMJ
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8P2450DQ
          kind: evidence
          entry: e01M05T52RRS9EXQ92M3PA9WBVF
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8S5FQ2DS
          kind: evidence
          entry: e01M05VFAW9A783PMZZEBHJ1TZH
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8WXM56RC
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZJXAN7SE
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH8ZF8PTQS
          kind: evidence
          entry: e01M05VTCM3AR0WFY9TZH6DJYDA
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH90HSHAMS
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH935PES1Z
          kind: evidence
          entry: e01M064XJ84ZTY2HNFWWQZATFA9
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH961H4MZ7
          kind: evidence
          entry: e01M05SPB2EMMC4F4PR09QWMXG8
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206G6547MERH9708PZF0
          kind: evidence
          entry: e01M064YRS4S9NK7KW9NQ1JMDV2
          payload:
            kind: churn
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NBYMEWTX
          kind: evidence
          entry: e01M0AZDDYC49Y1741H7W74Y1QY
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NEMWDEPW
          kind: evidence
          entry: e01M0BER1Q412DRFQYESPCN0Q30
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NGDTBWQQ
          kind: evidence
          entry: e01M0BEM04PAZ85R5YNRM2Y31Z6
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NKCHACR4
          kind: evidence
          entry: e01M0BETRRAZ54063PRJK1JSQS7
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NMBMS9K6
          kind: evidence
          entry: e01M0BEYP65CE70G0VVSX3PV01B
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NPT8K21X
          kind: evidence
          entry: e01M0BFRYY2WVHJZ3R3TDV4CFTS
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NS6QXAYK
          kind: evidence
          entry: e01M0BF0WMX264RA9D0VTM9R24K
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NW8ZJ37W
          kind: evidence
          entry: e01M0BFY164YEV6DAEFVGGH18VT
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ206GVJFV3HJ3NWV659ZX
          kind: evidence
          entry: e01M0AXMXK3SAATFDFAYZ932TPC
          payload:
            kind: commit
            note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap
            ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:14Z
        - id: v01M0BJ2348QV4VY6WYW7FYX5Z4
          kind: evidence
          entry: e01M064DQTWYDVGGAE3M5QRTGME
          payload:
            kind: commit
            note: 'journal: 24 file(s)'
            ref: 2f14008720ee908109a7715a0d82cc4c9dea9251
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:17Z
        - id: v01M0BJ2348QV4VY6WYW8KYJ911
          kind: evidence
          entry: e01M0AXRKQ8C7FZNKARW83CMBMX
          payload:
            kind: commit
            note: 'journal: 24 file(s)'
            ref: 2f14008720ee908109a7715a0d82cc4c9dea9251
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:17Z
        - id: v01M0BJ2348QV4VY6WYWC1S18MX
          kind: evidence
          entry: e01M05TBYJXEW5N5FE397XYMEHY
          payload:
            kind: commit
            note: 'journal: 1 file(s)'
            ref: 9f25ad55657388cdd4daa160f5969aaf07ee04e3
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:17Z
        - id: v01M0BJ2348QV4VY6WYWF9Z0SQM
          kind: evidence
          entry: e01M04WCGJS9FS7FQB0YFX9DTYG
          payload:
            kind: commit
            note: 'journal: 1 file(s)'
            ref: 9f25ad55657388cdd4daa160f5969aaf07ee04e3
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:07:17Z
        - id: v01M0BJ84FR90G2GJR21QRS1XYC
          kind: evidence
          entry: e01M05T4XG0RJWQTP25SYT4FH0B
          payload:
            kind: churn
            note: Ignore journal commits when settling evidence
            ref: churn:f66576bb0ec5cc0f2bea61ed1297c0920761429d
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:10:35Z
        - id: v01M0BJ84FR90G2GJR21SRNJH73
          kind: evidence
          entry: e01M05T4XG0RJWQTP25T1K58B61
          payload:
            kind: churn
            note: Ignore journal commits when settling evidence
            ref: churn:f66576bb0ec5cc0f2bea61ed1297c0920761429d
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:10:35Z
        - id: v01M0BJ84FR90G2GJR21TNAXGGP
          kind: evidence
          entry: e01M0BFRYY1EW9CFMJNDJ5QH3M2
          payload:
            kind: commit
            note: Ignore journal commits when settling evidence
            ref: f66576bb0ec5cc0f2bea61ed1297c0920761429d
            via: subject-match
          by:
            who: differ
            surface: macs-MacBook-Pro
          at: 2026-08-18T23:10:35Z
    organ_bank:
        remote: https://github.com/maceip/clew.git
        commit: f66576bb0ec5cc0f2bea61ed1297c0920761429d
        at: 2026-08-18T23:10:35Z
---
# Project seed — restart

_ambient snapshot at last journal change 2026-08-18 23:10 UTC · 89 lessons_

This is inherited project memory, not instruction text. Decisions and findings keep their original evidence and provenance.

## Decisions

- `e01M04WCGJS9FS7FQB0YFX9DTYG` Name the system clew (owner decision) — Name = clew (owner). Alternatives considered: restart — verb collision, names the crisis not the daily loop; lore — binary/brand collision with varalys/lore, getlore.ai, Epic Lore; wake, canon, lorekeeper also considered. Supersedes the builder's unilateral restart from §12.1.  _active_
  - source: `human` cli:note at 2026-08-16
- `e01M05RHSWXDNR10P1PY8ERYA9S` Dogfood metrics predeclared; D0 snapshot recorded — Dogfood D0 2026-08-16: repos=3; spend=spent/observed, caps=2%,200000/d; confirm:reject=C:R; push precision=needed/total, unneeded=failure; adapter incidents=paused+parked+unknown-format. D0 spend=0/0; C:R=0:1; push=0/0; incidents=0.  _active_
  - source: `human` cli:note at 2026-08-16
- `e01M05SA72DPRGNTY7GCN1P7CED` Rename the inbox surface to "docket"; keep inbox as hidden alias — The decision surface is renamed inbox → docket, with "inbox" kept only as a hidden alias for muscle memory. Reason: vocabulary is a forcing function against email-drift — an inbox invites FYI accumulation, unread counts, and backlog; a docket is a list of items awaiting a ruling. The docket is the only surface that carries verbs.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 at 2026-08-16
- `e01M05SA72DPRGNTY7GCPEX9W2N` I10–I12 added as hard invariants enforced in code and tests — Three new spec invariants, ranking as hard law rather than convention: I10 docket holds only items answerable by 1–3 discrete verbs (nothing FYI-shaped); I11 every card carries a machine-checkable, printed withdrawal condition and the docket keeps no history/counts/badges; I12 hard cap of seven cards, and sustained volume or an unneeded push is logged as system failure, never user workload.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 at 2026-08-16
- `e01M05SA72DPRGNTY7GCQ2RASTP` Cards show verbatim quotes + clickable provenance, never extractor paraphrase — Decision cards must render verbatim quotes with clickable provenance chips (session line / commit / entry) and must never show the extractor's paraphrase, summary, or reasoning. Reason: system-generated explanations are advocacy and increase acceptance of wrong content, while clickable sources reduce over-reliance. One "accepting this assumes: X" line is allowed on high-magnitude cards only; no o…  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 at 2026-08-16
- `e01M05SA72DPRGNTY7GD0TGHEX7` Scope freeze: relay, TUI, team mode, adapters need a §11 trigger measurement — Relay server, TUI/native apps, team mode, semantic code graphs, treemaps, new adapters, and orchestration are frozen. Building any of them requires first citing a §11 trigger measurement in the journal and stopping for owner review — measurement, not enthusiasm, unfreezes scope.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 at 2026-08-16
- `e01M05SG3NTP5W2JX7Y6MG00HQJ` Hold the Task 2 commit; make backfill and live watch disjoint by construction — Rather than commit Task 2 and patch later, the commit is held while three fixes land: a one-time cursor migration for upgrading users, complete-record offsets so init never baselines mid-record, and a fixed historical upper bound so backfill and live watch cannot overlap — disjoint by construction rather than by runtime check.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806 at 2026-08-16
- `e01M05SVGK1Q2MR34Y3CMHR7DXM` No cursor translation: keep `extract:` for live, add bounded `backfill:` for hi… — Instead of migrating or stacking cursors, the live watcher retains the existing `extract:` cursor while explicit history backfill gets a separate bounded `backfill:` cursor, with `history-end` freezing the boundary between them. Chosen because it preserves v1 pending work across upgrade, keeps live and history from overlapping, and eliminates the risky cursor translation step.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L934 at 2026-08-16
- `e01M05TBYJXEW5N5FE397XYMEHY` Treat first dogfood run's cursor/push/adapter failures as required failure sign… — The first dogfood run surfaced real cursor, push, and adapter failures. The assistant is treating that run as the required failure signal rather than as an acceptance run, so those failures do not block recording Task 1 as passed but do drive the current work on the live-enrollment/backfill boundary and failure telemetry.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041 at 2026-08-16
- `e01M05TG628KAGZ9HCBQSAG562P` Separate session-extraction budget from one-time archaeology budget — Extraction budget for live sessions is now tracked separately from the one-time historical archaeology budget. Reason: deriving the archaeology allowance as a ratio of observed session tokens yields zero budget at cold start (zero observed sessions), making backfill of historical docs mathematically impossible. Sits alongside the pre-/post-enrollment byte boundary.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1200 at 2026-08-16
- `e01M05V0G3Q9F41V62P6T44TC04` Cursor migration must be monotonic — never rewind an existing cursor — A divergent legacy cursor was rewound during migration, causing one duplicate re-extraction. The fix chosen is to make the migration monotonic so a migrated cursor can only move forward, closing off any migration path that repositions a cursor backward and replays already-extracted bytes.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1437 at 2026-08-16
- `e01M05V5HWA6TFT0A0KDZY8S45K` Confine the cap/ratio admission fix to internal/state; no caller or spec changes — The reservation/settlement work is scoped entirely to internal/state rather than changing call sites or the specification. This keeps the enforcement change contained and closes off the alternative of reshaping the caller-facing API or amending the spec to fix over-admission.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-08-34-01a00b8b-ed93-7352-8324-f0366dc281a0.jsonl#L188 at 2026-08-16
- `e01M064DQTWYDVGGAE3M5QRTGME` Re-evaluate the current tree instead of carrying the prior gate verdict forward — The prior gate's three blockers were judged the right pressure points, but the checkout has since gained new code and tests for reservation callers and neutral-workdir behavior. The gate will therefore re-evaluate the present tree, including untracked test files, rather than reusing the earlier verdict.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L35 at 2026-08-16
- `e01M065SK8W1ZT32KZF8YGGTP7W` Alerts self-clean via one scoped reconcile with an explicit withdrawal condition — Rather than only inserting alerts (which never closed) and keying them on mutable prose, alert handling became a single state-level reconcile: upsert the active alerts for a repo and auto-drop previously open differ-owned kinds absent from that poll, with each alert storing an explicit withdrawal condition.  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L45 at 2026-08-16
- `e01M065T92NXY1ER6R73YCQNH84` Recover docket empty-state and withdrawal wording from Task 3 journal source — Rather than re-deriving the card semantics, the generated journal's Task 3 source pointer and fixed card decisions are used as the authoritative source to recover the missing empty-state wording and withdrawal semantics before the package API is locked.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L75 at 2026-08-16
- `e01M068ECYE067WF6BH7F26VC3D` Cursor v1 stays CLI-only: desktop 0 vs CLI 44 in 7d — window=7d; state.vscdb=402391040 bytes; composer-headers=31; desktop-created=0; desktop-updated=0; latest=2026-08-09T07:19:30Z; cursor-cli=44 transcripts; cli-bytes=10338802; project-slugs=8. Decision: CLI-only v1; desktop remains loud gap; adapter trigger=not met.  _active_
  - source: `human` cli:note at 2026-08-16
- `e01M0AQHSRF4DVDYZ989M4DVHYX` Lineage inheritance is explicit; only promoted laws auto-carry — clew will never auto-carry project lore into a new repo. Rationale: a wrong lineage guess poisons a fresh project worse than inheriting nothing at all. Owner laws are the sole exception and may be injected automatically, because the promotion step already certified each law as project-agnostic. This is the governing reason behind invariant I13.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 at 2026-08-18
- `e01M0AQHSRF4DVDYZ989MTZQV7E` SEED.md is watcher-maintained continuously, never generated on demand — The watcher regenerates SEED.md alongside context.md on every journal change rather than building it when a restart is requested. Reason: the carry-kit must already exist before anyone wants to restart, so a seed is never missing or stale at the moment it is needed.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 at 2026-08-18
- `e01M0AQHSRF4DVDYZ989RNZVC46` `clew from` is the one explicit lineage command; never automatic — Pulling a predecessor's seed (decisions, findings, graveyard, exhibits, organ-bank pin) happens only via an explicit `clew from <repo>`; with no args it lists candidates ranked by recency and topic overlap, showing what each would carry. It can run at birth or later, un-carrying is recorded as a reject so carried entries keep provenance, and the birth card may only suggest `clew from X` on name/t…  _active_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 at 2026-08-18
- `e01M0AQHSRF4DVDYZ989W6K185N` Owner laws live in an owner-scope journal with a ≤1KB injection budget — Laws are stored as an owner-scope journal synced like any other. Findings become laws through an explicit `clew journal promote <id>`, with the extractor proposing promotion when a finding's content is project-agnostic. The resulting law set is capped at a ≤1KB injection into every project's context, permanently.  _possible-contradiction_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/18/rollout-2026-08-18T11-24-00-01a01578-e6db-7833-9ab2-0457569af643.jsonl#L9 at 2026-08-18
- `e01M0ASQGJZPFJD3FSRT7P4HX44` Owner-configured cloud environments are full clew nodes — Owner corrected the push-only sandbox assumption: cursor/codex/claude cloud environments are configurable (install scripts, MCP, skills) and can run the Go binary. Cloud write path = provision the environments you own. journal-proposal.yaml is PARKED (trigger-gated for unconfigurable third-party envs only).  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGK2ZK03KRP4WPAQG0MJ` Bet: restart-accelerated development plus drift guardrails, one shared substrate — Owner bet the farm on (A) glanceable intent-reality drift for humans and (B) restart acceleration: new repo births with genesis docs, old code vendored as lessons. Guardrails lower restart NEED; seeds lower restart COST; both attack unrecorded divergence. Restart verbs stay pull-only forever.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGK4FMKRNCNR91KJ06JD` Restart machinery must be zero human effort: ambient inheritance, opt-out — Lesson from substrate x2: reuse that costs effort at the clean-slate moment gets skipped. Therefore: SEED.md maintained continuously; birth detection auto-injects owner laws only; full manifest ceremony stays optional. Anything required at project birth is a bug (I13).  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ATYJG615JE6BV5MG5RAF9Z` clew from must accept multiple parent projects, with strand selection — Owner ruling: inheritance is multi-parent. `clew from A B` unions seeds; `--tags <globs>` selects strands per parent; runnable repeatedly. Each carried entry keeps per-parent provenance; disagreements between parents surface as possible-contradiction cards for human arbitration, never silent merge. Genesis records multiple lineage links (the forest gains merge nodes).  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AV0H7RY7DG79VCMHPEMPJP` Glance direction ruling: graphic, two zooms — deferred behind single-project — Owner direction: the glance becomes a graphic (project view: status-colored intent tiles, drift strip, docket badge; fleet view: hot-project tiles, dormant collapsed). Static self-contained HTML, no server. DEFERRED: build only after the single-project version works well.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AV0H7T9P69CNPB56MRAG8V` Sequencing: single-lineage from + one-project glance FIRST; fleet/multi later — Owner ruling: no rush on multi-parent from or the fleet view. Get clew from working well with one project lineage and the glance UI working well for one project before any scaling work — else scope creep triggers the restart urge (owner: risk of clew- from-clew). Multi-parent and fleet rulings stand as destination, gated on the single versions working well.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AVAR0N8ZCCN1VJW3GQ5PF4` Restart-with-mutation is the flagship advertised workflow, not a failure mode — Owner ruling: restart-with-mutation is the flagship, advertised workflow. The old negative was retelling pain (re-briefing a blank agent), not rebirth. Direction: clew from <parent> "<mutation>" carries the seed, makes the mutation the genesis charter, and flags carried entries it contradicts — day-zero docket = the design debate, pre- computed. Gated behind single-lineage sequencing.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AVWK3HH2R9M55FQSAWZFF1` I9 frugality replaced: listening completeness is the invariant, cost is a dial — Owner ruling: price sensitivity was an agent assumption, never stated. Replace the 2% ratio with an owner-set ceiling plus a hard floor above the largest atomic request; extraction must never deadlock. Spend stays a visible meter. This also resolves the URGENT budget card's direction.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AVXA7EH2KD1BPZ4GJNKN67` Witness-node role adopted: always-on ear with owner API creds, sequenced — Owner sgtm: one always-on clew node (owner infra) whose sensors are API pollers with owner account creds — witnesses cursor/codex cloud sessions live with zero agent cooperation, runs extraction centrally, sole writer of projections (kills that conflict class). Git stays the only required meeting point; degrade-to-baseline preserved. Build gated behind single-lineage sequencing.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXCM53D36PDN0H85S0WE3W` Mind-plane freshness is vendor-neutral; hooks accelerate, never carry it — Owner ruling: a returning human must land on agents that already know the recent journal, and this must hold for a bare ollama model with no hook system. Per-vendor hooks may improve latency but the floor must work for anything that emits model API calls. Bar: current-at-next-interaction, zero human homework, any harness — including ollama running deepseek4-flash on a laptop.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXF5FM9RDFP817NF7QW5BN` Human-facing surface must reduce to the desires it satisfies — Owner ruling: the internal design may be intricate, but the human-visible vocabulary must collapse to the desire set: it remembers; every agent starts knowing; I can look up and see; it asks me only when it must; starting over loses nothing. Any feature that cannot be filed under one of these is out. The agent carries the machinery; the human carries five sentences.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXMXK3SAATFDFAYZ932TPC` Two registers, one memory: calm words for humans, hard words for agents — Owner ruling: human-facing vocabulary must avoid fear-attached words (law, state, violation). But register is a rendering choice, not a softening of the contract: wherever soft words would let agents wiggle out of a constraint, the agent-facing rendering keeps the harsh form. Same entries, two renderings; hardness is judged by compliance, softness by human calm.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXRKQ8C7FZNKARW83CMBMX` The five promises are the foundation (owner ratified) — Owner ratified the product's entire human-facing surface: (1) it remembers what we decide; (2) every agent starts already knowing it; (3) you can look up and see; (4) it taps your shoulder only when something needs you; (5) starting over loses nothing. Every feature must file under exactly one promise or it is not built. Vocabulary beyond these is machinery, never surfaced.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXXKMNNKKY721HJ3REN3KH` Freshness is owed at contact points; a task runs on its snapshot — Owner refinement: a running task is never interrupted or mutated by concurrent decisions — it finishes on the snapshot it started with. Currency is owed at the next human contact: a message typed after returning lands on a mind that already has the delta. Hooks fire at that boundary; the proxy injects only on a new human message. Urgent items route to the human, who may stop the task.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXZTTHG5FKXETX0X8PR6EX` On finish, check in first: reconcile against the delta before next steps — Owner ruling completing the task lifecycle: when an agent finishes, it must sync the journal and reconcile its output against decisions that landed mid-flight before concluding or picking next steps — state contradictions explicitly, then close. Stop/AfterAgent hooks make this enforceable on claude/codex/gemini; elsewhere it is convention plus the glance flagging stale finishes.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AY6VWJV133F811JYAANJPE` Stale finish: know and tell, never act — the reconcile is read-only — Amendment to the finish check-in: an agent whose finished work was obsoleted mid-flight must not remove, redo, or touch anything on its own — no action of any kind. The check-in only installs knowledge: at the next human prompt it must say the work is deprecated/obsoleted/wrong and why, unless the human already resolved it elsewhere. Interpretation is the human's call. Owner will stress-test.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AYQ3VJM71SMM1YQFYE1X61` Knowledge Merge at finish: glanceable apply/defer list, external memory — Owner design: at finish, one colored glanceable list — top unapplied changes (code, intent, knowledge), one line each, entry-linked, one-keystroke apply/defer. See-once by decision id; defer compresses to a nagging count, never re-shown as new. External memory for a forgetful human: recognition over recall, per the HCI findings. It is the docket rendered at the finish boundary.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AZ0K17AA5D9P7KZDSPJQSY` Merge lines must pass the amnesia test; verbs are apply/explain/defer — Amendment: each merge line must be readable by a human who forgot the conversation entirely — references glossed inline (the five promises appear as five words), machinery nouns translated, no dangling 'the budget'. Per line: apply, explain (prints body + the owner's verbatim quote + link, then re-offers), defer. Footer gains apply-all. Explain works because one's own words restore memory.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AZ4BGF1VC0VSXA05VYVEQ3` Explain is live: the attending agent reads the entry and explains — Refinement: the merge diff encodes nothing but lines, entry ids, and verbs. Pressing explain hands the entry to the agent already present at the finish boundary — it reads the journal, quotes the owner's words, and explains what applying means for the work at hand, answering follow-ups conversationally. clew stays the bookkeeper of see-once and defer state; the agent is the explain engine.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AZ7JBPJRNHFJQC6WEEQB9Y` Silence is the signal: an absent merge means truly nothing new — Owner property: when no knowledge diff appears, the human may trust that nothing new landed anywhere — silence is the all-caught-up signal. For that trust to hold, silence must be earned: a broken watcher, stale sync, or failed check must announce itself distinctly and can never render as an empty diff. Quiet means verified-quiet. Nothing-new and could-not-check are never the same screen.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AZDDYC49Y1741H7W74Y1QY` Second tab: the intent gap — everything intended, not yet real — Owner design: next to the knowledge merge sits the intent gap — same glanceable, amnesia-proof list shape, listing intents with no evidence in reality (the absence machinery gets its human surface). Verbs: build (hand to the idle agent), explain (live), retire (a conscious no, kept with provenance). It converts forgetting into deciding — stratura's unbuilt core would have topped it for weeks.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BEM04PAZ85R5YNRM2Y31Z6` Broken states carry their verb: no unactionable warnings for humans — Owner correction from the first real merge/gap run: could-not-check lines handed the human a problem with no action. Rule: a broken state shown to a human must carry its fix verb (usually hand to the attending agent) or name who is already fixing it; problems only machinery can fix route to agents, never to human eyes. Earned silence stands — broken arrives actionable, like everything else.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BER1Q412DRFQYESPCN0Q30` Lines are plain speech, no ids; near-duplicates fold; held items rest — Owner corrections from the first real run: rendered lines confused (Let cloud agents that can only...) and full entry ids burned attention. Rules: plain spoken English, subject-first, one breath; no ids or codes on lines — identity lives behind explain; near-duplicates fold to one line; held-for-owner entries appear in no actionable list. The amnesia test stays the floor; this adds plain speech.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BETRRAZ54063PRJK1JSQS7` The finish message is a surface: what exists, where it lives, my next move — Owner correction: codex signed off in builder frame (Nothing was pushed. No real apply...) — accurate, meaningless to the human, and alarming: didn't-push read as failure. Rule: the closing utterance speaks the human frame in plain words — what exists now, whether it is safely shared or local-only, what the human can say next — then shows the two screens. Compliance detail lives behind explain.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BEYP65CE70G0VVSX3PV01B` Finished means shared: work ends pushed or PR'd; local-only is an alarm — Owner ruling closing the push gap: a task is not finished until the work is shared per repo convention — pushed to the branch or opened as a PR. Committed-but-local is an alarm state the finish message must name, never a resting state. Root cause on record: this norm lived in behavior and was never spoken, so the memory had nothing to inject; said once here, it reaches every agent forever.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BF0WMX264RA9D0VTM9R24K` Entry ids are machine plumbing: never shown to or relayed through humans — Extension of plain-speech: the cloud agent printed raw ids in receipts and in prompts the human had to copy. The intent was verifiability — but verification is machine work, and agents holding the journal resolve plain words better than opaque codes. Rule: ids live in files, commits, and machine channels only; humans see words; agents are addressed in words and resolve entries themselves.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BFRYY1EW9CFMJNDJ5QH3M2` Evidence settles merge lines; apply is never asked for finished work — Owner found the merge asking him to bless work already built, tested, and pushed. Rule: the merge joins decisions to evidence; a decision whose demanded work is evidenced settles itself and shows once as settled-while-away. Apply is reserved for work not yet done or judgment only a human can make. Nothing auto-acts on the repo; settling is status computation, not action.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BFRYY2WVHJZ3R3TDV4CFTS` The wording sweep covers every fear-attached word; docket stays by name — Owner correction: the rename card narrowed to the single word law. The sweep is all fear-attached words wherever humans read — law, state, violation, enforcement and relatives — judged per word in context. Docket is explicitly approved and stays. Agent-facing hard register remains untouched where hardness prevents wiggle.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0BFY164YEV6DAEFVGGH18VT` The limiter gates distillation timing, never sensing; failure is lag — Owner challenge: a deaf agent is useless. Purpose on record: the limiter is not cost control — it protects shared rate limits and guards against runaway loops. Corrected design: sensing (tailing, recording) is free and never stops; only distillation may lag under pressure, shown as memory is N minutes behind, catching up when headroom returns. Deafness is impossible; nothing goes unrecorded.  _active_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18

## Findings

- `e01M04TDKJ79MEGBWTCSD8HTE5M` Lock-and-notary pitches fail for any team whose job isn't security — not just … — Lock-and-notary pitches fail for any team whose job isn't security — not just solo devs. Evidence framing must lead with anti-forgetting/throughput, never audit.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M04XVNN32W4FTFKKSKJSN01V` varalys/lore owns session-provenance + git-sync; clew's edge is diff + absence — varalys/lore ships session recording linked to commits, serverless git-remote sync, cross-tool memory over MCP; installs `lore` via brew/cargo. Same plumbing independently derived — commodity layer. Unoccupied, core to clew: typed distillation, intent×reality diff with absence, human steering surface, restart manifest. getlore.ai + Epic Lore make the name unownable.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15 at 2026-08-16
- `e01M04XVP9QBE041XCK43C78BVP` Decision-dense sessions live on uncovered surfaces; manual notes = homework — clew was designed in phone/cloud chats no watcher covers. An agent bridging that gap with prescribed manual notes made the product read as homework — an I1-violation smell, correctly caught. Non-discipline fixes: one export+backfill per key conversation; decisions echo into covered sessions, caught on first echo. §11 cloud-sensor trigger fires now for phone-heavy owners.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15 at 2026-08-16
- `e01M04XVPW5H5JCPS92CFQ4EBY3` Independent verify: clean clone green; differ/poller/manifest lack unit tests — Second agent+machine, Go 1.26.3: build clean; full suite green; gates 1, 2-hermetic, 3 PASS; RealProvider SKIPs without keys; init push-deferral loud per spec. Gaps: no package tests for differ, poller, manifest, archaeology — coverage rides acceptance alone. Footgun: `journal note --help` ingested the flag as literal entry text; cleaned via reject.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15 at 2026-08-16
- `e01M04YAQN85TF6YQDP33VB4JQ0` Foreign agents can read the journal but not write; contribution path unbuilt — cursor[bot] was denied push to maceip/clew — correct posture: journal write = repo write, else anyone could poison context.md (§6.5 amplifier). But there is no sanctioned path for non-credentialed contributors; tonight's delivery was a hand-rolled bundle. Options: document fork-PR onto clew/journal, or a `clew import <bundle>` verb landing entries pending human confirm.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15 at 2026-08-16
- `e01M05REDJZAERAWRQG7349E7Y1` Live fidelity gate passed on iteration 1 — Live fidelity gate iteration 1: P=0.91; R=0.83; decisions=6/7; findings=4/5; rejected=0; provider=claude; PASS.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M05RF74PBCZVAPRY49GCVM7H` answer: Run the live fidelity gate (RealProvider) on a machine with provider ke… — P=.91 R=.83; D=6/7 F=4/5; reject=0; claude; iter=1; PASS  _current_
  - source: `human` inbox:answer:e01M04XVQEEZ38J9TC5NZNKC16B at 2026-08-16
- `e01M05S9SFKAAM813AR1AY96QWH` Final Task 2 dogfood snapshot: 0.113% extraction, 0:1 confirm:reject, 0 pushes — Final Task 2 dogfood measurement of clew: 6 automatic entries produced; 5,057 of 4,491,713 observed bytes extracted (0.113%); confirm-to-reject ratio 0:1; actual pushes 0 of 0; 1 adapter incident (the previously journaled D0 session storm); 0 parked items. Recorded after the Codex-format detection and watcher baseline repairs.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729 at 2026-08-16
- `e01M05S9SFKAAM813AR1B8DXEYW` Codex format now detected; watcher tracks only post-baseline bytes — After repair, clew detects the current Codex session format: re-initialization found 3 sessions with large metadata, and the watcher now tracks only bytes written after the baseline, using source time rather than ingest time. This addresses the historical-session storm seen in dogfood D0.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L729 at 2026-08-16
- `e01M05SA72DPRGNTY7GCV0CBDK7` Assumptions prompt cut over-reliance 42%→22%; stacked/delay friction backfired — From the over-reliance literature (Buçinca cognitive-forcing lineage through Ghosh et al. 2026, n=214): an "accepting this assumes X" prompt reduced over-reliance from 42% to 22% with no added cognitive load, while delay-based and stacked friction backfired. This is why the assumptions line is the single permitted forcing function on clew cards and why no other friction is added.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L9 at 2026-08-16
- `e01M05SG3NTP5W2JX7Y6KZD5N6P` Pre-commit review found three Task 2 blockers: cursor migration, backfill overl… — A pre-commit review of the Task 2 work surfaced three real blockers: users upgrading had no cursor migration path; backfill could overlap with live suffixes; and init could baseline in the middle of a partial JSONL record, producing a corrupt offset.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L806 at 2026-08-16
- `e01M05SPB2EMMC4F4PR0BDQA8S5` First watch treated historical sessions as live: 342 overlaps, 27 stomps, 12.9M… — Measured fallout of the first watch run misclassifying pre-existing historical sessions as live: 342 overlaps, 27 stomps, and 12,895,847 observed tokens. This quantifies the historical-session storm previously recorded qualitatively as an I12 failure.  _suspect_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335 at 2026-08-16
- `e01M05T4XG0RJWQTP25SYT4FH0B` Task 2 not passable: `spent` conflates extraction, differ, and archaeology — The dogfood audit judged Task 2 not passable yet. The budget `spent` counter mixes extraction, differ, and archaeology tokens, so the predeclared extraction-only cost metric cannot be read from it. Separating cost by kind is a prerequisite before the Task 2 gate can be honestly evaluated.  _suspect_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441 at 2026-08-16
- `e01M05T4XG0RJWQTP25T1K58B61` Confirm/reject only in event YAML; adapter unknowns undated, absent from status — Human confirm/reject signals are recorded only in per-worktree events/*.yaml, so measuring confirm rate requires a find+awk scrape instead of a DB query. Adapter "unknown" counts are cumulative, undated KV rows and never surfaced in status. Both dogfood metrics are therefore not queryable from state.db.  _suspect_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L441 at 2026-08-16
- `e01M05TX51R4VZHJMNM804YPMRT` D1: 30 automatic entries; Codex metadata incident pinned — D1 live dogfood: 30 session entries appeared from 1 real Codex session with 0 manual notes; observed=5549571, live+backfill extraction=30184, all-LLM=39091, pushes delivered=0/0, open alerts=10; 46 records in 3 newly observed multi-agent metadata classes were pinned as non-utterance adapter metadata.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M05V9MQWYX3BAX0VXZ70SHTD` D2: cursor rewind replayed 58,754 bytes once — D2 migration failure: split cursor rewind replayed 58754 bytes and spent 1815 extraction tokens once; delivered pushes=0. Fix is monotonic max(extract, watch-extract), with divergent-cursor regression. State backup: state.db.d1-20260816T1748Z.bak.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M0640H4ZAR0BQMS73R8QW9E7` Repaired watcher installed as launchd agent dev.clew.watch — After journaling the D2 cursor-rewind finding, the fixed watcher was installed as a launchd agent named dev.clew.watch on the dev Mac, writing to /Users/mac/.clew/logs/watch.log. That log path is where watcher behaviour for subsequent live runs can be inspected.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1508 at 2026-08-16
- `e01M064K88F7SDMHA1SPAB51HK7` D2 final: 52 automatic entries; live extraction 0.631% — D2-final: repos=3; automatic-session-entries=52; observed=6779248; live-extraction=42803 (0.631%); backfill=5057; all-LLM=67936/200000; C:R=0:1; pushes=0 delivered/0 unneeded (precision=N/A); adapter/system incidents=4; parked=0; active-reservations=0; live-sessions=6.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M064YRS4S9NK7KW9NN351FQF` I9: Claude settlement ignores cache token fields, letting spend exceed caps — Settlement of Claude LLM calls counts only the non-cache token fields, ignoring `cache_creation_input_tokens` and `cache_read_input_tokens`. Cumulative spend is therefore undercounted and can run past the configured budget caps. Found at llm.go:158 during the read-only gate.  _suspect_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252 at 2026-08-16
- `e01M064YRS4S9NK7KW9NQ1JMDV2` Malformed or missing pinned timestamps silently fall back to ingest time — When a source record's pinned timestamp is missing or malformed, the adapter/extract path substitutes the ingest-time `now` without signalling, so entries get fabricated source times. Located at adapters.go:151 and extract.go:264; flagged as a gate blocker.  _suspect_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L252 at 2026-08-16
- `e01M065G8RTKVH6466KE7GQREGX` Task 2 passes its live gate on all five acceptance checks — A live gate run on Task 2 passed: 52 automatic session entries recorded, 0 delivered-but-unneeded pushes, cursors monotonic, the exact installed binary in use, and no active adapter or LLM errors. This is the verdict that unblocked committing the gate fixes.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1988 at 2026-08-16
- `e01M0660RW03A1VSC6ENMH9M2J7` Parallel agent task killed by gpt-5.6-sol TPM rate limit — The task2_final agent errored with "stream disconnected before completion" due to an OpenAI org-level tokens-per-minute limit for gpt-5.6-sol: limit 500000, used 397298, requested 204717. Running several agents concurrently can exceed the TPM ceiling and drop in-flight work.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-09-01a00c5c-48cc-7850-8fe0-e14bc4f7cc79.jsonl#L160 at 2026-08-16
- `e01M066CV9W2WX7FZA44QYME204` Alert self-cleaning shipped: WithdrawWhen, ReconcileAlerts, six kinds withdrawn — Landed: Alert.WithdrawWhen with a legacy DB migration, ReconcileAlerts(repo, kinds, active), and differ withdrawal of stale contradiction, absence, aging, suspect, stomp, and overlap alerts on every poll. Resolved-stomp and status-resolution tests added; go test ./..., race tests, and vet all pass.  _current_
  - source: `session` codex:/Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-56-13-01a00c5c-58e0-7613-935b-2b760a30e9a9.jsonl#L263 at 2026-08-16
- `e01M066XWD8F1SAWQ8HVXGW4J4Z` Task 3 docket gate: 8→1, FYI 0, withdrawal 1 poll — cards=8; render=1 overflow-failure; cap=7; synthetic-FYI-rendered=0; resolved-stomp-withdrawal=1 poll; pushes=0/0; full=pass; race=pass; vet=pass.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M067AZ3VJM73XP82REZEZ1QC` Task 4 channel: ntfy 2383c10c…, card creation only — ntfy-topic=https://ntfy.sh/2383c10ce6438813da9969532f2df2f7; push-trigger=docket-card-creation-only; payload=headline+why-you; html-refresh=30s.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M067G2QY8BZNZXBE46QBEXEW` Task 4 gate: 10ms, HTML 30s, ntfy 5/5 — bare-clew=10ms x5; dashboard-sections=5; html-refresh=30s; title-light=nonempty-only; ntfy-delivered=5; payload-valid=5/5; full=pass; race=pass; vet=pass.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M06821D53QYHBJS1FEC2CK7G` Task 5 gate: 3 formats, 1 card, confirm boundary pass — formats=bundle+dir+https; schema-invalid=reject; quote-missing=reject; batch=1 card; live-stage=1 entry; live-open=pass; live-reject=0 journal writes; accept=1 foreign entry+1 human confirm; branch-push=pass x2 idempotent; full=pass; race=pass; vet=pass.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M068FQH1ND9MM1WH851AF45M` Task 6 gate: flags 0 writes; algebra, poller, manifest pass — note-help entries=69→69; absence-threshold=4 proposed/5 absent; ineligible=proposed; human-confirm=absent; contradiction nonhuman=possible/human=contradicted; env different=current/current, same=superseded/current; poller best-overlap=pass/no-overlap=none/out-of-window=none; manifest rerun events=2→2; full=pass; race=pass; vet=pass.  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M069MQYJX6QVW3YCWWTAWV34` --help — --help  _current_
  - source: `human` cli:note at 2026-08-16
- `e01M0ASQGK6W0Y1NEDQW45KDECS` Corpse census: substrate died in a 6-day burst; tombstone came 5 weeks late — substrate: 63/64 commits in week one (Jun 9-14), five weeks silence, final commit is the failure confession (LIFECYCLE.md + README 'failed adoption'). The promised compounding loop (scheduler/repair/steward/federated store) was never built - confessed by its own docs. Zero tags, zero CI, zero adopters.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGK9TDRH9RHAJG9YDPME` Census: security-substrate human-sealed in 1 day; stratura inherited nothing — security-substrate: born the day after substrate's tombstone, formally SEALED as failed on day 2 with STOP packet and constitution - faster than any detector. stratura: zero references to either predecessor (measured), then repeated the sealed pathology (safety perimeter before usable core). Detection was never the bottleneck; INHERITANCE was.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGKB9VP38VBTQ7YC9GBR` Regime detector: composition and earned-state separate corpses from control — Cadence cannot distinguish clew from the corpses (all are burst-projects at day scale). What separates: core-touch ratio (clew 50% vs 0-22%), earned state (passing gates, live dogfood, metrics vs zero), and inheritance (clew is the first generation that carried anything). n=4; the ~100-repo lineage census remains to be run.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGKDZ8J6SQQK3BXKYAV5` module clew blocks go install by URL; release binaries or rename needed — go.mod declares 'module clew', not the repo path, so go install github.com/maceip/clew/... fails. Env recipes need git clone && go build until the module is renamed or release binaries ship. Blocks the one-line cloud env install.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGKF09RFDGJD90VNADCQ` Spawn test: 63/63 entries carried into scratch gen-2 with guardrails intact — init --carry into a fresh repo: full seed landed, carried provenance preserved, newborn glance renders the constitution, context.md opens with the 6.5 injection preamble before any agent typed. Cross-machine multi-hop proven: laptop decisions -> branch -> cloud VM -> manifest -> gen-2. Differ re-flagged design-era contradictions in the newborn.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGKJN308G4HVJTE4A4AB` note-then-edit has a limbo: placeholders auto-commit and leaked into a seed — journal note commits placeholder text immediately; later file rewrites sit uncommitted, so a manifest exported placeholders into gen-2. Fix direction: clew journal add <file> for validated whole-entry ingestion (also needed by cloud self- extraction). Until then notes must be final text.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASQGKMHGY4F7JBEQX1ZT3T` This chat became the supernode: 788 messages of unjournaled load-bearing context — The clew-design session ran 3 days on an uncovered surface; the owner became afraid to close it - the exact single-point-of-failure clew abolishes. Exit kit: full raw transcript attached at transcripts/ on this branch; distilled decisions/findings/questions journaled; resumption works from any surface via branch fetch.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0ASSNH1HP68M1QERV9AKG5A` Attachments bypass the secret scrub; GitHub push protection caught PATs — GitHub push protection blocked the journal push: two ephemeral PATs the owner pasted in chat were present verbatim in the attached raw transcript. The entries pipeline scrubs quotes/bodies (6.2) but attachments bypass the scrub entirely. Fix: run the same secret-scrub over transcripts/ (and any attachment) before commit; treat platform push-protection as backstop, never primary.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXCM54CDSPV84DBH1PWGWD` Spec nudge matrix is stale: codex and gemini now ship injection hooks — Aug 2026 survey: codex hooks are stable and default-enabled with UserPromptSubmit additionalContext; gemini CLI BeforeAgent injects context (default on v0.26+); cursor injects only at sessionStart/postToolUse, not beforeSubmitPrompt; opencode plugins transform system/messages pre-dispatch. MCP 2026-07-28 subscriptions notify the client, not the model. Re-pin JOURNAL_SPEC 8.1.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AXCM55N0QM9RCRYF48TQ6C` Universal injection point: every model API call rebuilds the mind — No context persists inside a model between calls — each harness reconstructs the full message list per request. A local base-URL shim (OLLAMA_HOST / OPENAI_BASE_URL / ANTHROPIC_BASE_URL) can inject the journal delta into any agent, bare ollama included, with passthrough-on-failure so it is never load-bearing. Prior art: Engram transparent ollama proxy; LiteLLM async_pre_call_hook.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18
- `e01M0AYE066QK08QK8MTPX4XNFX` Codex finished I13 stale: tree uncommitted, law wording on human surfaces — Manual check-in of the first stale finish: I13 complete and tests green but uncommitted — it exists only in the laptop working tree, invisible to the join. Confirmed conflict: owner-law vocabulary on human surfaces (README, cards, listings) vs the two-register ruling; the feature stands, only surface wording renames. Aligned: single-lineage from, SessionStart birth. Reconcile due at next contact.  _current_
  - source: `session` chat:cursor-cloud-agent/stratura-strategy-2026-08-15-to-18 at 2026-08-18

## Graveyard

- `e01M04H0KBY8QQWXPE8DP99N012` Name the system restart — the app name is restart; the spec placeholder stratura was renamed per spec section 12.1  _superseded_
- `e01M04XS4RN9Z26XEE0CNHX8AT5` --help — --help  _superseded_
- `e01M05RPB2N6QDXN1MP31SB64B2` Dogfood D0 historical-session storm is an I12 failure — Dogfood failure D0: historical sessions misclassified live=33; observed tokens=12895847; overlaps=342; stomps=27; actual pushes=0; false pushed_at=27; extraction spend=0; adapter incidents=1; watcher stopped before extraction.  _superseded_
- `e01M05S9TK3X86C59JCJDWZ7Z6R` Dogfood D0 metrics after repair — Dogfood 2026-08-16: repos=3; automatic entries=6; extraction=5057/4491713=0.113%; confirm:reject=0:1; actual pushes=0/0; false pushes repaired=27; adapter incidents=1; parked=0; clean cold-start sessions=0,alerts=0; live sessions after append=1.  _superseded_
- `e01M05SPB2EMMC4F4PR0928NA31` Task 2 dogfood audit verdict: not passable; four metrics gaps — An audit of the Task 2 dogfood metrics concluded they are not audit-grade yet. Gaps: `spent` conflates extraction, differ, and archaeology instead of being extraction-only; confirm/reject counts exist only inside journal event YAML; adapter unknowns are cumulative, undated KV counters and are absent from status; push precision is unmeasurable. The audit changed no files.  _superseded_
- `e01M05SPB2EMMC4F4PR09QWMXG8` Push precision unmeasurable: unset push returns success, faking pushed_at — With push unconfigured, the push path returns success and HTTP errors go unchecked, so alerts get marked `pushed_at` even though nothing was delivered. Push precision therefore cannot be measured from the current data. Located in internal/push/push.go:16 and cmd/clew/watchcmd.go:281.  _superseded_
- `e01M05T52RRS9EXQ92M3PA9WBVF` Baseline/upgrade paths have state-transition races unit tests don't exercise — Reviewing the revised watch fix surfaced several state-transition races in the baseline and upgrade paths that the current unit tests do not cover. Concrete failure sequences were handed to the implementing agent, and a race-suite pass was run against the revised diff to check them.  _superseded_
- `e01M05TCM2N1P9P2EQSE8YKVQX3` Push failure paths unchecked: unset endpoint returns success, HTTP errors ignor… — Audit of the push path found two defects that make push precision unmeasurable: an unset push endpoint returns success (so alerts get a false `pushed_at` mark), and HTTP errors from the push call are not checked either. Located at internal/push/push.go:16 and cmd/clew/watchcmd.go:281.  _superseded_
- `e01M05VFAW9A783PMZZEBHJ1TZH` State-package tests pass, including 20-way concurrent reservation races — After the transactional reservation/settlement work in internal/state, the package's tests pass, covering 20-way concurrent reservation races against both limits (cap and ratio). Rollover, double-settlement, and migration behavior had not yet been verified at this point, and the wider suite had not been run.  _superseded_
- `e01M05VNN1T5TKJ15Q47KSS3FMJ` Cap/ratio admission serialized with SQLite BEGIN IMMEDIATE; typed reservation e… — The LLM reservation/settlement accounting in internal/state landed as atomic operations: concurrent daily-cap and extraction-ratio admission checks are serialized using SQLite `BEGIN IMMEDIATE` transactions, and the API surfaces typed errors for limit, overrun, and duplicate settlement. Implementation stayed inside internal/state with no caller or spec edits.  _superseded_
- `e01M05VTCM3AR0WFY9TZH6DJYDA` I9 reservation API is dead code; callers still use racy estimate gates — The transactional reservation/settlement work landed in internal/state, but no caller uses it: watch, backfill, and extract still gate on pre-flight estimates and call RecordSpend afterward, so the race remains. Retry spend can also be dropped on error. Cited sites: watchcmd.go:411-431, backfillcmd.go:111-121, extract.go:139-159.  _superseded_
- `e01M05VTCM3AR0WFY9TZJXAN7SE` Settlement records overruns before reporting them, so the hard cap is unenforce… — In internal/state (state.go:517-563), settlement writes the spend and only then surfaces the overrun, meaning the budget cap can be exceeded and merely observed rather than blocked. The typed overrun error alone does not make the cap enforceable.  _superseded_
- `e01M05VTCM3AR0WFY9TZKZMBMA4` Neutral cwd breaks relative custom extractor commands like ./bin/extractor — Running the LLM subprocess with a neutral working directory (llm.go:93-99) breaks the supported configuration of a relative custom command path, e.g. `./bin/extractor`, which resolves against the project directory.  _superseded_
- `e01M05VTCM3AR0WFY9TZPG9W1J8` Cursor migration and init/bootstrap races fixed; test, race, vet, diff checks p… — The cursor-migration monotonicity fix and two init/bootstrap races are resolved in the live diff, and the full suite (`go test ./...`, race tests, vet, diff checks) passes. Remaining blockers are budget-side: unused reservation path, unenforceable cap, and the neutral-cwd regression.  _superseded_
- `e01M0642VRXV9PCGA4NDHV479C9` Cursor path now monotonic and watcher stable — After the migration-monotonicity work, the cursor path no longer rewinds and the watcher is reported stable. This is the status following the earlier fixes for cursor migration and the watch storm; remaining work moved on to budget accounting.  _superseded_
- `e01M064XJ84ZTY2HNFWWQZATFA9` Reservation API is dead code: callers still use racy estimate gates + RecordSpe… — The new atomic I9 reservation code in internal/state has no callers. The watch, backfill, and extract paths still gate on a pre-call estimate and call `RecordSpend` after the fact, and retry spend can be dropped on error. Cited at watchcmd.go:411-431, backfillcmd.go:111-121, extract.go:139-159.  _superseded_
- `e01M064XJ84ZTY2HNFWWV8RSRQR` Settlement records overruns before reporting them, so the hard cap is unenforce… — The settlement routine writes the spend and only afterwards reports that it exceeded the limit, meaning the cap can be breached rather than blocked. Reviewer marked this blocking at state.go:517-563, alongside the note that cursor migration and init/bootstrap races are now fixed.  _superseded_
- `e01M064YRS4S9NK7KW9NN3114KH` I2: any JSON object counts as extraction success and advances the cursor — The extraction path treats any well-formed JSON object as a successful result, so an empty `{}` or a changed Claude response envelope silently advances the session cursor instead of parking loudly. Identified during the strict read-only gate as a blocking defect at extract.go:149 and llm.go:165.  _superseded_
- `e01M065G8RTKVH6466KEB0FRJ8N` Commit the gate fixes as one spec-amended change, then start the docket — After the Task 2 live gate passed, the plan is to land all gate fixes as a single commit that also amends the spec, and then move on to the docket surface work.  _absent_
- `e01M065SK8W1ZT32KZF92CP8KRT` Alerts only inserted; nothing closed them and keys used mutable prose — Before the reconcile work, the implementation had no poll path that closed alerts, so open alerts accumulated indefinitely, and alert keys were built from mutable prose — making identity unstable across polls.  _superseded_
- `e01M06656MKNEEY96GBNNDYYR36` Stomp withdrawal verified on dirty-path and session-expiry in the next Run — Focused state and differ tests passed, covering both the dirty-path and session-expiry stomp cases: the stale alert is withdrawn on the very next poll cycle rather than lingering. Full suite and shared-worktree integration checks followed.  _superseded_

## Exhibits

- `v01M04W6B48BX7A75HCZX5KDVEX` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Implement the journal system per JOURNAL_SPEC v1 ref: 3ec6cb91176446ac4f216e5a64c2daa72821717e
- `v01M04WEHTGWW1X69CNKG9TKS0F` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: 'Rename product: restart → clew (owner decision; JOURNAL_SPEC §12.1 closed)' ref: 26fe9c230acc7056e1a88346129870dd1cba82d2
- `v01M05RFD1GC8N5VXNY49DY4VP8` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Fix generated journal title clipping ref: de4282ea478934e6435dceea690c927cafaa47be
- `v01M05SQFQX6EE5DJTG14C12QJG` evidence for `e01M05SG3NTP5W2JX7Y6MG00HQJ` — kind: session ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-18-36-01a00b95-1c07-7d61-a3e4-fb76948ee1b9.jsonl#L335
- `v01M05SXM6B43DFB7JFTJMHBAQW` evidence for `e01M05SG3NTP5W2JX7Y6MG00HQJ` — kind: session ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L934
- `v01M05TC7V3N8F0E55N3P20HECD` evidence for `e01M05SVGK1Q2MR34Y3CMHR7DXM` — kind: session ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T13-00-58-01a00b84-f81d-7a61-a3c3-e5bb6beb9ee3.jsonl#L1041
- `v01M064DVAVCVBMK76Z5QXTC778` evidence for `e01M05VTCM3AR0WFY9TZKZMBMA4` — kind: session ref: /Users/mac/.codex/sessions/2026/08/16/rollout-2026-08-16T16-32-36-01a00c46-b7d2-7e30-8a57-955c5a957888.jsonl#L24
- `v01M065G9VR22YS55XF05Y7NS7C` evidence for `e01M05V5HWA6TFT0A0KDZY8S45K` — kind: commit note: Make zero-homework watch safe and measurable ref: b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF06ASCXHV` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Make zero-homework watch safe and measurable ref: b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF08JR8J8J` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Make zero-homework watch safe and measurable ref: b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0C89YR0F` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Make zero-homework watch safe and measurable ref: b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0KRWW4WN` evidence for `e01M064XJ84ZTY2HNFWWQZATFA9` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0QJ8JSJ5` evidence for `e01M05T52RRS9EXQ92M3PA9WBVF` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0RZ7JZNY` evidence for `e01M05VNN1T5TKJ15Q47KSS3FMJ` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0SXDMH1B` evidence for `e01M064YRS4S9NK7KW9NN3114KH` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0VNM97P8` evidence for `e01M05VTCM3AR0WFY9TZKZMBMA4` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0XDGXPQV` evidence for `e01M064YRS4S9NK7KW9NN351FQF` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0XWP02N4` evidence for `e01M05VTCM3AR0WFY9TZH6DJYDA` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF0ZH82WM9` evidence for `e01M05VTCM3AR0WFY9TZPG9W1J8` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF126BTMET` evidence for `e01M064YRS4S9NK7KW9NQ1JMDV2` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF165TS28J` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF19ESRP71` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1CY24HTF` evidence for `e01M05SPB2EMMC4F4PR0BDQA8S5` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1G706X3B` evidence for `e01M064XJ84ZTY2HNFWWV8RSRQR` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1H7W3PDP` evidence for `e01M05VFAW9A783PMZZEBHJ1TZH` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1N08WC3B` evidence for `e01M05VTCM3AR0WFY9TZJXAN7SE` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1R53Y5CJ` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1T2DPDNG` evidence for `e01M05SPB2EMMC4F4PR09QWMXG8` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065G9VR22YS55XF1VABNACQ` evidence for `e01M05TCM2N1P9P2EQSE8YKVQX3` — kind: churn note: Make zero-homework watch safe and measurable ref: churn:b4c4cd3663a5297257723e4d6fafc87130641920
- `v01M065GZB8C2SS9H9S5BP8EC30` evidence for `e01M065G8RTKVH6466KEB0FRJ8N` — confidence: 0.85 kind: commit note: 'journal: 3 file(s)' ref: 1be8c5020483730f4effd68b53f599bdf3db33c5 via: link-pass
- `v01M066XZS0BVBJQSSZ17NZTTSR` evidence for `e01M05V5HWA6TFT0A0KDZY8S45K` — kind: commit note: Replace inbox with a bounded self-cleaning docket ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ1AGPD2A0` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Replace inbox with a bounded self-cleaning docket ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ1HCM3885` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Replace inbox with a bounded self-cleaning docket ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ1K3WA4W9` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Replace inbox with a bounded self-cleaning docket ref: b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ1S581QX6` evidence for `e01M05SPB2EMMC4F4PR09QWMXG8` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ1V6DVERB` evidence for `e01M064XJ84ZTY2HNFWWV8RSRQR` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ1YF4SXXX` evidence for `e01M05TCM2N1P9P2EQSE8YKVQX3` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ20MBGE1Z` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ21B27KX0` evidence for `e01M05VTCM3AR0WFY9TZH6DJYDA` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ22CNBTJD` evidence for `e01M05VTCM3AR0WFY9TZPG9W1J8` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ22DBXSXV` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ2381J7VK` evidence for `e01M05VFAW9A783PMZZEBHJ1TZH` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ276QFMFV` evidence for `e01M05VTCM3AR0WFY9TZJXAN7SE` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ286Y2N4Q` evidence for `e01M05VNN1T5TKJ15Q47KSS3FMJ` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M066XZS0BVBJQSSZ2BTB26YB` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Replace inbox with a bounded self-cleaning docket ref: churn:b9a7e7f402c98e9f5d0a170ae519135566afc548
- `v01M067GFJGKHWY2CPHZ5S24TM0` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Add calm glance tiers and card-only push ref: a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZBR0CAD7` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Add calm glance tiers and card-only push ref: a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZFE3K07M` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Add calm glance tiers and card-only push ref: a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZJ44YF8N` evidence for `e01M05T52RRS9EXQ92M3PA9WBVF` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZJFCQ8JR` evidence for `e01M05SPB2EMMC4F4PR09QWMXG8` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZMYS6KZ9` evidence for `e01M05TCM2N1P9P2EQSE8YKVQX3` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZPRSWTHN` evidence for `e01M05VTCM3AR0WFY9TZH6DJYDA` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZSSTFBVY` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZTY3TBCS` evidence for `e01M05SPB2EMMC4F4PR0BDQA8S5` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZX049R86` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M067GFJGKHWY2CPHZXV43ZEQ` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Add calm glance tiers and card-only push ref: churn:a18078c4f17e970a8e9c47ef5a26b862cff47798
- `v01M068329RM9V5WNBWXP03ZJW5` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Add validated foreign proposal imports ref: 36350f8fa57bb8b838533665811bd1b42d93a71b
- `v01M068329RM9V5WNBWXPRJDB45` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Add validated foreign proposal imports ref: 36350f8fa57bb8b838533665811bd1b42d93a71b
- `v01M068329RM9V5WNBWXRFKHBSA` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Add validated foreign proposal imports ref: 36350f8fa57bb8b838533665811bd1b42d93a71b
- `v01M068329RM9V5WNBWXTACCYP5` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Add validated foreign proposal imports ref: churn:36350f8fa57bb8b838533665811bd1b42d93a71b
- `v01M068329RM9V5WNBWXVW0DCWZ` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Add validated foreign proposal imports ref: churn:36350f8fa57bb8b838533665811bd1b42d93a71b
- `v01M068329RM9V5WNBWXZ20P4S3` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Add validated foreign proposal imports ref: churn:36350f8fa57bb8b838533665811bd1b42d93a71b
- `v01M068G950YTD7884A946RZK55` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Harden note parsing attribution and dispositions ref: 1aefd6e112625763522d6af08115da7a96a1eaa1
- `v01M068G950YTD7884A949VJ0MD` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Harden note parsing attribution and dispositions ref: 1aefd6e112625763522d6af08115da7a96a1eaa1
- `v01M068G950YTD7884A95VR6D6N` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Harden note parsing attribution and dispositions ref: 1aefd6e112625763522d6af08115da7a96a1eaa1
- `v01M068G950YTD7884A9794NCXT` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Harden note parsing attribution and dispositions ref: churn:1aefd6e112625763522d6af08115da7a96a1eaa1
- `v01M068G950YTD7884A9AZ6MW7D` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Harden note parsing attribution and dispositions ref: churn:1aefd6e112625763522d6af08115da7a96a1eaa1
- `v01M068G950YTD7884A9DEJGNAT` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Harden note parsing attribution and dispositions ref: churn:1aefd6e112625763522d6af08115da7a96a1eaa1
- `v01M0ASAJ10QCE8YZ5CWQ7CT3WW` evidence for `e01M05SG3NTP5W2JX7Y6MG00HQJ` — kind: commit note: 'journal: session record — 4 surface intents + write-path ruling question (held construct) + projection-conflict finding' ref: f8facba76bfbce24f081a15ad55652a571d38670 via: subject-match
- `v01M0ASQWS8HS624DNQAERGH52G` evidence for `e01M068ECYE067WF6BH7F26VC3D` — kind: commit note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)' ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0 via: subject-match
- `v01M0ASQWS8HS624DNQAGGVN1VB` evidence for `e01M065T92NXY1ER6R73YCQNH84` — kind: commit note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)' ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0 via: subject-match
- `v01M0ASQWS8HS624DNQAKBS2VAS` evidence for `e01M05SA72DPRGNTY7GCQ2RASTP` — kind: commit note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)' ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0 via: subject-match
- `v01M0ASQWS8HS624DNQAKR74M0F` evidence for `e01M05RHSWXDNR10P1PY8ERYA9S` — kind: commit note: 'journal: full session record — 3 ratified decisions, 8 findings, 6 ruling questions + raw 788-message transcript (secrets scrubbed)' ref: 20e0d82b9ad522cfa410fcfbf8a3e538aa7ae0e0 via: subject-match
- `v01M0ATYJ301A7XEKWZF4PKCGWA` evidence for `e01M0AQHSRF4DVDYZ989W6K185N` — kind: commit note: 'journal: two owner rulings — multi-parent clew from with strand selection; graphic two-zoom glance' ref: a9615ce309edd28426f1f9cb83f59088c575b3d8 via: subject-match
- `v01M0ATYJ301A7XEKWZF595K1B8` evidence for `e01M0AQHSRF4DVDYZ989M4DVHYX` — kind: commit note: 'journal: two owner rulings — multi-parent clew from with strand selection; graphic two-zoom glance' ref: a9615ce309edd28426f1f9cb83f59088c575b3d8 via: subject-match
- `v01M0AVAQQR6FRF7A3AHG4B1QDN` evidence for `e01M0ASQGK4FMKRNCNR91KJ06JD` — kind: commit note: 'journal: owner ruling — restart-with-mutation is the flagship workflow' ref: 6d3b4f597e5e143524c95e04320266fba797679f via: subject-match
- `v01M0AVWK10B1RJVXDBZ7C7ZF96` evidence for `e01M0ASQGJZPFJD3FSRT7P4HX44` — kind: commit note: 'journal: two owner rulings — I9 replacement, witness-node adoption' ref: f905cc6c64ead7a14520f6245d3a5834c66c4dcf via: subject-match
- `v01M0AVX9FRR8116RPZ937TH2S5` evidence for `e01M0AVXA7EH2KD1BPZ4GJNKN67` — kind: commit note: 'journal: witness-node adoption ruling' ref: d03e76f1d2ebc2213271992191837d541cc33e95 via: subject-match
- `v01M0AXFHWGPXC1GDX69K031AHT` evidence for `e01M0ASQGK2ZK03KRP4WPAQG0MJ` — kind: commit note: 'journal: vendor-neutral freshness ruling + hook survey findings + ladder intent + vocabulary-reduction ruling' ref: cbfc0a87bef5ce45bb8c0e7b223faef0e19bbe29 via: subject-match
- `v01M0AXFHWGPXC1GDX69SZQR91Z` evidence for `e01M0AXF5FM9RDFP817NF7QW5BN` — kind: commit note: 'journal: vendor-neutral freshness ruling + hook survey findings + ladder intent + vocabulary-reduction ruling' ref: cbfc0a87bef5ce45bb8c0e7b223faef0e19bbe29 via: subject-match
- `v01M0AXXTXGQ78NG8YV5ZT8CFA3` evidence for `e01M0AXXKMNNKKY721HJ3REN3KH` — kind: commit note: 'journal: contact-point freshness ruling' ref: a83f0f04e011322aaa7801e563896243c74ec670 via: subject-match
- `v01M0AZ0J80BD0J5P9NNB3EHH4P` evidence for `e01M0AYQ3VJM71SMM1YQFYE1X61` — kind: commit note: 'journal: amnesia test + three-verb amendment to knowledge merge' ref: b6f6ab0cd6baf3ff7aba14dfcc4d2c79fcefca90 via: subject-match
- `v01M0AZ0J80BD0J5P9NNDF90J5S` evidence for `e01M0AY6VWJV133F811JYAANJPE` — kind: commit note: 'journal: amnesia test + three-verb amendment to knowledge merge' ref: b6f6ab0cd6baf3ff7aba14dfcc4d2c79fcefca90 via: subject-match
- `v01M0AZ7HW8KGA49R5THPW8KP8D` evidence for `e01M0AZ7JBPJRNHFJQC6WEEQB9Y` — kind: commit note: 'journal: silence-is-the-signal ruling' ref: bde49235a5a6fb0e11473a09e90a5951dd8018ea via: subject-match
- `v01M0AZMWMGD6SD3MGT9CDJDNQ1` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9DCW8YG7` evidence for `e01M05V5HWA6TFT0A0KDZY8S45K` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9H2WPVAG` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9NA4WEE7` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9NXC1RBW` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9QMV58RB` evidence for `e01M064XJ84ZTY2HNFWWQZATFA9` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9RMJWZB6` evidence for `e01M05SPB2EMMC4F4PR09QWMXG8` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGT9WGYXXKJ` evidence for `e01M05SPB2EMMC4F4PR0BDQA8S5` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA053QK3R` evidence for `e01M064XJ84ZTY2HNFWWV8RSRQR` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA0M6H25Z` evidence for `e01M064YRS4S9NK7KW9NN3114KH` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA1SE512J` evidence for `e01M064YRS4S9NK7KW9NQ1JMDV2` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA22CXAB9` evidence for `e01M05T52RRS9EXQ92M3PA9WBVF` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA3V187AT` evidence for `e01M05TCM2N1P9P2EQSE8YKVQX3` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA5TCMHHV` evidence for `e01M05VTCM3AR0WFY9TZH6DJYDA` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTA76GCDX7` evidence for `e01M05VTCM3AR0WFY9TZJXAN7SE` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTAAKQHGD1` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTADBEPG6F` evidence for `e01M05VFAW9A783PMZZEBHJ1TZH` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTAF6M4544` evidence for `e01M05VNN1T5TKJ15Q47KSS3FMJ` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTAFP3G8WH` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGD6SD3MGTAH5Y1001` evidence for `e01M05VTCM3AR0WFY9TZPG9W1J8` — kind: churn note: Ship I13 owner memory and explicit lineage ref: churn:5a56835ff13911a868bda153456b50e26b785574
- `v01M0AZMWMGEBE16SP0P9PNT451` evidence for `e01M0AV0H7T9P69CNPB56MRAG8V` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PBHC23RN` evidence for `e01M0AVXA7EH2KD1BPZ4GJNKN67` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PEJTYZPB` evidence for `e01M0AVAR0N8ZCCN1VJW3GQ5PF4` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PHT83B2D` evidence for `e01M0AQHSRF4DVDYZ989W6K185N` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PNMJ6RXN` evidence for `e01M0ATYJG615JE6BV5MG5RAF9Z` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PS13M6MD` evidence for `e01M0ASQGK4FMKRNCNR91KJ06JD` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PV79J7K0` evidence for `e01M0AQHSRF4DVDYZ989M4DVHYX` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZMWMGEBE16SP0PY76EGPR` evidence for `e01M0AQHSRF4DVDYZ989RNZVC46` — kind: commit note: Ship I13 owner memory and explicit lineage ref: 5a56835ff13911a868bda153456b50e26b785574 via: subject-match
- `v01M0AZS3D0R9F6JDW9DENN65PT` evidence for `e01M0AQHSRF4DVDYZ989RNZVC46` — kind: commit note: 'journal: restart tab hard half — what not to carry' ref: b9b17a6c6be36e3a76f1dd2a7de7f2d039ec84c8 via: subject-match
- `v01M0AZS3D0R9F6JDW9DGCSENKS` evidence for `e01M0AVAR0N8ZCCN1VJW3GQ5PF4` — kind: commit note: 'journal: restart tab hard half — what not to carry' ref: b9b17a6c6be36e3a76f1dd2a7de7f2d039ec84c8 via: subject-match
- `v01M0AZS3D0R9F6JDW9DN3M30D3` evidence for `e01M0AQHSRF4DVDYZ989MTZQV7E` — kind: commit note: 'journal: restart tab hard half — what not to carry' ref: b9b17a6c6be36e3a76f1dd2a7de7f2d039ec84c8 via: subject-match
- `v01M0B09808B4RKCT5KCXRXXCDK` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Add finish knowledge and intent screens ref: 1c9f2d28293be323727342c7bdb4ecad5646969e
- `v01M0B09808B4RKCT5KD1W0DPB5` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Add finish knowledge and intent screens ref: churn:1c9f2d28293be323727342c7bdb4ecad5646969e
- `v01M0B09808B4RKCT5KD1YWD5AJ` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Add finish knowledge and intent screens ref: churn:1c9f2d28293be323727342c7bdb4ecad5646969e
- `v01M0B09808B4RKCT5KD2J9WMDA` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Add finish knowledge and intent screens ref: churn:1c9f2d28293be323727342c7bdb4ecad5646969e
- `v01M0B09808FFQSS36NCHT943RR` evidence for `e01M0AY6VWJV133F811JYAANJPE` — kind: commit note: Add finish knowledge and intent screens ref: 1c9f2d28293be323727342c7bdb4ecad5646969e via: subject-match
- `v01M0B09808FFQSS36NCKNDV2ES` evidence for `e01M0AYQ3VJM71SMM1YQFYE1X61` — kind: commit note: Add finish knowledge and intent screens ref: 1c9f2d28293be323727342c7bdb4ecad5646969e via: subject-match
- `v01M0BER1JG15ZMPFYZF9KEE5Y7` evidence for `e01M0AXCM53D36PDN0H85S0WE3W` — kind: commit note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note' ref: 03453909c924fbd7d3d893642b3510654adc105a via: subject-match
- `v01M0BER1JG15ZMPFYZFD5YW7KF` evidence for `e01M0BEM04PAZ85R5YNRM2Y31Z6` — kind: commit note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note' ref: 03453909c924fbd7d3d893642b3510654adc105a via: subject-match
- `v01M0BER1JG15ZMPFYZFDWCMXK9` evidence for `e01M0AV0H7RY7DG79VCMHPEMPJP` — kind: commit note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note' ref: 03453909c924fbd7d3d893642b3510654adc105a via: subject-match
- `v01M0BER1JG15ZMPFYZFF1WW5Y0` evidence for `e01M05SA72DPRGNTY7GCN1P7CED` — kind: commit note: 'journal: plain-speech/no-ids ruling, 7+7 cap question, spoken-verbs interim note' ref: 03453909c924fbd7d3d893642b3510654adc105a via: subject-match
- `v01M0BEYNFR7B7XHM8Z8W9YJN02` evidence for `e01M0BETRRAZ54063PRJK1JSQS7` — kind: commit note: 'journal: finished-means-shared ruling' ref: 9df0522c4916bc1857e68e89a07f5b74d1239423 via: subject-match
- `v01M0BEYNFR7B7XHM8Z8XEGBZ4W` evidence for `e01M0AXZTTHG5FKXETX0X8PR6EX` — kind: commit note: 'journal: finished-means-shared ruling' ref: 9df0522c4916bc1857e68e89a07f5b74d1239423 via: subject-match
- `v01M0BEYNFR7B7XHM8Z8XTCCPRB` evidence for `e01M0BEYP65CE70G0VVSX3PV01B` — kind: commit note: 'journal: finished-means-shared ruling' ref: 9df0522c4916bc1857e68e89a07f5b74d1239423 via: subject-match
- `v01M0BF0VSRSCB9JV6YNSVKY381` evidence for `e01M0BER1Q412DRFQYESPCN0Q30` — kind: commit note: 'journal: ids-are-plumbing ruling' ref: d3c122da0c3be5edc4d8e2fa479eacced982986a via: subject-match
- `v01M0BF0VSRSCB9JV6YNTSSJK50` evidence for `e01M0BF0WMX264RA9D0VTM9R24K` — kind: commit note: 'journal: ids-are-plumbing ruling' ref: d3c122da0c3be5edc4d8e2fa479eacced982986a via: subject-match
- `v01M0BFFY78FSYSC382CJNMEFDC` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Make finish screens plain and spoken ref: b734c7e16971a7c67a901e5467e651a774eb313a
- `v01M0BFFY78FSYSC382CMKFW758` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Make finish screens plain and spoken ref: churn:b734c7e16971a7c67a901e5467e651a774eb313a
- `v01M0BFFY78FSYSC382CPVF31PZ` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Make finish screens plain and spoken ref: churn:b734c7e16971a7c67a901e5467e651a774eb313a
- `v01M0BFFY78FSYSC382CS5SKNBY` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Make finish screens plain and spoken ref: churn:b734c7e16971a7c67a901e5467e651a774eb313a
- `v01M0BFRYA06CNTJQN7ASDHEJM8` evidence for `e01M0AZ4BGF1VC0VSXA05VYVEQ3` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7AVXQ4JGY` evidence for `e01M0AZDDYC49Y1741H7W74Y1QY` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7AY81FBG8` evidence for `e01M0AV0H7T9P69CNPB56MRAG8V` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7AZ9J0VAD` evidence for `e01M0ATYJG615JE6BV5MG5RAF9Z` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7B2PEPA8R` evidence for `e01M0BFRYY2WVHJZ3R3TDV4CFTS` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7B9GXHMV7` evidence for `e01M0AXMXK3SAATFDFAYZ932TPC` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7BAEP60CE` evidence for `e01M0AZ0K17AA5D9P7KZDSPJQSY` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFRYA06CNTJQN7BBRE5XX5` evidence for `e01M0BFRYY1EW9CFMJNDJ5QH3M2` — kind: commit note: 'journal: evidence-settles-merge ruling + full wording-sweep scope' ref: a4992570440dee46a2794892cabb9d82e52d6522 via: subject-match
- `v01M0BFY0DGTZDK68V0F02X5J9C` evidence for `e01M0BFY164YEV6DAEFVGGH18VT` — kind: commit note: 'journal: lag-never-deaf limiter ruling' ref: 71bed05e27fcd9a52546a52882559cc84a250326 via: subject-match
- `v01M0BFYC4G9WTCPMV7SH2R3A6K` evidence for `e01M0AVWK3HH2R9M55FQSAWZFF1` — kind: commit note: 'journal: 2 file(s)' ref: d3f180a93e7beaf8e219e510ea4ca73d829b1055 via: subject-match
- `v01M0BJ206G6547MERH7NSAQN0P` evidence for `e01M05SA72DPRGNTY7GCPEX9W2N` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH7SA3NCZ7` evidence for `e01M05V5HWA6TFT0A0KDZY8S45K` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH7WY8PVEP` evidence for `e01M05SA72DPRGNTY7GD0TGHEX7` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH7ZR3EFAZ` evidence for `e01M04H0KBY8QQWXPE8DP99N012` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH81NP2GNY` evidence for `e01M05SPB2EMMC4F4PR0928NA31` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8511HF8E` evidence for `e01M05VTCM3AR0WFY9TZPG9W1J8` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH88HWFARN` evidence for `e01M05TCM2N1P9P2EQSE8YKVQX3` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH89NSAAYG` evidence for `e01M064YRS4S9NK7KW9NN3114KH` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8CFCPQV1` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8GAHM82W` evidence for `e01M064XJ84ZTY2HNFWWV8RSRQR` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8H4P4GP6` evidence for `e01M05SPB2EMMC4F4PR0BDQA8S5` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8JS9PXSK` evidence for `e01M05VNN1T5TKJ15Q47KSS3FMJ` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8P2450DQ` evidence for `e01M05T52RRS9EXQ92M3PA9WBVF` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8S5FQ2DS` evidence for `e01M05VFAW9A783PMZZEBHJ1TZH` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8WXM56RC` evidence for `e01M05VTCM3AR0WFY9TZJXAN7SE` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH8ZF8PTQS` evidence for `e01M05VTCM3AR0WFY9TZH6DJYDA` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH90HSHAMS` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH935PES1Z` evidence for `e01M064XJ84ZTY2HNFWWQZATFA9` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH961H4MZ7` evidence for `e01M05SPB2EMMC4F4PR09QWMXG8` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206G6547MERH9708PZF0` evidence for `e01M064YRS4S9NK7KW9NQ1JMDV2` — kind: churn note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: churn:e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5
- `v01M0BJ206GVJFV3HJ3NBYMEWTX` evidence for `e01M0AZDDYC49Y1741H7W74Y1QY` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NEMWDEPW` evidence for `e01M0BER1Q412DRFQYESPCN0Q30` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NGDTBWQQ` evidence for `e01M0BEM04PAZ85R5YNRM2Y31Z6` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NKCHACR4` evidence for `e01M0BETRRAZ54063PRJK1JSQS7` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NMBMS9K6` evidence for `e01M0BEYP65CE70G0VVSX3PV01B` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NPT8K21X` evidence for `e01M0BFRYY2WVHJZ3R3TDV4CFTS` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NS6QXAYK` evidence for `e01M0BF0WMX264RA9D0VTM9R24K` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NW8ZJ37W` evidence for `e01M0BFY164YEV6DAEFVGGH18VT` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ206GVJFV3HJ3NWV659ZX` evidence for `e01M0AXMXK3SAATFDFAYZ932TPC` — kind: commit note: Record lag; settle evidence; sweep wording; two registers; finish shared; plain speech; fix broken verbs; intent gap ref: e27073c609faf9ffc9bc2b0f5aacb3e67c4feec5 via: subject-match
- `v01M0BJ2348QV4VY6WYW7FYX5Z4` evidence for `e01M064DQTWYDVGGAE3M5QRTGME` — kind: commit note: 'journal: 24 file(s)' ref: 2f14008720ee908109a7715a0d82cc4c9dea9251 via: subject-match
- `v01M0BJ2348QV4VY6WYW8KYJ911` evidence for `e01M0AXRKQ8C7FZNKARW83CMBMX` — kind: commit note: 'journal: 24 file(s)' ref: 2f14008720ee908109a7715a0d82cc4c9dea9251 via: subject-match
- `v01M0BJ2348QV4VY6WYWC1S18MX` evidence for `e01M05TBYJXEW5N5FE397XYMEHY` — kind: commit note: 'journal: 1 file(s)' ref: 9f25ad55657388cdd4daa160f5969aaf07ee04e3 via: subject-match
- `v01M0BJ2348QV4VY6WYWF9Z0SQM` evidence for `e01M04WCGJS9FS7FQB0YFX9DTYG` — kind: commit note: 'journal: 1 file(s)' ref: 9f25ad55657388cdd4daa160f5969aaf07ee04e3 via: subject-match
- `v01M0BJ84FR90G2GJR21QRS1XYC` evidence for `e01M05T4XG0RJWQTP25SYT4FH0B` — kind: churn note: Ignore journal commits when settling evidence ref: churn:f66576bb0ec5cc0f2bea61ed1297c0920761429d
- `v01M0BJ84FR90G2GJR21SRNJH73` evidence for `e01M05T4XG0RJWQTP25T1K58B61` — kind: churn note: Ignore journal commits when settling evidence ref: churn:f66576bb0ec5cc0f2bea61ed1297c0920761429d
- `v01M0BJ84FR90G2GJR21TNAXGGP` evidence for `e01M0BFRYY1EW9CFMJNDJ5QH3M2` — kind: commit note: Ignore journal commits when settling evidence ref: f66576bb0ec5cc0f2bea61ed1297c0920761429d via: subject-match

## Organ-bank pin

- `https://github.com/maceip/clew.git` at `f66576bb0ec5cc0f2bea61ed1297c0920761429d`
