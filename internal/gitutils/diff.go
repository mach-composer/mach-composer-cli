package gitutils

import (
	"context"
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/rs/zerolog/log"
)

func pathFilter(paths []string) func(path string) bool {
	return func(path string) bool {
		if len(paths) == 0 {
			return true
		}

		for _, p := range paths {
			if strings.HasPrefix(path, p) {
				return true
			}
		}
		return false
	}
}

// commitsBetween returns the commits between revisions first and last. It should equal the functionality of
// `git log base..head`. See https://github.com/go-git/go-git/issues/69
func commitsBetween(ctx context.Context, repository *git.Repository, first, last *plumbing.Revision, paths []string) ([]*object.Commit, error) {
	log.Ctx(ctx).Debug().Msgf("Getting commits between %s and %s (paths = %s)", first, last, paths)
	if first != nil {
		_, err := repository.ResolveRevision(*first)
		if err != nil {
			return nil, ErrGitRevisionNotFound
		}
	}

	// Resolve the base commit in the repository. If it's not found, we'll
	// start from the beginning of the repository.
	var firstHash, lastHash *plumbing.Hash
	if first != nil {
		if val, err := repository.ResolveRevision(*first); err != nil {
			log.Ctx(ctx).Warn().Err(err).Msgf("failed to resolve %s in repository", first.String())
			return []*object.Commit{}, nil
		} else {
			firstHash = val
		}
	}

	firstCommit, err := repository.CommitObject(*firstHash)
	if err != nil {
		return nil, err
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
		Order:      git.LogOrderCommitterTime,
		PathFilter: pathFilter(paths),
		From:       *lastHash,
		//Since:      &firstCommit.Committer.When,
	})
	if err != nil {
		return nil, err
	}
	defer cIter.Close()

	var result []*object.Commit
	var found = first == nil
	err = cIter.ForEach(func(c *object.Commit) error {
		if first != nil && *firstHash == c.Hash {
			found = true
			return storer.ErrStop
		}
		if firstCommit.Committer.When.After(c.Committer.When) {
			found = false
			log.Ctx(ctx).Info().Msgf("Did not find commit %s in path but next older commit %s in paths %s", first, c.Hash, paths)
			return storer.ErrStop
		}

		result = append(result, c)
		return nil
	})
	if err != nil {
		return nil, err
	}
	if !found {
		log.Ctx(ctx).Warn().Msgf("found commit %s in %s but failed to find changes in paths %s", first, last, paths)
	}

	return result, nil
}
