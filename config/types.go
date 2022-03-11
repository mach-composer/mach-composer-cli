package config

type Root struct {
	Global     Global      `yaml:"global"`
	Sites      []Site      `yaml:"sites"`
	Components []Component `yaml:"components"`
}

type Global struct {
	Environment string `yaml:"environment"`
	Cloud       string `yaml:"cloud"`
}

// Site contains all configuration needed for a site.
type Site struct {
	Name       string
	Identifier string

	CommercetoolsSettings CommercetoolsSettings `yaml:"commercetools"`
	Components            []SiteComponent       `yaml:"components"`
}

type CommercetoolsSettings struct {
	ProjectKey   string `yaml:"project_key"`
	ClientID     string `yaml:"client_id"`
	ClientSecret string `yaml:"client_secret"`
	Scopes       string `yaml:"scopes"`
}

type CommercetoolsProjectSettings struct {
	Languages  []string `yaml:"languages"`
	Currencies []string `yaml:"currencies"`
	Countries  []string `yaml:"countries"`
}

type SiteComponent struct {
	Name      string
	Variables map[string]string
	Secrets   map[string]string
}

type Component struct {
	Name         string
	Source       string
	Version      string
	Integrations []string
}
