package config

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/creasty/defaults"
	"github.com/labd/mach-composer/internal/utils"
	"gopkg.in/yaml.v3"
)

func Load(filename string, varFilename string) (*MachConfig, error) {

	var vars *Variables
	if varFilename != "" {
		var err error
		vars, err = loadVariables(varFilename)
		if err != nil {
			panic(err)
		}
	}

	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	schemaVersion, err := GetSchemaVersion(body)
	if err != nil {
		return nil, err
	}

	if !ValidateConfig(body, schemaVersion) {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	cfg, err := Parse(body, vars)
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

func GetSchemaVersion(data []byte) (int, error) {

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

func Parse(data []byte, vars *Variables) (*MachConfig, error) {

	// Decode the yaml in an intermediate config file
	intermediate := &_RawMachConfig{}
	err := yaml.Unmarshal(data, intermediate)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall yaml: %w", err)
	}

	if vars == nil && intermediate.MachComposer.VariablesFile != "" {
		vars, err = loadVariables(intermediate.MachComposer.VariablesFile)
		if err != nil {
			panic(err)
		}
	}

	if vars == nil {
		vars = NewVariables()
	}

	varErr := InterpolateVars(intermediate, vars)
	if varErr != nil {
		return nil, varErr
	}

	cfg := &MachConfig{
		Filename:     intermediate.Filename,
		MachComposer: intermediate.MachComposer,
		Global:       intermediate.Global,
	}

	cfg.Variables = vars

	if intermediate.Sops.Kind == yaml.MappingNode {
		cfg.IsEncrypted = true
	}

	if err := intermediate.Sites.Decode(&cfg.Sites); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	if err := intermediate.Components.Decode(&cfg.Components); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	return cfg, nil
}
