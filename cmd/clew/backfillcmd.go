package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"clew/internal/adapters"
	"clew/internal/extract"
	"clew/internal/gitx"
	"clew/internal/state"
)

// cmdBackfill: retroactive extraction over existing session files (§5.3).
// Cost-capped: runs only inside explicit --budget bounds (§6.4); resumable
// by watermark. Original timestamps are preserved by the extractor.
func cmdBackfill(args []string) error {
	fs := flag.NewFlagSet("backfill", flag.ExitOnError)
	since := fs.String("since", "90d", "how far back (e.g. 30d, 12h)")
	budget := fs.Int("budget", 0, "token budget for this run (REQUIRED, §6.4)")
	fs.Parse(args)
	if *budget <= 0 {
		return fmt.Errorf("backfill runs only inside explicit --budget bounds (§6.4): pass --budget <tokens>")
	}
	dur, err := parseDur(*since)
	if err != nil {
		return err
	}
	cutoff := time.Now().Add(-dur)

	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	if a.db.Get(liveCursorMigration) == "" {
		lock := filepath.Join(gitx.Home(), "watch.lock")
		pid, adopted, err := claimWatchLock(lock)
		if err != nil {
			return err
		}
		if adopted {
			return fmt.Errorf("watcher pid %d is still running with unmigrated cursors; restart it before backfill", pid)
		}
		defer os.Remove(lock)
		if err := migrateLiveCursors(a); err != nil {
			return err
		}
	}
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	provider, note := a.provider()
	if provider == nil {
		return fmt.Errorf("no extraction provider: %s", note)
	}
	metered := newBudgetedProvider(provider, a.db, a.cfg, "backfill", *budget)
	provider = metered
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}

	spent := 0
	files := discoverSince(repo, cutoff)
	fmt.Printf("backfill: %d candidate session file(s) since %s, budget %d tokens\n", len(files), cutoff.Format("2006-01-02"), *budget)
	for _, cf := range files {
		spent = metered.Spent()
		if spent >= *budget {
			fmt.Printf("budget exhausted (%d tokens) — resumable: watermarks persist, re-run to continue\n", spent)
			break
		}
		end, ok := a.db.WatermarkOK("history-end:" + cf.path)
		if !ok {
			end, err = adapters.CompleteOffset(cf.path)
			if err != nil {
				fmt.Printf("  !! %s: cannot establish history boundary: %v\n", filepath.Base(cf.path), err)
				continue
			}
			if _, err := a.db.InitWatermarks(
				state.WatermarkInit{File: "tail:" + cf.path, Adapter: cf.adapter.ID(), Repo: repo, Offset: end},
				state.WatermarkInit{File: "extract:" + cf.path, Adapter: cf.adapter.ID(), Repo: repo, Offset: end},
				state.WatermarkInit{File: "history-end:" + cf.path, Adapter: cf.adapter.ID(), Repo: repo, Offset: end},
				state.WatermarkInit{File: "backfill:" + cf.path, Adapter: cf.adapter.ID(), Repo: repo, Offset: 0},
			); err != nil {
				return err
			}
			end = a.db.Watermark("history-end:" + cf.path)
		}
		currentEnd, err := adapters.CompleteOffset(cf.path)
		if err != nil {
			fmt.Printf("  !! %s: cannot read current boundary: %v\n", filepath.Base(cf.path), err)
			continue
		}
		if currentEnd < end {
			fmt.Printf("  !! %s: history truncated from %d to %d bytes; refusing to guess\n", filepath.Base(cf.path), end, currentEnd)
			continue
		}
		for spent < *budget {
			off := a.db.Watermark("backfill:" + cf.path)
			if end <= off {
				break
			}
			out, err := extract.RunUntil(j, provider, cf.adapter, cf.path, off, end, a.cfg.Surface, time.Now())
			spent = metered.Spent()
			if err != nil {
				fmt.Printf("  !! %s: %v\n", filepath.Base(cf.path), err)
				break
			}
			if out.Parked {
				if err := extract.ParkRawRange(a.db, cf.adapter, cf.path, off, end, out.ParkReason); err != nil {
					fmt.Printf("  !! failed to park %s: %v (cursor not advanced)\n", filepath.Base(cf.path), err)
					break
				}
				if err := a.db.SetWatermark("backfill:"+cf.path, cf.adapter.ID(), repo, end); err != nil {
					return err
				}
				fmt.Printf("  !! parked slice of %s: %s\n", filepath.Base(cf.path), out.ParkReason)
				break
			}
			if err := a.db.SetWatermark("backfill:"+cf.path, cf.adapter.ID(), repo, out.NewOffset); err != nil {
				return err
			}
			if len(out.Entries) > 0 {
				fmt.Printf("  %s: +%d entries (+%d events, %d rejected) [%d tokens]\n",
					filepath.Base(cf.path), len(out.Entries), len(out.Events), out.Rejected, out.Tokens)
			}
			if out.NewOffset >= end {
				break
			}
			if out.NewOffset <= off {
				fmt.Printf("  !! %s: backfill made no progress at offset %d; stopped\n", filepath.Base(cf.path), off)
				break
			}
		}
	}
	spent = metered.Spent()
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("backfill done: %d tokens spent\n", spent)
	return nil
}

type candidateFile struct {
	adapter adapters.Adapter
	path    string
}

// discoverSince finds historical session files for the repo: the standard
// adapter discovery plus codex's dated directories over the window.
func discoverSince(repo string, cutoff time.Time) []candidateFile {
	var out []candidateFile
	seen := map[string]bool{}
	add := func(a adapters.Adapter, p string) {
		if !seen[p] {
			seen[p] = true
			out = append(out, candidateFile{a, p})
		}
	}
	for _, a := range adapters.All() {
		for _, f := range a.Discover(repo) {
			if fi, err := os.Stat(f); err == nil && fi.ModTime().After(cutoff) {
				add(a, f)
			}
		}
	}
	// codex history beyond the live window: walk dated dirs in range.
	home, _ := os.UserHomeDir()
	cx := &adapters.Codex{}
	for d := 0; ; d++ {
		day := time.Now().AddDate(0, 0, -d)
		if day.Before(cutoff) {
			break
		}
		glob := filepath.Join(home, ".codex", "sessions", day.Format("2006"), day.Format("01"), day.Format("02"), "rollout-*.jsonl")
		m, _ := filepath.Glob(glob)
		for _, f := range m {
			if cwd := adapters.CodexCWD(f); cwd == repo || strings.HasPrefix(cwd, repo+string(os.PathSeparator)) {
				add(cx, f)
			}
		}
	}
	return out
}

func parseDur(s string) (time.Duration, error) {
	if strings.HasSuffix(s, "d") {
		n, err := strconv.Atoi(strings.TrimSuffix(s, "d"))
		if err != nil {
			return 0, fmt.Errorf("bad --since %q", s)
		}
		return time.Duration(n) * 24 * time.Hour, nil
	}
	return time.ParseDuration(s)
}
