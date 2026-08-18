// Package owner stores the owner's project-agnostic laws in a dedicated
// append-only journal. A law is not a fifth journal entry type: it is a
// finding copied with its original identity and evidence into owner scope,
// then certified by an explicit human promotion disposition.
package owner

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/scrub"
)

const (
	// LawCap is the complete ambient owner-law block budget, including its
	// heading and provenance ids. Admission and rendering share one formatter.
	LawCap = 1024

	ScopeOwner        = "owner"
	ActionPromote     = "promote"
	actionCollision   = "promotion-collision"
	lawsHeading       = "## Owner laws (human-promoted, project-agnostic)\n"
	defaultOwnerRepo  = "owner"
	verifiedRemoteKey = "clew.owner-verified-remote"
)

// BudgetError means a promotion would make the exact ambient block exceed
// LawCap. It is returned before either the entry or certification is written.
type BudgetError struct {
	Required int
	Limit    int
}

func (e *BudgetError) Error() string {
	return fmt.Sprintf("owner law budget exceeded: promotion requires %d bytes, limit is %d", e.Required, e.Limit)
}

// RenderResult is both the injection payload and the loud overflow signal.
// Markdown is always <= LawCap. RequiredBytes measures the uncapped block.
// On an impossible-by-normal-admission concurrent overflow, the oldest
// certified prefix remains ambient and Omitted names the newer laws withheld.
type RenderResult struct {
	Markdown      string
	Included      []string
	Omitted       []string
	RequiredBytes int
	Overflow      bool
}

// PromotionResult reports the durable owner-journal result. Added is false
// for an idempotent retry of an already-certified identical finding.
type PromotionResult struct {
	Journal            *journal.Journal
	Sync               *gitx.SyncResult
	Render             RenderResult
	Added              bool
	CertificationAdded bool
}

// Store is one owner-scope journal backed by a normal git repository. The
// repository is deliberately not part of state.repos: it must be synced, but
// never scanned as a project for sessions, commits, or archaeology.
type Store struct {
	RepoPath string
	Remote   string

	// beforeRewrite is a deterministic concurrency seam used by owner package
	// tests. Production stores leave it nil. It runs at most once, after the
	// redactor captured its sync lease and before its first root push.
	beforeRewrite func() error
}

// Default returns the machine owner's conventional store under CLEW_HOME.
func Default(remote string) *Store {
	return New(filepath.Join(gitx.Home(), defaultOwnerRepo), remote)
}

// New constructs a store at an explicit path (primarily useful to tests and
// embedders). Ensure performs all filesystem and git mutations.
func New(repoPath, remote string) *Store {
	return &Store{RepoPath: filepath.Clean(repoPath), Remote: strings.TrimSpace(remote)}
}

// Ensure creates/adopts the local owner repository, checks its configured
// origin without silently repointing it, and ensures the append-only journal
// worktree exists. With no remote the journal remains fully functional and
// local-only, matching ordinary project-journal behavior.
func (s *Store) Ensure() (string, error) {
	if s == nil || strings.TrimSpace(s.RepoPath) == "" || s.RepoPath == "." {
		return "", fmt.Errorf("owner journal: empty repository path")
	}
	if err := ensureRepository(s.RepoPath); err != nil {
		return "", err
	}
	if s.Remote != "" {
		current, err := gitx.Run(s.RepoPath, "remote", "get-url", "origin")
		switch {
		case err != nil:
			if _, addErr := gitx.Run(s.RepoPath, "remote", "add", "origin", s.Remote); addErr != nil {
				return "", fmt.Errorf("owner journal remote: %w", addErr)
			}
		case current != s.Remote:
			return "", fmt.Errorf("owner journal origin is %q, config owner.remote is %q; refusing to repoint it", current, s.Remote)
		}
	}
	wt, err := gitx.EnsureJournal(s.RepoPath)
	if err != nil {
		return "", fmt.Errorf("ensure owner journal: %w", err)
	}
	return wt, nil
}

// Open loads the local owner journal without a network sync.
func (s *Store) Open() (*journal.Journal, error) {
	wt, err := s.Ensure()
	if err != nil {
		return nil, err
	}
	return journal.Load(wt)
}

// Sync applies the same commit/fetch/rebase/push protocol as a project
// journal, then returns a freshly loaded owner journal.
func (s *Store) Sync() (*journal.Journal, *gitx.SyncResult, error) {
	if _, err := s.Ensure(); err != nil {
		return nil, nil, err
	}
	res, err := gitx.Sync(s.RepoPath, regenerate)
	if err != nil {
		return nil, res, fmt.Errorf("sync owner journal: %w", err)
	}
	j, err := s.Open()
	if err != nil {
		return nil, res, err
	}
	return j, res, nil
}

// Promote certifies one project finding as an owner law. It syncs before the
// admission check so the 1KB decision includes remote laws, preserves the
// finding's exact id/source/quote, writes a human disposition as the sole
// certification boundary, then syncs the addition. Source is never mutated.
func (s *Store) Promote(source *journal.Journal, id, fromRepo, surface string, at time.Time) (*PromotionResult, error) {
	if source == nil {
		return nil, fmt.Errorf("promote owner law: nil source journal")
	}
	e := source.Entries[id]
	if e == nil {
		return nil, fmt.Errorf("promote owner law: no entry %s", id)
	}
	if e.Type != model.Finding {
		return nil, fmt.Errorf("promote owner law: %s is a %s, not a finding", id, e.Type)
	}
	computed := journal.Compute(source, normalizedTime(at))[id]
	if computed == nil || !journal.Live(computed.Status) {
		return nil, fmt.Errorf("promote owner law: finding %s is not live", id)
	}
	if computed.Withheld {
		return nil, fmt.Errorf("promote owner law: finding %s is withheld as imperative content; rewrite it as a declarative finding before promotion", id)
	}
	// A project confirmation only releases an entry into that project's
	// context. Ambient law has a much wider blast radius, so promotion always
	// re-checks the immutable raw fields (including Title) and requires a
	// rewritten finding if they are directive-shaped.
	if journal.Imperative(e) {
		return nil, fmt.Errorf("promote owner law: finding %s contains raw imperative content; rewrite it as a declarative finding first", id)
	}
	if computed.Tainted || e.UtteranceBy == model.ByToolResult {
		return nil, fmt.Errorf("promote owner law: finding %s came from tool_result and cannot become ambient law", id)
	}
	if isRedacted(source, id) || e.Title == scrub.Mark || e.Body == scrub.Mark || e.Quote == scrub.Mark {
		return nil, fmt.Errorf("promote owner law: finding %s was redacted and cannot become ambient law", id)
	}

	ownerJournal, preSync, err := s.Sync()
	if err != nil {
		return nil, err
	}
	if err := s.RequireVerifiedSync(preSync, "promote owner law"); err != nil {
		return nil, err
	}
	if err := s.MarkVerifiedCache(); err != nil {
		return nil, err
	}
	if isRedacted(ownerJournal, id) {
		return nil, fmt.Errorf("promote owner law: %s has an owner-scope redaction tombstone and cannot be promoted", id)
	}
	copy := cloneEntry(e)
	if existing := ownerJournal.Entries[id]; existing != nil && !entriesEqual(existing, copy) {
		return nil, fmt.Errorf("promote owner law: owner entry %s already exists with different content", id)
	}
	if IsCertified(ownerJournal, id) {
		if promotionCollision(ownerJournal, id) {
			return nil, fmt.Errorf("promote owner law: %s has a prior distributed admission collision; create a new finding after making budget room", id)
		}
		st := journal.Compute(ownerJournal, normalizedTime(at))[id]
		if st == nil || !journal.Live(st.Status) {
			return nil, fmt.Errorf("promote owner law: owner entry %s was already certified but is no longer live; create a superseding finding", id)
		}
		rendered := Render(ownerJournal, normalizedTime(at))
		if err := verifyAmbient(id, rendered); err != nil {
			return nil, err
		}
		return &PromotionResult{Journal: ownerJournal, Sync: preSync, Render: rendered}, nil
	}
	if existing := ownerJournal.Entries[id]; existing != nil {
		st := journal.Compute(ownerJournal, normalizedTime(at))[id]
		if st == nil || !journal.Live(st.Status) {
			return nil, fmt.Errorf("promote owner law: owner entry %s is no longer live; create a superseding finding", id)
		}
	}

	stamp := normalizedTime(at)
	cert := &model.Event{
		ID: ids.NewEvent(stamp), Kind: model.EvDisposition, Entry: id,
		Payload: map[string]any{
			"action":     ActionPromote,
			"scope":      ScopeOwner,
			"from_entry": id,
			"from_repo":  fromRepo,
		},
		By: model.By{Who: "human", Surface: surface}, At: stamp,
	}
	for _, existing := range ownerJournal.Entries {
		if existing.Supersedes == copy.ID {
			return nil, fmt.Errorf("promote owner law: owner entry %s is already superseded by %s", copy.ID, existing.ID)
		}
	}
	prospective := renderProspective(ownerJournal, copy, cert)
	if prospective.RequiredBytes > LawCap {
		return nil, &BudgetError{Required: prospective.RequiredBytes, Limit: LawCap}
	}

	added := false
	if ownerJournal.Entries[id] == nil {
		if err := ownerJournal.AddEntry(copy); err != nil {
			return nil, err
		}
		added = true
	}
	if err := ownerJournal.AddEvent(cert); err != nil {
		return nil, err
	}
	postSync, err := gitx.Sync(s.RepoPath, regenerate)
	if err != nil {
		return nil, fmt.Errorf("sync promoted owner law: %w", err)
	}
	// Fetch once more after the attempted push. If another machine admitted a
	// law from the same pre-sync budget snapshot, this pass observes the union
	// (including a deferred non-fast-forward push) before the caller records the
	// source-side promotion disposition.
	ownerJournal, verifySync, err := s.Sync()
	if err != nil {
		return nil, fmt.Errorf("verify promoted owner law after sync: %w", err)
	}
	postSync = mergeSyncResults(preSync, postSync, verifySync)
	if isRedacted(ownerJournal, id) {
		// A redaction tombstone may have raced admission after preflight. The
		// candidate has already entered the local/remote union, so suppression
		// alone is insufficient: scrub it through the sanctioned root rewrite
		// before returning and before the source project records promotion.
		if _, redactErr := s.Redact(id, surface, stamp.Add(2*time.Nanosecond)); redactErr != nil {
			return nil, fmt.Errorf("promote owner law: %s raced an owner redaction tombstone; promotion was not recorded, and cleanup failed: %w", id, redactErr)
		}
		return nil, fmt.Errorf("promote owner law: %s raced an owner redaction tombstone; the owner copy was scrubbed and promotion was not recorded", id)
	}
	if err := s.RequireVerifiedSync(postSync, "promote owner law"); err != nil {
		return nil, s.quarantinePromotion(ownerJournal, id, stamp, err)
	}
	if err := s.MarkVerifiedCache(); err != nil {
		return nil, s.quarantinePromotion(ownerJournal, id, stamp, err)
	}
	rendered := Render(ownerJournal, stamp)
	if err := verifyAmbient(id, rendered); err != nil {
		return nil, s.quarantinePromotion(ownerJournal, id, stamp, err)
	}
	return &PromotionResult{
		Journal:            ownerJournal,
		Sync:               postSync,
		Render:             rendered,
		Added:              added,
		CertificationAdded: true,
	}, nil
}

// Render returns only live findings carrying a human owner-promotion
// disposition. Raw additions to the owner branch, extractor suggestions, and
// machine-authored dispositions can never cross the ambient-law boundary.
func Render(j *journal.Journal, now time.Time) RenderResult {
	if j == nil {
		return RenderResult{}
	}
	st := journal.Compute(j, normalizedTime(now))
	var laws []law
	for id, e := range j.Entries {
		cert, ok := certification(j, id)
		c := st[id]
		if !ok || e == nil || e.Type != model.Finding || c == nil || !journal.Live(c.Status) || c.Withheld || c.Tainted ||
			journal.Imperative(e) || isRedacted(j, id) || promotionCollision(j, id) {
			continue
		}
		laws = append(laws, law{Entry: e, CertifiedBy: cert.ID})
	}
	return renderLaws(laws)
}

type law struct {
	Entry       *model.Entry
	CertifiedBy string
}

func renderProspective(j *journal.Journal, candidate *model.Entry, cert *model.Event) RenderResult {
	st := journal.Compute(j, cert.At)
	laws := make([]law, 0, len(j.Entries)+1)
	for id, e := range j.Entries {
		ce, ok := certification(j, id)
		c := st[id]
		if !ok || e == nil || e.Type != model.Finding || c == nil || !journal.Live(c.Status) || c.Withheld || c.Tainted ||
			journal.Imperative(e) || isRedacted(j, id) || promotionCollision(j, id) {
			continue
		}
		// Adding candidate would immediately supersede this owner law.
		if candidate.Supersedes != "" && id == candidate.Supersedes {
			continue
		}
		laws = append(laws, law{Entry: e, CertifiedBy: ce.ID})
	}
	laws = append(laws, law{Entry: candidate, CertifiedBy: cert.ID})
	return renderLaws(laws)
}

func renderLaws(laws []law) RenderResult {
	if len(laws) == 0 {
		return RenderResult{}
	}
	sort.Slice(laws, func(i, k int) bool {
		if laws[i].CertifiedBy != laws[k].CertifiedBy {
			return laws[i].CertifiedBy < laws[k].CertifiedBy
		}
		return laws[i].Entry.ID < laws[k].Entry.ID
	})
	lines := make([]string, len(laws))
	required := len(lawsHeading)
	for i, item := range laws {
		lines[i] = lawLine(item.Entry)
		required += len(lines[i])
	}

	var out strings.Builder
	out.WriteString(lawsHeading)
	result := RenderResult{RequiredBytes: required, Overflow: required > LawCap}
	for i, item := range laws {
		line := lines[i]
		// Overflow keeps an oldest-certified prefix. Once one law cannot fit,
		// all newer laws are omitted even if a later short line could squeeze
		// in; newer promotions can never evict older ambient laws.
		if out.Len()+len(line) > LawCap || len(result.Omitted) > 0 {
			result.Omitted = append(result.Omitted, item.Entry.ID)
			continue
		}
		out.WriteString(line)
		result.Included = append(result.Included, item.Entry.ID)
	}
	if len(result.Included) > 0 {
		result.Markdown = out.String()
	}
	return result
}

func lawLine(e *model.Entry) string {
	// The docket presents Quote as exact, reversible evidence. Inject that same
	// byte sequence, not extractor-authored Title/Body prose, after the owner
	// certifies it. The wrapper is provenance; Quote itself is not normalized.
	line := fmt.Sprintf("- [owner:%s] %s", e.ID, e.Quote)
	if !strings.HasSuffix(line, "\n") {
		line += "\n"
	}
	return line
}

func verifyAmbient(id string, rendered RenderResult) error {
	if rendered.Overflow || len(rendered.Omitted) > 0 || !containsID(rendered.Included, id) {
		return fmt.Errorf("promote owner law: post-sync admission failed for %s: required=%d limit=%d included=%s omitted=%s; source promotion was not recorded",
			id, rendered.RequiredBytes, LawCap, strings.Join(rendered.Included, ","), strings.Join(rendered.Omitted, ","))
	}
	return nil
}

// RequireVerifiedSync is intentionally stricter than ordinary journal sync.
// Project journals remain useful offline, but owner-scope absence, redaction,
// and globally capped certification cannot be decided against an unknown
// remote law set.
func (s *Store) RequireVerifiedSync(result *gitx.SyncResult, operation string) error {
	if gitx.RemoteName(s.RepoPath) == "" {
		return nil // one local owner store has no distributed admission race
	}
	if strings.TrimSpace(operation) == "" {
		operation = "owner journal operation"
	}
	if result == nil {
		return fmt.Errorf("%s: remote owner state is unavailable", operation)
	}
	for _, note := range result.Notes {
		if strings.HasPrefix(note, "fetch failed:") || strings.HasPrefix(note, "push deferred:") {
			return fmt.Errorf("%s: remote owner state was not verified: %s", operation, note)
		}
	}
	return nil
}

// MarkVerifiedCache binds the local cached owner journal to its current
// remote URL after a verified sync. The marker lives in local git config, so
// deleting/recreating the store cannot inherit a stale machine-state flag.
func (s *Store) MarkVerifiedCache() error {
	remote := gitx.RemoteName(s.RepoPath)
	if remote == "" {
		return nil
	}
	url, err := gitx.Run(s.RepoPath, "remote", "get-url", remote)
	if err != nil {
		return fmt.Errorf("mark owner cache verified: %w", err)
	}
	if _, err := gitx.Run(s.RepoPath, "config", "--local", verifiedRemoteKey, url); err != nil {
		return fmt.Errorf("mark owner cache verified: %w", err)
	}
	return nil
}

// RequireVerifiedCache permits network-free birth only after this exact owner
// repository has successfully synchronized with its current remote at least
// once. Local-only owner stores need no marker.
func (s *Store) RequireVerifiedCache(operation string) error {
	remote := gitx.RemoteName(s.RepoPath)
	if remote == "" {
		return nil
	}
	url, err := gitx.Run(s.RepoPath, "remote", "get-url", remote)
	if err != nil {
		return fmt.Errorf("%s: read owner remote: %w", operation, err)
	}
	verified, err := gitx.Run(s.RepoPath, "config", "--local", "--get", verifiedRemoteKey)
	if err != nil || verified != url {
		return fmt.Errorf("%s: local owner cache has not been verified against %s", operation, url)
	}
	return nil
}

func promotionCollision(j *journal.Journal, id string) bool {
	if j == nil {
		return false
	}
	for _, event := range j.EventsFor(id) {
		if event.Kind == model.EvDisposition && event.PStr("scope") == ScopeOwner &&
			event.PStr("action") == actionCollision {
			return true
		}
	}
	return false
}

func (s *Store) quarantinePromotion(j *journal.Journal, id string, at time.Time, cause error) error {
	if j == nil {
		return cause
	}
	if !promotionCollision(j, id) {
		event := &model.Event{
			ID: ids.NewEvent(at.Add(time.Nanosecond)), Kind: model.EvDisposition, Entry: id,
			Payload: map[string]any{
				"scope": ScopeOwner, "action": actionCollision,
				"reason": "distributed owner-law admission was not verified",
			},
			By: model.By{Who: "clew", Surface: "owner-admission"}, At: at.Add(time.Nanosecond),
		}
		if err := j.AddEvent(event); err != nil {
			return fmt.Errorf("%w; additionally failed to quarantine %s: %v", cause, id, err)
		}
	}
	if _, err := gitx.Sync(s.RepoPath, regenerate); err != nil {
		return fmt.Errorf("%w; %s is locally quarantined but its collision marker could not sync: %v", cause, id, err)
	}
	return fmt.Errorf("%w; %s was quarantined and is not ambient; create a new finding after resolving the owner-law budget", cause, id)
}

func containsID(ids []string, want string) bool {
	for _, id := range ids {
		if id == want {
			return true
		}
	}
	return false
}

func mergeSyncResults(results ...*gitx.SyncResult) *gitx.SyncResult {
	var merged *gitx.SyncResult
	for _, result := range results {
		if result == nil {
			continue
		}
		if merged == nil {
			merged = &gitx.SyncResult{}
		}
		merged.Committed = merged.Committed || result.Committed
		merged.Pushed = merged.Pushed || result.Pushed
		merged.Pulled = merged.Pulled || result.Pulled
		merged.Adopted = merged.Adopted || result.Adopted
		merged.Notes = append(merged.Notes, result.Notes...)
	}
	return merged
}

func certification(j *journal.Journal, id string) (*model.Event, bool) {
	var first *model.Event
	for _, event := range j.EventsFor(id) {
		if event.Kind != model.EvDisposition || event.By.Who != "human" ||
			event.PStr("scope") != ScopeOwner || event.PStr("action") != ActionPromote {
			continue
		}
		if first == nil || event.ID < first.ID {
			first = event
		}
	}
	return first, first != nil
}

// IsCertified reports whether id carries the explicit human owner-promotion
// disposition. It intentionally says nothing about current liveness or
// whether a concurrent overflow can fit the law into the ambient block.
func IsCertified(j *journal.Journal, id string) bool {
	if j == nil {
		return false
	}
	_, ok := certification(j, id)
	return ok
}

func isRedacted(j *journal.Journal, id string) bool {
	if j == nil {
		return false
	}
	for _, event := range j.EventsFor(id) {
		if event.Kind == model.EvDisposition && event.PBool("redacted") {
			return true
		}
	}
	return false
}

// Redact removes a promoted finding's free text from owner scope as part of
// the same sanctioned history rewrite used by project journals. It runs
// before the project rewrite so a failure can never leave the broader,
// ambient copy intact after the source copy was reported clean.
func (s *Store) Redact(id, surface string, at time.Time) (bool, error) {
	j, syncResult, err := s.Sync()
	if err != nil {
		return false, fmt.Errorf("sync owner journal before redaction: %w", err)
	}
	if err := s.RequireVerifiedSync(syncResult, "redact owner law"); err != nil {
		return false, err
	}
	if err := s.MarkVerifiedCache(); err != nil {
		return false, err
	}
	e := j.Entries[id]
	lease, err := gitx.RewriteLeaseFromSync(s.RepoPath, syncResult)
	if err != nil {
		return false, fmt.Errorf("prepare owner redaction rewrite: %w", err)
	}
	stamp := normalizedTime(at)
	if !isRedacted(j, id) {
		if err := j.AddEvent(&model.Event{
			ID: ids.NewEvent(stamp), Kind: model.EvDisposition, Entry: id,
			Payload: map[string]any{"redacted": true, "scope": ScopeOwner},
			By:      model.By{Who: "human", Surface: surface}, At: stamp,
		}); err != nil {
			return false, err
		}
	}

	// Even an absent owner copy gets a durable tombstone. Otherwise a stale or
	// concurrent machine could promote the same project entry after this
	// command concluded that no ambient copy existed.
	if e == nil {
		_, err := gitx.Sync(s.RepoPath, regenerate)
		if err != nil {
			return false, fmt.Errorf("publish owner redaction tombstone: %w", err)
		}
		verified, secondSync, err := s.Sync()
		if err != nil {
			return false, fmt.Errorf("verify owner redaction tombstone: %w", err)
		}
		if err := s.RequireVerifiedSync(secondSync, "publish owner redaction tombstone"); err != nil {
			return false, err
		}
		if err := s.MarkVerifiedCache(); err != nil {
			return false, err
		}
		if !isRedacted(verified, id) {
			return false, fmt.Errorf("owner redaction tombstone postcondition failed for %s", id)
		}
		j = verified
		e = j.Entries[id]
		if e == nil {
			return false, nil
		}
		lease, err = gitx.RewriteLeaseFromSync(s.RepoPath, secondSync)
		if err != nil {
			return false, fmt.Errorf("prepare owner redaction rewrite after tombstone publication: %w", err)
		}
		// A promotion raced the tombstone publication. It is already suppressed
		// from Render; continue through the root rewrite so its sensitive bytes
		// do not survive in owner history either.
	}

	clean, err := s.rewriteRedactedRoot(j, id, lease)
	if err != nil {
		return false, err
	}
	if containsID(Render(clean, stamp).Included, id) {
		return false, fmt.Errorf("owner redaction postcondition failed: %s remains ambient", id)
	}
	return true, nil
}

const maxRedactionRewriteAttempts = 4

// rewriteRedactedRoot performs the sanctioned owner history rewrite as a
// compare-and-swap. A lease rejection means a new append-only entry arrived
// after our last sync; resyncing alone is not enough because an unrelated-root
// adoption can restore the remote's unsanitized target file. Every retry
// therefore reloads the union, verifies the tombstone, and scrubs again.
func (s *Store) rewriteRedactedRoot(j *journal.Journal, id string, lease gitx.RewriteLease) (*journal.Journal, error) {
	var lastLeaseErr error
	for attempt := 1; attempt <= maxRedactionRewriteAttempts; attempt++ {
		if !isRedacted(j, id) {
			return nil, fmt.Errorf("rewrite owner journal for redaction: owner tombstone for %s was lost before attempt %d", id, attempt)
		}
		e := j.Entries[id]
		if e == nil {
			return j, nil
		}
		if err := scrubEntryFile(j, e); err != nil {
			return nil, err
		}
		if err := regenerate(j.Dir); err != nil {
			return nil, fmt.Errorf("regenerate owner journal before redaction rewrite: %w", err)
		}
		if hook := s.beforeRewrite; hook != nil {
			s.beforeRewrite = nil
			if err := hook(); err != nil {
				return nil, fmt.Errorf("owner redaction concurrency hook: %w", err)
			}
		}
		err := gitx.RewriteRoot(s.RepoPath, "redact "+id, lease)
		if err == nil {
			return s.Open()
		}
		if !gitx.IsRewriteLeaseError(err) {
			return nil, fmt.Errorf("rewrite owner journal for redaction: %w", err)
		}
		lastLeaseErr = err
		if attempt == maxRedactionRewriteAttempts {
			break
		}

		union, syncResult, syncErr := s.Sync()
		if syncErr != nil {
			return nil, fmt.Errorf("resync owner journal after rewrite lease rejection: %w", syncErr)
		}
		if verifyErr := s.RequireVerifiedSync(syncResult, "retry owner redaction rewrite"); verifyErr != nil {
			return nil, verifyErr
		}
		if markErr := s.MarkVerifiedCache(); markErr != nil {
			return nil, markErr
		}
		if !isRedacted(union, id) {
			return nil, fmt.Errorf("retry owner redaction rewrite: owner tombstone for %s disappeared after sync", id)
		}
		lease, err = gitx.RewriteLeaseFromSync(s.RepoPath, syncResult)
		if err != nil {
			return nil, fmt.Errorf("prepare retry of owner redaction rewrite: %w", err)
		}
		j = union
	}
	return nil, fmt.Errorf("rewrite owner journal for redaction exhausted %d compare-and-swap attempts; the remote kept advancing and the source project was not rewritten: %w",
		maxRedactionRewriteAttempts, lastLeaseErr)
}

func scrubEntryFile(j *journal.Journal, e *model.Entry) error {
	if j == nil || e == nil {
		return fmt.Errorf("scrub owner entry: missing journal or entry")
	}
	e.Title = scrub.Mark
	e.Body = scrub.Mark
	e.Quote = scrub.Mark
	b, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(j.Dir, "entries", e.ID+".yaml"), b, 0o644)
}

func cloneEntry(e *model.Entry) *model.Entry {
	copy := *e
	copy.Tags = append([]string(nil), e.Tags...)
	copy.Affects = append([]string(nil), e.Affects...)
	if e.Env != nil {
		env := *e.Env
		copy.Env = &env
	}
	return &copy
}

func entriesEqual(a, b *model.Entry) bool {
	ab, aerr := yaml.Marshal(a)
	bb, berr := yaml.Marshal(b)
	return aerr == nil && berr == nil && bytes.Equal(ab, bb)
}

func normalizedTime(at time.Time) time.Time {
	if at.IsZero() {
		return time.Now().UTC()
	}
	return at.UTC()
}

func ensureRepository(path string) error {
	if err := os.MkdirAll(path, 0o755); err != nil {
		return fmt.Errorf("create owner journal repository: %w", err)
	}
	if root, err := gitx.Root(path); err == nil {
		pathInfo, pathErr := os.Stat(path)
		rootInfo, rootErr := os.Stat(root)
		if pathErr == nil && rootErr == nil && os.SameFile(pathInfo, rootInfo) {
			return nil
		}
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		return err
	}
	if len(entries) > 0 {
		return fmt.Errorf("owner journal path %s exists and is not an empty git repository", path)
	}
	if _, err := gitx.Run(path, "init", "-q"); err != nil {
		return fmt.Errorf("initialize owner journal repository: %w", err)
	}
	return nil
}

func regenerate(wt string) error {
	j, err := journal.Load(wt)
	if err != nil {
		return err
	}
	now := time.Now()
	return journal.WriteProjections(j, journal.Compute(j, now), now)
}

// IsBudgetError lets command integration distinguish an admission refusal
// without importing the concrete error type.
func IsBudgetError(err error) bool {
	var target *BudgetError
	return errors.As(err, &target)
}
