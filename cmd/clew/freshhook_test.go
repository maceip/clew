package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestFreshHookFilesPreserveOtherHooksAndAreIdempotent(t *testing.T) {
	for _, tc := range []struct {
		name  string
		codex bool
	}{
		{name: "claude"},
		{name: "codex", codex: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "hooks.json")
			initial := `{"hooks":{"Stop":[{"hooks":[{"type":"command","command":"keep-me"}]}]}}`
			if err := os.WriteFile(path, []byte(initial), 0o600); err != nil {
				t.Fatal(err)
			}
			changed, err := mutateFreshHookFile(path, "/opt/clew bin/clew", true, tc.codex)
			if err != nil || !changed {
				t.Fatalf("first install changed=%v err=%v", changed, err)
			}
			changed, err = mutateFreshHookFile(path, "/opt/clew bin/clew", true, tc.codex)
			if err != nil || changed {
				t.Fatalf("second install changed=%v err=%v", changed, err)
			}
			b, _ := os.ReadFile(path)
			var settings map[string]any
			if err := json.Unmarshal(b, &settings); err != nil {
				t.Fatal(err)
			}
			hooks := settings["hooks"].(map[string]any)
			for _, event := range []string{"Stop", "SessionStart", "UserPromptSubmit"} {
				if _, ok := hooks[event]; !ok {
					t.Errorf("missing %s after install: %s", event, b)
				}
			}
			changed, err = mutateFreshHookFile(path, "/new/location/clew", false, tc.codex)
			if err != nil || !changed {
				t.Fatalf("remove changed=%v err=%v", changed, err)
			}
			b, _ = os.ReadFile(path)
			if err := json.Unmarshal(b, &settings); err != nil {
				t.Fatal(err)
			}
			hooks = settings["hooks"].(map[string]any)
			if _, ok := hooks["Stop"]; !ok {
				t.Fatalf("unrelated hook removed: %s", b)
			}
			for _, event := range []string{"SessionStart", "UserPromptSubmit"} {
				if _, ok := hooks[event]; ok {
					t.Errorf("managed %s survived removal: %s", event, b)
				}
			}
		})
	}
}
