package updater

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/labd/mach-composer-go/config"
	"github.com/sirupsen/logrus"
)

type WorkerJob struct {
	component *config.Component
	cfg       *config.MachConfig
}

func UpdateFile(filename, componentName, componentVersion string) *UpdateSet {
	ctx := context.Background()
	cfg, err := config.Load(filename, "")
	if err != nil {
		panic(err)
	}

	// Find the component if defined in the config
	var component *config.Component
	if componentName != "" {
		for i, c := range cfg.Components {
			if strings.EqualFold(c.Name, componentName) {
				component = &cfg.Components[i]
				break
			}
		}

		if component == nil {
			logrus.Infof("No component found with name %s", componentName)
			return nil
		}
	}

	// Check for updates, we have three options:
	//  1. no component is defined (and thus no version). Update all components
	//  2. component is defined but no version. Get last version for component
	//  3. both component and version are defined. Just set it manually.
	var updateSet *UpdateSet
	if component != nil {

		// If no specific version is defined we auto-detect the last version
		if componentVersion != "" {
			updateSet = &UpdateSet{
				filename: cfg.Filename,
				updates: []ChangeSet{
					{
						LastVersion: componentVersion,
						Component:   component,
						Forced:      true,
					},
				},
			}
			logrus.Infof("Setting component %s to version %s", component.Name, componentVersion)
		} else {
			updateSet = FindSpecificUpdate(ctx, cfg, filename, component)
			if updateSet.HasChanges() {
				logrus.Infof("Updating component %s to version %s", component.Name, updateSet.updates[0].LastVersion)
			} else {
				logrus.Infof("No updates for component %s", component.Name)
			}
		}
	} else {
		updateSet = FindUpdates(ctx, cfg, filename, component)
		logrus.Infof("%d components have updates available", len(updateSet.updates))
	}

	if len(updateSet.updates) > 0 {
		WriteUpdates(ctx, cfg, updateSet)
	} else {
		logrus.Info("No changes detected")
	}

	return updateSet
}

func FindUpdates(ctx context.Context, cfg *config.MachConfig, filename string, component *config.Component) *UpdateSet {
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

		output := OutputChanges(changeset)
		fmt.Print(output)

		if changeset.HasChanges() {
			updates.updates = append(updates.updates, *changeset)
		}
	}

	return &updates
}

func FindSpecificUpdate(ctx context.Context, cfg *config.MachConfig, filename string, component *config.Component) *UpdateSet {
	changeset, err := GetLastVersion(ctx, component, cfg.Filename)
	if err != nil {
		panic(err)
	}

	output := OutputChanges(changeset)
	fmt.Print(output)

	updates := UpdateSet{
		filename: cfg.Filename,
		updates:  []ChangeSet{*changeset},
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
