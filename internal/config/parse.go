package config

import (
	"context"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"strings"

	"github.com/mach-composer/mach-composer-cli/internal/state"

	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type ConfigOptions struct {
	VarFilenames []string
	Plugins      *plugins.PluginRepository

	Validate bool

	NoResolveVars bool
}

// Open is the main entrypoint for this module. It opens the given yaml filename
// and reads it to construct the MachConfig.
// Note that you need to close the MachConfig via the Close() method in order
// to clean up.
func Open(ctx context.Context, filename string, opts *ConfigOptions) (*MachConfig, error) {
	var pluginRepo *plugins.PluginRepository
	if opts != nil {
		pluginRepo = opts.Plugins
	}

	//Take the relative path of the config file as the working directory
	cwd := path.Dir(filename)

	raw, err := loadConfig(ctx, filename, cwd, pluginRepo, opts.Validate)
	if err != nil {
		return nil, err
	}

	// Validate again
	if opts.Validate {
		isValid, err := validateCompleteConfig(raw)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, fmt.Errorf("failed to load config %s due to errors", filename)
		}
	}

	for _, f := range opts.VarFilenames {
		if err := raw.variables.Load(ctx, f, cwd); err != nil {
			return nil, err
		}
	}

	// For some actions we don't want to resolve variables since they then need
	// to be passed as argument.
	if !opts.NoResolveVars {
		if err := resolveVariables(ctx, raw, cwd); err != nil {
			var notFoundErr *NotFoundError
			if errors.As(err, &notFoundErr) {
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

func loadConfig(ctx context.Context, filename, cwd string, pr *plugins.PluginRepository, validate bool) (*rawConfig, error) {
	// Load the yaml file and do basic validation if the config file is valid
	// based on a json schema
	document, err := loadYamlFile(filename)
	if err != nil {
		return nil, err
	}

	// Initial validation. We validate the document twice, once only the
	// structure and later again when we loaded the plugins
	if validate {
		isValid, err := validateConfig(document)
		if err != nil {
			return nil, err
		}
		if !isValid {
			return nil, fmt.Errorf("failed to load config %s due to errors", filename)
		}
	}

	// Decode the yaml in an intermediate config file
	raw, err := newRawConfig(filename, document)
	if err != nil {
		return nil, err
	}

	// Load the plugins
	raw.plugins = pr
	if err := loadPlugins(ctx, raw); err != nil {
		return nil, err
	}
	return raw, nil
}

func loadPlugins(ctx context.Context, raw *rawConfig) error {
	if raw.plugins != nil {
		return nil
	}
	raw.plugins = plugins.NewPluginRepository()

	if len(raw.MachComposer.Plugins) == 0 {
		log.Debug().Msg("No plugins specified; loading default plugins")
		raw.MachComposer.Plugins = map[string]MachPluginConfig{
			"amplience": {
				Source:  "mach-composer/amplience",
				Version: "0.1.3",
			},
			"aws": {
				Source:  "mach-composer/aws",
				Version: "0.1.0",
			},
			"azure": {
				Source:  "mach-composer/azure",
				Version: "0.1.0",
			},
			"commercetools": {
				Source:  "mach-composer/commercetools",
				Version: "0.1.5",
			},
			"contentful": {
				Source:  "mach-composer/contentful",
				Version: "0.1.0",
			},
			"sentry": {
				Source:  "mach-composer/sentry",
				Version: "0.1.2",
			},
		}
	}

	for pluginName, pluginData := range raw.MachComposer.Plugins {
		pluginConfig := plugins.PluginConfig{
			Source:  pluginData.Source,
			Version: pluginData.Version,
			Replace: pluginData.Replace,
		}
		if err := raw.plugins.LoadPlugin(ctx, pluginName, pluginConfig); err != nil {
			return err
		}
	}
	return nil
}

// resolveConfig is responsible for parsing a mach composer yaml config file and creating the resulting MachConfig struct.
func resolveConfig(_ context.Context, intermediate *rawConfig) (*MachConfig, error) {
	if err := intermediate.validate(); err != nil {
		return nil, err
	}

	cfg := &MachConfig{
		StateRepository: state.NewRepository(),
		extraFiles:      make(map[string][]byte),
		Filename:        filepath.Base(intermediate.filename),
		MachComposer:    intermediate.MachComposer,
		Variables:       intermediate.variables,
		Plugins:         intermediate.plugins,
	}

	if err := parseGlobalNode(cfg, &intermediate.Global); err != nil {
		var pluginNotFoundError *plugins.PluginNotFoundError
		if errors.As(err, &pluginNotFoundError) {
			return nil, err
		}
		return nil, fmt.Errorf("failed to parse global node: %w", err)
	}

	if err := parseComponentsNode(cfg, &intermediate.Components); err != nil {
		return nil, fmt.Errorf("failed to parse components node: %w", err)
	}

	if err := parseSitesNode(cfg, &intermediate.Sites); err != nil {
		return nil, fmt.Errorf("failed to parse sites node: %w", err)
	}

	return cfg, nil
}

func resolveVariables(ctx context.Context, rawConfig *rawConfig, cwd string) error {
	vars := rawConfig.variables

	if rawConfig.MachComposer.VariablesFile != "" {
		if err := vars.Load(ctx, rawConfig.MachComposer.VariablesFile, cwd); err != nil {
			return err
		}
	}

	if err := vars.InterpolateNode(&rawConfig.Global); err != nil {
		return err
	}

	if err := vars.InterpolateNode(&rawConfig.Components); err != nil {
		return err
	}

	// TransformValue the variables per-site to keep track of which site uses which
	// variable.
	for _, node := range rawConfig.Sites.Content {
		mapping := MapYamlNodes(node.Content)
		if idNode, ok := mapping["identifier"]; ok {
			siteId := idNode.Value
			if err := vars.InterpolateSiteNode(siteId, node); err != nil {
				return err
			}
		}
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

	// Resolve $ref and ${include()} references
	if err := resolveReferences(document, filepath.Dir(filename)); err != nil {
		return nil, err
	}

	return document, nil
}

func resolveReferences(node *yaml.Node, baseDir string) error {
	if node.Kind == yaml.DocumentNode {
		for _, contentNode := range node.Content {
			if err := resolveReferences(contentNode, baseDir); err != nil {
				return err
			}
		}
	} else if node.Kind == yaml.MappingNode {
		for i := 0; i < len(node.Content); i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]
			if keyNode.Value == "$ref" && valueNode.Kind == yaml.ScalarNode {
				refFilename := filepath.Join(baseDir, valueNode.Value)
				refNode, err := loadYamlFile(refFilename)
				if err != nil {
					return err
				}
				contentNode := refNode.Content[0]
				*node = *contentNode
				return nil
			}
			if err := resolveReferences(valueNode, baseDir); err != nil {
				return err
			}
		}
	} else if node.Kind == yaml.SequenceNode {
		for _, contentNode := range node.Content {
			if err := resolveReferences(contentNode, baseDir); err != nil {
				return err
			}
		}
	} else if node.Kind == yaml.ScalarNode && strings.Contains(node.Value, "include") {
		refNode, _, err := LoadIncludeDocument(node, baseDir)
		if err != nil {
			return err
		}
		*node = *refNode
		return nil
	}
	return nil
}
