package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/labd/mach-composer-go/updater"
	"github.com/spf13/cobra"
)

var updateFlags struct {
	fileNames []string
	check     bool
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all (or a given) component.",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(generateFlags.fileNames) < 1 {
			matches, err := filepath.Glob("./*.yml")
			if err != nil {
				log.Fatal(err)
			}
			generateFlags.fileNames = matches
			if len(generateFlags.fileNames) < 1 {
				fmt.Println("No .yml files found")
				os.Exit(1)
			}
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := updateFunc(args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	updateCmd.Flags().StringArrayVarP(&updateFlags.fileNames, "file", "f", nil, "YAML file to update. If not set update all *.yml files.")
	updateCmd.Flags().BoolVarP(&updateFlags.check, "check", "", false, "Only checks for updates, doesnt change files.")
}

func updateFunc(args []string) error {
	for _, filename := range updateFlags.fileNames {
		updater.UpdateFile(filename)
	}
	return nil
}
