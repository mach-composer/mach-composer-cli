package cloud

import (
	"context"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"

	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
)

func RegisterComponentVersion(ctx context.Context, client ClientWrapper, repository gitutils.GitRepository, organization, project, componentKey, branch, version string, dryRun, auto, createComponent bool) error {
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
			if _, err := client.CreateComponent(ctx, organization, project, componentKey); err != nil {
				return err
			}
			log.Info().Msgf("Created component %s", componentKey)
		}
	}

	if auto {
		return autoRegisterVersion(ctx, client, repository, organization, project, componentKey, dryRun)
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

// autoRegisterVersion registers the current commit of the checked out branch as
// the new version of the component.
func autoRegisterVersion(ctx context.Context, client ClientWrapper, repository gitutils.GitRepository, organization, project, componentKey string, dryRun bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	branch, err := repository.GetCurrentBranch(ctx, cwd)
	if err != nil {
		return err
	}

	versionIdentifier, err := repository.GetLatestCommitHash(ctx, cwd, branch)
	if err != nil {
		return err
	}

	if dryRun {
		log.Info().Msgf("Would create new version: %s (branch=%s)", versionIdentifier, branch)
		return nil
	}

	createdVersion, err := client.CreateComponentVersion(ctx, organization, project, componentKey, versionIdentifier, branch)
	if err != nil {
		return err
	}
	log.Info().Msgf("Created new version: %s (branch=%s)", createdVersion.Version, branch)

	return nil
}
