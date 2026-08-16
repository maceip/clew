// Package wrapx is the PTY tee (JOURNAL_SPEC §5.1 `wrap`): a transcript-
// equivalent for agents with no accessible session store. Input bytes are
// logged as user turns, output as assistant turns, to ~/.restart/raw/.
// At a prompt boundary the wrapper can inject one nudge line (§8.1 matrix).
package wrapx

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/creack/pty"
	"golang.org/x/term"

	"restart/internal/gitx"
	"restart/internal/ids"
)

type Options struct {
	Argv      []string
	Nudge     bool // inject nudge lines at prompt boundaries
	Surface   string
	SessionID string // assigned by caller for registration; generated if empty
	LogPath   string // set by Run for the caller
}

type logger struct {
	mu   sync.Mutex
	f    *os.File
	dir  string // "in" | "out" accumulation
	buf  strings.Builder
	last time.Time
}

var ansiRe = regexp.MustCompile(`\x1b\[[0-9;?]*[a-zA-Z]|\x1b\][^\x07]*\x07|\x1b[()][A-Z0-9]|[\x00-\x08\x0b-\x1f\x7f]`)

func (l *logger) write(dir, text string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.dir != dir && l.buf.Len() > 0 {
		l.flushLocked()
	}
	l.dir = dir
	l.buf.WriteString(text)
	l.last = time.Now()
	if l.buf.Len() > 8192 || strings.Contains(text, "\n") {
		l.flushLocked()
	}
}

func (l *logger) flushLocked() {
	txt := ansiRe.ReplaceAllString(l.buf.String(), "")
	txt = strings.TrimSpace(strings.ReplaceAll(txt, "\r\n", "\n"))
	l.buf.Reset()
	if txt == "" {
		return
	}
	rec, _ := json.Marshal(map[string]string{
		"at": time.Now().UTC().Format(time.RFC3339), "dir": l.dir, "text": txt,
	})
	l.f.Write(append(rec, '\n'))
}

func (l *logger) flushIfIdle(d time.Duration) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.buf.Len() > 0 && time.Since(l.last) > d {
		l.flushLocked()
	}
}

// Run executes argv under a PTY tee. Returns the child's exit code.
func Run(o *Options) (int, error) {
	if len(o.Argv) == 0 {
		return 2, fmt.Errorf("wrap: no command given")
	}
	cwd, _ := os.Getwd()
	rawDir := filepath.Join(gitx.Home(), "raw")
	if err := os.MkdirAll(rawDir, 0o755); err != nil {
		return 2, err
	}
	session := o.SessionID
	if session == "" {
		session = ids.NewEntry(time.Now())[1:] // bare ULID
	}
	o.SessionID = session
	logPath := filepath.Join(rawDir, session+".jsonl")
	o.LogPath = logPath
	f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return 2, err
	}
	defer f.Close()
	meta, _ := json.Marshal(map[string]any{
		"kind": "meta", "argv": o.Argv, "cwd": cwd, "session": session,
		"surface": o.Surface, "at": time.Now().UTC().Format(time.RFC3339),
	})
	f.Write(append(meta, '\n'))

	cmd := exec.Command(o.Argv[0], o.Argv[1:]...)
	ptmx, err := pty.Start(cmd)
	if err != nil {
		return 2, err
	}
	defer ptmx.Close()

	// Window size propagation.
	winch := make(chan os.Signal, 1)
	signal.Notify(winch, syscall.SIGWINCH)
	go func() {
		for range winch {
			pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	winch <- syscall.SIGWINCH
	defer signal.Stop(winch)

	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err == nil {
		defer term.Restore(int(os.Stdin.Fd()), oldState)
	}

	lg := &logger{f: f}
	var lastOut atomic64

	// stdin → pty (tee "in")
	go func() {
		buf := make([]byte, 4096)
		for {
			n, err := os.Stdin.Read(buf)
			if n > 0 {
				ptmx.Write(buf[:n])
				lg.write("in", string(buf[:n]))
			}
			if err != nil {
				return
			}
		}
	}()

	// nudge injection: at a prompt boundary (output idle ≥ 2s), inject one
	// pending line from .restart/nudge.md, consuming it.
	stop := make(chan struct{})
	if o.Nudge {
		go func() {
			nudgePath := filepath.Join(cwd, ".restart", "nudge.md")
			tick := time.NewTicker(3 * time.Second)
			defer tick.Stop()
			for {
				select {
				case <-stop:
					return
				case <-tick.C:
					if time.Since(lastOut.get()) < 2*time.Second {
						continue // agent is mid-output; not a prompt boundary
					}
					line := takeFirstLine(nudgePath)
					if line == "" {
						continue
					}
					ptmx.Write([]byte(line + "\r"))
					lg.write("in", "[nudge] "+line+"\n")
				}
			}
		}()
	}

	// idle flusher
	go func() {
		tick := time.NewTicker(300 * time.Millisecond)
		defer tick.Stop()
		for {
			select {
			case <-stop:
				return
			case <-tick.C:
				lg.flushIfIdle(400 * time.Millisecond)
			}
		}
	}()

	// pty → stdout (tee "out")
	buf := make([]byte, 8192)
	for {
		n, err := ptmx.Read(buf)
		if n > 0 {
			os.Stdout.Write(buf[:n])
			lg.write("out", string(buf[:n]))
			lastOut.set(time.Now())
		}
		if err != nil {
			break
		}
	}
	close(stop)
	lg.flushIfIdle(0)
	err = cmd.Wait()
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), nil
	}
	if err != nil {
		return 1, err
	}
	return 0, nil
}

// takeFirstLine pops the first line off the nudge file (consume semantics).
func takeFirstLine(path string) string {
	b, err := os.ReadFile(path)
	if err != nil || len(b) == 0 {
		return ""
	}
	lines := strings.SplitN(string(b), "\n", 2)
	rest := ""
	if len(lines) == 2 {
		rest = lines[1]
	}
	os.WriteFile(path, []byte(rest), 0o644)
	return strings.TrimSpace(lines[0])
}

type atomic64 struct {
	mu sync.Mutex
	t  time.Time
}

func (a *atomic64) set(t time.Time) { a.mu.Lock(); a.t = t; a.mu.Unlock() }
func (a *atomic64) get() time.Time  { a.mu.Lock(); defer a.mu.Unlock(); return a.t }
