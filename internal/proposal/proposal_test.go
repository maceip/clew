package proposal

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/model"
)

func proposalEntry(at time.Time, title string) *model.Entry {
	return &model.Entry{
		ID: ids.NewEntry(at), Type: model.Decision, Title: title,
		Quote: "exact " + title, UtteranceBy: model.ByUser, Confidence: .9,
		Source: model.Source{Kind: model.SrcHuman, Ref: "notes.md:12", At: at},
	}
}

func writeBundle(t *testing.T, path string, entries ...*model.Entry) {
	t.Helper()
	b, err := yaml.Marshal(Bundle{Version: 1, Entries: entries})
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, b, 0o600); err != nil {
		t.Fatal(err)
	}
}

func TestStageBundleMarksForeignAndCreatesExactDiff(t *testing.T) {
	root := t.TempDir()
	source := filepath.Join(root, "bundle.yaml")
	entry := proposalEntry(time.Now().UTC(), "push, not poll")
	writeBundle(t, source, entry)
	m := New(filepath.Join(root, "store"))
	batch, err := m.Stage(context.Background(), "/repo", source)
	if err != nil {
		t.Fatal(err)
	}
	if batch.Status != "pending" || len(batch.Entries) != 1 || batch.Entries[0].Source.Kind != model.SrcForeign {
		t.Fatalf("batch = %#v", batch)
	}
	if batch.Entries[0].Quote != entry.Quote || !strings.Contains(batch.Entries[0].Source.Ref, "#human:notes.md:12") {
		t.Fatalf("foreign entry = %#v", batch.Entries[0])
	}
	diff, err := os.ReadFile(filepath.Join(batch.Dir, "diff.md"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(diff), "> "+entry.Quote) || !strings.Contains(string(diff), "entries: 1") {
		t.Fatalf("diff = %s", diff)
	}
	again, err := m.Stage(context.Background(), "/repo", source)
	if err != nil || again.ID != batch.ID {
		t.Fatalf("idempotent stage = %#v, %v", again, err)
	}
	if err := m.SetStatus("/repo", batch.ID, "accepted"); err != nil {
		t.Fatal(err)
	}
	loaded, err := m.Load("/repo", batch.ID)
	if err != nil || loaded.Status != "accepted" {
		t.Fatalf("loaded = %#v, %v", loaded, err)
	}
}

func TestStageRejectsUnknownSchemaAndMissingEvidence(t *testing.T) {
	root := t.TempDir()
	m := New(filepath.Join(root, "store"))
	badSchema := filepath.Join(root, "bad-schema.yaml")
	if err := os.WriteFile(badSchema, []byte("version: 1\nentries: []\nsurprise: true\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := m.Stage(context.Background(), "/repo", badSchema); err == nil {
		t.Fatal("unknown bundle field accepted")
	}
	entry := proposalEntry(time.Now().UTC(), "no evidence")
	entry.Quote = ""
	badEntry := filepath.Join(root, "bad-entry.yaml")
	writeBundle(t, badEntry, entry)
	if _, err := m.Stage(context.Background(), "/repo", badEntry); err == nil || !strings.Contains(err.Error(), "no quote") {
		t.Fatalf("missing evidence error = %v", err)
	}
}

func TestStageDirectoryAndURL(t *testing.T) {
	root := t.TempDir()
	dir := filepath.Join(root, "incoming", "entries")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	entry := proposalEntry(time.Now().UTC(), "from directory")
	b, _ := yaml.Marshal(entry)
	if err := os.WriteFile(filepath.Join(dir, entry.ID+".yaml"), b, 0o600); err != nil {
		t.Fatal(err)
	}
	m := New(filepath.Join(root, "store"))
	if batch, err := m.Stage(context.Background(), "/repo", filepath.Dir(dir)); err != nil || len(batch.Entries) != 1 {
		t.Fatalf("directory = %#v, %v", batch, err)
	}

	bundle, _ := yaml.Marshal(Bundle{Version: 1, Entries: []*model.Entry{proposalEntry(time.Now().UTC(), "from URL")}})
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write(bundle) }))
	defer server.Close()
	if batch, err := m.Stage(context.Background(), "/repo", server.URL+"/bundle.yaml"); err != nil || len(batch.Entries) != 1 {
		t.Fatalf("URL = %#v, %v", batch, err)
	}
}

func TestPushBranchBasesProposalOnJournalAndPublishesForeignEntry(t *testing.T) {
	root := t.TempDir()
	t.Setenv("CLEW_HOME", filepath.Join(root, "clew-home"))
	remote := filepath.Join(root, "remote.git")
	repo := filepath.Join(root, "repo")
	if err := os.MkdirAll(remote, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(repo, 0o755); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(remote, "init", "--bare", "-q"); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "init", "-q"); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.Run(repo, "remote", "add", "origin", remote); err != nil {
		t.Fatal(err)
	}
	if _, err := gitx.EnsureJournal(repo); err != nil {
		t.Fatal(err)
	}
	source := filepath.Join(root, "bundle.yaml")
	entry := proposalEntry(time.Now().UTC(), "branch proposal")
	writeBundle(t, source, entry)
	m := New(filepath.Join(root, "store"))
	batch, err := m.Stage(context.Background(), repo, source)
	if err != nil {
		t.Fatal(err)
	}
	branch := "github.com/maceip/clew/proposal-test"
	if err := m.PushBranch(repo, batch, branch); err != nil {
		t.Fatal(err)
	}
	got, err := gitx.Run(remote, "show", "refs/heads/"+branch+":entries/"+entry.ID+".yaml")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(got, "kind: foreign") || !strings.Contains(got, entry.Quote) {
		t.Fatalf("remote proposal entry = %s", got)
	}
	if err := m.PushBranch(repo, batch, branch); err != nil {
		t.Fatalf("idempotent proposal push: %v", err)
	}
}
