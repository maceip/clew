// Package journal is the append-only journal store (JOURNAL_SPEC §3–§4).
// Layout inside a journal directory (the orphan-branch worktree):
//
//	entries/<id>.yaml   immutable once written
//	events/<id>.yaml    immutable once written
//	journal.md          generated rollup (deterministic projection)
//	digest.md           generated ≤4KB digest
//
// Every writer only ever adds immutable ULID-named files; no file is ever
// written by two parties — that is the whole concurrency story (§4).
package journal

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"clew/internal/model"
)

type Journal struct {
	Dir     string
	Entries map[string]*model.Entry
	Events  []*model.Event

	byEntry map[string][]*model.Event
	// LoadErrors lists files that failed to parse (I2: loud, never guessed).
	LoadErrors []string
}

func Load(dir string) (*Journal, error) {
	j := &Journal{
		Dir:     dir,
		Entries: map[string]*model.Entry{},
		byEntry: map[string][]*model.Event{},
	}
	if err := j.reload(); err != nil {
		return nil, err
	}
	return j, nil
}

func (j *Journal) reload() error {
	j.Entries = map[string]*model.Entry{}
	j.Events = nil
	j.byEntry = map[string][]*model.Event{}
	j.LoadErrors = nil

	entDir := filepath.Join(j.Dir, "entries")
	files, err := os.ReadDir(entDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		p := filepath.Join(entDir, f.Name())
		e := &model.Entry{}
		if err := readYAML(p, e); err != nil || e.Validate() != nil {
			j.LoadErrors = append(j.LoadErrors, p)
			continue
		}
		j.Entries[e.ID] = e
	}

	evDir := filepath.Join(j.Dir, "events")
	files, err = os.ReadDir(evDir)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".yaml") {
			continue
		}
		p := filepath.Join(evDir, f.Name())
		v := &model.Event{}
		if err := readYAML(p, v); err != nil || v.Validate() != nil {
			j.LoadErrors = append(j.LoadErrors, p)
			continue
		}
		j.Events = append(j.Events, v)
	}
	sort.Slice(j.Events, func(a, b int) bool { return j.Events[a].ID < j.Events[b].ID })
	for _, v := range j.Events {
		j.byEntry[v.Entry] = append(j.byEntry[v.Entry], v)
	}
	return nil
}

// Reload re-reads the directory (after a sync/rebase).
func (j *Journal) Reload() error { return j.reload() }

// EventsFor returns events targeting an entry, id-ordered (time-ordered).
func (j *Journal) EventsFor(id string) []*model.Event { return j.byEntry[id] }

// AddEntry validates and writes a new immutable entry file.
func (j *Journal) AddEntry(e *model.Entry) error {
	if err := e.Validate(); err != nil {
		return err
	}
	if _, dup := j.Entries[e.ID]; dup {
		return fmt.Errorf("entry %s already exists (entries are immutable)", e.ID)
	}
	p := filepath.Join(j.Dir, "entries", e.ID+".yaml")
	if err := writeYAMLNew(p, e); err != nil {
		return err
	}
	j.Entries[e.ID] = e
	return nil
}

// AddEvent validates and writes a new immutable event file.
func (j *Journal) AddEvent(v *model.Event) error {
	if err := v.Validate(); err != nil {
		return err
	}
	p := filepath.Join(j.Dir, "events", v.ID+".yaml")
	if err := writeYAMLNew(p, v); err != nil {
		return err
	}
	j.Events = append(j.Events, v)
	j.byEntry[v.Entry] = append(j.byEntry[v.Entry], v)
	return nil
}

// HasEvent reports whether an equivalent event already exists (dedupe for
// differ re-runs across machines; duplicates are harmless but noisy).
func (j *Journal) HasEvent(kind model.EventKind, entry string, payloadKey, payloadVal string) bool {
	for _, v := range j.byEntry[entry] {
		if v.Kind != kind {
			continue
		}
		if payloadKey == "" || v.PStr(payloadKey) == payloadVal {
			return true
		}
	}
	return false
}

func readYAML(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, out)
}

// writeYAMLNew writes a file that must not already exist (append-only law).
func writeYAMLNew(path string, in any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	b, err := yaml.Marshal(in)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}
