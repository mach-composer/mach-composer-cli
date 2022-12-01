package mcsdk

type PluginRPCServer struct {
	Impl MachComposerPlugin
}

func (s *PluginRPCServer) Identifier(args any, resp *string) error {
	*resp = s.Impl.Identifier()
	return nil
}

func (s *PluginRPCServer) IsEnabled(args any, resp *bool) error {
	*resp = s.Impl.IsEnabled()
	return nil
}

func (s *PluginRPCServer) Configure(args ConfigureInput, resp *ErrorOutput) error {
	err := s.Impl.Configure(args.Environment, args.Provider)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetRemoteStateBackend(args SetRemoteStateBackendInput, resp *ErrorOutput) error {
	err := s.Impl.SetRemoteStateBackend(args.Data)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetGlobalConfig(args SetGlobalConfigInput, resp *ErrorOutput) error {
	err := s.Impl.SetGlobalConfig(args.Data)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetSiteConfig(args SetSiteConfigInput, resp *ErrorOutput) error {
	err := s.Impl.SetSiteConfig(args.Name, args.Data)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetSiteComponentConfig(args SetSiteComponentConfigInput, resp *SetSiteComponentConfigOutput) error {
	err := s.Impl.SetSiteComponentConfig(args.Site, args.Component, args.Data)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetSiteEndpointsConfig(args SetSiteComponentConfigInput, resp *SetSiteComponentConfigOutput) error {
	err := s.Impl.SetSiteEndpointsConfig(args.Site, args.Data)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetComponentConfig(args SetSiteComponentConfigInput, resp *SetSiteComponentConfigOutput) error {
	err := s.Impl.SetComponentConfig(args.Component, args.Data)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) SetComponentEndpointsConfig(args SetComponentEndpointsConfigInput, resp *SetSiteComponentConfigOutput) error {
	err := s.Impl.SetComponentEndpointsConfig(args.Component, args.Endpoints)
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) RenderTerraformStateBackend(
	args RenderTerraformStateBackendInput,
	resp *RenderTerraformStateBackendOutput) error {
	result, err := s.Impl.RenderTerraformStateBackend(args.Site)
	resp.Result = result
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) RenderTerraformProviders(
	args RenderTerraformProvidersInput,
	resp *RenderTerraformProvidersOutput) error {
	result, err := s.Impl.RenderTerraformProviders(args.Site)
	resp.Result = result
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) RenderTerraformResources(
	args RenderTerraformResourcesInput,
	resp *RenderTerraformResourcesOutput) error {
	result, err := s.Impl.RenderTerraformResources(args.Site)
	resp.Result = result
	resp.Err = err
	return nil
}

func (s *PluginRPCServer) RenderTerraformComponent(
	args RenderTerraformComponentInput,
	resp *RenderTerraformComponentOutput) error {
	result, err := s.Impl.RenderTerraformComponent(args.Site, args.Component)
	resp.Result = *result
	resp.Err = err
	return nil
}
