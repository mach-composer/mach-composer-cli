package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labd/mach-composer/config"
	"github.com/spf13/cobra"
)

type GenerateFlags struct {
	fileNames     []string
	siteName      string
	ignoreVersion bool
	outputPath    string
	varFile       string
}

var generateFlags GenerateFlags

func (gf GenerateFlags) ValidateSite(configs map[string]*config.MachConfig) {
	for _, filename := range generateFlags.fileNames {
		cfg := configs[filename]
		if gf.siteName != "" && !cfg.HasSite(gf.siteName) {
			fmt.Fprintf(os.Stderr, "No site found with identifier: %s\n", gf.siteName)
			os.Exit(1)
		}
	}
}

func registerGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringArrayVarP(&generateFlags.fileNames, "file", "f", nil, "YAML file to parse. If not set parse all *.yml files.")
	cmd.Flags().StringVarP(&generateFlags.varFile, "var-file", "", "", "Use a variable file to parse the configuration with.")
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

		for _, m := range matches {
			if generateFlags.varFile == "" && (m == "variables.yml" || m == "variables.yaml") {
				generateFlags.varFile = m
			} else {
				generateFlags.fileNames = append(generateFlags.fileNames, m)
			}
		}
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

// LoadConfigs loads all config files. This means it validates and parses
// the yaml file.
func LoadConfigs() map[string]*config.MachConfig {
	configs := make(map[string]*config.MachConfig)
	for _, filename := range generateFlags.fileNames {
		fmt.Printf("Processing %s\n", filename)

		cfg, err := config.Load(filename, generateFlags.varFile)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			continue
		}
		CheckDeprecations(cfg)
		configs[filename] = cfg
	}
	return configs
}

// CheckDeprecations warns if features have been deprecated
func CheckDeprecations(cfg *config.MachConfig) {
	for _, site := range cfg.Sites {
		if site.Commercetools != nil && site.Commercetools.Frontend != nil {
			fmt.Println("[WARN] Site", site.Identifier, "commercetools frontend block is deprecated and will be removed soon")
		}
	}
}
