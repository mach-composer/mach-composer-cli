package plugins

import (
	"context"
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/rs/zerolog/log"
)

type PluginNotFoundError struct {
	name string
}

func (e *PluginNotFoundError) Error() string {
	return fmt.Sprintf("plugin %s not found", e.name)
}

type PluginError struct {
	msg    string
	plugin string
}

func (e *PluginError) Error() string {
	return e.msg
}

type PluginRepository struct {
	plugins map[string]*Plugin
	schemas map[string]schema.ValidationSchema
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		plugins: make(map[string]*Plugin),
		schemas: make(map[string]schema.ValidationSchema),
	}
}

// Close kills alls the running plugins
func (p *PluginRepository) Close() {
	for _, rp := range p.plugins {
		rp.client.Kill()
		rp.isRunning = false
		rp.client = nil
	}
}

// All returns all the plugin names in the repository, ordered by the plugin name
func (p *PluginRepository) All() []Plugin {
	result := make([]Plugin, len(p.plugins))
	for i, key := range pie.Sort(pie.Keys(p.plugins)) {
		result[i] = *p.plugins[key]
	}
	return result
}

// Add an existing plugin to the repository. Only used for testing purposes
func (p *PluginRepository) Add(name string, plugin schema.MachComposerPlugin) error {
	p.plugins[name] = &Plugin{
		MachComposerPlugin: plugin,
		Name:               name,
		isRunning:          true,
	}
	return nil
}

// Get returns the plugin from the repository
func (p *PluginRepository) Get(name string) (*Plugin, error) {
	if name == "" {
		panic("plugin name is empty") // this should never happen
	}
	plugin, ok := p.plugins[name]
	if !ok {
		return nil, &PluginNotFoundError{name: name}
	}
	return plugin, nil
}

// LoadPlugin will load the plugin with the given name and start it
func (p *PluginRepository) LoadPlugin(ctx context.Context, name string, config PluginConfig) error {
	if _, ok := p.plugins[name]; ok {
		return fmt.Errorf("plugin %s is already loaded", name)
	}

	plug := &Plugin{
		Name:   name,
		config: config,
	}
	p.plugins[name] = plug

	if err := plug.start(ctx); err != nil {
		return err
	}
	return nil
}

func (p *PluginRepository) loadSchema(name string) (*schema.ValidationSchema, error) {
	plugin, err := p.Get(name)
	if err != nil {
		return nil, err
	}

	// Load validation schema
	schema, err := plugin.GetValidationSchema()
	if err != nil {
		return nil, fmt.Errorf("failed to load plugin schema %s: %w", name, err)
	}
	if schema != nil {
		p.schemas[name] = *schema
	}
	return schema, nil
}

func (p *PluginRepository) GetSchema(name string) (*schema.ValidationSchema, error) {
	if name == "" {
		panic("plugin name is empty") // this should never happen
	}

	// this should not happen in a regular use-case
	if schema, ok := p.schemas[name]; ok {
		return &schema, nil
	}

	return p.loadSchema(name)
}

func (p *PluginRepository) SetRemoteState(name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleError(name, plugin.SetRemoteStateBackend(data))
}

func (p *PluginRepository) SetGlobalConfig(name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleError(name, plugin.SetGlobalConfig(data))
}

func (p *PluginRepository) SetSiteEndpointConfig(name string, site string, key string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleError(name, plugin.SetSiteEndpointConfig(site, key, data))
}

func (p *PluginRepository) SetComponentConfig(name, component string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleError(name, plugin.SetComponentConfig(component, data))
}

func (p *PluginRepository) handleError(plugin string, err error) error {
	if err == nil {
		return nil
	}

	log.Error().Err(err).Stack().Msgf("plugin %s returned an error", plugin)

	if strings.Contains(err.Error(), "reading body EOF") {
		return &PluginError{
			plugin: plugin,
			msg:    fmt.Sprintf("the %s plugin crashed. This is a bug in the plugin.", plugin),
		}
	}
	return &PluginError{
		plugin: plugin,
		msg:    err.Error(),
	}
}
