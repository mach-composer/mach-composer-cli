package dependency

const (
	ProjectType       Type = "project"
	SiteType          Type = "site"
	SiteComponentType Type = "site-component"
)

type Type string

type Node interface {
	Path() string
	Identifier() string
	Type() Type
	Parent() Node
}

type node struct {
	path       string
	identifier string
	typ        Type
	parent     Node
}

func (n *node) Path() string {
	return n.path
}

func (n *node) Identifier() string {
	return n.identifier
}

func (n *node) Type() Type {
	return n.typ
}

func (n *node) Parent() Node {
	return n.parent
}
