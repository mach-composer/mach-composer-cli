package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/stretchr/testify/mock"
	"github.com/zclconf/go-cty/cty"
)

type NodeMock struct {
	mock.Mock
	tainted bool
}

func (n *NodeMock) Path() string {
	args := n.Called()
	return args.String(0)
}

func (n *NodeMock) Identifier() string {
	args := n.Called()
	return args.String(0)
}

func (n *NodeMock) Type() Type {
	args := n.Called()
	return args.Get(0).(Type)
}

func (n *NodeMock) Ancestor() Node {
	//TODO implement me
	panic("implement me")
}

func (n *NodeMock) Parents() ([]Node, error) {
	args := n.Called()
	return args.Get(0).([]Node), args.Error(1)
}

func (n *NodeMock) Independent() bool {
	//TODO implement me
	panic("implement me")
}

func (n *NodeMock) Tainted() bool {
	return n.tainted
}

func (n *NodeMock) Hash() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (n *NodeMock) Outputs() cty.Value {
	args := n.Called()
	return args.Get(0).(cty.Value)
}

func (n *NodeMock) SetOutputs(value cty.Value) {
	_ = n.Called(value)
}

func (n *NodeMock) SetTainted(tainted bool) {
	n.tainted = tainted
}

func (n *NodeMock) HasChanges() (bool, error) {
	args := n.Called()
	return args.Bool(0), args.Error(1)
}

func (n *NodeMock) resetGraph(g graph.Graph[string, Node]) {
	_ = n.Called(g)

	return
}
