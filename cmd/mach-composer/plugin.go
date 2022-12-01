package main

import (
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/plugins"
)

var pluginCmd = &cobra.Command{
	Use:   "plugin [name]",
	Short: "Start a plugin for mach-composer",
	RunE: func(cmd *cobra.Command, args []string) error {

		all := plugins.GetLocalPlugins()
		if serve, ok := all[args[0]]; ok {
			serve()

		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(pluginCmd)
}
