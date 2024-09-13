package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLocalRendererBackendFull(t *testing.T) {
	r := LocalRenderer{
		state: &LocalState{Path: "path"},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "local" {
		
		path = "path/test-1/component-1.tfstate"
		
	}
	`, b)
}

func TestLocalRendererRemoteStateFull(t *testing.T) {
	r := LocalRenderer{
		state: &LocalState{Path: "path"},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	rs, err := r.RemoteState()
	assert.NoError(t, err)
	assert.Equal(t, `
	data "terraform_remote_state" "component-1" {
	  backend = "local"
	
	  config = {
		
		path = "path/test-1/component-1.tfstate"
		
	  }
	}
	`, rs)
}
