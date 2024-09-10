//go:build testing

package cloud

import (
	"context"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
	"github.com/stretchr/testify/mock"
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
