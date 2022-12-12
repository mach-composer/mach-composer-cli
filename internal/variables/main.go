package variables

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/utils"
)

type NotFoundError struct {
	Name string
	Node *yaml.Node
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("variable %s not found", e.Name)
}

// Support both ${var.foobar} as well as ${env.foobar}
var varRegex = regexp.MustCompile(`\${((?:var|env)(?:\.[^\}]+)+)}`)

type Value struct {
	Origin    string
	Value     string
	Encrypted bool
}

type Variables struct {
	vars           map[string]Value
	EncryptedFiles []string

	// For now we only support loading one external file. To support multilpe
	// we need to track variable sources to check if they are encrypted or not
	loadedFile bool
}

func NewVariables() *Variables {
	v := &Variables{
		vars:           make(map[string]Value),
		EncryptedFiles: []string{},
	}
	return v
}

func (v *Variables) Get(key string) (string, error) {
	if strings.HasPrefix(key, "var.") {
		trimmedKey := key[4:]

		result, ok := v.vars[trimmedKey]
		if !ok {
			return "", &NotFoundError{Name: key}
		}

		if result.Encrypted {
			result := fmt.Sprintf(`${data.sops_external.variables.data["%s"]}`, trimmedKey)
			return result, nil
		}

		return result.Value, nil
	}

	if strings.HasPrefix(key, "env.") {
		trimmedKey := key[4:]
		return os.Getenv(trimmedKey), nil
	}

	log.Warn().Msgf("Unsupported variables type %s", key)
	return "", nil
}

func (v *Variables) Set(key string, value string) {
	v.vars[key] = Value{Value: value}
}

func (v *Variables) HasEncrypted() bool {
	return len(v.EncryptedFiles) > 0
}

// Recursive function to replace the
func (v *Variables) InterpolateNode(node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		val, err := v.interpolateValue(node.Value)
		if err != nil {
			if notFoundErr, ok := err.(*NotFoundError); ok {
				notFoundErr.Node = node
			}
			return err
		}
		node.Value = val
		return nil
	}

	// Loop through the content if available to update childs
	for i := range node.Content {
		// Skip over keys
		if node.Kind == yaml.MappingNode && i%2 == 0 {
			continue
		}

		err := v.InterpolateNode(node.Content[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Variables) interpolateValue(val string) (string, error) {
	matches := varRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return val, nil
	}

	for _, match := range matches {
		replacement, err := v.Get(match[1])
		if err != nil {
			return "", err
		}
		val = strings.ReplaceAll(val, match[0], replacement)
	}

	return val, nil
}

// newVariablesFromFile creates a new Variables struct based on the contents
// of the given file.
func (v *Variables) Load(ctx context.Context, filename string) error {
	if v.loadedFile {
		panic("Only one external file is supported currently")
	}
	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return err
	}

	values := make(map[string]any)
	if err := yaml.Unmarshal(body, &values); err != nil {
		return err
	}

	isEncrypted := false
	if _, ok := values["sops"]; ok {
		isEncrypted = true
		v.EncryptedFiles = append(v.EncryptedFiles, filename)
		delete(values, "sops")
	}

	dst := map[string]Value{}
	serializeNestedVariables(values, dst, "")

	for key, val := range dst {
		val.Encrypted = isEncrypted
		val.Origin = filename
		v.vars[key] = val
	}

	v.loadedFile = true
	return nil
}

// serializeNestedVariables reads a map recursively building a list of variable
// strings. It converst for example the following:
//
//	map[string]any{
//		"foo": "bar",
//		"my": map[string]any{
//			"var": 10,
//		},
//	}
//
// into:
//
//	map[string]string{
//		"foo": "bar",
//		"my.var": "10",
//	}
func serializeNestedVariables(in map[string]any, out map[string]Value, prefix string) {
	for k, v := range in {
		var key string
		if prefix != "" {
			key = fmt.Sprintf("%s.%s", prefix, k)
		} else {
			key = k
		}

		switch v := v.(type) {
		case string:
			out[key] = Value{Value: v}
		case int:
			out[key] = Value{Value: fmt.Sprint(v)}
		case map[string]any:
			serializeNestedVariables(v, out, key)
		}
	}
}
