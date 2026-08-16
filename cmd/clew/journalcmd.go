package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"

	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
	"clew/internal/scrub"
)

// cmdJournal: human edit verbs (§8.2). Editing is first-class writing (I6):
// every verb is an append-only event (or a new entry), never a mutation.
func cmdJournal(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("usage: clew journal show|edit|confirm|reject|supersede|answer|note …")
	}
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	verb, rest := args[0], args[1:]
	switch verb {
	case "show":
		if len(rest) < 1 {
			return fmt.Errorf("usage: clew journal show <entry-id>")
		}
		return journalShow(a, repo, rest[0])
	case "confirm":
		fs := flag.NewFlagSet("confirm", flag.ExitOnError)
		done := fs.Bool("done", false, "mark an intent completed (status → done)")
		contradicts := fs.String("contradicts", "", "confirm a real contradiction with this other decision id")
		if err := parseAfterID(fs, rest); err != nil {
			return err
		}
		payload := map[string]any{}
		if *done {
			payload["done"] = true
		}
		if *contradicts != "" {
			payload["contradicts"] = *contradicts
		}
		return humanEvent(a, repo, rest[0], model.EvConfirm, payload,
			"confirmed (confidence → 1.0"+ternary(*done, "; done", "")+ternary(*contradicts != "", "; contradiction ruled real", "")+")")
	case "reject":
		fs := flag.NewFlagSet("reject", flag.ExitOnError)
		reason := fs.String("reason", "", "why")
		if err := parseAfterID(fs, rest); err != nil {
			return err
		}
		return humanEvent(a, repo, rest[0], model.EvReject, map[string]any{"reason": *reason}, "rejected (superseded-by-human)")
	case "supersede":
		fs := flag.NewFlagSet("supersede", flag.ExitOnError)
		by := fs.String("by", "", "superseding entry id (required)")
		if err := parseAfterID(fs, rest); err != nil {
			return err
		}
		if *by == "" {
			return fmt.Errorf("usage: clew journal supersede <old-id> --by <new-id>")
		}
		return humanEvent(a, repo, rest[0], model.EvSupersede, map[string]any{"by": *by}, "superseded by "+*by)
	case "answer":
		fs := flag.NewFlagSet("answer", flag.ExitOnError)
		typ := fs.String("type", "finding", "entry type for the answer (decision|finding)")
		if len(rest) < 2 {
			return fmt.Errorf("usage: clew journal answer <question-id> \"text\" [--type decision]")
		}
		fs.Parse(rest[2:])
		return journalAnswer(a, repo, rest[0], rest[1], *typ)
	case "note":
		note, help, err := parseJournalNoteArgs(rest, os.Stdout)
		if err != nil || help {
			return err
		}
		return journalNote(a, repo, note.Text, note.Type, note.Title, note.Tags, note.Asks)
	case "edit":
		if len(rest) < 1 {
			return fmt.Errorf("usage: clew journal edit <entry-id>")
		}
		return journalEdit(a, repo, rest[0])
	default:
		return fmt.Errorf("unknown journal verb %q", verb)
	}
}

type journalNoteArgs struct {
	Text, Type, Title, Tags, Asks string
}

func parseJournalNoteArgs(args []string, out io.Writer) (journalNoteArgs, bool, error) {
	fs := flag.NewFlagSet("note", flag.ContinueOnError)
	fs.SetOutput(out)
	typ := fs.String("type", "finding", "decision|finding|question|intent")
	title := fs.String("title", "", "title (default: first 80 chars of text)")
	tags := fs.String("tags", "", "comma-separated path globs")
	asks := fs.String("asks", "any", "questions only: human|any")
	usage := func() {
		fmt.Fprintln(out, "usage: clew journal note \"text\" [--type …] [--title …] [--tags a/**,b]")
		fs.PrintDefaults()
	}
	if len(args) == 0 {
		return journalNoteArgs{}, false, fmt.Errorf("usage: clew journal note \"text\" [--type …] [--title …] [--tags a/**,b]")
	}
	if args[0] == "-h" || args[0] == "--help" {
		usage()
		return journalNoteArgs{}, true, nil
	}
	if err := fs.Parse(args[1:]); err != nil {
		if err == flag.ErrHelp {
			return journalNoteArgs{}, true, nil
		}
		return journalNoteArgs{}, false, err
	}
	if fs.NArg() != 0 {
		return journalNoteArgs{}, false, fmt.Errorf("unexpected note arguments: %s", strings.Join(fs.Args(), " "))
	}
	return journalNoteArgs{Text: args[0], Type: *typ, Title: *title, Tags: *tags, Asks: *asks}, false, nil
}

func parseAfterID(fs *flag.FlagSet, rest []string) error {
	if len(rest) < 1 || strings.HasPrefix(rest[0], "-") {
		return fmt.Errorf("first argument must be an entry id")
	}
	return fs.Parse(rest[1:])
}

func journalShow(a *app, repo, id string) error {
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	e, ok := j.Entries[id]
	if !ok {
		return fmt.Errorf("no entry %s", id)
	}
	st := journal.Compute(j, time.Now())[id]
	b, _ := yaml.Marshal(e)
	fmt.Print(string(b))
	fmt.Printf("--- computed (never persisted, §3.2) ---\nstatus: %s\nconfidence: %.2f\nevidence: %d\nlast_activity: %s\n",
		st.Status, st.Confidence, st.Evidence, st.LastActivity.Format(time.RFC3339))
	if len(st.Contradicts) > 0 {
		fmt.Printf("pairs-with: %s\n", strings.Join(st.Contradicts, ", "))
	}
	if st.Tainted {
		fmt.Println("taint: quote originated in tool_result (§6.5)")
	}
	if st.Withheld {
		fmt.Println("withheld from context.md pending confirm (imperative pattern, §6.5)")
	}
	evs := j.EventsFor(id)
	if len(evs) > 0 {
		fmt.Println("--- events ---")
		for _, v := range evs {
			pb, _ := yaml.Marshal(v.Payload)
			fmt.Printf("%s %s by %s@%s at %s: %s", v.ID, v.Kind, v.By.Who, v.By.Surface,
				v.At.Format("2006-01-02 15:04"), strings.ReplaceAll(string(pb), "\n", " "))
			fmt.Println()
		}
	}
	fmt.Printf("--- source jump ---\n%s\n", e.Source.Ref)
	return nil
}

func humanEvent(a *app, repo, id string, kind model.EventKind, payload map[string]any, what string) error {
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	if _, ok := j.Entries[id]; !ok {
		return fmt.Errorf("no entry %s", id)
	}
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(time.Now()), Kind: kind, Entry: id, Payload: payload,
		By: model.By{Who: "human", Surface: a.cfg.Surface}, At: time.Now().UTC(),
	}); err != nil {
		return err
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("%s: %s\n", id, what)
	return nil
}

func journalAnswer(a *app, repo, qid, text, typ string) error {
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	q, ok := j.Entries[qid]
	if !ok {
		return fmt.Errorf("no entry %s", qid)
	}
	if q.Type != model.Question {
		return fmt.Errorf("%s is a %s, not a question", qid, q.Type)
	}
	et := model.Finding
	if typ == "decision" {
		et = model.Decision
	}
	now := time.Now().UTC()
	text, _ = scrub.Scrub(text)
	ans := &model.Entry{
		ID: ids.NewEntry(now), Type: et,
		Title: clipStr("answer: "+q.Title, 80), Body: clipStr(text, 400), Quote: text,
		UtteranceBy: model.ByUser,
		Source: model.Source{Kind: model.SrcHuman, Ref: "docket:answer:" + qid,
			Surface: a.cfg.Surface, At: now},
		Confidence: 1.0, Tags: q.Tags,
	}
	if err := j.AddEntry(ans); err != nil {
		return err
	}
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(now), Kind: model.EvAnswer, Entry: qid,
		Payload: map[string]any{"by": ans.ID},
		By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: now,
	}); err != nil {
		return err
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("answered %s with %s — echoes to agents via context.md\n", qid, ans.ID)
	return nil
}

func journalNote(a *app, repo, text, typ, title, tags, asks string) error {
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	var et model.EntryType
	switch typ {
	case "decision":
		et = model.Decision
	case "finding":
		et = model.Finding
	case "question":
		et = model.Question
	case "intent":
		et = model.Intent
	default:
		return fmt.Errorf("--type must be decision|finding|question|intent")
	}
	if title == "" {
		title = clipStr(text, 80)
	}
	now := time.Now().UTC()
	text, _ = scrub.Scrub(text)
	e := &model.Entry{
		ID: ids.NewEntry(now), Type: et, Title: clipStr(title, 80),
		Body: clipStr(text, 400), Quote: text, UtteranceBy: model.ByUser,
		Source:     model.Source{Kind: model.SrcHuman, Ref: "cli:note", Surface: a.cfg.Surface, At: now},
		Confidence: 1.0,
	}
	if tags != "" {
		e.Tags = strings.Split(tags, ",")
	}
	if et == model.Question {
		e.Asks = asks
	}
	if err := j.AddEntry(e); err != nil {
		return err
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Println("noted", e.ID)
	return nil
}

// journalEdit: entries are immutable — editing opens a copy in $EDITOR and
// writes a NEW human entry superseding the old one (append-only law, §3.2).
func journalEdit(a *app, repo, id string) error {
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	old, ok := j.Entries[id]
	if !ok {
		return fmt.Errorf("no entry %s", id)
	}
	edit := *old
	now := time.Now().UTC()
	edit.ID = ids.NewEntry(now)
	edit.Source = model.Source{Kind: model.SrcHuman, Ref: "edit:" + id, Surface: a.cfg.Surface, At: now}
	edit.Confidence = 1.0
	b, _ := yaml.Marshal(&edit)
	tmp := filepath.Join(os.TempDir(), edit.ID+".yaml")
	header := "# Editing a COPY of " + id + ". Saving writes a new entry that supersedes it.\n"
	if err := os.WriteFile(tmp, append([]byte(header), b...), 0o600); err != nil {
		return err
	}
	defer os.Remove(tmp)
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi"
	}
	cmd := exec.Command(editor, tmp)
	cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}
	nb, err := os.ReadFile(tmp)
	if err != nil {
		return err
	}
	ne := &model.Entry{}
	if err := yaml.Unmarshal(nb, ne); err != nil {
		return fmt.Errorf("edited YAML invalid: %w", err)
	}
	ne.ID = edit.ID // id is not editable
	if err := j.AddEntry(ne); err != nil {
		return err
	}
	if err := j.AddEvent(&model.Event{
		ID: ids.NewEvent(now), Kind: model.EvSupersede, Entry: id,
		Payload: map[string]any{"by": ne.ID, "why": "human edit"},
		By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: now,
	}); err != nil {
		return err
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("wrote %s superseding %s\n", ne.ID, id)
	return nil
}

func ternary(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}
