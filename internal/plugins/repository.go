package plugins

import (
	"fmt"

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

func (e *PluginNotFoundError) Error() string {
	return fmt.Sprintf("plugin %s not found", e.Plugin)
}

type PluginRepository struct {
	Plugins map[string]schema.MachComposerPlugin
}

func NewPluginRepository() *PluginRepository {
	return &PluginRepository{
		Plugins: make(map[string]schema.MachComposerPlugin),
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
	return plugin.SetRemoteStateBackend(data)
}

func (p *PluginRepository) SetGlobalConfig(name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return plugin.SetGlobalConfig(data)
}

func (p *PluginRepository) SetSiteConfig(name string, site string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return plugin.SetSiteConfig(site, data)
}

func (p *PluginRepository) SetSiteComponentConfig(site, component, name string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return plugin.SetSiteComponentConfig(site, component, data)
}

func (p *PluginRepository) SetComponentConfig(name, component string, data map[string]any) error {
	plugin, err := p.Get(name)
	if err != nil {
		return err
	}
	return plugin.SetComponentConfig(component, data)
}
