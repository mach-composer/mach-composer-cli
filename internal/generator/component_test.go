package generator

import (
	"context"
	"strings"
	"testing"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/mach-composer/mach-composer-cli/internal/plugins"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockRenderer struct {
	mock.Mock
}

func (m *mockRenderer) Identifier() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRenderer) StateKey() string {
	args := m.Called()
	return args.String(0)
}

func (m *mockRenderer) Backend() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *mockRenderer) RemoteState() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func TestRenderRemoteSourcesCompactSources(t *testing.T) {
	rr := new(mockRenderer)
	rr.On("Identifier").Return("site/component").Times(2)
	rr.On("StateKey").Return("component")
	rr.On("RemoteState").Return("remote-state", nil)
	rr.On("Backend").Return("backend", nil)

	rr2 := new(mockRenderer)
	rr2.On("Identifier").Return("site/component-2")
	rr2.On("StateKey").Return("component-2")
	rr2.On("RemoteState").Return("remote-state-2", nil)
	rr.On("Backend").Return("backend-2", nil)

	r := state.NewRepository()
	_ = r.Add(rr)
	_ = r.Add(rr2)
	dup, _ := variable.NewScalarVariable("${component.component.value}")
	dup2, _ := variable.NewScalarVariable("${component.component-2.value}")

	sdup := variable.NewSliceVariable([]variable.Variable{dup2, dup})

	scc := config.SiteComponentConfig{
		Variables: map[string]variable.Variable{
			"value":   dup,
			"value_2": sdup,
		},
	}

	resp, err := renderRemoteSources(r, "site", scc)
	assert.NoError(t, err)
	assert.Equal(t, "remote-state\nremote-state-2", resp)
}

func TestRenderComponentModuleWithCount(t *testing.T) {
	ctx := context.Background()

	// Create mock state renderer
	rr := new(mockRenderer)
	rr.On("Identifier").Return("test-site/test-component")
	rr.On("StateKey").Return("test-component")
	rr.On("Backend").Return("backend {}", nil)

	// Create state repository
	r := state.NewRepository()
	_ = r.Add(rr)

	// Create plugin repository
	pr := plugins.NewPluginRepository()

	// Create test configuration
	cfg := &config.MachConfig{
		Global: config.GlobalConfig{
			Environment: "test",
		},
		Plugins:         pr,
		StateRepository: r,
	}

	// Create component definition
	source := config.Source("./modules/test-component")
	componentDef := &config.ComponentConfig{
		Name:    "test-component",
		Version: "1.0.0",
		Source:  source,
	}

	// Create site component with count
	siteComponent := config.SiteComponentConfig{
		Name:       "test-component",
		Definition: componentDef,
		Count:      "data.sops_external.variables.data[\"region\"] == \"foo\" ? 1 : 0",
		Variables:  make(variable.VariablesMap),
		Secrets:    make(variable.VariablesMap),
		Deployment: &config.Deployment{
			Type: config.DeploymentSiteComponent,
		},
	}

	// Create site config
	siteConfig := config.SiteConfig{
		Identifier: "test-site",
		Variables:  make(variable.VariablesMap),
		Secrets:    make(variable.VariablesMap),
	}

	// Create project config
	projectConfig := config.MachConfig{
		Global: config.GlobalConfig{
			Variables: make(variable.VariablesMap),
			Secrets:   make(variable.VariablesMap),
		},
	}

	// Create graph node
	node := &graph.SiteComponent{
		SiteComponentConfig: siteComponent,
		SiteConfig:          siteConfig,
		ProjectConfig:       projectConfig,
	}

	// Render the component module
	result, err := renderComponentModule(ctx, cfg, node)
	assert.NoError(t, err)

	// Check that count is included in the output
	assert.Contains(t, result, "count = data.sops_external.variables.data[\"region\"] == \"foo\" ? 1 : 0")
	assert.Contains(t, result, "module.test-component[*]")
}

func TestRenderComponentModuleWithoutCount(t *testing.T) {
	ctx := context.Background()

	// Create mock state renderer
	rr := new(mockRenderer)
	rr.On("Identifier").Return("test-site/test-component")
	rr.On("StateKey").Return("test-component")
	rr.On("Backend").Return("backend {}", nil)

	// Create state repository
	r := state.NewRepository()
	_ = r.Add(rr)

	// Create plugin repository
	pr := plugins.NewPluginRepository()

	// Create test configuration
	cfg := &config.MachConfig{
		Global: config.GlobalConfig{
			Environment: "test",
		},
		Plugins:         pr,
		StateRepository: r,
	}

	// Create component definition
	source := config.Source("./modules/test-component")
	componentDef := &config.ComponentConfig{
		Name:    "test-component",
		Version: "1.0.0",
		Source:  source,
	}

	// Create site component without count
	siteComponent := config.SiteComponentConfig{
		Name:       "test-component",
		Definition: componentDef,
		Count:      "", // No count specified
		Variables:  make(variable.VariablesMap),
		Secrets:    make(variable.VariablesMap),
		Deployment: &config.Deployment{
			Type: config.DeploymentSiteComponent,
		},
	}

	// Create site config
	siteConfig := config.SiteConfig{
		Identifier: "test-site",
		Variables:  make(variable.VariablesMap),
		Secrets:    make(variable.VariablesMap),
	}

	// Create project config
	projectConfig := config.MachConfig{
		Global: config.GlobalConfig{
			Variables: make(variable.VariablesMap),
			Secrets:   make(variable.VariablesMap),
		},
	}

	// Create graph node
	node := &graph.SiteComponent{
		SiteComponentConfig: siteComponent,
		SiteConfig:          siteConfig,
		ProjectConfig:       projectConfig,
	}

	// Render the component module
	result, err := renderComponentModule(ctx, cfg, node)
	assert.NoError(t, err)

	// Check that count is NOT included in the output
	assert.NotContains(t, result, "count =")
	// Check that output doesn't have [*] suffix
	lines := strings.Split(result, "\n")
	for _, line := range lines {
		if strings.Contains(line, "value =") && strings.Contains(line, "module.test-component") {
			assert.NotContains(t, line, "[*]")
		}
	}
}
