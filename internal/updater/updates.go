package updater

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"math"
	"runtime"
	"sync"

	"github.com/mach-composer/mach-composer-cli/internal/config"

	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
	"golang.org/x/sync/semaphore"
)

func findUpdates(ctx context.Context, cfg *PartialConfig, filename string) (*UpdateSet, error) {
	log.Ctx(ctx).Info().Msgf("Checking if there are updates for %d components\n", len(cfg.Components))
	if cfg.client == nil {
		return findUpdatesParallel(ctx, cfg, filename)
	}
	return findUpdatesSerial(ctx, cfg, filename)
}

func findUpdatesSerial(ctx context.Context, cfg *PartialConfig, filename string) (*UpdateSet, error) {
	updates := UpdateSet{
		filename: filename,
	}

	for i := range cfg.Components {
		cs, err := getLastVersion(ctx, cfg, &cfg.Components[i], cfg.filename)
		if err != nil {
			return nil, err
		}

		if cs == nil {
			continue
		}

		output := OutputChanges(cs)
		log.Ctx(ctx).Info().Msg(output)

		if cs.HasChanges() {
			updates.updates = append(updates.updates, *cs)
		}
	}
	return &updates, nil
}

func determineNumWorkers() int {
	var numWorkers = int(math.Ceil(float64(runtime.NumCPU() / 2)))

	if numWorkers < 2 {
		numWorkers = 2
	}

	return numWorkers
}

func findUpdatesParallel(ctx context.Context, cfg *PartialConfig, filename string) (*UpdateSet, error) {
	numUpdates := len(cfg.Components)
	resChan := make(chan *ChangeSet, numUpdates)
	errChan := make(chan error, numUpdates)

	var numWorkers = determineNumWorkers()
	var sem = semaphore.NewWeighted(int64(numWorkers))
	var wg sync.WaitGroup

	log.Info().Msgf("Running on %d workers", numWorkers)

	// Compute the output using up to maxWorkers goroutines at a time.
	for _, c := range cfg.Components {
		// When maxWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			log.Printf("Failed to acquire semaphore: %v", err)
			break
		}

		wg.Add(1)

		go func(c config.ComponentConfig) {
			defer sem.Release(1)
			defer wg.Done()

			logger := log.With().Str("component", c.Name).Logger()
			ctx = logger.WithContext(ctx)

			cs, err := getLastVersion(ctx, cfg, &c, cfg.filename)
			if err != nil {
				logger.Error().Msg(err.Error())
				errChan <- err
				return
			}

			if cs == nil {
				return
			}

			resChan <- cs
			return
		}(c)
	}

	wg.Wait()
	close(errChan)
	close(resChan)

	if n := len(errChan); n > 0 {
		return nil, fmt.Errorf("failed to update %d components", n)
	}

	// Process results as we receive them from the channel
	updates := UpdateSet{
		filename: filename,
	}
	for changeSet := range resChan {
		if changeSet == nil {
			continue
		}

		output := OutputChanges(changeSet)
		log.Ctx(ctx).Info().Msg(output)

		if changeSet.HasChanges() {
			updates.updates = append(updates.updates, *changeSet)
		}
	}

	return &updates, nil
}

func findSpecificUpdate(ctx context.Context, cfg *PartialConfig, filename string, component *config.ComponentConfig) (*UpdateSet, error) {
	changeSet, err := getLastVersion(ctx, cfg, component, filename)
	if err != nil {
		return nil, err
	}

	output := OutputChanges(changeSet)
	log.Ctx(ctx).Info().Msg(output)

	updates := UpdateSet{
		filename: cfg.filename,
		updates:  []ChangeSet{*changeSet},
	}
	return &updates, nil
}

func getLastVersion(ctx context.Context, cfg *PartialConfig, c *config.ComponentConfig, origin string) (*ChangeSet, error) {
	if c.Branch == "" {
		c.Branch = "main"
	}

	if cfg.client != nil {
		return getLastVersionCloud(ctx, cfg, c, origin)
	}

	if c.Source.IsType(config.SourceTypeGit) {
		return getLastVersionGit(ctx, c, origin)
	}

	err := &UpdateError{
		msg:       fmt.Sprintf("unrecognized component source for %s: %s", c.Name, c.Source),
		component: c.Name,
		source:    string(c.Source),
	}
	return nil, err
}

func getLastVersionCloud(ctx context.Context, cfg *PartialConfig, c *config.ComponentConfig, origin string) (*ChangeSet, error) {
	organization := cfg.MachComposer.Cloud.Organization
	project := cfg.MachComposer.Cloud.Project

	version, _, err := cfg.client.
		ComponentsApi.
		ComponentLatestVersion(ctx, organization, project, c.Name).
		Branch(c.Branch).
		Execute()

	if err != nil {
		if c.Source.IsType(config.SourceTypeGit) {
			log.Ctx(ctx).Warn().Msgf("Error checking for %s in MACH Composer Cloud, falling back to Git", c.Name)
			return getLastVersionGit(ctx, c, origin)
		}
		log.Ctx(ctx).Error().Err(err).Msgf("Error checking for latest version of %s", c.Name)
		return nil, nil
	}

	if version == nil {
		if c.Source.IsType(config.SourceTypeGit) {
			log.Ctx(ctx).Warn().Msgf("No version found for %s in MACH Composer Cloud, falling back to Git", c.Name)
			return getLastVersionGit(ctx, c, origin)
		}
		log.Ctx(ctx).Warn().Msgf("No version found for %s", c.Name)
		return nil, nil
	}

	cs := &ChangeSet{
		Changes:     []CommitData{},
		Component:   c,
		LastVersion: version.Version,
	}

	if c.Version != version.Version {
		paginator, _, err := cfg.client.
			ComponentsApi.
			ComponentCommitQuery(ctx, organization, project, c.Name).
			From(c.Version).
			To(version.Version).
			Offset(0).
			Limit(200).
			Execute()
		if err != nil {
			return nil, err
		}

		for _, record := range paginator.Results {
			change := CommitData{
				Commit:  record.Commit,
				Parents: record.Parents,
				Message: record.Subject,
				Author: CommitAuthor{
					Email: record.Author.Email,
					Name:  record.Author.Name,
					Date:  record.Author.Date,
				},
				Committer: CommitAuthor{
					Email: record.Committer.Email,
					Name:  record.Committer.Name,
					Date:  record.Committer.Date,
				},
			}
			cs.Changes = append(cs.Changes, change)
		}
	}

	return cs, nil
}

func getLastVersionGit(ctx context.Context, c *config.ComponentConfig, origin string) (*ChangeSet, error) {
	commits, err := gitutils.GetLastVersionGit(ctx, c, origin)
	if err != nil {
		return nil, err
	}

	cd := make([]CommitData, len(commits))
	for i := range commits {
		c := commits[i]

		cd[i].Commit = c.Commit
		cd[i].Parents = c.Parents
		cd[i].Message = c.Message

		cd[i].Author = CommitAuthor{
			Email: c.Author.Email,
			Name:  c.Author.Name,
			Date:  c.Author.Date,
		}
		cd[i].Committer = CommitAuthor{
			Email: c.Committer.Email,
			Name:  c.Committer.Name,
			Date:  c.Committer.Date,
		}
		cd[i].Tags = c.Tags
	}

	cs := &ChangeSet{
		Changes:   cd,
		Component: c,
	}

	if len(commits) < 1 {
		cs.LastVersion = c.Version
	} else {
		cs.LastVersion = commits[0].Commit
	}

	return cs, nil
}
