package generator

import (
	"embed"
	"github.com/mach-composer/mach-composer-cli/internal/dependency"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

//go:embed templates/hash_output.tmpl
var hashOutputTmpl embed.FS

// renderHashOutput uses templates/hash_output.tmpl to generate a terraform snippet for each node
func renderHashOutput(n dependency.Node) (string, error) {
	tpl, err := hashOutputTmpl.ReadFile("templates/hash_output.tmpl")
	if err != nil {
		return "", err
	}

	hash, err := n.Hash()
	if err != nil {
		return "", err
	}

	return utils.RenderGoTemplate(string(tpl), struct {
		NodeHash string
	}{
		NodeHash: hash,
	})
}
