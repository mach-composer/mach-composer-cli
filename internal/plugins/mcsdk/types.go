package mcsdk

type ConfigureInput struct {
	Environment string
	Provider    string
}

type SetRemoteStateBackendInput struct {
	Data map[string]any
}

type SetGlobalConfigInput struct {
	Data map[string]any
}

type SetSiteComponentConfigInput struct {
	Site      string
	Component string
	Data      map[string]any
}

type SetSiteComponentConfigOutput struct {
	Err error
}

type SetSiteEndpointsConfigInput struct {
	Site string
	Data map[string]any
}

type SetSiteEndpointsConfigOutput struct {
	Err error
}

type SetComponentConfigInput struct {
	Component string
	Data      map[string]any
}

type SetComponentConfigOutput struct {
	Err error
}

type SetComponentEndpointsConfigInput struct {
	Component string
	Endpoints map[string]string
}

type SetComponentEndpointsConfigOutput struct {
	Err error
}

type SetSiteConfigInput struct {
	Name string
	Data map[string]any
}

type RenderTerraformStateBackendInput struct {
	Site string
}

type RenderTerraformStateBackendOutput struct {
	Result string
	Err    error
}

type RenderTerraformProvidersInput struct {
	Site string
}

type RenderTerraformProvidersOutput struct {
	Result string
	Err    error
}

type RenderTerraformResourcesInput struct {
	Site string
}

type RenderTerraformResourcesOutput struct {
	Result string
	Err    error
}

type RenderTerraformComponentInput struct {
	Site      string
	Component string
}

type RenderTerraformComponentOutput struct {
	Result ComponentSnippets
	Err    error
}

type ErrorOutput struct {
	Err error
}
