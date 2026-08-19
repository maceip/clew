#!/bin/sh
# Installs the wiring for claude and codex, idempotently, with backups.
# Verify hook schema against your installed versions on first run (L5).
set -e
repo="$(cd "$(dirname "$0")/.." && pwd)"
echo "== codex: project-level .codex/hooks.json"
mkdir -p "$repo/.codex"
[ -f "$repo/.codex/hooks.json" ] && cp "$repo/.codex/hooks.json" "$repo/.codex/hooks.json.bak"
cp "$repo/hooks/codex-hooks.json" "$repo/.codex/hooks.json"
echo "== claude: merge hooks/claude-settings-snippet.json into ~/.claude/settings.json"
echo "   (manual merge on first run; claude validates on next start)"
echo "done. start a fresh session in each agent and ask: 'what are this project's laws?'"
