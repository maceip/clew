package manifest

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
)

func TestApplyPersistsCarryDropDispositionsAndIsIdempotent(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	j, err := journal.Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	add := func(typ model.EntryType, title string) *model.Entry {
		e := &model.Entry{
			ID: ids.NewEntry(now), Type: typ, Title: title, Quote: "exact " + title,
			UtteranceBy: model.ByUser, Confidence: .9,
			Source: model.Source{Kind: model.SrcSession, Ref: "session#L1", At: now},
		}
		if err := j.AddEntry(e); err != nil {
			t.Fatal(err)
		}
		return e
	}
	add(model.Decision, "keep this")
	add(model.Finding, "lose this deliberately")
	st := journal.Compute(j, now)
	manifestPath, err := Generate(j, st, "", nil, now)
	if err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(manifestPath)
	if err != nil {
		t.Fatal(err)
	}
	marked := strings.Replace(string(b), "disposition: [carry]", "disposition: [drop]", 1)
	if err := os.WriteFile(manifestPath, []byte(marked), 0o600); err != nil {
		t.Fatal(err)
	}
	out := t.TempDir()
	res, err := Apply(j, st, manifestPath, out, "test", now)
	if err != nil {
		t.Fatal(err)
	}
	if len(res.Carried) != 1 || len(res.Dropped) != 1 {
		t.Fatalf("result = %#v", res)
	}
	for _, test := range []struct{ id, want string }{{res.Carried[0], "carried"}, {res.Dropped[0], "dropped"}} {
		if !j.HasEvent(model.EvDisposition, test.id, "disposition", test.want) {
			t.Errorf("missing %s disposition for %s", test.want, test.id)
		}
	}
	if _, err := Apply(j, st, manifestPath, filepath.Join(out, "again"), "test", now); err != nil {
		t.Fatal(err)
	}
	if len(j.Events) != 2 {
		t.Fatalf("idempotent apply events = %d, want 2", len(j.Events))
	}
}
