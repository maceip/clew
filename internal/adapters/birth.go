package adapters

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

// BirthCandidate is a recent agent session whose cwd may identify a repository
// not yet known to the machine watcher. Bootstrap policy remains outside the
// adapters package.
type BirthCandidate struct {
	Adapter string
	File    string
	CWD     string
	At      time.Time
}

// BirthCandidates returns recent Claude, Codex, and clew-wrap session cwd
// candidates, newest first and deduplicated by cleaned absolute cwd. Cursor is
// excluded because its pinned transcript format carries no unambiguous cwd.
func BirthCandidates(since time.Time) []BirthCandidate {
	var candidates []BirthCandidate
	add := func(adapter string, files []string, cwd func(string) (string, error)) {
		for _, file := range files {
			fi, err := os.Stat(file)
			if err != nil || fi.IsDir() || fi.ModTime().Before(since) {
				continue
			}
			dir, err := cwd(file)
			if err != nil || dir == "" || !filepath.IsAbs(dir) {
				continue
			}
			candidates = append(candidates, BirthCandidate{
				Adapter: adapter,
				File:    file,
				CWD:     filepath.Clean(dir),
				At:      fi.ModTime().UTC(),
			})
		}
	}

	add((&Claude{}).ID(), recentClaudeFiles(), claudeCWDChecked)
	add((&Codex{}).ID(), recentCodexFiles(), codexCWDChecked)
	add((&Wrap{}).ID(), wrapFiles(), func(file string) (string, error) {
		meta, err := wrapMetaChecked(file)
		return meta.CWD, err
	})

	sort.Slice(candidates, func(i, j int) bool {
		if !candidates[i].At.Equal(candidates[j].At) {
			return candidates[i].At.After(candidates[j].At)
		}
		if candidates[i].CWD != candidates[j].CWD {
			return candidates[i].CWD < candidates[j].CWD
		}
		if candidates[i].Adapter != candidates[j].Adapter {
			return candidates[i].Adapter < candidates[j].Adapter
		}
		return candidates[i].File < candidates[j].File
	})

	seen := map[string]bool{}
	out := make([]BirthCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		if seen[candidate.CWD] {
			continue
		}
		seen[candidate.CWD] = true
		out = append(out, candidate)
	}
	return out
}

func recentClaudeFiles() []string {
	files, _ := filepath.Glob(filepath.Join(claudeConfigDir(), "projects", "*", "*.jsonl"))
	return files
}

// claudeCWDChecked scans only the opening envelope records. Claude commonly
// writes mode/queue metadata before the first record carrying cwd, so looking
// at the first line alone would miss real births.
func claudeCWDChecked(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 64*1024), 2*1024*1024)
	for line := 0; line < 256 && scanner.Scan(); line++ {
		var env struct {
			CWD string `json:"cwd"`
		}
		if json.Unmarshal(scanner.Bytes(), &env) == nil && env.CWD != "" {
			return env.CWD, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "", fmt.Errorf("no cwd in opening Claude session records")
}
