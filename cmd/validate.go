package cmd

import (
	"github.com/labd/mach-composer/generator"
	"github.com/labd/mach-composer/runner"
	"github.com/spf13/cobra"
)

var validateFlags struct {
	reuse bool
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate MACH composer configuration",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateFunc(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(validateCmd)
	validateCmd.MarkFlagRequired("site")
	validateCmd.Flags().BoolVarP(&validateFlags.reuse, "reuse", "", false, "Suppress a terraform init for improved speed (not recommended for production usage)")
}

func validateFunc(args []string) error {
	configs := LoadConfigs()
	allPaths := make(map[string]map[string]string)

	// Write the generate files for each config
	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths, err := generator.WriteFiles(cfg, genOptions)
		if err != nil {
			panic(err)
		}
		allPaths[filename] = paths
	}

	// Validate the generate files
	options := &runner.ValidateOptions{
		Reuse: validateFlags.reuse,
		Site:  generateFlags.siteName,
	}
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths := allPaths[filename]
		runner.TerraformValidate(cfg, paths, options)
	}

	return nil
}
