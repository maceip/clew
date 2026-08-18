package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"clew/internal/checkin"
	"clew/internal/gitx"
	"clew/internal/ids"
	"clew/internal/journal"
	"clew/internal/model"
)

func cmdMerge(args []string) error {
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, _, view, err := knowledgeMergeView(a)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return checkin.Render(os.Stdout, view)
	}
	if len(args) == 2 && args[0] == "apply" && args[1] == "all" {
		args = []string{"apply-all"}
	}
	verb := args[0]
	switch verb {
	case "apply-all":
		if len(args) != 1 {
			return fmt.Errorf("usage: clew merge apply-all")
		}
		if len(view.Items) == 0 {
			return fmt.Errorf("nothing new to apply")
		}
		ids := make([]string, 0, len(view.Items))
		for _, item := range view.Items {
			for _, id := range checkin.EntryIDs(item) {
				if err := a.db.Set(mergeStateKey(repo, id), "applied"); err != nil {
					return err
				}
				ids = append(ids, id)
			}
		}
		return checkin.RenderAgentHandoff(os.Stdout, "apply", ids)
	case "apply", "explain", "defer":
		if len(args) < 2 {
			return fmt.Errorf("name the change to %s", verb)
		}
		item, err := checkin.Resolve(view, strings.Join(args[1:], " "))
		if err != nil {
			return err
		}
		ids := checkin.EntryIDs(item)
		switch verb {
		case "apply":
			for _, id := range ids {
				if err := a.db.Set(mergeStateKey(repo, id), "applied"); err != nil {
					return err
				}
			}
			return checkin.RenderAgentHandoff(os.Stdout, "apply", ids)
		case "explain":
			return checkin.RenderAgentHandoff(os.Stdout, "explain", ids)
		case "defer":
			for _, id := range ids {
				if err := a.db.Set(mergeStateKey(repo, id), "deferred"); err != nil {
					return err
				}
			}
			fmt.Println("That change is deferred and remains in the count below the next list.")
			return nil
		}
	default:
		return fmt.Errorf("say apply, explain, or defer, then name the change")
	}
	return nil
}

func cmdGap(args []string) error {
	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, j, view, err := intentGapView(a)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		return checkin.Render(os.Stdout, view)
	}
	if len(args) < 2 {
		return fmt.Errorf("say build, explain, or retire, then name the change")
	}
	verb := args[0]
	item, err := checkin.Resolve(view, strings.Join(args[1:], " "))
	if err != nil {
		return err
	}
	entryIDs := checkin.EntryIDs(item)
	switch verb {
	case "build", "explain":
		return checkin.RenderAgentHandoff(os.Stdout, verb, entryIDs)
	case "retire":
		for _, id := range entryIDs {
			now := time.Now()
			if err := j.AddEvent(&model.Event{
				ID: ids.NewEvent(now), Kind: model.EvReject, Entry: id,
				Payload: map[string]any{"reason": "retired through intent gap"},
				By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: now.UTC(),
			}); err != nil {
				return err
			}
		}
		if _, err := a.syncAndMaterialize(repo); err != nil {
			return err
		}
		fmt.Println("That intent is retired; its source and this decision remain available.")
		return nil
	default:
		return fmt.Errorf("say build, explain, or retire, then name the change")
	}
}

func knowledgeMergeView(a *app) (string, *journal.Journal, checkin.View, error) {
	repo, err := a.repoFromCwd()
	if err != nil {
		return "", nil, checkin.View{}, err
	}
	j, err := journal.LoadForDisplay(gitx.WorktreeDir(repo))
	if err != nil {
		return "", nil, checkin.View{}, err
	}
	view := checkin.BuildMerge(j, mergeStates(a, repo))
	view.Repo = repoBase(repo)
	view.Repairs = checkinFailureRepairs(a, repo)
	return repo, j, view, nil
}

func intentGapView(a *app) (string, *journal.Journal, checkin.View, error) {
	repo, err := a.repoFromCwd()
	if err != nil {
		return "", nil, checkin.View{}, err
	}
	j, err := journal.LoadForDisplay(gitx.WorktreeDir(repo))
	if err != nil {
		return "", nil, checkin.View{}, err
	}
	view := checkin.BuildGap(j, a.db.OpenAlerts(repo, false))
	view.Repo = repoBase(repo)
	view.Repairs = checkinFailureRepairs(a, repo)
	return repo, j, view, nil
}

func mergeStateKey(repo, id string) string { return "checkin:merge:" + repo + ":" + id }

func mergeStates(a *app, repo string) map[string]string {
	prefix := "checkin:merge:" + repo + ":"
	out := map[string]string{}
	for _, pair := range a.db.KVPrefix(prefix) {
		out[strings.TrimPrefix(pair.Key, prefix)] = pair.Value
	}
	return out
}

func checkinFailureRepairs(a *app, repo string) []string {
	keys := []string{
		"sync-error:" + repo,
		"differ-error:" + repo,
		"materialize-error:" + repo,
		"extract-paused",
		"differ-paused",
		"watch-provider-error",
		"owner-sync-error",
	}
	n := 0
	spendFloor := false
	for _, key := range keys {
		value := a.db.Get(key)
		if value != "" {
			n++
			if key == "extract-paused" && strings.Contains(strings.ToLower(value), "budget") {
				spendFloor = true
			}
		}
	}
	for _, prefix := range []string{"extract-error:", "adapter-error:", "adapter-paused:", "birth-error:", "birth-hook-error:"} {
		for _, pair := range a.db.KVPrefix(prefix) {
			if pair.Value != "" {
				n++
			}
		}
	}
	var repairs []string
	if spendFloor {
		repairs = append(repairs, "Listening is paused until the spend floor is built — build")
		n--
	}
	if n == 1 {
		repairs = append(repairs, "The attending agent must fix one live check")
	} else if n > 1 {
		repairs = append(repairs, fmt.Sprintf("The attending agent must fix %d live checks", n))
	}
	return repairs
}
