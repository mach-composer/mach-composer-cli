package mcsdk

type ComponentSnippets struct {
	Resources string
	Variables string
	DependsOn []string
	Providers []string
}

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
	RenderTerraformStateBackend(site string) (string, error)

	// Render all required terraform providers
	RenderTerraformProviders(site string) (string, error)

	// Render global resources
	RenderTerraformResources(site string) (string, error)

	// Render component
	RenderTerraformComponent(site string, component string) (*ComponentSnippets, error)
}
