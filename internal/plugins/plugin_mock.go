package plugins

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
)

type MockPlugin struct {
	Environment string
	Provider    string
}

func NewMockPlugin() schema.MachComposerPlugin {
	state := MockPlugin{}

	return plugin.NewPlugin(&schema.PluginSchema{
		Configure: state.Configure,
		IsEnabled: func() bool { return true },
	})
}

func (p *MockPlugin) Configure(environment string, provider string) error {
	p.Environment = environment
	if provider != "" {
		p.Provider = provider
	}
	return nil
}
