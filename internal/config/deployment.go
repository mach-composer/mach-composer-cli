package config

type DeploymentType string

const (
	DeploymentSite          DeploymentType = "site"
	DeploymentSiteComponent DeploymentType = "site-component"
)

type Deployment struct {
	Type DeploymentType `yaml:"type" default:"site"`
}

type GlobalDeployment struct {
	Deployment
	Runners int `yaml:"runners" default:"1"`
}
