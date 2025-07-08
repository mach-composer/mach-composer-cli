package graph

import (
	"github.com/dominikbraun/graph"
	"github.com/mach-composer/mach-composer-cli/internal/config"
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

	//Children returns the direct children of the node. This is used to determine the nodes that are dependent on this node.
	Children() ([]Node, error)

	//Independent returns true if the node can be deployed independently, false otherwise
	Independent() bool

	//Tainted indicates if.
	Tainted() bool

	//Hash returns the hash of the node. The hash is based on the node's configuration as well as the configuration of any
	//related components. This can be compared to other hashes to determine whether a node has changed
	Hash() (string, error)

	//SetTainted sets the tainted status of the node
	SetTainted(tainted bool)

	//ResetGraph resets the graph of the node. If the graph the node belongs to the node graphs must also be reset,
	//as these are used to determine the parents of the node
	resetGraph(graph.Graph[string, Node])

	//SetOldHash sets the old hash of the node. This is used to determine if the node has changed
	SetOldHash(hash string)

	//GetOldHash returns the old hash of the node
	GetOldHash() string
}

type baseNode struct {
	graph          graph.Graph[string, Node]
	path           string
	identifier     string
	typ            Type
	ancestor       Node
	deploymentType config.DeploymentType
	tainted        bool
	oldHash        string
}

func newBaseNode(graph graph.Graph[string, Node], path string, identifier string, typ Type, ancestor Node, deploymentType config.DeploymentType) baseNode {
	return baseNode{graph: graph,
		path:           path,
		identifier:     identifier,
		typ:            typ,
		ancestor:       ancestor,
		deploymentType: deploymentType,
		tainted:        false,
	}
}

func (n *baseNode) resetGraph(ng graph.Graph[string, Node]) {
	n.graph = ng
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

func (n *baseNode) Children() ([]Node, error) {
	pm, err := n.graph.AdjacencyMap()
	if err != nil {
		return nil, err
	}

	eg := pm[n.Path()]

	var children []Node
	for _, pathElement := range eg {
		p, err := n.graph.Vertex(pathElement.Target)
		if err != nil {
			return nil, err
		}
		children = append(children, p)
	}

	return children, nil
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

func (n *baseNode) SetOldHash(hash string) {
	n.oldHash = hash
}

func (n *baseNode) GetOldHash() string {
	return n.oldHash
}
