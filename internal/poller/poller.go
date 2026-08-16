// Package poller watches registered repos (JOURNAL_SPEC §5.2): HEAD, branch
// list, dirty file set, ahead/behind, new commits with file lists.
// Ground truth for WHAT changed; transcripts are ground truth for WHO/WHY.
// Commit→session attribution: time window + footprint overlap + author;
// unattributed is a normal, displayed state — never guessed.
package poller

import (
	"path/filepath"
	"strings"
	"time"

	"clew/internal/gitx"
	"clew/internal/state"
)

type Snapshot struct {
	RepoPath   string
	Head       string
	Branch     string
	Dirty      []string // repo-relative dirty paths
	Ahead      int
	Behind     int
	NewCommits []state.Commit
}

const attributionWindow = 15 * time.Minute

// Poll gathers a snapshot and records unseen commits (with attribution).
func Poll(db *state.DB, repoPath string) (*Snapshot, error) {
	s := &Snapshot{RepoPath: repoPath}
	head, err := gitx.Run(repoPath, "rev-parse", "HEAD")
	if err == nil {
		s.Head = head
	}
	if br, err := gitx.Run(repoPath, "rev-parse", "--abbrev-ref", "HEAD"); err == nil {
		s.Branch = br
	}
	if out, err := gitx.Run(repoPath, "status", "--porcelain"); err == nil && out != "" {
		for _, l := range strings.Split(out, "\n") {
			if len(l) > 3 {
				p := strings.TrimSpace(l[3:])
				if i := strings.Index(p, " -> "); i >= 0 {
					p = p[i+4:]
				}
				s.Dirty = append(s.Dirty, p)
			}
		}
	}
	if out, err := gitx.Run(repoPath, "rev-list", "--left-right", "--count", "@{upstream}...HEAD"); err == nil {
		parts := strings.Fields(out)
		if len(parts) == 2 {
			s.Behind = atoi(parts[0])
			s.Ahead = atoi(parts[1])
		}
	}

	// New commits across all branches, last 14 days. The journal branch is
	// coordination data, not code reality — excluded from observation.
	out, err := gitx.Run(repoPath, "log",
		"--exclude=refs/heads/"+gitx.Branch, "--exclude=refs/remotes/*/"+gitx.Branch,
		"--all", "--since=14.days",
		"--date=iso-strict", "--pretty=%x01%H%x1f%an%x1f%aI%x1f%s", "--name-only")
	if err != nil {
		return s, nil // empty repo etc.: fine
	}
	for _, block := range strings.Split(out, "\x01") {
		block = strings.TrimSpace(block)
		if block == "" {
			continue
		}
		lines := strings.Split(block, "\n")
		fields := strings.Split(lines[0], "\x1f")
		if len(fields) != 4 {
			continue
		}
		c := state.Commit{RepoPath: repoPath, SHA: fields[0], Author: fields[1], Subject: fields[3]}
		c.At, _ = time.Parse(time.RFC3339, fields[2])
		for _, f := range lines[1:] {
			if f = strings.TrimSpace(f); f != "" {
				c.Files = append(c.Files, f)
			}
		}
		if db.CommitSeen(repoPath, c.SHA) {
			continue
		}
		c.SessionID = attribute(db, repoPath, &c)
		if err := db.AddCommit(c); err == nil {
			s.NewCommits = append(s.NewCommits, c)
		}
	}
	return s, nil
}

// attribute matches a commit to a live session: commit time inside the
// session's activity window (±15 min) AND footprint overlap. Best overlap
// wins; no match = unattributed ("" — displayed, never guessed).
func attribute(db *state.DB, repoPath string, c *state.Commit) string {
	sessions := db.LiveSessions(repoPath, 48*time.Hour)
	best, bestOverlap := "", 0
	for _, sess := range sessions {
		if c.At.Before(sess.StartedAt.Add(-attributionWindow)) || c.At.After(sess.LastActivity.Add(attributionWindow)) {
			continue
		}
		fp := db.Footprints(sess.ID)
		overlap := 0
		for _, cf := range c.Files {
			for _, sf := range fp {
				if cf == sf || strings.HasSuffix(sf, "/"+cf) || filepath.Base(sf) == filepath.Base(cf) && strings.HasSuffix(sf, cf) {
					overlap++
					break
				}
			}
		}
		if overlap > bestOverlap {
			best, bestOverlap = sess.ID, overlap
		}
	}
	return best
}

func atoi(s string) int {
	n := 0
	for _, r := range s {
		if r < '0' || r > '9' {
			return 0
		}
		n = n*10 + int(r-'0')
	}
	return n
}
