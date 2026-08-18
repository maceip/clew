# clew

Agents regenerate code cheaply; nothing regenerates the knowledge a project
earns — decisions with reasons, measured findings, questions never answered,
work never started. `clew` watches the sessions and repos you already
have, distills that knowledge into a durable journal on an orphan git branch,
joins it against observed reality, and shows both customers — the agent and
the human — where intent and reality diverge, including the work that was
never built.

This is the faithful implementation of [`JOURNAL_SPEC.md`](JOURNAL_SPEC.md)
(the founding spec; invariants I1–I13 in §1 are the constitution of this
repository). The spec's working name `stratura` was a placeholder to rename
before first tag; the owner named the system **clew** (§12.1 is closed — the
decision is journaled). Making restarts lossless remains the co-lead
feature (§9).

## Build & install

```bash
# install onto your PATH (single static binary; Go ≥ 1.22):
GOBIN="$HOME/.local/bin" go install ./cmd/clew   # any dir on your PATH works
clew help

# plain `go install ./cmd/clew` uses $(go env GOPATH)/bin — put that on PATH
# first, or you'll hit "command not found". For a local build: go build -o clew ./cmd/clew
```

## Quickstart

```bash
clew watch install      # once per machine: watcher + Claude SessionStart birth hook

# A new repo is born without a clew ceremony:
mkdir your-project && cd your-project && git init
claude                  # auto-registers; first context has owner laws, no lore

# Existing/ancient repos may still request archaeology explicitly:
clew init               # register; archaeology; install agent includes + nudge hook
clew                    # calm glance: DECIDED·LEARNED·OPEN·MAP·DOCKET
clew status             # diagnostic expansion
clew glance --html      # ~/.clew/glance.html, 30s refresh, pinned-tab status light
clew map                # intent × reality, absence highlighted
clew docket             # bounded decision cards only (`inbox` remains an alias)
clew import bundle.yaml # one validated foreign proposal batch; never direct injection
clew backfill --budget 100000   # months of existing session history → journal
```

The watcher tails Claude Code, Codex, and Cursor CLI session files (formats
pinned to observed versions — unknown line classes are parked and counted,
never guess-parsed), polls registered repos, extracts journal entries with
your own agent subscription (`claude -p` → `codex exec` → OpenAI-compatible,
first available; configurable in `~/.clew/config.yaml`), and syncs the
journal branch through the repo's own remote. It also keeps `.clew/SEED.md`
current whenever the durable journal changes and synchronizes the owner's
separate project-agnostic law journal. No server, no new credentials.
Optional phone delivery uses a unique ntfy topic (or plain webhook) configured in
`~/.clew/config.yaml`; only newly created docket cards push, with headline + why-you.

For agents with no session store: `clew wrap -- gemini …` (PTY tee).
Optional MCP surface: `clew mcp` (journal_search / journal_get /
journal_note).

Claude installations relocated with `CLAUDE_CONFIG_DIR` are honored by both
the hook installer and transcript discovery; the supervisor persists that
directory in its environment.

## Ambient birth and declared lineage (I13)

Birth carries only laws that a human has certified as project-agnostic:

```bash
clew journal promote <finding-id>   # certify its exact quoted evidence in owner scope
```

The owner-law block is injected before project lore in every `context.md`, is
independently capped at 1 KiB, and remains inside the existing 4 KiB total
context cap. The extractor can mark a project-agnostic finding for a human
docket ruling, but that proposal/card does not promote anything and never
enters an agent surface; the finding itself remains ordinary project-local
memory. The ambient law text is the exact Quote the docket showed, never the
extractor-authored Title or Body.
When `owner.remote` is configured, promotion requires a current remote budget
view; an offline fetch or deferred push refuses certification instead of
guessing that the global 1 KiB set still has room.

Project lore has the opposite boundary: lineage is never guessed or carried
automatically. Every watched repo already has a digest-checked ambient seed
containing its decisions, findings, graveyard, exhibits, and organ-bank pin.

```bash
clew from              # ranked by topic overlap + recency; read-only
clew from ../substrate # explicit import + durable lineage link
```

Carried entries retain their original quote, source ref, timestamp, and event
provenance. Graveyard entries also receive a durable human lineage-status
marker: derived absence, expiry, disposition-drop, and supersession cannot
accidentally become live merely because their source derivation was outside
the seed. New successor evidence can still revive an absent intent. To
un-carry one, reject it normally with
`clew journal reject <entry-id>`; the rejection is append-only and a repeated
`clew from` cannot resurrect it. A blatant name/topic match may be suggested
to the human, but is never acted on.

## The deliberate big-restart ceremony (§9)

```bash
clew manifest                    # explicit pull-only disposition pass → MANIFEST.md
$EDITOR .../MANIFEST.md             # human pass: mark [carry] / [drop]
clew manifest --out ./kit       # SEED.md (≤4KB) + genesis/ + journaled dispositions
cd ../successor && git init && clew init --carry ../old/kit/genesis
```

The manifest path remains useful when a human wants to disposition active work
and assemble a larger restart kit. It is never a prerequisite for ambient
`SEED.md`, costless birth, or `clew from`.

## Storage model (§4)

- Orphan branch `clew/journal` in the project's own remote; worktree
  checkout under `~/.clew/worktrees/<id>`. Reusing a path for a fresh
  `git init` allocates a persistent new incarnation for both the worktree and
  the newborn's local-only seed identity, leaving the moved predecessor
  available only through explicit lineage.
- `entries/<ulid>.yaml` + `events/<ulid>.yaml` + explicit
  `lineage/<ulid>.yaml`, all immutable once written;
  every writer only ever adds files → conflict-free by construction.
  `journal.md`/`digest.md`/`SEED.md` are regenerated projections (races
  harmless; seed bytes change only with the durable journal revision, which
  includes every canonical lineage link ID and selected seed digest).
- Owner laws live in a separate normal git repository at `~/.clew/owner`
  (optional `owner.remote`), never in the project registry or session scanner.
- Statuses are computed by each watcher from entries+events and **never
  persisted** (§3.2) — that is what makes the conflict-freedom claim true.
- `clew redact <id>` is the one sanctioned rewrite: scrub-in-place, fresh
  root, and force-with-lease against the exact pre-redaction sync tip. A
  concurrent append causes a bounded resync/re-scrub/retry, never erasure of
  unseen journal data. Other watchers adopt the accepted new root and cannot
  resurrect the secret (their re-add pass only restores files the remote
  *lacks*). If the entry was promoted, the owner journal is rewritten first
  and every reachable context is refreshed without that law.
- Machine-local working state in SQLite (`~/.clew/state.db`); never the
  source of truth.

### Proposing journal entries without journal credentials

`clew import <bundle.yaml|dir|https-url>` validates entry schemas and quotes, marks provenance
as foreign, and creates one owner docket card. The owner uses `clew docket open <proposal-id>`,
then `accept` or `reject`; accept is the human-confirm boundary.

PAT-authorized contributors can use the sanctioned fork/PR route:

```bash
CLEW_PROPOSAL_BRANCH=clew/proposal-my-change clew import bundle.yaml
# open a PR from clew/proposal-my-change with base branch clew/journal
```

The branch is based on `clew/journal`; merging that PR confirms the batch. Credentials push,
everyone else proposes.

## Acceptance (§10)

```bash
go test ./...                        # all five §10 gates across domain + command integration tests
CLEW_FIDELITY=1 go test ./acceptance -run RealProvider -v   # live extraction gate
```

1. **Absence detection** — synthetic agentdesk-shaped fixture: real git
   commits evidence 5 eligible sibling intents; the workload-runner intent
   flips to `absent` exactly then (sibling rule, not wall-clock), and lands
   in the docket.
2. **Extraction fidelity** — `fixtures/strategy-session/` carries the
   Appendix-A label set (D1–D7, F1–F5) with `ratified: false`; per §10.2 a
   human must review/amend once and set `ratified: true`, after which the
   live gate (precision ≥ 0.9, recall ≥ 0.75, ≤ 5 instruction iterations)
   is the written-down kill criterion. The hermetic run proves the pipeline
   and the I7 zero-fabrication gate.
3. **Restart round-trip** — 12 live entries → manifest (2 dropped, journaled)
   → `init --carry` into an empty repo → 100% of carried decisions/findings
   in the new `context.md`, provenance and timestamps intact.
4. **Costless birth** — in an isolated machine home, a certified owner law
   followed by `mkdir x && git init && claude`'s SessionStart hook produces a
   registered project and a first `context.md` containing that law, with zero
   project lore and zero clew commands typed for the project.
5. **Declared lineage** — ambient seeds round-trip with digest and provenance;
   ranking is deterministic; an explicit import carries lessons, graveyard,
   exhibits, and pin, records transitive ancestry, rejects cycles, and cannot
   resurrect an entry the successor rejected.

## Spec review — resolutions baked into this implementation

Ambiguities/underdeterminations found while implementing, and how they were
resolved (all within the spec's own invariants):

| # | Spec point | Resolution |
|---|---|---|
| 1 | §12.1 name before first tag | **`clew`** — owner's decision, superseding the builder's unilateral `restart` (journaled: `e01M04WCGJS9FS7FQB0YFX9DTYG` supersedes `e01M04H0KBY8QQWXPE8DP99N012`). Alternatives considered: `restart` (verb collision; names the crisis, not the daily loop), `lore` (binary/brand collision: varalys/lore, getlore.ai, Epic Lore), `wake`, `canon`, `lorekeeper`. Binary `clew`, branch `clew/journal`, `.clew/`, `~/.clew/` |
| 2 | §3.1 questions "carry who-can-answer" but §3.2's example (an intent) shows no field | added `asks: human\|any` on question entries |
| 3 | §7.1.3 intent "completion confirmation event" unnamed | `clew journal confirm <id> --done` |
| 4 | §7.1.3 "only a human confirm event sets contradicted" — verb unspecified | `clew journal confirm <id> --contradicts <other-id>` |
| 5 | docket `ack`/`drop` transport unspecified | persisted as `disposition` events on the branch so rulings propagate to every machine's docket |
| 6 | intent with evidence older than 7d, not done: no status enumerated | renders `proposed` (map shows evidence count, so it stays legible) |
| 7 | §12.3 cursor desktop store depth | v1 parses Cursor **CLI** transcripts (format pinned from live files); the desktop `state.vscdb` is detected and flagged "not parsed — lower fidelity" per I2, never guessed |
| 8 | §6.3 "default: cheapest configured provider" — cheapness unknowable | `auto` = first available of claude → codex → OpenAI-compatible; explicit `provider:` overrides |
| 9 | §4.1 relay | not built — §11's v1 core list excludes it; baseline git sync is the complete mode |
| 10 | §8.2 `answer`/`note` entry type | defaults to `finding`, `--type decision…` overrides |
| 11 | wrap-mode utterance attribution | stdin → `user`, terminal output → `assistant`; the extractor may still classify third-party content as `tool_result`, and the I7 verifier force-taints quotes found in tool output |
| 12 | §5.3 archaeology "distill" without a provider | mechanical passes (TODO/FIXME, ADR files, `gh` PR titles) always run; README/docs/commit distillation requires a provider and is loudly skipped otherwise |
| 13 | §9.2 manifest human pass | `[carry]`/`[drop]` marks in `MANIFEST.md`; `--yes` for non-interactive carry-all; spec-`covered` entries default to `[drop]` |
| 14 | journal `edit` vs immutable entries (§3.2) | `edit` opens a copy in `$EDITOR`, writes a **new** human entry + supersede event — append-only law preserved |
| 15 | §6 numbering (6.5 sits between 6.2 and 6.3) | noted; no action needed |
| 16 | I13 birth safety | owner laws are a separate certified journal layer (1 KiB); ambient project seeds are watcher-maintained; project lore crosses repos only through explicit `clew from` because a wrong lineage guess is more poisonous than no inheritance |

Everything else is as written: 4 entry types with computed statuses, the
absence rule with K=5 eligible siblings and the archaeology eligibility
guard, env-scoped finding supersession, 45d question expiry, 7d aging
alerts, 100KB slice cap with middle-elision, 4KB digest, 4KB context with
the fixed injection preamble + tainted-quote data fences + imperative
withholding, owner-law certification and its 1 KiB sub-cap, 2%/200k budget
with loud pause, the delivery matrix (Claude
hook / human-routed / PTY injection), ≤7 lines per glance section, 32KB
rollup over-fire flag, and the honesty lines in `status` (cursor desktop
fidelity, cloud-session gap).

## Layout

```
cmd/clew/          CLI: init watch status map docket journal from manifest
                      backfill wrap redact mcp
internal/model        entry/event schema + validation (§3)
internal/journal      store, status algebra, rollup/digest (§3, §7.1.3)
internal/gitx         orphan branch, sync, adoption, redact rewrite (§4)
internal/adapters     claude/codex/cursor/wrap sensors + adapter law (§5.1)
internal/poller       repo poller + commit attribution (§5.2)
internal/extract      triggers, provider call, schema+I7 gate, budget (§6)
internal/differ       mapping, auto-supersession, alerts, overlap radar (§7)
internal/materialize  context.md/SEED.md/nudge.md/journal.md (§8.1, I13)
internal/owner        certified project-agnostic owner journal + 1 KiB renderer
internal/seed         self-contained ambient seed, journal-revision gated
internal/lineage      ranking, provenance-preserving import, durable links
internal/manifest     MANIFEST/SEED/genesis/carry (§9)
internal/archaeology  cold-start distillation (§5.3)
internal/{llm,mcp,scrub,state,push,wrapx,globx,ids,config}
fixtures/             Appendix-A ground truth (awaiting human ratification)
acceptance/           original absence/extraction/restart gates; I13 gates live by their seams
```
