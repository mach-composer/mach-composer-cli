package graph

import (
	"context"
	"errors"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

type MockLoader struct {
	mock.Mock
}

func (m *MockLoader) Load(ctx context.Context, path string) (cty.Value, error) {
	args := m.Called(ctx, path)
	return args.Get(0).(cty.Value), args.Error(1)
}

func TestLoadOutputsSingleError(t *testing.T) {
	ctx := context.Background()

	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site1 := new(NodeMock)
	site1.On("Path").Return("main/site-1")

	mockLoader := new(MockLoader)
	mockLoader.On("Load", ctx, "main/site-1").Return(cty.NilVal, errors.New("error"))

	err := ig.AddVertex(site1)
	assert.NoError(t, err)

	g := &Graph{Graph: ig}

	err = LoadOutputs(ctx, g, mockLoader.Load)

	expected := &cli.GroupedError{}

	assert.ErrorAs(t, err, &expected)

	var gErr *cli.GroupedError
	errors.As(err, &gErr)
	assert.Equal(t, 1, len(gErr.Errors))
}

func TestLoadOutputsMultipleError(t *testing.T) {
	ctx := context.Background()

	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site1 := new(NodeMock)
	site1.On("Path").Return("main/site-1")
	site2 := new(NodeMock)
	site2.On("Path").Return("main/site-2")
	site2.On("SetOutputs", cty.StringVal("hello-world")).Return().Once()

	mockLoader := new(MockLoader)
	mockLoader.On("Load", ctx, "main/site-1").Return(cty.NilVal, errors.New("error")).Once()
	mockLoader.On("Load", ctx, "main/site-2").Return(cty.StringVal("hello-world"), nil).Once()

	_ = ig.AddVertex(site1)
	_ = ig.AddVertex(site2)

	g := &Graph{Graph: ig}

	err := LoadOutputs(ctx, g, mockLoader.Load)

	var gErr *cli.GroupedError
	errors.As(err, &gErr)
	assert.Equal(t, 1, len(gErr.Errors))
}
