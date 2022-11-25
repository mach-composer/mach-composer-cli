package cloudcmd

import (
	"os"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"
)

var listProjectCmd = &cobra.Command{
	Use:   "list-projects",
	Short: "List all Projects",
	RunE: func(cmd *cobra.Command, args []string) error {
		organization := MustGetString(cmd, "organization")

		client, ctx := getClient(cmd)
		paginator, _, err := (client.
			AccountManagementApi.
			ProjectQuery(ctx, organization).
			Execute())
		if err != nil {
			return handleError(err)
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			data[i] = []string{
				record.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				record.Name,
				record.Key,
			}
		}
		writeTable(os.Stdout, []string{"Created At", "Name", "Key"}, data)
		return nil
	},
}

var createProjectCmd = &cobra.Command{
	Use:   "create-project",
	Short: "Create a new Project",
	RunE: func(cmd *cobra.Command, args []string) error {
		name := MustGetString(cmd, "name")
		key := MustGetString(cmd, "key")
		organization := MustGetString(cmd, "organization")

		client, ctx := getClient(cmd)
		resource, _, err := (client.
			AccountManagementApi.
			ProjectCreate(ctx, organization).
			ProjectDraft(mccsdk.ProjectDraft{
				Name: name,
				Key:  key,
			}).
			Execute())
		if err != nil {
			return handleError(err)
		}

		cmd.Printf("Created new component: %s\n", resource.GetKey())
		return nil
	},
}

func init() {
	// Create project
	createProjectCmd.Flags().String("organization", "", "organization")
	createProjectCmd.Flags().String("name", "", "name of the project")
	createProjectCmd.Flags().String("key", "", "key for the project")
	Must(createProjectCmd.MarkFlagRequired("organization"))
	Must(createProjectCmd.MarkFlagRequired("name"))
	Must(createProjectCmd.MarkFlagRequired("key"))
	CloudCmd.AddCommand(createProjectCmd)

	// List projects
	listProjectCmd.Flags().String("organization", "", "The organization key to use")
	Must(listProjectCmd.MarkFlagRequired("organization"))
	CloudCmd.AddCommand(listProjectCmd)
}
