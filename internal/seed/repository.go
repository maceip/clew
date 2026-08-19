package seed

import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/journal"
)

// BuildForRepo joins a journal with the small set of repository facts needed
// by an ambient seed. Owner laws are deliberately absent: they are a separate
// ambient layer, not lineage lore.
func BuildForRepo(repoPath string, j *journal.Journal, st map[string]*journal.Computed) (*Snapshot, error) {
	in := BuildInputForRepo(repoPath, j)
	return Build(j, st, in)
}

// BuildInputForRepo is provider-free and deterministic for the same journal
// and checked-out code state. It never runs archaeology or reads transcripts.
func BuildInputForRepo(repoPath string, j *journal.Journal) BuildInput {
	abs, _ := filepath.Abs(repoPath)
	if abs != "" {
		repoPath = filepath.Clean(abs)
	}
	name := filepath.Base(repoPath)
	remote := ""
	if remoteName := gitx.RemoteName(repoPath); remoteName != "" {
		remote, _ = gitx.Run(repoPath, "remote", "get-url", remoteName)
	}
	identity := "path:" + repoPath
	if remote != "" {
		identity = "remote:" + remote
	} else if incarnation := gitx.ConfiguredJournalID(repoPath); incarnation != "" && incarnation != gitx.RepoID(repoPath) {
		// A path can be moved aside and reused by a fresh `git init`. gitx
		// persists a new journal incarnation in that case; include it in the
		// newborn seed identity so explicit lineage from the moved predecessor
		// is not misclassified as self-lineage. Existing seeds remain the
		// canonical identity in materialize.BuildSeedForRepo, and remote-backed
		// repositories continue to use their shared remote identity.
		identity = "journal-incarnation:" + incarnation
	}
	sum := sha256.Sum256([]byte(identity))
	repository := Repository{
		ID:     "r" + hex.EncodeToString(sum[:])[:24],
		Name:   name,
		Remote: remote,
	}

	// Ranking needs a small, high-signal vocabulary. Whole READMEs and entry
	// bodies dilute cosine overlap and bloat every ambient seed.
	topicText := []string{name}
	for _, readme := range []string{"README.md", "README", "readme.md"} {
		if b, err := os.ReadFile(filepath.Join(repoPath, readme)); err == nil {
			if len(b) > 256*1024 {
				b = b[:256*1024]
			}
			headings := 0
			for _, line := range strings.Split(string(b), "\n") {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "#") {
					topicText = append(topicText, strings.TrimSpace(strings.TrimLeft(line, "#")))
					headings++
					if headings == 32 {
						break
					}
				}
			}
			break
		}
	}
	if j != nil {
		for _, entry := range j.Entries {
			topicText = append(topicText, entry.Title, strings.Join(entry.Tags, " "))
		}
	}
	topics := Tokens(strings.Join(topicText, " "))
	if len(topics) > 128 {
		topics = topics[:128]
	}

	lifecycle := Lifecycle{State: "active"}
	if j != nil {
		for _, name := range []string{"TOMBSTONE.md", "TOMBSTONE"} {
			path := filepath.Join(j.Dir, name)
			if b, err := os.ReadFile(path); err == nil {
				lifecycle.State = "tombstoned"
				lifecycle.Note = strings.TrimSpace(string(b))
				if info, statErr := os.Stat(path); statErr == nil {
					lifecycle.At = info.ModTime().UTC()
				}
				if lifecycle.At.IsZero() {
					lifecycle.At = latestJournalTime(j)
				}
				break
			}
		}
	}

	var pin *OrganBankPin
	if head, err := gitx.Run(repoPath, "rev-parse", "--verify", "HEAD"); err == nil && head != "" {
		at := time.Time{}
		if stamp, stampErr := gitx.Run(repoPath, "show", "-s", "--format=%cI", head); stampErr == nil {
			at, _ = time.Parse(time.RFC3339, stamp)
		}
		dirty, _ := gitx.Run(repoPath, "status", "--porcelain", "--untracked-files=no")
		pin = &OrganBankPin{Remote: remote, Commit: head, Dirty: dirty != "", At: at.UTC()}
	}

	return BuildInput{Repository: repository, Lifecycle: lifecycle, Topics: topics, OrganBank: pin}
}

func latestJournalTime(j *journal.Journal) time.Time {
	var latest time.Time
	if j == nil {
		return latest
	}
	for _, entry := range j.Entries {
		if at := entry.Created(); at.After(latest) {
			latest = at
		}
	}
	for _, event := range j.Events {
		if event.At.After(latest) {
			latest = event.At
		}
	}
	return latest.UTC()
}
