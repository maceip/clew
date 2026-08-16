package main

import (
	"flag"
	"fmt"
	"sort"
	"strings"
	"time"

	"clew/internal/adapters"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/poller"
)

// cmdStatus renders the glance (§8.2): five fixed sections, ≤7 lines each
// (I8), every line carrying an entry id. Degradations are loud (I2).
func cmdStatus(args []string) error {
	fs := flag.NewFlagSet("status", flag.ExitOnError)
	all := fs.Bool("all", false, "all registered repos (cross-project join)")
	fs.Parse(args)
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repos, err := a.targetRepos(*all)
	if err != nil {
		return err
	}
	for i, repo := range repos {
		if i > 0 {
			fmt.Println()
		}
		if err := statusRepo(a, repo); err != nil {
			fmt.Printf("!! %s: %v\n", repoBase(repo), err)
		}
	}
	return statusMachine(a)
}

func statusRepo(a *app, repo string) error {
	now := time.Now()
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	st := journal.Compute(j, now)
	snap, _ := poller.Poll(a.db, repo)

	fmt.Printf("── %s ──\n", repoBase(repo))

	// SESSIONS: agent · surface · repo · branch · behind-by · live footprint
	fmt.Println("SESSIONS")
	sessions := a.db.LiveSessions(repo, 30*time.Minute)
	if len(sessions) == 0 {
		fmt.Println("  (none live in the last 30 min)")
	}
	for i, s := range sessions {
		if i >= 7 {
			break
		}
		fp := a.db.Footprints(s.ID)
		fpStr := "(no footprint yet)"
		if len(fp) > 0 {
			fpStr = strings.Join(firstN(fp, 2), ", ")
		}
		branch, behind := "", ""
		if snap != nil {
			branch = snap.Branch
			if snap.Behind > 0 {
				behind = fmt.Sprintf(" · behind %d", snap.Behind)
			}
		}
		fmt.Printf("  %s · %s · %s · %s%s · %s\n", s.Agent, s.Surface, repoBase(repo), branch, behind, clipStr(fpStr, 60))
	}

	section := func(name string, t model.EntryType, keep func(journal.Status) bool, line func(*model.Entry, *journal.Computed) string) {
		fmt.Println(name)
		var es []*model.Entry
		for id, e := range j.Entries {
			if e.Type == t && st[id] != nil && keep(st[id].Status) {
				es = append(es, e)
			}
		}
		sort.Slice(es, func(a, b int) bool { return es[a].ID > es[b].ID })
		if len(es) == 0 {
			fmt.Println("  (none)")
		}
		for i, e := range es {
			if i >= 7 { // I8: glance ≤ 7 items per section
				fmt.Printf("  … %d more (see .clew/journal.md)\n", len(es)-7)
				break
			}
			fmt.Println("  " + line(e, st[e.ID]))
		}
	}

	section("DECIDED", model.Decision,
		func(s journal.Status) bool {
			return s == journal.StActive || s == journal.StPossibleContradiction || s == journal.StContradicted
		},
		func(e *model.Entry, c *journal.Computed) string {
			flag := ""
			if c.Status != journal.StActive {
				flag = "  ⚠ " + string(c.Status)
			}
			return fmt.Sprintf("%s %s%s", e.ID, e.Title, flag)
		})
	section("LEARNED", model.Finding,
		func(s journal.Status) bool { return s == journal.StCurrent || s == journal.StSuspect },
		func(e *model.Entry, c *journal.Computed) string {
			flag := ""
			if c.Status == journal.StSuspect {
				flag = "  ⚠ suspect"
			}
			return fmt.Sprintf("%s %s%s", e.ID, e.Title, flag)
		})
	section("OPEN", model.Question,
		func(s journal.Status) bool { return s == journal.StOpen },
		func(e *model.Entry, c *journal.Computed) string {
			star := ""
			if e.Asks == "human" {
				star = " ★"
			}
			age := int(now.Sub(e.Created()).Hours() / 24)
			return fmt.Sprintf("%s %s · %dd%s", e.ID, e.Title, age, star)
		})

	fmt.Println("ALERTS")
	alerts := a.db.OpenAlerts(repo, false)
	if len(alerts) == 0 {
		fmt.Println("  (none)")
	}
	for i, al := range alerts {
		if i >= 7 {
			fmt.Printf("  … %d more (clew inbox)\n", len(alerts)-7)
			break
		}
		fmt.Printf("  [%s] %s\n", al.Kind, clipStr(al.Body, 110))
	}

	// Repo-level degradations (I2: loud lines, never silence).
	if v := a.db.Get("sync-error:" + repo); v != "" {
		fmt.Printf("!! sync: %s\n", clipStr(v, 120))
	}
	if v := a.db.Get("overfire:" + repo); v != "" {
		fmt.Printf("!! %s\n", v)
	}
	if len(j.LoadErrors) > 0 {
		fmt.Printf("!! %d journal file(s) failed to parse (never guessed): %s\n", len(j.LoadErrors), firstN(j.LoadErrors, 2))
	}
	return nil
}

func statusMachine(a *app) error {
	fmt.Println("\n── machine ──")
	p, note := a.provider()
	observed, spent := a.db.TokensToday("observed"), a.db.TokensToday("spent")
	pctBudget := int(a.cfg.Extractor.SessionPct / 100 * float64(observed))
	fmt.Printf("extraction: %s · spent %s of min(%s [2%% rule], %s [daily cap]) · observed %s today\n",
		providerName(p, note), kTok(spent), kTok(pctBudget), kTok(a.cfg.Extractor.DailyCapTokens), kTok(observed))
	if v := a.db.Get("extract-paused"); v != "" {
		fmt.Printf("!! extraction paused: %s\n", v)
	}
	if n := a.db.ParkedCount(); n > 0 {
		fmt.Printf("!! %d parked slice(s) awaiting a smarter parser: %s\n", n, strings.Join(a.db.ParkedRecent(2), "; "))
	}
	if v := a.db.Get("redactions"); v != "" && v != "0" {
		fmt.Printf("redactions applied to quotes/bodies: %s\n", v)
	}
	if adapters.DesktopStorePresent() {
		fmt.Println("cursor desktop store detected: not parsed in v1 (CLI transcripts only — lower fidelity, §5.1)")
	}
	// Cloud-session honesty (§11): visible gap, never a silent one.
	fmt.Println("cloud sessions (claude web / codex cloud): no sensor in v1 — those sessions are a visible gap here")
	return nil
}

func kTok(n int) string {
	if n >= 1000 {
		return fmt.Sprintf("%.1fk", float64(n)/1000)
	}
	return fmt.Sprintf("%d", n)
}

func firstN(s []string, n int) []string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}
