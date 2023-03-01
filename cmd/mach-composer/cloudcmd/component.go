package cloudcmd

import (
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"

	"github.com/labd/mach-composer/internal/cloud"
)

var componentCreateCmd = &cobra.Command{
	Use:   "create-component [name]",
	Short: "Register a new component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		resource, _, err := client.
			ComponentsApi.
			ComponentCreate(ctx, organization, project).
			ComponentDraft(mccsdk.ComponentDraft{
				Key: args[0],
			}).
			Execute()
		if err != nil {
			return err
		}
		cmd.Printf("Created new component: %s\n", resource.GetKey())
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
		paginator, _, err := (client.
			ComponentsApi.
			ComponentQuery(ctx, organization, project).
			Limit(250).
			Execute())

		if err != nil {
			return err
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			data[i] = []string{
				record.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				*record.Key,
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
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		componentKey := args[0]

		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		auto, err := cmd.Flags().GetBool("auto")
		if err != nil {
			return err
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		if !auto {
			if len(args) < 2 {
				cmd.Printf("Missing version argument")
				os.Exit(1)
			}
			version := args[1]
			resource, _, err := client.
				ComponentsApi.
				ComponentVersionCreate(ctx, organization, project, componentKey).
				ComponentVersionDraft(mccsdk.ComponentVersionDraft{
					Version: version,
				}).
				Execute()
			if err != nil {
				return err
			}
			cmd.Printf("Created new version for component %s: %s\n",
				resource.GetComponent(), resource.GetVersion())

		} else {
			_, err := cloud.AutoRegisterVersion(ctx, client, organization, project, componentKey)
			if err != nil {
				return err
			}
		}
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
	Short: "List all changes for an component version",
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

	CloudCmd.AddCommand(componentRegisterVersionCmd)
	registerContextFlags(componentRegisterVersionCmd)
	componentRegisterVersionCmd.Flags().Bool("auto", false, "Automate")

	CloudCmd.AddCommand(componentListVersionCmd)
	registerContextFlags(componentListVersionCmd)

	CloudCmd.AddCommand(componentDescribeVersionCmd)
	registerContextFlags(componentDescribeVersionCmd)
}
