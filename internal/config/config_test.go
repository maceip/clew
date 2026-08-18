package config

import (
	"os"
	"strings"
	"testing"
)

func TestLoadOwnerRemoteAndLegacyCompatibility(t *testing.T) {
	t.Run("legacy config has local owner journal", func(t *testing.T) {
		t.Setenv("CLEW_HOME", t.TempDir())
		legacy := "surface: desk\nextractor:\n  provider: off\n"
		if err := os.WriteFile(Path(), []byte(legacy), 0o644); err != nil {
			t.Fatal(err)
		}
		got := Load()
		if got.Owner.Remote != "" {
			t.Fatalf("legacy owner remote = %q, want local-only", got.Owner.Remote)
		}
		if got.Extractor.DailyCapTokens != 200_000 || got.Extractor.SessionPct != 2.0 {
			t.Fatalf("legacy defaults lost: %+v", got.Extractor)
		}
	})

	t.Run("owner remote is parsed", func(t *testing.T) {
		t.Setenv("CLEW_HOME", t.TempDir())
		body := "owner:\n  remote: git@example.test:me/clew-owner.git\n"
		if err := os.WriteFile(Path(), []byte(body), 0o644); err != nil {
			t.Fatal(err)
		}
		if got := Load().Owner.Remote; got != "git@example.test:me/clew-owner.git" {
			t.Fatalf("owner remote = %q", got)
		}
	})
}

func TestWriteSkeletonIncludesOwnerRemote(t *testing.T) {
	t.Setenv("CLEW_HOME", t.TempDir())
	if err := WriteSkeleton(); err != nil {
		t.Fatal(err)
	}
	b, err := os.ReadFile(Path())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(b), "owner:\n    remote: \"\"") {
		t.Fatalf("skeleton does not expose owner.remote:\n%s", b)
	}
}
