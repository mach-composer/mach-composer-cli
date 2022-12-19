package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize site directories Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	Run: func(cmd *cobra.Command, args []string) {
		handleError(initFunc(cmd.Context(), args))
	},
}

func init() {
	registerGenerateFlags(initCmd)
}

func initFunc(ctx context.Context, args []string) error {
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

	return runner.TerraformInit(ctx, cfg, paths, &runner.InitOptions{
		Site: generateFlags.siteName,
	})
}
