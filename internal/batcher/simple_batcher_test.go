package batcher

import (
	"github.com/dominikbraun/graph"
	internalgraph "github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSimpleBatchNodesDepth1(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	start := new(internalgraph.NodeMock)
	start.On("Path").Return("main/site-1")

	_ = ig.AddVertex(start)

	g := &internalgraph.Graph{Graph: ig, StartNode: start}

	batches, err := simpleBatchFunc()(g)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(batches))
}

func TestSimpleBatchNodesDepth2(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(internalgraph.NodeMock)
	site.On("Path").Return("main/site-1")

	component1 := new(internalgraph.NodeMock)
	component1.On("Path").Return("main/site-1/component-1")

	component2 := new(internalgraph.NodeMock)
	component2.On("Path").Return("main/site-1/component-2")

	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component1)
	_ = ig.AddVertex(component2)

	_ = ig.AddEdge("main/site-1", "main/site-1/component-1")
	_ = ig.AddEdge("main/site-1", "main/site-1/component-2")

	g := &internalgraph.Graph{Graph: ig, StartNode: site}

	batches, err := simpleBatchFunc()(g)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(batches))
	assert.Equal(t, 1, len(batches[0]))
	assert.Equal(t, "main/site-1", batches[0][0].Path())
	assert.Equal(t, 2, len(batches[1]))
	assert.Contains(t, batches[1][0].Path(), "component")
	assert.Contains(t, batches[1][1].Path(), "component")
}

func TestSimpleBatchNodesDepth3(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(internalgraph.NodeMock)
	site.On("Path").Return("main/site-1")

	component1 := new(internalgraph.NodeMock)
	component1.On("Path").Return("main/site-1/component-1")

	component2 := new(internalgraph.NodeMock)
	component2.On("Path").Return("main/site-1/component-2")

	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component1)
	_ = ig.AddVertex(component2)

	_ = ig.AddEdge("main/site-1", "main/site-1/component-1")
	_ = ig.AddEdge("main/site-1/component-1", "main/site-1/component-2")

	g := &internalgraph.Graph{Graph: ig, StartNode: site}

	batches, err := simpleBatchFunc()(g)

	assert.NoError(t, err)
	assert.Equal(t, 3, len(batches))
	assert.Equal(t, 1, len(batches[0]))
	assert.Equal(t, "main/site-1", batches[0][0].Path())
	assert.Equal(t, 1, len(batches[1]))
	assert.Contains(t, batches[1][0].Path(), "main/site-1/component-1")
	assert.Equal(t, 1, len(batches[2]))
	assert.Contains(t, batches[2][0].Path(), "main/site-1/component-2")
}
