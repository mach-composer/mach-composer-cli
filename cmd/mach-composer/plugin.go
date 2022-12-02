package main

import (
	"os"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/plugins"
)

// This is a temporary commando. Will be removed once all plugins are released
// separately from mach-composer
var pluginCmd = &cobra.Command{
	Use:   "plugin [name]",
	Short: "Start a plugin for mach-composer",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		all := plugins.GetLocalPlugins()
		if serve, ok := all[args[0]]; ok {
			serve()
			os.Exit(0)
		} else {
			cmd.Println("invalid plugin specified")
			os.Exit(1)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}
