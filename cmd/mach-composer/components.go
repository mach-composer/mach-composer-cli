package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var componentsCmd = &cobra.Command{
	Use:   "components",
	Short: "List all components.",
	PreRun: func(cmd *cobra.Command, args []string) {
		preprocessCommonFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		cfg := loadConfig(cmd, false)
		defer cfg.Close()

		fmt.Printf("%s:\n", commonFlags.configFile)
		for _, component := range cfg.Components {
			fmt.Printf(" - %s\n", component.Name)
			fmt.Printf("   version: %s\n", component.Version)
		}

		fmt.Println("")
	},
}

func init() {
	registerCommonFlags(componentsCmd)
}
