package plugins

type MachComposerPlugin interface {
	Identifier() string

	IsEnabled() bool

	Configure(environment string, provider string) error

	SetRemoteStateBackend(data map[string]any) error

	SetGlobalConfig(data map[string]any) error

	SetSiteConfig(site string, data map[string]any) error

	SetSiteComponentConfig(site string, component string, data map[string]any) error

	SetSiteEndpointsConfig(site string, data map[string]any) error

	SetComponentConfig(component string, data map[string]any) error

	SetComponentEndpointsConfig(component string, data map[string]string) error

	// Render remote state configuration
	TerraformRenderStateBackend(site string) (string, error)

	// Render all required terraform providers
	TerraformRenderProviders(site string) (string, error)

	// Render global resources
	TerraformRenderResources(site string) (string, error)

	// Render resources required per module
	TerraformRenderComponentResources(site string, component string) (string, error)

	// Render variables to pass to a component
	TerraformRenderComponentVars(site string, component string) (string, error)

	// Render depends_on clause for a component
	TerraformRenderComponentDependsOn(site string, component string) ([]string, error)

	TerraformRenderComponentProviders(site string, component string) ([]string, error)
}

type MachComposerPluginCloud interface {
	MachComposerPlugin
}
