package main

import (
	"crypto/sha256"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/maceip/clew/internal/adapters"
	"github.com/maceip/clew/internal/extract"
	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/llm"
)

func cmdWitness(args []string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: clew witness <transcript.jsonl>")
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	provider, note := a.provider()
	if provider == nil {
		return fmt.Errorf("cloud-session return needs an extraction provider: %s", note)
	}
	copyPath, err := preserveWitness(args[0])
	if err != nil {
		return err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	out, adapter, err := witnessTranscript(j, provider, copyPath, a.cfg.Surface, time.Now())
	if err != nil {
		return err
	}
	for i := 0; i < out.Redactions; i++ {
		a.db.Incr("redactions")
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("returned %d remembered change(s) from the %s cloud session\n", len(out.Entries), adapter.ID())
	return nil
}

func preserveWitness(source string) (string, error) {
	b, err := os.ReadFile(source)
	if err != nil {
		return "", err
	}
	if len(b) == 0 {
		return "", fmt.Errorf("cloud transcript is empty")
	}
	sum := sha256.Sum256(b)
	dir := filepath.Join(gitx.Home(), "witness")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return "", err
	}
	path := filepath.Join(dir, fmt.Sprintf("%x.jsonl", sum[:12]))
	if existing, readErr := os.ReadFile(path); readErr == nil {
		if string(existing) == string(b) {
			return path, nil
		}
		return "", fmt.Errorf("witness digest collision at %s", path)
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		return "", err
	}
	return path, nil
}

func witnessTranscript(j *journal.Journal, provider llm.Provider, file, surface string, now time.Time) (*extract.Outcome, adapters.Adapter, error) {
	adapter, err := strictTranscriptAdapter(file)
	if err != nil {
		return nil, nil, err
	}
	out, err := extract.Run(j, provider, adapter, file, 0, surface, now)
	if err != nil {
		return nil, nil, err
	}
	if out.Parked {
		return nil, nil, fmt.Errorf("cloud transcript was preserved but could not be distilled: %s", out.ParkReason)
	}
	return out, adapter, nil
}

// strictTranscriptAdapter accepts a transcript only when exactly one pinned
// parser consumes it without unknown record classes. That is deterministic
// format recognition, not heuristic guessing across cloud products.
func strictTranscriptAdapter(file string) (adapters.Adapter, error) {
	var matches []adapters.Adapter
	for _, adapter := range adapters.All() {
		delta, err := adapter.Parse(file, 0)
		if err != nil || delta == nil || len(delta.Messages) == 0 || len(delta.Unknown) != 0 {
			continue
		}
		matches = append(matches, adapter)
	}
	switch len(matches) {
	case 1:
		return matches[0], nil
	case 0:
		return nil, fmt.Errorf("cloud transcript matches no pinned clew session format")
	default:
		return nil, fmt.Errorf("cloud transcript matches more than one pinned format; refusing to guess")
	}
}
