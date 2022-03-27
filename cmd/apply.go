package cmd

import (
	"github.com/labd/mach-composer-go/config"
	"github.com/labd/mach-composer-go/generator"
	"github.com/labd/mach-composer-go/runner"
	"github.com/spf13/cobra"
)

var applyFlags struct {
	reuse bool
}

var applyCmd = &cobra.Command{
	Use:   "apply",
	Short: "Apply the configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := applyFunc(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(applyCmd)
	applyCmd.Flags().BoolVarP(&applyFlags.reuse, "reuse", "", false, "Supress a terraform init for improved speed (not recommended for production usage)")
}

func applyFunc(args []string) error {

	configs := make(map[string]*config.MachConfig)
	allPaths := make(map[string]map[string]string)

	// Note that we do this in multiple passes to minimize ending up with
	// half broken runs. We could in the future also run some parts in parallel

	// Load and parse all config files
	for _, filename := range generateFlags.fileNames {
		cfg, err := config.Load(filename)
		if err != nil {
			panic(err)
		}
		configs[filename] = cfg
	}

	// Write the generate files for each config
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths, err := generator.WriteFiles(cfg, generateFlags.outputPath)
		if err != nil {
			panic(err)
		}
		allPaths[filename] = paths
	}

	// Apply the generate files
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths := allPaths[filename]
		runner.TerraformApply(cfg, paths, applyFlags.reuse)
	}

	return nil
}
