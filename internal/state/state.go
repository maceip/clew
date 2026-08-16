// Package state is the machine-local working state (~/.clew/state.db,
// JOURNAL_SPEC §2). SQLite, never the source of truth — the journal lives in
// git. Holds: registered repos, tail watermarks, live sessions + footprints,
// seen commits + attribution, alerts, budget meters, parked slices, kv.
package state

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "modernc.org/sqlite"

	"clew/internal/gitx"
)

type DB struct{ *sql.DB }

func DefaultPath() string { return filepath.Join(gitx.Home(), "state.db") }

func Open(path string) (*DB, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}
	db, err := sql.Open("sqlite", path+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)")
	if err != nil {
		return nil, err
	}
	if err := migrate(db); err != nil {
		db.Close()
		return nil, err
	}
	return &DB{db}, nil
}

func migrate(db *sql.DB) error {
	_, err := db.Exec(`
CREATE TABLE IF NOT EXISTS repos(
  path TEXT PRIMARY KEY, remote TEXT, registered_at TEXT);
CREATE TABLE IF NOT EXISTS watermarks(
  file TEXT PRIMARY KEY, adapter TEXT, repo_path TEXT, offset INTEGER, updated_at TEXT);
CREATE TABLE IF NOT EXISTS sessions(
  id TEXT PRIMARY KEY, adapter TEXT, agent TEXT, file TEXT, repo_path TEXT,
  surface TEXT, title TEXT, ctl_sock TEXT, started_at TEXT, last_activity TEXT);
CREATE TABLE IF NOT EXISTS footprints(
  session_id TEXT, path TEXT, at TEXT, PRIMARY KEY(session_id, path));
CREATE TABLE IF NOT EXISTS commits_seen(
  repo_path TEXT, sha TEXT, author TEXT, at TEXT, subject TEXT, files TEXT,
  session_id TEXT, mapped INTEGER DEFAULT 0, PRIMARY KEY(repo_path, sha));
CREATE TABLE IF NOT EXISTS alerts(
  key TEXT PRIMARY KEY, repo_path TEXT, kind TEXT, body TEXT, entry_ids TEXT,
  blocking INTEGER, created_at TEXT, nudged_at TEXT, pushed_at TEXT,
  acked_at TEXT, dropped_at TEXT);
CREATE TABLE IF NOT EXISTS budget(
  day TEXT, kind TEXT, tokens INTEGER, PRIMARY KEY(day, kind));
CREATE TABLE IF NOT EXISTS parked(
  id INTEGER PRIMARY KEY AUTOINCREMENT, adapter TEXT, file TEXT,
  offset INTEGER, reason TEXT, raw_path TEXT, at TEXT);
CREATE TABLE IF NOT EXISTS kv(k TEXT PRIMARY KEY, v TEXT);
`)
	return err
}

func now() string { return time.Now().UTC().Format(time.RFC3339) }

// ---- repos ----

type Repo struct {
	Path   string
	Remote string
}

func (d *DB) RegisterRepo(path, remote string) error {
	_, err := d.Exec(`INSERT INTO repos(path, remote, registered_at) VALUES(?,?,?)
		ON CONFLICT(path) DO UPDATE SET remote=excluded.remote`, path, remote, now())
	return err
}

func (d *DB) Repos() ([]Repo, error) {
	rows, err := d.Query(`SELECT path, remote FROM repos ORDER BY path`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Repo
	for rows.Next() {
		var r Repo
		if err := rows.Scan(&r.Path, &r.Remote); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// RepoFor finds the registered repo containing dir ("" if none).
func (d *DB) RepoFor(dir string) string {
	abs, _ := filepath.Abs(dir)
	repos, _ := d.Repos()
	best := ""
	for _, r := range repos {
		if abs == r.Path || strings.HasPrefix(abs, r.Path+string(os.PathSeparator)) {
			if len(r.Path) > len(best) {
				best = r.Path
			}
		}
	}
	return best
}

// ---- watermarks (§6.1: byte offset per file, persisted) ----

func (d *DB) Watermark(file string) int64 {
	var off int64
	d.QueryRow(`SELECT offset FROM watermarks WHERE file=?`, file).Scan(&off)
	return off
}

func (d *DB) SetWatermark(file, adapter, repo string, offset int64) error {
	_, err := d.Exec(`INSERT INTO watermarks(file, adapter, repo_path, offset, updated_at)
		VALUES(?,?,?,?,?) ON CONFLICT(file) DO UPDATE SET offset=excluded.offset, updated_at=excluded.updated_at`,
		file, adapter, repo, offset, now())
	return err
}

// ---- sessions ----

type Session struct {
	ID, Adapter, Agent, File, RepoPath, Surface, Title, CtlSock string
	StartedAt, LastActivity                                     time.Time
}

func (d *DB) UpsertSession(s Session) error {
	_, err := d.Exec(`INSERT INTO sessions(id, adapter, agent, file, repo_path, surface, title, ctl_sock, started_at, last_activity)
		VALUES(?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET last_activity=excluded.last_activity,
		  title=CASE WHEN excluded.title!='' THEN excluded.title ELSE sessions.title END,
		  ctl_sock=CASE WHEN excluded.ctl_sock!='' THEN excluded.ctl_sock ELSE sessions.ctl_sock END`,
		s.ID, s.Adapter, s.Agent, s.File, s.RepoPath, s.Surface, s.Title, s.CtlSock,
		s.StartedAt.UTC().Format(time.RFC3339), s.LastActivity.UTC().Format(time.RFC3339))
	return err
}

func (d *DB) AddFootprints(sessionID string, paths []string) error {
	for _, p := range paths {
		if _, err := d.Exec(`INSERT INTO footprints(session_id, path, at) VALUES(?,?,?)
			ON CONFLICT(session_id, path) DO UPDATE SET at=excluded.at`, sessionID, p, now()); err != nil {
			return err
		}
	}
	return nil
}

func (d *DB) Footprints(sessionID string) []string {
	rows, err := d.Query(`SELECT path FROM footprints WHERE session_id=? ORDER BY at DESC`, sessionID)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var p string
		rows.Scan(&p)
		out = append(out, p)
	}
	return out
}

// LiveSessions returns sessions active within the window, newest first.
func (d *DB) LiveSessions(repo string, window time.Duration) []Session {
	cutoff := time.Now().Add(-window).UTC().Format(time.RFC3339)
	q := `SELECT id, adapter, agent, file, repo_path, surface, title, ctl_sock, started_at, last_activity
	      FROM sessions WHERE last_activity > ?`
	args := []any{cutoff}
	if repo != "" {
		q += ` AND repo_path = ?`
		args = append(args, repo)
	}
	q += ` ORDER BY last_activity DESC`
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Session
	for rows.Next() {
		var s Session
		var st, la string
		if rows.Scan(&s.ID, &s.Adapter, &s.Agent, &s.File, &s.RepoPath, &s.Surface, &s.Title, &s.CtlSock, &st, &la) == nil {
			s.StartedAt, _ = time.Parse(time.RFC3339, st)
			s.LastActivity, _ = time.Parse(time.RFC3339, la)
			out = append(out, s)
		}
	}
	return out
}

// ---- commits ----

type Commit struct {
	RepoPath, SHA, Author, Subject, SessionID string
	At                                        time.Time
	Files                                     []string
	Mapped                                    bool
}

func (d *DB) CommitSeen(repo, sha string) bool {
	var n int
	d.QueryRow(`SELECT COUNT(*) FROM commits_seen WHERE repo_path=? AND sha=?`, repo, sha).Scan(&n)
	return n > 0
}

func (d *DB) AddCommit(c Commit) error {
	_, err := d.Exec(`INSERT OR IGNORE INTO commits_seen(repo_path, sha, author, at, subject, files, session_id)
		VALUES(?,?,?,?,?,?,?)`,
		c.RepoPath, c.SHA, c.Author, c.At.UTC().Format(time.RFC3339), c.Subject, strings.Join(c.Files, "\n"), c.SessionID)
	return err
}

func (d *DB) AttributeCommit(repo, sha, sessionID string) error {
	_, err := d.Exec(`UPDATE commits_seen SET session_id=? WHERE repo_path=? AND sha=?`, sessionID, repo, sha)
	return err
}

func (d *DB) MarkCommitMapped(repo, sha string) error {
	_, err := d.Exec(`UPDATE commits_seen SET mapped=1 WHERE repo_path=? AND sha=?`, repo, sha)
	return err
}

func (d *DB) RecentCommits(repo string, since time.Time, onlyUnmapped bool) []Commit {
	q := `SELECT repo_path, sha, author, at, subject, files, session_id, mapped FROM commits_seen
	      WHERE repo_path=? AND at > ?`
	if onlyUnmapped {
		q += ` AND mapped=0`
	}
	q += ` ORDER BY at DESC LIMIT 500`
	rows, err := d.Query(q, repo, since.UTC().Format(time.RFC3339))
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Commit
	for rows.Next() {
		var c Commit
		var at, files string
		var mapped int
		if rows.Scan(&c.RepoPath, &c.SHA, &c.Author, &at, &c.Subject, &files, &c.SessionID, &mapped) == nil {
			c.At, _ = time.Parse(time.RFC3339, at)
			if files != "" {
				c.Files = strings.Split(files, "\n")
			}
			c.Mapped = mapped == 1
			out = append(out, c)
		}
	}
	return out
}

// ---- alerts ----

type Alert struct {
	Key, RepoPath, Kind, Body, EntryIDs    string
	Blocking                               bool
	CreatedAt                              time.Time
	NudgedAt, PushedAt, AckedAt, DroppedAt string
}

// UpsertAlert inserts if new; returns true if newly created.
func (d *DB) UpsertAlert(a Alert) (bool, error) {
	res, err := d.Exec(`INSERT OR IGNORE INTO alerts(key, repo_path, kind, body, entry_ids, blocking, created_at)
		VALUES(?,?,?,?,?,?,?)`,
		a.Key, a.RepoPath, a.Kind, a.Body, a.EntryIDs, boolInt(a.Blocking), now())
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (d *DB) OpenAlerts(repo string, blockingOnly bool) []Alert {
	q := `SELECT key, repo_path, kind, body, entry_ids, blocking, created_at,
	      COALESCE(nudged_at,''), COALESCE(pushed_at,''), COALESCE(acked_at,''), COALESCE(dropped_at,'')
	      FROM alerts WHERE acked_at IS NULL AND dropped_at IS NULL`
	args := []any{}
	if repo != "" {
		q += ` AND repo_path=?`
		args = append(args, repo)
	}
	if blockingOnly {
		q += ` AND blocking=1`
	}
	q += ` ORDER BY created_at DESC LIMIT 200`
	rows, err := d.Query(q, args...)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []Alert
	for rows.Next() {
		var a Alert
		var blocking int
		var created string
		if rows.Scan(&a.Key, &a.RepoPath, &a.Kind, &a.Body, &a.EntryIDs, &blocking, &created,
			&a.NudgedAt, &a.PushedAt, &a.AckedAt, &a.DroppedAt) == nil {
			a.Blocking = blocking == 1
			a.CreatedAt, _ = time.Parse(time.RFC3339, created)
			out = append(out, a)
		}
	}
	return out
}

func (d *DB) MarkAlert(key, column string) error {
	switch column {
	case "nudged_at", "pushed_at", "acked_at", "dropped_at":
	default:
		return fmt.Errorf("bad alert column %q", column)
	}
	_, err := d.Exec(`UPDATE alerts SET `+column+`=? WHERE key=?`, now(), key)
	return err
}

// ---- budget (I9) ----

func day() string { return time.Now().UTC().Format("2006-01-02") }

func (d *DB) AddTokens(kind string, n int) error {
	_, err := d.Exec(`INSERT INTO budget(day, kind, tokens) VALUES(?,?,?)
		ON CONFLICT(day, kind) DO UPDATE SET tokens = tokens + excluded.tokens`, day(), kind, n)
	return err
}

func (d *DB) TokensToday(kind string) int {
	var n int
	d.QueryRow(`SELECT tokens FROM budget WHERE day=? AND kind=?`, day(), kind).Scan(&n)
	return n
}

// ---- parked slices (I2) ----

func (d *DB) Park(adapter, file string, offset int64, reason, rawPath string) error {
	_, err := d.Exec(`INSERT INTO parked(adapter, file, offset, reason, raw_path, at) VALUES(?,?,?,?,?,?)`,
		adapter, file, offset, reason, rawPath, now())
	return err
}

func (d *DB) ParkedCount() int {
	var n int
	d.QueryRow(`SELECT COUNT(*) FROM parked`).Scan(&n)
	return n
}

func (d *DB) ParkedRecent(limit int) []string {
	rows, err := d.Query(`SELECT adapter || ': ' || reason || ' (' || file || ')' FROM parked ORDER BY id DESC LIMIT ?`, limit)
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []string
	for rows.Next() {
		var s string
		rows.Scan(&s)
		out = append(out, s)
	}
	return out
}

// ---- kv ----

func (d *DB) Get(k string) string {
	var v string
	d.QueryRow(`SELECT v FROM kv WHERE k=?`, k).Scan(&v)
	return v
}

func (d *DB) Set(k, v string) error {
	_, err := d.Exec(`INSERT INTO kv(k, v) VALUES(?,?) ON CONFLICT(k) DO UPDATE SET v=excluded.v`, k, v)
	return err
}

func (d *DB) Incr(k string) {
	d.Exec(`INSERT INTO kv(k, v) VALUES(?, '1') ON CONFLICT(k) DO UPDATE SET v = CAST(CAST(v AS INTEGER)+1 AS TEXT)`, k)
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
