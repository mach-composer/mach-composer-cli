package config

import (
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
