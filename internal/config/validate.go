package config

import (
	"embed"
	"fmt"
	"os"

	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*
var schemas embed.FS

func ValidateConfig(data []byte, version int) (bool, error) {
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
		return false, err
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

func loadSchema(version int) (*gojsonschema.JSONLoader, error) {
	body, err := schemas.ReadFile(fmt.Sprintf("schemas/schema-%d.yaml", version))
	if err != nil {
		return nil, err
	}
	return newYamlLoader(body)
}

func newYamlLoader(data []byte) (*gojsonschema.JSONLoader, error) {
	var document map[string]interface{}
	if err := yaml.Unmarshal(data, &document); err != nil {
		return nil, err
	}
	loader := gojsonschema.NewRawLoader(document)

	return &loader, nil
}
