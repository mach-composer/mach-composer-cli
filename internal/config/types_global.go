package config

const (
	Azure = "azure"
	AWS   = "aws"
)

type Global struct {
	Environment string `yaml:"environment"`
	Cloud       string `yaml:"cloud"`

	Azure           *GlobalAzureConfig  `yaml:"azure"`
	TerraformConfig TerraformConfig     `yaml:"terraform_config"`
	AmplienceConfig *AmplienceConfig    `yaml:"amplience"`
	SentryConfig    *GlobalSentryConfig `yaml:"sentry"`
}

type TerraformConfig struct {
	AzureRemoteState *AzureTFState     `yaml:"azure_remote_state"`
	AwsRemoteState   *AWSTFState       `yaml:"aws_remote_state"`
	Providers        TerraformProvider `yaml:"providers"`
}

// TerraformProvider version overwrites.
type TerraformProvider struct {
	AWS           string `yaml:"aws"`
	Azure         string `yaml:"azure"`
	Commercetools string `yaml:"commercetools"`
	Sentry        string `yaml:"sentry"`
	Contentful    string `yaml:"contentful"`
	Amplience     string `yaml:"amplience"`
}
