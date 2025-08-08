package batcher

import (
	"github.com/dominikbraun/graph"
	internalgraph "github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSiteBatchFuncReturnsErrorWhenNoProjects(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(internalgraph.NodeMock)
	site.On("Path").Return("main/site-1")
	site.On("Type").Return(internalgraph.SiteType)
	site.On("Identifier").Return("site-1")

	_ = ig.AddVertex(site)

	g := &internalgraph.Graph{Graph: ig, StartNode: site}

	batches, err := siteBatchFunc([]string{"site-1"})(g)
	assert.Nil(t, batches)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected 1 project")
}

func TestSiteBatchFuncReturnsErrorWhenMultipleProjects(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project1 := new(internalgraph.NodeMock)
	project1.On("Path").Return("main/project-1")
	project1.On("Type").Return(internalgraph.ProjectType)

	project2 := new(internalgraph.NodeMock)
	project2.On("Path").Return("main/project-2")
	project2.On("Type").Return(internalgraph.ProjectType)

	_ = ig.AddVertex(project1)
	_ = ig.AddVertex(project2)

	g := &internalgraph.Graph{Graph: ig, StartNode: project1}

	batches, err := siteBatchFunc([]string{})(g)
	assert.Nil(t, batches)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "expected 1 project")
}

func TestSiteBatchFuncReturnsErrorWhenSiteNotFound(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project := new(internalgraph.NodeMock)
	project.On("Path").Return("main/project-1")
	project.On("Type").Return(internalgraph.ProjectType)

	_ = ig.AddVertex(project)

	g := &internalgraph.Graph{Graph: ig, StartNode: project}
	batches, err := siteBatchFunc([]string{"site-unknown"})(g)
	assert.Nil(t, batches)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "site with identifier site-unknown not found")
}

func TestSiteBatchFuncBatchesNodesBySiteOrderAndDepth(t *testing.T) {
	ig := graph.New(func(n internalgraph.Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	project := new(internalgraph.NodeMock)
	project.On("Path").Return("main/project-1")
	project.On("Type").Return(internalgraph.ProjectType)

	site := new(internalgraph.NodeMock)
	site.On("Path").Return("main/site-1")
	site.On("Type").Return(internalgraph.SiteType)
	site.On("Identifier").Return("site-1")

	component := new(internalgraph.NodeMock)
	component.On("Path").Return("main/site-1/component-1")
	component.On("Type").Return(internalgraph.SiteComponentType)
	component.On("Children").Return([]internalgraph.Node{}, nil)

	site.On("Children").Return([]internalgraph.Node{component}, nil)

	_ = ig.AddVertex(project)
	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component)

	_ = ig.AddEdge("main/project-1", "main/site-1")
	_ = ig.AddEdge("main/site-1", "main/site-1/component-1")

	g := &internalgraph.Graph{Graph: ig, StartNode: project}

	batches, err := siteBatchFunc([]string{"site-1"})(g)
	assert.NoError(t, err)
	assert.NotNil(t, batches)
	assert.Equal(t, 3, len(batches))
	assert.Equal(t, "main/project-1", batches[0][0].Path())
	assert.Equal(t, "main/site-1", batches[1][0].Path())
}
