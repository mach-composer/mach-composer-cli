package config

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/utils"
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

func (v *Variables) Get(key string) (string, error) {
	if strings.HasPrefix(key, "var.") {
		trimmedKey := key[4:]

		result, ok := v.vars[trimmedKey]
		if !ok {
			return "", fmt.Errorf("missing variable %s", key)
		}

		if v.Encrypted {
			result := fmt.Sprintf(
				`${data.sops_external.variables.data["%s"]}`, trimmedKey)
			return result, nil
		}

		return result, nil
	}

	if strings.HasPrefix(key, "env.") {
		trimmedKey := key[4:]
		return os.Getenv(trimmedKey), nil
	}

	logrus.Warningf("Unsupported variables type %s", key)
	return "", nil
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
		return nil, err
	}

	isEncrypted, err := yamlIsEncrypted(body)
	if err != nil {
		return nil, err
	}
	if isEncrypted {
		logrus.Debug("Detected SOPS encryption; decrypting...")
		body, err = DecryptYaml(filename)
		if err != nil {
			return nil, err
		}

		vars.Encrypted = true
	}

	dst := make(map[string]interface{})
	if err := yaml.Unmarshal(body, &dst); err != nil {
		return nil, err
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
func yamlIsEncrypted(data []byte) (bool, error) {
	dst := make(map[string]interface{})
	err := yaml.Unmarshal(data, &dst)
	if err != nil {
		return false, err
	}
	if _, ok := dst["sops"]; ok {
		return true, nil
	}
	return false, nil
}

// InterpolateVars interpolates the variables within the values of the given
// MACH Config in-place.
func InterpolateVars(raw *_RawMachConfig, vars *Variables) error {
	if vars == nil {
		return errors.New("vars cannot be nil")
	}

	if err := interpolateNode(&raw.Sites, vars); err != nil {
		return err
	}
	if err := interpolateNode(&raw.Components, vars); err != nil {
		return err
	}

	return nil
}

// Recursive function to replace the
func interpolateNode(node *yaml.Node, vars *Variables) error {
	if node.Kind == yaml.ScalarNode {
		if val, replaced, err := interpolateValue(node.Value, vars); err != nil {
			return err
		} else if replaced {
			node.Value = val
		}
		return nil
	}

	// Loop through the content if available to update childs
	for i := range node.Content {
		// Skip over keys
		if node.Kind == yaml.MappingNode && i%2 == 0 {
			continue
		}

		err := interpolateNode(node.Content[i], vars)
		if err != nil {
			return err
		}
	}
	return nil
}

func interpolateValue(val string, vars *Variables) (string, bool, error) {
	matches := varRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return val, false, nil
	}

	for _, match := range matches {
		replacement, err := vars.Get(match[1])
		if err != nil {
			return "", false, err
		}
		val = strings.ReplaceAll(val, match[0], replacement)
	}

	return val, true, nil
}
