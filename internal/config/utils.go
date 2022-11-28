package config

import (
	"github.com/elliotchance/pie/v2"
	"gopkg.in/yaml.v3"
)

func mapYamlNodes(nodes []*yaml.Node) map[string]*yaml.Node {
	result := map[string]*yaml.Node{}
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
