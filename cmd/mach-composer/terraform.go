package main

import (
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
	RunE: func(cmd *cobra.Command, args []string) error {
		return handleError(terraformFunc(args))
	},
}

func init() {
	registerGenerateFlags(terraformCmd)
	if err := terraformCmd.MarkFlagRequired("site"); err != nil {
		panic(fmt.Errorf("terraformCmd.MarkFlagRequired: %v", err))
	}
}

func terraformFunc(args []string) error {
	cfg := LoadConfig()
	generateFlags.ValidateSite(cfg)

	fileLocations := generator.FileLocations(cfg, &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	})

	return runner.TerraformProxy(cfg, fileLocations, generateFlags.siteName, args)
}
