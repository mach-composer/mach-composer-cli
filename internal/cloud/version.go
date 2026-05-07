package cloud

import (
	"context"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"os"
	"path/filepath"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"

	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
)

func RegisterComponentVersion(ctx context.Context, client ClientWrapper, repository gitutils.GitRepository, organization, project, componentKey, branch, version string, dryRun, auto, createComponent bool, gitFilterPaths []string) error {
	lc, err := client.ListComponents(ctx, organization, project, 250)
	if err != nil {
		return err
	}

	var component *mccsdk.Component
	for _, c := range lc.Results {
		if c.GetKey() == componentKey {
			component = &c
			break
		}
	}

	if component == nil {
		if !createComponent {
			return fmt.Errorf("component %s does not exist, create it with `mach-composer cloud create-component` or use the `--create-component` flag to create it", componentKey)
		}

		if dryRun {
			log.Info().Msgf("Would create new component: %s", componentKey)
		} else {
			component, err = client.CreateComponent(ctx, organization, project, componentKey)
			if err != nil {
				return err
			}
			log.Info().Msgf("Created component %s", componentKey)
		}
	}

	if auto {
		return autoRegisterVersion(ctx, client, repository, organization, project, componentKey, dryRun, gitFilterPaths)
	} else {
		if dryRun {
			log.Info().Msgf("Would create new version: %s (branch=%s)", version, branch)
			return nil
		}

		resource, err := client.CreateComponentVersion(ctx, organization, project, componentKey, version, branch)
		if err != nil {
			return err
		}
		log.Info().Msgf("Created new version %s for component %s", resource.GetVersion(), resource.GetComponent())
		return nil
	}
}

func autoRegisterVersion(ctx context.Context, client ClientWrapper, repository gitutils.GitRepository, organization, project, componentKey string, dryRun bool, gitFilterPaths []string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	gitFilterPaths = pie.Map(gitFilterPaths, func(path string) string {
		if filepath.IsAbs(path) {
			return path
		}

		return filepath.Join(cwd, path)
	})

	branch, err := repository.GetCurrentBranch(ctx, cwd)
	if err != nil {
		return err
	}

	var lastVersion *mccsdk.ComponentVersion
	if !dryRun {
		lastVersion, err = client.GetLatestComponentVersion(ctx, organization, project, componentKey, branch)
		if err != nil {
			return err
		}
	}

	baseRef := ""
	if lastVersion != nil {
		baseRef = lastVersion.Version
	}

	newVersion, err := repository.GetVersionInfo(ctx, cwd, branch)
	if err != nil {
		return err
	}
	versionIdentifier := newVersion.Identifier()

	// Register new version
	if dryRun {
		log.Info().Msgf("Would create new version: %s (branch=%s)", versionIdentifier, branch)
	} else {
		createdVersion, err := client.CreateComponentVersion(ctx, organization, project, componentKey, versionIdentifier, branch)
		if err != nil {
			return err
		}
		log.Info().Msgf("Created new version: %s (branch=%s)", createdVersion.Version, branch)
	}

	commits, err := repository.GetRecentCommits(ctx, cwd, baseRef, branch, gitFilterPaths)
	if err != nil {
		if errors.Is(err, gitutils.ErrGitRevisionNotFound) {
			log.Info().Msgf("Failed to calculate changes, last version (%s) not found in the repository", baseRef)
			return nil
		}
		return err
	}

	if len(commits) == 0 {
		log.Info().Msgf("No new commits found since last version (%s)", baseRef)
		return nil
	}

	// Push commits
	newCommits := make([]mccsdk.CommitDraft, len(commits))
	for i := range pie.Reverse(commits) {
		c := commits[i]
		newCommits[i] = mccsdk.CommitDraft{
			Commit:  c.Commit,
			Subject: c.Message,
			Parents: c.Parents,
			Author: mccsdk.CommitAuthorDraft{
				Name:  c.Author.Name,
				Email: c.Author.Email,
				Date:  c.Author.Date,
			},
			Committer: mccsdk.CommitAuthorDraft{
				Name:  c.Committer.Name,
				Email: c.Committer.Email,
				Date:  c.Committer.Date,
			},
		}
	}

	if dryRun {
		log.Info().Msgf("Would add %d commits for version: %s", len(newCommits), versionIdentifier)
		for _, c := range newCommits {
			log.Info().Msgf("%s %s", c.Commit, c.Subject)
		}
	} else {
		err = client.PushComponentVersionCommits(ctx, organization, project, componentKey, versionIdentifier, newCommits)
		if err != nil {
			return err
		}
		log.Info().Msgf("Recorded %d commits for version: %s", len(newCommits), versionIdentifier)
	}
	return nil
}
