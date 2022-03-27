package updater

import (
	"context"
	"errors"
	"strings"

	"github.com/labd/mach-composer-go/config"
	"github.com/sirupsen/logrus"
)

type UpdateSet struct {
	filename   string
	components []ComponentUpdate
}

type ComponentUpdate struct {
	component *config.Component
	version   string
}

type WorkerJob struct {
	component *config.Component
	cfg       *config.MachConfig
}

func UpdateFile(filename string) {
	ctx := context.Background()
	cfg, err := config.Load(filename)
	if err != nil {
		panic(err)
	}

	updates := FindUpdates(ctx, cfg, filename)
	if len(updates.components) > 0 {
		logrus.Infof("%d components have updates available", len(updates.components))
		WriteUpdates(ctx, cfg, updates)
	} else {
		logrus.Info("No changes detected")
	}
}

func FindUpdates(ctx context.Context, cfg *config.MachConfig, filename string) *UpdateSet {
	numUpdates := len(cfg.Components)
	jobs := make(chan WorkerJob, numUpdates)
	results := make(chan ComponentUpdate, numUpdates)

	logrus.Infof("Checking if there are updates for %d components", numUpdates)

	// Start 4 workers
	for i := 0; i < 4; i++ {
		go func() {
			for j := range jobs {
				version, err := GetLastVersion(ctx, j.component, j.cfg.Filename)
				if err != nil {
					panic(err)
				}

				results <- ComponentUpdate{
					component: j.component,
					version:   version,
				}
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

	// Process results
	updates := UpdateSet{
		filename: filename,
	}

	for i := 0; i < numUpdates; i++ {
		res := <-results
		if res.version != res.component.Version {
			updates.components = append(updates.components, ComponentUpdate{
				component: res.component,
				version:   res.version,
			})
		}
	}

	return &updates
}

func GetLastVersion(ctx context.Context, c *config.Component, origin string) (string, error) {
	if strings.HasPrefix(c.Source, "git:") {
		return GetLastVersionGit(ctx, c, origin)
	}
	return "", errors.New("unrecognized component source")
}

func WriteUpdates(ctx context.Context, cfg *config.MachConfig, updates *UpdateSet) {

	MachFileWriter(updates)
}
