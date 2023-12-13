package dependency

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/zclconf/go-cty/cty"
)

const (
	ProjectType       Type = "project"
	SiteType          Type = "site"
	SiteComponentType Type = "site-component"
)

type Type string

type Node interface {
	//Path returns the directory path of the node, relative to the global output directory
	Path() string

	//Identifier returns the identifier of the node as set in the configurations
	Identifier() string

	//Type returns the type of the node
	Type() Type

	//Ancestor returns the ancestor of the node. The ancestor is specific to the type of the node. For example,
	//a site will have the project as ancestor, a site component will have the site as ancestor,
	//and project will have no ancestor
	Ancestor() Node

	//Parents returns the direct parents of the node
	Parents() ([]Node, error)

	//Independent returns true if the node can be deployed independently, false otherwise
	Independent() bool

	//Tainted indicates if.
	Tainted() bool

	//Hash returns the hash of the node. The hash is based on the node's configuration as well as the configuration of any
	//related components. This can be compared to other hashes to determine whether a node has changed
	Hash() (string, error)

	//Outputs returns the outputs of the node
	Outputs() cty.Value

	//SetOutputs sets the outputs of the node
	SetOutputs(cty.Value)

	//SetTainted sets the tainted status of the node
	SetTainted(tainted bool)

	//HasChanges returns true if the node has changes, false otherwise
	HasChanges() (bool, error)

	//ResetGraph resets the graph of the node. If the graph the node belongs to the node graphs must also be reset,
	//as these are used to determine the parents of the node
	resetGraph(graph.Graph[string, Node])
}

type baseNode struct {
	graph          graph.Graph[string, Node]
	path           string
	identifier     string
	typ            Type
	ancestor       Node
	deploymentType config.DeploymentType
	tainted        bool
	outputs        cty.Value
}

func newBaseNode(graph graph.Graph[string, Node], path string, identifier string, typ Type, ancestor Node, deploymentType config.DeploymentType) baseNode {
	return baseNode{graph: graph,
		path:           path,
		identifier:     identifier,
		typ:            typ,
		ancestor:       ancestor,
		deploymentType: deploymentType,
		tainted:        false,
		outputs:        cty.NilVal,
	}
}

func (n *baseNode) resetGraph(ng graph.Graph[string, Node]) {
	n.graph = ng
}

func (n *baseNode) Outputs() cty.Value {
	return n.outputs
}

func (n *baseNode) SetOutputs(val cty.Value) {
	n.outputs = val
}

func (n *baseNode) SetTainted(tainted bool) {
	n.tainted = tainted
}

func (n *baseNode) Tainted() bool {
	return n.tainted
}

func (n *baseNode) Path() string {
	return n.path
}

func (n *baseNode) Identifier() string {
	return n.identifier
}

func (n *baseNode) Type() Type {
	return n.typ
}

func (n *baseNode) Ancestor() Node {
	return n.ancestor
}

func (n *baseNode) Parents() ([]Node, error) {
	pm, err := n.graph.PredecessorMap()
	if err != nil {
		return nil, err
	}

	eg := pm[n.Path()]

	var parents []Node
	for _, pathElement := range eg {
		p, err := n.graph.Vertex(pathElement.Source)
		if err != nil {
			return nil, err
		}
		parents = append(parents, p)
	}

	return parents, nil
}

func (n *baseNode) Independent() bool {
	// Projects and sites are always independent elements
	if n.typ == ProjectType || n.typ == SiteType {
		return true
	}

	//If deployment type is set to site-component and node is a site-component we consider it independent also
	if n.typ == SiteComponentType && n.deploymentType == config.DeploymentSiteComponent {
		return true
	}

	return false
}
