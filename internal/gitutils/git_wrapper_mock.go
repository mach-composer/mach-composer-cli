//go:build testing

package gitutils

import (
	"context"
	"github.com/stretchr/testify/mock"
)

type GitRepositoryMock struct {
	mock.Mock
}

func (g *GitRepositoryMock) GetCurrentBranch(ctx context.Context, s string) (string, error) {
	args := g.Called(ctx, s)
	return args.String(0), args.Error(1)
}

func (g *GitRepositoryMock) GetVersionInfo(ctx context.Context, path, branch string) (*GitVersionInfo, error) {
	args := g.Called(ctx, path, branch)
	return args.Get(0).(*GitVersionInfo), args.Error(1)
}

func (g *GitRepositoryMock) GetRecentCommits(ctx context.Context, basePath, baseRevision, targetRevision string, filterPaths []string) ([]GitCommit, error) {
	args := g.Called(ctx, basePath, baseRevision, targetRevision, filterPaths)
	return args.Get(0).([]GitCommit), args.Error(1)
}
