package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"

	"clew/internal/adapters"
	"clew/internal/differ"
	"clew/internal/extract"
	"clew/internal/gitx"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/materialize"
	"clew/internal/poller"
	"clew/internal/push"
	"clew/internal/state"
)

func cmdWatch(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "install":
			return watchInstall()
		case "uninstall":
			return watchUninstall()
		}
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()

	// Idempotent adoption (§2): one watcher per machine.
	lock := filepath.Join(gitx.Home(), "watch.lock")
	if pid, ok := lockAlive(lock); ok {
		fmt.Printf("watcher already running (pid %d); adopting it\n", pid)
		return nil
	}
	if err := os.WriteFile(lock, []byte(strconv.Itoa(os.Getpid())), 0o644); err != nil {
		return err
	}
	defer os.Remove(lock)

	w := newWatcher(a)
	fmt.Printf("clew watcher: surface=%s provider=%s repos=%d\n",
		a.cfg.Surface, providerName(w.provider, w.providerNote), reposLen(a.db))

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	tickTail := time.NewTicker(2 * time.Second)
	tickPoll := time.NewTicker(30 * time.Second)
	defer tickTail.Stop()
	defer tickPoll.Stop()

	w.pollAll() // prime once at startup
	for {
		select {
		case <-stop:
			fmt.Println("watcher stopping")
			return nil
		case <-tickTail.C:
			w.tailAll()
			w.dueSyncs()
		case <-tickPoll.C:
			w.pollAll()
		}
	}
}

type pendState struct {
	lastSize int64
	lastGrew time.Time
}

type watcher struct {
	a            *app
	provider     llm.Provider
	providerNote string
	pend         map[string]*pendState
	syncDue      map[string]time.Time
	journals     map[string]*journal.Journal
}

func newWatcher(a *app) *watcher {
	p, note := a.provider()
	return &watcher{
		a: a, provider: p, providerNote: note,
		pend:     map[string]*pendState{},
		syncDue:  map[string]time.Time{},
		journals: map[string]*journal.Journal{},
	}
}

func (w *watcher) journal(repo string) *journal.Journal {
	if j, ok := w.journals[repo]; ok {
		return j
	}
	j, err := w.a.openJournal(repo)
	if err != nil {
		return nil
	}
	w.journals[repo] = j
	return j
}

func (w *watcher) repos() []string {
	rs, _ := w.a.db.Repos()
	var out []string
	for _, r := range rs {
		out = append(out, r.Path)
	}
	return out
}

// ---- 2s tick: tails + extraction triggers ----

func (w *watcher) tailAll() {
	for _, repo := range w.repos() {
		for _, ad := range adapters.All() {
			for _, file := range ad.Discover(repo) {
				w.tail(repo, ad, file)
				w.maybeExtract(repo, ad, file)
			}
		}
	}
}

func (w *watcher) tail(repo string, ad adapters.Adapter, file string) {
	db := w.a.db
	if db.Get("adapter-paused:"+file) != "" {
		return // paused loudly until format is re-pinned (I2)
	}
	fi, err := os.Stat(file)
	if err != nil {
		return
	}
	off := db.Watermark("tail:" + file)
	if fi.Size() == off {
		return
	}
	d, err := ad.Parse(file, off)
	if ferr, ok := err.(*adapters.FormatError); ok {
		db.Set("adapter-paused:"+file, ferr.Detail)
		db.UpsertAlert(state.Alert{
			Key: "adapter:" + file, RepoPath: repo, Kind: "adapter",
			Body:     fmt.Sprintf("adapter %s paused: %s (file %s) — update clew or report the format", ad.ID(), ferr.Detail, file),
			Blocking: true,
		})
		return
	}
	if err != nil || d == nil {
		return
	}
	db.SetWatermark("tail:"+file, ad.ID(), repo, d.NewOffset)
	db.AddTokens("observed", d.Bytes/4)
	for cls, n := range d.Unknown {
		for i := 0; i < n; i++ {
			db.Incr("unknown:" + ad.ID() + ":" + cls)
		}
	}
	if d.SessionID == "" {
		d.SessionID = strings.TrimSuffix(filepath.Base(file), ".jsonl")
	}
	title := ""
	for _, m := range d.Messages {
		if m.Role == "user" {
			title = clipStr(m.Text, 60)
			break
		}
	}
	db.UpsertSession(state.Session{
		ID: ad.ID() + ":" + d.SessionID, Adapter: ad.ID(), Agent: d.Agent,
		File: file, RepoPath: repo, Surface: w.a.cfg.Surface, Title: title,
		StartedAt: time.Now(), LastActivity: time.Now(),
	})
	if len(d.Footprints) > 0 {
		db.AddFootprints(ad.ID()+":"+d.SessionID, d.Footprints)
	}
}

func (w *watcher) maybeExtract(repo string, ad adapters.Adapter, file string) {
	if w.provider == nil {
		return // loud in status; sensors keep recording
	}
	db := w.a.db
	fi, err := os.Stat(file)
	if err != nil {
		return
	}
	exOff := db.Watermark("extract:" + file)
	pending := fi.Size() - exOff
	if pending <= 0 {
		return
	}
	p := w.pend[file]
	if p == nil {
		p = &pendState{lastSize: fi.Size(), lastGrew: time.Now()}
		w.pend[file] = p
	}
	if fi.Size() > p.lastSize {
		p.lastSize = fi.Size()
		p.lastGrew = time.Now()
	}
	// §6.1 triggers: idle ≥ 2 min OR ≥ 50 KB new OR rotation.
	rotated := exOff > fi.Size()
	if !rotated && pending < extract.BytesTrigger && time.Since(p.lastGrew) < extract.IdleTrigger {
		return
	}
	est := int(minI64(pending, extract.SliceCap))/4 + 1500
	if ok, reason := extract.Gate(db, w.a.cfg, est); !ok {
		db.Set("extract-paused", reason)
		db.UpsertAlert(state.Alert{
			Key: "budget:" + time.Now().UTC().Format("2006-01-02"), RepoPath: repo,
			Kind: "budget", Body: "extraction paused: " + reason + " — sensors keep recording; catch-up is automatic (I9)",
			Blocking: false,
		})
		return
	}
	db.Set("extract-paused", "")
	j := w.journal(repo)
	if j == nil {
		return
	}
	out, err := extract.Run(j, w.provider, ad, file, exOff, w.a.cfg.Surface, time.Now())
	if err != nil {
		return // transient provider failure: retry on next trigger
	}
	db.AddTokens("spent", out.Tokens)
	if out.Parked {
		extract.ParkSlice(db, ad, file, exOff, out.ParkReason)
		db.SetWatermark("extract:"+file, ad.ID(), repo, fi.Size())
	} else {
		db.SetWatermark("extract:"+file, ad.ID(), repo, out.NewOffset)
	}
	for i := 0; i < out.Redactions; i++ {
		db.Incr("redactions")
	}
	delete(w.pend, file)
	if len(out.Entries) > 0 || len(out.Events) > 0 {
		w.syncDue[repo] = time.Now().Add(5 * time.Second) // §4 debounce
	}
}

// ---- 30s tick: poller + differ + sync + materialize + push ----

func (w *watcher) pollAll() {
	for _, repo := range w.repos() {
		w.pollOne(repo)
	}
}

func (w *watcher) pollOne(repo string) {
	db := w.a.db
	snap, _ := poller.Poll(db, repo)
	j := w.journal(repo)
	if j == nil {
		return
	}
	res, err := differ.Run(db, &differ.Input{
		Repo: repo, Journal: j, Snapshot: snap, Surface: w.a.cfg.Surface,
		Provider: w.provider, LinkPass: w.a.cfg.LinkPass,
	}, time.Now())
	if err == nil && res.Tokens > 0 {
		db.AddTokens("spent", res.Tokens)
	}
	w.syncRepo(repo) // ≤30s fetch+rebase interval (§4)
	if res != nil {
		for _, al := range res.NewAlerts {
			if !al.Blocking {
				continue
			}
			if err := push.Send(w.a.cfg.Push, "clew: "+repoBase(repo), al.Body); err == nil {
				db.MarkAlert(al.Key, "pushed_at")
			}
		}
	}
}

func (w *watcher) dueSyncs() {
	for repo, due := range w.syncDue {
		if time.Now().After(due) {
			delete(w.syncDue, repo)
			w.syncRepo(repo)
		}
	}
}

func (w *watcher) syncRepo(repo string) {
	db := w.a.db
	res, err := gitx.Sync(repo, regen)
	if err != nil {
		db.Set("sync-error:"+repo, err.Error())
		return
	}
	db.Set("sync-error:"+repo, "")
	if res.Adopted || res.Pulled {
		delete(w.journals, repo) // branch moved: reload
	}
	j := w.journal(repo)
	if j == nil {
		return
	}
	j.Reload()
	now := time.Now()
	st := journal.Compute(j, now)
	materialize.Write(repo, j, st, db, now)
	// Cross-machine ack propagation: disposition events carrying ack keys.
	for _, v := range j.Events {
		if v.Kind == "disposition" {
			if key := v.PStr("ack"); key != "" {
				db.MarkAlert(key, "acked_at")
			}
		}
	}
	if journal.Overfiring(j, st, now) {
		db.Set("overfire:"+repo, "rollup exceeds 32KB — the extractor is over-firing (tune, don't scroll; §3.3)")
	} else {
		db.Set("overfire:"+repo, "")
	}
}

// ---- supervision install (§2: launchd / systemd-user) ----

func watchInstall() error {
	bin, err := os.Executable()
	if err != nil {
		return err
	}
	logDir := filepath.Join(gitx.Home(), "logs")
	os.MkdirAll(logDir, 0o755)
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>dev.clew.watch</string>
  <key>ProgramArguments</key><array><string>%s</string><string>watch</string></array>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
  <key>StandardOutPath</key><string>%s/watch.log</string>
  <key>StandardErrorPath</key><string>%s/watch.log</string>
</dict></plist>
`, bin, logDir, logDir)
		p := filepath.Join(home, "Library", "LaunchAgents", "dev.clew.watch.plist")
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return err
		}
		if err := os.WriteFile(p, []byte(plist), 0o644); err != nil {
			return err
		}
		exec.Command("launchctl", "unload", p).Run()
		if out, err := exec.Command("launchctl", "load", "-w", p).CombinedOutput(); err != nil {
			return fmt.Errorf("launchctl load: %v: %s", err, out)
		}
		fmt.Println("installed launchd agent dev.clew.watch (log:", logDir+"/watch.log)")
		return nil
	case "linux":
		home, _ := os.UserHomeDir()
		unit := fmt.Sprintf(`[Unit]
Description=clew journal watcher

[Service]
ExecStart=%s watch
Restart=always
RestartSec=5

[Install]
WantedBy=default.target
`, bin)
		dir := filepath.Join(home, ".config", "systemd", "user")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}
		p := filepath.Join(dir, "clew-watch.service")
		if err := os.WriteFile(p, []byte(unit), 0o644); err != nil {
			return err
		}
		if out, err := exec.Command("systemctl", "--user", "enable", "--now", "clew-watch.service").CombinedOutput(); err != nil {
			return fmt.Errorf("systemctl: %v: %s (unit written to %s)", err, out, p)
		}
		fmt.Println("installed systemd user unit clew-watch.service")
		return nil
	default:
		return fmt.Errorf("no supervisor template for %s — run `clew watch` under your own supervisor", runtime.GOOS)
	}
}

func watchUninstall() error {
	switch runtime.GOOS {
	case "darwin":
		home, _ := os.UserHomeDir()
		p := filepath.Join(home, "Library", "LaunchAgents", "dev.clew.watch.plist")
		exec.Command("launchctl", "unload", p).Run()
		os.Remove(p)
		fmt.Println("removed launchd agent")
	case "linux":
		exec.Command("systemctl", "--user", "disable", "--now", "clew-watch.service").Run()
		home, _ := os.UserHomeDir()
		os.Remove(filepath.Join(home, ".config", "systemd", "user", "clew-watch.service"))
		fmt.Println("removed systemd user unit")
	}
	return nil
}

func lockAlive(lock string) (int, bool) {
	b, err := os.ReadFile(lock)
	if err != nil {
		return 0, false
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil || pid <= 0 {
		return 0, false
	}
	if err := syscall.Kill(pid, 0); err != nil {
		return 0, false // stale lock
	}
	return pid, true
}

func providerName(p llm.Provider, note string) string {
	if p == nil {
		return "none (" + note + ")"
	}
	return p.Name()
}

func reposLen(db *state.DB) int {
	rs, _ := db.Repos()
	return len(rs)
}

func clipStr(s string, n int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	if n <= 0 {
		return ""
	}
	runes := []rune(s)
	if len(runes) <= n {
		return s
	}
	if n == 1 {
		return "…"
	}
	return string(runes[:n-1]) + "…"
}

func minI64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}
