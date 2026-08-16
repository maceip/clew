package main

import (
	"fmt"
	"os"
	"time"

	"restart/internal/state"
	"restart/internal/wrapx"
)

// cmdWrap: PTY tee for agents without session files (§5.1). The transcript
// lands in ~/.restart/raw/ where the wrap adapter tails it; nudges can be
// injected at prompt boundaries (§8.1 matrix).
func cmdWrap(args []string) error {
	nudge := true
	// Accept: restart wrap [--no-nudge] -- <argv…>   (also tolerate no --).
	var argv []string
	for i := 0; i < len(args); i++ {
		switch args[i] {
		case "--no-nudge":
			nudge = false
		case "--":
			argv = args[i+1:]
			i = len(args)
		default:
			argv = args[i:]
			i = len(args)
		}
	}
	if len(argv) == 0 {
		return fmt.Errorf("usage: restart wrap -- <agent argv…>")
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()

	opts := &wrapx.Options{Argv: argv, Nudge: nudge, Surface: a.cfg.Surface}
	cwd, _ := os.Getwd()
	repo := a.db.RepoFor(cwd)

	// Pre-register the session so status shows it while it runs.
	done := make(chan struct{})
	go func() {
		tick := time.NewTicker(10 * time.Second)
		defer tick.Stop()
		for {
			select {
			case <-done:
				return
			case <-tick.C:
				if opts.SessionID != "" && repo != "" {
					a.db.UpsertSession(state.Session{
						ID: "wrap:" + opts.SessionID, Adapter: "wrap",
						Agent: "wrap:" + argv[0], File: opts.LogPath, RepoPath: repo,
						Surface: a.cfg.Surface, StartedAt: time.Now(), LastActivity: time.Now(),
					})
				}
			}
		}
	}()
	code, err := wrapx.Run(opts)
	close(done)
	if err != nil {
		return err
	}
	os.Exit(code)
	return nil
}
