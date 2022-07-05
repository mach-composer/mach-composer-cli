package updater

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/labd/mach-composer/internal/config"
)

type WorkerJob struct {
	component *config.Component
	cfg       *config.MachConfig
}

func UpdateFile(ctx context.Context, filename, componentName, componentVersion string, writeChanges bool) (*UpdateSet, error) {
	cfg, err := config.Load(filename, "")
	if err != nil {
		return nil, err
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
			fmt.Fprintf(os.Stderr, "No component found with name %s", componentName)
			return nil, nil
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
			fmt.Printf("Setting component %s to version %s\n", component.Name, componentVersion)
		} else {
			updateSet = FindSpecificUpdate(ctx, cfg, filename, component)
			if updateSet.HasChanges() {
				fmt.Printf("Updating component %s to version %s\n", component.Name, updateSet.updates[0].LastVersion)
			} else {
				fmt.Printf("No updates for component %s\n", component.Name)
			}
		}
	} else {
		updateSet = FindUpdates(ctx, cfg, filename, component)
		fmt.Printf("%d components have updates available\n", len(updateSet.updates))
	}

	if writeChanges && len(updateSet.updates) > 0 {
		WriteUpdates(ctx, cfg, updateSet)
	}

	return updateSet, nil
}

func FindUpdates(ctx context.Context, cfg *config.MachConfig, filename string, component *config.Component) *UpdateSet {
	numUpdates := len(cfg.Components)
	jobs := make(chan WorkerJob, numUpdates)
	results := make(chan *ChangeSet, numUpdates)

	fmt.Printf("Checking if there are updates for %d components\n", numUpdates)

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
		changeSet := <-results

		output := OutputChanges(changeSet)
		fmt.Print(output)

		if changeSet.HasChanges() {
			updates.updates = append(updates.updates, *changeSet)
		}
	}

	return &updates
}

func FindSpecificUpdate(ctx context.Context, cfg *config.MachConfig, filename string, component *config.Component) *UpdateSet {
	changeSet, err := GetLastVersion(ctx, component, cfg.Filename)
	if err != nil {
		panic(err)
	}

	output := OutputChanges(changeSet)
	fmt.Print(output)

	updates := UpdateSet{
		filename: cfg.Filename,
		updates:  []ChangeSet{*changeSet},
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
	if cfg.IsEncrypted {
		SopsFileWriter(cfg, updates)
	} else {
		MachFileWriter(updates)
	}
}
