package config

import (
	"context"
	"fmt"
	"path"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
	"github.com/labd/mach-composer/internal/utils"
	"github.com/labd/mach-composer/internal/variables"
)

type ConfigOptions struct {
	VarFilenames []string
	Plugins      *plugins.PluginRepository

	NoResolveVars bool
}

func Open(ctx context.Context, filename string, opts *ConfigOptions) (*MachConfig, error) {
	raw, err := loadConfig(ctx, filename, opts.Plugins)
	if err != nil {
		return nil, err
	}

	// Validate again
	isValid, err := validateCompleteConfig(raw)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	for _, f := range opts.VarFilenames {
		if err := raw.variables.Load(ctx, f); err != nil {
			return nil, err
		}
	}

	// For some actions we don't want to resolve variables since they then need
	// to be passed as argument.
	if !opts.NoResolveVars {
		if err := resolveVariables(ctx, raw); err != nil {
			if notFoundErr, ok := err.(*variables.NotFoundError); ok {
				err = &SyntaxError{
					message:  fmt.Sprintf("unable to resolve variable %#v", notFoundErr.Name),
					line:     notFoundErr.Node.Line,
					filename: raw.filename,
					column:   notFoundErr.Node.Column,
				}
			}
			return nil, err
		}
	}

	return resolveConfig(ctx, raw)
}

func loadConfig(ctx context.Context, filename string, pr *plugins.PluginRepository) (*rawConfig, error) {
	// Load the yaml file and do basic validation if the config file is valid
	// based on a json schema
	document, err := loadYamlFile(filename)
	if err != nil {
		return nil, err
	}

	// Initial validation. We validate the document twice, once only the
	// structure and later again when we loaded the plugins
	isValid, err := validateConfig(document)
	if err != nil {
		return nil, err
	}
	if !isValid {
		return nil, fmt.Errorf("failed to load config %s due to errors", filename)
	}

	// Decode the yaml in an intermediate config file
	raw, err := newRawConfig(filename, document)
	if err != nil {
		return nil, err
	}

	if err := loadRefData(ctx, &raw.Components, path.Dir(filename)); err != nil {
		return nil, err
	}

	// Load the plugins
	raw.plugins = pr
	if raw.plugins == nil {
		raw.plugins = plugins.NewPluginRepository()
		if err := raw.plugins.Load(raw.MachComposer.Plugins); err != nil {
			return nil, err
		}
	}
	return raw, nil
}

// parseConfig is responsible for parsing a mach composer yaml config file and
// creating the resulting MachConfig struct.
func resolveConfig(ctx context.Context, intermediate *rawConfig) (*MachConfig, error) {
	if err := intermediate.validate(); err != nil {
		return nil, err
	}

	cfg := &MachConfig{
		extraFiles:   make(map[string][]byte, 0),
		Filename:     filepath.Base(intermediate.filename),
		MachComposer: intermediate.MachComposer,
		Variables:    intermediate.variables,
		Plugins:      intermediate.plugins,
	}

	for _, f := range cfg.Variables.EncryptedFiles {
		err := cfg.addFileToConfig(f)
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

	if err := parseComponentsNode(cfg, &intermediate.Components, intermediate.filename); err != nil {
		return nil, fmt.Errorf("failed to parse components node: %w", err)
	}

	if err := parseSitesNode(cfg, &intermediate.Sites); err != nil {
		return nil, fmt.Errorf("failed to parse sites node: %w", err)
	}

	return cfg, nil
}

func resolveVariables(ctx context.Context, rawConfig *rawConfig) error {
	vars := rawConfig.variables

	if rawConfig.MachComposer.VariablesFile != "" {
		if err := vars.Load(ctx, rawConfig.MachComposer.VariablesFile); err != nil {
			return err
		}
	}

	if err := vars.InterpolateNode(&rawConfig.Global); err != nil {
		return err
	}

	if err := vars.InterpolateNode(&rawConfig.Sites); err != nil {
		return err
	}

	if err := vars.InterpolateNode(&rawConfig.Components); err != nil {
		return err
	}

	return nil
}

func loadYamlFile(filename string) (*yaml.Node, error) {
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
	return document, nil
}
