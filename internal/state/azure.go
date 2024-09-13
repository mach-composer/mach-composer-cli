package state

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type AzureState struct {
	StorageAccount            string `mapstructure:"storage_account"`
	ContainerName             string `mapstructure:"container_name"`
	Environment               string `mapstructure:"environment"`
	Endpoint                  string `mapstructure:"endpoint"`
	MetadataHost              string `mapstructure:"metadata_host"`
	Snapshot                  bool   `mapstructure:"snapshot"`
	ResourceGroup             string `mapstructure:"resource_group"`
	MsiEndpoint               string `mapstructure:"msi_endpoint"`
	SubscriptionId            string `mapstructure:"subscription_id"`
	TenantId                  string `mapstructure:"tenant_id"`
	UseMsi                    bool   `mapstructure:"use_msi"`
	OidcRequestUrl            string `mapstructure:"oidc_request_url"`
	OidcRequestToken          string `mapstructure:"oidc_request_token"`
	OidcToken                 string `mapstructure:"oidc_token"`
	OidcTokenFilePath         string `mapstructure:"oidc_token_file_path"`
	UseOidc                   bool   `mapstructure:"use_oidc"`
	SasToken                  string `mapstructure:"sas_token"`
	AccessKey                 string `mapstructure:"access_key"`
	UseAzureAdAuth            bool   `mapstructure:"use_azuread_auth"`
	ClientId                  string `mapstructure:"client_id"`
	ClientSecret              string `mapstructure:"client_secret"`
	ClientCertificatePassword string `mapstructure:"client_certificate_password"`
	ClientCertificatePath     string `mapstructure:"client_certificate_path"`
	StateFolder               string `mapstructure:"state_folder"`
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
		Identifier string
	}{
		State:      ar.state,
		Identifier: ar.state.Identifier(ar.identifier),
	}

	tpl := `
	backend "azurerm" {
	  key = "{{ .Identifier }}"
	  storage_account_name = "{{ .State.StorageAccount }}"
	  container_name = "{{ .State.ContainerName }}"
	  {{- if .State.Environment }}
	  environment = "{{ .State.Environment }}"
	  {{- end }}
	  {{- if .State.Endpoint }}
	  endpoint = "{{ .State.Endpoint }}"
	  {{- end }}
	  {{- if .State.MetadataHost }}
	  metadata_host = "{{ .State.MetadataHost }}"
	  {{- end }}
	  {{- if .State.Snapshot }}
	  snapshot = true
	  {{- end }}
	  {{- if .State.ResourceGroup }}
	  resource_group_name = "{{ .State.ResourceGroup }}"
	  {{- end }}
	  {{- if .State.MsiEndpoint }}
	  msi_endpoint = "{{ .State.MsiEndpoint }}"
	  {{- end }}
	  {{- if .State.SubscriptionId }}
	  subscription_id = "{{ .State.SubscriptionId }}"
	  {{- end }}
	  {{- if .State.TenantId }}
	  tenant_id = "{{ .State.TenantId }}"
	  {{- end }}
	  {{- if .State.UseMsi }}
	  use_msi = true
	  {{- end }}
	  {{- if .State.OidcRequestUrl }}
	  oidc_request_url = "{{ .State.OidcRequestUrl }}"
	  {{- end }}
	  {{- if .State.OidcRequestToken }}
	  oidc_request_token = "{{ .State.OidcRequestToken }}"
	  {{- end }}
	  {{- if .State.OidcToken }}
	  oidc_token = "{{ .State.OidcToken }}"
	  {{- end }}
	  {{- if .State.OidcTokenFilePath }}
	  oidc_token_file_path = "{{ .State.OidcTokenFilePath }}"
	  {{- end }}
	  {{- if .State.UseOidc }}
	  use_oidc = true
	  {{- end }}
	  {{- if .State.SasToken }}
	  sas_token = "{{ .State.SasToken }}"
	  {{- end }}
	  {{- if .State.AccessKey }}
	  access_key = "{{ .State.AccessKey }}"
	  {{- end }}
	  {{- if .State.UseAzureAdAuth }}
	  use_azuread_auth = true
	  {{- end }}
	  {{- if .State.ClientId }}
	  client_id = "{{ .State.ClientId }}"
	  {{- end }}
	  {{- if .State.ClientSecret }}
	  client_secret = "{{ .State.ClientSecret }}"
	  {{- end }}
	  {{- if .State.ClientCertificatePassword }}
	  client_certificate_password = "{{ .State.ClientCertificatePassword }}"
	  {{- end }}
	  {{- if .State.ClientCertificatePath }}
	  client_certificate_path = "{{ .State.ClientCertificatePath }}"
	  {{- end }}
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (ar *AzureRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State      *AzureState
		Identifier string
		Key        string
	}{
		State:      ar.state,
		Identifier: ar.state.Identifier(ar.identifier),
		Key:        ar.stateKey,
	}

	template := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "azurerm"
	
	  config = {
		  key = "{{ .Identifier }}"
		  storage_account_name = "{{ .State.StorageAccount }}"
		  container_name = "{{ .State.ContainerName }}"
		  {{- if .State.Environment }}
		  environment = "{{ .State.Environment }}"
		  {{- end }}
		  {{- if .State.Endpoint }}
		  endpoint = "{{ .State.Endpoint }}"
		  {{- end }}
		  {{- if .State.MetadataHost }}
		  metadata_host = "{{ .State.MetadataHost }}"
		  {{- end }}
		  {{- if .State.Snapshot }}
		  snapshot = true
		  {{- end }}
		  {{- if .State.ResourceGroup }}
		  resource_group_name = "{{ .State.ResourceGroup }}"
		  {{- end }}
		  {{- if .State.MsiEndpoint }}
		  msi_endpoint = "{{ .State.MsiEndpoint }}"
		  {{- end }}
		  {{- if .State.SubscriptionId }}
		  subscription_id = "{{ .State.SubscriptionId }}"
		  {{- end }}
		  {{- if .State.TenantId }}
		  tenant_id = "{{ .State.TenantId }}"
		  {{- end }}
		  {{- if .State.UseMsi }}
		  use_msi = true
		  {{- end }}
		  {{- if .State.OidcRequestUrl }}
		  oidc_request_url = "{{ .State.OidcRequestUrl }}"
		  {{- end }}
		  {{- if .State.OidcRequestToken }}
		  oidc_request_token = "{{ .State.OidcRequestToken }}"
		  {{- end }}
		  {{- if .State.OidcToken }}
		  oidc_token = "{{ .State.OidcToken }}"
		  {{- end }}
		  {{- if .State.OidcTokenFilePath }}
		  oidc_token_file_path = "{{ .State.OidcTokenFilePath }}"
		  {{- end }}
		  {{- if .State.UseOidc }}
		  use_oidc = true
		  {{- end }}
		  {{- if .State.SasToken }}
		  sas_token = "{{ .State.SasToken }}"
		  {{- end }}
		  {{- if .State.AccessKey }}
		  access_key = "{{ .State.AccessKey }}"
		  {{- end }}
		  {{- if .State.UseAzureAdAuth }}
		  use_azuread_auth = true
		  {{- end }}
		  {{- if .State.ClientId }}
		  client_id = "{{ .State.ClientId }}"
		  {{- end }}
		  {{- if .State.ClientSecret }}
		  client_secret = "{{ .State.ClientSecret }}"
		  {{- end }}
		  {{- if .State.ClientCertificatePassword }}
		  client_certificate_password = "{{ .State.ClientCertificatePassword }}"
		  {{- end }}
		  {{- if .State.ClientCertificatePath }}
		  client_certificate_path = "{{ .State.ClientCertificatePath }}"
		  {{- end }}
	  }
	}
	`
	return utils.RenderGoTemplate(template, templateContext)
}
