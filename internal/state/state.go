// Package state is the machine-local working state (~/.clew/state.db,
// JOURNAL_SPEC §2). SQLite, never the source of truth — the journal lives in
// git. Holds: registered repos, tail watermarks, live sessions + footprints,
// seen commits + attribution, alerts, budget meters, parked slices, kv.
package state

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
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
  acked_at TEXT, dropped_at TEXT, withdraw_when TEXT);
CREATE TABLE IF NOT EXISTS budget(
  day TEXT, kind TEXT, tokens INTEGER, PRIMARY KEY(day, kind));
CREATE TABLE IF NOT EXISTS llm_budget_reservations(
  id TEXT PRIMARY KEY, day TEXT NOT NULL, kind TEXT NOT NULL,
  tokens INTEGER NOT NULL CHECK(tokens >= 0), created_at TEXT NOT NULL);
CREATE INDEX IF NOT EXISTS llm_budget_reservations_day_kind
  ON llm_budget_reservations(day, kind);
CREATE TABLE IF NOT EXISTS parked(
  id INTEGER PRIMARY KEY AUTOINCREMENT, adapter TEXT, file TEXT,
  offset INTEGER, reason TEXT, raw_path TEXT, at TEXT);
CREATE TABLE IF NOT EXISTS kv(k TEXT PRIMARY KEY, v TEXT);
`)
	if err != nil {
		return err
	}
	return migrateAlertWithdrawal(db)
}

// migrateAlertWithdrawal preserves databases created before alert withdrawal
// conditions became durable. CREATE TABLE IF NOT EXISTS does not add columns
// to an existing table, so keep the additive migration explicit.
func migrateAlertWithdrawal(db *sql.DB) error {
	rows, err := db.Query(`PRAGMA table_info(alerts)`)
	if err != nil {
		return err
	}
	defer rows.Close()
	found := false
	for rows.Next() {
		var cid, notNull, primaryKey int
		var name, typ string
		var defaultValue sql.NullString
		if err := rows.Scan(&cid, &name, &typ, &notNull, &defaultValue, &primaryKey); err != nil {
			return err
		}
		if name == "withdraw_when" {
			found = true
		}
	}
	if err := rows.Err(); err != nil {
		return err
	}
	if err := rows.Close(); err != nil {
		return err
	}
	if found {
		return nil
	}
	_, err = db.Exec(`ALTER TABLE alerts ADD COLUMN withdraw_when TEXT`)
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

func (d *DB) RepoRegistered(path string) bool {
	var n int
	d.QueryRow(`SELECT COUNT(*) FROM repos WHERE path=?`, path).Scan(&n)
	return n > 0
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
	off, _ := d.WatermarkOK(file)
	return off
}

func (d *DB) WatermarkOK(file string) (int64, bool) {
	var off int64
	if err := d.QueryRow(`SELECT offset FROM watermarks WHERE file=?`, file).Scan(&off); err != nil {
		return 0, false
	}
	return off, true
}

func (d *DB) SetWatermark(file, adapter, repo string, offset int64) error {
	_, err := d.Exec(`INSERT INTO watermarks(file, adapter, repo_path, offset, updated_at)
		VALUES(?,?,?,?,?) ON CONFLICT(file) DO UPDATE SET offset=excluded.offset, updated_at=excluded.updated_at`,
		file, adapter, repo, offset, now())
	return err
}

// InitWatermark records a starting offset without advancing an existing
// watermark. init uses this to make watch forward-only; explicit backfill
// keeps its independent extract: watermark.
func (d *DB) InitWatermark(file, adapter, repo string, offset int64) (bool, error) {
	n, err := d.InitWatermarks(WatermarkInit{File: file, Adapter: adapter, Repo: repo, Offset: offset})
	return n == 1, err
}

type WatermarkInit struct {
	File, Adapter, Repo string
	Offset              int64
}

// InitWatermarks enrolls all cursors for one session atomically. INSERT OR
// IGNORE makes concurrent watcher/backfill enrollment choose one coherent
// owner without moving established cursors.
func (d *DB) InitWatermarks(inits ...WatermarkInit) (int, error) {
	tx, err := d.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()
	added := 0
	stamp := now()
	for _, init := range inits {
		res, err := tx.Exec(`INSERT OR IGNORE INTO watermarks(file, adapter, repo_path, offset, updated_at)
			VALUES(?,?,?,?,?)`, init.File, init.Adapter, init.Repo, init.Offset, stamp)
		if err != nil {
			return added, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return added, err
		}
		added += int(n)
	}
	if err := tx.Commit(); err != nil {
		return added, err
	}
	return added, nil
}

// ---- sessions ----

type Session struct {
	ID, Adapter, Agent, File, RepoPath, Surface, Title, CtlSock string
	StartedAt, LastActivity                                     time.Time
}

func (d *DB) UpsertSession(s Session) error {
	_, err := d.Exec(`INSERT INTO sessions(id, adapter, agent, file, repo_path, surface, title, ctl_sock, started_at, last_activity)
		VALUES(?,?,?,?,?,?,?,?,?,?)
		ON CONFLICT(id) DO UPDATE SET
		  started_at=CASE WHEN excluded.started_at < sessions.started_at THEN excluded.started_at ELSE sessions.started_at END,
		  last_activity=CASE WHEN excluded.last_activity > sessions.last_activity THEN excluded.last_activity ELSE sessions.last_activity END,
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
	WithdrawWhen                           string
	Blocking                               bool
	CreatedAt                              time.Time
	NudgedAt, PushedAt, AckedAt, DroppedAt string
}

// UpsertAlert inserts if new; returns true if newly created.
func (d *DB) UpsertAlert(a Alert) (bool, error) {
	res, err := d.Exec(`INSERT OR IGNORE INTO alerts(key, repo_path, kind, body, entry_ids, blocking, created_at, withdraw_when)
		VALUES(?,?,?,?,?,?,?,?)`,
		a.Key, a.RepoPath, a.Kind, a.Body, a.EntryIDs, boolInt(a.Blocking), now(), a.WithdrawWhen)
	if err != nil {
		return false, err
	}
	n, _ := res.RowsAffected()
	return n > 0, nil
}

func (d *DB) OpenAlerts(repo string, blockingOnly bool) []Alert {
	q := `SELECT key, repo_path, kind, body, entry_ids, blocking, created_at,
	      COALESCE(nudged_at,''), COALESCE(pushed_at,''), COALESCE(acked_at,''), COALESCE(dropped_at,''),
	      COALESCE(withdraw_when,'')
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
			&a.NudgedAt, &a.PushedAt, &a.AckedAt, &a.DroppedAt, &a.WithdrawWhen) == nil {
			a.Blocking = blocking == 1
			a.CreatedAt, _ = time.Parse(time.RFC3339, created)
			out = append(out, a)
		}
	}
	return out
}

// ReconcileAlerts makes active the complete alert set for the supplied
// repo/kinds. Active alerts are inserted or refreshed; any previously open
// alert in the same scope whose key is absent is withdrawn via dropped_at.
// It returns only alerts inserted by this reconciliation, for delivery.
//
// WithdrawWhen is required here (but not by UpsertAlert, which also stores
// operational alerts owned by other subsystems). It is an opaque,
// machine-readable condition naming why the differ will stop emitting the
// alert on a later poll.
func (d *DB) ReconcileAlerts(repo string, kinds []string, active []Alert) ([]Alert, error) {
	managed := make(map[string]bool, len(kinds))
	for _, kind := range kinds {
		if kind == "" {
			return nil, fmt.Errorf("cannot reconcile empty alert kind")
		}
		managed[kind] = true
	}
	if len(managed) == 0 {
		return nil, nil
	}
	activeKeys := make(map[string]bool, len(active))
	for _, a := range active {
		switch {
		case a.Key == "":
			return nil, fmt.Errorf("cannot reconcile alert with empty key")
		case a.RepoPath != repo:
			return nil, fmt.Errorf("alert %s belongs to repo %q, reconcile scope is %q", a.Key, a.RepoPath, repo)
		case !managed[a.Kind]:
			return nil, fmt.Errorf("alert %s kind %q is outside reconcile scope", a.Key, a.Kind)
		case a.WithdrawWhen == "":
			return nil, fmt.Errorf("alert %s has no withdrawal condition", a.Key)
		case activeKeys[a.Key]:
			return nil, fmt.Errorf("duplicate active alert key %s", a.Key)
		}
		activeKeys[a.Key] = true
	}

	tx, err := d.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	stamp := now()
	created := make([]Alert, 0, len(active))
	for _, a := range active {
		res, err := tx.Exec(`INSERT OR IGNORE INTO alerts(
			key, repo_path, kind, body, entry_ids, blocking, created_at, withdraw_when)
			VALUES(?,?,?,?,?,?,?,?)`, a.Key, a.RepoPath, a.Kind, a.Body, a.EntryIDs,
			boolInt(a.Blocking), stamp, a.WithdrawWhen)
		if err != nil {
			return nil, err
		}
		n, err := res.RowsAffected()
		if err != nil {
			return nil, err
		}
		if n > 0 {
			a.CreatedAt, _ = time.Parse(time.RFC3339, stamp)
			created = append(created, a)
			continue
		}
		// Stable keys let mutable display prose (for example question age or
		// session labels) refresh without creating a new alert.
		if _, err := tx.Exec(`UPDATE alerts SET body=?, entry_ids=?, blocking=?, withdraw_when=?
			WHERE key=? AND repo_path=? AND kind=? AND acked_at IS NULL AND dropped_at IS NULL`,
			a.Body, a.EntryIDs, boolInt(a.Blocking), a.WithdrawWhen,
			a.Key, repo, a.Kind); err != nil {
			return nil, err
		}
	}

	args := make([]any, 0, 2+len(managed)+len(activeKeys))
	args = append(args, stamp, repo)
	kindMarks := make([]string, 0, len(managed))
	for kind := range managed {
		kindMarks = append(kindMarks, "?")
		args = append(args, kind)
	}
	q := `UPDATE alerts SET dropped_at=? WHERE repo_path=? AND kind IN (` + strings.Join(kindMarks, ",") + `)
		AND acked_at IS NULL AND dropped_at IS NULL`
	if len(activeKeys) > 0 {
		keyMarks := make([]string, 0, len(activeKeys))
		for key := range activeKeys {
			keyMarks = append(keyMarks, "?")
			args = append(args, key)
		}
		q += ` AND key NOT IN (` + strings.Join(keyMarks, ",") + `)`
	}
	if _, err := tx.Exec(q, args...); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return created, nil
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

// LLMBudgetLimits are checked when reserving a provider call. DailyCapTokens
// always applies to aggregate LLM spend. LiveSessionPct is optional (zero
// disables it) and applies only to the extraction cost center, whose spend is
// bounded against today's observed session tokens.
type LLMBudgetLimits struct {
	DailyCapTokens int
	LiveSessionPct float64
}

// LLMBudgetReservation is the durable claim a caller must settle after its
// provider call, including on provider failure (with actualTokens=0).
type LLMBudgetReservation struct {
	ID, Day, Kind string
	Tokens        int
}

// LLMBudgetLimitError is returned without creating a reservation. Limit is
// either "daily-cap" or "live-session-ratio".
type LLMBudgetLimitError struct {
	Limit                      string
	Kind                       string
	Requested, Spent, Reserved int
	LimitTokens, Observed      int
	LiveSessionPct             float64
}

func (e *LLMBudgetLimitError) Error() string {
	if e.Limit == "live-session-ratio" {
		return fmt.Sprintf("LLM budget reservation denied: extraction spend %d + reserved %d + requested %d exceeds %.2f%% of %d observed tokens (%d)",
			e.Spent, e.Reserved, e.Requested, e.LiveSessionPct, e.Observed, e.LimitTokens)
	}
	return fmt.Sprintf("LLM budget reservation denied: aggregate spend %d + reserved %d + requested %d exceeds daily cap %d",
		e.Spent, e.Reserved, e.Requested, e.LimitTokens)
}

// LLMBudgetOverrunError is returned after settlement has committed the actual
// spend. Callers must surface it loudly; accounting is intentionally honest.
type LLMBudgetOverrunError struct {
	ReservationID, Kind string
	Reserved, Actual    int
}

func (e *LLMBudgetOverrunError) Error() string {
	return fmt.Sprintf("LLM budget overrun: %s reservation %s reserved %d tokens but used %d; actual spend was recorded",
		e.Kind, e.ReservationID, e.Reserved, e.Actual)
}

type LLMBudgetReservationNotFoundError struct{ ReservationID string }

func (e *LLMBudgetReservationNotFoundError) Error() string {
	return fmt.Sprintf("LLM budget reservation %s not found (already settled or unknown)", e.ReservationID)
}

// ReserveLLMBudget atomically admits an estimated provider call. BEGIN
// IMMEDIATE serializes the read/check/insert across goroutines and processes,
// so concurrent callers cannot each observe the same remaining capacity.
// kind is a base cost-center name such as extraction, differ, or archaeology.
func (d *DB) ReserveLLMBudget(kind string, estimate int, limits LLMBudgetLimits) (*LLMBudgetReservation, error) {
	if kind == "" || kind == "spent" || kind == "observed" || strings.HasSuffix(kind, "-spent") {
		return nil, fmt.Errorf("invalid LLM budget kind %q (use a base cost-center name)", kind)
	}
	if estimate < 0 {
		return nil, fmt.Errorf("invalid LLM budget estimate %d", estimate)
	}
	if limits.DailyCapTokens <= 0 {
		return nil, fmt.Errorf("invalid LLM daily cap %d", limits.DailyCapTokens)
	}
	if limits.LiveSessionPct < 0 {
		return nil, fmt.Errorf("invalid live-session percentage %.2f", limits.LiveSessionPct)
	}
	id, err := newLLMBudgetReservationID()
	if err != nil {
		return nil, err
	}
	r := &LLMBudgetReservation{ID: id, Kind: kind, Tokens: estimate}
	err = d.withImmediate(func(conn *sql.Conn) error {
		r.Day = day()
		spent, err := budgetTokensConn(conn, r.Day, "spent")
		if err != nil {
			return err
		}
		reserved, err := reservedTokensConn(conn, r.Day, "")
		if err != nil {
			return err
		}
		if spent+reserved+estimate > limits.DailyCapTokens {
			return &LLMBudgetLimitError{
				Limit: "daily-cap", Kind: kind, Requested: estimate,
				Spent: spent, Reserved: reserved, LimitTokens: limits.DailyCapTokens,
			}
		}
		if kind == "extraction" && limits.LiveSessionPct > 0 {
			extractionSpent, err := budgetTokensConn(conn, r.Day, "extraction-spent")
			if err != nil {
				return err
			}
			extractionReserved, err := reservedTokensConn(conn, r.Day, "extraction")
			if err != nil {
				return err
			}
			observed, err := budgetTokensConn(conn, r.Day, "observed")
			if err != nil {
				return err
			}
			ratioCap := int(float64(observed) * limits.LiveSessionPct / 100)
			if extractionSpent+extractionReserved+estimate > ratioCap {
				return &LLMBudgetLimitError{
					Limit: "live-session-ratio", Kind: kind, Requested: estimate,
					Spent: extractionSpent, Reserved: extractionReserved,
					LimitTokens: ratioCap, Observed: observed, LiveSessionPct: limits.LiveSessionPct,
				}
			}
		}
		_, err = conn.ExecContext(context.Background(), `INSERT INTO llm_budget_reservations(id, day, kind, tokens, created_at)
			VALUES(?,?,?,?,?)`, r.ID, r.Day, r.Kind, r.Tokens, now())
		return err
	})
	if err != nil {
		return nil, err
	}
	return r, nil
}

// SettleLLMBudget atomically releases a reservation and records actual spend
// in both the aggregate "spent" meter and the per-kind "<kind>-spent" meter.
// If actual exceeds the reservation, accounting commits before the loud error
// is returned.
func (d *DB) SettleLLMBudget(reservationID string, actualTokens int) error {
	if reservationID == "" {
		return fmt.Errorf("empty LLM budget reservation id")
	}
	if actualTokens < 0 {
		return fmt.Errorf("invalid actual LLM token count %d", actualTokens)
	}
	var overrun *LLMBudgetOverrunError
	err := d.withImmediate(func(conn *sql.Conn) error {
		var reservation LLMBudgetReservation
		err := conn.QueryRowContext(context.Background(),
			`SELECT id, day, kind, tokens FROM llm_budget_reservations WHERE id=?`, reservationID).
			Scan(&reservation.ID, &reservation.Day, &reservation.Kind, &reservation.Tokens)
		if err == sql.ErrNoRows {
			return &LLMBudgetReservationNotFoundError{ReservationID: reservationID}
		}
		if err != nil {
			return err
		}
		if _, err := conn.ExecContext(context.Background(),
			`DELETE FROM llm_budget_reservations WHERE id=?`, reservationID); err != nil {
			return err
		}
		if err := addBudgetTokensConn(conn, reservation.Day, "spent", actualTokens); err != nil {
			return err
		}
		if err := addBudgetTokensConn(conn, reservation.Day, reservation.Kind+"-spent", actualTokens); err != nil {
			return err
		}
		if actualTokens > reservation.Tokens {
			overrun = &LLMBudgetOverrunError{
				ReservationID: reservation.ID, Kind: reservation.Kind,
				Reserved: reservation.Tokens, Actual: actualTokens,
			}
		}
		return nil
	})
	if err != nil {
		return err
	}
	if overrun != nil {
		return overrun
	}
	return nil
}

func newLLMBudgetReservationID() (string, error) {
	var raw [16]byte
	if _, err := rand.Read(raw[:]); err != nil {
		return "", fmt.Errorf("create LLM budget reservation id: %w", err)
	}
	return hex.EncodeToString(raw[:]), nil
}

// withImmediate is a small SQLite transaction primitive for budget admission
// and settlement. A deferred transaction would allow two readers to see the
// same capacity before either becomes the writer.
func (d *DB) withImmediate(fn func(*sql.Conn) error) (err error) {
	ctx := context.Background()
	conn, err := d.Conn(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()
	if _, err = conn.ExecContext(ctx, "BEGIN IMMEDIATE"); err != nil {
		return err
	}
	committed := false
	defer func() {
		if !committed {
			_, _ = conn.ExecContext(ctx, "ROLLBACK")
		}
	}()
	if err = fn(conn); err != nil {
		return err
	}
	if _, err = conn.ExecContext(ctx, "COMMIT"); err != nil {
		return err
	}
	committed = true
	return nil
}

func budgetTokensConn(conn *sql.Conn, budgetDay, kind string) (int, error) {
	var tokens int
	err := conn.QueryRowContext(context.Background(),
		`SELECT COALESCE((SELECT tokens FROM budget WHERE day=? AND kind=?), 0)`, budgetDay, kind).Scan(&tokens)
	return tokens, err
}

func reservedTokensConn(conn *sql.Conn, budgetDay, kind string) (int, error) {
	var tokens int
	query := `SELECT COALESCE(SUM(tokens), 0) FROM llm_budget_reservations WHERE day=?`
	args := []any{budgetDay}
	if kind != "" {
		query += ` AND kind=?`
		args = append(args, kind)
	}
	err := conn.QueryRowContext(context.Background(), query, args...).Scan(&tokens)
	return tokens, err
}

func addBudgetTokensConn(conn *sql.Conn, budgetDay, kind string, n int) error {
	_, err := conn.ExecContext(context.Background(), `INSERT INTO budget(day, kind, tokens) VALUES(?,?,?)
		ON CONFLICT(day, kind) DO UPDATE SET tokens = tokens + excluded.tokens`, budgetDay, kind, n)
	return err
}

func (d *DB) AddTokens(kind string, n int) error {
	_, err := d.Exec(`INSERT INTO budget(day, kind, tokens) VALUES(?,?,?)
		ON CONFLICT(day, kind) DO UPDATE SET tokens = tokens + excluded.tokens`, day(), kind, n)
	return err
}

// RecordSpend preserves the aggregate I9 meter while also making each LLM
// cost center independently measurable during dogfood.
func (d *DB) RecordSpend(kind string, n int) error {
	return d.withImmediate(func(conn *sql.Conn) error {
		budgetDay := day()
		if err := addBudgetTokensConn(conn, budgetDay, "spent", n); err != nil {
			return err
		}
		return addBudgetTokensConn(conn, budgetDay, kind+"-spent", n)
	})
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

type KVPair struct {
	Key, Value string
}

// KVPrefix returns non-empty machine degradations in stable key order. Status
// uses this instead of requiring operators to inspect state.db directly (I2).
func (d *DB) KVPrefix(prefix string) []KVPair {
	rows, err := d.Query(`SELECT k, v FROM kv WHERE k LIKE ? AND v != '' ORDER BY k`, prefix+"%")
	if err != nil {
		return nil
	}
	defer rows.Close()
	var out []KVPair
	for rows.Next() {
		var p KVPair
		if rows.Scan(&p.Key, &p.Value) == nil {
			out = append(out, p)
		}
	}
	return out
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
