package gitutils

import (
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-billy/v5/osfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/cache"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/storage/filesystem"
	"github.com/go-git/go-git/v5/storage/memory"
	"path/filepath"
	"time"
)

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

func (t *TestRepository) addTextFile(filename string, content string) error {
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
