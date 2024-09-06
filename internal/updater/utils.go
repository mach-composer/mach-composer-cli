package updater

import (
	"bufio"
	"context"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"gopkg.in/yaml.v3"
	"path"
	"strings"
)

func SplitLines(s string) []string {
	var lines []string
	sc := bufio.NewScanner(strings.NewReader(s))
	for sc.Scan() {
		lines = append(lines, sc.Text())
	}
	return lines
}

// LoadRefData will load referenced files and replace the node with the content of these files. It works both with the
// ${include()} syntax and the $ref syntax.
// Deprecated: Only used in the updater logic as we need the file name to update there
func loadRefData(_ context.Context, node *yaml.Node, cwd string) (*yaml.Node, string, error) {
	switch node.Kind {
	case yaml.ScalarNode:
		newNode, filePath, err := config.LoadIncludeDocument(node, cwd)
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
	data := config.MapYamlNodes(node.Content)
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
			nodes := config.MapYamlNodes(node.Content)
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
