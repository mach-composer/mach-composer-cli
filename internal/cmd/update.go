package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
	"github.com/mach-composer/mach-composer-cli/internal/updater"
)

var updateFlags struct {
	configFile    string
	check         bool
	components    []string
	commit        bool
	commitMessage string
	cloud         bool
	gitFallback   bool
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update all (or a given) component.",
	Args:  cobra.MaximumNArgs(2),
	PreRun: func(cmd *cobra.Command, args []string) {
		if _, err := os.Stat(updateFlags.configFile); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				fmt.Printf("%s: Config file not found\n", updateFlags.configFile)
				os.Exit(1)
			}
			fmt.Printf("error: %s\n", err.Error())
			os.Exit(1)
		}
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		return updateFunc(cmd.Context(), args)
	},
	ValidArgsFunction: AutocompleteComponentName,
}

func init() {
	updateCmd.Flags().StringVarP(&updateFlags.configFile, "file", "f", "main.yml", "YAML file to update.")
	updateCmd.Flags().BoolVarP(&updateFlags.check, "check", "", false, "Only checks for updates, doesnt change files.")
	updateCmd.Flags().BoolVarP(&updateFlags.commit, "commit", "c", false, "Automatically commits the change.")
	updateCmd.Flags().StringVarP(&updateFlags.commitMessage, "commit-message", "m", "", "Use a custom message for the commit.")
	updateCmd.Flags().BoolVar(&updateFlags.cloud, "cloud", false, "Use MACH composer cloud to check for updates.")
	updateCmd.Flags().BoolVar(&updateFlags.gitFallback, "git-fallback", false, "Fallback to git if composer cloud check fails.")
	updateCmd.Flags().StringArrayVarP(&updateFlags.components, "component", "", nil, "")

	handleError(updateCmd.MarkFlagFilename("file", "yml", "yaml"))
}

func updateFunc(ctx context.Context, args []string) error {
	if len(updateFlags.components) > 0 {
		log.Warn().Msgf("Components option not implemented")
	}

	componentName := ""
	componentVersion := ""

	if len(args) > 0 {
		componentName = args[0]
		if len(args) > 1 {
			componentVersion = args[1]
		}
	}

	writeChanges := !updateFlags.check

	u, err := updater.NewUpdater(ctx, updateFlags.configFile, updateFlags.cloud, updateFlags.gitFallback)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to update %s: %v\n", updateFlags.configFile, err.Error())
		os.Exit(1)
	}

	var changes string
	if componentName != "" {
		err = u.UpdateComponent(ctx, componentName, componentVersion)
		if err != nil {
			return err
		}

		updateSet := u.GetUpdateSet()
		if writeChanges && u.Write(ctx) {
			changes = updateSet.ComponentChangeLog(componentName)
		}
	} else {
		err := u.UpdateAllComponents(ctx)
		if err != nil {
			cli.PrintExitError(err.Error())
		}

		updateSet := u.GetUpdateSet()
		if writeChanges && u.Write(ctx) {
			changes = updateSet.ChangeLog()
		}
	}

	if changes == "" {
		return nil
	}

	// git commit
	if updateFlags.commit {
		commitMessage := updateFlags.commitMessage

		// Generate commit message if not passed
		if updateFlags.commitMessage == "" {
			commitMessage = generateCommitMessage(map[string]string{updateFlags.configFile: changes})
		}

		err = gitutils.Commit(ctx, []string{updateFlags.configFile}, commitMessage)
		if err != nil {
			return err
		}
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
