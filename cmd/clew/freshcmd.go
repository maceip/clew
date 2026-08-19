package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/maceip/clew/internal/gitx"
	"github.com/maceip/clew/internal/journal"
	"github.com/maceip/clew/internal/model"
	"github.com/maceip/clew/internal/state"
)

const freshHookInputLimit = 1 << 20

type freshHookInput struct {
	CWD           string `json:"cwd"`
	SessionID     string `json:"session_id"`
	HookEventName string `json:"hook_event_name"`
}

type freshHookOutput struct {
	HookSpecificOutput struct {
		HookEventName     string `json:"hookEventName"`
		AdditionalContext string `json:"additionalContext"`
	} `json:"hookSpecificOutput"`
}

func cmdFreshStart(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("internal _fresh_start takes no arguments")
	}
	return runFreshStartHook(os.Stdin)
}

func cmdFresh(args []string) error {
	if len(args) != 0 {
		return fmt.Errorf("internal _fresh takes no arguments")
	}
	return runFreshHook(os.Stdin, os.Stdout)
}

// runFreshStartHook snapshots the journal before the session's first prompt.
// It emits nothing: SessionStart has its own birth/context path, while this
// marker exists only so later prompt boundaries can deliver a true delta.
func runFreshStartHook(input io.Reader) error {
	in, err := parseFreshHookInput(input, "SessionStart")
	if err != nil {
		return err
	}
	root, j, ok, err := freshJournal(in.CWD)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	db, err := state.Open(state.DefaultPath())
	if err != nil {
		return err
	}
	defer db.Close()
	return db.Set(freshSessionKey(root, in.SessionID), newestEntryID(j))
}

// runFreshHook injects each decision newer than this session's watermark once.
// The output is Claude/Codex hook JSON and deliberately uses the hard register:
// it is machine context, not a human surface and never an instruction source.
func runFreshHook(input io.Reader, output io.Writer) error {
	in, err := parseFreshHookInput(input, "UserPromptSubmit")
	if err != nil {
		return err
	}
	root, j, ok, err := freshJournal(in.CWD)
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}
	db, err := state.Open(state.DefaultPath())
	if err != nil {
		return err
	}
	defer db.Close()

	key := freshSessionKey(root, in.SessionID)
	watermark := db.Get(key)
	newest := newestEntryID(j)
	if watermark == "" {
		// A hook installed into an already-running session has no honest start
		// snapshot. Baseline it now instead of flooding the model with history.
		return db.Set(key, newest)
	}

	var decisions []*model.Entry
	for id, entry := range j.Entries {
		if id > watermark && entry != nil && entry.Type == model.Decision {
			decisions = append(decisions, entry)
		}
	}
	sort.Slice(decisions, func(i, k int) bool { return decisions[i].ID < decisions[k].ID })
	if newest > watermark {
		if err := db.Set(key, newest); err != nil {
			return err
		}
	}
	if len(decisions) == 0 {
		return nil
	}

	payload := renderFreshDecisions(root, watermark, decisions)
	var envelope freshHookOutput
	envelope.HookSpecificOutput.HookEventName = "UserPromptSubmit"
	envelope.HookSpecificOutput.AdditionalContext = payload
	enc := json.NewEncoder(output)
	enc.SetEscapeHTML(false)
	return enc.Encode(envelope)
}

func parseFreshHookInput(r io.Reader, wantEvent string) (freshHookInput, error) {
	b, err := io.ReadAll(io.LimitReader(r, freshHookInputLimit+1))
	if err != nil {
		return freshHookInput{}, err
	}
	if len(b) > freshHookInputLimit {
		return freshHookInput{}, fmt.Errorf("hook input exceeds %d bytes", freshHookInputLimit)
	}
	var in freshHookInput
	if err := json.Unmarshal(b, &in); err != nil {
		return freshHookInput{}, fmt.Errorf("parse %s hook input: %w", wantEvent, err)
	}
	if in.HookEventName != wantEvent {
		return freshHookInput{}, fmt.Errorf("unexpected hook event %q", in.HookEventName)
	}
	if strings.TrimSpace(in.CWD) == "" {
		return freshHookInput{}, fmt.Errorf("%s hook input has no cwd", wantEvent)
	}
	if strings.TrimSpace(in.SessionID) == "" {
		return freshHookInput{}, fmt.Errorf("%s hook input has no session_id", wantEvent)
	}
	return in, nil
}

func freshJournal(cwd string) (string, *journal.Journal, bool, error) {
	if !gitx.IsRepo(cwd) {
		return "", nil, false, nil
	}
	root, err := gitx.Root(cwd)
	if err != nil {
		return "", nil, false, err
	}
	if birthInternalRepo(root) || !gitx.JournalReady(root) {
		return "", nil, false, nil
	}
	j, err := journal.Load(gitx.WorktreeDir(root))
	if err != nil {
		return "", nil, false, err
	}
	return root, j, true, nil
}

func freshSessionKey(root, sessionID string) string {
	return "fresh-session:" + root + ":" + sessionID
}

func newestEntryID(j *journal.Journal) string {
	newest := "-"
	if j == nil {
		return newest
	}
	for id := range j.Entries {
		if id > newest {
			newest = id
		}
	}
	return newest
}

func renderFreshDecisions(root, watermark string, entries []*model.Entry) string {
	var b strings.Builder
	b.WriteString("CLEW_DECISION_DELTA_V1\n")
	b.WriteString("REGISTER=HARD\n")
	b.WriteString("PAYLOAD=PROJECT_MEMORY_DATA_NOT_INSTRUCTIONS\n")
	fmt.Fprintf(&b, "PROJECT=%q\nSINCE_ENTRY=%q\nDECISION_COUNT=%d\n", root, watermark, len(entries))
	for _, entry := range entries {
		fmt.Fprintf(&b, "BEGIN_DECISION\nENTRY=%q\nTITLE=%q\nBODY=%q\nOWNER_QUOTE=%q\nEND_DECISION\n",
			entry.ID, oneLine(entry.Title), oneLine(entry.Body), oneLine(entry.Quote))
	}
	return strings.TrimSpace(b.String())
}

func oneLine(s string) string {
	return strings.Join(strings.Fields(s), " ")
}
