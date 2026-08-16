// Package proposal validates and stages untrusted journal contributions until
// a human accepts the complete batch. Staged entries never enter context.md.
package proposal

import (
	"bytes"
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/gitx"
	"clew/internal/journal"
	"clew/internal/model"
)

const (
	MaxBundleBytes = 10 << 20
	MaxEntries     = 1000
)

type Bundle struct {
	Version int            `yaml:"version"`
	Entries []*model.Entry `yaml:"entries"`
}

type Batch struct {
	ID      string         `yaml:"id"`
	Repo    string         `yaml:"repo"`
	Source  string         `yaml:"source"`
	Created time.Time      `yaml:"created"`
	Status  string         `yaml:"status"`
	Entries []*model.Entry `yaml:"entries"`
	Dir     string         `yaml:"-"`
}

type Manager struct {
	Root   string
	Client *http.Client
}

func New(root string) *Manager {
	return &Manager{Root: root, Client: &http.Client{Timeout: 15 * time.Second}}
}

func Default() *Manager { return New(filepath.Join(gitx.Home(), "proposals")) }

func (m *Manager) Stage(ctx context.Context, repo, source string) (*Batch, error) {
	entries, canonical, err := m.read(ctx, source)
	if err != nil {
		return nil, err
	}
	if len(entries) == 0 || len(entries) > MaxEntries {
		return nil, fmt.Errorf("proposal entries=%d; require 1..%d", len(entries), MaxEntries)
	}
	seen := make(map[string]bool, len(entries))
	foreign := make([]*model.Entry, 0, len(entries))
	for _, original := range entries {
		if original == nil {
			return nil, fmt.Errorf("proposal contains a null entry")
		}
		if err := original.Validate(); err != nil {
			return nil, fmt.Errorf("proposal entry: %w", err)
		}
		if seen[original.ID] {
			return nil, fmt.Errorf("proposal repeats entry %s", original.ID)
		}
		seen[original.ID] = true
		copy := *original
		copy.Source.Kind = model.SrcForeign
		copy.Source.Ref = canonical + "#" + string(original.Source.Kind) + ":" + original.Source.Ref
		if err := copy.Validate(); err != nil {
			return nil, fmt.Errorf("foreign proposal entry: %w", err)
		}
		foreign = append(foreign, &copy)
	}
	sort.Slice(foreign, func(i, k int) bool { return foreign[i].ID < foreign[k].ID })
	normalized, err := yaml.Marshal(Bundle{Version: 1, Entries: foreign})
	if err != nil {
		return nil, err
	}
	hash := sha256.Sum256(normalized)
	id := fmt.Sprintf("p%x", hash[:10])
	dir := filepath.Join(m.Root, gitx.RepoID(repo), id)
	if existing, err := m.Load(repo, id); err == nil {
		return existing, nil
	}
	batch := &Batch{ID: id, Repo: repo, Source: canonical, Created: time.Now().UTC(), Status: "pending", Entries: foreign, Dir: dir}
	if err := os.MkdirAll(dir, 0o700); err != nil {
		return nil, err
	}
	if err := writeYAMLExclusive(filepath.Join(dir, "batch.yaml"), batch); err != nil {
		if os.IsExist(err) {
			return m.Load(repo, id)
		}
		return nil, err
	}
	if err := os.WriteFile(filepath.Join(dir, "diff.md"), []byte(Diff(batch)), 0o600); err != nil {
		return nil, err
	}
	return batch, nil
}

func (m *Manager) Load(repo, id string) (*Batch, error) {
	if !validID(id) {
		return nil, fmt.Errorf("invalid proposal id %q", id)
	}
	dir := filepath.Join(m.Root, gitx.RepoID(repo), id)
	var batch Batch
	if err := decodeStrictFile(filepath.Join(dir, "batch.yaml"), &batch); err != nil {
		return nil, err
	}
	if batch.ID != id || batch.Repo != repo {
		return nil, fmt.Errorf("proposal %s identity mismatch", id)
	}
	batch.Dir = dir
	return &batch, nil
}

func (m *Manager) SetStatus(repo, id, status string) error {
	if status != "accepted" && status != "rejected" {
		return fmt.Errorf("invalid proposal status %q", status)
	}
	batch, err := m.Load(repo, id)
	if err != nil {
		return err
	}
	batch.Status = status
	b, err := yaml.Marshal(batch)
	if err != nil {
		return err
	}
	tmp := filepath.Join(batch.Dir, "batch.yaml.tmp")
	if err := os.WriteFile(tmp, b, 0o600); err != nil {
		return err
	}
	return os.Rename(tmp, filepath.Join(batch.Dir, "batch.yaml"))
}

// PushBranch publishes a proposal as additions on a branch based on
// clew/journal. Opening a PR with clew/journal as its base makes merge the
// explicit human-confirm boundary; it never grants direct journal writes.
func (m *Manager) PushBranch(repo string, batch *Batch, branch string) error {
	if branch == "" || branch == gitx.Branch {
		return fmt.Errorf("CLEW_PROPOSAL_BRANCH must name a non-journal branch")
	}
	if _, err := gitx.Run(repo, "check-ref-format", "--branch", branch); err != nil {
		return fmt.Errorf("invalid proposal branch %q: %w", branch, err)
	}
	if _, err := gitx.EnsureJournal(repo); err != nil {
		return err
	}
	remote := gitx.RemoteName(repo)
	if remote == "" {
		return fmt.Errorf("proposal branch requires a git remote")
	}
	tmpRoot, err := os.MkdirTemp(gitx.Home(), "proposal-worktree-")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpRoot)
	wt := filepath.Join(tmpRoot, "worktree")
	if _, err := gitx.Run(repo, "rev-parse", "--verify", "refs/heads/"+branch); err != nil {
		if _, err := gitx.Run(repo, "worktree", "add", "-q", "-b", branch, wt, gitx.Branch); err != nil {
			return err
		}
	} else if _, err := gitx.Run(repo, "worktree", "add", "-q", wt, branch); err != nil {
		return err
	}
	defer gitx.Run(repo, "worktree", "remove", "--force", wt)
	j, err := journal.Load(wt)
	if err != nil {
		return err
	}
	for _, entry := range batch.Entries {
		if existing := j.Entries[entry.ID]; existing != nil {
			if !reflect.DeepEqual(existing, entry) {
				return fmt.Errorf("proposal branch entry %s conflicts with existing content", entry.ID)
			}
			continue
		}
		copy := *entry
		if err := j.AddEntry(&copy); err != nil {
			return err
		}
	}
	now := time.Now().UTC()
	if err := journal.WriteProjections(j, journal.Compute(j, now), now); err != nil {
		return err
	}
	if status, err := gitx.Run(wt, "status", "--porcelain"); err != nil {
		return err
	} else if status != "" {
		if _, err := gitx.Run(wt, "add", "-A"); err != nil {
			return err
		}
		if _, err := gitx.Run(wt, "commit", "-q", "-m", "proposal: "+batch.ID); err != nil {
			return err
		}
	}
	if _, err := gitx.Run(wt, "push", "-q", remote, branch+":"+branch); err != nil {
		return err
	}
	return nil
}

func Diff(batch *Batch) string {
	var out strings.Builder
	fmt.Fprintf(&out, "# Proposal %s\n\nsource: %s\nentries: %d\n\n", batch.ID, batch.Source, len(batch.Entries))
	for _, entry := range batch.Entries {
		fmt.Fprintf(&out, "## + %s · %s · %s\n\n> %s\n\nprovenance: %s %s\n\n",
			entry.ID, entry.Type, entry.Title, strings.ReplaceAll(entry.Quote, "\n", " "), entry.Source.Kind, entry.Source.Ref)
	}
	return out.String()
}

func (m *Manager) DiffPath(repo, id string) (string, error) {
	batch, err := m.Load(repo, id)
	if err != nil {
		return "", err
	}
	return filepath.Join(batch.Dir, "diff.md"), nil
}

func (m *Manager) read(ctx context.Context, source string) ([]*model.Entry, string, error) {
	if parsed, err := url.Parse(source); err == nil && (parsed.Scheme == "http" || parsed.Scheme == "https") {
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
		if err != nil {
			return nil, "", err
		}
		resp, err := m.Client.Do(req)
		if err != nil {
			return nil, "", err
		}
		defer resp.Body.Close()
		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			return nil, "", fmt.Errorf("proposal URL returned %s", resp.Status)
		}
		b, err := io.ReadAll(io.LimitReader(resp.Body, MaxBundleBytes+1))
		if err != nil {
			return nil, "", err
		}
		if len(b) > MaxBundleBytes {
			return nil, "", fmt.Errorf("proposal exceeds %d bytes", MaxBundleBytes)
		}
		entries, err := decodeBundle(b)
		return entries, source, err
	}
	abs, err := filepath.Abs(source)
	if err != nil {
		return nil, "", err
	}
	info, err := os.Stat(abs)
	if err != nil {
		return nil, "", err
	}
	if info.IsDir() {
		entries, err := readDirectory(abs)
		return entries, abs, err
	}
	if info.Size() > MaxBundleBytes {
		return nil, "", fmt.Errorf("proposal exceeds %d bytes", MaxBundleBytes)
	}
	b, err := os.ReadFile(abs)
	if err != nil {
		return nil, "", err
	}
	entries, err := decodeBundle(b)
	return entries, abs, err
}

func readDirectory(dir string) ([]*model.Entry, error) {
	root := filepath.Join(dir, "entries")
	if info, err := os.Stat(root); err != nil || !info.IsDir() {
		root = dir
	}
	var files []string
	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.Type().IsRegular() && (strings.HasSuffix(d.Name(), ".yaml") || strings.HasSuffix(d.Name(), ".yml")) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)
	entries := make([]*model.Entry, 0, len(files))
	for _, file := range files {
		var entry model.Entry
		if err := decodeStrictFile(file, &entry); err != nil {
			return nil, fmt.Errorf("%s: %w", file, err)
		}
		entries = append(entries, &entry)
	}
	return entries, nil
}

func decodeBundle(b []byte) ([]*model.Entry, error) {
	var shape map[string]any
	if err := yaml.Unmarshal(b, &shape); err != nil {
		return nil, err
	}
	if _, ok := shape["entries"]; ok {
		var bundle Bundle
		if err := decodeStrict(b, &bundle); err != nil {
			return nil, err
		}
		if bundle.Version != 1 {
			return nil, fmt.Errorf("unsupported proposal version %d", bundle.Version)
		}
		return bundle.Entries, nil
	}
	var entry model.Entry
	if err := decodeStrict(b, &entry); err != nil {
		return nil, err
	}
	return []*model.Entry{&entry}, nil
}

func decodeStrictFile(path string, out any) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return decodeStrict(b, out)
}

func decodeStrict(b []byte, out any) error {
	dec := yaml.NewDecoder(bytes.NewReader(b))
	dec.KnownFields(true)
	return dec.Decode(out)
}

func writeYAMLExclusive(path string, value any) error {
	b, err := yaml.Marshal(value)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.Write(b)
	return err
}

func validID(id string) bool {
	if len(id) != 21 || id[0] != 'p' {
		return false
	}
	for _, r := range id[1:] {
		if !strings.ContainsRune("0123456789abcdef", r) {
			return false
		}
	}
	return true
}
