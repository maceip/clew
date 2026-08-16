package main

import (
	"errors"
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
	pid, adopted, err := claimWatchLock(lock)
	if err != nil {
		return err
	}
	if adopted {
		fmt.Printf("watcher already running (pid %d); adopting it\n", pid)
		return nil
	}
	defer os.Remove(lock)
	if err := migrateLiveCursors(a); err != nil {
		return err
	}

	w := newWatcher(a)
	if w.provider == nil {
		a.db.Set("watch-provider-error", w.providerNote)
	} else {
		a.db.Set("watch-provider-error", "")
	}
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
	retryAt  time.Time
	failures int
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
	// A file first discovered after the one-time enrollment/migration baseline
	// is live-only; explicit backfill must never replay it. Atomic enrollment
	// prevents a concurrent backfill from assigning half of the boundaries.
	if _, err := db.InitWatermarks(
		state.WatermarkInit{File: "tail:" + file, Adapter: ad.ID(), Repo: repo, Offset: 0},
		state.WatermarkInit{File: "extract:" + file, Adapter: ad.ID(), Repo: repo, Offset: 0},
		state.WatermarkInit{File: "history-end:" + file, Adapter: ad.ID(), Repo: repo, Offset: 0},
		state.WatermarkInit{File: "backfill:" + file, Adapter: ad.ID(), Repo: repo, Offset: 0},
	); err != nil {
		db.Set("adapter-error:"+file, err.Error())
		return
	}
	off := db.Watermark("tail:" + file)
	if fi.Size() == off {
		return
	}
	truncated := fi.Size() < off
	d, err := ad.Parse(file, off)
	if ferr, ok := err.(*adapters.FormatError); ok {
		parkEnd := off
		if d != nil {
			parkEnd = d.NewOffset
		}
		if parkEnd <= off {
			parkEnd, _ = adapters.CompleteOffset(file)
		}
		if parkEnd > off {
			if parkErr := extract.ParkRawRange(db, ad, file, off, parkEnd, ferr.Detail); parkErr != nil {
				db.Set("adapter-error:"+file, "format break; raw park failed: "+parkErr.Error())
			}
		}
		db.Set("adapter-paused:"+file, ferr.Detail)
		db.UpsertAlert(state.Alert{
			Key: "adapter:" + file, RepoPath: repo, Kind: "adapter",
			Body:     fmt.Sprintf("adapter %s paused: %s (file %s) — update clew or report the format", ad.ID(), ferr.Detail, file),
			Blocking: true,
		})
		return
	}
	if err != nil {
		w.pauseAdapter(repo, ad, file, off, "parse failed: "+err.Error())
		return
	}
	if d == nil {
		w.pauseAdapter(repo, ad, file, off, "parser returned no delta")
		return
	}
	unknown := 0
	for cls, n := range d.Unknown {
		unknown += n
		for i := 0; i < n; i++ {
			db.Incr("unknown:" + ad.ID() + ":" + cls)
		}
	}
	if unknown > 0 {
		reason := fmt.Sprintf("%d unrecognized %s record(s)", unknown, ad.ID())
		if err := extract.ParkRawRange(db, ad, file, off, d.NewOffset, reason); err != nil {
			db.Set("adapter-error:"+file, "unknown records; raw park failed: "+err.Error())
			return
		}
	}
	if err := db.SetWatermark("tail:"+file, ad.ID(), repo, d.NewOffset); err != nil {
		db.Set("adapter-error:"+file, "tail watermark failed: "+err.Error())
		return
	}
	if truncated {
		if exOff := db.Watermark("extract:" + file); exOff > d.NewOffset {
			if err := db.SetWatermark("extract:"+file, ad.ID(), repo, 0); err != nil {
				db.Set("extract-error:"+file, "rotation cursor reset failed: "+err.Error())
				return
			}
		}
	}
	db.Set("adapter-error:"+file, "")
	db.AddTokens("observed", d.Bytes/4)
	if d.Bytes == 0 {
		return
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
	startedAt, lastActivity := deltaTimes(d, fi.ModTime())
	db.UpsertSession(state.Session{
		ID: ad.ID() + ":" + d.SessionID, Adapter: ad.ID(), Agent: d.Agent,
		File: file, RepoPath: repo, Surface: w.a.cfg.Surface, Title: title,
		StartedAt: startedAt, LastActivity: lastActivity,
	})
	if len(d.Footprints) > 0 {
		db.AddFootprints(ad.ID()+":"+d.SessionID, d.Footprints)
	}
}

func (w *watcher) pauseAdapter(repo string, ad adapters.Adapter, file string, off int64, detail string) {
	end, _ := adapters.CompleteOffset(file)
	if end > off {
		if err := extract.ParkRawRange(w.a.db, ad, file, off, end, detail); err != nil {
			detail += "; raw park failed: " + err.Error()
		}
	}
	w.a.db.Set("adapter-paused:"+file, detail)
	w.a.db.Set("adapter-error:"+file, detail)
	w.a.db.UpsertAlert(state.Alert{
		Key: "adapter:" + file, RepoPath: repo, Kind: "adapter",
		Body: fmt.Sprintf("adapter %s paused: %s (file %s)", ad.ID(), detail, file), Blocking: true,
	})
}

const liveCursorMigration = "migration:live-cursors-v1"

func migrateLiveCursors(a *app) error {
	return migrateLiveCursorsDB(a.db, adapters.All())
}

func migrateLiveCursorsDB(db *state.DB, all []adapters.Adapter) error {
	if db.Get(liveCursorMigration) != "" {
		return nil
	}
	repos, err := db.Repos()
	if err != nil {
		return err
	}
	legacySkipped := 0
	for _, repo := range repos {
		for _, adapter := range all {
			for _, file := range adapter.Discover(repo.Path) {
				// A short-lived dogfood build used watch-extract: for the live
				// cursor while extract: remained backfill. Reconcile it explicitly
				// so upgrading that build cannot replay or duplicate the live tail.
				if legacyLive, legacyExists := db.WatermarkOK("watch-extract:" + file); legacyExists {
					historyEnd, historicalExists := db.WatermarkOK("extract:" + file)
					if !historicalExists {
						historyEnd = legacyLive
						legacySkipped++
					}
					// Never rewind a cursor that already consumed farther than the
					// short-lived split cursor. In that build extract: could advance
					// through explicit backfill while watch-extract: lagged.
					liveOffset := legacyLive
					if historicalExists && historyEnd > liveOffset {
						liveOffset = historyEnd
					}
					if _, tailExists := db.WatermarkOK("tail:" + file); !tailExists {
						if _, err := db.InitWatermarks(state.WatermarkInit{File: "tail:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: liveOffset}); err != nil {
							return err
						}
					}
					if _, err := db.InitWatermarks(
						state.WatermarkInit{File: "history-end:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: historyEnd},
						state.WatermarkInit{File: "backfill:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: historyEnd},
					); err != nil {
						return err
					}
					if err := db.SetWatermark("extract:"+file, adapter.ID(), repo.Path, liveOffset); err != nil {
						return err
					}
					continue
				}
				liveOffset, exists := db.WatermarkOK("extract:" + file)
				historyEnd, backfillOffset := liveOffset, liveOffset
				var inits []state.WatermarkInit
				if !exists {
					var err error
					liveOffset, err = adapters.CompleteOffset(file)
					if err != nil {
						return fmt.Errorf("migrate session %s: %w", file, err)
					}
					historyEnd, backfillOffset = liveOffset, 0
				}
				if _, tailExists := db.WatermarkOK("tail:" + file); !tailExists {
					inits = append(inits, state.WatermarkInit{File: "tail:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: liveOffset})
				}
				inits = append(inits,
					state.WatermarkInit{File: "extract:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: liveOffset},
					state.WatermarkInit{File: "history-end:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: historyEnd},
					state.WatermarkInit{File: "backfill:" + file, Adapter: adapter.ID(), Repo: repo.Path, Offset: backfillOffset},
				)
				if _, err := db.InitWatermarks(inits...); err != nil {
					return err
				}
			}
		}
	}
	if legacySkipped > 0 {
		if err := db.Set("migration-note", fmt.Sprintf("%d legacy dogfood session file(s) marked historically consumed to prevent duplicate replay", legacySkipped)); err != nil {
			return err
		}
	}
	return db.Set(liveCursorMigration, time.Now().UTC().Format(time.RFC3339))
}

func deltaTimes(d *adapters.Delta, fallback time.Time) (time.Time, time.Time) {
	var started, last time.Time
	for _, message := range d.Messages {
		if message.At.IsZero() {
			continue
		}
		if started.IsZero() || message.At.Before(started) {
			started = message.At
		}
		if last.IsZero() || message.At.After(last) {
			last = message.At
		}
	}
	if started.IsZero() {
		started = fallback
	}
	if last.IsZero() {
		last = fallback
	}
	return started.UTC(), last.UTC()
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
	exOff, tailOff, rotated, err := liveExtractionOffsets(db, ad, repo, file)
	if err != nil {
		db.Set("extract-error:"+file, err.Error())
		return
	}
	pending := tailOff - exOff
	if pending <= 0 {
		return
	}
	p := w.pend[file]
	if p == nil {
		p = &pendState{lastSize: fi.Size(), lastGrew: time.Now()}
		w.pend[file] = p
	}
	if time.Now().Before(p.retryAt) {
		return
	}
	if fi.Size() > p.lastSize {
		p.lastSize = fi.Size()
		p.lastGrew = time.Now()
	}
	// §6.1 triggers: idle ≥ 2 min OR ≥ 50 KB new OR rotation.
	if !rotated && pending < extract.BytesTrigger && time.Since(p.lastGrew) < extract.IdleTrigger {
		return
	}
	j := w.journal(repo)
	if j == nil {
		w.failExtraction(file, p, fmt.Errorf("journal unavailable"))
		return
	}
	metered := newBudgetedProvider(w.provider, db, w.a.cfg, "extraction", true, 0)
	out, err := extract.Run(j, metered, ad, file, exOff, w.a.cfg.Surface, time.Now())
	if err != nil {
		var limit *state.LLMBudgetLimitError
		if errors.As(err, &limit) {
			db.Set("extract-paused", limit.Error())
			db.UpsertAlert(state.Alert{
				Key: "budget:" + time.Now().UTC().Format("2006-01-02"), RepoPath: repo,
				Kind: "budget", Body: "extraction paused: " + limit.Error() + " — sensors keep recording; catch-up is automatic (I9)",
				Blocking: false,
			})
			return
		}
		w.failExtraction(file, p, err)
		return
	}
	db.Set("extract-paused", "")
	if out.Parked {
		if err := extract.ParkRawRange(db, ad, file, exOff, tailOff, out.ParkReason); err != nil {
			w.failExtraction(file, p, fmt.Errorf("park failed: %w", err))
			return
		}
		if err := db.SetWatermark("extract:"+file, ad.ID(), repo, tailOff); err != nil {
			w.failExtraction(file, p, fmt.Errorf("watermark failed after park: %w", err))
			return
		}
	} else {
		if err := db.SetWatermark("extract:"+file, ad.ID(), repo, out.NewOffset); err != nil {
			w.failExtraction(file, p, fmt.Errorf("watermark failed: %w", err))
			return
		}
	}
	db.Set("extract-error:"+file, "")
	for i := 0; i < out.Redactions; i++ {
		db.Incr("redactions")
	}
	delete(w.pend, file)
	if len(out.Entries) > 0 || len(out.Events) > 0 {
		w.syncDue[repo] = time.Now().Add(5 * time.Second) // §4 debounce
	}
}

func liveExtractionOffsets(db *state.DB, ad adapters.Adapter, repo, file string) (exOff, tailOff int64, rotated bool, err error) {
	exOff = db.Watermark("extract:" + file)
	tailOff = db.Watermark("tail:" + file)
	if exOff <= tailOff {
		return exOff, tailOff, false, nil
	}
	if err := db.SetWatermark("extract:"+file, ad.ID(), repo, 0); err != nil {
		return exOff, tailOff, true, fmt.Errorf("rotation cursor reset failed: %w", err)
	}
	return 0, tailOff, true, nil
}

func (w *watcher) failExtraction(file string, p *pendState, err error) {
	p.failures++
	shift := p.failures - 1
	if shift > 6 {
		shift = 6
	}
	delay := 5 * time.Second * time.Duration(1<<shift)
	p.retryAt = time.Now().Add(delay)
	w.a.db.Set("extract-error:"+file, fmt.Sprintf("%v; retry in %s", err, delay))
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
	differProvider := w.provider
	linkPass := w.a.cfg.LinkPass
	if linkPass && differProvider != nil {
		differProvider = newBudgetedProvider(differProvider, db, w.a.cfg, "differ", false, 0)
	}
	res, err := differ.Run(db, &differ.Input{
		Repo: repo, Journal: j, Snapshot: snap, Surface: w.a.cfg.Surface,
		Provider: differProvider, LinkPass: linkPass,
	}, time.Now())
	if err != nil {
		db.Set("differ-error:"+repo, err.Error())
	} else {
		db.Set("differ-error:"+repo, "")
	}
	w.syncRepo(repo) // ≤30s fetch+rebase interval (§4)
	if res != nil {
		for _, al := range res.NewAlerts {
			if !al.Blocking {
				continue
			}
			sent, err := push.Send(w.a.cfg.Push, "clew: "+repoBase(repo), al.Body)
			if err != nil {
				db.Set("push-error:"+repo, err.Error())
			} else {
				db.Set("push-error:"+repo, "")
				if sent {
					db.MarkAlert(al.Key, "pushed_at")
				}
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
		launchPath := supervisorPATH(bin)
		plist := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0"><dict>
  <key>Label</key><string>dev.clew.watch</string>
  <key>ProgramArguments</key><array><string>%s</string><string>watch</string></array>
  <key>EnvironmentVariables</key><dict><key>PATH</key><string>%s</string></dict>
  <key>RunAtLoad</key><true/>
  <key>KeepAlive</key><true/>
  <key>StandardOutPath</key><string>%s/watch.log</string>
  <key>StandardErrorPath</key><string>%s/watch.log</string>
</dict></plist>
`, xmlText(bin), xmlText(launchPath), xmlText(logDir), xmlText(logDir))
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

func supervisorPATH(bin string) string {
	path := os.Getenv("PATH")
	if path == "" {
		path = "/usr/local/bin:/usr/bin:/bin:/usr/sbin:/sbin"
	}
	binDir := filepath.Dir(bin)
	for _, dir := range filepath.SplitList(path) {
		if dir == binDir {
			return path
		}
	}
	return binDir + string(os.PathListSeparator) + path
}

func xmlText(s string) string {
	return strings.NewReplacer("&", "&amp;", "<", "&lt;", ">", "&gt;").Replace(s)
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

// claimWatchLock atomically owns the watcher/migration boundary. O_EXCL closes
// the check-then-write race between two watcher starts and between a legacy
// cursor migration and explicit backfill.
func claimWatchLock(lock string) (int, bool, error) {
	for attempt := 0; attempt < 2; attempt++ {
		f, err := os.OpenFile(lock, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err == nil {
			if _, err := f.WriteString(strconv.Itoa(os.Getpid())); err != nil {
				f.Close()
				os.Remove(lock)
				return 0, false, err
			}
			if err := f.Close(); err != nil {
				os.Remove(lock)
				return 0, false, err
			}
			return 0, false, nil
		}
		if !os.IsExist(err) {
			return 0, false, err
		}
		if pid, ok := lockAlive(lock); ok {
			return pid, true, nil
		}
		if err := os.Remove(lock); err != nil && !os.IsNotExist(err) {
			return 0, false, err
		}
	}
	return 0, false, fmt.Errorf("could not claim watcher lock %s", lock)
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
