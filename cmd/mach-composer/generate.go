package main

import (
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/generator"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return generateFunc(cmd, args)
	},
}

func init() {
	registerGenerateFlags(generateCmd)
}

func generateFunc(cmd *cobra.Command, args []string) error {
	cfg := loadConfig(cmd, true)
	defer cfg.Close()

	generateFlags.ValidateSite(cfg)

	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}

	_, err := generator.WriteFiles(cmd.Context(), cfg, genOptions)
	if err != nil {
		return err
	}

	return nil
}
