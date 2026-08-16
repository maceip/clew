package main

import (
	"bytes"
	"flag"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"clew/internal/docket"
	"clew/internal/gitx"
	"clew/internal/journal"
	"clew/internal/model"
)

type glanceItem struct {
	ID, Title, Age, Status string
	Star                   bool
}

type glanceSection struct {
	Name  string
	Items []glanceItem
}

type glanceView struct {
	Repo, Generated, Title string
	Sections               []glanceSection
	Intents, InFlight      int
	Absent, Proposed, Done int
	Docket                 int
}

func cmdGlance(args []string) error {
	fs := flag.NewFlagSet("glance", flag.ContinueOnError)
	htmlOut := fs.Bool("html", false, "write the self-contained pinned-tab glance")
	if err := fs.Parse(args); err != nil {
		return err
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	j, err := journal.Load(gitx.WorktreeDir(repo))
	if err != nil {
		return err
	}
	view := buildGlance(a, repo, j, time.Now())
	if *htmlOut {
		if err := a.db.Set("glance-repo", repo); err != nil {
			return err
		}
		path, err := writeGlanceHTML(view)
		if err != nil {
			return err
		}
		fmt.Println(path)
		return nil
	}
	renderGlanceTerminal(view)
	return nil
}

func buildGlance(a *app, repo string, j *journal.Journal, now time.Time) glanceView {
	st := journal.Compute(j, now)
	view := glanceView{Repo: repoBase(repo), Generated: now.Format("15:04:05")}
	section := func(name string, typ model.EntryType, keep func(journal.Status) bool) {
		var entries []*model.Entry
		for id, entry := range j.Entries {
			if entry.Type == typ && st[id] != nil && keep(st[id].Status) {
				entries = append(entries, entry)
			}
		}
		sort.Slice(entries, func(i, k int) bool { return entries[i].ID > entries[k].ID })
		items := make([]glanceItem, 0, 7)
		for _, entry := range entries {
			if len(items) == 7 {
				break
			}
			items = append(items, glanceItem{
				ID: entry.ID, Title: entry.Title, Age: glanceAge(now, entry.Created()),
				Status: string(st[entry.ID].Status), Star: entry.Type == model.Question && entry.Asks == "human",
			})
		}
		view.Sections = append(view.Sections, glanceSection{Name: name, Items: items})
	}
	section("DECIDED", model.Decision, func(s journal.Status) bool {
		return s == journal.StActive || s == journal.StPossibleContradiction || s == journal.StContradicted
	})
	section("LEARNED", model.Finding, func(s journal.Status) bool { return s == journal.StCurrent || s == journal.StSuspect })
	section("OPEN", model.Question, func(s journal.Status) bool { return s == journal.StOpen })
	for id, entry := range j.Entries {
		if entry.Type != model.Intent || st[id] == nil || !journal.Live(st[id].Status) {
			continue
		}
		view.Intents++
		switch st[id].Status {
		case journal.StInFlight:
			view.InFlight++
		case journal.StAbsent:
			view.Absent++
		case journal.StProposed:
			view.Proposed++
		case journal.StDone:
			view.Done++
		}
	}
	alerts := a.db.OpenAlerts(repo, true)
	cards := docket.Build(docket.Input{
		Journal: j, Alerts: alerts, Now: now, Assumptions: docketAssumptions(alerts),
	})
	view.Docket = len(cards)
	view.Title = "clew"
	if view.Docket > 0 {
		view.Title = "clew ●"
	}
	return view
}

func renderGlanceTerminal(view glanceView) {
	fmt.Printf("clew · %s · %s\n", view.Repo, view.Generated)
	for _, section := range view.Sections {
		fmt.Printf("%s %d\n", section.Name, len(section.Items))
		if len(section.Items) == 0 {
			fmt.Println("  —")
		}
		for _, item := range section.Items {
			star := ""
			if item.Star {
				star = " ★"
			}
			fmt.Printf("  %s · %s · %s%s\n", item.Title, item.Age, item.Status, star)
		}
	}
	fmt.Printf("MAP %d intents · %d in flight · %d absent · %d proposed · %d done\n",
		view.Intents, view.InFlight, view.Absent, view.Proposed, view.Done)
	fmt.Printf("DOCKET %d\n", view.Docket)
}

func glanceAge(now, then time.Time) string {
	d := now.Sub(then)
	if d < 0 {
		d = 0
	}
	if d < time.Hour {
		return fmt.Sprintf("%dm", int(d/time.Minute))
	}
	if d < 24*time.Hour {
		return fmt.Sprintf("%dh", int(d/time.Hour))
	}
	return fmt.Sprintf("%dd", int(d/(24*time.Hour)))
}

var glanceHTML = template.Must(template.New("glance").Funcs(template.FuncMap{
	"lower": strings.ToLower,
}).Parse(`<!doctype html>
<html lang="en"><head><meta charset="utf-8"><meta name="viewport" content="width=device-width,initial-scale=1">
<meta http-equiv="refresh" content="30"><title>{{.Title}}</title>
<style>
:root{color-scheme:light dark;--paper:#f6f3eb;--ink:#17211d;--muted:#66706b;--line:#d8d3c7;--amber:#d88b17;--alert:#a53b2a}
@media(prefers-color-scheme:dark){:root{--paper:#151917;--ink:#e9ece7;--muted:#9ba39e;--line:#343a36;--amber:#f0a72b;--alert:#ed7664}}
*{box-sizing:border-box}body{margin:0;background:var(--paper);color:var(--ink);font:15px/1.4 ui-monospace,SFMono-Regular,Menlo,monospace}
main{max-width:980px;margin:auto;padding:28px}header{display:flex;justify-content:space-between;align-items:baseline;border-bottom:1px solid var(--line);padding-bottom:14px}
h1{font:600 22px/1.2 system-ui;margin:0}.light{color:var(--amber)}.quiet{color:var(--muted);font-size:12px}.grid{display:grid;grid-template-columns:repeat(3,1fr);gap:16px;margin-top:18px}
section,.strip{border:1px solid var(--line);border-radius:10px;padding:14px;background:color-mix(in srgb,var(--paper) 92%,var(--ink) 8%)}h2{font-size:12px;letter-spacing:.12em;margin:0 0 10px;color:var(--muted)}
ul{list-style:none;margin:0;padding:0}li{padding:7px 0;border-top:1px solid var(--line)}li:first-child{border-top:0}.meta{color:var(--muted);font-size:12px}.star,.absent{color:var(--alert)}
.strip{margin-top:16px;display:flex;gap:22px;flex-wrap:wrap}.docket{margin-left:auto}.docket.hot{color:var(--amber);font-weight:700}@media(max-width:720px){.grid{grid-template-columns:1fr}main{padding:16px}}
</style></head><body><main><header><h1>clew {{if .Docket}}<span class="light">●</span>{{end}} · {{.Repo}}</h1><span class="quiet">updated {{.Generated}} · refresh 30s</span></header>
<div class="grid">{{range .Sections}}<section><h2>{{.Name}} · {{len .Items}}</h2><ul>{{if not .Items}}<li class="quiet">Nothing current</li>{{end}}{{range .Items}}<li>{{.Title}}{{if .Star}} <span class="star">★</span>{{end}}<div class="meta">{{.Age}} · {{.Status}} · {{.ID}}</div></li>{{end}}</ul></section>{{end}}</div>
<div class="strip"><span>MAP · {{.Intents}} intents</span><span>{{.InFlight}} in flight</span><span class="{{if .Absent}}absent{{end}}">{{.Absent}} absent</span><span>{{.Proposed}} proposed</span><span>{{.Done}} done</span><span class="docket {{if .Docket}}hot{{end}}">DOCKET · {{.Docket}}</span></div>
</main></body></html>`))

func writeGlanceHTML(view glanceView) (string, error) {
	var rendered bytes.Buffer
	if err := glanceHTML.Execute(&rendered, view); err != nil {
		return "", err
	}
	dir := gitx.Home()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp(dir, "glance-*.html")
	if err != nil {
		return "", err
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.Write(rendered.Bytes()); err != nil {
		tmp.Close()
		return "", err
	}
	if err := tmp.Close(); err != nil {
		return "", err
	}
	path := filepath.Join(dir, "glance.html")
	if err := os.Rename(tmpPath, path); err != nil {
		return "", err
	}
	return path, nil
}
