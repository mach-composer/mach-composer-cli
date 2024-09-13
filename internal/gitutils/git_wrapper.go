package gitutils

import (
	"context"
)

type GitRepository interface {
	GetCurrentBranch(context.Context, string) (string, error)
	GetVersionInfo(ctx context.Context, path string, branch string) (*GitVersionInfo, error)
	GetRecentCommits(ctx context.Context, basePath string, baseRevision, targetRevision string, filterPaths []string) ([]GitCommit, error)
}

type GitRepositoryWrapper struct {
}

func NewGitRepositoryWrapper() *GitRepositoryWrapper {
	return &GitRepositoryWrapper{}
}

func (g GitRepositoryWrapper) GetCurrentBranch(ctx context.Context, s string) (string, error) {
	return GetCurrentBranch(ctx, s)
}

func (g GitRepositoryWrapper) GetVersionInfo(ctx context.Context, path string, branch string) (*GitVersionInfo, error) {
	return GetVersionInfo(ctx, path, branch)
}

func (g GitRepositoryWrapper) GetRecentCommits(ctx context.Context, basePath string, baseRevision, targetRevision string, filterPaths []string) ([]GitCommit, error) {
	return GetRecentCommits(ctx, basePath, baseRevision, targetRevision, filterPaths)
}
