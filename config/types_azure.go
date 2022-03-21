package config

// Azure storage account state backend configuration.
type AzureTFState struct {
	ResourceGroup  string `yaml:"resource_group"`
	StorageAccount string `yaml:"storage_account"`
	ContainerName  string `yaml:"container_name"`
	StateFolder    string `yaml:"state_folder"`
}

type AzureEndpoint struct {
	SessionAffinityEnabled bool   `yaml:"session_affinity_enabled"`
	SessionAffinityTTL     int    `yaml:"session_affinity_ttl_seconds"`
	WAFPolicyID            string `yaml:"waf_policy_id"`
	InternalName           string `yaml:"internal_name"`
}
