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
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	provider, note := a.provider()
	if provider == nil {
		return fmt.Errorf("no extraction provider: %s", note)
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}

	spent := 0
	files := discoverSince(repo, cutoff)
	fmt.Printf("backfill: %d candidate session file(s) since %s, budget %d tokens\n", len(files), cutoff.Format("2006-01-02"), *budget)
	for _, cf := range files {
		if spent >= *budget {
			fmt.Printf("budget exhausted (%d tokens) — resumable: watermarks persist, re-run to continue\n", spent)
			break
		}
		for spent < *budget {
			off := a.db.Watermark("extract:" + cf.path)
			fi, err := os.Stat(cf.path)
			if err != nil || fi.Size() <= off {
				break
			}
			out, err := extract.Run(j, provider, cf.adapter, cf.path, off, a.cfg.Surface, time.Now())
			if err != nil {
				fmt.Printf("  !! %s: %v\n", filepath.Base(cf.path), err)
				break
			}
			spent += out.Tokens
			a.db.AddTokens("spent", out.Tokens)
			if out.Parked {
				extract.ParkSlice(a.db, cf.adapter, cf.path, off, out.ParkReason)
				a.db.SetWatermark("extract:"+cf.path, cf.adapter.ID(), repo, fi.Size())
				fmt.Printf("  !! parked slice of %s: %s\n", filepath.Base(cf.path), out.ParkReason)
				break
			}
			a.db.SetWatermark("extract:"+cf.path, cf.adapter.ID(), repo, out.NewOffset)
			if len(out.Entries) > 0 {
				fmt.Printf("  %s: +%d entries (+%d events, %d rejected) [%d tokens]\n",
					filepath.Base(cf.path), len(out.Entries), len(out.Events), out.Rejected, out.Tokens)
			}
			if out.NewOffset >= fi.Size() {
				break
			}
		}
	}
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
