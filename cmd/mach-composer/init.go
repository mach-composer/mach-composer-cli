package main

import (
	"context"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize site directories Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		return handleError(initFunc(cmd.Context(), args))
	},
}

func init() {
	registerGenerateFlags(initCmd)
}

func initFunc(ctx context.Context, args []string) error {
	cfg := LoadConfig(ctx)
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
