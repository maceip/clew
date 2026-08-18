// Package gitx implements journal storage on git (JOURNAL_SPEC §4):
// an orphan branch in the project's own remote, checked out into a
// per-repo worktree under ~/.clew/worktrees/. Git is the wire, never
// the witness (I5).
package gitx

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// Branch is the orphan journal branch name. A plain branch is chosen over a
// custom ref namespace deliberately (§4): it survives every host and tool.
const Branch = "clew/journal"

var identity = []string{"-c", "user.name=clew", "-c", "user.email=journal@clew.invalid"}

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
	if h := os.Getenv("CLEW_HOME"); h != "" {
		return h
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".clew")
}

// WorktreeDir is where the journal branch is checked out for a repo.
func WorktreeDir(repoPath string) string {
	id := RepoID(repoPath)
	if configured := ConfiguredJournalID(repoPath); configured != "" {
		id = configured
	}
	return filepath.Join(Home(), "worktrees", id)
}

// ConfiguredJournalID returns the validated, persistent journal incarnation
// selected when a fresh git repository reuses a previously enrolled path. It
// is empty for the ordinary path-derived namespace. The value lives in the
// repository's local git config, so ambient seed identity can distinguish
// incarnations without consulting machine state or the predecessor worktree.
func ConfiguredJournalID(repoPath string) string {
	configured, err := Run(repoPath, "config", "--local", "--get", "clew.journal-id")
	if err != nil || !validJournalID(configured) {
		return ""
	}
	return configured
}

func validJournalID(id string) bool {
	if id == "" || id == "." || id == ".." {
		return false
	}
	for _, r := range id {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' || r == '.' {
			continue
		}
		return false
	}
	return true
}

// JournalReady proves that the external journal worktree belongs to the
// current checkout incarnation. A pathname alone is not repository identity:
// a directory can be moved away and replaced by a fresh `git init` at the
// same path. Birth must never bind that newborn to the predecessor's lore.
func JournalReady(repoPath string) bool {
	if !IsRepo(repoPath) {
		return false
	}
	wt := WorktreeDir(repoPath)
	if _, err := os.Stat(filepath.Join(wt, ".git")); err != nil {
		return false
	}
	repoCommon, err := gitCommonDir(repoPath)
	if err != nil {
		return false
	}
	worktreeCommon, err := gitCommonDir(wt)
	return err == nil && samePath(repoCommon, worktreeCommon)
}

// BirthReady is the local, network-free proof used by SessionStart. The
// marker is written only after registration and first atomic materialization;
// pairing it with JournalReady prevents a fresh repo at a reused path from
// consuming the predecessor's context.
func BirthReady(repoPath string) bool {
	if !JournalReady(repoPath) {
		return false
	}
	value, err := Run(repoPath, "config", "--local", "--get", "clew.birth-ready")
	return err == nil && value == "true"
}

// HasBirthIncarnation reports whether the current .git contains any durable
// evidence that clew has already enrolled this repository incarnation. It is
// deliberately independent of the external worktree's health: a missing or
// damaged worktree is a repair case, not proof that the path now names a new
// repository whose disposable machine state should be reset.
func HasBirthIncarnation(repoPath string) bool {
	if !IsRepo(repoPath) {
		return false
	}
	for _, key := range []string{"clew.birth-ready", "clew.journal-id"} {
		if _, err := Run(repoPath, "config", "--local", "--get", key); err == nil {
			return true
		}
	}
	_, err := Run(repoPath, "show-ref", "--verify", "--quiet", "refs/heads/"+Branch)
	return err == nil
}

func MarkBirthReady(repoPath string) error {
	_, err := Run(repoPath, "config", "--local", "clew.birth-ready", "true")
	return err
}

func gitCommonDir(dir string) (string, error) {
	out, err := Run(dir, "rev-parse", "--path-format=absolute", "--git-common-dir")
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(out) {
		out = filepath.Join(dir, out)
	}
	abs, err := filepath.Abs(out)
	if err != nil {
		return "", err
	}
	if resolved, resolveErr := filepath.EvalSymlinks(abs); resolveErr == nil {
		abs = resolved
	}
	return filepath.Clean(abs), nil
}

func samePath(a, b string) bool {
	ai, aerr := os.Stat(a)
	bi, berr := os.Stat(b)
	if aerr == nil && berr == nil {
		return os.SameFile(ai, bi)
	}
	return filepath.Clean(a) == filepath.Clean(b)
}

func chooseFreshJournalID(repoPath string) (string, error) {
	for attempt := 0; attempt < 8; attempt++ {
		var nonce [6]byte
		if _, err := rand.Read(nonce[:]); err != nil {
			return "", fmt.Errorf("journal incarnation id: %w", err)
		}
		id := RepoID(repoPath) + "-" + hex.EncodeToString(nonce[:])
		if _, err := os.Stat(filepath.Join(Home(), "worktrees", id)); os.IsNotExist(err) {
			if _, err := Run(repoPath, "config", "--local", "clew.journal-id", id); err != nil {
				return "", err
			}
			return id, nil
		}
	}
	return "", fmt.Errorf("could not allocate a journal worktree id for %s", repoPath)
}

// AssignFreshJournalIncarnation persists a new external-worktree namespace in
// the current repository. Birth calls this only after proving that a
// registered pathname now contains a fresh git repository. Keeping the
// assignment explicit means the newborn receives a distinct seed identity
// even when the predecessor worktree has already been manually removed.
func AssignFreshJournalIncarnation(repoPath string) (string, error) {
	if !IsRepo(repoPath) {
		return "", fmt.Errorf("%s is not a git repository", repoPath)
	}
	return chooseFreshJournalID(repoPath)
}

// lookupRemoteBranch returns the remote journal tip. An empty result means the
// branch is absent; transport/authentication failures remain distinguishable.
func lookupRemoteBranch(repoPath, remote string) (string, error) {
	if remote == "" {
		return "", nil
	}
	out, err := Run(repoPath, "ls-remote", remote, "refs/heads/"+Branch)
	if err != nil {
		return "", err
	}
	if out == "" {
		return "", nil
	}
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return "", fmt.Errorf("git ls-remote %s returned no journal tip", remote)
	}
	return fields[0], nil
}

// remoteBranchSHA returns the remote journal tip, or "" if absent/no remote.
// Bootstrap deliberately treats an unavailable remote like an absent one so
// the journal remains useful offline; destructive rewrites use the strict
// lookup above through their verified-sync lease instead.
func remoteBranchSHA(repoPath, remote string) string {
	sha, err := lookupRemoteBranch(repoPath, remote)
	if err != nil {
		return ""
	}
	return sha
}

// genesisCommit creates the orphan root commit via plumbing (works even in a
// repo with an unborn HEAD) and returns its sha.
func genesisCommit(repoPath string) (string, error) {
	readme := "# clew journal\n\nAppend-only project journal (entries/, events/) plus generated\nprojections (journal.md, digest.md). This orphan branch shares no\nhistory with your code and never appears in PRs. See JOURNAL_SPEC.md.\n"
	blob, err := runStdin(repoPath, readme, "hash-object", "-w", "--stdin")
	if err != nil {
		return "", err
	}
	tree, err := runStdin(repoPath, fmt.Sprintf("100644 blob %s\tREADME.md\n", blob), "mktree")
	if err != nil {
		return "", err
	}
	return commitTree(repoPath, tree, "", "clew: journal genesis")
}

func commitTree(repoPath, tree, parent, msg string) (string, error) {
	args := []string{"commit-tree", tree, "-m", msg}
	if parent != "" {
		args = []string{"commit-tree", tree, "-p", parent, "-m", msg}
	}
	cmd := exec.Command("git", append(append([]string{"-C", repoPath}, identity...), args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=clew", "GIT_AUTHOR_EMAIL=journal@clew.invalid",
		"GIT_COMMITTER_NAME=clew", "GIT_COMMITTER_EMAIL=journal@clew.invalid")
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
		if JournalReady(repoPath) {
			return wt, nil
		}
		// Preserve the stale predecessor worktree as an explicit future
		// lineage candidate; allocate a new namespace for this incarnation.
		if _, err := AssignFreshJournalIncarnation(repoPath); err != nil {
			return "", err
		}
		wt = WorktreeDir(repoPath)
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

	// Remote/RemoteURL identify the wire endpoint used by this sync. RemoteTip
	// is populated only when that endpoint's journal tip was observed and all
	// required local commits were successfully published. Root rewrites bind a
	// force-with-lease to this exact observation.
	Remote    string
	RemoteURL string
	RemoteTip string
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
	// Maintain the local journal projection before any network operation. An
	// offline fetch is normal, and must not leave SEED.md one journal change
	// behind merely because the remote could not be contacted.
	if regenerate != nil {
		if err := regenerate(wt); err != nil {
			return res, err
		}
		if err := commitAll(wt, res); err != nil {
			return res, err
		}
	}
	remote := RemoteName(repoPath)
	res.Remote = remote
	if remote == "" {
		res.Notes = append(res.Notes, "no remote: journal is local-only")
		return res, nil
	}
	res.RemoteURL, _ = Run(repoPath, "remote", "get-url", remote)
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
			return res, nil
		} else {
			res.Pushed = true
		}
		fetched, err = Run(wt, "rev-parse", "HEAD")
		if err != nil {
			return res, err
		}
	}
	res.RemoteTip = fetched
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
			if base != "journal.md" && base != "digest.md" && base != "SEED.md" && base != "MANIFEST.md" {
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
	for _, sub := range []string{"entries", "events", "lineage"} {
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

// RewriteLease binds a destructive root rewrite to the exact remote journal
// tip observed by a completed Sync. Its fields are intentionally opaque: a
// caller must obtain one through RewriteLeaseFromSync rather than guessing at
// whatever tip happens to exist immediately before the force-push.
type RewriteLease struct {
	remote    string
	remoteURL string
	expected  string
}

// RewriteLeaseError means the remote advanced after the sync that authorized
// a root rewrite. The caller must re-sync the append-only union, re-apply its
// scrub to that union, and derive a fresh lease before retrying.
type RewriteLeaseError struct {
	Remote   string
	Expected string
	Actual   string
	Err      error
}

func (e *RewriteLeaseError) Error() string {
	actual := e.Actual
	if actual == "" {
		actual = "<absent>"
	}
	return fmt.Sprintf("journal root rewrite lease rejected: remote %s advanced from %s to %s; re-sync and re-scrub before retrying: %v",
		e.Remote, e.Expected, actual, e.Err)
}

func (e *RewriteLeaseError) Unwrap() error { return e.Err }

// IsRewriteLeaseError reports the one retryable RewriteRoot failure. Other
// failures (including inability to verify the remote after a failed push) are
// deliberately loud and must not be treated as a concurrency retry.
func IsRewriteLeaseError(err error) bool {
	var target *RewriteLeaseError
	return errors.As(err, &target)
}

// RewriteLeaseFromSync turns a fully published/observed sync result into the
// compare-and-swap token required by RewriteRoot. Local-only journals receive
// a valid empty lease; configured remotes require an exact endpoint and tip.
func RewriteLeaseFromSync(repoPath string, synced *SyncResult) (RewriteLease, error) {
	if synced == nil {
		return RewriteLease{}, fmt.Errorf("journal root rewrite requires a completed sync")
	}
	remote := RemoteName(repoPath)
	if remote == "" {
		if synced.Remote != "" {
			return RewriteLease{}, fmt.Errorf("journal remote changed after sync: observed %s, now local-only", synced.Remote)
		}
		return RewriteLease{}, nil
	}
	if synced.Remote != remote {
		return RewriteLease{}, fmt.Errorf("journal remote changed after sync: observed %s, now %s", synced.Remote, remote)
	}
	remoteURL, err := Run(repoPath, "remote", "get-url", remote)
	if err != nil {
		return RewriteLease{}, fmt.Errorf("verify journal remote before root rewrite: %w", err)
	}
	if synced.RemoteURL == "" || synced.RemoteURL != remoteURL {
		return RewriteLease{}, fmt.Errorf("journal remote URL changed after sync: observed %q, now %q", synced.RemoteURL, remoteURL)
	}
	if synced.RemoteTip == "" {
		return RewriteLease{}, fmt.Errorf("journal root rewrite requires verified remote state; the last sync did not publish or observe a usable journal tip")
	}
	return RewriteLease{remote: remote, remoteURL: remoteURL, expected: synced.RemoteTip}, nil
}

// RewriteRoot creates a fresh single-commit root of the worktree's current
// content and force-pushes it — the one sanctioned rewrite (§4, redaction).
// The force is a compare-and-swap against lease: a concurrent append can never
// be discarded. The local branch is reset only after the remote accepts the
// new root, leaving a rejected caller on its related pre-rewrite history so a
// normal Sync can recover the union before the caller re-scrubs and retries.
func RewriteRoot(repoPath, msg string, lease RewriteLease) error {
	wt, err := EnsureJournal(repoPath)
	if err != nil {
		return err
	}
	remote := RemoteName(repoPath)
	if remote != lease.remote {
		return fmt.Errorf("journal remote changed after rewrite lease was issued: leased %s, now %s", lease.remote, remote)
	}
	if remote != "" {
		remoteURL, err := Run(repoPath, "remote", "get-url", remote)
		if err != nil {
			return fmt.Errorf("verify journal rewrite remote: %w", err)
		}
		if remoteURL != lease.remoteURL {
			return fmt.Errorf("journal remote URL changed after rewrite lease was issued: leased %q, now %q", lease.remoteURL, remoteURL)
		}
		if lease.expected == "" {
			return fmt.Errorf("journal root rewrite refused: configured remote has no verified expected tip")
		}
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
	if remote != "" {
		ref := "refs/heads/" + Branch
		_, pushErr := Run(wt, "push", "-q", "--force-with-lease="+ref+":"+lease.expected, remote, sha+":"+ref)
		if pushErr != nil {
			actual, lookupErr := lookupRemoteBranch(repoPath, remote)
			if lookupErr == nil && actual != lease.expected {
				return &RewriteLeaseError{Remote: remote, Expected: lease.expected, Actual: actual, Err: pushErr}
			}
			if lookupErr != nil {
				return fmt.Errorf("redaction force-with-lease failed and remote state could not be verified (secret still on remote!): %w; verify remote: %v", pushErr, lookupErr)
			}
			return fmt.Errorf("redaction force-with-lease failed (secret still on remote!): %w", pushErr)
		}
	}
	if _, err := Run(wt, "reset", "--hard", sha); err != nil {
		return fmt.Errorf("root rewrite was accepted but the local journal could not reset to it: %w", err)
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
