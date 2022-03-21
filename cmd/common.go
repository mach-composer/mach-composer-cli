package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var generateFlags struct {
	fileNames     []string
	siteName      string
	ignoreVersion bool
	outputPath    string
}

func registerGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&generateFlags.fileNames, "file", "f", nil, "YAML file to parse. If not set parse all *.yml files.")
	cmd.Flags().StringVarP(&generateFlags.siteName, "site", "s", "", "Site to parse. If not set parse all sites.")
	cmd.Flags().BoolVarP(&generateFlags.ignoreVersion, "ignore-version", "", false, "Skip MACH composer version check")
	cmd.Flags().StringVarP(&generateFlags.outputPath, "output-path", "", "", "Output path, defaults to `cwd`/deployments.")
}

func preprocessGenerateFlags() {
	if len(generateFlags.fileNames) < 1 {
		matches, err := filepath.Glob("./*.yml")
		if err != nil {
			log.Fatal(err)
		}
		generateFlags.fileNames = matches
		if len(generateFlags.fileNames) < 1 {
			fmt.Println("No .yml files found")
			os.Exit(1)
		}
	}

	if generateFlags.outputPath == "" {
		var err error
		value, err := os.Getwd()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		} else {
			generateFlags.outputPath = filepath.Join(value, "deployments")
		}
	}
}
