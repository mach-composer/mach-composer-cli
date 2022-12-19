package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Execute terraform commands directly",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleError(terraformFunc(cmd.Context(), args))
	},
}

func init() {
	registerGenerateFlags(terraformCmd)
}

func terraformFunc(ctx context.Context, args []string) error {
	cfg := loadConfig(ctx, true)
	defer cfg.Close()

	generateFlags.ValidateSite(cfg)

	paths, err := generator.WriteFiles(cfg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})
	if err != nil {
		return err
	}

	return runner.TerraformProxy(ctx, cfg, paths, &runner.ProxyOptions{
		Site:    generateFlags.siteName,
		Command: args,
	})
}
