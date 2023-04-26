package gitutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/format/diff"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/rs/zerolog"
)

// commitsBetween returns the commits between revisions first and last. It
// should equal the functionality of `git log base..head`
// See https://github.com/go-git/go-git/issues/69
func commitsBetween(ctx context.Context, name string, repository *git.Repository, first, last *plumbing.Revision, paths []string) ([]*object.Commit, error) {
	zerolog.Ctx(ctx).Debug().Msgf("Getting %s commits between %s and %s (paths = %s)", name, first, last, paths)
	if first != nil {
		_, err := repository.ResolveRevision(*first)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).Msgf("%s: failed to find commit %s in repository", name, first)
			return nil, ErrGitRevisionNotFound
		}
	}

	// Resolve the base commit in the repository. If it's not found, we'll
	// start from the beginning of the repository.
	var firstHash, lastHash *plumbing.Hash
	if first != nil {
		if val, err := repository.ResolveRevision(*first); err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Msgf("%s: failed to resolve %s in repository", name, first.String())
			return []*object.Commit{}, nil
		} else {
			firstHash = val
		}
	}

	if last == nil {
		last = asRevision("HEAD")
	}

	if val, err := repository.ResolveRevision(*last); err != nil {
		remoteRevision := asRevision(fmt.Sprintf("remotes/origin/%s", *last))
		if val, err := repository.ResolveRevision(*remoteRevision); err != nil {
			return []*object.Commit{}, fmt.Errorf("%s: failed to resolve %s in repository", name, last.String())
		} else {
			lastHash = val
		}
	} else {
		lastHash = val
	}

	cIter, err := repository.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		From:  *lastHash,
	})
	if err != nil {
		return nil, err
	}

	result := []*object.Commit{}
	err = cIter.ForEach(func(c *object.Commit) error {
		if first != nil && *firstHash == c.Hash {
			return storer.ErrStop
		}

		if matchPaths(ctx, repository, c, paths) {
			result = append(result, c)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// matchPaths returns true if the commit matches the paths provided. If no
// paths are provided, then it always returns true.
func matchPaths(ctx context.Context, r *git.Repository, commit *object.Commit, paths []string) bool {
	if len(paths) == 0 {
		return true
	}

	for _, ph := range commit.ParentHashes {
		p, err := r.CommitObject(ph)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).Msgf("failed to find commit %s in repository", ph)
			return false
		}

		patch, err := p.Patch(commit)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).Msgf("failed to generate patch")
			return false
		}

		for _, fp := range patch.FilePatches() {
			from, to := fp.Files()

			if matchFilePatch(from, paths) {
				return true
			}

			if matchFilePatch(to, paths) {
				return true
			}
		}
	}

	// If there are no parents, then this is the first commit. We need to check
	// the tree.
	if len(commit.ParentHashes) == 0 {
		fIter, err := commit.Files()
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).Msgf("failed to get files for commit %s", commit.Hash)
			return false
		}

		for {
			f, err := fIter.Next()
			if err != nil {
				break
			}
			for _, p := range paths {
				if strings.HasPrefix(f.Name, p) {
					return true
				}
			}
		}
	}

	return false
}

func matchFilePatch(f diff.File, paths []string) bool {
	if f == nil {
		return false
	}
	path := f.Path()
	return pie.Any(paths, func(p string) bool { return strings.HasPrefix(path, p) })
}
