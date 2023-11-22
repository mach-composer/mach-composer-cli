package dependency

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/zclconf/go-cty/cty"
)

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
	Independent() bool
	HasChanges() (bool, error)
	Tainted() bool
	SetTainted(tainted bool)
	Hash() (string, error)
	LoadOutputs(ctx context.Context) error
	Outputs() cty.Value
}

type baseNode struct {
	path           string
	identifier     string
	typ            Type
	parent         Node
	deploymentType config.DeploymentType
	tainted        bool
	outputs        cty.Value
}

// LoadOutputs fetches all the outputs for the given state file. It will return a cty.NilVal if no outputs are present.
// The outputs are cached in the node.
func (n *baseNode) LoadOutputs(ctx context.Context) error {
	tfPath := fmt.Sprintf("deployments/%s", n.Path())
	val, err := utils.GetTerraformOutputs(ctx, tfPath)
	if err != nil {
		return err
	}
	n.outputs = val
	return nil
}

func (n *baseNode) Outputs() cty.Value {
	return n.outputs
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

func (n *baseNode) Parent() Node {
	return n.parent
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
