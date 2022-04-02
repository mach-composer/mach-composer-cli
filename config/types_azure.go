package config

import "strings"

// Azure storage account state backend configuration.
type AzureTFState struct {
	ResourceGroup  string `yaml:"resource_group"`
	StorageAccount string `yaml:"storage_account"`
	ContainerName  string `yaml:"container_name"`
	StateFolder    string `yaml:"state_folder"`
}

type GlobalAzureConfig struct {
	TenantID       string `yaml:"tenant_id"`
	SubscriptionID string `yaml:"subscription_id"`
	Region         string

	Frontdoor        *AzureFrontdoorSettings     `yaml:"frontdoor"`
	ResourcesPrefix  string                      `yaml:"resources_prefix"`
	ServiceObjectIds map[string]string           `yaml:"service_object_ids"`
	ServicePlans     map[string]AzureServicePlan `yaml:"service_plans"`
}

// Site-specific Azure settings
type SiteAzureSettings struct {
	Frontdoor  *AzureFrontdoorSettings `yaml:"frontdoor"`
	AlertGroup *AzureAlertGroup        `yaml:"alert_group"`

	// Can overwrite values from AzureConfig
	ResourceGroup  string
	TenantID       string `yaml:"tenant_id"`
	SubscriptionID string `yaml:"subscription_id"`

	Region           string
	ServiceObjectIds map[string]string           `yaml:"service_object_ids"`
	ServicePlans     map[string]AzureServicePlan `yaml:"service_plans"`
}

type ComponentAzureConfig struct {
	ServicePlan string `yaml:"service_plan"`
	ShortName   string `yaml:"short_name"`
}

// func (c *ComponentAzureConfig) Merge() {
// 	if c.ServicePlan == "" {
// 		c.ServicePlan = config.ServicePlan
// 	}
// 	if c.ShortName == "" {
// 		c.ShortName = config.ShortName
// 	}
// }

func (a *SiteAzureSettings) Merge(c *GlobalAzureConfig) {
	if a.Frontdoor == nil {
		a.Frontdoor = c.Frontdoor
	}
	if a.TenantID == "" {
		a.TenantID = c.TenantID
	}
	if a.SubscriptionID == "" {
		a.SubscriptionID = c.SubscriptionID
	}
	if a.Region == "" {
		a.Region = c.Region
	}

	if len(a.ServiceObjectIds) == 0 {
		a.ServiceObjectIds = c.ServiceObjectIds
	}

	for k, v := range c.ServicePlans {
		a.ServicePlans[k] = v
	}
}

type AzureEndpoint struct {
	SessionAffinityEnabled bool   `yaml:"session_affinity_enabled"`
	SessionAffinityTTL     int    `yaml:"session_affinity_ttl_seconds"`
	WAFPolicyID            string `yaml:"waf_policy_id"`
	InternalName           string `yaml:"internal_name"`
}

type AzureFrontdoorSettings struct {
	DnsResourceGroup string                   `yaml:"dns_resource_group"`
	SslKeyVault      *AzureFrontdoorSslConfig `yaml:"ssl_key_vault"`

	// Undocumented option to workaround some tenacious issues
	// with using Frontdoor in the Azure Terraform provider
	SupressChanges bool `yaml:"supress_changes"`
}

type AzureFrontdoorSslConfig struct {
	Name          string
	ResourceGroup string `yaml:"resource_group"`
	SecretName    string `yaml:"secret_name"`
}

type AzureServicePlan struct {
	Kind                   string
	Tier                   string
	Size                   string
	Capacity               int
	DedicatedResourceGroup bool `yaml:"dedicated_resource_group"`
	PerSiteScaling         bool `yaml:"per_site_scaling"`
}

type AzureAlertGroup struct {
	Name        string
	AlertEmails []string `yaml:"alert_emals"`
	WebhookURL  string   `yaml:"webhook_url"`
	LogicApp    string   `yaml:"logic_app"`
}

func (a *AzureAlertGroup) LogicAppName() string {
	if a.LogicApp != "" {
		parts := strings.Split(a.LogicApp, ".")
		return parts[1]
	}
	return ""
}

func (a *AzureAlertGroup) LogicAppResourceGroup() string {
	if a.LogicApp != "" {
		parts := strings.Split(a.LogicApp, ".")
		return parts[0]
	}
	return ""
}
