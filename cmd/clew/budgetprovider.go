package main

import (
	"errors"
	"fmt"

	"clew/internal/config"
	"clew/internal/llm"
	"clew/internal/state"
)

// A provider response should be a compact extraction/link JSON object, never
// prose or a transcript. The byte ceiling is included in the prompt and its
// worst-case token cost is reserved before the process/API call starts.
const (
	llmOutputByteCeiling = 16 * 1024
	llmTokenOverhead     = 256
)

type budgetedProvider struct {
	base      llm.Provider
	db        *state.DB
	cfg       *config.Config
	kind      string
	liveRatio bool
	runCap    int // zero means only the machine limits apply
	charged   int
}

func newBudgetedProvider(base llm.Provider, db *state.DB, cfg *config.Config, kind string, liveRatio bool, runCap int) *budgetedProvider {
	return &budgetedProvider{base: base, db: db, cfg: cfg, kind: kind, liveRatio: liveRatio, runCap: runCap}
}

func (p *budgetedProvider) Name() string { return p.base.Name() }
func (p *budgetedProvider) Spent() int   { return p.charged }

func (p *budgetedProvider) Call(prompt string) (*llm.Result, error) {
	if disabled := p.db.Get("llm-disabled:" + p.kind); disabled != "" {
		return nil, fmt.Errorf("%s provider paused after budget-contract failure: %s", p.kind, disabled)
	}
	boundedPrompt := prompt + fmt.Sprintf("\n\nHARD OUTPUT BUDGET: return at most %d UTF-8 bytes of strict JSON. If the answer would exceed it, return the same schema with fewer items.", llmOutputByteCeiling)
	// BPE tokens are byte sequences, so prompt bytes + bounded output bytes is
	// a conservative token reservation (plus fixed envelope overhead).
	reserve := len([]byte(boundedPrompt)) + llmOutputByteCeiling + llmTokenOverhead
	if p.runCap > 0 && p.charged+reserve > p.runCap {
		err := fmt.Errorf("explicit run budget: charged %d + worst-case call %d > %d", p.charged, reserve, p.runCap)
		p.db.Set("llm-error:"+p.kind, err.Error())
		return nil, err
	}
	limits := state.LLMBudgetLimits{DailyCapTokens: p.cfg.Extractor.DailyCapTokens}
	if p.liveRatio {
		limits.LiveSessionPct = p.cfg.Extractor.SessionPct
	}
	claim, err := p.db.ReserveLLMBudget(p.kind, reserve, limits)
	if err != nil {
		p.db.Set("llm-error:"+p.kind, err.Error())
		return nil, err
	}
	res, callErr := p.base.Call(boundedPrompt)
	if callErr != nil {
		// A failed external CLI/API call does not provide trustworthy usage.
		// Charge its reserved upper bound so repeated failures cannot bypass I9.
		settleErr := p.db.SettleLLMBudget(claim.ID, claim.Tokens)
		p.charged += claim.Tokens
		if settleErr != nil {
			err := fmt.Errorf("provider failed: %v; budget settlement failed: %w", callErr, settleErr)
			p.db.Set("llm-error:"+p.kind, err.Error())
			return nil, err
		}
		err := fmt.Errorf("provider failed (charged %d-token reservation): %w", claim.Tokens, callErr)
		p.db.Set("llm-error:"+p.kind, err.Error())
		return nil, err
	}
	actual := res.Tokens
	if actual <= 0 {
		actual = (len(boundedPrompt) + len(res.Text)) / 4
	}
	settleErr := p.db.SettleLLMBudget(claim.ID, actual)
	p.charged += actual
	if settleErr != nil {
		var overrun *state.LLMBudgetOverrunError
		if errors.As(settleErr, &overrun) {
			p.db.Set("llm-disabled:"+p.kind, overrun.Error())
		}
		p.db.Set("llm-error:"+p.kind, settleErr.Error())
		return nil, settleErr
	}
	if len([]byte(res.Text)) > llmOutputByteCeiling {
		err := fmt.Errorf("provider returned %d bytes, exceeding the %d-byte output contract", len([]byte(res.Text)), llmOutputByteCeiling)
		p.db.Set("llm-disabled:"+p.kind, err.Error())
		p.db.Set("llm-error:"+p.kind, err.Error())
		return nil, err
	}
	p.db.Set("llm-error:"+p.kind, "")
	return res, nil
}
