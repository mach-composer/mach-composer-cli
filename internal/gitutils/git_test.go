package gitutils

import (
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseGitSource(t *testing.T) {
	params := []struct {
		source   string
		expected *gitSource
	}{
		{
			source: "git::https://github.com/labd/mach-composer",
			expected: &gitSource{
				Repository: "https://github.com/labd/mach-composer",
				Path:       "",
				Name:       "mach-composer",
			},
		},
		{
			source: "git::https://github.com/labd/mach-composer//terraform",
			expected: &gitSource{
				Repository: "https://github.com/labd/mach-composer",
				Path:       "terraform",
				Name:       "mach-composer",
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

func TestCommitsBetwee(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)
	require.NotNil(t, repo)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	// First commit
	err = addNewFile(worktree, "test-1.txt")
	require.NoError(t, err)

	firsthash, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	// Second commit
	err = addNewFile(worktree, "test-2.txt")
	require.NoError(t, err)

	secondhash, err := worktree.Commit("Second commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	commits, err := commitsBetween(repo, nil, &secondhash, []string{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(commits))

	commits, err = commitsBetween(repo, &firsthash, &secondhash, []string{})
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))
}

func TestCommitsBetweenFilterPath(t *testing.T) {
	repo, err := git.Init(memory.NewStorage(), memfs.New())
	require.NoError(t, err)
	require.NotNil(t, repo)

	worktree, err := repo.Worktree()
	require.NoError(t, err)

	// First commit
	err = addNewFile(worktree, "test-1.txt")
	require.NoError(t, err)

	firsthash, err := worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)
	require.NotNil(t, firsthash)

	// Second commit
	err = addNewFile(worktree, "test-2.txt")
	require.NoError(t, err)

	err = addNewFile(worktree, "wanted/test-2.txt")
	require.NoError(t, err)

	_, err = worktree.Commit("Second commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	// Third commit
	err = addNewFile(worktree, "test-3.txt")
	require.NoError(t, err)

	err = addNewFile(worktree, "wantedotherdir/test-3.txt")
	require.NoError(t, err)

	thirdhash, err := worktree.Commit("third commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.org",
			When:  time.Now(),
		},
	})
	require.NoError(t, err)

	// Check results
	commits, err := commitsBetween(repo, nil, &thirdhash, []string{})
	require.NoError(t, err)
	assert.Equal(t, 3, len(commits))

	commits, err = commitsBetween(repo, nil, &thirdhash, []string{"wanted/"})
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))
}

func addNewFile(w *git.Worktree, filename string) error {
	file, err := w.Filesystem.Create(filename)
	if err != nil {
		return err
	}

	text := []byte("Text file 1")
	if _, err = file.Write(text); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	_, err = w.Add(filename)
	return err
}
