package main

import (
	"github.com/elliotchance/pie/v2"
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/config"
)

func AutocompleteSiteName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	cfg := loadConfig(cmd, false)

	identifiers := pie.Map(cfg.Sites, func(s config.SiteConfig) string {
		return s.Identifier

	})
	return identifiers, cobra.ShellCompDirectiveNoFileComp
}

func AutocompleteComponentName(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	name, err := cmd.Flags().GetString("file")
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	generateFlags.configFile = name
	cfg := loadConfig(cmd, false)

	identifiers := pie.Map(cfg.Components, func(s config.Component) string {
		return s.Name

	})
	return identifiers, cobra.ShellCompDirectiveNoFileComp
}
