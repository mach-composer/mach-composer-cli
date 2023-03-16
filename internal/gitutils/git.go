package gitutils

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/elliotchance/pie/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/labd/mach-composer/internal/config"
	"github.com/labd/mach-composer/internal/utils"
)

type gitSource struct {
	URL        string
	Repository string
	Path       string
	Name       string
}

type gitCommit struct {
	Commit    string
	Parents   []string
	Author    gitCommitAuthor
	Committer gitCommitAuthor
	Message   string
}

type gitCommitAuthor struct {
	Name  string
	Email string
	Date  time.Time
}

func OpenRepository(path string) (*git.Repository, error) {
	gitPath, err := getGitPath(path)
	if err != nil {
		return nil, err
	}
	repo, err := git.PlainOpen(gitPath)
	if err != nil {
		return nil, err
	}

	// Hack to make resolving short hashes work
	// https://github.com/go-git/go-git/issues/148#issuecomment-989635832
	if _, err = repo.CommitObjects(); err != nil {
		return nil, err
	}

	return repo, nil
}

func Commit(ctx context.Context, fileNames []string, message string) error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	repository, err := git.PlainOpen(path)
	if err != nil {
		return err
	}

	w, err := repository.Worktree()
	if err != nil {
		return err
	}

	for _, filename := range fileNames {
		if _, err := w.Add(filename); err != nil {
			return err
		}
	}

	if _, err := w.Commit(message, &git.CommitOptions{}); err != nil {
		return err
	}

	return nil
}

func GetLastVersionGit(ctx context.Context, c *config.Component, origin string) ([]gitCommit, error) {
	cacheDir, err := getGitCachePath(origin)
	if err != nil {
		return nil, err
	}

	source, err := parseGitSource(c.Source)
	if err != nil {
		return nil, fmt.Errorf("cannot check %s component since it doesn't have a Git source defined", c.Name)
	}

	branch := ""
	if c.Branch != "" {
		branch = c.Branch
	}
	fetchGitRepository(ctx, source, cacheDir)
	path := filepath.Join(cacheDir, source.Name)
	commits, err := GetRecentCommits(ctx, path, branch, c.Version, []string{})
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func GetCurrentBranch(ctx context.Context, path string) (string, error) {
	repository, err := OpenRepository(path)
	if err != nil {
		return "", err
	}
	branchRef, err := repository.Head()
	if err != nil {
		err = fmt.Errorf("failed to resolve HEAD in repository (%s): %w", path, err)
		return "", err
	}
	return branchRef.Name().Short(), nil
}

// GetRecentCommits returns all commits in descending order (newest first)
// baseRef is the commit to start from, if empty the current HEAD is used
func GetRecentCommits(ctx context.Context, basePath string, branch string, baseRef string, extraPaths []string) ([]gitCommit, error) {
	gitPath, err := getGitPath(basePath)
	if err != nil {
		return nil, err
	}

	repository, err := OpenRepository(gitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve the base commit in the repository. If it's not found, we'll
	// start from the beginning of the repository.
	var baseRevision *plumbing.Hash
	if baseRef != "" {
		baseRevision, err = repository.ResolveRevision(plumbing.Revision(baseRef))
		if err != nil {
			zerolog.Ctx(ctx).Warn().Err(err).Msgf("failed to find commit %s in repository %s", baseRef, gitPath)
			baseRevision = nil
		}
	}

	// Resolve the last commit in the branch
	var branchRevision *plumbing.Hash
	if branch == "" {
		branchRef, err := repository.Head()
		if err != nil {
			err = fmt.Errorf("failed to resolve HEAD in repository (%s): %w", gitPath, err)
			return nil, err
		}
		hash := branchRef.Hash()
		branchRevision = &hash
	} else {
		branchRef := plumbing.NewBranchReferenceName(branch)
		branchRevision, err = repository.ResolveRevision(plumbing.Revision(branchRef))
		if err != nil {
			return nil, fmt.Errorf("failed to find branch %s in repository: %w", branch, err)
		}
	}

	relevantPaths := []string{basePath}
	relevantPaths = append(relevantPaths, extraPaths...)
	paths, err := filterPaths(gitPath, relevantPaths)
	if err != nil {
		return nil, err
	}

	commits, err := commitsBetween(repository, baseRevision, branchRevision, paths)
	if err != nil {
		return nil, err
	}

	result := make([]gitCommit, len(commits))
	for i, c := range commits {
		fields := strings.Split(c.Message, "\n")
		subject := strings.TrimSpace(fields[0])
		parents := make([]string, len(c.ParentHashes))
		for i, parent := range c.ParentHashes {
			parents[i] = parent.String()[:7]
		}

		result[i] = gitCommit{
			Commit:  c.Hash.String()[:7],
			Parents: parents,
			Author: gitCommitAuthor{
				Name:  c.Author.Name,
				Email: c.Author.Email,
				Date:  c.Author.When,
			},
			Committer: gitCommitAuthor{
				Name:  c.Committer.Name,
				Email: c.Committer.Email,
				Date:  c.Committer.When,
			},
			Message: subject,
		}
	}
	return result, nil
}

// getGitCachePath returns the path to the directory used to clone all the
// git repositories for components referenced from the config file. It's
// used only for checking the last version of a component when running
// `mach composer update`.
func getGitCachePath(origin string) (string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	base := strings.TrimSuffix(origin, filepath.Ext(origin))
	path := filepath.Join(cwd, ".mach", base)
	if err := os.MkdirAll(path, 0700); err != nil {
		return "", err
	}
	return path, nil
}

// Parse a git url and return a gitSource reference
func parseGitSource(source string) (*gitSource, error) {
	re := regexp.MustCompile("^git::(?P<repo>https://.*?)(?://(?P<path>.*))?$")
	match := re.FindStringSubmatch(source)

	if match == nil {
		return nil, errors.New("invalid Git source defined")
	}

	result := &gitSource{
		URL: source,
	}
	for i, name := range re.SubexpNames() {
		if name == "repo" {
			result.Repository = match[i]
		}
		if name == "path" {
			result.Path = match[i]
		}
	}

	parts, err := url.Parse(result.Repository)
	if err != nil {
		return nil, err
	}
	result.Name = filepath.Base(parts.Path)
	return result, nil
}

// fetchGitRepository clones or updates the repository. We only need the history
// so clone using --bare
func fetchGitRepository(ctx context.Context, source *gitSource, cacheDir string) {
	dest := filepath.Join(cacheDir, source.Name)

	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		output := runGit(ctx, ".", "clone", "--bare", source.Repository, dest)
		log.Debug().Msgf(string(output))
	} else {
		output := runGit(ctx, dest, "fetch", "-f", "origin", "*:*")
		log.Debug().Msgf(string(output))
	}
}

func getGitPath(path string) (string, error) {
	// Walk upwards to find a .git directory from the current path
	// This is needed because the current working directory is not always the
	// same as the root of the repository
	if path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	for {
		gitDir := filepath.Join(path, ".git")
		_, err := os.Stat(gitDir)
		if err == nil {
			return path, nil
		}
		if os.IsNotExist(err) {
			path = filepath.Dir(path)
			if path == "/" {
				return "", errors.New("could not find .git directory")
			}
			continue
		}
		return "", err
	}
}

// filterPaths returns the paths to filter on for git commits. It cretes
// relative paths from the gitPath to the paths provided.
func filterPaths(gitPath string, paths []string) ([]string, error) {
	var err error
	gitPath, err = filepath.Abs(gitPath)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, p := range paths {

		absPath, err := filepath.Abs(p)
		if err != nil {
			return nil, err
		}

		// If the path is the same as the git path, then all paths
		// are targets
		if absPath == gitPath {
			return []string{}, nil
		}

		rel, err := filepath.Rel(gitPath, absPath)
		if err != nil {
			return nil, err
		}

		fi, err := os.Stat(absPath)
		if err == nil && fi.IsDir() {
			rel = rel + string(filepath.Separator)
		}

		result = append(result, rel)
	}
	return result, nil
}

// commitsBetween returns the commits from x to y. It should equal the
// functionality of `git log base..head`
// See https://github.com/go-git/go-git/issues/69
// FIXME: very naive implementation
func commitsBetween(repository *git.Repository, first, last *plumbing.Hash, paths []string) ([]*object.Commit, error) {
	cIter, err := repository.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		From:  *last,
		PathFilter: func(path string) bool {
			if len(paths) == 0 {
				return true
			}
			return pie.Any(paths, func(p string) bool {
				return strings.HasPrefix(path, p)
			})
		},
	})
	if err != nil {
		return nil, err
	}

	result := []*object.Commit{}
	err = cIter.ForEach(func(c *object.Commit) error {
		if first != nil && *first == c.Hash {
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

// runGit executes the git command
func runGit(ctx context.Context, cwd string, args ...string) []byte {
	log.Debug().Msgf("Running: git %s\n", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"git",
		args...,
	)
	cmd.Dir = cwd
	utils.CmdSetForeground(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintln(os.Stderr, string(output))
		os.Exit(1)
	}

	return output
}
