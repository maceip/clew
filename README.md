# clew

Agents regenerate code cheaply; nothing regenerates the knowledge a project
earns — decisions with reasons, measured findings, questions never answered,
work never started. `clew` watches the sessions and repos you already
have, distills that knowledge into a durable journal on an orphan git branch,
joins it against observed reality, and shows both customers — the agent and
the human — where intent and reality diverge, including the work that was
never built.

This is the faithful implementation of [`JOURNAL_SPEC.md`](JOURNAL_SPEC.md)
(the founding spec; invariants I1–I9 in §1 are the constitution of this
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
cd your-project            # any git repo (empty or ancient)
clew init               # register; archaeology; installs CLAUDE.md/AGENTS.md
                           #   include + Claude nudge hook + .gitignore entry
clew watch install      # launchd/systemd-user supervised watcher
clew status             # the glance: SESSIONS·DECIDED·LEARNED·OPEN·ALERTS
clew map                # intent × reality, absence highlighted
clew inbox              # human-blocking items only
clew backfill --budget 100000   # months of existing session history → journal
```

The watcher tails Claude Code, Codex, and Cursor CLI session files (formats
pinned to observed versions — unknown line classes are parked and counted,
never guess-parsed), polls registered repos, extracts journal entries with
your own agent subscription (`claude -p` → `codex exec` → OpenAI-compatible,
first available; configurable in `~/.clew/config.yaml`), and syncs the
journal branch through the repo's own remote. No server, no new credentials.

For agents with no session store: `clew wrap -- gemini …` (PTY tee).
Optional MCP surface: `clew mcp` (journal_search / journal_get /
journal_note).

## The restart moment (§9)

```bash
clew manifest                    # disposition pass → MANIFEST.md
$EDITOR .../MANIFEST.md             # human pass: mark [carry] / [drop]
clew manifest --out ./kit       # SEED.md (≤4KB) + genesis/ + journaled dispositions
cd ../successor && git init && clew init --carry ../old/kit/genesis
```

The successor's first `context.md` already contains the predecessor's earned
knowledge, provenance intact — restart becomes a compile from source instead
of amnesia.

## Storage model (§4)

- Orphan branch `clew/journal` in the project's own remote; worktree
  checkout under `~/.clew/worktrees/<id>`.
- `entries/<ulid>.yaml` + `events/<ulid>.yaml`, both immutable once written;
  every writer only ever adds files → conflict-free by construction.
  `journal.md`/`digest.md` are regenerated projections (races harmless).
- Statuses are computed by each watcher from entries+events and **never
  persisted** (§3.2) — that is what makes the conflict-freedom claim true.
- `clew redact <id>` is the one sanctioned rewrite: scrub-in-place, fresh
  root, force-push; other watchers adopt the new root and cannot resurrect
  the secret (their re-add pass only restores files the remote *lacks*).
- Machine-local working state in SQLite (`~/.clew/state.db`); never the
  source of truth.

## Acceptance (§10)

```bash
go test ./...                        # includes the three §10 gates (hermetic)
CLEW_FIDELITY=1 go test ./acceptance -run RealProvider -v   # live extraction gate
```

1. **Absence detection** — synthetic agentdesk-shaped fixture: real git
   commits evidence 5 eligible sibling intents; the workload-runner intent
   flips to `absent` exactly then (sibling rule, not wall-clock), and lands
   in the inbox.
2. **Extraction fidelity** — `fixtures/strategy-session/` carries the
   Appendix-A label set (D1–D7, F1–F5) with `ratified: false`; per §10.2 a
   human must review/amend once and set `ratified: true`, after which the
   live gate (precision ≥ 0.9, recall ≥ 0.75, ≤ 5 instruction iterations)
   is the written-down kill criterion. The hermetic run proves the pipeline
   and the I7 zero-fabrication gate.
3. **Restart round-trip** — 12 live entries → manifest (2 dropped, journaled)
   → `init --carry` into an empty repo → 100% of carried decisions/findings
   in the new `context.md`, provenance and timestamps intact.

## Spec review — resolutions baked into this implementation

Ambiguities/underdeterminations found while implementing, and how they were
resolved (all within the spec's own invariants):

| # | Spec point | Resolution |
|---|---|---|
| 1 | §12.1 name before first tag | **`clew`** — owner's decision, superseding the builder's unilateral `restart` (journaled: `e01M04WCGJS9FS7FQB0YFX9DTYG` supersedes `e01M04H0KBY8QQWXPE8DP99N012`). Alternatives considered: `restart` (verb collision; names the crisis, not the daily loop), `lore` (binary/brand collision: varalys/lore, getlore.ai, Epic Lore), `wake`, `canon`, `lorekeeper`. Binary `clew`, branch `clew/journal`, `.clew/`, `~/.clew/` |
| 2 | §3.1 questions "carry who-can-answer" but §3.2's example (an intent) shows no field | added `asks: human\|any` on question entries |
| 3 | §7.1.3 intent "completion confirmation event" unnamed | `clew journal confirm <id> --done` |
| 4 | §7.1.3 "only a human confirm event sets contradicted" — verb unspecified | `clew journal confirm <id> --contradicts <other-id>` |
| 5 | inbox `ack`/`drop` transport unspecified | persisted as `disposition` events on the branch so acks propagate to every machine's inbox |
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

Everything else is as written: 4 entry types with computed statuses, the
absence rule with K=5 eligible siblings and the archaeology eligibility
guard, env-scoped finding supersession, 45d question expiry, 7d aging
alerts, 100KB slice cap with middle-elision, 4KB digest, 4KB context with
the fixed injection preamble + tainted-quote data fences + imperative
withholding, 2%/200k budget with loud pause, the delivery matrix (Claude
hook / human-routed / PTY injection), ≤7 lines per glance section, 32KB
rollup over-fire flag, and the honesty lines in `status` (cursor desktop
fidelity, cloud-session gap).

## Layout

```
cmd/clew/          CLI: init watch status map inbox journal manifest
                      backfill wrap redact mcp
internal/model        entry/event schema + validation (§3)
internal/journal      store, status algebra, rollup/digest (§3, §7.1.3)
internal/gitx         orphan branch, sync, adoption, redact rewrite (§4)
internal/adapters     claude/codex/cursor/wrap sensors + adapter law (§5.1)
internal/poller       repo poller + commit attribution (§5.2)
internal/extract      triggers, provider call, schema+I7 gate, budget (§6)
internal/differ       mapping, auto-supersession, alerts, overlap radar (§7)
internal/materialize  context.md/nudge.md/journal.md (§8.1)
internal/manifest     MANIFEST/SEED/genesis/carry (§9)
internal/archaeology  cold-start distillation (§5.3)
internal/{llm,mcp,scrub,state,push,wrapx,globx,ids,config}
fixtures/             Appendix-A ground truth (awaiting human ratification)
acceptance/           the three §10 gates
```
