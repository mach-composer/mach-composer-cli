package sentry

// SentryConfigBase is the base sentry config.
type BaseConfig struct {
	DSN             string `mapstructure:"dsn"`
	RateLimitWindow *int   `mapstructure:"rate_limit_window"`
	RateLimitCount  *int   `mapstructure:"rate_limit_count"`
}

// GlobalConfig global Sentry configuration.
type GlobalConfig struct {
	BaseConfig   `mapstructure:",squash"`
	AuthToken    string `mapstructure:"auth_token"`
	BaseURL      string `mapstructure:"base_url"`
	Project      string `mapstructure:"project"`
	Organization string `mapstructure:"organization"`
}

// SentryConfig is for site specific sentry DSN settings
type SiteConfig struct {
	BaseConfig `mapstructure:",squash"`
	Project    string                     `mapstructure:"project"`
	Components map[string]ComponentConfig `mapstructure:"-"`
}

// SentryConfig is for site specific sentry DSN settings
type ComponentConfig struct {
	BaseConfig `mapstructure:",squash"`
	Project    string `mapstructure:"project"`
}

func (c *SiteConfig) extendGlobalConfig(g *GlobalConfig) *SiteConfig {
	cfg := &SiteConfig{
		BaseConfig: g.BaseConfig,
		Project:    g.Project,
		Components: c.Components,
	}
	if c.DSN != "" {
		cfg.DSN = c.DSN
	}
	if c.RateLimitWindow != nil {
		cfg.RateLimitWindow = c.RateLimitWindow
	}
	if c.RateLimitCount != nil {
		cfg.RateLimitCount = c.RateLimitCount
	}
	if c.Project != "" {
		cfg.Project = c.Project
	}
	return cfg
}

func (c *ComponentConfig) extendSiteConfig(s *SiteConfig) *ComponentConfig {
	cfg := &ComponentConfig{
		BaseConfig: s.BaseConfig,
		Project:    s.Project,
	}

	if c.DSN != "" {
		cfg.DSN = c.DSN
	}
	if c.RateLimitWindow != nil {
		cfg.RateLimitWindow = c.RateLimitWindow
	}
	if c.RateLimitCount != nil {
		cfg.RateLimitCount = c.RateLimitCount
	}
	if c.Project != "" {
		cfg.Project = c.Project
	}
	return cfg
}

func (c *SiteConfig) getComponentSiteConfig(name string) *ComponentConfig {
	compConfig, ok := c.Components[name]
	if !ok {
		compConfig = ComponentConfig{}
	}
	return compConfig.extendSiteConfig(c)
}
