package updater

import (
	"context"
	"errors"
	"fmt"
	"runtime"
	"sync"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"

	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
	"golang.org/x/sync/semaphore"
)

func determineNumWorkers() int {
	return max(2, runtime.NumCPU()/2)
}

// findUpdates checks every component for a newer version. Both the MACH
// composer Cloud and the git lookup are a single independent check per
// component, so they share the same worker pool.
func findUpdates(ctx context.Context, cfg *PartialConfig) ([]ChangeSet, error) {
	log.Ctx(ctx).Info().Msgf("Checking if there are updates for %d components\n", len(cfg.Components))

	numUpdates := len(cfg.Components)
	resChan := make(chan *ChangeSet, numUpdates)
	errChan := make(chan error, numUpdates)

	var numWorkers = determineNumWorkers()
	var sem = semaphore.NewWeighted(int64(numWorkers))
	var wg sync.WaitGroup

	log.Info().Msgf("Running on %d workers", numWorkers)

	// Compute the output using up to numWorkers goroutines at a time.
	var acquireErr error
	for i := range cfg.Components {
		// When numWorkers goroutines are in flight, Acquire blocks until one of the
		// workers finishes.
		if err := sem.Acquire(ctx, 1); err != nil {
			acquireErr = fmt.Errorf("failed to check all components for updates: %w", err)
			break
		}

		wg.Add(1)

		// Each goroutine gets its own component, so they never write to the same
		// element of cfg.Components.
		go func(c *config.ComponentConfig) {
			defer sem.Release(1)
			defer wg.Done()

			logger := log.With().Str("component", c.Name).Logger()
			cctx := logger.WithContext(ctx)

			cs, err := getLastVersion(cctx, cfg, c)
			if err != nil {
				logger.Error().Msg(err.Error())
				errChan <- err
				return
			}

			if cs == nil {
				return
			}

			resChan <- cs
		}(&cfg.Components[i])
	}

	wg.Wait()
	close(errChan)
	close(resChan)

	// A failed Acquire means the context was cancelled, so the remaining
	// components were never checked. Report that instead of returning a partial
	// result that reads as complete.
	if acquireErr != nil {
		return nil, acquireErr
	}

	if n := len(errChan); n > 0 {
		return nil, fmt.Errorf("failed to update %d components", n)
	}

	// Process results as we receive them from the channel
	var updates []ChangeSet
	for changeSet := range resChan {
		if changeSet == nil {
			continue
		}

		output := OutputChanges(changeSet)
		log.Ctx(ctx).Info().Msg(output)

		if changeSet.HasChanges() {
			updates = append(updates, *changeSet)
		}
	}

	return updates, nil
}

func findSpecificUpdate(ctx context.Context, cfg *PartialConfig, component *config.ComponentConfig) (*ChangeSet, error) {
	changeSet, err := getLastVersion(ctx, cfg, component)
	if err != nil {
		return nil, err
	}

	if changeSet == nil {
		return nil, nil
	}

	output := OutputChanges(changeSet)
	log.Ctx(ctx).Info().Msg(output)

	return changeSet, nil
}

func getLastVersion(ctx context.Context, cfg *PartialConfig, c *config.ComponentConfig) (*ChangeSet, error) {
	if c.Branch == "" {
		c.Branch = "main"
	}

	if cfg.client != nil {
		return getLastVersionCloud(ctx, cfg, c)
	}

	if c.Source.IsType(config.SourceTypeGit) {
		return getLastVersionGit(ctx, c, cfg.filename)
	}

	return nil, &UpdateError{
		msg: fmt.Sprintf("unrecognized component source for %s: %s", c.Name, c.Source),
	}
}

func getLastVersionCloud(ctx context.Context, cfg *PartialConfig, c *config.ComponentConfig) (*ChangeSet, error) {
	organization := cfg.MachComposer.Cloud.Organization
	project := cfg.MachComposer.Cloud.Project

	if c.Version == cloud.LatestVersion {
		log.Ctx(ctx).Warn().
			Str("msg", "When using $LATEST the difference between latest available and currently configured cannot be determined.").
			Msgf("Version set to %s for %s", cloud.LatestVersion, c.Name)
		return nil, nil
	}

	if c.Version == cloud.VersionNotApplicable {
		log.Ctx(ctx).Debug().Msgf("Version set to %s for %s. Ignoring", cloud.VersionNotApplicable, c.Name)
		return nil, nil
	}

	version, _, err := cfg.client.
		ComponentsApi.
		ComponentLatestVersion(ctx, organization, project, c.Name).
		Branch(c.Branch).
		Execute()

	if err != nil {
		if cfg.gitFallback && c.Source.IsType(config.SourceTypeGit) {
			log.Ctx(ctx).Err(err).Msgf("Error checking for %s in MACH Composer Cloud, falling back to Git", c.Name)
			return getLastVersionGit(ctx, c, cfg.filename)
		}
		log.Ctx(ctx).Error().Err(err).Msgf("Error checking for latest version of %s", c.Name)
		return nil, nil
	}

	if version == nil {
		if cfg.gitFallback && c.Source.IsType(config.SourceTypeGit) {
			log.Ctx(ctx).Warn().Msgf("No version found for %s in MACH Composer Cloud, falling back to Git", c.Name)
			return getLastVersionGit(ctx, c, cfg.filename)
		}
		log.Ctx(ctx).Warn().Msgf("No version found for %s", c.Name)
		return nil, nil
	}

	// Changes is left empty: commits are no longer registered with a version, so
	// there is nothing to list between the configured and the latest version. The
	// update itself only depends on the latest version.
	return &ChangeSet{
		Component:   c,
		LastVersion: version.Version,
	}, nil
}

func getLastVersionGit(ctx context.Context, c *config.ComponentConfig, origin string) (*ChangeSet, error) {
	commits, err := gitutils.GetLastVersionGit(ctx, c, origin)
	if err != nil {
		// The configured version is not in the repository (anymore), so the
		// changes cannot be determined. Skip the component instead of failing
		// the complete run.
		if errors.Is(err, gitutils.ErrGitRevisionNotFound) {
			log.Ctx(ctx).Warn().Msgf("Could not determine changes for %s, version %s not found in the repository",
				c.Name, c.Version)
			return nil, nil
		}
		return nil, err
	}

	cd := make([]CommitData, len(commits))
	for i := range commits {
		c := commits[i]

		cd[i].Commit = c.Commit
		cd[i].Message = c.Message

		cd[i].Author = CommitAuthor{
			Email: c.Author.Email,
			Name:  c.Author.Name,
			Date:  c.Author.Date,
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
