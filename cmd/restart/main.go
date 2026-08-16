// restart — the journal system (JOURNAL_SPEC.md). One binary, one git branch.
package main

import (
	"fmt"
	"os"
)

const usage = `restart — durable project journal for humans steering agents

usage:
  restart init [--carry <dir>] [--no-archaeology]   register repo; archaeology; install snippets
  restart watch [install|uninstall]                 start/adopt the machine's watcher
  restart status [--all]                            the glance
  restart map [--html <file>]                       intent × reality with absence
  restart inbox [answer|ack|drop] …                 human-blocking items only
  restart journal show|edit|confirm|reject|supersede|answer|note …
  restart manifest [--spec <file>] [--out <dir>] [--yes]   restart kit
  restart backfill [--since 90d] --budget <tokens>  retroactive extraction
  restart wrap -- <agent argv…>                     PTY tee for hookless agents
  restart redact <entry-id>                         scrub + rewrite journal branch
  restart mcp                                       stdio MCP server (optional)
`

func main() {
	if len(os.Args) < 2 {
		fmt.Print(usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]
	var err error
	switch cmd {
	case "init":
		err = cmdInit(args)
	case "watch":
		err = cmdWatch(args)
	case "status":
		err = cmdStatus(args)
	case "map":
		err = cmdMap(args)
	case "inbox":
		err = cmdInbox(args)
	case "journal":
		err = cmdJournal(args)
	case "manifest":
		err = cmdManifest(args)
	case "backfill":
		err = cmdBackfill(args)
	case "wrap":
		err = cmdWrap(args)
	case "redact":
		err = cmdRedact(args)
	case "mcp":
		err = cmdMCP(args)
	case "help", "-h", "--help":
		fmt.Print(usage)
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n%s", cmd, usage)
		os.Exit(2)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "restart:", err)
		os.Exit(1)
	}
}
