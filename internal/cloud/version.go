package cloud

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"

	"github.com/labd/mach-composer/internal/gitutils"
)

func AutoRegisterVersion(ctx context.Context, client *mccsdk.APIClient, organization, project, componentKey string, dryRun bool, gitFilterPaths []string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	branch, err := gitutils.GetCurrentBranch(ctx, cwd)
	if err != nil {
		return "", err
	}

	lastVersion, _, err := client.
		ComponentsApi.
		ComponentLatestVersion(ctx, organization, project, componentKey).
		Branch(branch).
		Execute() //nolint:bodyclose
	if err != nil {
		return "", err
	}

	baseRef := ""
	if lastVersion != nil {
		baseRef = lastVersion.Version
	}

	newVersion, err := gitutils.GetVersionInfo(ctx, cwd, branch)
	if err != nil {
		return "", err
	}
	versionIdentifier := newVersion.Identifier()

	// Register new version
	if dryRun {
		fmt.Printf("Would create new version: %s (branch=%s)\n", versionIdentifier, branch)
	} else {
		createdVersion, _, err := client.
			ComponentsApi.
			ComponentVersionCreate(ctx, organization, project, componentKey).
			ComponentVersionDraft(mccsdk.ComponentVersionDraft{
				Version: versionIdentifier,
				Branch:  branch,
			}).
			Execute() //nolint:bodyclose
		if err != nil {
			return "", err
		}
		fmt.Printf("Created new version: %s (branch=%s)\n", createdVersion.Version, branch)
	}

	commits, err := gitutils.GetRecentCommits(ctx, componentKey, cwd, baseRef, branch, gitFilterPaths)
	if err != nil {
		if errors.Is(err, gitutils.ErrGitRevisionNotFound) {
			fmt.Printf("Failed to calculate changes, last version (%s) not found in the repository\n", baseRef)
			return "", nil
		}
		return "", err
	}

	if len(commits) == 0 {
		fmt.Printf("No new commits found since last version (%s)\n", baseRef)
		return "'", nil
	}

	// Push commits
	newCommits := make([]mccsdk.CommitData, len(commits))
	for i := range pie.Reverse(commits) {
		c := commits[i]
		newCommits[i] = mccsdk.CommitData{
			Commit:  c.Commit,
			Subject: c.Message,
			Parents: c.Parents,
			Author: mccsdk.CommitDataAuthor{
				Name:  c.Author.Name,
				Email: c.Author.Email,
				Date:  c.Author.Date,
			},
			Committer: mccsdk.CommitDataAuthor{
				Name:  c.Committer.Name,
				Email: c.Committer.Email,
				Date:  c.Committer.Date,
			},
		}
	}

	if !dryRun {
		_, err = client.
			ComponentsApi.
			ComponentVersionPushCommits(ctx, organization, project, componentKey, versionIdentifier).
			ComponentVersionCommits(mccsdk.ComponentVersionCommits{
				Commits: newCommits,
			}).
			Execute() //nolint:bodyclose
		if err != nil {
			return versionIdentifier, err
		}
		fmt.Printf("Recorded %d commits for version: %s\n", len(newCommits), versionIdentifier)
	} else {
		fmt.Printf("Found %d commits for version: %s\n", len(newCommits), versionIdentifier)
		for _, c := range newCommits {
			fmt.Printf("%s %s\n", c.Commit, c.Subject)
		}
	}
	return versionIdentifier, nil
}
