package config

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/utils"
)

func Load(ctx context.Context, filename string, varFilename string) (*MachConfig, error) {
	var vars *Variables
	if varFilename != "" {
		var err error
		vars, err = loadVariables(ctx, varFilename)
		if err != nil {
			return nil, err
		}
	}

	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	schemaVersion, err := GetSchemaVersion(body)
	if err != nil {
		return nil, err
	}

	isValid, err := ValidateConfig(body, schemaVersion)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	cfg, err := parseConfig(ctx, body, vars, filename)
	if err != nil {
		return nil, err
	}

	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	cfg.Filename = filepath.Base(filename)
	if err := ProcessConfig(cfg); err != nil {
		return nil, err
	}

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

// parseConfig is responsible for parsing a mach composer yaml config file and
// creating the resulting MachConfig struct.
func parseConfig(ctx context.Context, data []byte, vars *Variables, filename string) (*MachConfig, error) {
	// Decode the yaml in an intermediate config file
	intermediate := &_RawMachConfig{}
	err := yaml.Unmarshal(data, intermediate)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall yaml: %w", err)
	}
	if vars == nil && intermediate.MachComposer.VariablesFile != "" {
		vars, err = loadVariables(ctx, intermediate.MachComposer.VariablesFile)
		if err != nil {
			return nil, err
		}
	}

	if vars == nil {
		vars = NewVariables()
	}

	varErr := InterpolateVars(intermediate, vars)
	if varErr != nil {
		return nil, varErr
	}

	cfg := NewMachConfig()
	cfg.Filename = intermediate.Filename
	cfg.MachComposer = intermediate.MachComposer
	cfg.Global = intermediate.Global
	cfg.Variables = vars

	if intermediate.Sops.Kind == yaml.MappingNode {
		cfg.IsEncrypted = true
	}

	if vars.Encrypted {
		err := addFileToConfig(cfg, intermediate.MachComposer.VariablesFile)
		if err != nil {
			return nil, err
		}
	}

	if err := intermediate.Sites.Decode(&cfg.Sites); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	if err := parseComponentsNode(intermediate.Components, filename, &cfg.Components); err != nil {
		return nil, fmt.Errorf("decoding error: %w", err)
	}

	return cfg, nil
}

func addFileToConfig(cfg *MachConfig, filename string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("error reading variables file: %w", err)
	}
	filename = filepath.Base(filename)
	cfg.ExtraFiles[filename] = b
	return nil
}

func parseComponentsNode(node yaml.Node, source string, target *[]Component) error {
	if node.Tag != "!!str" {
		if err := node.Decode(target); err != nil {
			return fmt.Errorf("decoding error: %w", err)
		}
		return nil
	}

	re := regexp.MustCompile(`\$\{include\(([^)]+)\)\}`)
	data := re.FindStringSubmatch(node.Value)
	if len(data) != 2 {
		return fmt.Errorf("failed to parse ${include()} tag")
	}
	filename := filepath.Join(filepath.Dir(source), data[1])

	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(body, &target)
	if err != nil {
		return err
	}
	return nil
}
