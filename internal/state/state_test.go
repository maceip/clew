package state

import (
	"errors"
	"path/filepath"
	"sync"
	"testing"
	"time"
)

func TestSessionActivityRangeNeverRegresses(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	now := time.Now().UTC().Truncate(time.Second)
	first := Session{ID: "s1", RepoPath: "/repo", StartedAt: now.Add(-time.Hour), LastActivity: now}
	if err := db.UpsertSession(first); err != nil {
		t.Fatal(err)
	}
	stale := first
	stale.StartedAt = now.Add(-2 * time.Hour)
	stale.LastActivity = now.Add(-30 * time.Minute)
	if err := db.UpsertSession(stale); err != nil {
		t.Fatal(err)
	}

	var started, last string
	if err := db.QueryRow(`SELECT started_at, last_activity FROM sessions WHERE id='s1'`).Scan(&started, &last); err != nil {
		t.Fatal(err)
	}
	if started != stale.StartedAt.Format(time.RFC3339) {
		t.Fatalf("started_at = %s, want %s", started, stale.StartedAt.Format(time.RFC3339))
	}
	if last != first.LastActivity.Format(time.RFC3339) {
		t.Fatalf("last_activity regressed to %s, want %s", last, first.LastActivity.Format(time.RFC3339))
	}
}

func TestInitWatermarkDoesNotAdvanceExistingOffset(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	added, err := db.InitWatermark("tail:file", "fixture", "/repo", 4)
	if err != nil || !added {
		t.Fatalf("first InitWatermark() = %v, %v; want true, nil", added, err)
	}
	added, err = db.InitWatermark("tail:file", "fixture", "/repo", 8)
	if err != nil || added {
		t.Fatalf("second InitWatermark() = %v, %v; want false, nil", added, err)
	}
	if got := db.Watermark("tail:file"); got != 4 {
		t.Fatalf("watermark advanced to %d, want 4", got)
	}
}

func TestKVPrefixReturnsOnlyNonEmptyMatches(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Set("unknown:codex:item:a", "2"); err != nil {
		t.Fatal(err)
	}
	if err := db.Set("unknown:codex:item:b", ""); err != nil {
		t.Fatal(err)
	}
	if err := db.Set("other", "9"); err != nil {
		t.Fatal(err)
	}
	got := db.KVPrefix("unknown:")
	if len(got) != 1 || got[0].Key != "unknown:codex:item:a" || got[0].Value != "2" {
		t.Fatalf("KVPrefix = %#v", got)
	}
}

func TestReserveLLMBudgetEnforcesAggregateAndExtractionLimits(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.AddTokens("observed", 10_000); err != nil {
		t.Fatal(err)
	}
	if err := db.RecordSpend("extraction", 100); err != nil {
		t.Fatal(err)
	}
	limits := LLMBudgetLimits{DailyCapTokens: 1_000, LiveSessionPct: 2}

	if _, err := db.ReserveLLMBudget("extraction", 80, limits); err != nil {
		t.Fatalf("reserve below ratio: %v", err)
	}
	if _, err := db.ReserveLLMBudget("extraction", 21, limits); err == nil {
		t.Fatal("reservation exceeding extraction ratio was admitted")
	} else {
		var limit *LLMBudgetLimitError
		if !errors.As(err, &limit) || limit.Limit != "live-session-ratio" {
			t.Fatalf("ratio error = %T %v", err, err)
		}
		if limit.Spent != 100 || limit.Reserved != 80 || limit.LimitTokens != 200 {
			t.Fatalf("ratio diagnostics = %+v", limit)
		}
	}
	if _, err := db.ReserveLLMBudget("extraction", 20, limits); err != nil {
		t.Fatalf("reserve exactly to ratio: %v", err)
	}
	// Non-extraction calls share the aggregate cap but not the live-session
	// denominator. Existing spend 100 + reservations 100 + 800 = cap 1000.
	if _, err := db.ReserveLLMBudget("differ", 800, limits); err != nil {
		t.Fatalf("reserve exactly to daily cap: %v", err)
	}
	if _, err := db.ReserveLLMBudget("differ", 1, limits); err == nil {
		t.Fatal("reservation exceeding aggregate cap was admitted")
	} else {
		var limit *LLMBudgetLimitError
		if !errors.As(err, &limit) || limit.Limit != "daily-cap" {
			t.Fatalf("daily-cap error = %T %v", err, err)
		}
		if limit.Spent != 100 || limit.Reserved != 900 || limit.LimitTokens != 1_000 {
			t.Fatalf("daily diagnostics = %+v", limit)
		}
	}
}

func TestSettleLLMBudgetRecordsActualAndReportsOverrun(t *testing.T) {
	db, err := Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	limits := LLMBudgetLimits{DailyCapTokens: 1_000}

	under, err := db.ReserveLLMBudget("differ", 40, limits)
	if err != nil {
		t.Fatal(err)
	}
	if err := db.SettleLLMBudget(under.ID, 30); err != nil {
		t.Fatalf("settle under reservation: %v", err)
	}
	if got := db.TokensToday("spent"); got != 30 {
		t.Fatalf("aggregate actual = %d, want 30", got)
	}
	if got := db.TokensToday("differ-spent"); got != 30 {
		t.Fatalf("per-kind actual = %d, want 30", got)
	}

	over, err := db.ReserveLLMBudget("extraction", 20, limits)
	if err != nil {
		t.Fatal(err)
	}
	err = db.SettleLLMBudget(over.ID, 25)
	var overrun *LLMBudgetOverrunError
	if !errors.As(err, &overrun) {
		t.Fatalf("overrun settlement error = %T %v", err, err)
	}
	if overrun.Reserved != 20 || overrun.Actual != 25 || overrun.Kind != "extraction" {
		t.Fatalf("overrun diagnostics = %+v", overrun)
	}
	if got := db.TokensToday("spent"); got != 55 {
		t.Fatalf("aggregate did not record overrun honestly: %d", got)
	}
	if got := db.TokensToday("extraction-spent"); got != 25 {
		t.Fatalf("per-kind did not record overrun honestly: %d", got)
	}
	var active int
	if err := db.QueryRow(`SELECT COUNT(*) FROM llm_budget_reservations`).Scan(&active); err != nil {
		t.Fatal(err)
	}
	if active != 0 {
		t.Fatalf("settlement left %d active reservations", active)
	}
	err = db.SettleLLMBudget(over.ID, 25)
	var missing *LLMBudgetReservationNotFoundError
	if !errors.As(err, &missing) {
		t.Fatalf("second settlement = %T %v, want not-found", err, err)
	}
	if got := db.TokensToday("spent"); got != 55 {
		t.Fatalf("double settlement changed aggregate to %d", got)
	}
}

func TestConcurrentLLMBudgetReservationsCannotOverAdmit(t *testing.T) {
	tests := []struct {
		name         string
		setup        func(*testing.T, *DB)
		kind         string
		limits       LLMBudgetLimits
		estimate     int
		wantAdmitted int
		wantLimit    string
		wantReserved int
	}{
		{
			name: "aggregate daily cap",
			setup: func(t *testing.T, db *DB) {
				t.Helper()
				if err := db.RecordSpend("differ", 10); err != nil {
					t.Fatal(err)
				}
			},
			kind: "differ", limits: LLMBudgetLimits{DailyCapTokens: 100},
			estimate: 30, wantAdmitted: 3, wantLimit: "daily-cap", wantReserved: 90,
		},
		{
			name: "extraction live-session ratio",
			setup: func(t *testing.T, db *DB) {
				t.Helper()
				if err := db.AddTokens("observed", 5_000); err != nil {
					t.Fatal(err)
				}
				if err := db.RecordSpend("extraction", 10); err != nil {
					t.Fatal(err)
				}
			},
			kind: "extraction", limits: LLMBudgetLimits{DailyCapTokens: 10_000, LiveSessionPct: 2},
			estimate: 30, wantAdmitted: 3, wantLimit: "live-session-ratio", wantReserved: 90,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db, err := Open(filepath.Join(t.TempDir(), "state.db"))
			if err != nil {
				t.Fatal(err)
			}
			defer db.Close()
			db.SetMaxOpenConns(24)
			tt.setup(t, db)

			const contenders = 20
			type result struct {
				reservation *LLMBudgetReservation
				err         error
			}
			start := make(chan struct{})
			results := make(chan result, contenders)
			var wg sync.WaitGroup
			for i := 0; i < contenders; i++ {
				wg.Add(1)
				go func() {
					defer wg.Done()
					<-start
					r, err := db.ReserveLLMBudget(tt.kind, tt.estimate, tt.limits)
					results <- result{r, err}
				}()
			}
			close(start)
			wg.Wait()
			close(results)

			admitted := 0
			for result := range results {
				if result.err == nil {
					admitted++
					continue
				}
				var limit *LLMBudgetLimitError
				if !errors.As(result.err, &limit) || limit.Limit != tt.wantLimit {
					t.Fatalf("contender error = %T %v", result.err, result.err)
				}
			}
			if admitted != tt.wantAdmitted {
				t.Fatalf("admitted %d, want %d", admitted, tt.wantAdmitted)
			}
			var reserved int
			if err := db.QueryRow(`SELECT COALESCE(SUM(tokens),0) FROM llm_budget_reservations`).Scan(&reserved); err != nil {
				t.Fatal(err)
			}
			if reserved != tt.wantReserved {
				t.Fatalf("durable reservations = %d, want %d", reserved, tt.wantReserved)
			}
		})
	}
}
