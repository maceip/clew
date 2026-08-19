package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/maceip/clew/internal/calm"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
)

type mapRow struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Tags     string `json:"tags"`
	Status   string `json:"status"`
	Evidence int    `json:"evidence"`
	LastAct  string `json:"last_activity"`
	Session  string `json:"session"`
}

// cmdMap renders the reflexion view (§8.2): one row per live intent, absence
// highlighted; expired questions and unmapped commits listed (displayed, not
// guessed). --html emits a single self-contained page.
func cmdMap(args []string) error {
	fs := flag.NewFlagSet("map", flag.ExitOnError)
	htmlOut := fs.String("html", "", "write a self-contained HTML page")
	fs.Parse(args)
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
	now := time.Now()
	st := journal.Compute(j, now)

	var rows []mapRow
	for id, e := range j.Entries {
		if e.Type != model.Intent {
			continue
		}
		c := st[id]
		if c.Status == journal.StDropped {
			continue
		}
		attributed := ""
		for _, v := range j.EventsFor(id) {
			if v.Kind == model.EvEvidence {
				if s := v.PStr("session"); s != "" {
					attributed = s
				}
			}
		}
		rows = append(rows, mapRow{
			ID: id, Title: calm.Text(e.Title), Tags: strings.Join(e.Tags, ","),
			Status: calm.Text(string(c.Status)), Evidence: c.Evidence,
			LastAct: c.LastActivity.Format("2006-01-02"), Session: attributed,
		})
	}
	sort.Slice(rows, func(a, b int) bool {
		// absent first (that is the point of the map), then newest.
		if (rows[a].Status == "absent") != (rows[b].Status == "absent") {
			return rows[a].Status == "absent"
		}
		return rows[a].ID > rows[b].ID
	})

	var expired []*model.Entry
	for id, e := range j.Entries {
		if e.Type == model.Question && st[id].Status == journal.StExpired {
			expired = append(expired, e)
		}
	}
	var unmapped []string
	json.Unmarshal([]byte(a.db.Get("unmapped:"+repo)), &unmapped)

	if *htmlOut != "" {
		if err := writeMapHTML(*htmlOut, repoBase(repo), rows, expired, unmapped, now); err != nil {
			return err
		}
		fmt.Println("wrote", *htmlOut)
		return nil
	}

	fmt.Printf("INTENT × REALITY — %s (%d intents)\n\n", repoBase(repo), len(rows))
	wTitle := 44
	fmt.Printf("%-28s %-*s %-22s %-9s %3s  %-10s %s\n", "id", wTitle, "title", "tags", "status", "ev", "last", "session")
	for _, r := range rows {
		status := r.Status
		if r.Status == "absent" {
			status = "◀ ABSENT" // the flagship signal (§0, §7.1.3)
		}
		fmt.Printf("%-28s %-*s %-22s %-9s %3d  %-10s %s\n",
			r.ID, wTitle, clipStr(r.Title, wTitle), clipStr(r.Tags, 22), status, r.Evidence, r.LastAct, r.Session)
	}
	if len(expired) > 0 {
		fmt.Printf("\nEXPIRED QUESTIONS (%d — aged out after 45d, §3.3)\n", len(expired))
		for _, e := range expired {
			fmt.Printf("  %s %s\n", e.ID, calm.Text(e.Title))
		}
	}
	if len(unmapped) > 0 {
		fmt.Printf("\nUNMAPPED COMMITS (%d — no intent claims them; unmapped is itself signal, §7.1)\n", len(unmapped))
		for i, sha := range unmapped {
			if i >= 10 {
				fmt.Printf("  … %d more\n", len(unmapped)-10)
				break
			}
			fmt.Printf("  %s\n", sha[:10])
		}
	}
	for _, al := range a.db.OpenAlerts(repo, false) {
		if al.Kind == "overlap" {
			fmt.Printf("\nOVERLAP: %s\n", calm.Text(al.Body))
		}
	}
	return nil
}

func writeMapHTML(path, repo string, rows []mapRow, expired []*model.Entry, unmapped []string, now time.Time) error {
	rowsJSON, _ := json.Marshal(rows)
	var exp []string
	for _, e := range expired {
		exp = append(exp, e.ID+" "+e.Title)
	}
	expJSON, _ := json.Marshal(exp)
	unmJSON, _ := json.Marshal(unmapped)
	page := `<!doctype html><meta charset="utf-8"><title>clew map — ` + repo + `</title>
<style>
body{font:14px/1.5 -apple-system,system-ui,sans-serif;margin:2rem;color:#1a1a1a;max-width:1100px}
h1{font-size:1.2rem} table{border-collapse:collapse;width:100%}
th,td{text-align:left;padding:.35rem .6rem;border-bottom:1px solid #eee;font-size:.85rem}
tr.absent{background:#fff2f0} tr.absent td.status{color:#c0392b;font-weight:700}
tr.in_flight td.status{color:#2e7d32} .muted{color:#777} code{background:#f5f5f5;padding:0 .3em}
</style>
<h1>intent × reality — ` + repo + ` <span class="muted">` + now.UTC().Format("2006-01-02 15:04 MST") + `</span></h1>
<table id="t"><thead><tr><th>id</th><th>title</th><th>tags</th><th>status</th><th>evidence</th><th>last activity</th><th>session</th></tr></thead><tbody></tbody></table>
<div id="extra"></div>
<script>
const rows=` + string(rowsJSON) + `,expired=` + string(expJSON) + `,unmapped=` + string(unmJSON) + `;
const tb=document.querySelector('#t tbody');
for(const r of rows){const tr=document.createElement('tr');tr.className=r.status;
tr.innerHTML='<td><code>'+r.id.slice(0,10)+'…</code></td><td>'+esc(r.title)+'</td><td>'+esc(r.tags)+'</td>'+
'<td class="status">'+(r.status==='absent'?'◀ ABSENT':r.status)+'</td><td>'+r.evidence+'</td><td>'+r.last_activity+'</td><td>'+esc(r.session||'')+'</td>';
tb.appendChild(tr);}
const ex=document.querySelector('#extra');let h='';
if(expired.length){h+='<h1>expired questions</h1><ul>'+expired.map(q=>'<li>'+esc(q)+'</li>').join('')+'</ul>';}
if(unmapped.length){h+='<h1>unmapped commits <span class=muted>(no intent claims them)</span></h1><ul>'+unmapped.map(c=>'<li><code>'+esc(c.slice(0,10))+'</code></li>').join('')+'</ul>';}
ex.innerHTML=h;
function esc(s){return String(s).replace(/[&<>"]/g,c=>({'&':'&amp;','<':'&lt;','>':'&gt;','"':'&quot;'}[c]))}
</script>`
	return os.WriteFile(path, []byte(page), 0o644)
}
