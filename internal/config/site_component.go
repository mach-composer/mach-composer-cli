package config

import (
	"fmt"
	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
	"path/filepath"
	"regexp"
	"slices"
	"strings"
)

var varComponentRegex = regexp.MustCompile(`\${(component(?:\.[^}]+)+)}`)

type SiteComponentConfigs []SiteComponentConfig

func (s *SiteComponentConfigs) Get(name string) (*SiteComponentConfig, error) {
	for _, site := range *s {
		if site.Name == name {
			return &site, nil
		}
	}
	return nil, fmt.Errorf("site component %s not found", name)
}

type SiteComponentConfig struct {
	Name       string            `yaml:"name"`
	Definition *ComponentConfig  `yaml:"-"`
	Variables  SiteComponentVars `yaml:"variables"`
	Secrets    SiteComponentVars `yaml:"secrets"`

	DependsOn []string `yaml:"depends_on"`
}

type SiteComponentVars map[string]*SiteComponentVar

func (scv *SiteComponentVars) Interpolated() map[string]any {
	var original = map[string]any{}

	for K, v := range *scv {
		original[K] = v.Interpolated
	}

	return original
}

func (scv *SiteComponentVars) ListComponents() []string {
	var references []string

	for _, v := range *scv {
		references = append(references, v.References...)
	}

	return slices.Compact(references)
}

type SiteComponentVar struct {
	Value        any
	Interpolated any
	References   []string
}

func (s *SiteComponentVar) UnmarshalYAML(value *yaml.Node) error {
	var val any
	err := value.Decode(&val)
	if err != nil {
		return err
	}

	s.Value = val
	s.Interpolated, s.References, err = interpolateComponentVar(val)
	if err != nil {
		return err
	}
	return nil
}

type ComponentConfig struct {
	Name         string            `yaml:"name"`
	Source       string            `yaml:"source"`
	Paths        []string          `yaml:"paths"`
	Version      string            `yaml:"version"`
	Branch       string            `yaml:"branch"`
	Integrations []string          `yaml:"integrations"`
	Endpoints    map[string]string `yaml:"endpoints"`
}

type TerraformConfig struct {
	Providers   map[string]string `yaml:"providers"`
	RemoteState map[string]any    `yaml:"remote_state"`
}

func (sc *SiteComponentConfig) HasCloudIntegration(g *GlobalConfig) bool {
	if sc.Definition == nil {
		log.Fatal().Msgf("ComponentConfig %s was not resolved properly (missing definition)", sc.Name)
	}
	return pie.Contains(sc.Definition.Integrations, g.Cloud)
}

// IsGitSource indicates if the source definition refers to Git.
func (c *ComponentConfig) IsGitSource() bool {
	return strings.HasPrefix(c.Source, "git")
}

func interpolateComponentVar(value any) (any, []string, error) {
	var result any
	var refs []string

	switch v := value.(type) {
	case string:
		parsed, ref, err := parseComponentVariable(v)
		if err != nil {
			return nil, refs, err
		}
		result = parsed
		refs = append(refs, ref...)
	case []any:
		var slice []any
		for _, element := range v {
			ipElement, newRefs, err := interpolateComponentVar(element)
			if err != nil {
				return nil, refs, err
			}
			slice = append(slice, ipElement)
			refs = append(refs, slices.Compact(newRefs)...)
		}

		result = slice
	case map[string]any:
		var ipMap = map[string]any{}
		for key, element := range v {
			ipElement, newRefs, err := interpolateComponentVar(element)
			if err != nil {
				return nil, refs, err
			}
			ipMap[key] = ipElement
			refs = append(refs, slices.Compact(newRefs)...)
		}

		result = ipMap
	default:
		result = v
	}

	return result, refs, nil
}

// parseTemplateVariable replaces `${component.my-component.my-var}` with the
// terraform string `module.my-component.my-var` so that components can uses
// output of one component for another component.
func parseComponentVariable(raw string) (string, []string, error) {
	val := strings.TrimSpace(raw)
	var refs []string
	matches := varComponentRegex.FindAllStringSubmatch(val, 20)
	if len(matches) == 0 {
		result := val
		return result, refs, nil
	}

	for _, match := range matches {
		parts := strings.SplitN(match[1], ".", 3)
		if len(parts) < 3 {
			return "", refs, fmt.Errorf(
				"invalid variable '%s'; "+
					"When using a ${component...} variable it has to consist of 2 parts; "+
					"component-name.output-name",
				match[1])
		}

		replacement := fmt.Sprintf("${data.terraform_remote_state.%s.outputs.%s}", parts[1], parts[2])
		val = strings.ReplaceAll(val, match[0], replacement)
		refs = append(refs, parts[1])
	}

	return val, refs, nil
}

func parseComponentsNode(cfg *MachConfig, node *yaml.Node) error {
	if err := node.Decode(&cfg.Components); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	if err := verifyComponents(cfg); err != nil {
		return fmt.Errorf("verify of components failed: %w", err)
	}
	knownKeys := []string{
		"name", "source", "version", "branch", "integrations", "endpoints", "paths",
	}
	for _, component := range node.Content {
		nodes := mapYamlNodes(component.Content)
		componentName := nodes["name"].Value
		err := iterateYamlNodes(nodes, knownKeys, func(key string, data map[string]any) error {
			return cfg.Plugins.SetComponentConfig(key, componentName, data)
		})
		if err != nil {
			return err
		}
	}

	if err := registerComponentEndpoints(cfg); err != nil {
		return fmt.Errorf("register of components failed: %w", err)
	}

	return nil
}

func registerComponentEndpoints(cfg *MachConfig) error {
	var cloudPlugin schema.MachComposerPlugin
	if cfg.Global.Cloud != "" {
		var err error
		cloudPlugin, err = cfg.Plugins.Get(cfg.Global.Cloud)
		if err != nil {
			return err
		}
	}

	for i := range cfg.Components {
		c := &cfg.Components[i]
		if cloudPlugin == nil {
			if len(c.Endpoints) > 0 {
				log.Error().Msg("Unable to register component endpoints when no cloud provider is configured")
			}
			continue
		}
		if len(c.Endpoints) > 0 {
			err := cloudPlugin.SetComponentEndpointsConfig(c.Name, c.Endpoints)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Verify the components config and set default values where needed.
func verifyComponents(cfg *MachConfig) error {
	var seen []string
	for i := range cfg.Components {
		c := &cfg.Components[i]

		if c.Integrations == nil {
			c.Integrations = make([]string, 0)
		}

		// Make sure the component names are unique. Otherwise raise an error
		if pie.Contains(seen, c.Name) {
			return fmt.Errorf("component %s is duplicate", c.Name)
		}

		// If the component has no integrations (or now called plugins)
		// specified then set it to the cloud integration
		if cfg.Global.Cloud != "" && len(c.Integrations) < 1 {
			c.Integrations = append(c.Integrations, cfg.Global.Cloud)
		}

		// If the source is a relative locale path then transform it to an
		// absolute path (required for Terraform)
		if strings.HasPrefix(c.Source, ".") {
			if val, err := filepath.Abs(c.Source); err == nil {
				c.Source = val
			} else {
				return err
			}
		}

		seen = append(seen, c.Name)
	}

	return nil
}
