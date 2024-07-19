package config

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"path"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func nodeAsMap(n *yaml.Node) (map[string]any, error) {
	target := map[string]any{}
	if err := n.Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}

// mapYamlNodes creates a map[key]node from a slice of yaml.Node's. It assumes
// that the nodes are pairs, e.g. [key, value, key, value]
func mapYamlNodes(nodes []*yaml.Node) map[string]*yaml.Node {
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
		aliasData := mapYamlNodes(val.Alias.Content)
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

// LoadRefData will load referenced files and replace the node with the content of these files. It works both with the
// ${include()} syntax and the $ref syntax.
func LoadRefData(_ context.Context, node *yaml.Node, cwd string) (*yaml.Node, string, error) {
	switch node.Kind {
	case yaml.ScalarNode:
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

		newNode, filePath, err := loadIncludeDocument(node, cwd)
		if err != nil {
			return nil, "", err
		}

		return newNode, filePath, nil
	case yaml.MappingNode:
		newNode, filePath, err := loadRefDocument(node, cwd)
		if err != nil {
			return nil, "", err
		}

		return newNode, filePath, nil
	default:
		return node, "", nil
	}
}

func loadRefDocument(node *yaml.Node, cwd string) (*yaml.Node, string, error) {
	data := mapYamlNodes(node.Content)
	ref, ok := data["$ref"]
	if !ok {
		return node, "", nil
	}

	parts := strings.SplitN(ref.Value, "#", 2)
	fileName := parts[0]

	body, err := utils.AFS.ReadFile(path.Join(cwd, fileName))
	if err != nil {
		return nil, "", err
	}
	result := &yaml.Node{}
	if err = yaml.Unmarshal(body, result); err != nil {
		return nil, "", err
	}

	root := result.Content[0]

	if len(parts) > 1 {
		p := strings.TrimPrefix(parts[1], "/")
		node := root
		for _, p := range strings.Split(p, "/") {
			nodes := mapYamlNodes(node.Content)
			n, ok := nodes[p]
			if !ok {
				return nil, "", fmt.Errorf("unable to resolve node %s", parts[1])
			}
			node = n
		}
		root = node
	}

	return root, fileName, nil
}

func loadIncludeDocument(node *yaml.Node, cwd string) (*yaml.Node, string, error) {
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
