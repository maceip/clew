# Return cloud work to the journal

Use the same setup in every cloud environment you control:

```bash
export PATH="$HOME/.local/bin:$PATH"
GOBIN="$HOME/.local/bin" go install github.com/maceip/clew/cmd/clew@main
clew init
```

The environment needs the repository's normal Git credentials, including
permission to push the `clew/journal` branch. Source code may still return by a
pull request; the journal branch is independent of that choice.

## Cursor cloud

Put the install lines in the environment setup command. Run `clew init` after
the checkout is available. If the environment retains Cursor CLI transcripts,
the watcher can read them normally. Otherwise export the session before the
environment is discarded and run `clew witness` on that export.

## Codex cloud

Put the install lines in the Codex environment setup script. Run `clew init`
inside the task checkout. A Codex JSONL export can be returned with:

```bash
clew witness /path/to/codex-session.jsonl
```

This is useful when the code itself is delivered through a pull request. The
command writes learned project memory to the journal branch, not to the source
branch.

## Claude cloud

Put the install lines in the Claude environment setup command and run
`clew init` after checkout. Claude JSONL exports use the same return command:

```bash
clew witness /path/to/claude-session.jsonl
```

## What happens when a surface is narrower

`clew witness` recognizes only formats already pinned by clew. It refuses an
unknown or ambiguous export and keeps the failure visible. It does not invent a
second journal representation. If the environment cannot push the journal
branch at all, the command fails rather than claiming the session returned.
