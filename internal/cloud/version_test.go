package cloud

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/gitutils"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
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

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, false, false)
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

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, true, false, true)
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

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, false, true)
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

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, false, true)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
}

// In auto mode the version is taken from the current commit of the checked out
// branch. No commits are looked up or registered.
func TestRegisterComponentVersionComponentFoundAutoOK(t *testing.T) {
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
	gitRepo.On("GetCurrentBranch", mock.Anything, mock.Anything).Return("main", nil)
	gitRepo.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(&gitutils.GitVersionInfo{}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, false, true, true)
	assert.NoError(t, err)
	assert.True(t, client.AssertCalled(t, "CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything))
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	gitRepo.AssertExpectations(t)
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

	gitRepo := &gitutils.GitRepositoryMock{}
	gitRepo.On("GetCurrentBranch", mock.Anything, mock.Anything).Return("main", nil)
	gitRepo.On("GetVersionInfo", mock.Anything, mock.Anything, mock.Anything).Return(&gitutils.GitVersionInfo{}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"

	err := RegisterComponentVersion(ctx, client, gitRepo, organization, project, componentKey, branch, version, true, true, true)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "CreateComponentVersion"))
}
