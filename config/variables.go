package config

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/labd/mach-composer-go/utils"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

var varRegex = regexp.MustCompile(`\${((?:var|env)(?:\.[^\}]+)+)}`)

type Variables struct {
	vars      map[string]string
	Filepath  string
	Encrypted bool
}

func NewVariables() *Variables {
	vars := &Variables{}
	vars.vars = make(map[string]string)
	return vars
}

func (v *Variables) Get(key string) string {

	if strings.HasPrefix(key, "var.") {
		trimmedKey := key[4:]
		result, ok := v.vars[trimmedKey]
		if !ok {
			logrus.Warningf("no variable %s", trimmedKey)
		}

		if v.Encrypted {
			return fmt.Sprintf(
				`${data.sops_external.variables.data["%s"]}`, trimmedKey)
		}

		return result
	}

	if strings.HasPrefix(key, "env.") {
		trimmedKey := key[4:]
		return os.Getenv(trimmedKey)
	}

	logrus.Warningf("Unsupported variables type %s", key)
	return ""
}

func (v *Variables) Set(key string, value string) {
	if v.vars == nil {
		v.vars = make(map[string]string)
	}
	v.vars[key] = value
}

func loadVariables(filename string) (*Variables, error) {
	body, err := utils.AFS.ReadFile(filename)
	vars := NewVariables()
	vars.Filepath = filename

	if err != nil {
		panic(err)
	}

	if yamlIsEncrypted(body) {
		logrus.Debug("Detected SOPS encryption; decrypting...")
		body, err = DecryptYaml(filename)
		if err != nil {
			panic(err)
		}

		vars.Encrypted = true
	}

	dst := make(map[string]interface{})
	err = yaml.Unmarshal(body, &dst)
	if err != nil {
		panic(err)
	}

	if vars.Encrypted {
		delete(dst, "sops")
	}

	processVariablesYaml(dst, vars.vars, "")

	return vars, nil
}

func processVariablesYaml(in map[string]interface{}, out map[string]string, prefix string) {
	for k, v := range in {
		var key string
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, k)
		} else {
			key = k
		}

		switch v := v.(type) {
		case string:
			out[key] = v
		case int:
			out[key] = fmt.Sprint(v)
		case map[string]interface{}:
			processVariablesYaml(v, out, key)
		}
	}
}

// Check if the file is encrypted with sops
func yamlIsEncrypted(data []byte) bool {
	dst := make(map[string]interface{})
	err := yaml.Unmarshal(data, &dst)
	if err != nil {
		panic(err)
	}
	if _, ok := dst["sops"]; ok {
		return true
	}
	return false
}

// InterpolateVars interpolates the variables within the values of the given
// MACH Config in-place.
func InterpolateVars(raw *_RawMachConfig, vars *Variables) {
	if vars == nil {
		panic("Vars cannot be nil")
	}
	interpolateNode(&raw.Sites, vars)
	interpolateNode(&raw.Components, vars)
}

// Recursive function to replace the
func interpolateNode(node *yaml.Node, vars *Variables) {
	if node.Kind == yaml.ScalarNode {
		if val, replaced := interpolateValue(node.Value, vars); replaced {
			node.Value = val
		}
		return
	}

	// Loop through the content if available to update childs
	for i := range node.Content {

		// Skip over keys
		if node.Kind == yaml.MappingNode && i%2 == 0 {
			continue
		}

		interpolateNode(node.Content[i], vars)
	}
}

func interpolateValue(val string, vars *Variables) (string, bool) {
	matches := varRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return val, false
	}

	for _, match := range matches {
		replacement := vars.Get(match[1])
		val = strings.ReplaceAll(val, match[0], replacement)
	}

	return val, true
}
