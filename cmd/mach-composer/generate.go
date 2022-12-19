package main

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/generator"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		handleError(generateFunc(cmd.Context(), args))
	},
}

func init() {
	registerGenerateFlags(generateCmd)
}

func generateFunc(ctx context.Context, args []string) error {
	cfg := loadConfig(ctx, true)
	generateFlags.ValidateSite(cfg)

	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}

	_, err := generator.WriteFiles(cfg, genOptions)
	if err != nil {
		return err
	}

	return nil
}
