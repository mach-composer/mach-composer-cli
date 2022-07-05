package main

import (
	"github.com/labd/mach-composer/generator"
	"github.com/labd/mach-composer/internal/runner"
	"github.com/spf13/cobra"
)

var applyFlags struct {
	reuse       bool
	autoApprove bool
	destroy     bool
	components  []string
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
	applyCmd.Flags().BoolVarP(&applyFlags.autoApprove, "auto-approve", "", false, "Supress a terraform init for improved speed (not recommended for production usage)")
	applyCmd.Flags().BoolVarP(&applyFlags.destroy, "destroy", "", false, "Destroy option is a convenient way to destroy all remote objects managed by this mach config")
	applyCmd.Flags().StringArrayVarP(&applyFlags.components, "component", "c", []string{}, "")
}

func applyFunc(args []string) error {

	allPaths := make(map[string]map[string]string)
	configs := LoadConfigs()

	generateFlags.ValidateSite(configs)

	// Note that we do this in multiple passes to minimize ending up with
	// half broken runs. We could in the future also run some parts in parallel

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

	// Apply the generate files
	options := &runner.ApplyOptions{
		Destroy:     applyFlags.destroy,
		Reuse:       applyFlags.reuse,
		AutoApprove: applyFlags.autoApprove,
		Site:        generateFlags.siteName,
		Components:  applyFlags.components,
	}
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths := allPaths[filename]
		runner.TerraformApply(cfg, paths, options)
	}

	return nil
}
