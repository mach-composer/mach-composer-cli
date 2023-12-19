package graph

import (
	"fmt"
	"github.com/dominikbraun/graph"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTaintNodesNoChanges(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("main")
	site.On("HasChanges").Return(false, nil).Once()

	_ = ig.AddVertex(site)

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.NoError(t, err)
	assert.False(t, site.Tainted())
}

func TestTaintNodesHasChanges(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("site")
	site.On("HasChanges").Return(true, nil).Once()

	_ = ig.AddVertex(site)

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.NoError(t, err)
	assert.True(t, site.Tainted())
}

func TestTaintNodesHasChanges_Error(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("site")
	site.On("HasChanges").Return(false, fmt.Errorf("error")).Once()

	_ = ig.AddVertex(site)

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.Error(t, err)
}

func TestTaintNodesChildHasChanges(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("site").Times(2)
	site.On("HasChanges").Return(false, nil).Once()

	component := new(NodeMock)
	component.On("Path").Return("component").Times(2)
	component.On("HasChanges").Return(true, nil).Once()

	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component)

	_ = ig.AddEdge("site", "component")

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.NoError(t, err)
	assert.False(t, site.Tainted())
	assert.True(t, component.Tainted())
}

func TestTaintNodesChildHasChanges_Error(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("site").Times(2)
	site.On("HasChanges").Return(false, nil).Once()

	component := new(NodeMock)
	component.On("Path").Return("component").Times(2)
	component.On("HasChanges").Return(true, fmt.Errorf("error")).Once()

	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component)

	_ = ig.AddEdge("site", "component")

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.Error(t, err)
}

func TestTaintNodesParentHasChanges(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("site").Times(2)
	site.On("HasChanges").Return(true, nil).Once()

	component := new(NodeMock)
	component.On("Path").Return("component").Times(2)
	component.On("HasChanges").Return(false, nil).Once()

	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component)

	_ = ig.AddEdge("site", "component")

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.NoError(t, err)
	assert.True(t, site.Tainted())
	assert.True(t, component.Tainted())

	component.AssertNumberOfCalls(t, "HasChanges", 0)
}

func TestTaintNodesOneParentHasChanges(t *testing.T) {
	ig := graph.New(func(n Node) string { return n.Path() }, graph.Directed(), graph.Tree(), graph.PreventCycles())

	site := new(NodeMock)
	site.On("Path").Return("site").Times(2)
	site.On("HasChanges").Return(false, nil).Once()

	component1 := new(NodeMock)
	component1.On("Path").Return("component-1").Times(2)
	component1.On("HasChanges").Return(true, nil).Once()

	component2 := new(NodeMock)
	component2.On("Path").Return("component-2").Times(2)
	component2.On("HasChanges").Return(false, nil).Once()

	component3 := new(NodeMock)
	component3.On("Path").Return("component-3").Times(2)
	component3.On("HasChanges").Return(false, nil).Once()

	_ = ig.AddVertex(site)
	_ = ig.AddVertex(component1)
	_ = ig.AddVertex(component2)
	_ = ig.AddVertex(component3)

	_ = ig.AddEdge("site", "component-1")
	_ = ig.AddEdge("site", "component-2")
	_ = ig.AddEdge("component-1", "component-3")
	_ = ig.AddEdge("component-2", "component-3")

	g := &Graph{Graph: ig, StartNode: site}

	err := TaintNodes(g)
	assert.NoError(t, err)

	assert.False(t, site.Tainted())
	assert.True(t, component1.Tainted())
	assert.False(t, component2.Tainted())
	assert.True(t, component3.Tainted())
}
