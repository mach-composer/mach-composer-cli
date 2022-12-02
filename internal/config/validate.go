package config

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*
var schemas embed.FS

func validateConfig(data []byte) (bool, error) {
	version, err := getSchemaVersion(data)
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

	docLoader, err := newYamlLoader(data)
	if err != nil {
		return false, err
	}

	result, err := gojsonschema.Validate(*schemaLoader, *docLoader)
	if err != nil {
		return false, fmt.Errorf("configuration file is invalid: %w", err)
	}

	// Deal with result
	if !result.Valid() {
		fmt.Fprintln(os.Stderr, "The config is not valid:")
		for _, desc := range result.Errors() {
			fmt.Fprintf(os.Stderr, " - %s\n", desc)
		}
		return false, nil
	}
	return true, nil
}

func getSchemaVersion(data []byte) (int, error) {
	type PartialMachConfig struct {
		MachComposer MachComposer `yaml:"mach_composer"`
	}

	// Decode the yaml in an intermediate config file
	intermediate := &PartialMachConfig{}
	err := yaml.Unmarshal(data, intermediate)
	if err != nil {
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
	return newYamlLoader(body)
}

func newYamlLoader(data []byte) (*gojsonschema.JSONLoader, error) {
	var document map[string]interface{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, fmt.Errorf("yaml unmarshalling failed: %w", err)
	}
	loader := gojsonschema.NewRawLoader(document)

	return &loader, nil
}
