package graph

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/zclconf/go-cty/cty"
	"testing"
)

func TestHasMissingParentOutputsSiteType(t *testing.T) {
	n := new(NodeMock)
	n.On("Type").Return(SiteType).Once()

	missing, err := HasMissingParentOutputs(n)
	assert.NoError(t, err)
	assert.False(t, missing)
}

func TestHasMissingParentOutputsNoParents(t *testing.T) {
	n := new(NodeMock)
	n.On("Type").Return(SiteComponentType).Once()
	n.On("Parents").Return([]Node{}, fmt.Errorf("error")).Once()

	_, err := HasMissingParentOutputs(n)
	assert.Error(t, err)
}

func TestHasMissingParentOutputsParentWithNilVal(t *testing.T) {
	parent := new(NodeMock)
	parent.On("Outputs").Return(cty.NilVal).Once()

	n := new(NodeMock)
	n.On("Type").Return(SiteComponentType).Once()
	n.On("Parents").Return([]Node{parent}, nil).Once()

	missing, err := HasMissingParentOutputs(n)
	assert.NoError(t, err)
	assert.True(t, missing)
}

func TestHasMissingParentOutputsSingle_False(t *testing.T) {
	parent := new(NodeMock)
	parent.On("Outputs").Return(cty.StringVal("value")).Once()

	n := new(NodeMock)
	n.On("Type").Return(SiteComponentType).Once()
	n.On("Parents").Return([]Node{parent}, nil).Once()

	missing, err := HasMissingParentOutputs(n)
	assert.NoError(t, err)
	assert.False(t, missing)
}

func TestHasMissingParentOutputsMultiple_False(t *testing.T) {
	parent1 := new(NodeMock)
	parent1.On("Outputs").Return(cty.StringVal("value")).Once()
	parent2 := new(NodeMock)
	parent2.On("Outputs").Return(cty.StringVal("value")).Once()

	n := new(NodeMock)
	n.On("Type").Return(SiteComponentType).Once()
	n.On("Parents").Return([]Node{parent1, parent2}, nil).Once()

	missing, err := HasMissingParentOutputs(n)
	assert.NoError(t, err)
	assert.False(t, missing)
}

func TestHasMissingParentOutputsMultiple_True(t *testing.T) {
	parent1 := new(NodeMock)
	parent1.On("Outputs").Return(cty.StringVal("value")).Once()
	parent2 := new(NodeMock)
	parent2.On("Outputs").Return(cty.NilVal).Once()

	n := new(NodeMock)
	n.On("Type").Return(SiteComponentType).Once()
	n.On("Parents").Return([]Node{parent1, parent2}, nil).Once()

	missing, err := HasMissingParentOutputs(n)
	assert.NoError(t, err)
	assert.True(t, missing)
}
