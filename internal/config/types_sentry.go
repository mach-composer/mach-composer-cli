package config

// SentryConfigBase is the base sentry config.
type SentryConfigBase struct {
	DSN             string `yaml:"dsn"`
	RateLimitWindow int    `yaml:"rate_limit_window"`
	RateLimitCount  int    `yaml:"rate_limit_count"`
}

// GlobalSentryConfig global Sentry configuration.
type GlobalSentryConfig struct {
	SentryConfigBase `yaml:",inline"`
	AuthToken        string `yaml:"auth_token"`
	BaseURL          string `yaml:"base_url"`
	Project          string `yaml:"project"`
	Organization     string `yaml:"organization"`
}

// SentryConfig is for site specific sentry DSN settings
type SentryConfig struct {
	SentryConfigBase `yaml:",inline"`
	Project          string `yaml:"project"`
}

func newSentryConfigFromGlobal(c *GlobalSentryConfig) *SentryConfig {
	return &SentryConfig{
		Project: c.Project,
		SentryConfigBase: SentryConfigBase{
			DSN:             c.DSN,
			RateLimitCount:  c.RateLimitCount,
			RateLimitWindow: c.RateLimitWindow,
		},
	}
}

func newSentryConfig(c *SentryConfig) *SentryConfig {
	return &SentryConfig{
		Project: c.Project,
		SentryConfigBase: SentryConfigBase{
			DSN:             c.DSN,
			RateLimitCount:  c.RateLimitCount,
			RateLimitWindow: c.RateLimitWindow,
		},
	}
}

func (sc *SentryConfig) merge(c *SentryConfig) {
	if sc.DSN == "" {
		sc.DSN = c.DSN
	}
	if sc.Project == "" {
		sc.Project = c.Project
	}
	if sc.RateLimitCount == 0 {
		sc.RateLimitCount = c.RateLimitCount
	}
	if sc.RateLimitWindow == 0 {
		sc.RateLimitWindow = c.RateLimitWindow
	}
}

func (sc *SentryConfig) mergeGlobal(c *GlobalSentryConfig) {
	if sc.DSN == "" {
		sc.DSN = c.DSN
	}
	if sc.Project == "" {
		sc.Project = c.Project
	}
	if sc.RateLimitCount == 0 {
		sc.RateLimitCount = c.RateLimitCount
	}
	if sc.RateLimitWindow == 0 {
		sc.RateLimitWindow = c.RateLimitWindow
	}
}
