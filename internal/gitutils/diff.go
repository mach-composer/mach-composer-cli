package gitutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/rs/zerolog"
)

// commitsBetween returns the commits between revisions first and last. It
// should equal the functionality of `git log base..head`
// See https://github.com/go-git/go-git/issues/69
func commitsBetween(ctx context.Context, repository *git.Repository, first, last *plumbing.Revision, paths []string) ([]*object.Commit, error) {
	zerolog.Ctx(ctx).Debug().Msgf("Getting commits between %s and %s (paths = %s)", first, last, paths)
	if first != nil {
		_, err := repository.ResolveRevision(*first)
		if err != nil {
			zerolog.Ctx(ctx).Debug().Err(err).Msgf("failed to find commit %s in repository", first)
			return nil, ErrGitRevisionNotFound
		}
	}

	// Resolve the base commit in the repository. If it's not found, we'll
	// start from the beginning of the repository.
	var firstHash, lastHash *plumbing.Hash
	if first != nil {
		if val, err := repository.ResolveRevision(*first); err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Msgf("failed to resolve %s in repository", first.String())
			return []*object.Commit{}, nil
		} else {
			firstHash = val
		}
	}

	if last == nil {
		last = asRevision("HEAD")
	}

	if val, err := repository.ResolveRevision(*last); err != nil {
		return []*object.Commit{}, fmt.Errorf("failed to resolve %s in repository", last.String())
	} else {
		lastHash = val
	}

	cIter, err := repository.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		PathFilter: func(path string) bool {
			for _, p := range paths {
				if strings.HasPrefix(path, p) {
					return true
				}
			}
			return false
		},
		From: *lastHash,
	})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	var result []*object.Commit
	err = cIter.ForEach(func(c *object.Commit) error {
		if first != nil && *firstHash == c.Hash {
			return storer.ErrStop
		}

		result = append(result, c)
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}
