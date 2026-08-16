package adapters

import (
	"os"
	"path/filepath"
	"testing"
)

func write(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// Fixture lines mirror the observed on-disk formats (pinned 2026-08).
const claudeFixture = `{"type":"mode","mode":"normal","sessionId":"s1"}
{"type":"user","message":{"role":"user","content":"let's use SQLite for state"},"uuid":"u1","timestamp":"2026-08-11T14:02:11Z","cwd":"/w","sessionId":"s1","version":"2.1.220"}
{"type":"assistant","message":{"role":"assistant","content":[{"type":"text","text":"Agreed, SQLite it is."},{"type":"tool_use","id":"t1","name":"Write","input":{"file_path":"/w/store/db.go","content":"x"}},{"type":"tool_use","id":"t2","name":"Bash","input":{"command":"go test ./...","description":"run tests"}}]},"timestamp":"2026-08-11T14:03:00Z","cwd":"/w","sessionId":"s1"}
{"type":"user","message":{"role":"user","content":[{"type":"tool_result","tool_use_id":"t2","content":[{"type":"text","text":"ok clew/internal/journal 0.3s"}]}]},"timestamp":"2026-08-11T14:03:05Z","cwd":"/w","sessionId":"s1"}
{"type":"brand-new-thing","payload":{"x":1}}
`

func TestClaudeAdapterPinned(t *testing.T) {
	dir := t.TempDir()
	f := write(t, dir, "s1.jsonl", claudeFixture)
	d, err := (&Claude{}).Parse(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	if d.SessionID != "s1" || d.CWD != "/w" {
		t.Errorf("meta: %q %q", d.SessionID, d.CWD)
	}
	if len(d.Messages) != 4 { // user, assistant text, bash cmd, tool_result
		t.Fatalf("messages: got %d want 4: %+v", len(d.Messages), d.Messages)
	}
	if d.Messages[0].Role != "user" || d.Messages[0].Line != 2 {
		t.Errorf("first message: %+v", d.Messages[0])
	}
	if d.Messages[3].Role != "tool_result" {
		t.Errorf("tool result role: %+v", d.Messages[3])
	}
	if len(d.Footprints) != 1 || d.Footprints[0] != "/w/store/db.go" {
		t.Errorf("footprints: %v", d.Footprints)
	}
	if len(d.Commands) != 1 {
		t.Errorf("commands: %v", d.Commands)
	}
	// Adapter law: unknown line class parked and counted, not guessed.
	if d.Unknown["brand-new-thing"] != 1 {
		t.Errorf("unknown classes: %v", d.Unknown)
	}
	// Incremental: second parse from watermark yields nothing.
	d2, err := (&Claude{}).Parse(f, d.NewOffset)
	if err != nil {
		t.Fatal(err)
	}
	if len(d2.Messages) != 0 || d2.NewOffset != d.NewOffset {
		t.Errorf("watermark: %+v", d2)
	}
}

const codexFixture = `{"timestamp":"2026-08-11T18:47:09.146Z","type":"session_meta","payload":{"id":"c-1","cwd":"/w","originator":"codex-tui","cli_version":"0.128.0"}}
{"timestamp":"2026-08-11T18:47:10.000Z","type":"response_item","payload":{"type":"message","role":"user","content":[{"type":"input_text","text":"measure the p95 please"}]}}
{"timestamp":"2026-08-11T18:47:11.000Z","type":"response_item","payload":{"type":"reasoning","summary":[],"encrypted_content":"xxx"}}
{"timestamp":"2026-08-11T18:47:12.000Z","type":"response_item","payload":{"type":"function_call","name":"exec_command","arguments":"{\"cmd\":\"hey -n 200 http://localhost:8080\",\"workdir\":\"/w\"}","call_id":"call1"}}
{"timestamp":"2026-08-11T18:47:13.000Z","type":"response_item","payload":{"type":"function_call_output","call_id":"call1","output":"p95: 340ms"}}
{"timestamp":"2026-08-11T18:47:14.000Z","type":"response_item","payload":{"type":"custom_tool_call","status":"completed","call_id":"call2","name":"apply_patch","input":"*** Begin Patch\n*** Update File: /w/server/handler.go\n@@\n+x\n*** End Patch"}}
{"timestamp":"2026-08-11T18:47:15.000Z","type":"response_item","payload":{"type":"message","role":"assistant","content":[{"type":"output_text","text":"p95 is 340ms on the emulator."}]}}
{"timestamp":"2026-08-11T18:47:16.000Z","type":"event_msg","payload":{"type":"token_count"}}
{"timestamp":"2026-08-11T18:47:17.000Z","type":"mystery_type","payload":{}}
`

func TestCodexAdapterPinned(t *testing.T) {
	dir := t.TempDir()
	f := write(t, dir, "rollout-1.jsonl", codexFixture)
	d, err := (&Codex{}).Parse(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	if d.SessionID != "c-1" || d.CWD != "/w" {
		t.Errorf("meta: %+v", d)
	}
	if len(d.Messages) != 4 { // user, $cmd, tool_result, assistant
		t.Fatalf("messages: got %d: %+v", len(d.Messages), d.Messages)
	}
	if d.Messages[2].Role != "tool_result" {
		t.Errorf("output must be tool_result: %+v", d.Messages[2])
	}
	if len(d.Footprints) != 1 || d.Footprints[0] != "/w/server/handler.go" {
		t.Errorf("apply_patch footprint: %v", d.Footprints)
	}
	if d.Unknown["mystery_type"] != 1 {
		t.Errorf("unknown: %v", d.Unknown)
	}
}

const cursorFixture = `{"role":"user","message":{"content":[{"type":"text","text":"<timestamp>x</timestamp>\n<user_query>\nplease add caching\n</user_query>"}]}}
{"role":"assistant","message":{"content":[{"type":"text","text":"Adding a cache layer."},{"type":"tool_use","name":"Write","input":{"path":"/w/cache/lru.go","contents":"x"}},{"type":"tool_use","name":"Shell","input":{"command":"go build ./..."}}]}}
{"role":"assistant","message":{"content":[{"type":"weird_block","x":1}]}}
`

func TestCursorAdapterPinned(t *testing.T) {
	dir := t.TempDir()
	f := write(t, dir, "abc.jsonl", cursorFixture)
	d, err := (&Cursor{}).Parse(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(d.Messages) != 3 {
		t.Fatalf("messages: %+v", d.Messages)
	}
	if d.Messages[0].Text != "please add caching" {
		t.Errorf("user_query unwrap failed: %q", d.Messages[0].Text)
	}
	if len(d.Footprints) != 1 || d.Footprints[0] != "/w/cache/lru.go" {
		t.Errorf("footprints: %v", d.Footprints)
	}
	if d.Unknown["block:weird_block"] != 1 {
		t.Errorf("unknown block: %v", d.Unknown)
	}
}

func TestFormatBreakPausesAdapter(t *testing.T) {
	dir := t.TempDir()
	f := write(t, dir, "s.jsonl", "garbage\nmore garbage\nnot json\nnope\nstill no\n")
	_, err := (&Claude{}).Parse(f, 0)
	if _, ok := err.(*FormatError); !ok {
		t.Fatalf("want FormatError (I2: pause loudly), got %v", err)
	}
}

func TestSlugify(t *testing.T) {
	if got := slugify("/Users/x/my-app"); got != "-Users-x-my-app" {
		t.Errorf("claude slug: %q", got)
	}
}

func TestWrapRoundTrip(t *testing.T) {
	dir := t.TempDir()
	f := write(t, dir, "w.jsonl",
		`{"kind":"meta","argv":["gemini","chat"],"cwd":"/w","session":"w1","at":"2026-08-11T10:00:00Z"}
{"at":"2026-08-11T10:00:05Z","dir":"in","text":"why is the build red?"}
{"at":"2026-08-11T10:00:09Z","dir":"out","text":"the linker needs -lz on darwin"}
`)
	d, err := (&Wrap{}).Parse(f, 0)
	if err != nil {
		t.Fatal(err)
	}
	if d.Agent != "wrap:gemini" || d.CWD != "/w" || d.SessionID != "w1" {
		t.Errorf("meta: %+v", d)
	}
	if len(d.Messages) != 2 || d.Messages[0].Role != "user" || d.Messages[1].Role != "assistant" {
		t.Errorf("messages: %+v", d.Messages)
	}
}
