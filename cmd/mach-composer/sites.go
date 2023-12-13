package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List all sites.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := loadConfig(cmd, false)
		defer cfg.Close()

		fmt.Printf("%s:\n", commonFlags.configFile)
		for _, site := range cfg.Sites {
			fmt.Printf(" - %s\n", site.Identifier)
			fmt.Println("   components:")
			for _, component := range site.Components {
				fmt.Printf("     %s\n", component.Name)
			}
		}

		fmt.Println("")

		return nil
	},
}

func init() {
	registerCommonFlags(sitesCmd)
}
