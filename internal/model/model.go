// Package model defines the journal entry and event schema (JOURNAL_SPEC §3).
// Entry files are immutable once written; everything that happens to an entry
// afterwards is a separate append-only event file.
package model

import (
	"fmt"
	"time"
	"unicode/utf8"

	"clew/internal/ids"
)

type EntryType string

const (
	Decision EntryType = "decision"
	Finding  EntryType = "finding"
	Question EntryType = "question"
	Intent   EntryType = "intent"
)

type SourceKind string

const (
	SrcSession     SourceKind = "session"
	SrcCommit      SourceKind = "commit"
	SrcDoc         SourceKind = "doc"
	SrcHuman       SourceKind = "human"
	SrcCarried     SourceKind = "carried"
	SrcArchaeology SourceKind = "archaeology"
	SrcForeign     SourceKind = "foreign"
)

type UtteranceBy string

const (
	ByUser       UtteranceBy = "user"
	ByAssistant  UtteranceBy = "assistant"
	ByToolResult UtteranceBy = "tool_result" // tainted (§6.5)
)

type Source struct {
	Kind    SourceKind `yaml:"kind"`
	Ref     string     `yaml:"ref"`
	Agent   string     `yaml:"agent,omitempty"`
	Surface string     `yaml:"surface,omitempty"`
	At      time.Time  `yaml:"at"`
}

// Env records where a finding was measured (findings only).
type Env struct {
	Host    string `yaml:"host,omitempty"`
	HW      string `yaml:"hw,omitempty"`
	Dataset string `yaml:"dataset,omitempty"`
}

func (e *Env) Equal(o *Env) bool {
	if e == nil || o == nil {
		return e == o
	}
	return e.Host == o.Host && e.HW == o.HW && e.Dataset == o.Dataset
}

type Entry struct {
	ID          string      `yaml:"id"`
	Type        EntryType   `yaml:"type"`
	Title       string      `yaml:"title"`
	Body        string      `yaml:"body"`
	Quote       string      `yaml:"quote"`
	UtteranceBy UtteranceBy `yaml:"utterance_by"`
	Source      Source      `yaml:"source"`
	Confidence  float64     `yaml:"confidence"`
	Tags        []string    `yaml:"tags"`
	Env         *Env        `yaml:"env"`
	Affects     []string    `yaml:"affects"`
	// Asks is questions-only: who can answer, "human" or "any" (§3.1).
	Asks string `yaml:"asks,omitempty"`
	// Supersedes chains findings to the measurement they replace (§3.3).
	Supersedes string `yaml:"supersedes,omitempty"`
	// PromotionCandidate is a durable extractor proposal, never promotion
	// authority. Watchers can reconstruct the human docket card from the
	// journal after a state-db loss; only a human owner-scope disposition can
	// make the finding ambient.
	PromotionCandidate bool `yaml:"promotion_candidate,omitempty"`
}

// Created is the entry's creation instant: the source timestamp when present
// (preserves original time for carried/backfilled entries), else ULID time.
func (e *Entry) Created() time.Time {
	if !e.Source.At.IsZero() {
		return e.Source.At
	}
	return ids.Time(e.ID)
}

const (
	MaxTitle = 80  // runes (§3.2)
	MaxBody  = 400 // runes (§3.2)
)

func (e *Entry) Validate() error {
	if !ids.ValidEntry(e.ID) {
		return fmt.Errorf("entry %q: invalid id", e.ID)
	}
	switch e.Type {
	case Decision, Finding, Question, Intent:
	default:
		return fmt.Errorf("entry %s: unknown type %q", e.ID, e.Type)
	}
	if e.Title == "" {
		return fmt.Errorf("entry %s: empty title", e.ID)
	}
	if utf8.RuneCountInString(e.Title) > MaxTitle {
		return fmt.Errorf("entry %s: title exceeds %d chars", e.ID, MaxTitle)
	}
	if utf8.RuneCountInString(e.Body) > MaxBody {
		return fmt.Errorf("entry %s: body exceeds %d chars", e.ID, MaxBody)
	}
	// I7: no entry without evidence of utterance.
	if e.Quote == "" {
		return fmt.Errorf("entry %s: no quote (I7: no quote, no entry)", e.ID)
	}
	switch e.UtteranceBy {
	case ByUser, ByAssistant, ByToolResult:
	default:
		return fmt.Errorf("entry %s: unknown utterance_by %q", e.ID, e.UtteranceBy)
	}
	switch e.Source.Kind {
	case SrcSession, SrcCommit, SrcDoc, SrcHuman, SrcCarried, SrcArchaeology, SrcForeign:
	default:
		return fmt.Errorf("entry %s: unknown source.kind %q", e.ID, e.Source.Kind)
	}
	if e.Source.Ref == "" {
		return fmt.Errorf("entry %s: empty source.ref (I7)", e.ID)
	}
	if e.Confidence < 0 || e.Confidence > 1 {
		return fmt.Errorf("entry %s: confidence %v out of range", e.ID, e.Confidence)
	}
	if e.Source.Kind == SrcArchaeology && e.Confidence > 0.6 {
		return fmt.Errorf("entry %s: archaeology confidence capped at 0.6 (§5.3)", e.ID)
	}
	if e.Type != Finding && (e.Env != nil || len(e.Affects) > 0) {
		return fmt.Errorf("entry %s: env/affects are findings-only", e.ID)
	}
	if e.PromotionCandidate && e.Type != Finding {
		return fmt.Errorf("entry %s: promotion_candidate is findings-only", e.ID)
	}
	if e.Type != Question && e.Asks != "" {
		return fmt.Errorf("entry %s: asks is questions-only", e.ID)
	}
	if e.Type == Question && e.Asks != "" && e.Asks != "human" && e.Asks != "any" {
		return fmt.Errorf("entry %s: asks must be human|any", e.ID)
	}
	return nil
}

type EventKind string

const (
	EvEvidence          EventKind = "evidence"
	EvConfirm           EventKind = "confirm"
	EvReject            EventKind = "reject"
	EvSupersede         EventKind = "supersede"
	EvAnswer            EventKind = "answer"
	EvDisposition       EventKind = "disposition"
	EvEvidenceWithdrawn EventKind = "evidence-withdrawn"
)

type By struct {
	Who     string `yaml:"who"` // human | differ | extractor | <agent id>
	Surface string `yaml:"surface,omitempty"`
}

type Event struct {
	ID      string         `yaml:"id"`
	Kind    EventKind      `yaml:"kind"`
	Entry   string         `yaml:"entry"`
	Payload map[string]any `yaml:"payload"`
	By      By             `yaml:"by"`
	At      time.Time      `yaml:"at"`
}

func (v *Event) Validate() error {
	if !ids.ValidEvent(v.ID) {
		return fmt.Errorf("event %q: invalid id", v.ID)
	}
	switch v.Kind {
	case EvEvidence, EvConfirm, EvReject, EvSupersede, EvAnswer, EvDisposition, EvEvidenceWithdrawn:
	default:
		return fmt.Errorf("event %s: unknown kind %q", v.ID, v.Kind)
	}
	if v.Entry == "" {
		return fmt.Errorf("event %s: no target entry", v.ID)
	}
	if v.By.Who == "" {
		return fmt.Errorf("event %s: empty by.who", v.ID)
	}
	if v.At.IsZero() {
		return fmt.Errorf("event %s: zero timestamp", v.ID)
	}
	return nil
}

// PStr reads a string field out of an event payload.
func (v *Event) PStr(key string) string {
	if v.Payload == nil {
		return ""
	}
	if s, ok := v.Payload[key].(string); ok {
		return s
	}
	return ""
}

// PBool reads a bool field out of an event payload.
func (v *Event) PBool(key string) bool {
	if v.Payload == nil {
		return false
	}
	if b, ok := v.Payload[key].(bool); ok {
		return b
	}
	return false
}
