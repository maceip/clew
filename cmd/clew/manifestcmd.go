package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"clew/internal/journal"
	"clew/internal/manifest"
)

// cmdManifest: the restart kit (§9). First run generates MANIFEST.md for the
// human pass; run with --out to apply marks and emit SEED.md + genesis/.
func cmdManifest(args []string) error {
	fs := flag.NewFlagSet("manifest", flag.ExitOnError)
	spec := fs.String("spec", "", "new spec file for the disposition pass")
	out := fs.String("out", "", "output dir: apply MANIFEST.md marks → SEED.md + genesis/")
	yes := fs.Bool("yes", false, "skip the human pass: carry everything as marked by default")
	fs.Parse(args)

	a, err := load()
	if err != nil {
		return err
	}
	defer a.close()
	repo, err := a.repoFromCwd()
	if err != nil {
		return err
	}
	j, err := a.openJournal(repo)
	if err != nil {
		return err
	}
	now := time.Now()
	st := journal.Compute(j, now)
	mpath := filepath.Join(j.Dir, "MANIFEST.md")

	if *out == "" {
		p, _ := a.provider()
		if p != nil {
			p = newBudgetedProvider(p, a.db, a.cfg, "manifest", false, 0)
		}
		path, err := manifest.Generate(j, st, *spec, p, now)
		if err != nil {
			return err
		}
		fmt.Printf("wrote %s\n\nnext: edit the [carry]/[drop] marks (the human pass, §9.2),\nthen: clew manifest --out <dir>\n", path)
		return nil
	}

	if _, err := os.Stat(mpath); os.IsNotExist(err) || *yes {
		p, _ := a.provider()
		if p != nil {
			p = newBudgetedProvider(p, a.db, a.cfg, "manifest", false, 0)
		}
		if _, err := manifest.Generate(j, st, *spec, p, now); err != nil {
			return err
		}
		if !*yes {
			return fmt.Errorf("MANIFEST.md was just generated at %s — do the human pass first (or --yes)", mpath)
		}
	}
	res, err := manifest.Apply(j, st, mpath, *out, a.cfg.Surface, now)
	if err != nil {
		return err
	}
	if _, err := a.syncAndMaterialize(repo); err != nil {
		return err
	}
	fmt.Printf("carried %d, dropped %d (dispositions journaled — the loss is deliberate and dated)\n", len(res.Carried), len(res.Dropped))
	fmt.Printf("  %s   (≤4KB, paste-ready first prompt)\n  %s\nimport in the successor repo: clew init --carry %s\n",
		res.SeedPath, res.GenesisDir, res.GenesisDir)
	return nil
}
