package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"clew/internal/adapters"
	"clew/internal/config"
	"clew/internal/extract"
	"clew/internal/journal"
	"clew/internal/llm"
	"clew/internal/state"
)

type budgetProviderStub struct {
	result *llm.Result
	err    error
	seq    []struct {
		result *llm.Result
		err    error
	}
}

func (p *budgetProviderStub) Name() string { return "budget-stub" }
func (p *budgetProviderStub) Call(string) (*llm.Result, error) {
	if len(p.seq) > 0 {
		next := p.seq[0]
		p.seq = p.seq[1:]
		return next.result, next.err
	}
	return p.result, p.err
}

func budgetTestDB(t *testing.T) *state.DB {
	t.Helper()
	db, err := state.Open(filepath.Join(t.TempDir(), "state.db"))
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { db.Close() })
	return db
}

func TestBudgetedProviderSettlesActualAtomically(t *testing.T) {
	db := budgetTestDB(t)
	cfg := config.Default()
	p := newBudgetedProvider(&budgetProviderStub{result: &llm.Result{Text: `{}`, Tokens: 37}}, db, cfg, "differ", false, 0)
	if _, err := p.Call("small prompt"); err != nil {
		t.Fatal(err)
	}
	if p.Spent() != 37 || db.TokensToday("spent") != 37 || db.TokensToday("differ-spent") != 37 {
		t.Fatalf("settlement: wrapper=%d aggregate=%d differ=%d", p.Spent(), db.TokensToday("spent"), db.TokensToday("differ-spent"))
	}
}

func TestBudgetedProviderChargesReservationOnUnknownFailure(t *testing.T) {
	db := budgetTestDB(t)
	cfg := config.Default()
	p := newBudgetedProvider(&budgetProviderStub{err: errors.New("transport broke")}, db, cfg, "extraction", false, 0)
	_, err := p.Call("small prompt")
	if err == nil || !strings.Contains(err.Error(), "charged") {
		t.Fatalf("failure = %v, want loud conservative charge", err)
	}
	if p.Spent() == 0 || db.TokensToday("spent") != p.Spent() {
		t.Fatalf("failed call escaped meter: wrapper=%d db=%d", p.Spent(), db.TokensToday("spent"))
	}
}

func TestBudgetedProviderHonorsExplicitRunCapBeforeCall(t *testing.T) {
	db := budgetTestDB(t)
	cfg := config.Default()
	base := &budgetProviderStub{result: &llm.Result{Text: `{}`, Tokens: 1}}
	p := newBudgetedProvider(base, db, cfg, "backfill", false, 100)
	if _, err := p.Call("prompt"); err == nil || !strings.Contains(err.Error(), "explicit run budget") {
		t.Fatalf("small run cap admitted call: %v", err)
	}
	if db.TokensToday("spent") != 0 {
		t.Fatal("denied call changed spend")
	}
}

func TestExtractorRetryFailureStillChargesFirstAttemptAndFailedReservation(t *testing.T) {
	db := budgetTestDB(t)
	cfg := config.Default()
	base := &budgetProviderStub{seq: []struct {
		result *llm.Result
		err    error
	}{
		{result: &llm.Result{Text: "not json", Tokens: 10}},
		{err: errors.New("second attempt transport failure")},
	}}
	p := newBudgetedProvider(base, db, cfg, "backfill", false, 100_000)
	file := filepath.Join(t.TempDir(), "session.jsonl")
	raw := `{"type":"user","message":{"role":"user","content":"remember this"},"timestamp":"2026-08-16T12:00:00Z","cwd":"/repo","sessionId":"s1"}` + "\n"
	if err := os.WriteFile(file, []byte(raw), 0o600); err != nil {
		t.Fatal(err)
	}
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	if _, err := extract.Run(j, p, &adapters.Claude{}, file, 0, "test", time.Now()); err == nil {
		t.Fatal("retry transport failure was hidden")
	}
	if p.Spent() <= 10 || db.TokensToday("backfill-spent") != p.Spent() {
		t.Fatalf("partial retry spend escaped meter: wrapper=%d db=%d", p.Spent(), db.TokensToday("backfill-spent"))
	}
}
