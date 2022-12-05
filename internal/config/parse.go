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

type ParseOptions struct {
	Variables *variables.Variables
	Plugins   *plugins.PluginRepository
	Filename  string
}

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

	// Read the yaml nodes
	document := &yaml.Node{}
	if err := yaml.Unmarshal(body, document); err != nil {
		return nil, err
	}

	// Basic validation if the config file is valid based on a json schema
	isValid, err := validateConfig(document)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	return ParseConfig(ctx, document, ParseOptions{
		Variables: vars,
		Filename:  filename,
	})
}

// parseConfig is responsible for parsing a mach composer yaml config file and
// creating the resulting MachConfig struct.
func ParseConfig(ctx context.Context, document *yaml.Node, options ParseOptions) (*MachConfig, error) {
	// Decode the yaml in an intermediate config file
	intermediate := &_RawMachConfig{}
	if err := document.Decode(intermediate); err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %w", err)
	}

	vars, err := processVariables(ctx, options.Variables, intermediate)
	if err != nil {
		if notFoundErr, ok := err.(*variables.NotFoundError); ok {
			err = &SyntaxError{
				message:  fmt.Sprintf("unable to resolve variable %#v", notFoundErr.Name),
				line:     notFoundErr.Node.Line,
				filename: options.Filename,
				column:   notFoundErr.Node.Column,
			}
		}
		return nil, err
	}

	cfg := NewMachConfig()
	cfg.Filename = filepath.Base(options.Filename)
	cfg.MachComposer = intermediate.MachComposer
	cfg.Variables = vars

	cfg.Plugins = options.Plugins
	if cfg.Plugins == nil {
		cfg.Plugins = plugins.NewPluginRepository()
	}
	if err := cfg.Plugins.Load(intermediate.MachComposer.Plugins); err != nil {
		return nil, err
	}

	if vars.Encrypted {
		err := cfg.addFileToConfig(vars.Filepath)
		if err != nil {
			return nil, err
		}
	}

	if err := parseGlobalNode(cfg, &intermediate.Global); err != nil {
		if _, ok := err.(*plugins.PluginNotFoundError); ok {
			return nil, err
		}
		return nil, fmt.Errorf("failed to parse global node: %w", err)
	}

	if err := parseComponentsNode(cfg, &intermediate.Components, options.Filename); err != nil {
		return nil, fmt.Errorf("failed to parse components node: %w", err)
	}

	if err := parseSitesNode(cfg, &intermediate.Sites); err != nil {
		return nil, fmt.Errorf("failed to parse sites node: %w", err)
	}

	return cfg, nil
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
