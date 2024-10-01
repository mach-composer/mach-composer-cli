package cloud

import (
	"context"
	"github.com/mach-composer/mcc-sdk-go/mccsdk"
)

var _ ClientWrapper = (*MccSdkClientWrapper)(nil)

// ClientWrapper is an interface wrapper for the mccsdk.APIClient so we can more easily write tests
type ClientWrapper interface {
	ListComponents(ctx context.Context, organization, project string, limit int32) (*mccsdk.ComponentPaginator, error)
	CreateComponent(ctx context.Context, organization, project, key string) (*mccsdk.Component, error)
	CreateComponentVersion(ctx context.Context, organization, project, key, version, branch string) (*mccsdk.ComponentVersion, error)
	GetLatestComponentVersion(ctx context.Context, organization, project, key, branch string) (*mccsdk.ComponentVersion, error)
	PushComponentVersionCommits(ctx context.Context, organization, project, componentKey, versionIdentifier string, commits []mccsdk.CommitDraft) error
}

func NewClientWrapper(client *mccsdk.APIClient) *MccSdkClientWrapper {
	return &MccSdkClientWrapper{client: client}
}

type MccSdkClientWrapper struct {
	client *mccsdk.APIClient
}

func (m *MccSdkClientWrapper) PushComponentVersionCommits(ctx context.Context, organization, project, componentKey, versionIdentifier string, commits []mccsdk.CommitDraft) error {
	_, _, err := m.client.
		ComponentsApi.
		ComponentVersionPushCommits(ctx, organization, project, componentKey, versionIdentifier).
		ComponentCommitCreateDraft(mccsdk.ComponentCommitCreateDraft{
			Commits: commits,
		}).
		Execute()

	return err
}

func (m *MccSdkClientWrapper) GetLatestComponentVersion(ctx context.Context, organization, project, componentKey, branch string) (*mccsdk.ComponentVersion, error) {
	r, _, err := m.client.
		ComponentsApi.
		ComponentLatestVersion(ctx, organization, project, componentKey).
		Branch(branch).
		Execute()

	return r, err
}

func (m *MccSdkClientWrapper) CreateComponentVersion(ctx context.Context, organization, project, componentKey, version, branch string) (*mccsdk.ComponentVersion, error) {
	r, _, err := m.client.
		ComponentsApi.
		ComponentVersionCreate(ctx, organization, project, componentKey).
		ComponentVersionDraft(mccsdk.ComponentVersionDraft{
			Version: version,
			Branch:  &branch,
		}).
		Execute()
	return r, err
}

func (m *MccSdkClientWrapper) CreateComponent(ctx context.Context, organization, project, componentKey string) (*mccsdk.Component, error) {
	r, _, err := m.client.
		ComponentsApi.
		ComponentCreate(ctx, organization, project).
		ComponentDraft(mccsdk.ComponentDraft{
			Name: componentKey,
			Key:  componentKey,
		}).
		Execute()

	return r, err
}

func (m *MccSdkClientWrapper) ListComponents(ctx context.Context, organization, project string, limit int32) (*mccsdk.ComponentPaginator, error) {
	res, _, err := m.client.ComponentsApi.ComponentQuery(ctx, organization, project).Limit(limit).Execute()
	return res, err
}
