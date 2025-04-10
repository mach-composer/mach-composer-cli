package cloud

import (
	"context"
	"fmt"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

const (
	LatestVersion        = "$LATEST"
	VersionNotApplicable = "NOT_APPLICABLE"
)

func ResolveComponentsData(ctx context.Context, cfg *config.MachConfig) error {
	for i := range cfg.Components {
		if err := resolveComponentVersion(ctx, cfg, &cfg.Components[i]); err != nil {
			return err
		}
	}
	return nil
}

func resolveComponentVersion(ctx context.Context, cfg *config.MachConfig, c *config.ComponentConfig) error {
	if c.Version != LatestVersion {
		return nil
	}

	if cfg.MachComposer.Cloud.Client == nil {
		client, err := NewClient(ctx)
		if err != nil {
			return err
		}
		cfg.MachComposer.Cloud.Client = client
	}

	client := cfg.MachComposer.Cloud.Client
	organization := cfg.MachComposer.Cloud.Organization
	project := cfg.MachComposer.Cloud.Project

	version, res, err := client.
		ComponentsApi.
		ComponentLatestVersion(ctx, organization, project, c.Name).
		Branch(c.Branch).
		Execute()
	if err != nil {
		if res != nil && res.StatusCode == 404 {
			return fmt.Errorf("component %s (branch %s) could not be found", c.Name, c.Branch)
		}

		return fmt.Errorf("failed to resolve latest version for component %s (branch %s): %s", c.Name, c.Branch, err.Error())
	}

	c.Version = version.Version
	return nil
}
