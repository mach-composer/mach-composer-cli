package plugins

import (
	"fmt"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/rs/zerolog/log"
)

type PluginNotFoundError struct {
	Plugin string
}

type Plugin struct {
	schema.MachComposerPlugin
	Name string
}

type PluginError struct {
	msg    string
	plugin string
}

func (e *PluginError) Error() string {
	return e.msg
}

func (e *PluginNotFoundError) Error() string {
	return fmt.Sprintf("plugin %s not found", e.Plugin)
}

type PluginRepository struct {
	Plugins map[string]schema.MachComposerPlugin
	Schemas map[string]schema.ValidationSchema
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		Plugins: make(map[string]schema.MachComposerPlugin),
		Schemas: make(map[string]schema.ValidationSchema),
	}
}

// resolvePluginConfig loads the plugins
func (p *PluginRepository) Load(data map[string]map[string]string) error {
	if data == nil {
		log.Debug().Msg("No plugins specified; loading default plugins")
		return p.LoadDefault()
	} else {
		for name, properties := range data {
			err := p.LoadPlugin(name, properties)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// LoadDefault loads the default plugins, for backwards compatibility
func (p *PluginRepository) LoadDefault() error {
	// Don't load default plugins if we already have
	if len(p.Plugins) > 0 {
		return nil
	}

	for _, name := range LocalPluginNames {
		if err := p.LoadPlugin(name, map[string]string{}); err != nil {
			return err
		}
	}
	return nil
}

func (p *PluginRepository) LoadPlugin(name string, properties map[string]string) error {
	if _, ok := p.Plugins[name]; ok {
		return fmt.Errorf("plugin %s is already loaded", name)
	}

	plugin, err := startPlugin(name)
	if err != nil {
		return fmt.Errorf("failed to start plugin %s: %w", name, err)
	}
	p.Plugins[name] = plugin
	p.LoadSchema(name)
	return nil
}

func (p *PluginRepository) LoadSchema(name string) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}

	// Load validation schema
	schema, err := plugin.GetValidationSchema()
	if err != nil {
		return fmt.Errorf("failed to load plugin schema %s: %w", name, err)
	}
	if schema != nil {
		p.Schemas[name] = *schema
	}
	return nil
}

func (p *PluginRepository) Get(name string) (schema.MachComposerPlugin, error) {
	if name == "" {
		panic("plugin name is empty") // this should never happen
	}
	plugin, ok := p.Plugins[name]
	if !ok {
		return nil, &PluginNotFoundError{Plugin: name}
	}
	return plugin, nil
}

func (p *PluginRepository) GetSchema(name string) (*schema.ValidationSchema, error) {
	if name == "" {
		panic("plugin name is empty") // this should never happen
	}

	// this should not happen in a regular use-case
	schema, ok := p.Schemas[name]
	if !ok {
		return nil, fmt.Errorf("No schema found for %s (internal error)", name)
	}
	return &schema, nil
}

func (p *PluginRepository) All() map[string]schema.MachComposerPlugin {
	return p.Plugins
}

func (p *PluginRepository) Enabled() []Plugin {
	plugins := []Plugin{}
	for _, name := range pie.Sort(pie.Keys(p.Plugins)) {
		plugins = append(plugins, Plugin{Name: name, MachComposerPlugin: p.Plugins[name]})
	}
	return pie.Filter(plugins, func(p Plugin) bool { return p.IsEnabled() })
}

func (p *PluginRepository) SetRemoteState(name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleErr(name, plugin.SetRemoteStateBackend(data))
}

func (p *PluginRepository) SetGlobalConfig(name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleErr(name, plugin.SetGlobalConfig(data))
}

func (p *PluginRepository) SetSiteConfig(name string, site string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleErr(name, plugin.SetSiteConfig(site, data))
}

func (p *PluginRepository) SetSiteComponentConfig(site, component, name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleErr(name, plugin.SetSiteComponentConfig(site, component, data))
}

func (p *PluginRepository) SetComponentConfig(name, component string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return p.handleErr(name, plugin.SetComponentConfig(component, data))
}

func (p *PluginRepository) handleErr(plugin string, err error) error {
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
