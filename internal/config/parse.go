package config

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/creasty/defaults"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
	"github.com/labd/mach-composer/internal/utils"
	"github.com/labd/mach-composer/internal/variables"
)

func Load(ctx context.Context, filename string, varFilename string) (*MachConfig, error) {
	var vars *variables.Variables
	if varFilename != "" {
		var err error
		vars, err = variables.NewVariablesFromFile(ctx, varFilename)
		if err != nil {
			return nil, err
		}
	}

	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	schemaVersion, err := getSchemaVersion(body)
	if err != nil {
		return nil, err
	}

	isValid, err := validateConfig(body, schemaVersion)
	if err != nil {
		return nil, err
	}

	if !isValid {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	cfg, err := ParseConfig(ctx, body, vars, filename)
	if err != nil {
		return nil, err
	}

	if err := defaults.Set(cfg); err != nil {
		return nil, err
	}

	cfg.Filename = filepath.Base(filename)
	return cfg, nil
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

// parseConfig is responsible for parsing a mach composer yaml config file and
// creating the resulting MachConfig struct.
func ParseConfig(ctx context.Context, data []byte, vars *variables.Variables, filename string) (*MachConfig, error) {
	// Decode the yaml in an intermediate config file
	intermediate := &_RawMachConfig{}
	err := yaml.Unmarshal(data, intermediate)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall yaml: %w", err)
	}

	pluginRepo := plugins.NewPluginRepository()
	if len(intermediate.MachComposer.Plugins) > 0 {
		for name, version := range intermediate.MachComposer.Plugins {
			err := pluginRepo.Load(name, version)
			if err != nil {
				return nil, err
			}
		}
	} else {
		loadDefaultPlugins(pluginRepo)
	}

	vars, err = processVariables(ctx, vars, intermediate)
	if err != nil {
		return nil, err
	}

	cfg := NewMachConfig()
	cfg.Plugins = pluginRepo
	cfg.Filename = intermediate.Filename
	cfg.MachComposer = intermediate.MachComposer
	cfg.Variables = vars

	if intermediate.Sops.Kind == yaml.MappingNode {
		cfg.IsEncrypted = true
	}

	if vars.Encrypted {
		err := cfg.addFileToConfig(intermediate.MachComposer.VariablesFile)
		if err != nil {
			return nil, err
		}
	}

	if err := parseGlobalNode(cfg, &intermediate.Global); err != nil {
		return nil, fmt.Errorf("failed to parse global node: %w", err)
	}

	if err := parseComponentsNode(cfg, &intermediate.Components, filename); err != nil {
		return nil, fmt.Errorf("failed to parse components node: %w", err)
	}

	if err := parseSitesNode(cfg, &intermediate.Sites); err != nil {
		return nil, fmt.Errorf("failed to parse sites node: %w", err)
	}

	return cfg, nil
}

func parseGlobalNode(cfg *MachConfig, globalNode *yaml.Node) error {
	if err := globalNode.Decode(&cfg.Global); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	knownKeys := []string{"cloud", "terraform_config", "environment"}
	nodes := mapYamlNodes(globalNode.Content)

	err := iterateYamlNodes(nodes, knownKeys, func(key string, data map[string]any) error {
		return cfg.Plugins.SetGlobalConfig(key, data)
	})
	if err != nil {
		return err
	}

	if node, ok := nodes["terraform_config"]; ok {
		childs := mapYamlNodes(node.Content)

		// Backwards compat
		if child, ok := childs["aws_remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}
			if err := cfg.Plugins.SetRemoteState("aws", data); err != nil {
				return err
			}
			cfg.Global.TerraformStateProvider = "aws"
		}
		if child, ok := childs["azure_remote_state"]; ok {
			data, err := nodeAsMap(child)
			if err != nil {
				return err
			}
			if err := cfg.Plugins.SetRemoteState("azure", data); err != nil {
				return err
			}
			cfg.Global.TerraformStateProvider = "azure"
		}
	}

	return nil
}

func nodeAsMap(n *yaml.Node) (map[string]any, error) {
	target := map[string]any{}
	if err := n.Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}

func loadDefaultPlugins(r *plugins.PluginRepository) {
	plugins := []string{"amplience", "aws", "azure", "contentful", "commercetools", "sentry"}
	for _, name := range plugins {
		if err := r.Load(name, "internal"); err != nil {
			panic(err)
		}
	}
}

func processVariables(ctx context.Context, vars *variables.Variables, rawConfig *_RawMachConfig) (*variables.Variables, error) {
	if vars == nil && rawConfig.MachComposer.VariablesFile != "" {
		var err error
		vars, err = variables.NewVariablesFromFile(ctx, rawConfig.MachComposer.VariablesFile)
		if err != nil {
			return nil, err
		}
	}

	if vars == nil {
		vars = variables.NewVariables()
	}

	if err := vars.InterpolateNode(&rawConfig.Sites); err != nil {
		return nil, err
	}
	if err := vars.InterpolateNode(&rawConfig.Components); err != nil {
		return nil, err
	}

	return vars, nil
}
