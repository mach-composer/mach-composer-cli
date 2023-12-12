package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
)

type ProxyOptions struct {
	Site    string
	Command []string
}

func TerraformProxy(ctx context.Context, cfg *config.MachConfig, locations map[string]string, options *ProxyOptions) error {
	for i := range cfg.Sites {
		site := cfg.Sites[i]

		if options.Site != "" && site.Identifier != options.Site {
			continue
		}

		err := defaultRunTerraform(ctx, locations[site.Identifier], options.Command...)
		if err != nil {
			return err
		}
	}
	return nil
}
