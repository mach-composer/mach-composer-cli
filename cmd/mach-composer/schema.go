package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/cli"
	"github.com/labd/mach-composer/internal/config"
)

var schemaCmd = &cobra.Command{
	Use:   "schema",
	Short: "Generate a JSON schema for your config based on the plugins.",
	Run: func(cmd *cobra.Command, args []string) {
		filename, err := cmd.Flags().GetString("file")
		if err != nil {
			cli.PrintExitError(err.Error())
			return
		}
		if _, err := os.Stat(filename); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				cli.PrintExitError(fmt.Sprintf("%s: Config file not found\n", filename))
				return
			}
			cli.PrintExitError(err.Error())
		}

		schema, err := config.GenerateSchema(cmd.Context(), filename, nil)
		if err != nil {
			cli.PrintExitError(err.Error())
		}
		fmt.Println(schema)
	},
}

func init() {
	schemaCmd.Flags().StringP("file", "f", "main.yml", "YAML file to update.")
}
