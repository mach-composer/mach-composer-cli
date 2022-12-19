package main

import (
	"os"

	amplience "github.com/mach-composer/mach-composer-plugin-amplience/plugin"
	aws "github.com/mach-composer/mach-composer-plugin-aws/plugin"
	azure "github.com/mach-composer/mach-composer-plugin-azure/plugin"
	commercetools "github.com/mach-composer/mach-composer-plugin-commercetools/plugin"
	contentful "github.com/mach-composer/mach-composer-plugin-contentful/plugin"
	sentry "github.com/mach-composer/mach-composer-plugin-sentry/plugin"
	"github.com/spf13/cobra"
)

// This is a temporary commando. Will be removed once all plugins are released
// separately from mach-composer
var pluginCmd = &cobra.Command{
	Use:   "plugin [name]",
	Short: "Start a plugin for mach-composer",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		plugins := map[string]func(){
			"amplience":     amplience.Serve,
			"aws":           aws.Serve,
			"azure":         azure.Serve,
			"contentful":    contentful.Serve,
			"commercetools": commercetools.Serve,
			"sentry":        sentry.Serve,
		}

		if serve, ok := plugins[args[0]]; ok {
			serve()
			os.Exit(0)
		} else {
			cmd.Println("invalid plugin specified")
			os.Exit(1)
		}
	},
}

func init() {
	pluginCmd.Hidden = true
	rootCmd.AddCommand(pluginCmd)
}
