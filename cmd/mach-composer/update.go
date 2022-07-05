package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/labd/mach-composer/updater"
	"github.com/spf13/cobra"
)

var updateFlags struct {
	fileNames     []string
	check         bool
	components    []string
	commit        bool
	commitMessage string
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all (or a given) component.",
	Args:  cobra.MaximumNArgs(2),
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
	updateCmd.Flags().BoolVarP(&updateFlags.commit, "commit", "c", false, "Automatically commits the change.")
	updateCmd.Flags().StringVarP(&updateFlags.commitMessage, "commit-message", "m", "", "Use a custom message for the commit.")
}

func updateFunc(args []string) error {
	ctx := context.Background()
	changes := map[string]string{}

	componentName := ""
	componentVersion := ""

	if len(args) > 0 {
		componentName = args[0]
		if len(args) > 1 {
			componentVersion = args[1]
		}
	}

	// Iterate through all yaml files and update them all.
	writeChanges := !updateFlags.check
	for _, filename := range updateFlags.fileNames {

		updateSet, err := updater.UpdateFile(ctx, filename, componentName, componentVersion, writeChanges)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to update %s: %v\n", filename, err.Error())
			os.Exit(1)
		}

		if writeChanges && updateSet.HasChanges() {
			if componentName == "" {
				changes[filename] = updateSet.ChangeLog()
			} else {
				changes[filename] = updateSet.ComponentChangeLog(componentName)
			}
		}

	}

	if len(changes) < 1 {
		return nil
	}

	// git commit
	if updateFlags.commit {
		filenames := []string{}
		commitMessage := updateFlags.commitMessage

		for fn := range changes {
			filenames = append(filenames, fn)
		}

		// Generate commit message if not passed
		if updateFlags.commitMessage == "" {
			commitMessage = generateCommitMessage(changes)
		}

		ctx := context.Background()
		updater.Commit(ctx, filenames, commitMessage)
	}
	return nil
}

func generateCommitMessage(changes map[string]string) string {
	multipleFiles := len(changes) > 1
	var cm strings.Builder

	for fn, msg := range changes {
		if multipleFiles {
			fmt.Fprintf(&cm, "Changes for %s:\n", fn)
			cm.WriteString(msg)
			fmt.Fprintln(&cm, "")
		} else {
			cm.WriteString(msg)
		}
	}
	return cm.String()
}
