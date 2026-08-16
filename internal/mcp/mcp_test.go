package mcp

import (
	"bufio"
	"encoding/json"
	"strings"
	"testing"
	"time"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
)

func TestServerRoundTrip(t *testing.T) {
	j, _ := journal.Load(t.TempDir())
	e := &model.Entry{
		ID: ids.NewEntry(time.Now()), Type: model.Decision,
		Title: "push over polling", Body: "battery", Quote: "push over polling because battery",
		UtteranceBy: model.ByUser,
		Source:      model.Source{Kind: model.SrcSession, Ref: "s#L1", At: time.Now()},
		Confidence:  0.9, Tags: []string{"sync/**"},
	}
	if err := j.AddEntry(e); err != nil {
		t.Fatal(err)
	}
	wrote := false
	srv := &Server{Journal: j, Surface: "test", AfterWrite: func() { wrote = true }}

	in := strings.Join([]string{
		`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"testclient"}}}`,
		`{"jsonrpc":"2.0","method":"notifications/initialized"}`,
		`{"jsonrpc":"2.0","id":2,"method":"tools/list"}`,
		`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"journal_search","arguments":{"query":"polling"}}}`,
		`{"jsonrpc":"2.0","id":4,"method":"tools/call","params":{"name":"journal_note","arguments":{"type":"finding","title":"emulator p95","quote":"p95 = 340ms on emulator"}}}`,
		`{"jsonrpc":"2.0","id":5,"method":"nope"}`,
	}, "\n") + "\n"
	var out strings.Builder
	if err := srv.Serve(strings.NewReader(in), &out); err != nil {
		t.Fatal(err)
	}
	sc := bufio.NewScanner(strings.NewReader(out.String()))
	var resps []map[string]any
	for sc.Scan() {
		var m map[string]any
		if err := json.Unmarshal(sc.Bytes(), &m); err != nil {
			t.Fatal(err)
		}
		resps = append(resps, m)
	}
	if len(resps) != 5 { // notification gets no response
		t.Fatalf("want 5 responses, got %d", len(resps))
	}
	if r := resps[1]["result"].(map[string]any); len(r["tools"].([]any)) != 3 {
		t.Fatal("tools/list must expose exactly journal_search, journal_get, journal_note")
	}
	search := resps[2]["result"].(map[string]any)["content"].([]any)[0].(map[string]any)["text"].(string)
	if !strings.Contains(search, "push over polling") {
		t.Errorf("search miss: %q", search)
	}
	note := resps[3]["result"].(map[string]any)["content"].([]any)[0].(map[string]any)["text"].(string)
	if !strings.Contains(note, "recorded e") || !wrote {
		t.Errorf("journal_note failed: %q (AfterWrite=%v)", note, wrote)
	}
	if resps[4]["error"] == nil {
		t.Error("unknown method must return a JSON-RPC error")
	}
	// The note is a real journal entry with mcp provenance.
	found := false
	for _, en := range j.Entries {
		if en.Title == "emulator p95" && en.Source.Ref == "mcp:testclient" {
			found = true
		}
	}
	if !found {
		t.Error("journal_note entry not persisted with mcp provenance")
	}
}
