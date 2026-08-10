package gitutils

import (
	"context"
)

type GitRepository interface {
	GetCurrentBranch(context.Context, string) (string, error)
	GetLatestCommitHash(ctx context.Context, path string, branch string) (string, error)
}

type GitRepositoryWrapper struct {
}

func NewGitRepositoryWrapper() *GitRepositoryWrapper {
	return &GitRepositoryWrapper{}
}

func (g GitRepositoryWrapper) GetCurrentBranch(ctx context.Context, s string) (string, error) {
	return GetCurrentBranch(ctx, s)
}

func (g GitRepositoryWrapper) GetLatestCommitHash(ctx context.Context, path string, branch string) (string, error) {
	return GetLatestCommitHash(ctx, path, branch)
}
