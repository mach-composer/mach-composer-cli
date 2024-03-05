package cli

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const GithubCIKey = "github-ci"

func ContextWithGithubCI(ctx context.Context) context.Context {
	return context.WithValue(ctx, GithubCIKey, true)
}

// GithubCIFromContext returns whether the context is running in a GitHub CI
func GithubCIFromContext(ctx context.Context) bool {
	if v := ctx.Value(GithubCIKey); v != nil {
		return v.(bool)
	}

	return false
}

type GitHubHook struct {
	identifier string
}

func NewGitHubHook(identifier string) *GitHubHook {
	return &GitHubHook{
		identifier: identifier,
	}
}

func (g *GitHubHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if GithubCIFromContext(e.GetCtx()) {
		log.WithLevel(level).Msgf("::group::{%s}\n%s\n::endgroup::", g.identifier, message)
		e.Discard()
		return
	}

	log.WithLevel(level).Msg(message)
	e.Discard()
	return
}
