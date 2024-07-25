package cloudcmd

import (
	"os"
	"path/filepath"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/spf13/cobra"

	"github.com/mach-composer/mach-composer-cli/internal/cloud"
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
					Branch:  branch,
				}).
				Execute()
			if err != nil {
				return err
			}
			cmd.Printf("Created new version for component %s: %s\n",
				resource.GetComponent(), resource.GetVersion())

		} else {
			dryRun, err := cmd.Flags().GetBool("dry-run")
			if err != nil {
				return err
			}

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}
			gitFilterPaths = pie.Map(gitFilterPaths, func(path string) string {
				if filepath.IsAbs(path) {
					return path
				}

				return filepath.Join(cwd, path)
			})
			_, err = cloud.AutoRegisterVersion(ctx, client, organization, project, componentKey, dryRun, gitFilterPaths)
			if err != nil {
				return err
			}
		}
		return nil
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

		newKey, err := cmd.Flags().GetString("key")
		if err != nil {
			return err
		}
		patches := []mccsdk.PatchRequestInner{}

		if newKey != "" {
			patches = append(patches, mccsdk.PatchRequestInner{
				JSONPatchRequestAddReplaceTest: &mccsdk.JSONPatchRequestAddReplaceTest{
					Path:  "/key",
					Op:    "replace",
					Value: newKey,
				},
			})
		}

		if len(patches) == 0 {
			return nil
		}

		client, err := cloud.NewClient(ctx)
		if err != nil {
			return err
		}
		resource, _, err := client.
			ComponentsApi.
			ComponentPatch(ctx, organization, project, componentKey).
			PatchRequestInner(patches).
			Execute()
		if err != nil {
			return err
		}

		cmd.Printf("Updated component: %s\n", resource.GetKey())
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
	componentRegisterVersionCmd.Flags().Bool("auto", false, "Automate")
	componentRegisterVersionCmd.Flags().Bool("dry-run", false, "Dry run")
	componentRegisterVersionCmd.Flags().StringArray("git-filter-path", nil, "Filter commits based on given paths")
	componentRegisterVersionCmd.Flags().String("branch", "", "The branch to use for the version. Defaults to the backend default if not set")

	CloudCmd.AddCommand(componentListVersionCmd)
	registerContextFlags(componentListVersionCmd)

	CloudCmd.AddCommand(componentDescribeVersionCmd)
	registerContextFlags(componentDescribeVersionCmd)
}
