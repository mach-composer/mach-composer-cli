package config

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/cli"
	"github.com/mach-composer/mach-composer-cli/internal/plugins"
)

//go:embed schemas/*
var schemas embed.FS

func GenerateSchema(ctx context.Context, filename string, pr *plugins.PluginRepository) (string, error) {
	raw, err := loadConfig(ctx, filename, pr)
	if err != nil {
		return "", err
	}
	defer raw.plugins.Close()

	data, err := createFullSchema(raw.plugins, &raw.Global)
	if err != nil {
		return "", err
	}

	result, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return "", err
	}

	return string(result), nil
}

func validateConfig(document *yaml.Node) (bool, error) {
	version, err := getSchemaVersion(document)
	if err != nil {
		return false, err
	}

	if version != 1 {
		err := fmt.Errorf("Config version %d is unsupported. Only version 1 is supported.\n", version)
		return false, err
	}

	schemaLoader, err := loadSchema(version)
	if err != nil {
		return false, err
	}

	docLoader, err := newYamlLoader(document)
	if err != nil {
		return false, err
	}

	result, err := gojsonschema.Validate(*schemaLoader, *docLoader)
	if err != nil {
		return false, fmt.Errorf("configuration file is invalid: %w", err)
	}

	// Deal with result
	if !result.Valid() {
		err := &ValidationError{
			errors: []string{},
		}
		for _, desc := range result.Errors() {
			err.errors = append(err.errors, fmt.Sprintf("%s\n", desc))
		}
		return false, err
	}
	return true, nil
}

func validateCompleteConfig(raw *rawConfig) (bool, error) {
	schemaData, err := createFullSchema(raw.plugins, &raw.Global)
	if err != nil {
		return false, err
	}
	schemaLoader := gojsonschema.NewRawLoader(schemaData)

	docLoader, err := newYamlLoader(raw.document)
	if err != nil {
		return false, err
	}

	result, err := gojsonschema.Validate(schemaLoader, *docLoader)
	if err != nil {
		return false, fmt.Errorf("configuration file is invalid: %w", err)
	}

	// Deal with result
	if !result.Valid() {
		err := &ValidationError{
			errors: []string{},
		}
		for _, desc := range result.Errors() {
			err.errors = append(err.errors, fmt.Sprintf("%s\n", desc))
		}
		return false, err
	}
	return true, nil
}

func getSchemaVersion(document *yaml.Node) (int, error) {
	type PartialMachConfig struct {
		MachComposer MachComposer `yaml:"mach_composer"`
	}

	// Decode the yaml in an intermediate config file
	intermediate := &PartialMachConfig{}
	if err := document.Decode(intermediate); err != nil {
		return 0, err
	}

	v := intermediate.MachComposer.Version
	if _, ok := v.(int); !ok {
		cli.DeprecationWarning(&cli.DeprecationOptions{
			Message: fmt.Sprintf("mach_composer.version should be an integer (not %v)", v),
		})
	}

	vs := fmt.Sprintf("%v", v)
	if val, err := strconv.Atoi(vs); err == nil {
		return val, err
	}

	parts := strings.SplitN(vs, ".", 2)
	if val, err := strconv.Atoi(parts[0]); err == nil {
		return val, err
	}

	return 0, errors.New("no valid version identifier found")
}

func loadSchema(version int) (*gojsonschema.JSONLoader, error) {
	body, err := schemas.ReadFile(fmt.Sprintf("schemas/schema-%d.yaml", version))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}
	return newFileLoader(body)
}

func newFileLoader(document []byte) (*gojsonschema.JSONLoader, error) {
	var data map[string]any
	if err := yaml.Unmarshal(document, &data); err != nil {
		return nil, fmt.Errorf("yaml unmarshalling failed: %w", err)
	}
	loader := gojsonschema.NewRawLoader(data)

	return &loader, nil
}

// newYamlLoader allows us to validate yaml file with the gojsonschema. First we
// convert the nodes to a map[string]any and then we serialize it to a json
// string for validation. This extra serialization helps validation since
// yaml is rather lax regarding data-types (ints vs strings vs floats)
func newYamlLoader(document *yaml.Node) (*gojsonschema.JSONLoader, error) {
	var data map[string]any
	if err := document.Decode(&data); err != nil {
		return nil, fmt.Errorf("yaml unmarshalling failed: %w", err)
	}

	transformed := transformYamlData(data)

	body, err := json.Marshal(transformed)
	if err != nil {
		panic(err)
	}

	loader := gojsonschema.NewStringLoader(string(body))
	return &loader, nil
}

func createFullSchema(pr *plugins.PluginRepository, globalNode *yaml.Node) (map[string]any, error) {
	g := GlobalConfig{}
	if err := globalNode.Decode(&g); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	body, err := schemas.ReadFile(fmt.Sprintf("schemas/schema-%d.yaml", 1))
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	var data map[string]any
	if err := yaml.Unmarshal(body, &data); err != nil {
		return nil, fmt.Errorf("yaml unmarshalling failed: %w", err)
	}

	definitions, ok := data["definitions"].(map[string]any)
	if !ok {
		return nil, fmt.Errorf("unable to read schema")
	}

	statePluginName, ok := g.TerraformConfig.RemoteState["plugin"]
	if ok {
		stateSchema, err := state.GetSchema(state.Type(statePluginName.(string)))
		if err != nil {
			return nil, fmt.Errorf("unable to get state schema")
		}
		definitions["RemoteState"] = *stateSchema
	}

	// Disable additionalProperties
	setAdditionalProperties(definitions["GlobalConfig"], false)
	setAdditionalProperties(definitions["SiteConfig"], false)
	setAdditionalProperties(definitions["SiteComponentConfig"], false)
	setAdditionalProperties(definitions["SiteEndpointConfig"], false)
	setAdditionalProperties(definitions["ComponentConfig"], false)
	setAdditionalProperties(definitions["ComponentEndpointConfig"], false)

	// Site config
	for _, plugin := range pr.All() {
		schema, err := pr.GetSchema(plugin.Name)
		if err != nil {
			return nil, err
		}

		setObjectProperties(definitions["GlobalConfig"], plugin.Name, schema.GlobalConfigSchema)
		setObjectProperties(definitions["SiteConfig"], plugin.Name, schema.SiteConfigSchema)
		setObjectProperties(definitions["SiteComponentConfig"], plugin.Name, schema.SiteComponentConfigSchema)
		setObjectProperties(definitions["SiteEndpointConfig"], plugin.Name, schema.SiteEndpointConfig)
		setObjectProperties(definitions["ComponentConfig"], plugin.Name, schema.ComponentConfigSchema)
		setObjectProperties(definitions["ComponentEndpointConfig"], plugin.Name, schema.ComponentEndpointsConfigSchema)
	}

	return data, nil
}

func setAdditionalProperties(values any, value bool) {
	items, ok := values.(map[string]any)
	if !ok {
		panic("error parsing schema") // Program error
	}
	items["additionalProperties"] = value
}

func setObjectProperties(values any, name string, p map[string]any) {
	if len(p) < 1 {
		return
	}

	item, ok := values.(map[string]any)
	if !ok {
		panic("error parsing schema") // Program error
	}

	properties, ok := item["properties"].(map[string]any)
	if !ok {
		properties = map[string]any{}
		item["properties"] = properties
	}
	properties[name] = p
}

// transformYamlData returns the given data whereby keys of maps are always
// strings so that i can be serialized to json
func transformYamlData(d any) any {
	switch t := d.(type) {
	case map[string]any:
		result := map[string]any{}
		for k, v := range t {
			result[k] = transformYamlData(v)
		}
		return result
	case map[any]any:
		result := map[string]any{}
		for k, v := range t {
			key := fmt.Sprintf("%v", k)
			result[key] = transformYamlData(v)
		}
		return result
	case []any:
		for i, v := range t {
			t[i] = transformYamlData(v)
		}
	}
	return d
}
