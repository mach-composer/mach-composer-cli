package config

import (
	"embed"

	"github.com/sirupsen/logrus"
	"github.com/xeipuuv/gojsonschema"
	"gopkg.in/yaml.v3"
)

//go:embed schemas/*
var schemas embed.FS

func ValidateConfig(data []byte) bool {

	schemaLoader, err := loadSchema()
	if err != nil {
		panic(err)
	}

	docLoader, err := newYamlLoader(data)
	if err != nil {
		panic(err)
	}

	result, err := gojsonschema.Validate(*schemaLoader, *docLoader)
	if err != nil {
		panic(err)
	}

	// Deal with result
	if !result.Valid() {
		logrus.Error("The document is not valid:")
		for _, desc := range result.Errors() {
			logrus.Errorf(" - %s\n", desc)
		}
		return false
	}
	return true
}

func loadSchema() (*gojsonschema.JSONLoader, error) {

	body, err := schemas.ReadFile("schemas/schema-1.yaml")
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
