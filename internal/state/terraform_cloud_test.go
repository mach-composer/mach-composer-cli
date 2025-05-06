package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTerraformCloudRendererBackendFull(t *testing.T) {
	r := TerraformCloudRenderer{
		state: &TerraformCloudState{
			Organization: "organization",
			Hostname:     "hostname",
			Token:        "token",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "remote" {
	  organization = "organization"
	  
	  hostname = "hostname"
	  
	  
	  token = "token"
	  
	  
	  workspaces {
		
		
	  }
	  
	}
	`, b)
}

func TestTerraformCloudRendererRemoteStateFull(t *testing.T) {
	r := TerraformCloudRenderer{
		state: &TerraformCloudState{
			Organization: "organization",
			Hostname:     "hostname",
			Token:        "token",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	rs, err := r.RemoteState()
	assert.NoError(t, err)
	assert.Equal(t, `
	data "terraform_remote_state" "component-1" {
	  backend = "remote"
	
	  config = {
	    organization = "organization"
	    
	    hostname = "hostname"
	    
	    
	    workspaces {
		  
		  
	    }
	    
	  }
	}
	`, rs)
}
