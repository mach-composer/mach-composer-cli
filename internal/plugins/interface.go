package plugins

type MachComposerPlugin interface {
	Identifier() string

	IsEnabled() bool

	SetRemoteStateBackend(data map[string]any) error

	SetGlobalConfig(data map[string]any) error

	SetSiteConfig(site string, data map[string]any) error

	SetSiteComponentConfig(site string, component string, data map[string]any) error

	SetSiteEndpointsConfig(site string, data map[string]any) error

	SetComponentConfig(component string, data map[string]any) error

	SetComponentEndpointsConfig(component string, data map[string]string) error

	// Render remote state configuration
	TerraformRenderStateBackend(site string) string

	// Render all required terraform providers
	TerraformRenderProviders(site string) string

	// Render global resources
	TerraformRenderResources(site string) string

	// Render resources required per module
	TerraformRenderComponentResources(site string, component string) string

	// Render variables to pass to a component
	TerraformRenderComponentVars(site string, component string) string

	// Render depends_on clause for a component
	TerraformRenderComponentDependsOn(site string, component string) []string

	TerraformRenderComponentProviders(site string, component string) []string
}

type MachComposerPluginCloud interface {
	MachComposerPlugin
}
