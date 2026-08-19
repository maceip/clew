package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
	"github.com/maceip/clew/internal/owner"
	"github.com/maceip/clew/internal/scrub"
)

// cmdRedact: the one sanctioned rewrite (§4). Scrubs the entry file in place
// (the file survives, scrubbed — that is what prevents other machines from
// re-adding their poisoned copy), rewrites the branch to a fresh root, and
// force-pushes. The redaction itself is journaled (minus the secret).
func cmdRedact(args []string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: clew redact <entry-id>")
	}
	id := args[0]
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	syncResult, err := gitx.Sync(repo, regenFor(repo))
	if err != nil {
		return fmt.Errorf("sync project journal before redaction: %w", err)
	}
	lease, err := gitx.RewriteLeaseFromSync(repo, syncResult)
	if err != nil {
		return fmt.Errorf("prepare project redaction rewrite: %w", err)
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	if _, ok := j.Entries[id]; !ok {
		return fmt.Errorf("no entry %s", id)
	}

	stamp := time.Now().UTC()
	// Remove an ambient copy first. If the owner rewrite or force-push fails,
	// stop before reporting the narrower project source as clean.
	ownerRedacted, err := owner.Default(a.cfg.Owner.Remote).Redact(id, a.cfg.Surface, stamp)
	if err != nil {
		return err
	}

	if err := rewriteProjectRedactedRoot(repo, j, id, a.cfg.Surface, stamp, lease); err != nil {
		return err
	}
	a.db.Incr("redactions")
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	if ownerRedacted {
		laws, err := a.ownerLawBlock(false, stamp)
		if err != nil {
			return err
		}
		if err := fanoutOwnerLaws(a, laws, stamp); err != nil {
			return err
		}
	}
	fmt.Printf("redacted %s: branch rewritten and force-pushed; other watchers re-sync from the remote root.\n", id)
	if ownerRedacted {
		fmt.Println("the promoted owner-journal copy was scrubbed and removed from every reachable ambient context too.")
	}
	fmt.Println("note: local raw session files and parked slices on source machines are not touched — clean those by hand if the secret lives there too.")
	return nil
}

const maxProjectRedactionRewriteAttempts = 4

// rewriteProjectRedactedRoot is the project-scope counterpart to the owner
// CAS rewrite. A concurrent journal append invalidates the lease; the retry
// first syncs that append into the union, then scrubs and regenerates again.
// This keeps both immutable entries and generated projections out of the new
// root unless they were part of the verified union being rewritten.
func rewriteProjectRedactedRoot(repo string, j *journal.Journal, id, surface string, stamp time.Time, lease gitx.RewriteLease) error {
	var lastLeaseErr error
	for attempt := 1; attempt <= maxProjectRedactionRewriteAttempts; attempt++ {
		e := j.Entries[id]
		if e == nil {
			return fmt.Errorf("project redaction retry lost entry %s after sync", id)
		}
		// Scrub every free-text field on every attempt. An unrelated-root adoption
		// deliberately trusts the remote's copy of an existing file, so a sync
		// after a lease rejection may have restored the unsanitized bytes.
		e.Title = scrub.Mark
		e.Body = scrub.Mark
		e.Quote = scrub.Mark
		b, err := yaml.Marshal(e)
		if err != nil {
			return err
		}
		p := filepath.Join(j.Dir, "entries", id+".yaml")
		if err := os.WriteFile(p, b, 0o644); err != nil {
			return err
		}
		if !hasProjectRedaction(j, id) {
			at := stamp.Add(time.Duration(attempt-1) * time.Nanosecond)
			if err := j.AddEvent(&model.Event{
				ID: ids.NewEvent(at), Kind: model.EvDisposition, Entry: id,
				Payload: map[string]any{"redacted": true},
				By:      model.By{Who: "human", Surface: surface}, At: at,
			}); err != nil {
				return err
			}
		}
		if err := regenFor(repo)(j.Dir); err != nil {
			return fmt.Errorf("regenerate project journal before redaction rewrite: %w", err)
		}
		err = gitx.RewriteRoot(repo, "redact "+id, lease)
		if err == nil {
			return nil
		}
		if !gitx.IsRewriteLeaseError(err) {
			return fmt.Errorf("rewrite project journal for redaction: %w", err)
		}
		lastLeaseErr = err
		if attempt == maxProjectRedactionRewriteAttempts {
			break
		}
		syncResult, syncErr := gitx.Sync(repo, regenFor(repo))
		if syncErr != nil {
			return fmt.Errorf("resync project journal after rewrite lease rejection: %w", syncErr)
		}
		lease, err = gitx.RewriteLeaseFromSync(repo, syncResult)
		if err != nil {
			return fmt.Errorf("prepare retry of project redaction rewrite: %w", err)
		}
		wt, err := gitx.EnsureJournal(repo)
		if err != nil {
			return err
		}
		j, err = journal.Load(wt)
		if err != nil {
			return err
		}
	}
	return fmt.Errorf("rewrite project journal for redaction exhausted %d compare-and-swap attempts; the remote kept advancing and the redaction was not reported complete: %w",
		maxProjectRedactionRewriteAttempts, lastLeaseErr)
}

func hasProjectRedaction(j *journal.Journal, id string) bool {
	for _, event := range j.EventsFor(id) {
		if event.Kind == model.EvDisposition && event.PBool("redacted") {
			return true
		}
	}
	return false
}
