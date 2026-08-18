// Package lineage implements the data layer behind the explicit
// `clew from <repo>` operation. It deliberately contains no discovery,
// watcher, manifest, or command wiring: listing is a pure ranking operation,
// and importing mutates only the successor journal plus one append-only link.
package lineage

import (
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/seed"
)

const linkDir = "lineage"

var (
	linkIDPattern = regexp.MustCompile(`^l[0-9A-HJKMNP-TV-Z]{26}$`)
	digestPattern = regexp.MustCompile(`^sha256:[0-9a-f]{64}$`)
)

// Candidate is one already-generated ambient seed. Location is an opaque
// resolver token (normally a canonical path, cached seed path, or repo id).
type Candidate struct {
	Snapshot *seed.Snapshot
	Location string
}

type Target struct {
	RepositoryID string
	Name         string
	Topics       []string
}

type RankedCandidate struct {
	Candidate
	Score        float64
	TopicOverlap float64
	Recency      float64
	Blatant      bool // safe only as a suggestion signal; never an action
}

// Rank combines binary-cosine topic overlap (65%) with the same 30-day
// recency curve used by context ranking (35%). Stable tie-breaks make output
// testable. The current repository is always excluded.
func Rank(target Target, candidates []Candidate, now time.Time) ([]RankedCandidate, error) {
	targetTopics := topicSet(append(append([]string(nil), target.Topics...), target.Name))
	var ranked []RankedCandidate
	for _, candidate := range candidates {
		if candidate.Snapshot == nil {
			return nil, fmt.Errorf("lineage candidate %q has no seed", candidate.Location)
		}
		if err := seed.Validate(candidate.Snapshot); err != nil {
			return nil, fmt.Errorf("lineage candidate %q: %w", candidate.Location, err)
		}
		if candidate.Snapshot.Repository.ID == target.RepositoryID {
			continue
		}
		candidateTopics := topicSet(append(append([]string(nil), candidate.Snapshot.Topics...), candidate.Snapshot.Repository.Name))
		intersection := intersectCount(targetTopics, candidateTopics)
		overlap := cosine(intersection, len(targetTopics), len(candidateTopics))
		at := candidate.Snapshot.ChangedAt
		if candidate.Snapshot.Lifecycle.At.After(at) {
			at = candidate.Snapshot.Lifecycle.At
		}
		ageDays := now.Sub(at).Hours() / 24
		if at.IsZero() {
			ageDays = math.Inf(1)
		} else if ageDays < 0 {
			ageDays = 0
		}
		recency := 0.0
		if !math.IsInf(ageDays, 1) {
			recency = 1 / (1 + ageDays/30)
		}
		nameMatch := normalizedName(target.Name) != "" && normalizedName(target.Name) == normalizedName(candidate.Snapshot.Repository.Name)
		ranked = append(ranked, RankedCandidate{
			Candidate: candidate, TopicOverlap: overlap, Recency: recency,
			Score:   .65*overlap + .35*recency,
			Blatant: nameMatch || (intersection >= 2 && overlap >= .6),
		})
	}
	sort.SliceStable(ranked, func(i, k int) bool {
		a, b := ranked[i], ranked[k]
		if a.Score != b.Score {
			return a.Score > b.Score
		}
		if a.TopicOverlap != b.TopicOverlap {
			return a.TopicOverlap > b.TopicOverlap
		}
		if a.Recency != b.Recency {
			return a.Recency > b.Recency
		}
		if a.Snapshot.Repository.Name != b.Snapshot.Repository.Name {
			return a.Snapshot.Repository.Name < b.Snapshot.Repository.Name
		}
		return a.Snapshot.Repository.ID < b.Snapshot.Repository.ID
	})
	return ranked, nil
}

// Summary is the compact no-argument candidate line. Tombstoned lifecycle is
// explicit; inactivity alone is never rendered as death.
func (c RankedCandidate) Summary() string {
	if c.Snapshot == nil {
		return "unavailable · 0 lessons · no seed"
	}
	lessons := len(c.Snapshot.Decisions) + len(c.Snapshot.Findings)
	word := "lessons"
	if lessons == 1 {
		word = "lesson"
	}
	name := strings.Join(strings.Fields(c.Snapshot.Repository.Name), " ")
	if c.Snapshot.Lifecycle.State == "tombstoned" {
		return fmt.Sprintf("%s · %d %s · died %s · tombstoned", name, lessons, word,
			c.Snapshot.Lifecycle.At.Format("Jan 2"))
	}
	changed := "unknown"
	if !c.Snapshot.ChangedAt.IsZero() {
		changed = c.Snapshot.ChangedAt.Format("Jan 2")
	}
	return fmt.Sprintf("%s · %d %s · changed %s · active", name, lessons, word, changed)
}

// Link is the durable declaration written to lineage/<id>.yaml in the
// successor's journal branch. The source repository is not mutated.
type Link struct {
	ID                  string             `yaml:"id"`
	TargetRepository    string             `yaml:"target_repository"`
	From                seed.Repository    `yaml:"from"`
	FromAncestors       []string           `yaml:"from_ancestors,omitempty"`
	SeedJournalRevision string             `yaml:"seed_journal_revision,omitempty"`
	SeedDigest          string             `yaml:"seed_digest"`
	SeedChangedAt       time.Time          `yaml:"seed_changed_at,omitempty"`
	ImportedEntries     []string           `yaml:"imported_entries,omitempty"`
	ImportedEvents      []string           `yaml:"imported_events,omitempty"`
	OrganBank           *seed.OrganBankPin `yaml:"organ_bank,omitempty"`
	By                  model.By           `yaml:"by"`
	At                  time.Time          `yaml:"at"`
}

func NewLinkID(at time.Time) string { return "l" + ids.NewEvent(at)[1:] }

func (l *Link) Validate() error {
	if l == nil {
		return fmt.Errorf("nil lineage link")
	}
	if !linkIDPattern.MatchString(l.ID) {
		return fmt.Errorf("invalid lineage link id %q", l.ID)
	}
	if strings.TrimSpace(l.TargetRepository) == "" {
		return fmt.Errorf("lineage link %s has no target repository", l.ID)
	}
	if strings.TrimSpace(l.From.ID) == "" || strings.TrimSpace(l.From.Name) == "" {
		return fmt.Errorf("lineage link %s has incomplete predecessor identity", l.ID)
	}
	l.FromAncestors = cleanIDs(l.FromAncestors)
	if l.From.ID == l.TargetRepository {
		return fmt.Errorf("lineage link %s points a repository to itself", l.ID)
	}
	for _, ancestor := range l.FromAncestors {
		if ancestor == l.TargetRepository {
			return fmt.Errorf("lineage link %s creates a cycle through %s", l.ID, ancestor)
		}
	}
	if !digestPattern.MatchString(l.SeedDigest) {
		return fmt.Errorf("lineage link %s has invalid seed digest", l.ID)
	}
	// Empty is accepted only for links written by the short-lived pre-I13
	// implementation. Every new link records the predecessor journal revision,
	// which is the semantic idempotency key across machines.
	if l.SeedJournalRevision != "" && !digestPattern.MatchString(l.SeedJournalRevision) {
		return fmt.Errorf("lineage link %s has invalid seed journal revision", l.ID)
	}
	if l.By.Who != "human" {
		return fmt.Errorf("lineage link %s was not explicitly declared by a human", l.ID)
	}
	if l.At.IsZero() {
		return fmt.Errorf("lineage link %s has zero timestamp", l.ID)
	}
	return nil
}

// AddLink appends a link with O_EXCL. Link files are immutable and ULID-named,
// matching the journal's conflict-free entry/event write law.
func AddLink(journalDir string, link *Link) error {
	if err := link.Validate(); err != nil {
		return err
	}
	dir := filepath.Join(journalDir, linkDir)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(link)
	if err != nil {
		return err
	}
	path := filepath.Join(dir, link.ID+".yaml")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		f.Close()
		return err
	}
	return f.Close()
}

// LoadLinks loads every durable link or returns an error naming the malformed
// file. Invalid lineage is never silently omitted from candidate/cycle logic.
func LoadLinks(journalDir string) ([]*Link, error) {
	dir := filepath.Join(journalDir, linkDir)
	files, err := os.ReadDir(dir)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	var out []*Link
	seen := map[string]string{}
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".yaml") {
			continue
		}
		path := filepath.Join(dir, file.Name())
		b, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}
		link := &Link{}
		if err := yaml.Unmarshal(b, link); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if err := link.Validate(); err != nil {
			return nil, fmt.Errorf("%s: %w", path, err)
		}
		if file.Name() != link.ID+".yaml" {
			return nil, fmt.Errorf("%s: filename does not match lineage id %s", path, link.ID)
		}
		if prior := seen[link.ID]; prior != "" {
			return nil, fmt.Errorf("duplicate lineage id %s in %s and %s", link.ID, prior, path)
		}
		seen[link.ID] = path
		out = append(out, link)
	}
	sort.Slice(out, func(i, k int) bool { return out[i].ID < out[k].ID })
	return out, nil
}

// AncestorIDs returns the transitive predecessor IDs that a successor seed
// must publish. A later project can then reject A→B→A using only B's
// already-generated seed; no predecessor contact or on-demand generation is
// required.
func AncestorIDs(journalDir string) ([]string, error) {
	links, err := LoadLinks(journalDir)
	if err != nil {
		return nil, err
	}
	var ids []string
	for _, link := range links {
		ids = append(ids, link.From.ID)
		ids = append(ids, link.FromAncestors...)
	}
	return cleanIDs(ids), nil
}

// RevisionTokens returns the exact append-only link identities that a seed's
// JournalRevision must cover. Flattened ancestor IDs alone lose the fact that
// a particular predecessor seed digest was explicitly selected.
func RevisionTokens(links []*Link) []string {
	var tokens []string
	for _, link := range links {
		if link == nil {
			continue
		}
		tokens = append(tokens, strings.Join([]string{
			link.ID, link.TargetRepository, link.From.ID, link.SeedJournalRevision, link.SeedDigest,
		}, "\x1f"))
	}
	sort.Strings(tokens)
	return tokens
}

type ImportRequest struct {
	Journal          *journal.Journal
	TargetRepository string
	Snapshot         *seed.Snapshot
	By               model.By
	At               time.Time
}

type ImportResult struct {
	Link            *Link
	NewEntries      int
	NewEvents       int
	AlreadyImported bool
}

// Import preflights the entire seed, writes missing immutable entries/events,
// and appends the lineage link last. A crash before the link is harmless: a
// rerun recognizes exact files and completes the link. An ID/content conflict
// fails before any write.
func Import(req ImportRequest) (*ImportResult, error) {
	if req.Journal == nil {
		return nil, fmt.Errorf("nil successor journal")
	}
	if strings.TrimSpace(req.TargetRepository) == "" {
		return nil, fmt.Errorf("target repository id is required")
	}
	if req.By.Who != "human" {
		return nil, fmt.Errorf("lineage import must be explicitly declared by a human")
	}
	if req.At.IsZero() {
		return nil, fmt.Errorf("lineage import timestamp is required")
	}
	if err := seed.Validate(req.Snapshot); err != nil {
		return nil, err
	}
	if req.Snapshot.Repository.ID == req.TargetRepository {
		return nil, fmt.Errorf("repository %s cannot inherit from itself", req.TargetRepository)
	}
	for _, ancestor := range req.Snapshot.Ancestors {
		if ancestor == req.TargetRepository {
			return nil, fmt.Errorf("lineage cycle: predecessor %s already descends from %s", req.Snapshot.Repository.ID, req.TargetRepository)
		}
	}
	digest, err := seed.Digest(req.Snapshot)
	if err != nil {
		return nil, err
	}
	links, err := LoadLinks(req.Journal.Dir)
	if err != nil {
		return nil, err
	}
	for _, link := range links {
		if link.TargetRepository == req.TargetRepository && link.From.ID == req.Snapshot.Repository.ID &&
			((link.SeedJournalRevision != "" && link.SeedJournalRevision == req.Snapshot.JournalRevision) ||
				(link.SeedJournalRevision == "" && link.SeedDigest == digest)) {
			return &ImportResult{Link: link, AlreadyImported: true}, nil
		}
	}

	entries := seedEntries(req.Snapshot)
	if err := preflightEntries(req.Journal, entries); err != nil {
		return nil, err
	}
	markers, err := lineageStatusEvents(req.Journal, req.Snapshot, digest, req.By, req.At)
	if err != nil {
		return nil, err
	}
	events := mergeEvents(seedEvents(req.Snapshot), markers)
	if err := preflightEvents(req.Journal, events); err != nil {
		return nil, err
	}

	result := &ImportResult{}
	for _, entry := range entries {
		if existing := req.Journal.Entries[entry.ID]; existing != nil {
			continue
		}
		clone, err := cloneEntry(entry)
		if err != nil {
			return nil, err
		}
		if err := req.Journal.AddEntry(clone); err != nil {
			return nil, err
		}
		result.NewEntries++
	}
	existingEvents := eventIndex(req.Journal)
	for _, event := range events {
		if existingEvents[event.ID] != nil {
			continue
		}
		clone, err := cloneEvent(event)
		if err != nil {
			return nil, err
		}
		if err := req.Journal.AddEvent(clone); err != nil {
			return nil, err
		}
		result.NewEvents++
	}

	link := &Link{
		ID: NewLinkID(req.At), TargetRepository: req.TargetRepository,
		From: req.Snapshot.Repository, FromAncestors: cleanIDs(req.Snapshot.Ancestors),
		SeedJournalRevision: req.Snapshot.JournalRevision, SeedDigest: digest,
		SeedChangedAt:   req.Snapshot.ChangedAt,
		ImportedEntries: entryIDs(entries), ImportedEvents: eventIDs(events),
		OrganBank: clonePin(req.Snapshot.OrganBank), By: req.By, At: req.At.UTC(),
	}
	if err := AddLink(req.Journal.Dir, link); err != nil {
		return nil, err
	}
	result.Link = link
	return result, nil
}

// lineageStatusEvents turns every negative seed grave into a durable,
// human-authored disposition. Source terminal events are still copied, but
// they are insufficient for derived states such as absence/expiry and for a
// supersession whose successor was outside the seed. These markers make the
// selected grave status self-contained in the successor journal.
func lineageStatusEvents(j *journal.Journal, snapshot *seed.Snapshot, digest string, by model.By, at time.Time) ([]*model.Event, error) {
	var out []*model.Event
	for i := range snapshot.Graveyard {
		grave := &snapshot.Graveyard[i]
		var matched *model.Event
		for _, existing := range j.EventsFor(grave.Entry.ID) {
			if existing.Kind != model.EvDisposition || existing.By.Who != "human" ||
				existing.PStr("lineage_from_repository") != snapshot.Repository.ID ||
				existing.PStr("lineage_seed_revision") != snapshot.JournalRevision {
				continue
			}
			if existing.PStr("lineage_status") != string(grave.Status) {
				return nil, fmt.Errorf("entry %s already has a different lineage status for predecessor revision %s", grave.Entry.ID, snapshot.JournalRevision)
			}
			matched = existing
			break
		}
		if matched != nil {
			out = append(out, matched)
			continue
		}
		out = append(out, &model.Event{
			ID: ids.NewEvent(at), Kind: model.EvDisposition, Entry: grave.Entry.ID,
			Payload: map[string]any{
				"lineage_status":          string(grave.Status),
				"lineage_from_repository": snapshot.Repository.ID,
				"lineage_seed_revision":   snapshot.JournalRevision,
				"lineage_seed_digest":     digest,
			},
			By: by, At: at.UTC(),
		})
	}
	return out, nil
}

func mergeEvents(groups ...[]*model.Event) []*model.Event {
	byID := map[string]*model.Event{}
	for _, group := range groups {
		for _, event := range group {
			byID[event.ID] = event
		}
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]*model.Event, 0, len(ids))
	for _, id := range ids {
		out = append(out, byID[id])
	}
	return out
}

func preflightEntries(j *journal.Journal, entries []*model.Entry) error {
	for _, incoming := range entries {
		existing := j.Entries[incoming.ID]
		if existing == nil {
			continue
		}
		equal, err := equivalentEntry(existing, incoming)
		if err != nil {
			return err
		}
		if !equal {
			return fmt.Errorf("entry %s already exists with different content; lineage import wrote nothing", incoming.ID)
		}
	}
	return nil
}

func preflightEvents(j *journal.Journal, events []*model.Event) error {
	existing := eventIndex(j)
	for _, incoming := range events {
		current := existing[incoming.ID]
		if current == nil {
			continue
		}
		equal, err := canonicalEqual(current, incoming)
		if err != nil {
			return err
		}
		if !equal {
			return fmt.Errorf("event %s already exists with different content; lineage import wrote nothing", incoming.ID)
		}
	}
	return nil
}

func equivalentEntry(a, b *model.Entry) (bool, error) {
	equal, err := canonicalEqual(a, b)
	if err != nil || equal {
		return equal, err
	}
	// Compatibility with the original manifest importer, which changed only
	// source.kind to carried while retaining the source ref/time/agent/surface.
	ac, err := cloneEntry(a)
	if err != nil {
		return false, err
	}
	bc, err := cloneEntry(b)
	if err != nil {
		return false, err
	}
	ac.Source.Kind, bc.Source.Kind = model.SrcCarried, model.SrcCarried
	return canonicalEqual(ac, bc)
}

func canonicalEqual(a, b any) (bool, error) {
	ab, err := yaml.Marshal(a)
	if err != nil {
		return false, err
	}
	bb, err := yaml.Marshal(b)
	if err != nil {
		return false, err
	}
	return string(ab) == string(bb), nil
}

func seedEntries(s *seed.Snapshot) []*model.Entry {
	byID := map[string]*model.Entry{}
	for i := range s.Decisions {
		entry := &s.Decisions[i].Entry
		byID[entry.ID] = entry
	}
	for i := range s.Findings {
		entry := &s.Findings[i].Entry
		byID[entry.ID] = entry
	}
	for i := range s.Graveyard {
		entry := &s.Graveyard[i].Entry
		byID[entry.ID] = entry
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]*model.Entry, 0, len(ids))
	for _, id := range ids {
		out = append(out, byID[id])
	}
	return out
}

func seedEvents(s *seed.Snapshot) []*model.Event {
	byID := map[string]*model.Event{}
	for i := range s.Graveyard {
		for k := range s.Graveyard[i].Events {
			event := &s.Graveyard[i].Events[k]
			byID[event.ID] = event
		}
	}
	for i := range s.Exhibits {
		event := &s.Exhibits[i]
		byID[event.ID] = event
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	out := make([]*model.Event, 0, len(ids))
	for _, id := range ids {
		out = append(out, byID[id])
	}
	return out
}

func eventIndex(j *journal.Journal) map[string]*model.Event {
	out := make(map[string]*model.Event, len(j.Events))
	for _, event := range j.Events {
		out[event.ID] = event
	}
	return out
}

func entryIDs(entries []*model.Entry) []string {
	out := make([]string, 0, len(entries))
	for _, entry := range entries {
		out = append(out, entry.ID)
	}
	return out
}

func eventIDs(events []*model.Event) []string {
	out := make([]string, 0, len(events))
	for _, event := range events {
		out = append(out, event.ID)
	}
	return out
}

func topicSet(values []string) map[string]bool {
	out := map[string]bool{}
	for _, value := range values {
		for _, token := range seed.Tokens(value) {
			out[token] = true
		}
	}
	return out
}

func intersectCount(a, b map[string]bool) int {
	if len(a) > len(b) {
		a, b = b, a
	}
	n := 0
	for value := range a {
		if b[value] {
			n++
		}
	}
	return n
}

func cosine(intersection, a, b int) float64 {
	if intersection == 0 || a == 0 || b == 0 {
		return 0
	}
	return float64(intersection) / math.Sqrt(float64(a*b))
}

func normalizedName(value string) string {
	return strings.Join(seed.Tokens(value), "-")
}

func cleanIDs(values []string) []string {
	set := map[string]bool{}
	for _, value := range values {
		if value = strings.TrimSpace(value); value != "" {
			set[value] = true
		}
	}
	out := make([]string, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Strings(out)
	return out
}

func cloneEntry(in *model.Entry) (*model.Entry, error) {
	b, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	out := &model.Entry{}
	if err := yaml.Unmarshal(b, out); err != nil {
		return nil, err
	}
	return out, nil
}

func cloneEvent(in *model.Event) (*model.Event, error) {
	b, err := yaml.Marshal(in)
	if err != nil {
		return nil, err
	}
	out := &model.Event{}
	if err := yaml.Unmarshal(b, out); err != nil {
		return nil, err
	}
	return out, nil
}

func clonePin(in *seed.OrganBankPin) *seed.OrganBankPin {
	if in == nil {
		return nil
	}
	out := *in
	return &out
}
