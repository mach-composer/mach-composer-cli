package config

import (
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"
)

func Load(filename string) (*MachConfig, error) {

	body, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	cfg, err := Parse(body)
	if err != nil {
		panic(err)
	}

	if err := defaults.Set(cfg); err != nil {
		panic(err)
	}

	cfg.Filename = filepath.Base(filename)
	Process(cfg)

	return cfg, nil
}

func Parse(data []byte) (*MachConfig, error) {
	cfg := MachConfig{}

	err := yaml.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return &cfg, nil
}
