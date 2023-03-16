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

func TestCommitsBetween(t *testing.T) {
	tr := NewTestRepository("")

	// First commit
	err := tr.addTextfile("test-1.txt", "test")
	require.NoError(t, err)

	firsthash, err := tr.commit("Initial commit")
	require.NoError(t, err)

	// Second commit
	err = tr.addTextfile("test-2.txt", "test")
	require.NoError(t, err)

	secondhash, err := tr.commit("Second commit")
	require.NoError(t, err)

	ctx := context.Background()

	targetRev := plumbing.Revision(secondhash.String())
	baseRev := plumbing.Revision(firsthash.String())

	commits, err := commitsBetween(ctx, tr.repository(), nil, &targetRev, []string{})
	require.NoError(t, err)
	assert.Equal(t, 2, len(commits))

	commits, err = commitsBetween(ctx, tr.repository(), &baseRev, &targetRev, []string{})
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))
}

func TestCommitsBetweenFilterPath(t *testing.T) {
	tr := NewTestRepository("")

	// First commit
	err := tr.addTextfile("test-1.txt", "test")
	require.NoError(t, err)

	firsthash, err := tr.commit("Initial commit")
	require.NoError(t, err)
	require.NotNil(t, firsthash)

	// Second commit
	err = tr.addTextfile("test-2.txt", "test")
	require.NoError(t, err)

	err = tr.addTextfile("wanted/test-2.txt", "test")
	require.NoError(t, err)

	_, err = tr.commit("Second commit")
	require.NoError(t, err)

	// Third commit
	err = tr.addTextfile("test-3.txt", "test")
	require.NoError(t, err)

	err = tr.addTextfile("wantedotherdir/test-3.txt", "test")
	require.NoError(t, err)

	thirdhash, err := tr.commit("third commit")
	require.NoError(t, err)

	ctx := context.Background()
	targetRev := plumbing.Revision(thirdhash.String())

	// Check results
	commits, err := commitsBetween(ctx, tr.repository(), nil, &targetRev, []string{})
	require.NoError(t, err)
	assert.Equal(t, 3, len(commits))

	commits, err = commitsBetween(ctx, tr.repository(), nil, &targetRev, []string{"wanted/"})
	require.NoError(t, err)
	assert.Equal(t, 1, len(commits))
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

func TestFilterPaths(t *testing.T) {
	type input struct {
		gitPath string
		paths   []string
	}
	type output struct {
		relPaths []string
		err      error
	}
	testCases := []struct {
		name  string
		input input
		want  output
	}{
		{
			name: "happy path",
			input: input{
				gitPath: "/home/user/repo",
				paths:   []string{"/home/user/repo/file1.txt", "/home/user/repo/dir1/file2.txt"},
			},
			want: output{
				relPaths: []string{"file1.txt", "dir1/file2.txt"},
				err:      nil,
			},
		},
		{
			name: "empty path list",
			input: input{
				gitPath: "/home/user/repo",
				paths:   []string{},
			},
			want: output{
				relPaths: []string{},
				err:      nil,
			},
		},
		{
			name: "git path is not absolute",
			input: input{
				gitPath: "../",
				paths:   []string{"../plugins/config.go", "../generator", "dir1/file2.txt"},
			},
			want: output{
				relPaths: []string{
					"plugins/config.go",
					"generator/",
					"gitutils/dir1/file2.txt",
				},
			},
		},
		{
			name: "no paths",
			input: input{
				gitPath: "/home/user/repo",
				paths:   []string{"/home/user/repo"},
			},
			want: output{
				relPaths: []string{},
				err:      nil,
			},
		},
		{
			name: "no paths (ignore extra)",
			input: input{
				gitPath: "/home/user/repo",
				paths:   []string{"/home/user/repo", "/home/user/repo/dir1"},
			},
			want: output{
				relPaths: []string{},
				err:      nil,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := filterPaths(tc.input.gitPath, tc.input.paths)
			assert.Equal(t, tc.want.err, err)
			assert.EqualValues(t, tc.want.relPaths, got)
		})
	}
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

	h := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.ReferenceName("refs/heads/main"))
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
