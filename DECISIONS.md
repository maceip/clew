# DECISIONS

This file is the project's memory. Agents: read it fully before planning; treat entries
as constraints unless new evidence contradicts one — then say so explicitly. Humans and
agents append; nobody edits or deletes (strike through with a dated note instead).

## The product (ratified 2026-08-18)

Five promises. Every feature files under exactly one, or it is not built.

1. It remembers what we decide.
2. Every agent starts already knowing it.
3. You can look up and see.
4. It taps your shoulder only when something needs you.
5. Starting over loses nothing.

> "yeah those 5 things are the foundation, agreed"

## Laws

**L1 — A direct owner order is never refused.** Rules and defaults govern autonomy only;
the freshest explicit order wins over everything, including older orders. State collisions
in one line while obeying. Refusing a direct order is a product-ending failure.
> "when i tell the agent to do something 'build all' and it doesnt do it, then im angry and were done"

**L2 — Tokens are not scarce. Owner attention is the only scarce resource.** No rule,
meter, ceiling, floor, or pause may exist to conserve spend. The only limiter is loop
detection (runaway protection), never budget.
> "i keep telling you we have an unlimited token budget and not to be frugal"

**L3 — Every gate must cite the owner's words.** Any rule that withholds, delays, or
gates work carries a verbatim owner quote authorizing it, or it is invalid by construction.
> "when did i authorize or approve _not doing work_ because it was oversized?"

**L4 — Wire the pulse before the organs.** The load-bearing loop ships first and gets the
acceptance gate; features queue behind it. Gates test the promises, not the machinery.
> "so lets build the one piece thats missing then? why is that not top of the list of things to do ?"

**L5 — Verify load-bearing facts the day you build on them.** The parent died partly of a
spec that froze "codex has no hooks" — false, and falsifiable with one search at write time.
> "the simple (and easy via one google search 'does codex have hooks')"

**L6 — Auto-absorb knowledge; auto-build intents; humans see genuine conflicts only.**
Absorption is the default. A human is interrupted only by a conflict the system cannot
resolve — never by a queue, a ceremony, or a consent ritual.
> "it should auto-merge all incoming knoweldge items and build all pending intent gaps"

**L7 — Freshness contract.** A running task finishes on its snapshot, undisturbed. At the
human's next message the agent is current. On finishing: sync, reconcile, say what changed
— and never revert or rework its own finished work uninvited (know and tell, never act).
> "what i do expect is that after it finishes, to check in first and then figure out what to do based on the new synced state"

**L8 — Finished means shared.** Work ends pushed (or PR'd per repo convention).
Committed-but-local is an alarm the finish message must name. The finish message speaks
the human frame: what exists, where it lives, what you can say next.
> "why do i need to tell it to push? with all the journals and dockets and intents and knowledge that should never be the case?"

**L9 — Plain speech where humans read.** No ids, codes, or machinery nouns on human
surfaces; hard register only where softness lets agents wiggle. Everything human-facing
passes the amnesia test: readable a day later, cold, no conversation context needed.
> "these diff entries need to be something the human can read after maybe a day where he totally forgot the convo he had with you"

**L10 — Silence must be earned.** An empty status means verified-nothing, never
could-not-check. Broken states name their fix or their fixer — no bare warnings.
> "if there is no knowledge diff shown, the human knows nothing new has been added somewhere else"

**L11 — Capped views disclose totals.** Never a silent window (the 88 lesson).
> "88 outstanding knowledge merges is unacceptable"

**L12 — Intents enter pre-sliced.** Under auto-build, an oversized intent self-blocks.
Write work items the size they can ship.

**L13 — Restarts carry the seed, never blank.** What is not carried gets a grave with
provenance (see GRAVES.md) so no generation re-derives or re-dies the same way.
> "restarting often isnt discussed but everyone does it"

**L14 — Complexity re-earns itself.** The human carries five sentences; machinery must
justify each noun against a promise and an owner quote before it exists. When in doubt,
the floor beats the tower.
> "its already too much to understand given how it should reduce down to the simple set of human desires it satisfies"
