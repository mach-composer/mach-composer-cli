package updater

import (
	"context"
	"errors"
	"strings"

	"github.com/labd/mach-composer-go/config"
	"github.com/sirupsen/logrus"
)

type UpdateSet struct {
	filename string
	updates  []ChangeSet
}

type WorkerJob struct {
	component *config.Component
	cfg       *config.MachConfig
}

func UpdateFile(filename string) {
	ctx := context.Background()
	cfg, err := config.Load(filename, "")
	if err != nil {
		panic(err)
	}

	updateSet := FindUpdates(ctx, cfg, filename)
	if len(updateSet.updates) > 0 {
		logrus.Infof("%d components have updates available", len(updateSet.updates))
		WriteUpdates(ctx, cfg, updateSet)
	} else {
		logrus.Info("No changes detected")
	}
}

func FindUpdates(ctx context.Context, cfg *config.MachConfig, filename string) *UpdateSet {
	numUpdates := len(cfg.Components)
	jobs := make(chan WorkerJob, numUpdates)
	results := make(chan *ChangeSet, numUpdates)

	logrus.Infof("Checking if there are updates for %d components", numUpdates)

	// Start 4 workers
	for i := 0; i < 4; i++ {
		go func() {
			for j := range jobs {
				cs, err := GetLastVersion(ctx, j.component, j.cfg.Filename)
				if err != nil {
					panic(err)
				}

				results <- cs
			}
		}()
	}

	// Send work
	for i := range cfg.Components {
		jobs <- WorkerJob{
			component: &cfg.Components[i],
			cfg:       cfg,
		}
	}
	close(jobs)

	// Process results as we receive them from the channel
	updates := UpdateSet{
		filename: filename,
	}
	for i := 0; i < numUpdates; i++ {
		changeset := <-results

		OutputChanges(changeset)

		if changeset.HasChanges() {
			updates.updates = append(updates.updates, *changeset)
		}
	}

	return &updates
}

func GetLastVersion(ctx context.Context, c *config.Component, origin string) (*ChangeSet, error) {
	if strings.HasPrefix(c.Source, "git:") {
		return GetLastVersionGit(ctx, c, origin)
	}
	return nil, errors.New("unrecognized component source")
}

func WriteUpdates(ctx context.Context, cfg *config.MachConfig, updates *UpdateSet) {

	MachFileWriter(updates)
}
