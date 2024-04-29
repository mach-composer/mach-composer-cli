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

	version, _, err := client.
		ComponentsApi.ComponentLatestVersion(ctx, organization, project, c.Name).
		Branch(c.Branch).
		Execute()
	if err != nil {
		return err
	}

	if version == nil {
		return fmt.Errorf("failed to resolve latest version for component %s (branch %s)", c.Name, c.Branch)
	}

	c.Version = version.Version
	return nil
}
