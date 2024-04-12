package state

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type AzureState struct {
	ResourceGroup  string `mapstructure:"resource_group"`
	StorageAccount string `mapstructure:"storage_account"`
	ContainerName  string `mapstructure:"container_name"`
	StateFolder    string `mapstructure:"state_folder"`
}

func (a AzureState) Identifier(identifier string) string {
	if a.StateFolder == "" {
		return identifier
	}
	return fmt.Sprintf("%s/%s", a.StateFolder, identifier)
}

type AzureRenderer struct {
	BaseRenderer
	state *AzureState
}

func (ar *AzureRenderer) Backend() (string, error) {
	templateContext := struct {
		State      *AzureState
		Site       string
		Identifier string
	}{
		State:      ar.state,
		Site:       ar.key,
		Identifier: ar.state.Identifier(ar.identifier),
	}

	tpl := `
	backend "azurerm" {
	  resource_group_name  = "{{ .State.ResourceGroup }}"
	  storage_account_name = "{{ .State.StorageAccount }}"
	  container_name       = "{{ .State.ContainerName }}"
	  key                  = "{{ .Identifier }}"
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (ar *AzureRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State      *AzureState
		Site       string
		Identifier string
		Key        string
	}{
		State:      ar.state,
		Site:       ar.key,
		Identifier: ar.state.Identifier(ar.identifier),
		Key:        ar.key,
	}

	template := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "azurerm"
	
	  config = {
		  resource_group_name  = "{{ .State.ResourceGroup }}"
		  storage_account_name = "{{ .State.StorageAccount }}"
		  container_name       = "{{ .State.ContainerName }}"
		  key                  = "{{ .Identifier }}"
	  }
	}
	`
	return utils.RenderGoTemplate(template, templateContext)
}
