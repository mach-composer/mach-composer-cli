package config

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"regexp"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func nodeAsMap(n *yaml.Node) (map[string]any, error) {
	target := map[string]any{}
	if err := n.Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}

// MapYamlNodes creates a map[key]node from a slice of yaml.Node's. It assumes
// that the nodes are pairs, e.g. [key, value, key, value]
func MapYamlNodes(nodes []*yaml.Node) map[string]*yaml.Node {
	result := map[string]*yaml.Node{}

	// Check if there are an even number of nodes as we expect a
	// key, value nodes
	if len(nodes)%2 != 0 {
		return nil
	}
	for i := 0; i < len(nodes); i += 2 {
		key := nodes[i].Value
		value := nodes[i+1]
		result[key] = value
	}

	// Check if there is an alias node and resolve it
	val, ok := result["<<"]
	if ok {
		aliasData := MapYamlNodes(val.Alias.Content)
		delete(result, "<<")

		// We only add the aliased data if the key is not already present
		for k, v := range aliasData {
			if _, ok := result[k]; ok {
				continue
			}
			result[k] = v
		}
	}

	return result
}

func LoadIncludeDocument(node *yaml.Node, cwd string) (*yaml.Node, string, error) {
	cli.DeprecationWarning(&cli.DeprecationOptions{
		Message: "the '${include()}' syntax is deprecated and will be removed in version 3.0",
		Details: `
				For example instead of:
					components: ${include(components.yml)}

				You should use:
					components:
						$ref: "components.yml"
			`,
	})

	re := regexp.MustCompile(`\$\{include\(([^)]+)\)\}`)
	data := re.FindStringSubmatch(node.Value)
	if len(data) != 2 {
		return nil, "", fmt.Errorf("failed to parse ${include()} tag")
	}
	filename := filepath.Join(cwd, data[1])
	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}

	result := yaml.Node{}
	if err = yaml.Unmarshal(body, &result); err != nil {
		return nil, "", err
	}
	if len(result.Content) != 1 {
		return nil, "", fmt.Errorf("invalid yaml file")
	}
	return result.Content[0], filename, nil
}
