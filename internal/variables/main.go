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

type Variables struct {
	vars           map[string]string
	EncryptedFiles []string

	// For now we only support loading one external file. To support multilpe
	// we need to track variable sources to check if they are encrypted or not
	loadedFile bool
}

func NewVariables() *Variables {
	v := &Variables{
		vars:           make(map[string]string),
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

		// FIXME: we should check the source of the var for this
		if len(v.EncryptedFiles) > 0 {
			result := fmt.Sprintf(`${data.sops_external.variables.data["%s"]}`, trimmedKey)
			return result, nil
		}

		return result, nil
	}

	if strings.HasPrefix(key, "env.") {
		trimmedKey := key[4:]
		return os.Getenv(trimmedKey), nil
	}

	log.Warn().Msgf("Unsupported variables type %s", key)
	return "", nil
}

func (v *Variables) Set(key string, value string) {
	if v.vars == nil {
		v.vars = make(map[string]string)
	}
	v.vars[key] = value
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

	isEncrypted, err := yamlIsEncrypted(body)
	if err != nil {
		return err
	}
	if isEncrypted {
		log.Info().Msgf("Detected SOPS encryption; decrypting...")
		body, err = utils.DecryptYaml(ctx, filename)
		if err != nil {
			return err
		}

		v.EncryptedFiles = append(v.EncryptedFiles, filename)
	}

	dst := make(map[string]any)
	if err := yaml.Unmarshal(body, &dst); err != nil {
		return err
	}

	if isEncrypted {
		delete(dst, "sops")
	}

	serializeNestedVariables(dst, v.vars, "")

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
func serializeNestedVariables(in map[string]interface{}, out map[string]string, prefix string) {
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
		case map[string]any:
			serializeNestedVariables(v, out, key)
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
