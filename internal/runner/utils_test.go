package runner

import (
	"context"
	"github.com/mach-composer/mach-composer-cli/internal/graph"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func TestTerraformCanPlanNoParents(t *testing.T) {
	n := new(graph.NodeMock)
	n.On("Parents").Return([]graph.Node{}, nil).Once()

	canPlan, err := terraformCanPlan(context.Background(), n)
	assert.NoError(t, err)
	assert.True(t, canPlan)
}

func TestTerraformCanPlanWithSite(t *testing.T) {
	p := new(graph.NodeMock)
	dir, _ := os.Getwd()
	p.On("Path").Return(path.Join(dir, "testdata/empty")).Once()
	p.On("Identifier").Return("main").Once()

	n := new(graph.NodeMock)
	n.On("Parents").Return([]graph.Node{p}, nil).Once()
	n.On("Type").Return(graph.SiteType).Once()

	canPlan, err := terraformCanPlan(context.Background(), n)
	assert.NoError(t, err)
	assert.True(t, canPlan)
}

func TestTerraformCanPlanWithComponentParentEmptyOutput(t *testing.T) {
	p := new(graph.NodeMock)
	dir, _ := os.Getwd()
	p.On("Path").Return(path.Join(dir, "testdata/empty")).Once()
	p.On("Identifier").Return("main").Once()

	n := new(graph.NodeMock)
	n.On("Parents").Return([]graph.Node{p}, nil).Once()
	n.On("Type").Return(graph.SiteComponentType).Once()

	canPlan, err := terraformCanPlan(context.Background(), n)
	assert.NoError(t, err)
	assert.False(t, canPlan)
}

func TestTerraformCanPlanWithParentOutput(t *testing.T) {
	p := new(graph.NodeMock)
	dir, _ := os.Getwd()
	p.On("Path").Return(path.Join(dir, "testdata/initialized")).Once()
	p.On("Identifier").Return("main").Once()

	n := new(graph.NodeMock)
	n.On("Parents").Return([]graph.Node{p}, nil).Once()
	n.On("Type").Return(graph.SiteComponentType).Once()

	canPlan, err := terraformCanPlan(context.Background(), n)
	assert.NoError(t, err)
	assert.True(t, canPlan)
}
