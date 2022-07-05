package config

type AmplienceConfig struct {
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	HubID        string `yaml:"hub_id"`
}
