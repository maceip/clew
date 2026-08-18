package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"clew/internal/config"
	"clew/internal/gitx"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/materialize"
	"clew/internal/owner"
	"clew/internal/seed"
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

// regenFor is the projection regenerator passed to project git sync. The
// ambient SEED.md is written on every journal union, so the carry-kit exists
// before a restart is contemplated. Owner laws are intentionally excluded.
func regenFor(repo string) func(string) error {
	return func(wt string) error {
		j, err := journal.Load(wt)
		if err != nil {
			return err
		}
		now := time.Now()
		st := journal.Compute(j, now)
		if err := journal.WriteProjections(j, st, now); err != nil {
			return err
		}
		snapshot, err := materialize.BuildSeedForRepo(repo, j, st)
		if err != nil {
			return err
		}
		_, err = seed.WriteOnJournalChange(filepath.Join(wt, "SEED.md"), snapshot)
		return err
	}
}

// ownerLawBlock loads or synchronizes the owner journal and returns the exact
// independently capped injection block. A concurrent remote overflow is
// rendered deterministically and made loud; it is never silently accepted.
func (a *app) ownerLawBlock(sync bool, now time.Time) (string, error) {
	store := owner.Default(a.cfg.Owner.Remote)
	var (
		j          *journal.Journal
		syncResult *gitx.SyncResult
		err        error
	)
	if sync {
		j, syncResult, err = store.Sync()
	} else {
		j, err = store.Open()
	}
	if err != nil {
		_ = a.db.Set("owner-sync-error", err.Error())
		return "", err
	}
	if sync {
		if err := store.RequireVerifiedSync(syncResult, "sync owner laws"); err != nil {
			_ = a.db.Set("owner-sync-error", err.Error())
			return "", err
		}
		if err := store.MarkVerifiedCache(); err != nil {
			_ = a.db.Set("owner-sync-error", err.Error())
			return "", err
		}
		_ = a.db.Set("owner-sync-error", "")
	} else if err := store.RequireVerifiedCache("load cached owner laws"); err != nil {
		_ = a.db.Set("owner-sync-error", err.Error())
		return "", err
	}
	rendered := owner.Render(j, now)
	if rendered.Overflow {
		_ = a.db.Set("owner-law-overflow", fmt.Sprintf("owner laws require %d bytes; cap=%d; omitted=%s", rendered.RequiredBytes, owner.LawCap, strings.Join(rendered.Omitted, ",")))
	} else {
		_ = a.db.Set("owner-law-overflow", "")
	}
	return rendered.Markdown, nil
}

// syncAndMaterialize pushes/pulls the journal and refreshes .clew/.
func (a *app) syncAndMaterialize(repo string) (*gitx.SyncResult, error) {
	res, err := gitx.Sync(repo, regenFor(repo))
	if err != nil {
		return res, err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return res, err
	}
	now := time.Now()
	st := journal.Compute(j, now)
	laws, err := a.ownerLawBlock(true, now)
	if err != nil {
		return res, err
	}
	if err := materialize.WriteWithOwner(repo, j, st, a.db, laws, now); err != nil {
		return res, err
	}
	return res, nil
}

func (a *app) provider() (llm.Provider, string) { return llm.Pick(a.cfg) }

func repoBase(repo string) string { return filepath.Base(repo) }
