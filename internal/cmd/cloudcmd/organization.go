package cloudcmd

import (
	"os"
	"strings"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
)

var listOrganizationCmd = &cobra.Command{
	Use:   "list-organizations",
	Short: "List all organizations",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		paginator, _, err := (client.
			OrganizationManagementApi.
			OrganizationQuery(ctx).
			Execute())
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
		writeTable(os.Stdout, []string{"Created At", "Name", "Key"}, data)
		return nil
	},
}

var createOrganizationCmd = &cobra.Command{
	Use:   "create-organization",
	Short: "Create a new organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organizationDraft := mccsdk.OrganizationDraft{
			Name: MustGetString(cmd, "name"),
			Key:  MustGetString(cmd, "key"),
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		resource, _, err := (client.
			OrganizationManagementApi.
			OrganizationCreate(ctx).
			OrganizationDraft(organizationDraft).
			Execute())
		if err != nil {
			return err
		}

		cmd.Printf("Created new organization: %s\n", resource.GetKey())
		return nil
	},
}

var listOrganizationUsersCmd = &cobra.Command{
	Use:   "list-organization-users",
	Short: "List all users in an organization",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		paginator, _, err := (client.
			OrganizationManagementApi.
			OrganizationUserQuery(ctx, organization).
			Execute())
		if err != nil {
			return err
		}

		for i, record := range paginator.Results {
			if i == 0 {
				cmd.Printf("The following users are member of %s\n\n", *record.OrganizationName)
			}
			cmd.Println("  - Name:", *record.Name)
			cmd.Println("    Email:", record.Email)
			cmd.Println("    Scope:", strings.Join(record.OrganizationScopes, " "))
			for j, project := range record.Projects {
				if j == 0 {
					cmd.Println("     Projects:")
				}
				cmd.Println("      - Name:", project.Name)
				cmd.Println("        Key:", project.Key)
				cmd.Println("        Scope:", strings.Join(project.Scopes, " "))
			}
			cmd.Println()
		}
		return nil
	},
}

var addOrganizationUsersCmd = &cobra.Command{
	Use:   "add-organization-user [email-address]",
	Short: "Invite a user to a specific organization",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		email := args[0]
		organization := MustGetString(cmd, "organization")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}

		_, _, err = (client.
			OrganizationManagementApi.
			OrganizationUserInvite(ctx, organization).
			OrganizationUserInviteDraft(mccsdk.OrganizationUserInviteDraft{
				Email: email,
			}).
			Execute())
		if err != nil {
			return err
		}

		cmd.Println("User added to organization")
		return nil
	},
}

func init() {
	// Create organization
	CloudCmd.AddCommand(createOrganizationCmd)
	createOrganizationCmd.Flags().String("name", "", "name of the organization")
	createOrganizationCmd.Flags().String("key", "", "key for the organization")
	Must(createOrganizationCmd.MarkFlagRequired("name"))
	Must(createOrganizationCmd.MarkFlagRequired("key"))

	// List organizations
	CloudCmd.AddCommand(listOrganizationCmd)

	// List organization users
	listOrganizationUsersCmd.Flags().String("organization", "", "The organization key to use")
	Must(listOrganizationUsersCmd.MarkFlagRequired("organization"))
	CloudCmd.AddCommand(listOrganizationUsersCmd)

	// Add organization users
	addOrganizationUsersCmd.Flags().String("organization", "", "The organization key to use")
	Must(addOrganizationUsersCmd.MarkFlagRequired("organization"))
	CloudCmd.AddCommand(addOrganizationUsersCmd)
}
