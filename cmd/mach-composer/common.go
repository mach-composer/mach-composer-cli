package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/config"
)

type GenerateFlags struct {
	configFile    string
	siteName      string
	ignoreVersion bool
	outputPath    string
	varFile       string
}

var generateFlags GenerateFlags

func (gf GenerateFlags) ValidateSite(cfg *config.MachConfig) {
	if gf.siteName != "" && !cfg.HasSite(gf.siteName) {
		fmt.Fprintf(os.Stderr, "No site found with identifier: %s\n", gf.siteName)
		os.Exit(1)
	}
}

func registerGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&generateFlags.configFile, "file", "f", "main.yml", "YAML file to parse.")
	cmd.Flags().StringVarP(&generateFlags.varFile, "var-file", "", "", "Use a variable file to parse the configuration with.")
	cmd.Flags().StringVarP(&generateFlags.siteName, "site", "s", "", "Site to parse. If not set parse all sites.")
	cmd.Flags().BoolVarP(&generateFlags.ignoreVersion, "ignore-version", "", false, "Skip MACH composer version check")
	cmd.Flags().StringVarP(&generateFlags.outputPath, "output-path", "", "", "Output path, defaults to `cwd`/deployments.")
}

func preprocessGenerateFlags() {
	if _, err := os.Stat(generateFlags.configFile); err != nil {
		if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%s: Config file not found\n", generateFlags.configFile)
			os.Exit(1)
		}
		fmt.Printf("error: %s\n", err.Error())
		os.Exit(1)
	}
	if generateFlags.varFile != "" {
		if _, err := os.Stat(generateFlags.varFile); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("%s: Variables file not found\n", generateFlags.varFile)
				os.Exit(1)
			}
			fmt.Printf("error: %s\n", err.Error())
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

func handleError(err error) error {
	if err == nil {
		return nil
	}
	fmt.Printf("Error encountered: %s", err)
	os.Exit(1)
	return nil
}

// LoadConfig parses and validates the given config file path.
func LoadConfig(ctx context.Context) *config.MachConfig {
	cfg, err := config.Load(ctx, generateFlags.configFile, generateFlags.varFile)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	CheckDeprecations(cfg)
	return cfg
}

// CheckDeprecations warns if features have been deprecated
func CheckDeprecations(cfg *config.MachConfig) {
	for _, site := range cfg.Sites {
		if site.Commercetools != nil && site.Commercetools.Frontend != nil {
			fmt.Println("[WARN] Site", site.Identifier, "commercetools frontend block is deprecated and will be removed soon")
		}
	}
}
