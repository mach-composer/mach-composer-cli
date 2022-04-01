package cmd

import (
	"github.com/labd/mach-composer-go/config"
	"github.com/labd/mach-composer-go/generator"
	"github.com/labd/mach-composer-go/runner"
	"github.com/spf13/cobra"
)

var planFlags struct {
	reuse      bool
	components []string
}

var planCmd = &cobra.Command{
	Use:   "plan",
	Short: "Plan the configuration.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},

	RunE: func(cmd *cobra.Command, args []string) error {
		if err := planFunc(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(planCmd)
	planCmd.Flags().BoolVarP(&planFlags.reuse, "reuse", "", false, "Supress a terraform init for improved speed (not recommended for production usage)")
	planCmd.Flags().StringArrayVarP(&planFlags.components, "component", "c", []string{}, "")
}

func planFunc(args []string) error {

	configs := make(map[string]*config.MachConfig)
	allPaths := make(map[string]map[string]string)

	// Note that we do this in multiple passes to minimize ending up with
	// half broken runs. We could in the future also run some parts in parallel

	// Load and parse all config files
	for _, filename := range generateFlags.fileNames {
		cfg, err := config.Load(filename, generateFlags.varFile)
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

	// Plan the generate files
	options := &runner.PlanOptions{
		Reuse: planFlags.reuse,
	}
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths := allPaths[filename]
		runner.TerraformPlan(cfg, paths, options)
	}

	return nil
}
