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
			if err := a.db.Set(mergeStateKey(repo, item.ID), "applied"); err != nil {
				return err
			}
			ids = append(ids, item.ID)
		}
		return checkin.RenderAgentHandoff(os.Stdout, "apply", ids)
	case "apply", "explain", "defer":
		if len(args) != 2 {
			return fmt.Errorf("usage: clew merge %s <entry-id>", verb)
		}
		if !viewHas(view, args[1]) {
			return fmt.Errorf("%s is not on the current knowledge merge", args[1])
		}
		switch verb {
		case "apply":
			if err := a.db.Set(mergeStateKey(repo, args[1]), "applied"); err != nil {
				return err
			}
			return checkin.RenderAgentHandoff(os.Stdout, "apply", []string{args[1]})
		case "explain":
			return checkin.RenderAgentHandoff(os.Stdout, "explain", []string{args[1]})
		case "defer":
			if err := a.db.Set(mergeStateKey(repo, args[1]), "deferred"); err != nil {
				return err
			}
			fmt.Println("Deferred. It will remain in the count below the next list.")
			return nil
		}
	default:
		return fmt.Errorf("usage: clew merge [apply|explain|defer <entry-id>|apply-all]")
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
	if len(args) != 2 {
		return fmt.Errorf("usage: clew gap build|explain|retire <entry-id>")
	}
	verb, id := args[0], args[1]
	if !viewHas(view, id) {
		return fmt.Errorf("%s is not on the current intent gap", id)
	}
	switch verb {
	case "build", "explain":
		return checkin.RenderAgentHandoff(os.Stdout, verb, []string{id})
	case "retire":
		if err := j.AddEvent(&model.Event{
			ID: ids.NewEvent(time.Now()), Kind: model.EvReject, Entry: id,
			Payload: map[string]any{"reason": "retired through intent gap"},
			By:      model.By{Who: "human", Surface: a.cfg.Surface}, At: time.Now().UTC(),
		}); err != nil {
			return err
		}
		if _, err := a.syncAndMaterialize(repo); err != nil {
			return err
		}
		fmt.Printf("Retired %s. Its source and this decision remain available.\n", id)
		return nil
	default:
		return fmt.Errorf("usage: clew gap build|explain|retire <entry-id>")
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
	view.Checks = checkinFailureCount(a, repo)
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
	view.Checks = checkinFailureCount(a, repo)
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

func viewHas(view checkin.View, id string) bool {
	for _, item := range view.Items {
		if item.ID == id {
			return true
		}
	}
	return false
}

func checkinFailureCount(a *app, repo string) int {
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
	for _, key := range keys {
		if a.db.Get(key) != "" {
			n++
		}
	}
	for _, prefix := range []string{"extract-error:", "adapter-error:", "adapter-paused:", "birth-error:", "birth-hook-error:"} {
		n += len(a.db.KVPrefix(prefix))
	}
	return n
}
