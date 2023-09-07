package gitutils

import (
	"context"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseGitSource(t *testing.T) {
	params := []struct {
		source   string
		expected *gitSource
	}{
		{
			source: "git::https://github.com/mach-composer/mach-composer-cli",
			expected: &gitSource{
				Repository: "https://github.com/mach-composer/mach-composer-cli",
				Path:       "",
				Name:       "mach-composer-cli",
			},
		},
		{
			source: "git::https://github.com/mach-composer/mach-composer-cli//terraform",
			expected: &gitSource{
				Repository: "https://github.com/mach-composer/mach-composer-cli",
				Path:       "terraform",
				Name:       "mach-composer-cli",
			},
		},
	}
	for _, p := range params {
		p.expected.URL = p.source

		res, err := parseGitSource(p.source)
		assert.NoError(t, err)
		assert.EqualValues(t, p.expected, res)
	}
}

func TestGetRecentCommitsInvalidStart(t *testing.T) {
	lastVersion := "2f52d4"
	path := t.TempDir()
	tr := NewTestRepository(path)

	err := tr.addTextFile("test-1.txt", "test-1")
	require.NoError(t, err)
	_, err = tr.commit("Initial commit")
	require.NoError(t, err)

	err = tr.addTextFile("test-2.txt", "test-2")
	require.NoError(t, err)
	_, err = tr.commit("Second commit")
	require.NoError(t, err)

	head, err := tr.r.Head()
	require.NoError(t, err)

	commits, err := GetRecentCommits(context.Background(), path, lastVersion, head.Name().Short(), []string{})
	require.Error(t, err)
	require.Nil(t, commits)
}

func TestGetRecentCommitsValidStart(t *testing.T) {
	path := t.TempDir()
	tr := NewTestRepository(path)

	err := tr.addTextFile("test-1.txt", "test-1")
	require.NoError(t, err)
	h, err := tr.commit("Initial commit")
	require.NoError(t, err)

	err = tr.addTextFile("test-2.txt", "test-2")
	require.NoError(t, err)
	_, err = tr.commit("Second commit")
	require.NoError(t, err)

	head, err := tr.r.Head()
	require.NoError(t, err)

	lastVersion := h.String()

	commits, err := GetRecentCommits(context.Background(), path, lastVersion, head.Name().Short(), []string{})
	require.NoError(t, err)

	require.NotNil(t, commits)
	assert.Equal(t, 1, len(commits))
}

func TestCommitExistsTrue(t *testing.T) {
	path := t.TempDir()
	tr := NewTestRepository(path)

	err := tr.addTextFile("test-1.txt", "test-1")
	require.NoError(t, err)
	h, err := tr.commit("Initial commit")
	require.NoError(t, err)

	head, err := tr.r.Head()
	require.NoError(t, err)

	lastVersion := h.String()

	err = commitExists(context.Background(), path, lastVersion, head.Name().Short())
	require.NoError(t, err)
}

func TestCommitExistsFalse(t *testing.T) {
	path := t.TempDir()
	tr := NewTestRepository(path)

	err := tr.addTextFile("test-1.txt", "test-1")
	require.NoError(t, err)

	h, err := tr.commit("Initial commit")
	require.NoError(t, err)

	mainHead, err := tr.r.Head()

	err = tr.repository().CreateBranch(&config.Branch{
		Name:  "test",
		Merge: plumbing.NewBranchReferenceName("test"),
	})
	require.NoError(t, err)

	err = tr.w.Checkout(&git.CheckoutOptions{
		Branch: plumbing.NewBranchReferenceName("test"),
		Create: true,
	})
	require.NoError(t, err)

	err = tr.addTextFile("test-2.txt", "test-2")
	require.NoError(t, err)

	h, err = tr.commit("Commit in test")
	require.NoError(t, err)

	_, err = tr.r.Head()
	require.NoError(t, err)

	lastVersion := h.String()

	err = commitExists(context.Background(), path, lastVersion, mainHead.Name().Short())
	require.ErrorContains(t, err, "in branch main (exists in * test)")
}
