package cloudcmd

import (
	"os"
	"strings"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
)

var listApiClientCmd = &cobra.Command{
	Use:   "list-api-clients",
	Short: "List existing API clients (without secret)",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		paginator, _, err := client.
			APIClientsApi.
			ApiClientQuery(ctx, organization, project).
			Execute()

		if err != nil {
			return err
		}

		data := make([][]string, len(paginator.Results))
		for i, record := range paginator.Results {
			lastUsed := "never"
			if record.LastUsedAt.IsSet() {
				lastUsed = record.GetLastUsedAt().Format("2006-01-02 15:04:05")
			}

			var clientId string
			if record.ClientId != nil {
				clientId = *record.ClientId
			} else {
				clientId = "N/A"
			}
			var clientSecret string
			if record.ClientSecret != nil {
				clientSecret = *record.ClientSecret
			} else {
				clientSecret = "N/A"
			}

			data[i] = []string{
				record.CreatedAt.Local().Format("2006-01-02 15:04:05"),
				clientId,
				clientSecret,
				lastUsed,
				record.GetDescription(),
				strings.Join(record.Scope, " "),
			}
		}

		writeTable(os.Stdout,
			[]string{"Created At", "ClientWrapper ID", "ClientWrapper Secret", "Last Used", "Description", "Scopes"},
			data,
		)
		return nil
	},
}

var createApiClientCmd = &cobra.Command{
	Use:   "create-api-client",
	Short: "Manage your components",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()
		organization := MustGetString(cmd, "organization")
		project := MustGetString(cmd, "project")

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		resource, _, err := client.
			APIClientsApi.
			ApiClientCreate(ctx, organization, project).
			ApiClientDraft(mccsdk.ApiClientDraft{
				Scope: []string{"project:manage"},
			}).
			Execute()

		if err != nil {
			return err
		}

		cmd.Printf("ClientWrapper ID: %s\n", *resource.ClientId)
		cmd.Printf("ClientWrapper Secret: %s\n", *resource.ClientSecret)
		cmd.Printf("Scopes: %s\n", strings.Join(resource.Scope, " "))
		return nil
	},
}

func init() {
	CloudCmd.AddCommand(listApiClientCmd)
	registerContextFlags(listApiClientCmd)

	CloudCmd.AddCommand(createApiClientCmd)
	registerContextFlags(createApiClientCmd)
}
