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

func (a AzureState) Key(site string) string {
	if a.StateFolder == "" {
		return site
	}
	return fmt.Sprintf("%s/%s", a.StateFolder, site)
}

type AzureRenderer struct {
	key   string
	state *AzureState
}

func (ar *AzureRenderer) Backend() (string, error) {
	templateContext := struct {
		State *AzureState
		Site  string
		Key   string
	}{
		State: ar.state,
		Site:  ar.key,
		Key:   ar.state.Key(ar.key),
	}

	tpl := `
	backend "azurerm" {
	  resource_group_name  = "{{ .State.ResourceGroup }}"
	  storage_account_name = "{{ .State.StorageAccount }}"
	  container_name       = "{{ .State.ContainerName }}"
	  Key                  = "{{ .Key }}"
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (ar *AzureRenderer) Key() string {
	return ar.key
}

func (ar *AzureRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State *AzureState
		Site  string
		Key   string
	}{
		State: ar.state,
		Site:  ar.key,
		Key:   ar.state.Key(ar.key),
	}

	template := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "azurerm"
	
	  config = {
		  resource_group_name  = "{{ .State.ResourceGroup }}"
		  storage_account_name = "{{ .State.StorageAccount }}"
		  container_name       = "{{ .State.ContainerName }}"
		  Key                  = "{{ .Key }}"
	  }
	}
	`
	return utils.RenderGoTemplate(template, templateContext)
}
