package updater

import (
	"fmt"
	"log"
	"strings"

	"github.com/labd/mach-composer/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

type MachConfig struct {
	Components yaml.Node
}

func MachConfigUpdater(src []byte, updateSet *UpdateSet) []byte {
	data := MachConfig{}
	err := yaml.Unmarshal(src, &data)
	if err != nil {
		log.Fatalln(err)
	}

	// Create a mapping where the key is the component name and the value the
	// yaml node of the component. We do this by iterating all component children
	// and searching for the `name` tag.
	nodes := map[string]*yaml.Node{}
	for _, cn := range data.Components.Content {
		for i, n := range cn.Content {
			if n.Tag == "!!str" && n.Value == "name" {
				name := cn.Content[i+1].Value
				nodes[name] = cn
				break
			}
		}
	}

	// Walk through the updated components and search the corresponding yaml node
	// via the previously created mapping. Withing the node search for the
	// `version` tag and use the line number to change the value in the source
	// document (lines list)
	lines := SplitLines(string(src))
	for _, c := range updateSet.updates {
		node, ok := nodes[c.Component.Name]
		if !ok {
			logrus.Warn("Component with update not found in yaml file")
			continue
		}

		for i, n := range node.Content {
			if n.Tag == "!!str" && n.Value == "version" {

				// The value is in the node after this node. Assume it's always
				// sequential
				vn := node.Content[i+1]
				if vn.Value != c.Component.Version {
					log.Fatal("Unexpected version")
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

// MachFileWriter updates the contents of a mach file with the updated
// version of the components
func MachFileWriter(updates *UpdateSet) {

	input, err := utils.AFS.ReadFile(updates.filename)
	if err != nil {
		log.Fatalln(err)
	}

	output := MachConfigUpdater(input, updates)

	err = utils.AFS.WriteFile(updates.filename, output, 0644)
	if err != nil {
		log.Fatalln(err)
	}
}
