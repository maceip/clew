// Package gitx implements journal storage on git (JOURNAL_SPEC §4):
// an orphan branch in the project's own remote, checked out into a
// per-repo worktree under ~/.restart/worktrees/. Git is the wire, never
// the witness (I5).
package gitx

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Branch is the orphan journal branch name. A plain branch is chosen over a
// custom ref namespace deliberately (§4): it survives every host and tool.
const Branch = "restart/journal"

var identity = []string{"-c", "user.name=restart", "-c", "user.email=journal@restart.invalid"}

// Run executes git in dir and returns trimmed combined output. Errors carry
// the command line and output (failure-trace discipline).
func Run(dir string, args ...string) (string, error) {
	full := append(append([]string{"-C", dir}, identity...), args...)
	cmd := exec.Command("git", full...)
	cmd.Env = append(os.Environ(), "GIT_TERMINAL_PROMPT=0")
	out, err := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		return s, fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, s)
	}
	return s, nil
}

func IsRepo(dir string) bool {
	out, err := Run(dir, "rev-parse", "--is-inside-work-tree")
	return err == nil && out == "true"
}

func Root(dir string) (string, error) {
	return Run(dir, "rev-parse", "--show-toplevel")
}

// RemoteName returns the repo's remote ("origin" preferred) or "" if none.
func RemoteName(dir string) string {
	out, err := Run(dir, "remote")
	if err != nil || out == "" {
		return ""
	}
	names := strings.Fields(out)
	for _, n := range names {
		if n == "origin" {
			return "origin"
		}
	}
	return names[0]
}

// RepoID is a stable short id for a repo path (worktree dir naming).
func RepoID(repoPath string) string {
	abs, err := filepath.Abs(repoPath)
	if err != nil {
		abs = repoPath
	}
	h := sha1.Sum([]byte(abs))
	return hex.EncodeToString(h[:])[:12]
}

func Home() string {
	if h := os.Getenv("RESTART_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".restart")
}

// WorktreeDir is where the journal branch is checked out for a repo.
func WorktreeDir(repoPath string) string {
	return filepath.Join(Home(), "worktrees", RepoID(repoPath))
}

// remoteBranchSHA returns the remote journal tip, or "" if absent/no remote.
func remoteBranchSHA(repoPath, remote string) string {
	if remote == "" {
		return ""
	}
	out, err := Run(repoPath, "ls-remote", remote, "refs/heads/"+Branch)
	if err != nil || out == "" {
		return ""
	}
	return strings.Fields(out)[0]
}

// genesisCommit creates the orphan root commit via plumbing (works even in a
// repo with an unborn HEAD) and returns its sha.
func genesisCommit(repoPath string) (string, error) {
	readme := "# restart journal\n\nAppend-only project journal (entries/, events/) plus generated\nprojections (journal.md, digest.md). This orphan branch shares no\nhistory with your code and never appears in PRs. See JOURNAL_SPEC.md.\n"
	blob, err := runStdin(repoPath, readme, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}
	tree, err := runStdin(repoPath, fmt.Sprintf("100644 blob %s\tREADME.md\n", blob), "mktree")
	if err != nil {
		return "", err
	}
	return commitTree(repoPath, tree, "", "restart: journal genesis")
}

func commitTree(repoPath, tree, parent, msg string) (string, error) {
	args := []string{"commit-tree", tree, "-m", msg}
	if parent != "" {
		args = []string{"commit-tree", tree, "-p", parent, "-m", msg}
	}
	cmd := exec.Command("git", append(append([]string{"-C", repoPath}, identity...), args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=restart", "GIT_AUTHOR_EMAIL=journal@restart.invalid",
		"GIT_COMMITTER_NAME=restart", "GIT_COMMITTER_EMAIL=journal@restart.invalid")
	out, err := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		return "", fmt.Errorf("git commit-tree: %w\n%s", err, s)
	}
	return s, nil
}

func runStdin(dir, stdin string, args ...string) (string, error) {
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Stdin = strings.NewReader(stdin)
	out, err := cmd.CombinedOutput()
	s := strings.TrimSpace(string(out))
	if err != nil {
		return s, fmt.Errorf("git %s: %w\n%s", strings.Join(args, " "), err, s)
	}
	return s, nil
}

// EnsureJournal makes the journal branch + worktree exist and returns the
// worktree dir. Bootstrap per §4: check the remote for an existing branch
// before creating one.
func EnsureJournal(repoPath string) (string, error) {
	if !IsRepo(repoPath) {
		return "", fmt.Errorf("%s is not a git repository", repoPath)
	}
	wt := WorktreeDir(repoPath)
	if _, err := os.Stat(filepath.Join(wt, ".git")); err == nil {
		return wt, nil
	}
	remote := RemoteName(repoPath)

	if _, err := Run(repoPath, "rev-parse", "--verify", "refs/heads/"+Branch); err != nil {
		// No local branch yet: adopt remote's if present, else genesis.
		if sha := remoteBranchSHA(repoPath, remote); sha != "" {
			if _, err := Run(repoPath, "fetch", "-q", remote, Branch); err != nil {
				return "", err
			}
			if _, err := Run(repoPath, "update-ref", "refs/heads/"+Branch, sha); err != nil {
				return "", err
			}
		} else {
			sha, err := genesisCommit(repoPath)
			if err != nil {
				return "", err
			}
			if _, err := Run(repoPath, "update-ref", "refs/heads/"+Branch, sha); err != nil {
				return "", err
			}
			if remote != "" {
				if _, err := Run(repoPath, "push", "-q", remote, Branch+":"+Branch); err != nil {
					// Race: someone pushed first. Adopt theirs (§4 bootstrap).
					if rsha := remoteBranchSHA(repoPath, remote); rsha != "" {
						Run(repoPath, "fetch", "-q", remote, Branch)
						Run(repoPath, "update-ref", "refs/heads/"+Branch, rsha)
					}
				}
			}
		}
	}
	if err := os.MkdirAll(filepath.Dir(wt), 0o755); err != nil {
		return "", err
	}
	if _, err := Run(repoPath, "worktree", "add", "-q", wt, Branch); err != nil {
		return "", err
	}
	return wt, nil
}

// SyncResult reports what a sync did (surfaced in status; I2 loudness).
type SyncResult struct {
	Committed bool
	Pushed    bool
	Pulled    bool
	Adopted   bool // unrelated remote root adopted (init race or redaction)
	Notes     []string
}

// Sync commits local additions, reconciles with the remote (pull --rebase +
// union; §4 conflict-freedom), and pushes. Safe to call repeatedly.
func Sync(repoPath string, regenerate func(wt string) error) (*SyncResult, error) {
	res := &SyncResult{}
	wt, err := EnsureJournal(repoPath)
	if err != nil {
		return nil, err
	}
	if err := commitAll(wt, res); err != nil {
		return res, err
	}
	remote := RemoteName(repoPath)
	if remote == "" {
		res.Notes = append(res.Notes, "no remote: journal is local-only")
		return res, nil
	}
	if _, err := Run(wt, "fetch", "-q", remote, Branch); err != nil {
		res.Notes = append(res.Notes, "fetch failed: "+firstLine(err.Error()))
		return res, nil // offline is normal operation, not degraded (§4.1)
	}
	fetched, err := Run(wt, "rev-parse", "FETCH_HEAD")
	if err != nil {
		return res, err
	}
	head, _ := Run(wt, "rev-parse", "HEAD")
	if head != fetched {
		if _, err := Run(wt, "merge-base", "HEAD", "FETCH_HEAD"); err != nil {
			// Unrelated roots: redaction rewrite or a lost init race.
			if err := adoptRemote(wt, res); err != nil {
				return res, err
			}
		} else if err := rebaseUnion(wt, res); err != nil {
			return res, err
		}
		res.Pulled = true
	}
	// Regenerate projections on the union, commit if changed.
	if regenerate != nil {
		if err := regenerate(wt); err != nil {
			return res, err
		}
		if err := commitAll(wt, res); err != nil {
			return res, err
		}
	}
	ahead, _ := Run(wt, "rev-list", "--count", "FETCH_HEAD..HEAD")
	if ahead != "0" && ahead != "" {
		if _, err := Run(wt, "push", "-q", remote, Branch+":"+Branch); err != nil {
			res.Notes = append(res.Notes, "push deferred: "+firstLine(err.Error()))
		} else {
			res.Pushed = true
		}
	}
	return res, nil
}

func commitAll(wt string, res *SyncResult) error {
	out, err := Run(wt, "status", "--porcelain")
	if err != nil {
		return err
	}
	if out == "" {
		return nil
	}
	if _, err := Run(wt, "add", "-A"); err != nil {
		return err
	}
	n := len(strings.Split(out, "\n"))
	if _, err := Run(wt, "commit", "-q", "-m", fmt.Sprintf("journal: %d file(s)", n)); err != nil {
		return err
	}
	res.Committed = true
	return nil
}

// rebaseUnion rebases local commits onto FETCH_HEAD. Only the generated
// projections can ever conflict (both sides regenerate them); entry/event
// files are ULID-named and single-writer by construction (§4).
func rebaseUnion(wt string, res *SyncResult) error {
	_, err := Run(wt, "rebase", "FETCH_HEAD")
	for err != nil {
		conflicts, cerr := Run(wt, "diff", "--name-only", "--diff-filter=U")
		if cerr != nil || conflicts == "" {
			Run(wt, "rebase", "--abort")
			return fmt.Errorf("journal rebase failed (should be impossible by construction): %w", err)
		}
		for _, f := range strings.Split(conflicts, "\n") {
			base := filepath.Base(f)
			if base != "journal.md" && base != "digest.md" && base != "MANIFEST.md" {
				Run(wt, "rebase", "--abort")
				return fmt.Errorf("unexpected conflict on %s (append-only law violated?)", f)
			}
			// Either side is fine: it is regenerated after sync.
			if _, err := Run(wt, "checkout", "--theirs", "--", f); err != nil {
				Run(wt, "checkout", "--ours", "--", f)
			}
			Run(wt, "add", "--", f)
		}
		_, err = Run(wt, "-c", "core.editor=true", "rebase", "--continue")
	}
	res.Notes = append(res.Notes, "rebased onto remote")
	return nil
}

// adoptRemote handles an unrelated remote root: snapshot local entry/event
// files, hard-reset to the remote, re-add files the remote lacks (§4). Files
// present remotely (e.g. a scrubbed entry after redaction) are NOT re-added,
// which is what keeps a redacted secret from resurrecting.
func adoptRemote(wt string, res *SyncResult) error {
	local := map[string][]byte{}
	for _, sub := range []string{"entries", "events"} {
		files, _ := os.ReadDir(filepath.Join(wt, sub))
		for _, f := range files {
			if f.IsDir() {
				continue
			}
			p := filepath.Join(sub, f.Name())
			b, err := os.ReadFile(filepath.Join(wt, p))
			if err == nil {
				local[p] = b
			}
		}
	}
	if _, err := Run(wt, "reset", "--hard", "FETCH_HEAD"); err != nil {
		return err
	}
	readd := 0
	for p, b := range local {
		abs := filepath.Join(wt, p)
		if _, err := os.Stat(abs); os.IsNotExist(err) {
			os.MkdirAll(filepath.Dir(abs), 0o755)
			if err := os.WriteFile(abs, b, 0o644); err != nil {
				return err
			}
			readd++
		}
	}
	res.Adopted = true
	res.Notes = append(res.Notes, fmt.Sprintf("adopted remote journal root, re-added %d local file(s)", readd))
	return commitAll(wt, res)
}

// RewriteRoot creates a fresh single-commit root of the worktree's current
// content and force-pushes it — the one sanctioned rewrite (§4, redaction).
func RewriteRoot(repoPath, msg string) error {
	wt, err := EnsureJournal(repoPath)
	if err != nil {
		return err
	}
	if _, err := Run(wt, "add", "-A"); err != nil {
		return err
	}
	tree, err := Run(wt, "write-tree")
	if err != nil {
		return err
	}
	sha, err := commitTree(wt, tree, "", msg)
	if err != nil {
		return err
	}
	if _, err := Run(wt, "reset", "--hard", sha); err != nil {
		return err
	}
	if remote := RemoteName(repoPath); remote != "" {
		if _, err := Run(wt, "push", "-q", "--force", remote, Branch+":"+Branch); err != nil {
			return fmt.Errorf("redaction force-push failed (secret still on remote!): %w", err)
		}
	}
	return nil
}

// Show returns a file's content at the journal branch tip without a worktree
// (used by read-only surfaces on machines that only fetched).
func Show(repoPath, file string) (string, error) {
	return Run(repoPath, "show", Branch+":"+file)
}

func firstLine(s string) string {
	if i := strings.IndexByte(s, '\n'); i >= 0 {
		return s[:i]
	}
	return s
}
