package differ

import (
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/poller"
	"clew/internal/state"
)

func testJournal(t *testing.T) *journal.Journal {
	t.Helper()
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	return j
}

func testEntry(t *testing.T, j *journal.Journal, typ model.EntryType, title string, at time.Time, mutate ...func(*model.Entry)) *model.Entry {
	t.Helper()
	e := &model.Entry{
		ID: ids.NewEntry(at), Type: typ, Title: title,
		Quote: "verbatim: " + title, UtteranceBy: model.ByUser,
		Source:     model.Source{Kind: model.SrcSession, Ref: "test:session#L1", At: at},
		Confidence: 0.9,
	}
	for _, fn := range mutate {
		fn(e)
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	return e
}

func testEvent(t *testing.T, j *journal.Journal, kind model.EventKind, entry string, at time.Time, payload map[string]any) {
	t.Helper()
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(at), Kind: kind, Entry: entry, At: at,
		Payload: payload, By: model.By{Who: "human", Surface: "test"},
	}); err != nil {
		t.Fatal(err)
	}
}

func TestJournalAlertsWithdrawWhenStatusesResolve(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	j := testJournal(t)
	now := time.Now().UTC().Truncate(time.Second)
	repo := filepath.Join(t.TempDir(), "repo")

	decisionA := testEntry(t, j, model.Decision, "push updates", now.Add(-3*24*time.Hour), func(e *model.Entry) {
		e.Tags = []string{"sync/**"}
	})
	decisionB := testEntry(t, j, model.Decision, "poll updates", now.Add(-2*24*time.Hour), func(e *model.Entry) {
		e.Tags = []string{"sync/**"}
	})
	question := testEntry(t, j, model.Question, "which channel?", now.Add(-8*24*time.Hour), func(e *model.Entry) {
		e.Asks = "human"
	})
	finding := testEntry(t, j, model.Finding, "cache is stale", now.Add(-4*24*time.Hour), func(e *model.Entry) {
		e.Affects = []string{"cache/**"}
	})
	testEvent(t, j, model.EvEvidence, finding.ID, now.Add(-24*time.Hour), map[string]any{
		"kind": "churn", "ref": "churn:abc123",
	})
	target := testEntry(t, j, model.Intent, "ship the dormant path", now.Add(-30*24*time.Hour))
	for i := 0; i < journal.AbsenceSiblings; i++ {
		sibling := testEntry(t, j, model.Intent, "progress sibling "+string(rune('A'+i)), now.Add(-29*24*time.Hour+time.Duration(i)*time.Minute))
		testEvent(t, j, model.EvEvidence, sibling.ID, now.Add(-10*24*time.Hour+time.Duration(i)*time.Minute), map[string]any{
			"kind": "commit", "ref": "sibling-" + string(rune('A'+i)),
		})
	}

	first, err := Run(db, &Input{Repo: repo, Journal: j, Snapshot: &poller.Snapshot{RepoPath: repo}}, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(first.NewAlerts) != 4 {
		t.Fatalf("new alerts = %#v, want contradiction, absence, aging, and suspect", first.NewAlerts)
	}
	open := db.OpenAlerts(repo, false)
	if len(open) != 4 {
		t.Fatalf("open alerts = %#v, want 4", open)
	}
	wantCondition := map[string]string{
		"contradiction": "contradiction:status-resolved:",
		"absence":       "absence:status-resolved:",
		"aging":         "aging:status-resolved:",
		"suspect":       "suspect:status-resolved:",
	}
	for _, alert := range open {
		prefix, ok := wantCondition[alert.Kind]
		if !ok {
			t.Fatalf("unexpected alert: %#v", alert)
		}
		if !strings.HasPrefix(alert.WithdrawWhen, prefix) {
			t.Errorf("%s withdrawal condition = %q, want prefix %q", alert.Kind, alert.WithdrawWhen, prefix)
		}
		delete(wantCondition, alert.Kind)
	}
	if len(wantCondition) != 0 {
		t.Fatalf("missing alert kinds: %v", wantCondition)
	}

	resolvedAt := now.Add(time.Minute)
	testEvent(t, j, model.EvSupersede, decisionA.ID, resolvedAt, map[string]any{"by": decisionB.ID})
	testEvent(t, j, model.EvAnswer, question.ID, resolvedAt, map[string]any{"by": decisionB.ID})
	testEvent(t, j, model.EvSupersede, finding.ID, resolvedAt, map[string]any{"by": "new-measurement"})
	testEvent(t, j, model.EvConfirm, target.ID, resolvedAt, map[string]any{"done": true})

	second, err := Run(db, &Input{Repo: repo, Journal: j, Snapshot: &poller.Snapshot{RepoPath: repo}}, resolvedAt)
	if err != nil {
		t.Fatal(err)
	}
	if len(second.NewAlerts) != 0 {
		t.Fatalf("resolution created alerts: %#v", second.NewAlerts)
	}
	if got := db.OpenAlerts(repo, false); len(got) != 0 {
		t.Fatalf("resolved journal alerts survived one poll: %#v", got)
	}
	var dropped int
	if err := db.QueryRow(`SELECT COUNT(*) FROM alerts WHERE repo_path=? AND dropped_at IS NOT NULL`, repo).Scan(&dropped); err != nil {
		t.Fatal(err)
	}
	if dropped != 4 {
		t.Fatalf("withdrawn alerts = %d, want 4", dropped)
	}
}

func TestCommitSubjectSettlesMatchingDecisionWithoutModelCall(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	j := testJournal(t)
	repo := filepath.Join(t.TempDir(), "repo")
	now := time.Now().UTC().Truncate(time.Second)
	matched := testEntry(t, j, model.Decision, "Evidence settles merge lines for verified work", now.Add(-time.Hour))
	unrelated := testEntry(t, j, model.Decision, "The tablet shows decision cards", now.Add(-time.Hour))
	crossClause := testEntry(t, j, model.Decision, "Verified phone status", now.Add(-time.Hour))
	commit := state.Commit{
		RepoPath: repo, SHA: strings.Repeat("a", 40), Author: "test", At: now,
		Subject: "Settle verified evidence in the merge; update phone card", Files: []string{"internal/state/state.go"},
	}
	if err := db.AddCommit(commit); err != nil {
		t.Fatal(err)
	}
	if err := db.MarkCommitMapped(repo, commit.SHA); err != nil {
		t.Fatal(err)
	}
	res, err := Run(db, &Input{Repo: repo, Journal: j, Snapshot: &poller.Snapshot{RepoPath: repo}}, now)
	if err != nil {
		t.Fatal(err)
	}
	if res.EventsAdded != 1 || !j.HasEvent(model.EvEvidence, matched.ID, "ref", commit.SHA) {
		t.Fatalf("matching decision did not receive commit evidence: result=%+v events=%#v", res, j.Events)
	}
	if j.HasEvent(model.EvEvidence, unrelated.ID, "ref", commit.SHA) {
		t.Fatal("unrelated decision received guessed evidence")
	}
	if j.HasEvent(model.EvEvidence, crossClause.ID, "ref", commit.SHA) {
		t.Fatal("words split across commit clauses guessed evidence")
	}
}

func TestJournalCommitSubjectNeverProvesCodeReality(t *testing.T) {
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	j := testJournal(t)
	repo := filepath.Join(t.TempDir(), "repo")
	now := time.Now().UTC().Truncate(time.Second)
	entry := testEntry(t, j, model.Decision, "Evidence settles merge lines", now.Add(-time.Hour))
	commit := state.Commit{
		RepoPath: repo, SHA: strings.Repeat("b", 40), Author: "test", At: now,
		Subject: "journal: settle merge evidence ruling", Files: []string{"entries/decision.yaml"},
	}
	if err := db.AddCommit(commit); err != nil {
		t.Fatal(err)
	}
	if _, err := Run(db, &Input{Repo: repo, Journal: j, Snapshot: &poller.Snapshot{RepoPath: repo}}, now); err != nil {
		t.Fatal(err)
	}
	if j.HasEvent(model.EvEvidence, entry.ID, "ref", commit.SHA) {
		t.Fatal("journal coordination commit was attached as code evidence")
	}
}

func TestWorkspaceAlertsWithdrawWithinOnePoll(t *testing.T) {
	tests := []struct {
		name        string
		firstDirty  []string
		secondDirty []string
		expire      bool
		firstKind   string
		wantAfter   string
	}{
		{
			name: "stomp becomes clean overlap", firstDirty: []string{"shared.go"},
			firstKind: "stomp", wantAfter: "overlap",
		},
		{
			name: "stomp sessions diverge", firstDirty: []string{"shared.go"}, secondDirty: []string{"shared.go"},
			expire: true, firstKind: "stomp",
		},
		{
			name: "overlap sessions diverge", expire: true, firstKind: "overlap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			j := testJournal(t)
			repo := filepath.Join(t.TempDir(), "repo")
			observed := time.Now().UTC().Truncate(time.Second)
			for _, session := range []state.Session{
				{
					ID: "session-a", Agent: "codex", RepoPath: repo, Surface: "desk",
					StartedAt: observed.Add(-10 * time.Minute), LastActivity: observed,
				},
				{
					ID: "session-b", Agent: "claude", RepoPath: repo, Surface: "phone",
					StartedAt: observed.Add(-9 * time.Minute), LastActivity: observed.Add(-time.Second),
				},
			} {
				if err := db.UpsertSession(session); err != nil {
					t.Fatal(err)
				}
				if err := db.AddFootprints(session.ID, []string{filepath.Join(repo, "shared.go")}); err != nil {
					t.Fatal(err)
				}
			}

			first, err := Run(db, &Input{
				Repo: repo, Journal: j,
				Snapshot: &poller.Snapshot{RepoPath: repo, Dirty: tt.firstDirty},
			}, observed)
			if err != nil {
				t.Fatal(err)
			}
			if len(first.NewAlerts) != 1 || first.NewAlerts[0].Kind != tt.firstKind {
				t.Fatalf("first alerts = %#v, want one %s", first.NewAlerts, tt.firstKind)
			}
			firstAlert := first.NewAlerts[0]
			if firstAlert.WithdrawWhen == "" {
				t.Fatal("workspace alert has no withdrawal condition")
			}
			if tt.expire {
				if _, err := db.Exec(`UPDATE sessions SET last_activity=? WHERE id='session-b'`,
					time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)); err != nil {
					t.Fatal(err)
				}
			}

			second, err := Run(db, &Input{
				Repo: repo, Journal: j,
				Snapshot: &poller.Snapshot{RepoPath: repo, Dirty: tt.secondDirty},
			}, observed.Add(time.Minute))
			if err != nil {
				t.Fatal(err)
			}
			open := db.OpenAlerts(repo, false)
			if tt.wantAfter == "" {
				if len(open) != 0 || len(second.NewAlerts) != 0 {
					t.Fatalf("resolved alert survived one poll: open=%#v new=%#v", open, second.NewAlerts)
				}
			} else if len(open) != 1 || open[0].Kind != tt.wantAfter || len(second.NewAlerts) != 1 || second.NewAlerts[0].Kind != tt.wantAfter {
				t.Fatalf("after resolution: open=%#v new=%#v, want one new %s", open, second.NewAlerts, tt.wantAfter)
			}
			var dropped string
			if err := db.QueryRow(`SELECT COALESCE(dropped_at,'') FROM alerts WHERE key=?`, firstAlert.Key).Scan(&dropped); err != nil {
				t.Fatal(err)
			}
			if dropped == "" {
				t.Fatalf("resolved %s alert was not withdrawn", tt.firstKind)
			}
		})
	}
}
