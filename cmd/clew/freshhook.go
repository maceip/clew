package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
)

const freshHookTimeout = 10

func installMachineFreshnessHooks(executable string) error {
	claudePath, err := claudeUserSettingsPath()
	if err != nil {
		return err
	}
	if _, err := mutateFreshHookFile(claudePath, executable, true, false); err != nil {
		return fmt.Errorf("install Claude freshness hooks: %w", err)
	}
	codexPath, err := codexUserHooksPath()
	if err != nil {
		return err
	}
	if _, err := mutateFreshHookFile(codexPath, executable, true, true); err != nil {
		return fmt.Errorf("install Codex freshness hooks: %w", err)
	}
	if codex, err := exec.LookPath("codex"); err == nil {
		if out, runErr := exec.Command(codex, "features", "enable", "hooks").CombinedOutput(); runErr != nil {
			return fmt.Errorf("enable Codex hooks: %v: %s", runErr, strings.TrimSpace(string(out)))
		}
	}
	return nil
}

func removeMachineFreshnessHooks(executable string) error {
	claudePath, err := claudeUserSettingsPath()
	if err != nil {
		return err
	}
	if _, err := mutateFreshHookFile(claudePath, executable, false, false); err != nil {
		return fmt.Errorf("remove Claude freshness hooks: %w", err)
	}
	codexPath, err := codexUserHooksPath()
	if err != nil {
		return err
	}
	if _, err := mutateFreshHookFile(codexPath, executable, false, true); err != nil {
		return fmt.Errorf("remove Codex freshness hooks: %w", err)
	}
	return nil
}

func codexUserHooksPath() (string, error) {
	if dir := os.Getenv("CODEX_HOME"); dir != "" {
		abs, err := filepath.Abs(dir)
		if err != nil {
			return "", fmt.Errorf("resolve CODEX_HOME: %w", err)
		}
		return filepath.Join(abs, "hooks.json"), nil
	}
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return "", fmt.Errorf("resolve home for Codex hooks")
	}
	return filepath.Join(home, ".codex", "hooks.json"), nil
}

func mutateFreshHookFile(path, executable string, install, codex bool) (bool, error) {
	if executable == "" {
		return false, fmt.Errorf("freshness hook executable is empty")
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
	case err != nil:
		return false, err
	}

	original := deepJSONCopy(settings)
	hooks, err := objectField(settings, "hooks")
	if err != nil {
		return false, fmt.Errorf("%s: %w", path, err)
	}
	if hooks == nil {
		hooks = map[string]any{}
	}
	for event, arg := range map[string]string{
		"SessionStart":     "_fresh_start",
		"UserPromptSubmit": "_fresh",
	} {
		groups, err := arrayField(hooks, event)
		if err != nil {
			return false, fmt.Errorf("%s: %w", path, err)
		}
		groups = removeFreshHookHandlers(groups, executable, arg)
		if install {
			groups = append(groups, map[string]any{
				"hooks": []any{freshHookHandler(executable, arg, codex)},
			})
		}
		if len(groups) == 0 {
			delete(hooks, event)
		} else {
			hooks[event] = groups
		}
	}
	if len(hooks) == 0 {
		delete(settings, "hooks")
	} else {
		settings["hooks"] = hooks
	}
	if reflect.DeepEqual(original, settings) {
		return false, nil
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

func freshHookHandler(executable, arg string, codex bool) map[string]any {
	if codex {
		return map[string]any{
			"type":          "command",
			"command":       shellWord(executable) + " " + arg,
			"timeout":       float64(freshHookTimeout),
			"statusMessage": "checking project memory",
		}
	}
	return map[string]any{
		"type":    "command",
		"command": executable,
		"args":    []any{arg},
		"timeout": float64(freshHookTimeout),
	}
}

func removeFreshHookHandlers(groups []any, executable, arg string) []any {
	var out []any
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
		var handlers []any
		for _, raw := range rawHandlers {
			if !isFreshHookHandler(raw, executable, arg) {
				handlers = append(handlers, raw)
			}
		}
		if len(handlers) == 0 {
			continue
		}
		copyGroup := make(map[string]any, len(group))
		for key, value := range group {
			copyGroup[key] = value
		}
		copyGroup["hooks"] = handlers
		out = append(out, copyGroup)
	}
	return out
}

func isFreshHookHandler(raw any, executable, arg string) bool {
	handler, ok := raw.(map[string]any)
	if !ok || handler["type"] != "command" {
		return false
	}
	command, _ := handler["command"].(string)
	if args, ok := handler["args"].([]any); ok {
		return len(args) == 1 && args[0] == arg && sameExecutableName(command, executable)
	}
	return strings.Contains(command, filepath.Base(executable)) && strings.HasSuffix(strings.TrimSpace(command), " "+arg)
}

func shellWord(value string) string {
	return "'" + strings.ReplaceAll(value, "'", "'\\''") + "'"
}

func deepJSONCopy(in map[string]any) map[string]any {
	b, _ := json.Marshal(in)
	out := map[string]any{}
	_ = json.Unmarshal(b, &out)
	return out
}
