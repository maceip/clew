package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var errClaudeSettingsChanged = errors.New("Claude settings changed during clew hook update")

const (
	claudeBirthHookArg     = "_birth"
	claudeBirthHookTimeout = 30
	claudeHookInputLimit   = 1 << 20
)

// installClaudeBirthHook installs the machine-scope Claude SessionStart hook.
// Bootstrap and dispatch are deliberately owned by the caller; this helper
// only performs a preserving settings-file update.
func installClaudeBirthHook(executable string) error {
	path, err := claudeUserSettingsPath()
	if err != nil {
		return err
	}
	_, err = mutateClaudeBirthHook(path, executable, true)
	return err
}

// removeClaudeBirthHook removes only clew's startup hook, leaving every other
// user setting and hook intact. Matching by executable basename also removes a
// hook left by an older installation path of the same binary.
func removeClaudeBirthHook(executable string) error {
	path, err := claudeUserSettingsPath()
	if err != nil {
		return err
	}
	_, err = mutateClaudeBirthHook(path, executable, false)
	return err
}

func claudeUserSettingsPath() (string, error) {
	if dir := os.Getenv("CLAUDE_CONFIG_DIR"); dir != "" {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return "", fmt.Errorf("resolve CLAUDE_CONFIG_DIR: %w", err)
		}
		return filepath.Join(abs, "settings.json"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("resolve home for Claude settings: %w", err)
	}
	if home == "" {
		return "", fmt.Errorf("resolve home for Claude settings: empty home directory")
	}
	return filepath.Join(home, ".claude", "settings.json"), nil
}

func mutateClaudeBirthHook(path, executable string, install bool) (bool, error) {
	for attempt := 0; attempt < 4; attempt++ {
		changed, err := mutateClaudeBirthHookOnce(path, executable, install)
		if !errors.Is(err, errClaudeSettingsChanged) {
			return changed, err
		}
	}
	return false, fmt.Errorf("%w after 4 merge retries", errClaudeSettingsChanged)
}

func mutateClaudeBirthHookOnce(path, executable string, install bool) (bool, error) {
	if executable == "" {
		return false, fmt.Errorf("Claude birth hook executable is empty")
	}

	settings := map[string]any{}
	mode := os.FileMode(0o644)
	b, err := os.ReadFile(path)
	existed := err == nil
	switch {
	case err == nil:
		if len(bytes.TrimSpace(b)) == 0 {
			return false, fmt.Errorf("%s is empty, not touching it", path)
		}
		if err := json.Unmarshal(b, &settings); err != nil {
			return false, fmt.Errorf("%s is not valid JSON, not touching it: %w", path, err)
		}
		if fi, statErr := os.Stat(path); statErr == nil {
			mode = fi.Mode().Perm()
		}
	case os.IsNotExist(err):
		if !install {
			return false, nil
		}
	default:
		return false, err
	}

	hooks, err := objectField(settings, "hooks")
	if err != nil {
		return false, fmt.Errorf("%s: %w", path, err)
	}
	if hooks == nil {
		hooks = map[string]any{}
	}

	groups, err := arrayField(hooks, "SessionStart")
	if err != nil {
		return false, fmt.Errorf("%s: %w", path, err)
	}

	exact, managed := birthHandlerCounts(groups, executable)
	if install && exact == 1 && managed == 1 {
		return false, nil
	}
	filtered, removed := removeBirthHandlers(groups, executable)
	changed := removed
	if install {
		filtered = append(filtered, map[string]any{
			"matcher": "startup",
			"hooks": []any{map[string]any{
				"type":    "command",
				"command": executable,
				"args":    []any{claudeBirthHookArg},
				"timeout": claudeBirthHookTimeout,
			}},
		})
		changed = true
	}

	if !changed {
		return false, nil
	}
	if len(filtered) == 0 {
		delete(hooks, "SessionStart")
	} else {
		hooks["SessionStart"] = filtered
	}
	if len(hooks) == 0 {
		delete(settings, "hooks")
	} else {
		settings["hooks"] = hooks
	}

	out, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return false, err
	}
	out = append(out, '\n')
	if err := atomicWriteClaudeSettings(path, out, mode, b, existed); err != nil {
		return false, err
	}
	return true, nil
}

func objectField(parent map[string]any, key string) (map[string]any, error) {
	v, ok := parent[key]
	if !ok {
		return nil, nil
	}
	obj, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a JSON object", key)
	}
	return obj, nil
}

func arrayField(parent map[string]any, key string) ([]any, error) {
	v, ok := parent[key]
	if !ok {
		return nil, nil
	}
	arr, ok := v.([]any)
	if !ok {
		return nil, fmt.Errorf("%s must be a JSON array", key)
	}
	return arr, nil
}

func removeBirthHandlers(groups []any, executable string) ([]any, bool) {
	out := make([]any, 0, len(groups))
	removed := false
	for _, rawGroup := range groups {
		group, ok := rawGroup.(map[string]any)
		if !ok {
			out = append(out, rawGroup)
			continue
		}
		rawHandlers, ok := group["hooks"].([]any)
		if !ok {
			out = append(out, rawGroup)
			continue
		}
		handlers := make([]any, 0, len(rawHandlers))
		for _, handler := range rawHandlers {
			if isBirthHandler(handler, executable) {
				removed = true
				continue
			}
			handlers = append(handlers, handler)
		}
		if len(handlers) == 0 {
			if len(rawHandlers) > 0 {
				continue
			}
			out = append(out, rawGroup)
			continue
		}
		if len(handlers) != len(rawHandlers) {
			copyGroup := make(map[string]any, len(group))
			for k, v := range group {
				copyGroup[k] = v
			}
			copyGroup["hooks"] = handlers
			out = append(out, copyGroup)
		} else {
			out = append(out, rawGroup)
		}
	}
	return out, removed
}

func birthHandlerCounts(groups []any, executable string) (exact, managed int) {
	for _, rawGroup := range groups {
		group, ok := rawGroup.(map[string]any)
		if !ok {
			continue
		}
		handlers, ok := group["hooks"].([]any)
		if !ok {
			continue
		}
		for _, raw := range handlers {
			if !isBirthHandler(raw, executable) {
				continue
			}
			managed++
			handler := raw.(map[string]any)
			command, _ := handler["command"].(string)
			if group["matcher"] == "startup" && command == executable &&
				birthTimeout(handler["timeout"]) && !birthAsync(handler) {
				exact++
			}
		}
	}
	return exact, managed
}

func birthTimeout(raw any) bool {
	switch value := raw.(type) {
	case int:
		return value == claudeBirthHookTimeout
	case float64: // encoding/json's representation in a map[string]any
		return value == claudeBirthHookTimeout
	default:
		return false
	}
}

func birthAsync(handler map[string]any) bool {
	async, _ := handler["async"].(bool)
	rewake, _ := handler["asyncRewake"].(bool)
	_, conditional := handler["if"]
	return async || rewake || conditional
}

func isBirthHandler(raw any, executable string) bool {
	handler, ok := raw.(map[string]any)
	if !ok || handler["type"] != "command" || !birthArgs(handler["args"]) {
		return false
	}
	command, _ := handler["command"].(string)
	return sameExecutableName(command, executable)
}

func birthArgs(raw any) bool {
	args, ok := raw.([]any)
	return ok && len(args) == 1 && args[0] == claudeBirthHookArg
}

func sameExecutableName(a, b string) bool {
	return a != "" && b != "" && filepath.Base(a) == filepath.Base(b)
}

func atomicWriteClaudeSettings(path string, data []byte, mode os.FileMode, expected []byte, existed bool) (err error) {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	f, err := os.CreateTemp(dir, ".settings.json.clew-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer func() {
		f.Close()
		os.Remove(tmp)
	}()
	if err := f.Chmod(mode.Perm()); err != nil {
		return err
	}
	if _, err := f.Write(data); err != nil {
		return err
	}
	if err := f.Sync(); err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	current, readErr := os.ReadFile(path)
	if existed {
		if readErr != nil || !bytes.Equal(current, expected) {
			return errClaudeSettingsChanged
		}
	} else if readErr == nil || !os.IsNotExist(readErr) {
		return errClaudeSettingsChanged
	}
	return os.Rename(tmp, path)
}

type claudeSessionStartInput struct {
	CWD           string `json:"cwd"`
	HookEventName string `json:"hook_event_name"`
	Source        string `json:"source"`
}

// parseBirthHookCWD parses Claude's SessionStart JSON from stdin. Validation
// is intentionally limited to the hook envelope; git-root and birth policy
// belong to the bootstrap integration.
func parseBirthHookCWD(r io.Reader) (string, error) {
	b, err := io.ReadAll(io.LimitReader(r, claudeHookInputLimit+1))
	if err != nil {
		return "", err
	}
	if len(b) > claudeHookInputLimit {
		return "", fmt.Errorf("Claude SessionStart input exceeds %d bytes", claudeHookInputLimit)
	}
	var in claudeSessionStartInput
	if err := json.Unmarshal(b, &in); err != nil {
		return "", fmt.Errorf("parse Claude SessionStart input: %w", err)
	}
	if in.HookEventName != "" && in.HookEventName != "SessionStart" {
		return "", fmt.Errorf("unexpected Claude hook event %q", in.HookEventName)
	}
	if in.Source != "" && in.Source != "startup" {
		return "", fmt.Errorf("unexpected Claude SessionStart source %q", in.Source)
	}
	if in.CWD == "" {
		return "", fmt.Errorf("Claude SessionStart input has no cwd")
	}
	return in.CWD, nil
}
