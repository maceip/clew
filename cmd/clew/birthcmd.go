package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"clew/internal/adapters"
	"clew/internal/config"
	"clew/internal/gitx"
	"clew/internal/journal"
	"clew/internal/materialize"
	"clew/internal/state"
)

// cmdBirth is the quiet SessionStart entrypoint installed in the owner's
// Claude settings. Its stdout is agent context, so all operational detail and
// failures stay off stdout. Non-git sessions are deliberately a no-op.
func cmdBirth(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("internal _birth takes no arguments")
	}
	return runBirthHook(os.Stdin, os.Stdout)
}

func runBirthHook(input io.Reader, output io.Writer) error {
	cwd, err := parseBirthHookCWD(input)
	if err != nil {
		return err
	}
	if !gitx.IsRepo(cwd) {
		return nil
	}
	root, err := gitx.Root(cwd)
	if err != nil {
		return nil
	}
	if birthInternalRepo(root) {
		return nil
	}
	// The ordinary SessionStart path consumes the watcher-maintained atomic
	// context directly. It performs no DB migration, owner load, journal sync,
	// or materialization, so a transient background failure cannot suppress a
	// context that was already valid before this session existed.
	if gitx.BirthReady(root) {
		if context, readErr := os.ReadFile(filepath.Join(root, ".clew", "context.md")); readErr == nil {
			_, err = output.Write(context)
			return err
		}
	}
	var context []byte
	err = withBirthMachineLock(func() error {
		if err := config.WriteSkeleton(); err != nil {
			return err
		}
		a, err := load()
		if err != nil {
			return err
		}
		defer a.close()
		if err := autoBirth(a, root); err != nil {
			_ = a.db.Set("birth-error:"+root, err.Error())
			return err
		}
		_ = a.db.Set("birth-error:"+root, "")
		context, err = os.ReadFile(filepath.Join(root, ".clew", "context.md"))
		return err
	})
	if err != nil {
		return err
	}
	_, err = output.Write(context)
	return err
}

// autoBirth enrolls a repository without archaeology, carry, lineage, or a
// decision card. The only cross-project content materialized here is the
// separately human-certified owner-law block. It is safe to call repeatedly.
func autoBirth(a *app, root string) error {
	return autoBirthFromCandidate(a, root, nil)
}

// autoBirthFromCandidate preserves the already-written beginning of a
// daemon-discovered session. Older transcripts remain forward-only, while the
// concrete session that revealed this birth starts at offset zero.
func autoBirthFromCandidate(a *app, root string, candidate *adapters.BirthCandidate) error {
	if a == nil || a.db == nil {
		return fmt.Errorf("auto-birth: application state is unavailable")
	}
	resolved, err := gitx.Root(root)
	if err != nil {
		return err
	}
	root = resolved
	if birthInternalRepo(root) {
		return nil
	}
	registered := a.db.RepoRegistered(root)
	journalReady := gitx.JournalReady(root)
	if registered && journalReady {
		if _, err := os.Stat(filepath.Join(root, ".clew", "context.md")); err == nil {
			if !gitx.BirthReady(root) {
				if err := gitx.MarkBirthReady(root); err != nil {
					return err
				}
			}
			return nil // steady-state SessionStart is a single local file read
		}
	}

	lockDir := filepath.Join(gitx.Home(), "birth-locks")
	if err := os.MkdirAll(lockDir, 0o755); err != nil {
		return err
	}
	lock := filepath.Join(lockDir, gitx.RepoID(root)+".lock")
	_, adopted, err := claimWatchLock(lock)
	if err != nil {
		return err
	}
	if adopted {
		// The machine watcher and SessionStart hook can observe the same birth.
		// Wait only for the owning process to publish the atomic context file.
		deadline := time.Now().Add(15 * time.Second)
		for time.Now().Before(deadline) {
			if a.db.RepoRegistered(root) && gitx.JournalReady(root) {
				if _, err := os.Stat(filepath.Join(root, ".clew", "context.md")); err == nil {
					return nil
				}
			}
			time.Sleep(25 * time.Millisecond)
		}
		return fmt.Errorf("auto-birth for %s is still owned by another process", root)
	}
	defer os.Remove(lock)

	// An enrollment missing its materialized context is repaired synchronously.
	registered = a.db.RepoRegistered(root)
	journalReady = gitx.JournalReady(root)
	if registered && !journalReady && !gitx.HasBirthIncarnation(root) {
		// Registration is keyed by pathname, so it can outlive the repository
		// that originally occupied that path. Only the absence of local
		// incarnation evidence distinguishes that case from a damaged external
		// journal worktree. The old journal itself remains intact as an explicit
		// future lineage candidate.
		if err := a.db.ResetRepoIncarnation(root); err != nil {
			return fmt.Errorf("reset stale birth state for %s: %w", root, err)
		}
		if _, err := gitx.AssignFreshJournalIncarnation(root); err != nil {
			// Keep registration as a retry marker if persisting the incarnation
			// unexpectedly fails after the disposable-state transaction. A later
			// birth will re-prove the same fresh-repository condition rather than
			// silently falling back to the predecessor's path-derived identity.
			_ = a.db.RegisterRepo(root, "")
			return fmt.Errorf("assign fresh journal incarnation for %s: %w", root, err)
		}
		registered = false
	}
	needsEnrollment := !registered || !journalReady
	if needsEnrollment {
		if _, err := baselineBirthSessionFiles(a.db, root, adapters.All(), candidate); err != nil {
			return err
		}
		if _, err := gitx.EnsureJournal(root); err != nil {
			return err
		}
		for _, name := range []string{"CLAUDE.md", "AGENTS.md"} {
			if _, err := ensureSnippet(filepath.Join(root, name)); err != nil {
				return err
			}
		}
		if err := installClaudeHook(root); err != nil {
			// Nudges are secondary to birth context; keep enrollment working and
			// surface this setup failure through machine status.
			_ = a.db.Set("birth-hook-error:"+root, err.Error())
		} else {
			_ = a.db.Set("birth-hook-error:"+root, "")
		}
		if err := ensureGitignore(root); err != nil {
			return err
		}
	}

	if needsEnrollment {
		if _, err := gitx.Sync(root, regenFor(root)); err != nil {
			return err
		}
	}
	j, err := a.openJournal(root)
	if err != nil {
		return err
	}
	if err := j.Reload(); err != nil {
		return err
	}
	now := time.Now()
	laws, err := a.ownerLawBlock(false, now)
	if err != nil {
		return err
	}
	if err := materialize.WriteWithOwner(root, j, journal.Compute(j, now), a.db, laws, now); err != nil {
		return err
	}
	if err := a.db.RegisterRepo(root, gitx.RemoteName(root)); err != nil {
		return err
	}
	return gitx.MarkBirthReady(root)
}

func baselineBirthSessionFiles(db interface {
	InitWatermarks(...state.WatermarkInit) (int, error)
}, repo string, all []adapters.Adapter, live *adapters.BirthCandidate) (int, error) {
	baselined := 0
	livePath := ""
	if live != nil {
		livePath, _ = filepath.Abs(live.File)
	}
	seenLive := false
	for _, adapter := range all {
		for _, file := range adapter.Discover(repo) {
			abs, _ := filepath.Abs(file)
			offset := int64(0)
			if livePath == "" || filepath.Clean(abs) != filepath.Clean(livePath) {
				var err error
				offset, err = adapters.CompleteOffset(file)
				if err != nil {
					return baselined, fmt.Errorf("baseline %s: %w", file, err)
				}
			} else {
				seenLive = true
			}
			added, err := db.InitWatermarks(
				state.WatermarkInit{File: "tail:" + file, Adapter: adapter.ID(), Repo: repo, Offset: offset},
				state.WatermarkInit{File: "extract:" + file, Adapter: adapter.ID(), Repo: repo, Offset: offset},
				state.WatermarkInit{File: "history-end:" + file, Adapter: adapter.ID(), Repo: repo, Offset: offset},
				state.WatermarkInit{File: "backfill:" + file, Adapter: adapter.ID(), Repo: repo, Offset: 0},
			)
			if err != nil {
				return baselined, err
			}
			if added > 0 {
				baselined++
			}
		}
	}
	if live != nil && !seenLive {
		added, err := db.InitWatermarks(
			state.WatermarkInit{File: "tail:" + live.File, Adapter: live.Adapter, Repo: repo, Offset: 0},
			state.WatermarkInit{File: "extract:" + live.File, Adapter: live.Adapter, Repo: repo, Offset: 0},
			state.WatermarkInit{File: "history-end:" + live.File, Adapter: live.Adapter, Repo: repo, Offset: 0},
			state.WatermarkInit{File: "backfill:" + live.File, Adapter: live.Adapter, Repo: repo, Offset: 0},
		)
		if err != nil {
			return baselined, err
		}
		if added > 0 {
			baselined++
		}
	}
	return baselined, nil
}

// withBirthMachineLock serializes cold state migration and owner bootstrap.
// Steady-state SessionStart still returns via the fast path above, but the
// first two simultaneous births cannot race SQLite migration or git init.
func withBirthMachineLock(fn func() error) error {
	if err := os.MkdirAll(gitx.Home(), 0o755); err != nil {
		return err
	}
	lock := filepath.Join(gitx.Home(), "birth-bootstrap.lock")
	deadline := time.Now().Add(30 * time.Second)
	for {
		_, adopted, err := claimWatchLock(lock)
		if err != nil {
			return err
		}
		if !adopted {
			defer os.Remove(lock)
			return fn()
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("birth bootstrap is still owned by another process")
		}
		time.Sleep(25 * time.Millisecond)
	}
}

func primeBirthMachine() error {
	return withBirthMachineLock(func() error {
		if err := config.WriteSkeleton(); err != nil {
			return err
		}
		a, err := load()
		if err != nil {
			return err
		}
		defer a.close()
		_, err = a.ownerLawBlock(true, time.Now())
		return err
	})
}

func birthInternalRepo(root string) bool {
	root, _ = filepath.Abs(root)
	for _, internal := range []string{
		filepath.Join(gitx.Home(), "worktrees"),
		filepath.Join(gitx.Home(), "owner"),
	} {
		internal, _ = filepath.Abs(internal)
		if root == internal || pathWithin(root, internal) {
			return true
		}
	}
	return false
}

func pathWithin(path, parent string) bool {
	rel, err := filepath.Rel(parent, path)
	return err == nil && rel != "." && rel != ".." && !filepath.IsAbs(rel) &&
		!strings.HasPrefix(rel, ".."+string(os.PathSeparator))
}
