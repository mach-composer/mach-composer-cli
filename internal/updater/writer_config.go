package updater

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func machConfigUpdater(ctx context.Context, cfg *PartialConfig, src []byte, updateSet *UpdateSet) []byte {
	nodes := mapComponentNodes(cfg.ComponentsNode)

	// Walk through the updated components and search the corresponding yaml node
	// via the previously created mapping. Withing the node search for the
	// `version` tag and use the line number to change the value in the source
	// document (lines list)
	lines := SplitLines(string(src))
	for _, c := range updateSet.updates {
		node, ok := nodes[c.Component.Name]
		if !ok {
			log.Warn().Msgf("Component %s with update not found in yaml file", c.Component.Name)
			continue
		}

		for i, n := range node.Content {
			if n.Tag == "!!str" && n.Value == "version" {
				// The value is in the node after this node. Assume it's always
				// sequential
				vn := node.Content[i+1]
				if vn.Value != c.Component.Version {
					log.Fatal().Msg("Unexpected version")
				}

				// Make sure the version is always quoted. This is currently not
				// super robust.
				replacement := c.LastVersion
				firstChar := lines[vn.Line-1][vn.Column-1]
				if firstChar != '"' && firstChar != '\'' {
					replacement = fmt.Sprintf(`"%s"`, replacement)
				}

				key := lines[vn.Line-1][:vn.Column-1]
				value := lines[vn.Line-1][vn.Column-1:]

				lines[vn.Line-1] = key + strings.Replace(value, vn.Value, replacement, 1)
				break
			}
		}
	}

	output := strings.Join(lines, "\n") + "\n"
	return []byte(output)
}

// mapComponentNodes creates a mapping where the key is the component name and
// the value the yaml node of the component. We do this by iterating all
// component children and searching for the `name` tag.
func mapComponentNodes(node *yaml.Node) map[string]*yaml.Node {
	nodes := map[string]*yaml.Node{}
	for _, cn := range node.Content {
		for i, n := range cn.Content {
			if n.Tag == "!!str" && n.Value == "name" {
				name := cn.Content[i+1].Value
				nodes[name] = cn
				break
			}
		}
	}
	return nodes
}

// machFileWriter updates the contents of a mach file with the updated
// version of the components
func machFileWriter(ctx context.Context, cfg *PartialConfig, updates *UpdateSet) {
	input, err := utils.AFS.ReadFile(updates.filename)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to read file")
	}

	output := machConfigUpdater(ctx, cfg, input, updates)

	err = utils.AFS.WriteFile(updates.filename, output, 0644)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to write file")
	}
}
