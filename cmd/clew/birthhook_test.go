package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func readJSONMap(t *testing.T, path string) map[string]any {
	t.Helper()
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

func countBirthHooks(settings map[string]any, executable string) int {
	hooks, _ := settings["hooks"].(map[string]any)
	groups, _ := hooks["SessionStart"].([]any)
	count := 0
	for _, rawGroup := range groups {
		group, _ := rawGroup.(map[string]any)
		handlers, _ := group["hooks"].([]any)
		for _, handler := range handlers {
			if isBirthHandler(handler, executable) {
				count++
			}
		}
	}
	return count
}

func TestClaudeBirthHookInstallPreservesSettingsAndRemoveIsSurgical(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	original := map[string]any{
		"env": map[string]any{"KEEP": "yes"},
		"hooks": map[string]any{
			"PreToolUse": []any{map[string]any{
				"matcher": "Bash",
				"hooks":   []any{map[string]any{"type": "command", "command": "/usr/bin/check"}},
			}},
			"SessionStart": []any{map[string]any{
				"matcher": "resume",
				"hooks":   []any{map[string]any{"type": "command", "command": "/usr/bin/resume"}},
			}},
		},
	}
	b, _ := json.MarshalIndent(original, "", "  ")
	if err := os.WriteFile(path, append(b, '\n'), 0o600); err != nil {
		t.Fatal(err)
	}

	executable := "/Applications/Clew Tools/clew"
	if err := installClaudeBirthHook(executable); err != nil {
		t.Fatal(err)
	}
	installed := readJSONMap(t, path)
	if got := countBirthHooks(installed, executable); got != 1 {
		t.Fatalf("installed birth hook count = %d, want 1", got)
	}
	var installedTimeout float64
	for _, rawGroup := range installed["hooks"].(map[string]any)["SessionStart"].([]any) {
		group, _ := rawGroup.(map[string]any)
		handlers, _ := group["hooks"].([]any)
		for _, rawHandler := range handlers {
			handler, _ := rawHandler.(map[string]any)
			if isBirthHandler(handler, executable) {
				installedTimeout, _ = handler["timeout"].(float64)
			}
		}
	}
	if installedTimeout != 30 {
		t.Fatalf("installed birth hook timeout = %v, want 30 seconds", installedTimeout)
	}
	if !reflect.DeepEqual(installed["env"], original["env"]) {
		t.Fatalf("unrelated env changed: %#v", installed["env"])
	}
	hooks := installed["hooks"].(map[string]any)
	if !reflect.DeepEqual(hooks["PreToolUse"], original["hooks"].(map[string]any)["PreToolUse"]) {
		t.Fatal("unrelated PreToolUse hook changed")
	}
	fi, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if fi.Mode().Perm() != 0o600 {
		t.Fatalf("settings mode = %v, want 0600", fi.Mode().Perm())
	}

	// Reinstall is idempotent.
	if err := installClaudeBirthHook(executable); err != nil {
		t.Fatal(err)
	}
	if got := countBirthHooks(readJSONMap(t, path), executable); got != 1 {
		t.Fatalf("birth hook count after reinstall = %d, want 1", got)
	}

	if err := removeClaudeBirthHook(executable); err != nil {
		t.Fatal(err)
	}
	removed := readJSONMap(t, path)
	if got := countBirthHooks(removed, executable); got != 0 {
		t.Fatalf("birth hook count after removal = %d, want 0", got)
	}
	if !reflect.DeepEqual(removed, original) {
		got, _ := json.MarshalIndent(removed, "", "  ")
		want, _ := json.MarshalIndent(original, "", "  ")
		t.Fatalf("remove did not restore unrelated settings\ngot:  %s\nwant: %s", got, want)
	}
	if matches, _ := filepath.Glob(filepath.Join(filepath.Dir(path), ".settings.json.clew-*")); len(matches) != 0 {
		t.Fatalf("atomic-write temp files leaked: %v", matches)
	}
}

func TestClaudeBirthHookInstallUpdatesOlderExecutablePath(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	old := `{"hooks":{"SessionStart":[{"matcher":"startup","hooks":[{"type":"command","command":"/old/bin/clew","args":["_birth"]}]}]}}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	current := "/new/bin/clew"
	if err := installClaudeBirthHook(current); err != nil {
		t.Fatal(err)
	}
	settings := readJSONMap(t, path)
	if got := countBirthHooks(settings, current); got != 1 {
		t.Fatalf("updated birth hook count = %d, want 1", got)
	}
	b, _ := os.ReadFile(path)
	if strings.Contains(string(b), "/old/bin/clew") {
		t.Fatalf("old executable path survived update: %s", b)
	}
}

func TestClaudeBirthHookInstallRepairsUnboundedOrAsyncHandler(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	executable := "/usr/local/bin/clew"
	old := `{"hooks":{"SessionStart":[{"matcher":"startup","hooks":[{"type":"command","command":"/usr/local/bin/clew","args":["_birth"],"async":true}]}]}}`
	if err := os.WriteFile(path, []byte(old), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := installClaudeBirthHook(executable); err != nil {
		t.Fatal(err)
	}
	settings := readJSONMap(t, path)
	hooks := settings["hooks"].(map[string]any)["SessionStart"].([]any)
	if len(hooks) != 1 {
		t.Fatalf("SessionStart group count = %d, want 1", len(hooks))
	}
	handlers := hooks[0].(map[string]any)["hooks"].([]any)
	if len(handlers) != 1 {
		t.Fatalf("SessionStart handler count = %d, want 1", len(handlers))
	}
	handler := handlers[0].(map[string]any)
	if got, _ := handler["timeout"].(float64); got != claudeBirthHookTimeout {
		t.Fatalf("repaired birth hook timeout = %v, want %d", got, claudeBirthHookTimeout)
	}
	if _, exists := handler["async"]; exists {
		t.Fatalf("repaired birth hook remained asynchronous: %#v", handler)
	}
}

func TestBirthHandlerEquivalenceRequiresExactStartupExecShape(t *testing.T) {
	executable := "/current/bin/clew"
	validHandler := func() map[string]any {
		return map[string]any{
			"type": "command", "command": executable,
			"args": []any{"_birth"}, "timeout": float64(30),
		}
	}
	tests := []struct {
		name    string
		matcher string
		mutate  func(map[string]any)
	}{
		{name: "wrong matcher", matcher: "resume"},
		{name: "old executable", matcher: "startup", mutate: func(h map[string]any) { h["command"] = "/old/bin/clew" }},
		{name: "wrong args", matcher: "startup", mutate: func(h map[string]any) { h["args"] = []any{"_birth", "extra"} }},
		{name: "missing timeout", matcher: "startup", mutate: func(h map[string]any) { delete(h, "timeout") }},
		{name: "wrong timeout", matcher: "startup", mutate: func(h map[string]any) { h["timeout"] = float64(29) }},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := validHandler()
			if test.mutate != nil {
				test.mutate(handler)
			}
			groups := []any{map[string]any{
				"matcher": test.matcher,
				"hooks":   []any{handler},
			}}
			if exact, _ := birthHandlerCounts(groups, executable); exact != 0 {
				t.Fatalf("non-equivalent handler counted as exact: %#v", groups)
			}
		})
	}

	groups := []any{map[string]any{
		"matcher": "startup",
		"hooks":   []any{validHandler()},
	}}
	if exact, managed := birthHandlerCounts(groups, executable); exact != 1 || managed != 1 {
		t.Fatalf("valid handler counts = exact %d, managed %d; want 1, 1", exact, managed)
	}
}

func TestClaudeBirthHookInstallCollapsesManagedDuplicates(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	duplicate := `{"hooks":{"SessionStart":[` +
		`{"matcher":"startup","hooks":[{"type":"command","command":"/new/bin/clew","args":["_birth"]}]},` +
		`{"matcher":"startup","hooks":[{"type":"command","command":"/old/bin/clew","args":["_birth"]}]}` +
		`]}}`
	if err := os.WriteFile(path, []byte(duplicate), 0o644); err != nil {
		t.Fatal(err)
	}
	executable := "/new/bin/clew"
	if err := installClaudeBirthHook(executable); err != nil {
		t.Fatal(err)
	}
	if got := countBirthHooks(readJSONMap(t, path), executable); got != 1 {
		t.Fatalf("managed hook count after repair = %d, want 1", got)
	}
}

func TestClaudeBirthHookInvalidSettingsAreUntouched(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	path := filepath.Join(home, ".claude", "settings.json")
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	original := []byte("{not-json}\n")
	if err := os.WriteFile(path, original, 0o644); err != nil {
		t.Fatal(err)
	}
	if err := installClaudeBirthHook("/usr/local/bin/clew"); err == nil {
		t.Fatal("invalid settings were silently overwritten")
	}
	got, _ := os.ReadFile(path)
	if string(got) != string(original) {
		t.Fatalf("invalid settings changed: %q", got)
	}
}

func TestRemoveClaudeBirthHookMissingSettingsIsNoOp(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := removeClaudeBirthHook("/usr/local/bin/clew"); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, ".claude", "settings.json")
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Fatalf("remove created missing settings: %v", err)
	}
}

func TestClaudeBirthHookHonorsClaudeConfigDir(t *testing.T) {
	home := t.TempDir()
	custom := filepath.Join(t.TempDir(), "claude-custom")
	t.Setenv("HOME", home)
	t.Setenv("CLAUDE_CONFIG_DIR", custom)
	if err := installClaudeBirthHook("/usr/local/bin/clew"); err != nil {
		t.Fatal(err)
	}
	if got := countBirthHooks(readJSONMap(t, filepath.Join(custom, "settings.json")), "/usr/local/bin/clew"); got != 1 {
		t.Fatalf("custom-dir birth hook count = %d, want 1", got)
	}
	if _, err := os.Stat(filepath.Join(home, ".claude", "settings.json")); !os.IsNotExist(err) {
		t.Fatalf("default Claude settings were touched: %v", err)
	}
}

func TestParseBirthHookCWD(t *testing.T) {
	got, err := parseBirthHookCWD(strings.NewReader(`{
  "hook_event_name":"SessionStart",
  "source":"startup",
  "cwd":"/Users/test/new repo"
}`))
	if err != nil {
		t.Fatal(err)
	}
	if got != "/Users/test/new repo" {
		t.Fatalf("cwd = %q", got)
	}

	for name, input := range map[string]string{
		"malformed":    `{`,
		"missing cwd":  `{"hook_event_name":"SessionStart","source":"startup"}`,
		"wrong event":  `{"hook_event_name":"PreToolUse","source":"startup","cwd":"/w"}`,
		"wrong source": `{"hook_event_name":"SessionStart","source":"resume","cwd":"/w"}`,
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := parseBirthHookCWD(strings.NewReader(input)); err == nil {
				t.Fatalf("parse accepted %s input", name)
			}
		})
	}
}
