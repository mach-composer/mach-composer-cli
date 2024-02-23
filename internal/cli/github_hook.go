package cli

import (
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type GitHubHook struct {
	identifier string
}

func NewGitHubHook(identifier string) *GitHubHook {
	return &GitHubHook{
		identifier: identifier,
	}
}

func (g *GitHubHook) Run(e *zerolog.Event, level zerolog.Level, message string) {
	if OutputFromContext(e.GetCtx()) == OutputTypeGitHub {
		log.WithLevel(level).Msgf("::group::{%s}\n%s\n::endgroup::", g.identifier, message)
		e.Discard()
		return
	}

	log.WithLevel(level).Msg(message)
	e.Discard()
	return
}
