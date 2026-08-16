package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/model"
	"clew/internal/scrub"
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
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	e, ok := j.Entries[id]
	if !ok {
		return fmt.Errorf("no entry %s", id)
	}

	// Scrub all free-text fields fully (the secret's location is unknown).
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
	// Record the redaction as an event BEFORE the rewrite so it survives it.
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(time.Now()), Kind: model.EvDisposition, Entry: id,
		Payload: map[string]any{"redacted": true},
		By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: time.Now().UTC(),
	}); err != nil {
		return err
	}
	if err := gitx.RewriteRoot(repo, "redact "+id); err != nil {
		return err
	}
	a.db.Incr("redactions")
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("redacted %s: branch rewritten and force-pushed; other watchers re-sync from the remote root.\n", id)
	fmt.Println("note: local raw session files and parked slices on source machines are not touched — clean those by hand if the secret lives there too.")
	return nil
}
