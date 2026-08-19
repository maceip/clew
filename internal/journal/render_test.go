package journal

import (
	"strings"
	"testing"
	"time"

	"github.com/maceip/clew/internal/ids"
	"github.com/maceip/clew/internal/model"
)

func TestRollupStartsWithGitHubDashboardAndHighlightsAbsence(t *testing.T) {
	now := time.Date(2026, 8, 16, 18, 0, 0, 0, time.UTC)
	j, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	add := func(typ model.EntryType, title string, at time.Time) *model.Entry {
		e := &model.Entry{
			ID: ids.NewEntry(at), Type: typ, Title: title, Quote: "exact " + title,
			UtteranceBy: model.ByUser, Source: model.Source{Kind: model.SrcSession, Ref: "session#L1", At: at},
			Confidence: 0.9,
		}
		if err := j.AddEntry(e); err != nil {
			t.Fatal(err)
		}
		return e
	}
	add(model.Decision, "ship push", now.Add(-48*time.Hour))
	add(model.Finding, "battery measurement", now.Add(-3*time.Hour))
	q := add(model.Question, "Owner choice?", now.Add(-8*24*time.Hour))
	q.Asks = "human"
	intent := add(model.Intent, "unbuilt | path", now.Add(-30*24*time.Hour))

	st := Compute(j, now)
	st[intent.ID].Status = StAbsent
	got := Rollup(j, st, now)
	for _, want := range []string{"## DECIDED", "## LEARNED", "## OPEN", "## ALERTS", "## Intent × reality", "★", "**ABSENT**", "unbuilt \\| path"} {
		if !strings.Contains(got, want) {
			t.Errorf("rollup missing %q:\n%s", want, got)
		}
	}
	if strings.Index(got, "## DECIDED") > strings.Index(got, "## Decisions") {
		t.Fatal("dashboard must precede detailed journal")
	}
}

func TestRollupDashboardEmptySectionsRemainVisible(t *testing.T) {
	j, err := Load(t.TempDir())
	if err != nil {
		t.Fatal(err)
	}
	got := Rollup(j, Compute(j, time.Now()), time.Now())
	if strings.Count(got, "_None._") < 4 || !strings.Contains(got, "| _None_ | — | — | — |") {
		t.Fatalf("empty dashboard is not explicit:\n%s", got)
	}
}
