package cloud

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestRegisterComponentVersionComponentNotFoundWithoutCreateComponent(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, false, false, gitFilterPaths)
	assert.ErrorContains(t, err, "component test-component does not exist")
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "CreateComponentVersion"))
}

func TestRegisterComponentVersionComponentNotFoundWithCreateComponentDryRun(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, true, false, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "CreateComponentVersion"))
}

func TestRegisterComponentVersionComponentNotFoundWithCreateComponent(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{}, nil)
	client.On("CreateComponent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.Component{
		Key: "test-component",
	}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, false, true, gitFilterPaths)
	assert.NoError(t, err)
}

func TestRegisterComponentVersionComponentFound(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{
		Results: []mccsdk.Component{
			{
				Key: "test-component",
			},
		},
	}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, false, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
}

func TestRegisterComponentVersionComponentFoundAutoGitRevisionNotFound(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{
		Results: []mccsdk.Component{
			{
				Key: "test-component",
			},
		},
	}, nil)
	client.On("GetLatestComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "previous-test-component-version"}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}
	gitRepo.On("GetCurrentBranch", mock.Anything, mock.Anything).Return("main", nil)
	gitRepo.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(&gitutils.GitVersionInfo{}, nil)
	gitRepo.On("GetRecentCommits", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]gitutils.GitCommit{}, gitutils.ErrGitRevisionNotFound)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, true, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "PushComponentVersionCommits"))
}

func TestRegisterComponentVersionComponentFoundAutoNoCommits(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{
		Results: []mccsdk.Component{
			{
				Key: "test-component",
			},
		},
	}, nil)
	client.On("GetLatestComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "previous-test-component-version"}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}
	gitRepo.On("GetCurrentBranch", mock.Anything, mock.Anything).Return("main", nil)
	gitRepo.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(&gitutils.GitVersionInfo{}, nil)
	gitRepo.On("GetRecentCommits", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]gitutils.GitCommit{}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, true, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "PushComponentVersionCommits"))
}

func TestRegisterComponentVersionComponentFoundAutoOK(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{
		Results: []mccsdk.Component{
			{
				Key: "test-component",
			},
		},
	}, nil)
	client.On("GetLatestComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "previous-test-component-version"}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)
	client.On("PushComponentVersionCommits", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.MatchedBy(func(s []mccsdk.CommitDraft) bool {
		assert.Len(t, s, 1)
		assert.Equal(t, "test-commit", s[0].Commit)
		assert.Equal(t, "test-subject", s[0].Subject)
		assert.Equal(t, "test-parent", s[0].Parents[0])
		assert.Equal(t, "test-author", s[0].Author.Name)
		assert.Equal(t, "test-committer", s[0].Committer.Name)

		return true
	},
	)).Return(nil)

	gitRepo := &gitutils.GitRepositoryMock{}
	gitRepo.On("GetCurrentBranch", mock.Anything, mock.Anything).Return("main", nil)
	gitRepo.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(&gitutils.GitVersionInfo{}, nil)
	gitRepo.On("GetRecentCommits", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]gitutils.GitCommit{
		{
			Commit:  "test-commit",
			Message: "test-subject",
			Parents: []string{"test-parent"},
			Author: gitutils.GitCommitAuthor{
				Name:  "test-author",
				Email: "test-author@email.com",
				Date:  time.Now(),
			},
			Committer: gitutils.GitCommitAuthor{
				Name:  "test-committer",
				Email: "test-committer@email.com",
				Date:  time.Now(),
			},
			Tags: []string{"test-tag"},
		},
	}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, true, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
}

func TestRegisterComponentVersionComponentFoundAutoDryRun(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{
		Results: []mccsdk.Component{
			{
				Key: "test-component",
			},
		},
	}, nil)
	client.On("GetLatestComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "previous-test-component-version"}, nil)

	gitRepo := &gitutils.GitRepositoryMock{}
	gitRepo.On("GetCurrentBranch", mock.Anything, mock.Anything).Return("main", nil)
	gitRepo.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(&gitutils.GitVersionInfo{}, nil)
	gitRepo.On("GetRecentCommits", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]gitutils.GitCommit{
		{
			Commit:  "test-commit",
			Message: "test-subject",
			Parents: []string{"test-parent"},
			Author: gitutils.GitCommitAuthor{
				Name:  "test-author",
				Email: "test-author@email.com",
				Date:  time.Now(),
			},
			Committer: gitutils.GitCommitAuthor{
				Name:  "test-committer",
				Email: "test-committer@email.com",
				Date:  time.Now(),
			},
			Tags: []string{"test-tag"},
		},
	}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, true, true, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "CreateComponentVersion"))
	assert.True(t, client.AssertNotCalled(t, "PushComponentVersionCommits"))
}
