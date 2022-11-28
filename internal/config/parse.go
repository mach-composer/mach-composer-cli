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
	"github.com/elliotchance/pie/v2"
	"gopkg.in/yaml.v3"

	"github.com/labd/mach-composer/internal/plugins"
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
func parseConfig(ctx context.Context, data []byte, vars *Variables, filename string) (*MachConfig, error) {
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
		err := addFileToConfig(cfg, intermediate.MachComposer.VariablesFile)
		if err != nil {
			return nil, err
		}
	}

	if err := parseGlobalNode(cfg, &intermediate.Global); err != nil {
		return nil, fmt.Errorf("failed to parse global node: %w", err)
	}

	if err := parseSitesNode(cfg, &intermediate.Sites); err != nil {
		return nil, fmt.Errorf("failed to parse sites node: %w", err)
	}

	if err := parseComponentsNode(cfg, &intermediate.Components, filename); err != nil {
		return nil, fmt.Errorf("failed to parse components node: %w", err)
	}

	return cfg, nil
}

func parseGlobalNode(cfg *MachConfig, globalNode *yaml.Node) error {
	if err := globalNode.Decode(&cfg.Global); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	knownKeys := []string{"cloud", "terraform_config", "environment"}
	nodes := mapYamlNodes(globalNode.Content)
	for key, node := range nodes {
		if pie.Contains(knownKeys, key) {
			continue
		}

		data, err := nodeAsMap(node)
		if err != nil {
			return err
		}
		if err := cfg.Plugins.SetGlobalConfig(key, data); err != nil {
			return err
		}
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

func parseSitesNode(cfg *MachConfig, sitesNode *yaml.Node) error {
	if err := sitesNode.Decode(&cfg.Sites); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	cloudPlugin, err := cfg.Plugins.Get(cfg.Global.Cloud)
	if err != nil {
		return err
	}

	knownKeys := []string{
		"name", "identifier", "endpoints", "components",
	}
	for _, site := range sitesNode.Content {
		nodes := mapYamlNodes(site.Content)
		siteId := nodes["identifier"].Value

		for key, node := range nodes {
			if pie.Contains(knownKeys, key) {
				continue
			}
			data, err := nodeAsMap(node)
			if err != nil {
				return err
			}

			if err := cfg.Plugins.SetSiteConfig(key, siteId, data); err != nil {
				return err
			}
		}

		if node, ok := nodes["endpoints"]; ok {
			data, err := nodeAsMap(node)
			if err != nil {
				return err
			}
			if err := cloudPlugin.SetSiteEndpointsConfig(siteId, data); err != nil {
				return err
			}
		}

		if err := parseSiteComponentsNode(cfg, siteId, nodes["components"]); err != nil {
			return err
		}
	}

	return nil
}

func parseSiteComponentsNode(cfg *MachConfig, site string, node *yaml.Node) error {
	knownKeys := []string{
		"name", "variables", "secrets",
		"store_variables", "store_secrets",
	}
	for _, component := range node.Content {
		nodes := mapYamlNodes(component.Content)
		identifier := nodes["name"].Value
		migrateCommercetools(identifier, nodes)

		for key, node := range nodes {
			if pie.Contains(knownKeys, key) {
				continue
			}
			data, err := nodeAsMap(node)
			if err != nil {
				return err
			}
			if err := cfg.Plugins.SetSiteComponentConfig(site, identifier, key, data); err != nil {
				return err
			}
		}
	}
	return nil
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

func parseComponentsNode(cfg *MachConfig, node *yaml.Node, source string) error {
	if node.Tag == "!!str" {
		path := filepath.Dir(source)
		var err error
		node, err = loadComponentsNode(node, path)
		if err != nil {
			return err
		}
	}

	if err := node.Decode(&cfg.Components); err != nil {
		return fmt.Errorf("decoding error: %w", err)
	}

	cloudPlugin, err := cfg.Plugins.Get(cfg.Global.Cloud)
	if err != nil {
		return err
	}
	for i := range cfg.Components {
		c := &cfg.Components[i]
		cloudPlugin.SetComponentEndpointsConfig(c.Name, c.Endpoints)
	}

	knownKeys := []string{
		"name", "source", "version", "branch",
		"integrations", "endpoints",
	}
	for _, component := range node.Content {
		nodes := mapYamlNodes(component.Content)
		identifier := nodes["name"].Value

		for key, node := range nodes {
			if pie.Contains(knownKeys, key) {
				continue
			}

			data, err := nodeAsMap(node)
			if err != nil {
				return err
			}

			if err := cfg.Plugins.SetComponentConfig(key, identifier, data); err != nil {
				return err
			}
		}
	}

	return nil
}

func loadComponentsNode(node *yaml.Node, path string) (*yaml.Node, error) {
	re := regexp.MustCompile(`\$\{include\(([^)]+)\)\}`)
	data := re.FindStringSubmatch(node.Value)
	if len(data) != 2 {
		return nil, fmt.Errorf("failed to parse ${include()} tag")
	}
	filename := filepath.Join(path, data[1])
	body, err := utils.AFS.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	result := yaml.Node{}
	if err = yaml.Unmarshal(body, &result); err != nil {
		return nil, err
	}
	if len(result.Content) != 1 {
		return nil, fmt.Errorf("Invalid yaml file")
	}
	return result.Content[0], nil
}

func nodeAsMap(n *yaml.Node) (map[string]any, error) {
	target := map[string]any{}
	if err := n.Decode(&target); err != nil {
		return nil, err
	}
	return target, nil
}

// migrateCommercetools moves the store_variables and store_secrets under the
// commercetools node. Needed to say backwards compatible
func migrateCommercetools(name string, nodes map[string]*yaml.Node) {
	needsMigrate := false
	if _, ok := nodes["store_variables"]; ok {
		needsMigrate = true
	}
	if _, ok := nodes["store_secrets"]; ok {
		needsMigrate = true
	}
	if !needsMigrate {
		return
	}

	fmt.Printf("Warning: %s move store_variables and store_secrets to commercetools node\n", name)

	if _, ok := nodes["commercetools"]; !ok {
		nodes["commercetools"] = &yaml.Node{
			Kind:    yaml.MappingNode,
			Tag:     "!!map",
			Content: []*yaml.Node{},
		}

		if val, ok := nodes["store_variables"]; ok {
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "store_variables",
			}
			nodes["commercetools"].Content = append(
				nodes["commercetools"].Content,
				keyNode,
				val,
			)
		}
		if val, ok := nodes["store_secrets"]; ok {
			keyNode := &yaml.Node{
				Kind:  yaml.ScalarNode,
				Tag:   "!!str",
				Value: "store_secrets",
			}
			nodes["commercetools"].Content = append(
				nodes["commercetools"].Content,
				keyNode,
				val,
			)
		}
	}
}

func loadDefaultPlugins(r *plugins.PluginRepository) {
	plugins := []string{"amplience", "aws", "azure", "contentful", "commercetools", "sentry"}
	for _, name := range plugins {
		if err := r.Load(name, "internal"); err != nil {
			panic(err)
		}
	}
}
