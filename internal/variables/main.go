package variables

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
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

const globalNodeContext = "__global__"

type FileSource struct {
	Filename  string
	Encrypted bool
}

type Value struct {
	val        string
	fileSource *FileSource
}

type Variables struct {
	vars map[string]Value

	// External files from which variables are used which need
	// to be copied to the destination
	fileSources []FileSource

	// Mapping of used file sources. The key of the map is the site identifier
	// when the variable is used in a specific site, or __global__ when used
	// in non-site specific nodes
	usedFileSources map[string][]*FileSource

	// For now we only support loading one external file. To support multilpe
	// we need to track variable sources to check if they are encrypted or not
	loadedFile bool
}

func NewVariables() *Variables {
	v := &Variables{
		vars:            make(map[string]Value),
		fileSources:     []FileSource{},
		usedFileSources: map[string][]*FileSource{},
	}
	return v
}

func (v *Variables) getValue(nc string, key string) (string, error) {
	if strings.HasPrefix(key, "var.") {
		trimmedKey := key[4:]

		variable, ok := v.vars[trimmedKey]
		if !ok {
			return "", &NotFoundError{Name: key}
		}

		if variable.fileSource == nil {
			return variable.val, nil
		}

		if variable.fileSource.Encrypted {
			if _, ok := v.usedFileSources[nc]; !ok {
				v.usedFileSources[nc] = []*FileSource{}
			}

			if !pie.Any(v.usedFileSources[nc], func(f *FileSource) bool {
				return f == variable.fileSource
			}) {
				v.usedFileSources[nc] = append(v.usedFileSources[nc], variable.fileSource)
			}
			result := fmt.Sprintf(`${data.sops_external.variables.data["%s"]}`, trimmedKey)
			return result, nil
		}

		return variable.val, nil
	}

	if strings.HasPrefix(key, "env.") {
		trimmedKey := key[4:]
		return os.Getenv(trimmedKey), nil
	}

	log.Warn().Msgf("Unsupported variables type %s", key)
	return "", nil
}

func (v *Variables) Set(key string, value string) {
	v.vars[key] = Value{val: value}
}

func (v *Variables) HasEncrypted(site string) bool {
	return pie.Any(v.GetEncryptedSources(site), func(f FileSource) bool { return f.Encrypted })
}

func (v *Variables) GetEncryptedSources(site string) []FileSource {
	items := []*FileSource{}
	if fs, ok := v.usedFileSources[globalNodeContext]; ok {
		items = append(items, fs...)
	}
	if fs, ok := v.usedFileSources[site]; ok {
		items = append(items, fs...)
	}

	return pie.Map(pie.Unique(items), func(f *FileSource) FileSource {
		return *f
	})
}

func (v *Variables) InterpolateNode(node *yaml.Node) error {
	return v.interpolateNodeContext(globalNodeContext, node)
}

func (v *Variables) InterpolateSiteNode(site string, node *yaml.Node) error {
	if site == globalNodeContext {
		return fmt.Errorf("invalid site identifier")
	}
	return v.interpolateNodeContext(site, node)
}

func (v *Variables) interpolateNodeContext(nc string, node *yaml.Node) error {
	if node.Kind == yaml.ScalarNode {
		val, err := v.interpolateValue(nc, node.Value)
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

		err := v.interpolateNodeContext(nc, node.Content[i])
		if err != nil {
			return err
		}
	}
	return nil
}

func (v *Variables) interpolateValue(nc string, val string) (string, error) {
	matches := varRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		return val, nil
	}

	for _, match := range matches {
		replacement, err := v.getValue(nc, match[1])
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
		delete(values, "sops")
	}

	fs := FileSource{
		Filename:  filename,
		Encrypted: isEncrypted,
	}
	v.fileSources = append(v.fileSources, fs)

	dst := map[string]Value{}
	serializeNestedVariables(values, dst, "")
	for key, val := range dst {
		val.fileSource = &fs
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
			out[key] = Value{val: v}
		case int:
			out[key] = Value{val: fmt.Sprint(v)}
		case map[string]any:
			serializeNestedVariables(v, out, key)
		}
	}
}
