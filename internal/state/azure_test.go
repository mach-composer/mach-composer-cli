package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAzureRendererBackendFull(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			ResourceGroup:  "resource_group",
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			StateFolder:    "state_folder",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			key:        "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  resource_group_name  = "resource_group"
	  storage_account_name = "storage_account"
	  container_name       = "container_name"
	  key                  = "state_folder/test-1/component-1"
	}
	`, b)
}

func TestAzureRendererRemoteStateFull(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			ResourceGroup:  "resource_group",
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			StateFolder:    "state_folder",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			key:        "component-1",
		},
	}

	rs, err := r.RemoteState()
	assert.NoError(t, err)
	assert.Equal(t, `
	data "terraform_remote_state" "component-1" {
	  backend = "azurerm"
	
	  config = {
		  resource_group_name  = "resource_group"
		  storage_account_name = "storage_account"
		  container_name       = "container_name"
		  key                  = "state_folder/test-1/component-1"
	  }
	}
	`, rs)
}
