package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var sitesCmd = &cobra.Command{
	Use:   "sites",
	Short: "List all sites.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := LoadConfig(cmd.Context())
		generateFlags.ValidateSite(cfg)

		fmt.Printf("%s:\n", generateFlags.configFile)
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
	registerGenerateFlags(sitesCmd)
}
