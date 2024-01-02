package plugins

import (
	"encoding/json"

	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
)

type MockPlugin struct {
	Environment string
	Provider    string
}

func NewMockPluginV1() schema.MachComposerPlugin {
	state := MockPlugin{}

	return plugin.NewPlugin(&schema.PluginSchema{
		Configure:                state.Configure,
		GetValidationSchema:      state.GetValidationSchema,
		RenderTerraformComponent: state.RenderTerraformComponent,
		IsEnabled:                func() bool { return true },
	})
}

func (p *MockPlugin) Configure(environment string, provider string) error {
	p.Environment = environment
	if provider != "" {
		p.Provider = provider
	}
	return nil
}

func (p *MockPlugin) GetValidationSchema() (*schema.ValidationSchema, error) {
	content := []byte(`
	{
		"type": "object",
		"required": ["requiredValue"],
		"additionalProperties": false,
		"properties": {
		  "boolValue": {
			"type": "boolean"
		  },
		  "stringValue": {
			"type": "string"
		  },
		  "requiredValue": {
			"type": "string"
		  },
		  "intValue": {
			"type": "number"
		  }
		}
	  }
	`)
	result := schema.ValidationSchema{}
	if err := json.Unmarshal(content, &result.SiteConfigSchema); err != nil {
		panic(err)
	}

	return &result, nil
}

func (p *MockPlugin) RenderTerraformComponent(site string, component string) (*schema.ComponentSchema, error) {
	return &schema.ComponentSchema{}, nil
}
