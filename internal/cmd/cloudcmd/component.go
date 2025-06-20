package cloudcmd

import (
	"fmt"
	"github.com/gosimple/slug"
	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"os"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
)

var componentCreateCmd = &cobra.Command{
	Use:   "create-component [name] [key]",
	Short: "Register a new component. If key is not provided it will be generated from the name",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		var key string
		if len(args) == 1 {
			key = slug.Make(args[0])
		} else {
			key = args[1]
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		resource, _, err := client.
			ComponentsApi.
			ComponentCreate(ctx, organization, project).
			ComponentDraft(mccsdk.ComponentDraft{
				Name: args[0],
				Key:  key,
			}).
			Execute()
		if err != nil {
			return err
		}
		log.Info().Msgf("Created new component: %s (key: %s)\n", resource.GetName(), resource.GetKey())
		return nil
	},
}

var componentListCmd = &cobra.Command{
	Use:   "list-components",
	Short: "List your components",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		paginator, _, err := client.
			ComponentsApi.
			ComponentQuery(ctx, organization, project).
			Limit(250).
			Execute()

		if err != nil {
			return err
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			data[i] = []string{
				record.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				record.Key,
			}
		}

		writeTable(os.Stdout,
			[]string{"Created At", "Key"},
			data,
		)
		return nil
	},
}

var componentRegisterVersionCmd = &cobra.Command{
	Use:   "register-component-version [name] [version]",
	Short: "Register a new version for an existing component",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		componentKey := args[0]

		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")
		gitFilterPaths, err := cmd.Flags().GetStringArray("git-filter-path")
		if err != nil {
			return err
		}

		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}

		branch, err := cmd.Flags().GetString("branch")
		if err != nil {
			return err
		}

		dryRun, err := cmd.Flags().GetBool("dry-run")
		if err != nil {
			return err
		}

		createComponent, err := cmd.Flags().GetBool("create-component")
		if err != nil {
			return err
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		var version string
		if len(args) == 1 && !auto {
			return fmt.Errorf("you must provide a version when not using the --auto flag")
		} else if len(args) >= 2 {
			if auto {
				log.Warn().Msgf("ignoring --auto flag, version will be set to %s", args[1])
			}
			version = args[1]
		}

		return cloud.RegisterComponentVersion(ctx, cloud.NewClientWrapper(client), gitutils.NewGitRepositoryWrapper(), organization, project, componentKey, branch, version, dryRun, auto, createComponent, gitFilterPaths)
	},
}

var componentUpdateCmd = &cobra.Command{
	Use:   "update-component [key]",
	Short: "Update component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")
		componentKey := args[0]
		var patchedDraft = mccsdk.PatchedComponentDraft{}

		newKey, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}
		if newKey != "" {
			patchedDraft.Key = &newKey
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		resource, _, err := client.
			ComponentsApi.
			ComponentPatch(ctx, organization, project, componentKey).
			PatchedComponentDraft(patchedDraft).
			Execute()
		if err != nil {
			return err
		}

		log.Info().Msgf("Updated component: %s\n", resource.GetKey())
		return nil
	},
}

var componentListVersionCmd = &cobra.Command{
	Use:   "list-component-versions [name]",
	Short: "List all version for an existing component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		componentKey := args[0]

		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		paginator, _, err := client.
			ComponentsApi.
			ComponentVersionQuery(ctx, organization, project, componentKey).
			Execute()
		if err != nil {
			return err
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			data[i] = []string{
				record.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				record.GetVersion(),
			}
		}

		writeTable(os.Stdout,
			[]string{"Created At", "Key"},
			data,
		)

		return nil
	},
}

var componentDescribeVersionCmd = &cobra.Command{
	Use:   "describe-component-versions [name] [version]",
	Short: "List all changes for a component version",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		key := args[0]
		version := args[1]

		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		paginator, _, err := client.
			ComponentsApi.
			ComponentVersionQueryCommits(ctx, organization, project, key, version).
			Execute()
		if err != nil {
			return err
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			data[i] = []string{
				record.Author.Date.Local().Format("2006-01-02 15:04:05"),
				record.Commit,
				record.Author.Name,
				record.Subject,
			}
		}

		writeTable(os.Stdout,
			[]string{"Date", "Commit", "Author", "Message"},
			data,
		)

		return nil
	},
}

func init() {
	CloudCmd.AddCommand(componentCreateCmd)
	registerContextFlags(componentCreateCmd)

	CloudCmd.AddCommand(componentListCmd)
	registerContextFlags(componentListCmd)

	CloudCmd.AddCommand(componentUpdateCmd)
	registerContextFlags(componentUpdateCmd)
	componentUpdateCmd.Flags().String("key", "", "Set new component key")

	CloudCmd.AddCommand(componentRegisterVersionCmd)
	registerContextFlags(componentRegisterVersionCmd)
	componentRegisterVersionCmd.Flags().Bool("auto", false, "Add the version commits automatically based on the current branch")
	componentRegisterVersionCmd.Flags().Bool("dry-run", false, "Dry run")
	componentRegisterVersionCmd.Flags().StringArray("git-filter-path", nil, "Filter commits based on given paths")
	componentRegisterVersionCmd.Flags().String("branch", "", "The branch to use for the version. Defaults to the backend default if not set")
	componentRegisterVersionCmd.Flags().Bool("create-component", false, "Will create the component if it does not already exist")

	CloudCmd.AddCommand(componentListVersionCmd)
	registerContextFlags(componentListVersionCmd)

	CloudCmd.AddCommand(componentDescribeVersionCmd)
	registerContextFlags(componentDescribeVersionCmd)
}
