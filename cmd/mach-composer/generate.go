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
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFunc(cmd)
	},
}

func init() {
	registerGenerateFlags(generateCmd)
}

func generateFunc(cmd *cobra.Command) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	generateFlags.ValidateSite(cfg)

	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}

	gd, err := dependency.ToDeploymentGraph(cfg)
	if err != nil {
		return err
	}

	err = generator.Write(cmd.Context(), cfg, gd, genOptions)
	if err != nil {
		return err
	}

	return nil
}
