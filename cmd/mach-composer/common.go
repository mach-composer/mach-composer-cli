package main

import (
	"errors"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cloud"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/config"
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
		_, _ = fmt.Fprintf(os.Stderr, "No site found with identifier: %s\n", gf.siteName)
		os.Exit(1)
	}
}

func registerGenerateFlags(cmd *cobra.Command) {
	cmd.Flags().StringVarP(&generateFlags.configFile, "file", "f", "main.yml", "YAML file to parse.")
	cmd.Flags().StringVarP(&generateFlags.varFile, "var-file", "", "", "Use a variable file to parse the configuration with.")
	cmd.Flags().StringVarP(&generateFlags.siteName, "site", "s", "", "DeploymentSite to parse. If not set parse all sites.")
	cmd.Flags().BoolVarP(&generateFlags.ignoreVersion, "ignore-version", "", false, "Skip MACH composer version check")
	cmd.Flags().StringVarP(&generateFlags.outputPath, "output-path", "", "", "Output path, defaults to `cwd`/deployments.")

	handleError(cmd.MarkFlagFilename("var-file", "yml", "yaml"))
	handleError(cmd.MarkFlagFilename("file", "yml", "yaml"))

	_ = cmd.RegisterFlagCompletionFunc("site", AutocompleteSiteName)
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

func handleError(err error) {
	if err != nil {
		cli.PrintExitError("An error occurred:", err.Error())
	}
}

// loadConfig parses and validates the given config file path.
func loadConfig(cmd *cobra.Command, resolveVars bool) *config.MachConfig {
	opts := &config.ConfigOptions{
		NoResolveVars: !resolveVars,
		Validate:      true,
	}
	if generateFlags.varFile != "" {
		opts.VarFilenames = []string{generateFlags.varFile}
	}

	configFile, err := cmd.Flags().GetString("file")
	if err != nil {
		cli.PrintExitError("Missing config filename", err.Error())
	}

	cfg, err := config.Open(cmd.Context(), configFile, opts)
	if err != nil {
		cli.PrintExitError("An error occurred while loading the config file", err.Error())
	}

	if cfg.MachComposer.CloudEnabled() {
		if err := cloud.ResolveComponentsData(cmd.Context(), cfg); err != nil {
			cli.PrintExitError("An error occurred while fetching cloud component info", err.Error())
		}
	}

	cfg.ConfigHash, err = config.ComputeHash(cfg)
	if err != nil {
		cli.PrintExitError("An error occurred while computing hash", err.Error())
	}

	return cfg
}
