// Package seed builds and reads the ambient, self-contained project seed.
//
// SEED.md has two audiences. YAML front matter is the exact, digest-checked
// carry payload; the Markdown body is a calm human projection. Keeping both in
// one file means a watched project's carry-kit is portable without consulting
// its live journal at the moment somebody decides to restart.
package seed

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
	"unicode"

	"gopkg.in/yaml.v3"

	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
)

const Format = "clew.seed/v1"

// Repository is the stable identity carried by a seed. ID should be stable
// across machines (normally derived from the canonical remote); Path is only a
// local resolution hint and is never used as identity.
type Repository struct {
	ID     string `yaml:"id"`
	Name   string `yaml:"name"`
	Path   string `yaml:"path,omitempty"`
	Remote string `yaml:"remote,omitempty"`
}

// Lifecycle is explicit metadata. A project is never inferred to be dead from
// inactivity: only State=tombstoned permits the human projection to say died.
type Lifecycle struct {
	State string    `yaml:"state"` // active | tombstoned
	At    time.Time `yaml:"at,omitempty"`
	Note  string    `yaml:"note,omitempty"`
}

// OrganBankPin points back to reusable source without claiming that a dirty
// working tree is present at the pinned commit.
type OrganBankPin struct {
	Remote string    `yaml:"remote,omitempty"`
	Commit string    `yaml:"commit"`
	Dirty  bool      `yaml:"dirty,omitempty"`
	At     time.Time `yaml:"at,omitempty"`
}

// Lesson is a current decision or finding. Status is captured so warnings such
// as suspect and possible-contradiction survive the carry boundary.
type Lesson struct {
	Entry  model.Entry    `yaml:"entry"`
	Status journal.Status `yaml:"status"`
}

// Grave records negative knowledge plus the terminal events that made it
// non-live. The complete entry is retained so source quote/ref/time survive.
type Grave struct {
	Entry  model.Entry    `yaml:"entry"`
	Status journal.Status `yaml:"status"`
	Events []model.Event  `yaml:"events,omitempty"`
}

// Snapshot is the machine payload embedded in SEED.md. Ancestors contains
// stable repository IDs and lets an explicit import reject lineage cycles.
type Snapshot struct {
	Repository      Repository    `yaml:"repository"`
	JournalRevision string        `yaml:"journal_revision"`
	ChangedAt       time.Time     `yaml:"changed_at,omitempty"`
	Lifecycle       Lifecycle     `yaml:"lifecycle"`
	Topics          []string      `yaml:"topics,omitempty"`
	Ancestors       []string      `yaml:"ancestors,omitempty"`
	Decisions       []Lesson      `yaml:"decisions,omitempty"`
	Findings        []Lesson      `yaml:"findings,omitempty"`
	Graveyard       []Grave       `yaml:"graveyard,omitempty"`
	Exhibits        []model.Event `yaml:"exhibits,omitempty"`
	OrganBank       *OrganBankPin `yaml:"organ_bank,omitempty"`
}

// BuildInput supplies the non-journal facts a watcher already knows. Build is
// deterministic: ChangedAt is derived from journal content, not wall clock.
type BuildInput struct {
	Repository       Repository
	Lifecycle        Lifecycle
	Topics           []string
	Ancestors        []string
	LineageRevision  []string
	LineageChangedAt time.Time
	OrganBank        *OrganBankPin
}

type envelope struct {
	Format   string   `yaml:"format"`
	Digest   string   `yaml:"digest"`
	Snapshot Snapshot `yaml:"snapshot"`
}

// Build selects the ambient carry surface: current decisions and findings,
// terminal negative knowledge, evidence exhibits, and an optional source pin.
// Active questions/intents remain project work; the deliberate manifest path
// can carry those during a large restart.
func Build(j *journal.Journal, st map[string]*journal.Computed, in BuildInput) (*Snapshot, error) {
	if j == nil {
		return nil, fmt.Errorf("nil journal")
	}
	s := &Snapshot{
		Repository: in.Repository,
		Lifecycle:  in.Lifecycle,
		Topics:     normalizeSet(in.Topics),
		Ancestors:  cleanSet(in.Ancestors),
		OrganBank:  clonePin(in.OrganBank),
	}
	if s.Lifecycle.State == "" {
		s.Lifecycle.State = "active"
	}
	revision, err := journalRevision(j, s.Ancestors, in.LineageRevision, s.Lifecycle)
	if err != nil {
		return nil, err
	}
	s.JournalRevision = revision

	carried := map[string]bool{}
	ids := make([]string, 0, len(j.Entries))
	for id := range j.Entries {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	for _, id := range ids {
		e := j.Entries[id]
		c := st[id]
		if c == nil {
			return nil, fmt.Errorf("entry %s has no computed state", id)
		}
		if t := e.Created(); t.After(s.ChangedAt) {
			s.ChangedAt = t.UTC()
		}
		switch {
		case graveworthy(j, id, c.Status):
			g := Grave{Entry: *cloneEntry(e), Status: carriedGraveStatus(j, e, c.Status)}
			for _, event := range j.EventsFor(id) {
				if terminal(event) {
					g.Events = append(g.Events, *cloneEvent(event))
				}
			}
			s.Graveyard = append(s.Graveyard, g)
			carried[id] = true
		case e.Type == model.Decision && journal.Live(c.Status):
			s.Decisions = append(s.Decisions, Lesson{Entry: *cloneEntry(e), Status: c.Status})
			carried[id] = true
		case e.Type == model.Finding && journal.Live(c.Status):
			s.Findings = append(s.Findings, Lesson{Entry: *cloneEntry(e), Status: c.Status})
			carried[id] = true
		}
	}

	for _, event := range j.Events {
		if event.At.After(s.ChangedAt) {
			s.ChangedAt = event.At.UTC()
		}
		if event.Kind == model.EvEvidence && carried[event.Entry] {
			s.Exhibits = append(s.Exhibits, *cloneEvent(event))
		}
	}
	if in.LineageChangedAt.After(s.ChangedAt) {
		s.ChangedAt = in.LineageChangedAt.UTC()
	}
	normalize(s)
	if err := Validate(s); err != nil {
		return nil, err
	}
	return s, nil
}

func graveworthy(j *journal.Journal, id string, status journal.Status) bool {
	switch status {
	case journal.StSuperseded, journal.StExpired, journal.StDropped, journal.StAbsent:
		return true
	}
	for _, event := range j.EventsFor(id) {
		if event.Kind == model.EvReject ||
			(event.Kind == model.EvDisposition && event.PStr("disposition") == "dropped") {
			return true
		}
	}
	return false
}

// carriedGraveStatus converts the generic manifest disposition "dropped" to
// the terminal status that belongs to the entry's type. A seed must never
// describe a decision as an intent-only "dropped" status (or vice versa),
// because the successor algebra could otherwise interpret the grave as live.
func carriedGraveStatus(j *journal.Journal, e *model.Entry, status journal.Status) journal.Status {
	for _, event := range j.EventsFor(e.ID) {
		if event.Kind == model.EvDisposition && event.PStr("disposition") == "dropped" {
			switch e.Type {
			case model.Decision, model.Finding:
				return journal.StSuperseded
			case model.Question:
				return journal.StExpired
			case model.Intent:
				return journal.StDropped
			}
		}
	}
	return status
}

func terminal(event *model.Event) bool {
	if event == nil {
		return false
	}
	switch event.Kind {
	case model.EvReject, model.EvSupersede, model.EvAnswer:
		return true
	case model.EvDisposition:
		return event.PStr("disposition") == "dropped"
	}
	return false
}

// Validate rejects malformed or ambiguous carry payloads before an importer
// can write any journal files.
func Validate(s *Snapshot) error {
	if s == nil {
		return fmt.Errorf("nil seed")
	}
	if strings.TrimSpace(s.Repository.ID) == "" {
		return fmt.Errorf("seed repository has no stable id")
	}
	if strings.TrimSpace(s.Repository.Name) == "" {
		return fmt.Errorf("seed repository has no name")
	}
	if !validDigest(s.JournalRevision) {
		return fmt.Errorf("seed has invalid journal revision %q", s.JournalRevision)
	}
	switch s.Lifecycle.State {
	case "active":
	case "tombstoned":
		if s.Lifecycle.At.IsZero() {
			return fmt.Errorf("tombstoned seed has no lifecycle timestamp")
		}
	default:
		return fmt.Errorf("unknown seed lifecycle %q", s.Lifecycle.State)
	}

	entries := map[string]bool{}
	add := func(e *model.Entry, want model.EntryType) error {
		if err := e.Validate(); err != nil {
			return err
		}
		if want != "" && e.Type != want {
			return fmt.Errorf("entry %s is %s, want %s", e.ID, e.Type, want)
		}
		if entries[e.ID] {
			return fmt.Errorf("entry %s appears in more than one seed section", e.ID)
		}
		entries[e.ID] = true
		return nil
	}
	for i := range s.Decisions {
		if !decisionStatus(s.Decisions[i].Status) {
			return fmt.Errorf("decision %s has invalid status %q", s.Decisions[i].Entry.ID, s.Decisions[i].Status)
		}
		if err := add(&s.Decisions[i].Entry, model.Decision); err != nil {
			return fmt.Errorf("decision %d: %w", i, err)
		}
	}
	for i := range s.Findings {
		if !findingStatus(s.Findings[i].Status) {
			return fmt.Errorf("finding %s has invalid status %q", s.Findings[i].Entry.ID, s.Findings[i].Status)
		}
		if err := add(&s.Findings[i].Entry, model.Finding); err != nil {
			return fmt.Errorf("finding %d: %w", i, err)
		}
	}
	events := map[string]bool{}
	checkEvent := func(event *model.Event) error {
		if err := event.Validate(); err != nil {
			return err
		}
		if events[event.ID] {
			return fmt.Errorf("event %s appears more than once in seed", event.ID)
		}
		events[event.ID] = true
		return nil
	}
	for i := range s.Graveyard {
		g := &s.Graveyard[i]
		if !graveStatus(g.Entry.Type, g.Status) {
			return fmt.Errorf("grave %s (%s) has incompatible status %q", g.Entry.ID, g.Entry.Type, g.Status)
		}
		if err := add(&g.Entry, ""); err != nil {
			return fmt.Errorf("grave %d: %w", i, err)
		}
		for k := range g.Events {
			event := &g.Events[k]
			if event.Entry != g.Entry.ID {
				return fmt.Errorf("grave %s contains event %s for %s", g.Entry.ID, event.ID, event.Entry)
			}
			if !terminal(event) {
				return fmt.Errorf("grave %s contains non-terminal event %s", g.Entry.ID, event.ID)
			}
			if err := checkEvent(event); err != nil {
				return err
			}
		}
	}
	for i := range s.Exhibits {
		event := &s.Exhibits[i]
		if event.Kind != model.EvEvidence {
			return fmt.Errorf("exhibit %s is %s, want evidence", event.ID, event.Kind)
		}
		if !entries[event.Entry] {
			return fmt.Errorf("exhibit %s targets entry %s outside seed", event.ID, event.Entry)
		}
		if err := checkEvent(event); err != nil {
			return err
		}
	}
	if s.OrganBank != nil && strings.TrimSpace(s.OrganBank.Commit) == "" {
		return fmt.Errorf("organ-bank pin has no commit")
	}
	return nil
}

// Render emits digest-checked YAML front matter followed by a readable
// projection. It is deterministic for a normalized Snapshot.
func Render(snapshot *Snapshot) ([]byte, error) {
	if snapshot == nil {
		return nil, fmt.Errorf("nil seed")
	}
	s := cloneSnapshot(snapshot)
	normalize(s)
	if err := Validate(s); err != nil {
		return nil, err
	}
	digest, err := Digest(s)
	if err != nil {
		return nil, err
	}
	front, err := yaml.Marshal(envelope{Format: Format, Digest: digest, Snapshot: *s})
	if err != nil {
		return nil, err
	}
	var out bytes.Buffer
	out.WriteString("---\n")
	out.Write(front)
	out.WriteString("---\n")
	out.WriteString(renderHuman(s))
	return out.Bytes(), nil
}

// Parse reads SEED.md and verifies that its machine payload matches its
// digest. The Markdown projection is deliberately not parsed or trusted.
func Parse(data []byte) (*Snapshot, error) {
	if !bytes.HasPrefix(data, []byte("---\n")) {
		return nil, fmt.Errorf("seed has no YAML front matter")
	}
	rest := data[len("---\n"):]
	end := bytes.Index(rest, []byte("\n---\n"))
	if end < 0 {
		return nil, fmt.Errorf("seed YAML front matter is not terminated")
	}
	var env envelope
	if err := yaml.Unmarshal(rest[:end+1], &env); err != nil {
		return nil, fmt.Errorf("parse seed front matter: %w", err)
	}
	if env.Format != Format {
		return nil, fmt.Errorf("unsupported seed format %q", env.Format)
	}
	normalize(&env.Snapshot)
	if err := Validate(&env.Snapshot); err != nil {
		return nil, err
	}
	digest, err := Digest(&env.Snapshot)
	if err != nil {
		return nil, err
	}
	if env.Digest != digest {
		return nil, fmt.Errorf("seed digest mismatch: got %s want %s", env.Digest, digest)
	}
	return &env.Snapshot, nil
}

func Read(path string) (*Snapshot, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return Parse(b)
}

// Write writes atomically only when the complete SEED.md bytes changed. The
// bool is false for an unchanged journal, allowing watcher tests to prove that
// seed generation is event-driven rather than a periodic rewrite.
func Write(path string, snapshot *Snapshot) (bool, error) {
	b, err := Render(snapshot)
	if err != nil {
		return false, err
	}
	if old, err := os.ReadFile(path); err == nil && bytes.Equal(old, b) {
		return false, nil
	} else if err != nil && !os.IsNotExist(err) {
		return false, err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return false, err
	}
	f, err := os.CreateTemp(filepath.Dir(path), ".seed-*.tmp")
	if err != nil {
		return false, err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if err := f.Chmod(0o644); err != nil {
		f.Close()
		return false, err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return false, err
	}
	if err := f.Sync(); err != nil {
		f.Close()
		return false, err
	}
	if err := f.Close(); err != nil {
		return false, err
	}
	if err := os.Rename(tmp, path); err != nil {
		return false, err
	}
	return true, nil
}

// WriteOnJournalChange is the watcher-facing write seam for I13. Repository
// facts such as README topics, HEAD, or dirty state may be sampled while
// building a Snapshot, but they cannot rewrite SEED.md until the durable
// journal revision (entries, events, lineage ancestry, lifecycle) changes.
// This is what makes seed maintenance journal-triggered rather than poll- or
// demand-triggered.
func WriteOnJournalChange(path string, snapshot *Snapshot) (bool, error) {
	if snapshot == nil {
		return false, fmt.Errorf("nil seed")
	}
	if err := Validate(snapshot); err != nil {
		return false, err
	}
	if old, err := Read(path); err == nil {
		if old.Repository.ID == snapshot.Repository.ID && old.JournalRevision == snapshot.JournalRevision {
			return false, nil
		}
	} else if !os.IsNotExist(err) {
		return false, fmt.Errorf("read existing seed before journal-gated write: %w", err)
	}
	return Write(path, snapshot)
}

// Digest returns the content identity used by durable lineage links.
func Digest(snapshot *Snapshot) (string, error) {
	if snapshot == nil {
		return "", fmt.Errorf("nil seed")
	}
	s := cloneSnapshot(snapshot)
	normalize(s)
	b, err := yaml.Marshal(s)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func journalRevision(j *journal.Journal, ancestors, lineageRevision []string, lifecycle Lifecycle) (string, error) {
	type revisionPayload struct {
		Entries         []model.Entry `yaml:"entries,omitempty"`
		Events          []model.Event `yaml:"events,omitempty"`
		Ancestors       []string      `yaml:"ancestors,omitempty"`
		LineageRevision []string      `yaml:"lineage_revision,omitempty"`
		Lifecycle       Lifecycle     `yaml:"lifecycle"`
	}
	payload := revisionPayload{
		Ancestors: cleanSet(ancestors), LineageRevision: cleanSet(lineageRevision), Lifecycle: lifecycle,
	}
	entryIDs := make([]string, 0, len(j.Entries))
	for id := range j.Entries {
		entryIDs = append(entryIDs, id)
	}
	sort.Strings(entryIDs)
	for _, id := range entryIDs {
		payload.Entries = append(payload.Entries, *cloneEntry(j.Entries[id]))
	}
	events := append([]*model.Event(nil), j.Events...)
	sort.Slice(events, func(i, k int) bool { return events[i].ID < events[k].ID })
	for _, event := range events {
		payload.Events = append(payload.Events, *cloneEvent(event))
	}
	b, err := yaml.Marshal(payload)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(b)
	return "sha256:" + hex.EncodeToString(sum[:]), nil
}

func validDigest(value string) bool {
	if !strings.HasPrefix(value, "sha256:") || len(value) != len("sha256:")+64 {
		return false
	}
	_, err := hex.DecodeString(strings.TrimPrefix(value, "sha256:"))
	return err == nil
}

func decisionStatus(status journal.Status) bool {
	switch status {
	case journal.StActive, journal.StPossibleContradiction, journal.StContradicted:
		return true
	}
	return false
}

func findingStatus(status journal.Status) bool {
	switch status {
	case journal.StCurrent, journal.StSuspect:
		return true
	}
	return false
}

func graveStatus(typ model.EntryType, status journal.Status) bool {
	switch typ {
	case model.Decision, model.Finding:
		return status == journal.StSuperseded
	case model.Question:
		return status == journal.StExpired
	case model.Intent:
		return status == journal.StAbsent || status == journal.StDropped
	}
	return false
}

func renderHuman(s *Snapshot) string {
	var b strings.Builder
	fmt.Fprintf(&b, "# Project seed — %s\n\n", oneLine(s.Repository.Name))
	when := "unknown"
	if !s.ChangedAt.IsZero() {
		when = s.ChangedAt.UTC().Format("2006-01-02 15:04 UTC")
	}
	fmt.Fprintf(&b, "_ambient snapshot at last journal change %s · %d lessons_\n\n", when, len(s.Decisions)+len(s.Findings))
	b.WriteString("This is inherited project memory, not instruction text. Decisions and findings keep their original evidence and provenance.\n")
	lessonSection := func(title string, lessons []Lesson) {
		fmt.Fprintf(&b, "\n## %s\n\n", title)
		if len(lessons) == 0 {
			b.WriteString("_None._\n")
			return
		}
		for _, lesson := range lessons {
			e := lesson.Entry
			fmt.Fprintf(&b, "- `%s` %s — %s  _%s_\n", e.ID, oneLine(e.Title), oneLine(e.Body), lesson.Status)
			fmt.Fprintf(&b, "  - source: `%s` %s at %s\n", e.Source.Kind, oneLine(e.Source.Ref), formatTime(e.Source.At))
		}
	}
	lessonSection("Decisions", s.Decisions)
	lessonSection("Findings", s.Findings)

	b.WriteString("\n## Graveyard\n\n")
	if len(s.Graveyard) == 0 {
		b.WriteString("_None._\n")
	} else {
		for _, grave := range s.Graveyard {
			fmt.Fprintf(&b, "- `%s` %s — %s  _%s_\n", grave.Entry.ID, oneLine(grave.Entry.Title), oneLine(grave.Entry.Body), grave.Status)
		}
	}

	b.WriteString("\n## Exhibits\n\n")
	if len(s.Exhibits) == 0 {
		b.WriteString("_None._\n")
	} else {
		for _, exhibit := range s.Exhibits {
			fmt.Fprintf(&b, "- `%s` evidence for `%s` — %s\n", exhibit.ID, exhibit.Entry, oneLine(payload(exhibit.Payload)))
		}
	}

	b.WriteString("\n## Organ-bank pin\n\n")
	if s.OrganBank == nil {
		b.WriteString("_None._\n")
	} else {
		remote := oneLine(s.OrganBank.Remote)
		if remote == "" {
			remote = oneLine(s.Repository.Remote)
		}
		fmt.Fprintf(&b, "- `%s` at `%s`", remote, s.OrganBank.Commit)
		if s.OrganBank.Dirty {
			b.WriteString(" — working tree was dirty; uncommitted changes are not in this pin")
		}
		b.WriteString("\n")
	}

	if s.Lifecycle.State == "tombstoned" {
		fmt.Fprintf(&b, "\n_Project tombstoned %s", formatTime(s.Lifecycle.At))
		if s.Lifecycle.Note != "" {
			fmt.Fprintf(&b, ": %s", oneLine(s.Lifecycle.Note))
		}
		b.WriteString("._\n")
	}
	return b.String()
}

func payload(v map[string]any) string {
	if len(v) == 0 {
		return "evidence recorded"
	}
	b, err := yaml.Marshal(v)
	if err != nil {
		return "evidence recorded"
	}
	return string(b)
}

func normalize(s *Snapshot) {
	if s == nil {
		return
	}
	s.Topics = normalizeSet(s.Topics)
	s.Ancestors = cleanSet(s.Ancestors)
	sort.Slice(s.Decisions, func(i, k int) bool { return s.Decisions[i].Entry.ID < s.Decisions[k].Entry.ID })
	sort.Slice(s.Findings, func(i, k int) bool { return s.Findings[i].Entry.ID < s.Findings[k].Entry.ID })
	sort.Slice(s.Graveyard, func(i, k int) bool { return s.Graveyard[i].Entry.ID < s.Graveyard[k].Entry.ID })
	for i := range s.Graveyard {
		sort.Slice(s.Graveyard[i].Events, func(a, b int) bool {
			return s.Graveyard[i].Events[a].ID < s.Graveyard[i].Events[b].ID
		})
	}
	sort.Slice(s.Exhibits, func(i, k int) bool { return s.Exhibits[i].ID < s.Exhibits[k].ID })
}

func normalizeSet(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		for _, token := range Tokens(value) {
			set[token] = true
		}
	}
	return sortedKeys(set)
}

func cleanSet(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			set[value] = true
		}
	}
	return sortedKeys(set)
}

func sortedKeys(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

var stopWords = map[string]bool{
	"and": true, "are": true, "for": true, "from": true, "into": true,
	"project": true, "repo": true, "repository": true, "the": true,
	"this": true, "with": true,
}

// Tokens is the deterministic, provider-free topic vocabulary shared by seed
// generation and lineage ranking.
func Tokens(text string) []string {
	fields := strings.FieldsFunc(strings.ToLower(text), func(r rune) bool {
		return !unicode.IsLetter(r) && !unicode.IsDigit(r)
	})
	set := map[string]bool{}
	for _, field := range fields {
		if len([]rune(field)) >= 3 && !stopWords[field] {
			set[field] = true
		}
	}
	return sortedKeys(set)
}

func cloneSnapshot(in *Snapshot) *Snapshot {
	if in == nil {
		return nil
	}
	b, _ := yaml.Marshal(in)
	out := &Snapshot{}
	_ = yaml.Unmarshal(b, out)
	return out
}

func cloneEntry(in *model.Entry) *model.Entry {
	if in == nil {
		return nil
	}
	out := *in
	out.Tags = append([]string(nil), in.Tags...)
	out.Affects = append([]string(nil), in.Affects...)
	if in.Env != nil {
		env := *in.Env
		out.Env = &env
	}
	return &out
}

func cloneEvent(in *model.Event) *model.Event {
	if in == nil {
		return nil
	}
	b, _ := yaml.Marshal(in)
	out := &model.Event{}
	_ = yaml.Unmarshal(b, out)
	return out
}

func clonePin(in *OrganBankPin) *OrganBankPin {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

func formatTime(at time.Time) string {
	if at.IsZero() {
		return "unknown time"
	}
	return at.UTC().Format("2006-01-02")
}
