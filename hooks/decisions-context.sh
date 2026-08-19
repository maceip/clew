#!/bin/sh
# Prints the project memory. Used by claude + codex hooks at session start
# and prompt submit. Pull is best-effort and quiet; reading never fails loud.
d="${CLAUDE_PROJECT_DIR:-${CODEX_PROJECT_DIR:-$PWD}}"
git -C "$d" pull -q --rebase --autostash 2>/dev/null || true
[ -f "$d/DECISIONS.md" ] && cat "$d/DECISIONS.md"
