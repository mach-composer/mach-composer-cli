package contentful

type ContentfulConfig struct {
	Space          string `mapstructure:"space"`
	DefaultLocale  string `mapstructure:"default_locale"`
	CMAToken       string `mapstructure:"default_locale"`
	OrganizationID string `mapstructure:"organization_id"`
}

func (c *ContentfulConfig) extendConfig(o *ContentfulConfig) *ContentfulConfig {
	cfg := &ContentfulConfig{
		Space:          o.Space,
		DefaultLocale:  o.DefaultLocale,
		CMAToken:       o.CMAToken,
		OrganizationID: o.OrganizationID,
	}
	if c.Space != "" {
		cfg.Space = c.Space
	}
	if c.DefaultLocale != "" {
		cfg.DefaultLocale = c.DefaultLocale
	}
	if c.CMAToken != "" {
		cfg.CMAToken = c.CMAToken
	}
	if c.OrganizationID != "" {
		cfg.OrganizationID = c.OrganizationID
	}
	return cfg
}
