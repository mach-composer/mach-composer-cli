package gitutils

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
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

	err := tr.addTextfile("test-1.txt", "test-1")
	require.NoError(t, err)
	_, err = tr.commit("Initial commit")
	require.NoError(t, err)

	err = tr.addTextfile("test-2.txt", "test-2")
	require.NoError(t, err)
	_, err = tr.commit("Second commit")
	require.NoError(t, err)

	_, err = tr.r.Head()
	require.NoError(t, err)

	commits, err := GetRecentCommits(context.Background(), path, lastVersion, "", []string{})
	require.ErrorIs(t, err, ErrGitRevisionNotFound)
	require.Nil(t, commits)
}

func TestGetRecentCommitsValidStart(t *testing.T) {
	path := t.TempDir()
	tr := NewTestRepository(path)

	err := tr.addTextfile("test-1.txt", "test-1")
	require.NoError(t, err)
	h, err := tr.commit("Initial commit")
	require.NoError(t, err)

	err = tr.addTextfile("test-2.txt", "test-2")
	require.NoError(t, err)
	_, err = tr.commit("Second commit")
	require.NoError(t, err)

	_, err = tr.r.Head()
	require.NoError(t, err)

	lastVersion := h.String()

	commits, err := GetRecentCommits(context.Background(), path, lastVersion, "", []string{})
	require.NoError(t, err)

	require.NotNil(t, commits)
	assert.Equal(t, 1, len(commits))
}

type TestRepository struct {
	r *git.Repository
	w *git.Worktree
}

func NewTestRepository(path string) *TestRepository {
	var r *git.Repository
	if path == "" {
		var err error
		r, err = git.Init(memory.NewStorage(), memfs.New())
		if err != nil {
			panic(err)
		}
	} else {
		var err error
		fs := osfs.New(path)
		fss := osfs.New(filepath.Join(path, ".git"))
		r, err = git.Init(filesystem.NewStorage(fss, cache.NewObjectLRUDefault()), fs)
		if err != nil {
			panic(err)
		}

	}

	if err := r.CreateBranch(&config.Branch{
		Name:   "main",
		Remote: "origin",
		Merge:  "refs/heads/main",
	}); err != nil {
		panic(err)
	}

	h := plumbing.NewSymbolicReference(plumbing.HEAD, "refs/heads/main")
	if err := r.Storer.SetReference(h); err != nil {
		panic(err)
	}

	w, err := r.Worktree()
	if err != nil {
		panic(err)
	}

	return &TestRepository{
		r: r,
		w: w,
	}
}

func (t *TestRepository) repository() *git.Repository {
	return t.r
}

func (t *TestRepository) commit(message string) (plumbing.Hash, error) {
	return t.w.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test",
			Email: "test@example.org",
			When:  time.Now(),
		},
	})
}

func (t *TestRepository) addTextfile(filename string, content string) error {
	file, err := t.w.Filesystem.Create(filename)
	if err != nil {
		return err
	}

	text := []byte(content)
	if _, err = file.Write(text); err != nil {
		return err
	}

	if err := file.Close(); err != nil {
		return err
	}

	_, err = t.w.Add(filename)
	return err
}
