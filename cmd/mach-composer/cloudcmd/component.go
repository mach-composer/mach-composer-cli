package cloudcmd

import (
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"
)

var componentCmd = &cobra.Command{
	Use:   "component",
	Short: "Manage your components",
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

var componentCreateCmd = &cobra.Command{
	Use:   "create-component [name]",
	Short: "Register a new component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, ctx := getClient(cmd)
		resource, _, err := client.
			ComponentsApi.
			ComponentCreate(ctx, organization, project).
			ComponentDraft(mccsdk.ComponentDraft{
				Key: args[0],
			}).
			Execute()
		if err != nil {
			return handleError(err)
		}
		cmd.Printf("Created new component: %s\n", resource.GetKey())
		return nil
	},
}

var componentListCmd = &cobra.Command{
	Use:   "list-components",
	Short: "List your components",
	RunE: func(cmd *cobra.Command, args []string) error {
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, ctx := getClient(cmd)
		paginator, _, err := (client.
			ComponentsApi.
			ComponentQuery(ctx, organization, project).
			Execute())

		if err != nil {
			panic(err)
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
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentKey := args[0]
		version := args[1]

		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, ctx := getClient(cmd)
		resource, _, err := client.
			ComponentsApi.
			ComponentVersionCreate(ctx, organization, project, componentKey).
			ComponentVersionDraft(mccsdk.ComponentVersionDraft{
				Version: version,
			}).
			Execute()
		if err != nil {
			return handleError(err)
		}
		cmd.Printf("Created new version for component %s: %s\n",
			resource.GetComponent(), resource.GetVersion())
		return nil
	},
}

var componentListVersionCmd = &cobra.Command{
	Use:   "list-component-versions [name]",
	Short: "List all version for an existing component",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		componentKey := args[0]

		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, ctx := getClient(cmd)
		paginator, _, err := client.
			ComponentsApi.
			ComponentVersionQuery(ctx, organization, project, componentKey).
			Execute()
		if err != nil {
			return handleError(err)
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

func init() {

	CloudCmd.AddCommand(componentCreateCmd)
	registerContextFlags(componentCreateCmd)

	CloudCmd.AddCommand(componentListCmd)
	registerContextFlags(componentListCmd)

	CloudCmd.AddCommand(componentRegisterVersionCmd)
	registerContextFlags(componentRegisterVersionCmd)

	CloudCmd.AddCommand(componentListVersionCmd)
	registerContextFlags(componentListVersionCmd)
}
