package runner

import (
	"context"
	"errors"
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	internalgraph "github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBatchNodesDepth1(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	start := new(internalgraph.NodeMock)
	start.On("Path").Return("main/site-1")

	_ = ig.AddVertex(start)

	g := &internalgraph.Graph{Graph: ig, StartNode: start}

	batches := batchNodes(g)

	assert.Equal(t, 1, len(batches))
}

func TestBatchNodesDepth2(t *testing.T) {
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

	batches := batchNodes(g)

	assert.Equal(t, 2, len(batches))
	assert.Equal(t, 1, len(batches[0]))
	assert.Equal(t, "main/site-1", batches[0][0].Path())
	assert.Equal(t, 2, len(batches[1]))
	assert.Contains(t, batches[1][0].Path(), "component")
	assert.Contains(t, batches[1][1].Path(), "component")
}

func TestBatchNodesDepth3(t *testing.T) {
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

	batches := batchNodes(g)

	assert.Equal(t, 3, len(batches))
	assert.Equal(t, 1, len(batches[0]))
	assert.Equal(t, "main/site-1", batches[0][0].Path())
	assert.Equal(t, 1, len(batches[1]))
	assert.Contains(t, batches[1][0].Path(), "main/site-1/component-1")
	assert.Equal(t, 1, len(batches[2]))
	assert.Contains(t, batches[2][0].Path(), "main/site-1/component-2")
}

func TestGraphRunnerMultipleLevels(t *testing.T) {
	project := new(internalgraph.NodeMock)
	project.On("Identifier").Return("main")

	site := new(internalgraph.NodeMock)
	site.On("Identifier").Return("site-1")
	site.SetTainted(false)

	component1 := new(internalgraph.NodeMock)
	component1.On("Identifier").Return("component-1")
	component1.SetTainted(false)

	component2 := new(internalgraph.NodeMock)
	component2.On("Identifier").Return("component-2")
	component2.SetTainted(true)

	component3 := new(internalgraph.NodeMock)
	component3.On("Identifier").Return("component-3")
	component3.SetTainted(true)

	runner := NewGraphRunner(1)
	runner.taint = func(ctx context.Context, g *internalgraph.Graph) error {
		return nil
	}
	runner.batch = func(g *internalgraph.Graph) map[int][]internalgraph.Node {
		return map[int][]internalgraph.Node{
			0: {project},
			1: {site},
			2: {component1, component2},
			3: {component3},
		}
	}

	var called []string

	err := runner.run(context.Background(), &internalgraph.Graph{}, func(ctx context.Context, node internalgraph.Node) (string, error) {
		called = append(called, node.Identifier())
		return "", nil
	}, "test", false)

	assert.NoError(t, err)
	assert.Equal(t, []string{"component-2", "component-3"}, called)
}

func TestGraphRunnerError(t *testing.T) {
	project := new(internalgraph.NodeMock)
	project.On("Identifier").Return("main")

	site := new(internalgraph.NodeMock)
	site.On("Identifier").Return("site-1")
	site.SetTainted(true)

	component1 := new(internalgraph.NodeMock)
	component1.On("Identifier").Return("component-1")
	component1.SetTainted(true)

	component2 := new(internalgraph.NodeMock)
	component2.On("Identifier").Return("component-2")
	component2.SetTainted(true)

	runner := NewGraphRunner(1)
	runner.taint = func(ctx context.Context, g *internalgraph.Graph) error {
		return nil
	}
	runner.batch = func(g *internalgraph.Graph) map[int][]internalgraph.Node {
		return map[int][]internalgraph.Node{
			0: {project},
			1: {site},
			2: {component1, component2},
		}
	}

	err := runner.run(context.Background(), &internalgraph.Graph{}, func(ctx context.Context, node internalgraph.Node) (string, error) {
		if node.Identifier() == "component-2" {
			return "", assert.AnError
		}
		return "", nil
	}, "test", false)

	cliErr := &cli.GroupedError{}

	assert.ErrorAs(t, err, &cliErr)

	errors.As(err, &cliErr)
	assert.Len(t, cliErr.Errors, 1)
	assert.Equal(t, assert.AnError, cliErr.Errors[0])
}

func TestGraphRunnerForce(t *testing.T) {
	project := new(internalgraph.NodeMock)
	project.On("Identifier").Return("main")

	site := new(internalgraph.NodeMock)
	site.On("Identifier").Return("site-1")
	site.SetTainted(false)

	runner := NewGraphRunner(1)
	runner.taint = func(ctx context.Context, g *internalgraph.Graph) error {
		return nil
	}
	runner.batch = func(g *internalgraph.Graph) map[int][]internalgraph.Node {
		return map[int][]internalgraph.Node{
			0: {project},
			1: {site},
		}
	}

	var called []string

	err := runner.run(context.Background(), &internalgraph.Graph{}, func(ctx context.Context, node internalgraph.Node) (string, error) {
		called = append(called, node.Identifier())
		return "", nil
	}, true)

	assert.NoError(t, err)
	assert.Equal(t, []string{"site-1"}, called)
}
