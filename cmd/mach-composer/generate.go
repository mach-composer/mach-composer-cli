package main

import (
	"github.com/labd/mach-composer/generator"
	"github.com/spf13/cobra"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate the Terraform files.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := generateFunc(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(generateCmd)
}

func generateFunc(args []string) error {
	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}

	configs := LoadConfigs()
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]

		_, err := generator.WriteFiles(cfg, genOptions)
		if err != nil {
			panic(err)
		}
	}
	return nil
}
