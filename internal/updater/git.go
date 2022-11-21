package updater

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

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/sirupsen/logrus"

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

func GetLastVersionGit(ctx context.Context, c *config.Component, origin string) (*ChangeSet, error) {
	cacheDir := getGitCachePath(origin)
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
	commits, err := GetRecentCommits(ctx, path, branch, c.Version)
	if err != nil {
		return nil, err
	}

	cs := &ChangeSet{
		Changes:   commits,
		Component: c,
	}

	if len(commits) < 1 {
		cs.LastVersion = c.Version
	} else {
		cs.LastVersion = commits[0].Commit
	}

	return cs, nil
}

func getGitCachePath(origin string) string {
	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	base := strings.TrimSuffix(origin, filepath.Ext(origin))
	path := filepath.Join(cwd, ".mach", base)
	if err := os.MkdirAll(path, 0700); err != nil {
		panic(err)
	}
	return path
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
		panic(err)
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
		logrus.Debug(string(output))
	} else {
		output := runGit(ctx, dest, "fetch", "-f", "origin", "*:*")
		logrus.Debug(string(output))
	}
}

func GetCurrentBranch(ctx context.Context, path string) (string, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		return "", err
	}
	branchRef, err := repository.Head()
	if err != nil {
		return "", fmt.Errorf("failed to resolve HEAD in repository: %w", err)
	}
	return branchRef.Name().Short(), nil
}

// GetRecentCommits returns all commits in descending order (newest first)
func GetRecentCommits(ctx context.Context, path string, branch string, baseRef string) ([]gitCommit, error) {
	repository, err := git.PlainOpen(path)
	if err != nil {
		return nil, err
	}

	// Hack to make resolving short hashes work
	// https://github.com/go-git/go-git/issues/148#issuecomment-989635832
	_, err = repository.CommitObjects()
	if err != nil {
		return nil, err
	}

	var baseRevision *plumbing.Hash
	if baseRef != "" {
		baseRevision, err = repository.ResolveRevision(plumbing.Revision(baseRef))
		if err != nil {
			return nil, fmt.Errorf("failed to find commit %s in repository %s: %w", baseRef, path, err)
		}
	}

	var branchRevision *plumbing.Hash
	if branch == "" {
		branchRef, err := repository.Head()
		if err != nil {
			return nil, fmt.Errorf("failed to resolve HEAD in repository: %w", err)
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

	commits, err := commitsBetween(repository, baseRevision, branchRevision)
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

// commitsBetween returns the commits from x to y. It should equal the
// functionality of `git log base..head`
// See https://github.com/go-git/go-git/issues/69
// FIXME: very naive implementation
func commitsBetween(repository *git.Repository, first, last *plumbing.Hash) ([]*object.Commit, error) {
	cIter, err := repository.Log(&git.LogOptions{
		Order: git.LogOrderCommitterTime,
		From:  *last,
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

// runGit executes the git command
func runGit(ctx context.Context, cwd string, args ...string) []byte {
	logrus.Debugf("Running: git %s\n", strings.Join(args, " "))
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
