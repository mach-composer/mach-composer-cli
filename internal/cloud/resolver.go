package cloud

import (
	"context"

	"github.com/labd/mach-composer/internal/config"
)

func ResolveComponentsData(ctx context.Context, cfg *config.MachConfig) error {
	for i := range cfg.Components {
		if err := resolveComponentVersion(ctx, cfg, &cfg.Components[i]); err != nil {
			return err
		}
	}
	return nil
}

func resolveComponentVersion(ctx context.Context, cfg *config.MachConfig, c *config.Component) error {
	if c.Version != "$LATEST" {
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

	c.Version = version.Version
	return nil
}
