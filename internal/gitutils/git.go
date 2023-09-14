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

var downloadFiles = utils.OnceMap[string]{}

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

type GitVersionInfo struct {
	Hash     plumbing.Hash
	Tag      string
	Revision plumbing.Revision
}

func (g *GitVersionInfo) Identifier() string {
	return g.Hash.String()[0:7]
}

func OpenRepository(path string) (*git.Repository, error) {
	repo, err := git.PlainOpen(path)
	if errors.Is(err, git.ErrRepositoryNotExists) {
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

func Commit(_ context.Context, fileNames []string, message string) error {
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

func GetLastVersionGit(ctx context.Context, c *config.ComponentConfig, origin string) ([]gitCommit, error) {
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

	path := filepath.Join(cacheDir, source.Name)

	downloadFiles.Get(path).Do(func() {
		fetchGitRepository(ctx, source, path)
	})

	commits, err := GetRecentCommits(ctx, path, c.Version, branch, c.Paths)
	if err != nil {
		return nil, err
	}

	return commits, nil
}

func GetCurrentBranch(_ context.Context, path string) (string, error) {
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
func GetVersionInfo(_ context.Context, path string, branch string) (*GitVersionInfo, error) {
	repository, err := OpenRepository(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	// Resolve the last commit in the branch
	info := &GitVersionInfo{}
	if branch == "" {
		branchRef, err := repository.Head()
		if err != nil {
			err = fmt.Errorf("failed to resolve HEAD in repository (%s): %w", path, err)
			return nil, err
		}
		info.Hash = branchRef.Hash()
		info.Revision = "HEAD"
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

func commitExists(ctx context.Context, gitPath, baseRevision, targetRevision string) error {
	containsCommit, err := branchContainsCommit(ctx, gitPath, targetRevision, baseRevision)
	if err != nil {
		return err
	}

	if !containsCommit {
		branches, err := branchesContainingCommit(ctx, gitPath, baseRevision)
		if err != nil {
			return err
		}
		found := "not found in any branch"
		if len(branches) > 0 {
			found = fmt.Sprintf("exists in %s", strings.Join(branches, ", "))
		}
		return fmt.Errorf("failed to find commit %s in branch %s (%s)", baseRevision, targetRevision, found)
	}

	return nil
}

// GetRecentCommits returns all commits in descending order (newest first)
// baseRef is the commit to start from, if empty the current HEAD is used
func GetRecentCommits(ctx context.Context, basePath string, baseRevision, targetRevision string, filterPaths []string) ([]gitCommit, error) {
	gitPath, err := searchGitPath(basePath)
	if err != nil {
		return nil, err
	}

	repository, err := OpenRepository(gitPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	baseRev := asRevision(baseRevision)
	targetRev := asRevision(targetRevision)

	if err = commitExists(ctx, gitPath, baseRevision, targetRevision); err != nil {
		return nil, err
	}

	commits, err := commitsBetween(ctx, repository, baseRev, targetRev, filterPaths)
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

// fetchGitRepository clones or updates the repository. We only need the history so clone using --bare
func fetchGitRepository(ctx context.Context, source *gitSource, dest string) {
	_, err := os.Stat(dest)
	if os.IsNotExist(err) {
		output, _ := runGit(ctx, ".", "clone", "--bare", source.Repository, dest)
		log.Ctx(ctx).Debug().
			Msgf("downloaded new repository %s at destination %s: %s", source.Name, dest, string(output))
	} else {
		output, _ := runGit(ctx, dest, "fetch", "-f", "origin", "*:*")
		log.Ctx(ctx).Debug().
			Msgf("updated existing repository %s at destination %s: %s", source.Name, dest, string(output))
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

		if !errors.Is(err, git.ErrRepositoryNotExists) {
			return "", fmt.Errorf("failed to find open repository: %w", err)
		}

		path = filepath.Dir(path)
		if path == "/" {
			break
		}
	}

	return "", errors.New("could not find .git directory")
}

// runGit executes the git command
func runGit(ctx context.Context, cwd string, args ...string) ([]byte, error) {
	log.Debug().Msgf("Running: git %s", strings.Join(args, " "))
	cmd := exec.CommandContext(
		ctx,
		"git",
		args...,
	)
	cmd.Dir = cwd
	utils.CmdSetForeground(cmd)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("failed to run git command: %s", output)
	}

	return output, nil
}

func asRevision(s string) *plumbing.Revision {
	if s == "" {
		return nil
	}
	r := plumbing.Revision(s)
	return &r
}

func branchContainsCommit(ctx context.Context, gitPath, targetRev, baseRev string) (bool, error) {
	output, err := runGit(ctx, gitPath, "branch", "-a", targetRev, "--contains", baseRev)
	if err != nil {
		return false, fmt.Errorf(strings.TrimSpace(err.Error()))
	}

	return len(output) > 0, nil
}

func branchesContainingCommit(ctx context.Context, gitPath, baseRev string) ([]string, error) {
	output, err := runGit(ctx, gitPath, "branch", "--contains", baseRev)
	if err != nil {
		return nil, err
	}

	var branches []string
	outputBranches := strings.Split(string(output), "\n")
	for _, branch := range outputBranches {
		if branch == "" {
			continue
		}
		branches = append(branches, strings.TrimSpace(branch))
	}

	return branches, nil
}
