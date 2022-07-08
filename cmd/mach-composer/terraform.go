package main

import (
	"fmt"

	"github.com/labd/mach-composer/internal/generator"
	"github.com/labd/mach-composer/internal/runner"
	"github.com/spf13/cobra"
)

var terraformCmd = &cobra.Command{
	Use:   "terraform",
	Short: "Execute terraform commands directly",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := terraformFunc(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	registerGenerateFlags(terraformCmd)
	if err := terraformCmd.MarkFlagRequired("site"); err != nil {
		panic(fmt.Errorf("terraformCmd.MarkFlagRequired: %v", err))
	}
}

func terraformFunc(args []string) error {
	allPaths := make(map[string]map[string]string)
	configs := LoadConfigs()

	generateFlags.ValidateSite(configs)

	// Write the generate files for each config
	genOptions := &generator.GenerateOptions{
		OutputPath: generateFlags.outputPath,
		Site:       generateFlags.siteName,
	}

	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		allPaths[filename] = generator.FileLocations(cfg, genOptions)
	}

	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		paths := allPaths[filename]
		runner.TerraformProxy(cfg, paths, generateFlags.siteName, args)
	}

	return nil
}
