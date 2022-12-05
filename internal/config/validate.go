package config

import (
	"embed"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*
var schemas embed.FS

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
	if val, err := strconv.Atoi(v); err == nil {
		return val, err
	}

	parts := strings.SplitN(v, ".", 2)
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

func newYamlLoader(document *yaml.Node) (*gojsonschema.JSONLoader, error) {
	var data map[string]any
	if err := document.Decode(&data); err != nil {
		return nil, fmt.Errorf("yaml unmarshalling failed: %w", err)
	}
	loader := gojsonschema.NewRawLoader(data)

	return &loader, nil
}
