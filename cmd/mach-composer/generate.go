package main

import (
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/generator"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFunc(cmd)
	},
}

func init() {
	registerCommonFlags(generateCmd)
}

func generateFunc(cmd *cobra.Command) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	gd, err := dependency.ToDeploymentGraph(cfg, commonFlags.outputPath)
	if err != nil {
		return err
	}

	return generator.Write(cmd.Context(), cfg, gd, nil)
}
