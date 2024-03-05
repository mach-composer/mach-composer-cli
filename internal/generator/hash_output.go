package generator

import (
	"embed"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

//go:embed templates/hash_output.tmpl
var hashOutputTmpl embed.FS

// renderHashOutput uses templates/hash_output.tmpl to generate a terraform snippet for each node
func renderHashOutput(n graph.Node, siteComponents []config.SiteComponentConfig) (string, error) {
	tpl, err := hashOutputTmpl.ReadFile("templates/hash_output.tmpl")
	if err != nil {
		return "", err
	}

	hash, err := n.Hash()
	if err != nil {
		return "", err
	}

	var componentNames []string
	for _, component := range siteComponents {
		componentNames = append(componentNames, component.Name)
	}

	return utils.RenderGoTemplate(string(tpl), struct {
		NodeHash       string
		ComponentNames []string
	}{
		NodeHash:       hash,
		ComponentNames: componentNames,
	})
}
