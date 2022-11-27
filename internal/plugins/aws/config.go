package aws

import (
	"log"
	"net/url"
	"strings"
)

type SiteConfig struct {
	AccountID string `mapstructure:"account_id"`
	Region    string `mapstructure:"region"`

	DeployRoleName string            `mapstructure:"deploy_role_name"`
	ExtraProviders []AWSProvider     `mapstructure:"extra_providers"`
	DefaultTags    map[string]string `mapstructure:"default_tags"`
}

type ComponentConfig struct {
	Endpoints map[string]string
}

// AWSTFState AWS S3 bucket state backend configuration.
type AWSTFState struct {
	Bucket    string `mapstructure:"bucket"`
	KeyPrefix string `mapstructure:"key_prefix"`
	Region    string `mapstructure:"region"`
	RoleARN   string `mapstructure:"role_arn"`
	LockTable string `mapstructure:"lock_table"`
	Encrypt   bool   `mapstructure:"encrypt" default:"true"`
}

type AWSProvider struct {
	Name        string
	Region      string
	DefaultTags map[string]string `mapstructure:"default_tags"`
}

type EndpointConfig struct {
	// Key is the id of the endpoint (.e.g internal)
	Key string `mapstructure:"key"`

	URL                  string `mapstructure:"url"`
	ThrottlingBurstLimit *int   `mapstructure:"throttling_burst_limit"`
	ThrottlingRateLimit  *int   `mapstructure:"throttling_rate_limit"`
	EnableCDN            bool   `mapstructure:"enable_cdn"`

	AssociatedComponents []string `mapstructure:"-"`

	// Zone is the route53 DNS Zone
	Zone string `mapstructure:"zone"`
}

func (e *EndpointConfig) SetDefaults() {
	e.URL = StripProtocol(e.URL)

	if e.Zone == "" && e.URL != "" {
		e.Zone = ZoneFromURL(e.URL)
	}
}

func ZoneFromURL(value string) string {
	u, err := url.Parse(value)
	if err != nil {
		log.Fatal(err)
	}

	var domains []string
	if !strings.Contains(value, "://") {
		parts := strings.SplitN(value, "/", 2)
		domains = strings.Split(parts[0], ".")
	} else {
		domains = strings.Split(u.Hostname(), ".")
	}

	if len(domains) < 3 {
		return strings.Join(domains, ".")
	} else {
		return strings.Join(domains[1:], ".")
	}
}
