package state

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAzureRendererBackendBase(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			Environment:    "public",
			Endpoint:       "https://management.azure.com",
			MetadataHost:   "http://metadata",
			Snapshot:       true,
			StateFolder:    "state_folder",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "state_folder/test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  environment = "public"
	  endpoint = "https://management.azure.com"
	  metadata_host = "http://metadata"
	  snapshot = true
	}
	`, b)
}

func TestAzureRendererBackendMsi(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			ResourceGroup:  "resource_group",
			MsiEndpoint:    "http://msi",
			SubscriptionId: "subscription_id",
			TenantId:       "tenant_id",
			UseMsi:         true,
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  resource_group_name = "resource_group"
	  msi_endpoint = "http://msi"
	  subscription_id = "subscription_id"
	  tenant_id = "tenant_id"
	  use_msi = true
	}
	`, b)
}

func TestAzureRendererBackendOIDC(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount:    "storage_account",
			ContainerName:     "container_name",
			OidcRequestUrl:    "http://oidc",
			OidcRequestToken:  "request-token",
			OidcToken:         "token",
			OidcTokenFilePath: "token-file-path",
			UseOidc:           true,
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  oidc_request_url = "http://oidc"
	  oidc_request_token = "request-token"
	  oidc_token = "token"
	  oidc_token_file_path = "token-file-path"
	  use_oidc = true
	}
	`, b)
}

func TestAzureRendererBackendSas(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			SasToken:       "sas-token",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  sas_token = "sas-token"
	}
	`, b)
}

func TestAzureRendererBackendAccessKey(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			AccessKey:      "access-key",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  access_key = "access-key"
	}
	`, b)
}

func TestAzureRendererBackendAzureAD(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			UseAzureAdAuth: true,
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  use_azuread_auth = true
	}
	`, b)
}

func TestAzureRendererBackendServicePrincipalWithClientCertificate(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount:            "storage_account",
			ContainerName:             "container_name",
			ResourceGroup:             "resource_group",
			ClientId:                  "client_id",
			ClientCertificatePassword: "client_certificate_password",
			ClientCertificatePath:     "client_certificate_path",
			SubscriptionId:            "subscription_id",
			TenantId:                  "tenant_id",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  resource_group_name = "resource_group"
	  subscription_id = "subscription_id"
	  tenant_id = "tenant_id"
	  client_id = "client_id"
	  client_certificate_password = "client_certificate_password"
	  client_certificate_path = "client_certificate_path"
	}
	`, b)
}

func TestAzureRendererBackendServicePrincipalWithClientSecret(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			ResourceGroup:  "resource_group",
			ClientId:       "client_id",
			ClientSecret:   "client_secret",
			SubscriptionId: "subscription_id",
			TenantId:       "tenant_id",
		},
		BaseRenderer: BaseRenderer{
			identifier: "test-1/component-1",
			stateKey:   "component-1",
		},
	}

	b, err := r.Backend()
	assert.NoError(t, err)
	assert.Equal(t, `
	backend "azurerm" {
	  key = "test-1/component-1"
	  storage_account_name = "storage_account"
	  container_name = "container_name"
	  resource_group_name = "resource_group"
	  subscription_id = "subscription_id"
	  tenant_id = "tenant_id"
	  client_id = "client_id"
	  client_secret = "client_secret"
	}
	`, b)
}

func TestAzureRendererRemoteStateBase(t *testing.T) {
	r := AzureRenderer{
		state: &AzureState{
			StorageAccount: "storage_account",
			ContainerName:  "container_name",
			Environment:    "public",
			Endpoint:       "https://management.azure.com",
			MetadataHost:   "http://metadata",
			Snapshot:       true,
			StateFolder:    "state_folder",
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
	  backend = "azurerm"
	
	  config = {
		  key = "state_folder/test-1/component-1"
		  storage_account_name = "storage_account"
		  container_name = "container_name"
		  environment = "public"
		  endpoint = "https://management.azure.com"
		  metadata_host = "http://metadata"
		  snapshot = true
	  }
	}
	`, rs)
}
