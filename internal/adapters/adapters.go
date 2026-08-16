// Package adapters implements the v1 session sensors (JOURNAL_SPEC §5.1):
// claude, codex, cursor (CLI transcripts), wrap. Adapter law (I2): parsers
// are pinned to observed format versions; unknown line classes are parked
// raw and counted; a format break pauses the adapter with a red status line.
// No heuristic parsing of unknown formats, ever.
package adapters

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Message is one utterance in a transcript slice.
type Message struct {
	Role string // user | assistant | tool_result
	Text string
	At   time.Time
	Line int // 1-based line number in the session file (source spans, I7)
}

// Delta is what a tail pass yields for one session file.
type Delta struct {
	Adapter    string
	File       string
	SessionID  string
	CWD        string
	Agent      string // adapter id per §3.2 source.agent
	Title      string
	Messages   []Message
	Footprints []string // files touched via attributed tool calls
	Commands   []string // shell commands run in the session
	NewOffset  int64
	StartLine  int            // line number of first new line
	Unknown    map[string]int // unknown line classes, parked+counted (I2)
	Bytes      int            // new bytes consumed (budget observation, I9)
}

// FormatError means the adapter must pause loudly (I2), not guess.
type FormatError struct {
	Adapter, File, Detail string
}

func (e *FormatError) Error() string {
	return fmt.Sprintf("%s adapter: format break in %s: %s", e.Adapter, e.File, e.Detail)
}

// readNew returns complete new lines since offset (a trailing partial line is
// left for the next pass), the new offset, and the 1-based first line number.
func readNew(file string, offset int64) ([][]byte, int64, int, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, offset, 0, err
	}
	defer f.Close()
	fi, err := f.Stat()
	if err != nil {
		return nil, offset, 0, err
	}
	if fi.Size() < offset { // truncated/rotated: start over (watermark reset)
		offset = 0
	}
	if fi.Size() == offset {
		return nil, offset, 0, nil
	}
	// Count lines before offset for line numbering.
	startLine := 1
	if offset > 0 {
		head := make([]byte, offset)
		if _, err := f.ReadAt(head, 0); err != nil {
			return nil, offset, 0, err
		}
		startLine = bytes.Count(head, []byte{'\n'}) + 1
	}
	buf := make([]byte, fi.Size()-offset)
	if _, err := f.ReadAt(buf, offset); err != nil {
		return nil, offset, 0, err
	}
	last := bytes.LastIndexByte(buf, '\n')
	if last < 0 {
		return nil, offset, 0, nil // no complete line yet
	}
	buf = buf[:last+1]
	var lines [][]byte
	for _, l := range bytes.Split(buf, []byte{'\n'}) {
		if len(bytes.TrimSpace(l)) > 0 {
			lines = append(lines, l)
		}
	}
	return lines, offset + int64(last+1), startLine, nil
}

// slugify reproduces the session-store cwd slug: every non-alphanumeric
// byte becomes '-' (observed in both claude and cursor stores).
func slugify(path string) string {
	var b strings.Builder
	for _, r := range path {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') {
			b.WriteRune(r)
		} else {
			b.WriteByte('-')
		}
	}
	return b.String()
}

func parseTS(s string) time.Time {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return time.Time{}
	}
	return t
}

func home() string {
	h, _ := os.UserHomeDir()
	return h
}

// An Adapter can discover session files for a repo and parse new content.
type Adapter interface {
	ID() string
	// Discover returns session files whose cwd matches repoPath.
	Discover(repoPath string) []string
	// Parse consumes complete lines from offset and returns a delta.
	Parse(file string, offset int64) (*Delta, error)
}

// All returns the v1 adapter set.
func All() []Adapter {
	return []Adapter{&Claude{}, &Codex{}, &Cursor{}, &Wrap{}}
}

// ---------- claude (~/.claude/projects/<slug>/*.jsonl) ----------
// Pinned to the observed 2.x envelope: {"type":"user|assistant",...,
// "message":{"role","content"}, "timestamp","cwd","sessionId"}.

type Claude struct{}

func (a *Claude) ID() string { return "claude-code" }

// knownClaudeTypes are non-message line classes we deliberately skip.
var knownClaudeTypes = map[string]bool{
	"system": true, "summary": true, "attachment": true, "mode": true,
	"permission-mode": true, "ai-title": true, "last-prompt": true,
	"queue-operation": true, "file-history-snapshot": true,
	"file-history-delta": true, "progress": true, "thinking-block": true,
	"compact-boundary": true, "diagnostic": true,
}

func (a *Claude) Discover(repoPath string) []string {
	dir := filepath.Join(home(), ".claude", "projects", slugify(repoPath))
	m, _ := filepath.Glob(filepath.Join(dir, "*.jsonl"))
	return m
}

func (a *Claude) Parse(file string, offset int64) (*Delta, error) {
	lines, newOff, startLine, err := readNew(file, offset)
	if err != nil {
		return nil, err
	}
	d := &Delta{Adapter: a.ID(), Agent: a.ID(), File: file, NewOffset: newOff,
		StartLine: startLine, Unknown: map[string]int{}}
	bad := 0
	for i, raw := range lines {
		d.Bytes += len(raw) + 1
		var env struct {
			Type      string          `json:"type"`
			Timestamp string          `json:"timestamp"`
			CWD       string          `json:"cwd"`
			SessionID string          `json:"sessionId"`
			Message   json.RawMessage `json:"message"`
		}
		if json.Unmarshal(raw, &env) != nil || env.Type == "" {
			bad++
			d.Unknown["unparseable"]++
			continue
		}
		if env.CWD != "" {
			d.CWD = env.CWD
		}
		if env.SessionID != "" {
			d.SessionID = env.SessionID
		}
		at := parseTS(env.Timestamp)
		line := startLine + i
		switch env.Type {
		case "user", "assistant":
			a.message(d, env.Type, env.Message, at, line)
		default:
			if !knownClaudeTypes[env.Type] {
				d.Unknown[env.Type]++
			}
		}
	}
	if len(lines) > 4 && bad*2 > len(lines) {
		return d, &FormatError{a.ID(), file, fmt.Sprintf("%d/%d lines unparseable (version drift?)", bad, len(lines))}
	}
	return d, nil
}

func (a *Claude) message(d *Delta, role string, raw json.RawMessage, at time.Time, line int) {
	var m struct {
		Role    string          `json:"role"`
		Content json.RawMessage `json:"content"`
	}
	if json.Unmarshal(raw, &m) != nil {
		d.Unknown["message-shape"]++
		return
	}
	// content: plain string
	var s string
	if json.Unmarshal(m.Content, &s) == nil {
		if t := strings.TrimSpace(s); t != "" {
			d.Messages = append(d.Messages, Message{Role: role, Text: t, At: at, Line: line})
		}
		return
	}
	// content: block array
	var blocks []map[string]any
	if json.Unmarshal(m.Content, &blocks) != nil {
		d.Unknown["content-shape"]++
		return
	}
	for _, b := range blocks {
		switch b["type"] {
		case "text":
			if t, _ := b["text"].(string); strings.TrimSpace(t) != "" {
				d.Messages = append(d.Messages, Message{Role: role, Text: t, At: at, Line: line})
			}
		case "tool_use":
			name, _ := b["name"].(string)
			input, _ := b["input"].(map[string]any)
			switch name {
			case "Edit", "Write", "MultiEdit", "NotebookEdit":
				if p, _ := input["file_path"].(string); p != "" {
					d.Footprints = append(d.Footprints, p)
				}
			case "Bash":
				if c, _ := input["command"].(string); c != "" {
					d.Commands = append(d.Commands, c)
					d.Messages = append(d.Messages, Message{Role: "assistant", Text: "$ " + c, At: at, Line: line})
				}
			}
		case "tool_result":
			txt := flattenToolResult(b["content"])
			if txt != "" {
				d.Messages = append(d.Messages, Message{Role: "tool_result", Text: clip(txt, 2000), At: at, Line: line})
			}
		case "thinking", "image", "document":
		default:
			if t, _ := b["type"].(string); t != "" {
				d.Unknown["block:"+t]++
			}
		}
	}
}

func flattenToolResult(c any) string {
	switch v := c.(type) {
	case string:
		return v
	case []any:
		var parts []string
		for _, it := range v {
			if m, ok := it.(map[string]any); ok {
				if m["type"] == "text" {
					if t, _ := m["text"].(string); t != "" {
						parts = append(parts, t)
					}
				}
			}
		}
		return strings.Join(parts, "\n")
	}
	return ""
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + " …[clipped]"
}

// ---------- codex (~/.codex/sessions/YYYY/MM/DD/rollout-*.jsonl) ----------
// Pinned to the observed envelope: {"timestamp","type","payload"} with types
// session_meta | response_item | event_msg | turn_context.

type Codex struct{}

func (a *Codex) ID() string { return "codex" }

func (a *Codex) Discover(repoPath string) []string {
	root := filepath.Join(home(), ".codex", "sessions")
	var out []string
	// Only recent files can be cheaply cwd-matched: meta line holds the cwd,
	// checked in Parse; Discover returns candidates from the last 3 days.
	for d := 0; d < 3; d++ {
		day := time.Now().AddDate(0, 0, -d)
		m, _ := filepath.Glob(filepath.Join(root, day.Format("2006"), day.Format("01"), day.Format("02"), "rollout-*.jsonl"))
		out = append(out, m...)
	}
	var match []string
	for _, f := range out {
		if cwd := codexCWD(f); cwd != "" && sameOrUnder(cwd, repoPath) {
			match = append(match, f)
		}
	}
	return match
}

func sameOrUnder(cwd, repo string) bool {
	return cwd == repo || strings.HasPrefix(cwd, repo+string(os.PathSeparator))
}

// CodexCWD peeks at the session_meta line (first line) for the cwd.
// Exported for backfill's dated-directory discovery.
func CodexCWD(file string) string { return codexCWD(file) }

// codexCWD peeks at the session_meta line (first line) for the cwd.
func codexCWD(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 8192)
	n, _ := f.Read(buf)
	i := bytes.IndexByte(buf[:n], '\n')
	if i < 0 {
		i = n
	}
	var env struct {
		Type    string `json:"type"`
		Payload struct {
			CWD string `json:"cwd"`
		} `json:"payload"`
	}
	if json.Unmarshal(buf[:i], &env) != nil || env.Type != "session_meta" {
		return ""
	}
	return env.Payload.CWD
}

var knownCodexTypes = map[string]bool{
	"event_msg": true, "turn_context": true, "compacted": true,
}

var knownCodexItems = map[string]bool{
	"reasoning": true, "web_search_call": true, "function_call_output": false,
}

func (a *Codex) Parse(file string, offset int64) (*Delta, error) {
	lines, newOff, startLine, err := readNew(file, offset)
	if err != nil {
		return nil, err
	}
	d := &Delta{Adapter: a.ID(), Agent: a.ID(), File: file, NewOffset: newOff,
		StartLine: startLine, Unknown: map[string]int{}}
	bad := 0
	for i, raw := range lines {
		d.Bytes += len(raw) + 1
		var env struct {
			Timestamp string          `json:"timestamp"`
			Type      string          `json:"type"`
			Payload   json.RawMessage `json:"payload"`
		}
		if json.Unmarshal(raw, &env) != nil || env.Type == "" {
			bad++
			d.Unknown["unparseable"]++
			continue
		}
		at := parseTS(env.Timestamp)
		line := startLine + i
		switch env.Type {
		case "session_meta":
			var p struct {
				ID  string `json:"id"`
				CWD string `json:"cwd"`
			}
			if json.Unmarshal(env.Payload, &p) == nil {
				d.SessionID, d.CWD = p.ID, p.CWD
			}
		case "response_item":
			a.item(d, env.Payload, at, line)
		default:
			if !knownCodexTypes[env.Type] {
				d.Unknown[env.Type]++
			}
		}
	}
	if len(lines) > 4 && bad*2 > len(lines) {
		return d, &FormatError{a.ID(), file, fmt.Sprintf("%d/%d lines unparseable", bad, len(lines))}
	}
	return d, nil
}

func (a *Codex) item(d *Delta, raw json.RawMessage, at time.Time, line int) {
	var p map[string]any
	if json.Unmarshal(raw, &p) != nil {
		d.Unknown["payload-shape"]++
		return
	}
	typ, _ := p["type"].(string)
	switch typ {
	case "message":
		role, _ := p["role"].(string)
		if role != "user" && role != "assistant" {
			return // developer/system instruction blobs are not utterances
		}
		var parts []string
		if content, ok := p["content"].([]any); ok {
			for _, c := range content {
				if m, ok := c.(map[string]any); ok {
					if t, _ := m["text"].(string); t != "" {
						parts = append(parts, t)
					}
				}
			}
		}
		txt := strings.TrimSpace(strings.Join(parts, "\n"))
		// codex user messages wrap environment/system context in tags.
		if role == "user" && (strings.HasPrefix(txt, "<environment_context>") || strings.HasPrefix(txt, "<user_instructions>")) {
			return
		}
		if txt != "" {
			d.Messages = append(d.Messages, Message{Role: role, Text: txt, At: at, Line: line})
		}
	case "function_call":
		name, _ := p["name"].(string)
		argsStr, _ := p["arguments"].(string)
		var args map[string]any
		json.Unmarshal([]byte(argsStr), &args)
		switch name {
		case "exec_command", "shell", "container.exec", "local_shell":
			cmd := ""
			if c, ok := args["cmd"].(string); ok {
				cmd = c
			} else if arr, ok := args["command"].([]any); ok {
				var w []string
				for _, x := range arr {
					if s, ok := x.(string); ok {
						w = append(w, s)
					}
				}
				cmd = strings.Join(w, " ")
			}
			if cmd != "" {
				d.Commands = append(d.Commands, cmd)
				d.Messages = append(d.Messages, Message{Role: "assistant", Text: "$ " + clip(cmd, 500), At: at, Line: line})
			}
		}
	case "custom_tool_call":
		name, _ := p["name"].(string)
		if name == "apply_patch" {
			input, _ := p["input"].(string)
			for _, fp := range patchPaths(input) {
				d.Footprints = append(d.Footprints, fp)
			}
		}
	case "function_call_output", "custom_tool_call_output":
		if out, ok := p["output"].(string); ok && strings.TrimSpace(out) != "" {
			d.Messages = append(d.Messages, Message{Role: "tool_result", Text: clip(out, 2000), At: at, Line: line})
		}
	case "reasoning", "web_search_call":
	default:
		if typ != "" {
			d.Unknown["item:"+typ]++
		}
	}
}

// patchPaths extracts file paths from an apply_patch envelope.
func patchPaths(patch string) []string {
	var out []string
	for _, l := range strings.Split(patch, "\n") {
		for _, pfx := range []string{"*** Update File: ", "*** Add File: ", "*** Delete File: "} {
			if p, ok := strings.CutPrefix(l, pfx); ok {
				out = append(out, strings.TrimSpace(p))
			}
		}
	}
	return out
}

// ---------- cursor (~/.cursor/projects/<slug>/agent-transcripts/<id>/<id>.jsonl) ----------
// Pinned to the observed CLI transcript: {"role","message":{"content":[...]}}
// Anthropic-style blocks. No per-line timestamps (file mtime is used); the
// desktop composer store (state.vscdb) is detected but NOT parsed in v1 —
// flagged as lower fidelity in status (§5.1, open decision 3).

type Cursor struct{}

func (a *Cursor) ID() string { return "cursor-agent" }

func (a *Cursor) Discover(repoPath string) []string {
	slug := strings.TrimPrefix(slugify(repoPath), "-")
	dir := filepath.Join(home(), ".cursor", "projects", slug, "agent-transcripts")
	m, _ := filepath.Glob(filepath.Join(dir, "*", "*.jsonl"))
	return m
}

// DesktopStorePresent reports the unparsed desktop store (status line, I2).
func DesktopStorePresent() bool {
	p := filepath.Join(home(), "Library", "Application Support", "Cursor", "User", "globalStorage", "state.vscdb")
	_, err := os.Stat(p)
	return err == nil
}

func (a *Cursor) Parse(file string, offset int64) (*Delta, error) {
	lines, newOff, startLine, err := readNew(file, offset)
	if err != nil {
		return nil, err
	}
	fi, _ := os.Stat(file)
	at := time.Time{}
	if fi != nil {
		at = fi.ModTime().UTC()
	}
	d := &Delta{Adapter: a.ID(), Agent: a.ID(), File: file, NewOffset: newOff,
		StartLine: startLine, Unknown: map[string]int{}}
	d.SessionID = strings.TrimSuffix(filepath.Base(file), ".jsonl")
	bad := 0
	for i, raw := range lines {
		d.Bytes += len(raw) + 1
		var env struct {
			Role    string `json:"role"`
			Message struct {
				Content []map[string]any `json:"content"`
			} `json:"message"`
		}
		if json.Unmarshal(raw, &env) != nil || env.Role == "" {
			bad++
			d.Unknown["unparseable"]++
			continue
		}
		line := startLine + i
		for _, b := range env.Message.Content {
			switch b["type"] {
			case "text":
				if t, _ := b["text"].(string); strings.TrimSpace(t) != "" {
					d.Messages = append(d.Messages, Message{Role: env.Role, Text: stripCursorWrapper(t), At: at, Line: line})
				}
			case "tool_use":
				name, _ := b["name"].(string)
				input, _ := b["input"].(map[string]any)
				for _, k := range []string{"path", "file_path", "target_notebook"} {
					if p, _ := input[k].(string); p != "" {
						d.Footprints = append(d.Footprints, p)
					}
				}
				if name == "Shell" {
					if c, _ := input["command"].(string); c != "" {
						d.Commands = append(d.Commands, c)
						d.Messages = append(d.Messages, Message{Role: "assistant", Text: "$ " + clip(c, 500), At: at, Line: line})
					}
				}
			case "tool_result":
				if txt := flattenToolResult(b["content"]); txt != "" {
					d.Messages = append(d.Messages, Message{Role: "tool_result", Text: clip(txt, 2000), At: at, Line: line})
				}
			case "thinking", "image":
			default:
				if t, _ := b["type"].(string); t != "" {
					d.Unknown["block:"+t]++
				}
			}
		}
	}
	if len(lines) > 4 && bad*2 > len(lines) {
		return d, &FormatError{a.ID(), file, fmt.Sprintf("%d/%d lines unparseable", bad, len(lines))}
	}
	return d, nil
}

// stripCursorWrapper removes the <user_query> envelope cursor wraps around
// typed user text.
func stripCursorWrapper(s string) string {
	if i := strings.Index(s, "<user_query>"); i >= 0 {
		if j := strings.Index(s, "</user_query>"); j > i {
			return strings.TrimSpace(s[i+len("<user_query>") : j])
		}
	}
	return s
}

// ---------- wrap (~/.restart/raw/*.jsonl, written by `restart wrap`) ----------
// Our own format, v1: first line {"kind":"meta",...}; then
// {"at","dir":"in|out","text"}. Terminal input = user; output = assistant
// (the extractor decides taint for quoted third-party content, §6.5).

type Wrap struct{}

func (a *Wrap) ID() string { return "wrap" }

func (a *Wrap) Discover(repoPath string) []string {
	m, _ := filepath.Glob(filepath.Join(home(), ".restart", "raw", "*.jsonl"))
	if rh := os.Getenv("RESTART_HOME"); rh != "" {
		m2, _ := filepath.Glob(filepath.Join(rh, "raw", "*.jsonl"))
		m = append(m, m2...)
	}
	var out []string
	for _, f := range m {
		if cwd := wrapMeta(f); cwd != "" && sameOrUnder(cwd, repoPath) {
			out = append(out, f)
		}
	}
	return out
}

func wrapMeta(file string) string {
	f, err := os.Open(file)
	if err != nil {
		return ""
	}
	defer f.Close()
	buf := make([]byte, 4096)
	n, _ := f.Read(buf)
	i := bytes.IndexByte(buf[:n], '\n')
	if i < 0 {
		i = n
	}
	var meta struct {
		Kind string `json:"kind"`
		CWD  string `json:"cwd"`
	}
	if json.Unmarshal(buf[:i], &meta) != nil || meta.Kind != "meta" {
		return ""
	}
	return meta.CWD
}

func (a *Wrap) Parse(file string, offset int64) (*Delta, error) {
	lines, newOff, startLine, err := readNew(file, offset)
	if err != nil {
		return nil, err
	}
	d := &Delta{Adapter: a.ID(), File: file, NewOffset: newOff,
		StartLine: startLine, Unknown: map[string]int{}}
	for i, raw := range lines {
		d.Bytes += len(raw) + 1
		var env struct {
			Kind, Dir, Text, At, Argv, CWD, Session string
		}
		var generic map[string]any
		if json.Unmarshal(raw, &generic) != nil {
			d.Unknown["unparseable"]++
			continue
		}
		env.Kind, _ = generic["kind"].(string)
		env.Dir, _ = generic["dir"].(string)
		env.Text, _ = generic["text"].(string)
		env.At, _ = generic["at"].(string)
		env.CWD, _ = generic["cwd"].(string)
		env.Session, _ = generic["session"].(string)
		line := startLine + i
		switch {
		case env.Kind == "meta":
			d.CWD = env.CWD
			d.SessionID = env.Session
			if argv, ok := generic["argv"].([]any); ok && len(argv) > 0 {
				if s, ok := argv[0].(string); ok {
					d.Agent = "wrap:" + filepath.Base(s)
				}
			}
		case env.Dir == "in":
			d.Messages = append(d.Messages, Message{Role: "user", Text: env.Text, At: parseTS(env.At), Line: line})
		case env.Dir == "out":
			d.Messages = append(d.Messages, Message{Role: "assistant", Text: env.Text, At: parseTS(env.At), Line: line})
		default:
			d.Unknown["line"]++
		}
	}
	if d.Agent == "" {
		d.Agent = "wrap"
	}
	return d, nil
}
