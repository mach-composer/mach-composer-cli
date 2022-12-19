package main

import (
	"context"
	"fmt"

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
	if err := terraformCmd.MarkFlagRequired("site"); err != nil {
		panic(fmt.Errorf("terraformCmd.MarkFlagRequired: %v", err))
	}
}

func terraformFunc(ctx context.Context, args []string) error {
	cfg := loadConfig(ctx, true)
	defer cfg.Close()

	generateFlags.ValidateSite(cfg)

	fileLocations := generator.FileLocations(cfg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})

	return runner.TerraformProxy(ctx, cfg, fileLocations, generateFlags.siteName, args)
}
