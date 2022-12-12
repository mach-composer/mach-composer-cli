package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/cli"
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
			cli.PrintExitError(fmt.Sprintf("Config file %s does not exist", generateFlags.configFile))
		}
		cli.PrintExitError(err.Error())
	}
	if generateFlags.varFile != "" {
		if _, err := os.Stat(generateFlags.varFile); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				cli.PrintExitError(fmt.Sprintf("Variable file %s does not exist", generateFlags.varFile))
			}
			log.Error().Msgf("error: %s\n", err.Error())
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
	cli.PrintExitError("An error occured:", err.Error())
	return nil
}

// loadConfig parses and validates the given config file path.
func loadConfig(ctx context.Context, resolveVars bool) *config.MachConfig {
	opts := &config.ConfigOptions{
		NoResolveVars: !resolveVars,
	}
	if generateFlags.varFile != "" {
		opts.VarFilenames = []string{generateFlags.varFile}
	}

	cfg, err := config.Open(ctx, generateFlags.configFile, opts)
	if err != nil {
		cli.PrintExitError("An error occured while loading the config file", err.Error())
	}
	return cfg
}
