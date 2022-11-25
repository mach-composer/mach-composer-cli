package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "List all components.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessGenerateFlags()
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg := LoadConfig(cmd.Context())
		generateFlags.ValidateSite(cfg)

		fmt.Printf("%s:\n", generateFlags.configFile)
		for _, component := range cfg.Components {
			fmt.Printf(" - %s\n", component.Name)
			fmt.Printf("   version: %s\n", component.Version)
		}

		fmt.Println("")
		return nil
	},
}

func init() {
	registerGenerateFlags(componentsCmd)
}
