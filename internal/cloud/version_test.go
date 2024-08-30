package cloud

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type ClientWrapperMock struct {
	mock.Mock
}

func (c *ClientWrapperMock) ListComponents(ctx context.Context, organization, project string, limit int32) (*mccsdk.ComponentPaginator, error) {
	args := c.Called(ctx, organization, project, limit)
	return args.Get(0).(*mccsdk.ComponentPaginator), args.Error(1)
}

func (c *ClientWrapperMock) CreateComponent(ctx context.Context, organization, project, key string) (*mccsdk.Component, error) {
	args := c.Called(ctx, organization, project, key)
	return args.Get(0).(*mccsdk.Component), args.Error(1)
}

func (c *ClientWrapperMock) CreateComponentVersion(ctx context.Context, organization, project, key, version, branch string) (*mccsdk.ComponentVersion, error) {
	args := c.Called(ctx, organization, project, key, version, branch)
	return args.Get(0).(*mccsdk.ComponentVersion), args.Error(1)
}

func (c *ClientWrapperMock) GetLatestComponentVersion(ctx context.Context, organization, project, key, branch string) (*mccsdk.ComponentVersion, error) {
	args := c.Called(ctx, organization, project, key, branch)
	return args.Get(0).(*mccsdk.ComponentVersion), args.Error(1)
}

func (c *ClientWrapperMock) PushComponentVersionCommits(ctx context.Context, organization, project, componentKey, versionIdentifier string, commits []mccsdk.CommitData) error {
	args := c.Called(ctx, organization, project, componentKey, versionIdentifier, commits)
	return args.Error(0)
}

func TestRegisterComponentVersionComponentNotFoundWithoutCreateComponent(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, organization, project, componentKey, branch, version, false, false, false, gitFilterPaths)
	assert.ErrorContains(t, err, "Component test-component does not exist")
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "CreateComponentVersion"))
}

func TestRegisterComponentVersionComponentNotFoundWithCreateComponentDryRun(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, organization, project, componentKey, branch, version, true, false, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
	assert.True(t, client.AssertNotCalled(t, "CreateComponentVersion"))
}

func TestRegisterComponentVersionComponentNotFoundWithCreateComponent(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{}, nil)
	client.On("CreateComponent", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.Component{
		Key: utils.Ref("test-component"),
	}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, organization, project, componentKey, branch, version, false, false, true, gitFilterPaths)
	assert.NoError(t, err)
}

func TestRegisterComponentVersionComponentFound(t *testing.T) {
	client := &ClientWrapperMock{}
	client.On("ListComponents", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentPaginator{
		Results: []mccsdk.Component{
			{
				Key: utils.Ref("test-component"),
			},
		},
	}, nil)
	client.On("CreateComponentVersion", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(&mccsdk.ComponentVersion{
		Version: "test-component-version"}, nil)

	ctx := context.Background()
	organization := "test-org"
	project := "test-project"
	componentKey := "test-component"
	branch := "main"
	version := "1.0.0"
	var gitFilterPaths []string

	err := RegisterComponentVersion(ctx, client, organization, project, componentKey, branch, version, false, false, true, gitFilterPaths)
	assert.NoError(t, err)
	assert.True(t, client.AssertNotCalled(t, "CreateComponent"))
}
