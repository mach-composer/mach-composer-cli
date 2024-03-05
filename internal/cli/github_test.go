package cli

import (
	"context"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRun_Discarded_NoGithubOutput(t *testing.T) {
	g := NewGitHubHook("test")

	e := &zerolog.Event{}
	e.Ctx(ContextWithOutput(context.Background(), OutputTypeConsole))

	sink := &LogSink{}

	defer SetSinkLogger(sink)()

	g.Run(e, zerolog.InfoLevel, "test-message")

	assert.Equal(t, 1, len(sink.logs))
	assert.Equal(t, "test-message", sink.Index(0).Message)
	assert.Equal(t, "info", sink.Index(0).Level)
}

func TestRun_Enabled_GithubOutput(t *testing.T) {
	g := NewGitHubHook("test")

	e := &zerolog.Event{}
	e.Ctx(ContextWithGithubCI(context.Background()))

	sink := &LogSink{}

	defer SetSinkLogger(sink)()

	g.Run(e, zerolog.InfoLevel, "test-message")

	assert.Equal(t, 1, len(sink.logs))
	assert.Equal(t, "::group::{test}\ntest-message\n::endgroup::", sink.Index(0).Message)
	assert.Equal(t, "info", sink.Index(0).Level)
}

func TestGithubCIFromContext_Disabled(t *testing.T) {
	assert.False(t, GithubCIFromContext(context.Background()))
}

func TestGithubCIFromContext_Enabled(t *testing.T) {
	assert.True(t, GithubCIFromContext(ContextWithGithubCI(context.Background())))
}
