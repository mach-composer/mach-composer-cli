package config

import (
	"context"
	"fmt"
	"path/filepath"

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

	// Read the config file from the given filename
	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	// Basic validation if the config file is valid based on a json schema
	isValid, err := validateConfig(body)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	return ParseConfig(ctx, body, vars, filename)
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

	pluginRepo, err := resolvePluginConfig(intermediate.MachComposer.Plugins)
	if err != nil {
		return nil, err
	}

	vars, err = processVariables(ctx, vars, intermediate)
	if err != nil {
		return nil, err
	}

	cfg := NewMachConfig()
	cfg.Filename = filepath.Base(filename)
	cfg.MachComposer = intermediate.MachComposer
	cfg.Plugins = pluginRepo
	cfg.Variables = vars

	if vars.Encrypted {
		err := cfg.addFileToConfig(vars.Filepath)
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

// resolvePluginConfig loads the plugins
func resolvePluginConfig(data map[string]string) (*plugins.PluginRepository, error) {
	pluginRepo := plugins.NewPluginRepository()
	if len(data) > 0 {
		for name, version := range data {
			err := pluginRepo.Load(name, version)
			if err != nil {
				return nil, err
			}
		}
	} else {
		pluginRepo.LoadDefault()
	}
	pluginRepo.StartPlugins()

	return pluginRepo, nil
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

	if err := vars.InterpolateNode(&rawConfig.Global); err != nil {
		return nil, err
	}

	if err := vars.InterpolateNode(&rawConfig.Sites); err != nil {
		return nil, err
	}

	if err := vars.InterpolateNode(&rawConfig.Components); err != nil {
		return nil, err
	}

	return vars, nil
}
