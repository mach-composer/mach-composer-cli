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

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

var ErrGitRevisionNotFound = errors.New("git revision not found")

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

type gitVersionInfo struct {
	Hash     plumbing.Hash
	Tag      string
	Revision plumbing.Revision
}

func (g *gitVersionInfo) Identifier() string {
	return g.Hash.String()[0:7]
}

func OpenRepository(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if err == git.ErrRepositoryNotExists {
		gitPath, err := searchGitPath(path)
		if err != nil {
			return nil, err
		}
		repo, err = git.PlainOpen(gitPath)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
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
	commits, err := GetRecentCommits(ctx, path, c.Version, branch, []string{})
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

// GetVersionInfo returns the latest commit hash of a specific branch
func GetVersionInfo(ctx context.Context, path string, branch string) (*gitVersionInfo, error) {
	repository, err := OpenRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve the last commit in the branch
	info := &gitVersionInfo{}
	if branch == "" {
		branchRef, err := repository.Head()
		if err != nil {
			err = fmt.Errorf("failed to resolve HEAD in repository (%s): %w", path, err)
			return nil, err
		}
		info.Hash = branchRef.Hash()
		info.Revision = plumbing.Revision("HEAD")
	} else {
		branchRef := plumbing.NewBranchReferenceName(branch)
		branchRevision, err := repository.ResolveRevision(plumbing.Revision(branchRef))
		if err != nil {
			return nil, fmt.Errorf("failed to find branch %s in repository: %w", branch, err)
		}
		info.Hash = *branchRevision
		info.Revision = plumbing.Revision(branchRef)
	}
	return info, nil
}

// GetRecentCommits returns all commits in descending order (newest first)
// baseRef is the commit to start from, if empty the current HEAD is used
func GetRecentCommits(ctx context.Context, basePath string, baseRevision, targetRevision string, extraPaths []string) ([]gitCommit, error) {
	gitPath, err := searchGitPath(basePath)
	if err != nil {
		return nil, err
	}

	repository, err := OpenRepository(gitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	paths, err := filterPaths(gitPath, extraPaths)
	if err != nil {
		return nil, err
	}

	baseRev := asRevision(baseRevision)
	targetRev := asRevision(targetRevision)
	commits, err := commitsBetween(ctx, repository, baseRev, targetRev, paths)
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

// Walk upwards to find a .git directory from the current path
// This is needed because the current working directory is not always the
// same as the root of the repository
func searchGitPath(path string) (string, error) {

	if path == "." {
		var err error
		path, err = os.Getwd()
		if err != nil {
			return "", err
		}
	}

	for {
		_, err := git.PlainOpen(path)
		if err == nil {
			return path, nil
		}

		if err != git.ErrRepositoryNotExists {
			return "", fmt.Errorf("failed to find open repository: %w", err)
		}

		path = filepath.Dir(path)
		if path == "/" {
			break
		}
	}

	return "", errors.New("could not find .git directory")
}

// filterPaths returns the paths to filter on for git commits. It creates
// relative paths from the gitPath to the paths provided.
func filterPaths(gitPath string, paths []string) ([]string, error) {
	var err error
	gitPath, err = filepath.Abs(gitPath)
	if err != nil {
		return nil, err
	}

	result := []string{}
	for _, p := range paths {
		var absPath string
		if filepath.IsAbs(p) {
			absPath = p
		} else {
			ap, err := filepath.Abs(p)
			if err != nil {
				return nil, err
			}
			absPath = ap
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

func asRevision(s string) *plumbing.Revision {
	if s == "" {
		return nil
	}
	r := plumbing.Revision(s)
	return &r
}
