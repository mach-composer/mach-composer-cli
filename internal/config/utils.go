package config

import (
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"
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

// mapYamlNodes cretes a map[key]node from a slice of yaml.Node's. It assumes
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
	return result
}

func iterateYamlNodes(
	nodes map[string]*yaml.Node,
	skipKeys []string,
	cb func(key string, data map[string]any) error) error {
	for key, node := range nodes {
		if pie.Contains(skipKeys, key) {
			continue
		}

		data, err := nodeAsMap(node)
		if err != nil {
			return err
		}

		if err := cb(key, data); err != nil {
			return err
		}
	}

	return nil
}

// LoadRefData will load referenced files and replace the node with the content of these files. It works both with the
// ${include()} syntax and the $ref syntax.
func LoadRefData(ctx context.Context, node *yaml.Node, cwd string) (*yaml.Node, string, error) {
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
		data := mapYamlNodes(node.Content)
		ref, ok := data["$ref"]
		if !ok {
			return nil, "", nil
		}

		newNode, err := loadRefDocument(ctx, ref.Value, cwd)
		if err != nil {
			return nil, "", err
		}

		return newNode, ref.Value, nil
	default:
		return node, "", nil
	}
}

func loadRefDocument(_ context.Context, filename, cwd string) (*yaml.Node, error) {
	parts := strings.SplitN(filename, "#", 2)
	fname := parts[0]

	body, err := utils.AFS.ReadFile(fname)
	if err != nil {
		return nil, err
	}
	result := &yaml.Node{}
	if err = yaml.Unmarshal(body, result); err != nil {
		return nil, err
	}

	root := result.Content[0]

	if len(parts) > 1 {
		path := strings.TrimPrefix(parts[1], "/")
		node := root
		for _, p := range strings.Split(path, "/") {
			nodes := mapYamlNodes(node.Content)
			n, ok := nodes[p]
			if !ok {
				return nil, fmt.Errorf("unable to resolve node %s", parts[1])
			}
			node = n
		}
		root = node
	}

	return root, nil
}

func loadIncludeDocument(node *yaml.Node, path string) (*yaml.Node, string, error) {
	re := regexp.MustCompile(`\$\{include\(([^)]+)\)\}`)
	data := re.FindStringSubmatch(node.Value)
	if len(data) != 2 {
		return nil, "", fmt.Errorf("failed to parse ${include()} tag")
	}
	filename := filepath.Join(path, data[1])
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
