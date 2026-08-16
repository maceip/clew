package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"clew/internal/config"
	"clew/internal/gitx"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/materialize"
	"clew/internal/state"
)

// app bundles the shared runtime: config + machine state.
type app struct {
	cfg *config.Config
	db  *state.DB
}

func load() (*app, error) {
	cfg := config.Load()
	db, err := state.Open(state.DefaultPath())
	if err != nil {
		return nil, fmt.Errorf("open state.db: %w", err)
	}
	return &app{cfg: cfg, db: db}, nil
}

func (a *app) close() { a.db.Close() }

// repoFromCwd resolves the registered repo containing the cwd.
func (a *app) repoFromCwd() (string, error) {
	cwd, _ := os.Getwd()
	if r := a.db.RepoFor(cwd); r != "" {
		return r, nil
	}
	// Not registered yet, but inside a git repo → helpful error.
	if gitx.IsRepo(cwd) {
		root, _ := gitx.Root(cwd)
		return "", fmt.Errorf("%s is not registered — run `clew init` first", root)
	}
	return "", fmt.Errorf("not inside a registered repo (run `clew init` in one, or see `clew status --all`)")
}

// targetRepos returns the repo of cwd, or all registered when all=true or
// cwd is outside any registered repo.
func (a *app) targetRepos(all bool) ([]string, error) {
	if !all {
		if r, err := a.repoFromCwd(); err == nil {
			return []string{r}, nil
		}
	}
	repos, err := a.db.Repos()
	if err != nil {
		return nil, err
	}
	if len(repos) == 0 {
		return nil, fmt.Errorf("no repos registered — run `clew init` inside one")
	}
	var out []string
	for _, r := range repos {
		out = append(out, r.Path)
	}
	return out, nil
}

// openJournal ensures the branch/worktree and loads the journal.
func (a *app) openJournal(repo string) (*journal.Journal, error) {
	wt, err := gitx.EnsureJournal(repo)
	if err != nil {
		return nil, err
	}
	return journal.Load(wt)
}

// regen is the projection regenerator passed to gitx.Sync.
func regen(wt string) error {
	j, err := journal.Load(wt)
	if err != nil {
		return err
	}
	st := journal.Compute(j, time.Now())
	return journal.WriteProjections(j, st, time.Now())
}

// syncAndMaterialize pushes/pulls the journal and refreshes .clew/.
func (a *app) syncAndMaterialize(repo string) (*gitx.SyncResult, error) {
	res, err := gitx.Sync(repo, regen)
	if err != nil {
		return res, err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return res, err
	}
	now := time.Now()
	st := journal.Compute(j, now)
	if err := materialize.Write(repo, j, st, a.db, now); err != nil {
		return res, err
	}
	return res, nil
}

func (a *app) provider() (llm.Provider, string) { return llm.Pick(a.cfg) }

func repoBase(repo string) string { return filepath.Base(repo) }
