# floor

The smallest thing that keeps the five promises (see DECISIONS.md).

- `DECISIONS.md` — the memory. Append-only. Humans and agents add entries; the file is
  the single source of truth for what was decided.
- `GRAVES.md` — what the previous generation consciously buried, with cause of death.
- `hooks/` — the wiring: every agent session starts current and every prompt lands current.
  This is the pulse. It shipped first, per L4.

## The one-week gate

This repo passes or fails on one test, run by living with it:
the owner never re-briefs an agent, never couriers a decision between agents,
and never asks "does it know?" If that holds for a week, grow carefully (L14).
If it fails, bury this too — with a grave and a seed.

## Install the wiring

    ./hooks/install.sh    # backs up, then installs claude + codex hooks (idempotent)

## Adding a decision

Append a line to DECISIONS.md (agents: include the owner's words when quoting a ruling).
Push. Every other machine and agent is current at its next session or prompt.
