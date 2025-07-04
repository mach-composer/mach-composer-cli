package cloudcmd

import (
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
)

var listProjectCmd = &cobra.Command{
	Use:   "list-projects",
	Short: "List all Projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		paginator, _, err := client.
			OrganizationManagementApi.
			ProjectQuery(ctx, organization).
			Execute()
		if err != nil {
			return err
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			data[i] = []string{
				record.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				record.Name,
				record.Key,
			}
		}
		return writeTable(os.Stdout, []string{"Created At", "Name", "Key"}, data)
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create-project [key] [name]",
	Short: "Create a new Project",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		key := args[0]
		name := args[1]

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		resource, _, err := client.
			OrganizationManagementApi.
			ProjectCreate(ctx, organization).
			ProjectDraft(mccsdk.ProjectDraft{
				Name: name,
				Key:  key,
			}).
			Execute()
		if err != nil {
			return err
		}

		cmd.Printf("Created new project: %s\n", resource.GetKey())
		return nil
	},
}

func init() {
	// Create project
	createProjectCmd.Flags().String("organization", "", "organization")
	Must(createProjectCmd.MarkFlagRequired("organization"))
	CloudCmd.AddCommand(createProjectCmd)

	// List projects
	listProjectCmd.Flags().String("organization", "", "The organization key to use")
	Must(listProjectCmd.MarkFlagRequired("organization"))
	CloudCmd.AddCommand(listProjectCmd)
}
