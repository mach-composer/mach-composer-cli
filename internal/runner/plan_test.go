package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func resetTerraformRunner() func() {
	return func() {
		defaultRunTerraform = runTerraform
	}
}

func TestTerraformPlanWithLock(t *testing.T) {
	defer resetTerraformRunner()()

	cwd := "path/to/test"

	mockedRunner := new(TerraformRunnerMock)
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"init"}).Return(nil).Once()
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"plan", "-out=terraform.plan"}).Return(nil).Once()

	defaultRunTerraform = mockedRunner.RunTerraform

	err := TerraformPlan(context.Background(), &config.MachConfig{
		Sites: []config.SiteConfig{
			{Identifier: "my-test"},
		},
	}, map[string]string{
		"my-test": cwd,
	}, &PlanOptions{Lock: true})

	assert.NoError(t, err)
}

func TestTerraformPlanWithoutLock(t *testing.T) {
	defer resetTerraformRunner()()
	cwd := "path/to/test"

	mockedRunner := new(TerraformRunnerMock)
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"init"}).
		Return(nil).Once()
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"plan", "-lock=false", "-out=terraform.plan"}).
		Return(nil).Once()

	defaultRunTerraform = mockedRunner.RunTerraform

	err := TerraformPlan(context.Background(), &config.MachConfig{
		Sites: []config.SiteConfig{
			{Identifier: "my-test"},
		},
	}, map[string]string{
		"my-test": cwd,
	}, &PlanOptions{Lock: false})

	assert.NoError(t, err)
}

func TestTerraformPlanWithReuse(t *testing.T) {
	defer resetTerraformRunner()()
	cwd := "path/to/test"

	mockedRunner := new(TerraformRunnerMock)
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"plan", "-out=terraform.plan"}).
		Return(nil).Once()

	defaultRunTerraform = mockedRunner.RunTerraform

	err := TerraformPlan(context.Background(), &config.MachConfig{
		Sites: []config.SiteConfig{
			{Identifier: "my-test"},
		},
	}, map[string]string{
		"my-test": cwd,
	}, &PlanOptions{Lock: true, Reuse: true})

	assert.NoError(t, err)
}

func TestTerraformPlanWithComponents(t *testing.T) {
	defer resetTerraformRunner()()
	cwd := "path/to/test"
	siteIdentifier := "my-site"
	componentIdentifier := "my-component"

	mockedRunner := new(TerraformRunnerMock)
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"init"}).
		Return(nil).Once()
	mockedRunner.
		On("runTerraform", mock.Anything, cwd, []string{"plan", "-target=module.my-component", "-out=terraform.plan"}).
		Return(nil).Once()

	defaultRunTerraform = mockedRunner.RunTerraform

	err := TerraformPlan(context.Background(),
		&config.MachConfig{Sites: []config.SiteConfig{{Identifier: siteIdentifier}}},
		map[string]string{siteIdentifier: cwd},
		&PlanOptions{Lock: true, Components: []string{componentIdentifier}},
	)

	assert.NoError(t, err)
}

func TestTerraformPlanWithSingleSite(t *testing.T) {
	defer resetTerraformRunner()()
	sitePath := "path/to/test"
	siteIdentifier := "my-site"
	otherSiteIdentifier := "my-other-site"

	mockedRunner := new(TerraformRunnerMock)
	mockedRunner.
		On("runTerraform", mock.Anything, sitePath, []string{"init"}).
		Return(nil).Once()
	mockedRunner.
		On("runTerraform", mock.Anything, sitePath, []string{"plan", "-out=terraform.plan"}).
		Return(nil).Once()

	defaultRunTerraform = mockedRunner.RunTerraform

	err := TerraformPlan(context.Background(),
		&config.MachConfig{Sites: []config.SiteConfig{{Identifier: siteIdentifier}}},
		map[string]string{siteIdentifier: sitePath, otherSiteIdentifier: "path/to/other/test"},
		&PlanOptions{Lock: true, Site: siteIdentifier},
	)

	assert.NoError(t, err)
}
