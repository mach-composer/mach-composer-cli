package config

import (
	"log"

	"gopkg.in/yaml.v3"
)

func Parse(data []byte) (*Root, error) {
	root := Root{}

	err := yaml.Unmarshal(data, &root)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &root, nil
}
