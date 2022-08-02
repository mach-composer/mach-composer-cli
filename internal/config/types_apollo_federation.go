package config

type ApolloFederationConfig struct {
	APIKey       string `yaml:"api_key"`
	Graph        string `yaml:"graph"`
	GraphVariant string `yaml:"graph_variant"`
}
