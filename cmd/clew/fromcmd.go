package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/maceip/clew/internal/config"
	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/lineage"
	"github.com/maceip/clew/internal/model"
	"github.com/maceip/clew/internal/seed"
	"github.com/maceip/clew/internal/state"
)

// cmdFrom is the single explicit lineage ceremony. With no argument it only
// ranks already-maintained seeds; with one argument it imports exactly that
// seed and records an immutable lineage link. It never generates or refreshes
// a predecessor seed.
func cmdFrom(args []string) error {
	if len(args) > 1 {
		return fmt.Errorf("usage: clew from [<repo-path|name|id>]")
	}
	cwd, _ := os.Getwd()
	if !gitx.IsRepo(cwd) {
		return fmt.Errorf("clew from must run inside the successor git repository")
	}
	repo, err := gitx.Root(cwd)
	if err != nil {
		return err
	}
	// A truly new machine has no candidate registry. Listing there must be a
	// literal read: do not create ~/.clew, config.yaml, SQLite, a project
	// journal, or workspace files merely to print an empty result.
	if len(args) == 0 {
		if _, statErr := os.Stat(state.DefaultPath()); os.IsNotExist(statErr) {
			fmt.Println("lineage candidates (topic overlap + recency; nothing is carried):")
			fmt.Println("(none with a maintained SEED.md)")
			return nil
		}
	}
	if len(args) == 1 {
		if err := config.WriteSkeleton(); err != nil {
			return err
		}
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	var target *seed.Snapshot
	if !a.db.RepoRegistered(repo) && len(args) == 0 {
		// Candidate inspection is genuinely non-carrying even before birth. A
		// name/README vocabulary is enough to rank; no journal or seed is made.
		in := seed.BuildInputForRepo(repo, nil)
		target = &seed.Snapshot{Repository: in.Repository, Topics: in.Topics, Lifecycle: seed.Lifecycle{State: "active"}}
	} else {
		if !a.db.RepoRegistered(repo) {
			if err := autoBirth(a, repo); err != nil {
				return fmt.Errorf("initialize successor before lineage: %w", err)
			}
		}
		target, err = maintainedSeed(repo)
		if err != nil {
			return fmt.Errorf("read successor's ambient seed: %w", err)
		}
	}
	candidates, issues := lineageCandidates(a, repo)

	if len(args) == 0 {
		ranked, err := lineage.Rank(lineage.Target{
			RepositoryID: target.Repository.ID,
			Name:         target.Repository.Name,
			Topics:       target.Topics,
		}, candidates, time.Now())
		if err != nil {
			return err
		}
		fmt.Println("lineage candidates (topic overlap + recency; nothing is carried):")
		if len(ranked) == 0 {
			fmt.Println("(none with a maintained SEED.md)")
		}
		for _, candidate := range ranked {
			fmt.Printf("%s — %s\n", candidate.Summary(), candidate.Location)
		}
		for _, issue := range issues {
			fmt.Printf("!! %s\n", issue)
		}
		return nil
	}

	predecessor, location, err := resolveLineageSource(args[0], candidates)
	if err != nil {
		return err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	result, err := lineage.Import(lineage.ImportRequest{
		Journal:          j,
		TargetRepository: target.Repository.ID,
		Snapshot:         predecessor,
		By:               model.By{Who: "human", Surface: a.cfg.Surface},
		At:               time.Now().UTC(),
	})
	if err != nil {
		return err
	}
	if result.AlreadyImported {
		if _, err := a.syncAndMaterialize(repo); err != nil {
			return err
		}
		fmt.Printf("lineage already declared from %s at seed %s; carried material was not resurrected\n", predecessor.Repository.Name, result.Link.SeedDigest)
		return nil
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("declared lineage %s ← %s (%s): carried %d entries and %d exhibits/terminal/status events; link %s\n",
		target.Repository.Name, predecessor.Repository.Name, location,
		result.NewEntries, result.NewEvents, result.Link.ID)
	return nil
}

func lineageCandidates(a *app, targetRepo string) ([]lineage.Candidate, []string) {
	repos, err := a.db.Repos()
	if err != nil {
		return nil, []string{err.Error()}
	}
	var candidates []lineage.Candidate
	var issues []string
	for _, registered := range repos {
		if registered.Path == targetRepo {
			continue
		}
		snapshot, err := maintainedSeed(registered.Path)
		if err != nil {
			issues = append(issues, fmt.Sprintf("%s: %v", registered.Path, err))
			continue
		}
		candidates = append(candidates, lineage.Candidate{Snapshot: snapshot, Location: registered.Path})
	}
	return candidates, issues
}

func resolveLineageSource(arg string, candidates []lineage.Candidate) (*seed.Snapshot, string, error) {
	// An explicit directory or SEED.md path need not be registered, but its
	// seed must already exist. This is still a read, never demand generation.
	if info, err := os.Stat(arg); err == nil {
		path := arg
		location := arg
		if info.IsDir() {
			abs, _ := filepath.Abs(arg)
			location = filepath.Clean(abs)
			snapshot, err := maintainedSeed(location)
			return snapshot, location, err
		}
		snapshot, err := seed.Read(path)
		return snapshot, path, err
	}

	abs, _ := filepath.Abs(arg)
	var matches []lineage.Candidate
	for _, candidate := range candidates {
		s := candidate.Snapshot
		if s == nil {
			continue
		}
		if filepath.Clean(candidate.Location) == filepath.Clean(abs) ||
			arg == candidate.Location || arg == s.Repository.ID || arg == s.Repository.Remote ||
			strings.EqualFold(arg, s.Repository.Name) || strings.EqualFold(arg, filepath.Base(candidate.Location)) {
			matches = append(matches, candidate)
		}
	}
	switch len(matches) {
	case 0:
		return nil, "", fmt.Errorf("no maintained predecessor seed matches %q; run `clew from` to list candidates or pass a repo path", arg)
	case 1:
		return matches[0].Snapshot, matches[0].Location, nil
	default:
		var paths []string
		for _, match := range matches {
			paths = append(paths, match.Location)
		}
		return nil, "", fmt.Errorf("predecessor %q is ambiguous: %s", arg, strings.Join(paths, ", "))
	}
}

func maintainedSeed(repo string) (*seed.Snapshot, error) {
	// A ready journal worktree is the canonical, potentially newer projection:
	// sync commits it before workspace materialization. Prefer it so a crash in
	// that small window cannot hide new ancestry and weaken cycle detection.
	if gitx.JournalReady(repo) {
		path := filepath.Join(gitx.WorktreeDir(repo), "SEED.md")
		snapshot, err := seed.Read(path)
		if err == nil {
			return snapshot, nil
		}
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
	}
	// A moved repository may no longer have a usable linked worktree, while
	// its branch ref remains perfectly readable. This is still a pure read.
	projection, err := gitx.Show(repo, "SEED.md")
	if err == nil {
		snapshot, parseErr := seed.Parse([]byte(projection))
		if parseErr != nil {
			return nil, fmt.Errorf("journal-branch SEED.md in %s: %w", repo, parseErr)
		}
		return snapshot, nil
	}
	workspace := filepath.Join(repo, ".clew", "SEED.md")
	snapshot, workspaceErr := seed.Read(workspace)
	if workspaceErr == nil {
		return snapshot, nil
	}
	if !os.IsNotExist(workspaceErr) {
		return nil, fmt.Errorf("%s: %w", workspace, workspaceErr)
	}
	return nil, fmt.Errorf("no watcher-maintained SEED.md in %s", repo)
}
