package main

import (
	"github.com/maceip/clew/internal/mcp"
)

// cmdMCP: optional stdio MCP surface (§8.1) — journal_search, journal_get,
// journal_note. Never load-bearing (I1).
func cmdMCP(args []string) error {
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	srv := &mcp.Server{
		Journal: j,
		Surface: a.cfg.Surface,
		AfterWrite: func() {
			a.syncAndMaterialize(repo)
		},
	}
	return srv.ServeStdio()
}
