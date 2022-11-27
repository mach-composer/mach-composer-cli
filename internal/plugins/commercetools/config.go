package commercetools

import (
	"github.com/creasty/defaults"
)

type SiteConfig struct {
	ProjectKey      string                        `mapstructure:"project_key"`
	ClientID        string                        `mapstructure:"client_id"`
	ClientSecret    string                        `mapstructure:"client_secret"`
	Scopes          string                        `mapstructure:"scopes"`
	TokenURL        string                        `mapstructure:"token_url" default:"https://auth.europe-west1.gcp.commercetools.com"`
	APIURL          string                        `mapstructure:"api_url" default:"https://api.europe-west1.gcp.commercetools.com"`
	ProjectSettings *CommercetoolsProjectSettings `mapstructure:"project_settings"`

	Frontend *CommercetoolsFrontendSettings `mapstructure:"frontend"`

	Channels      []CommercetoolsChannel
	Taxes         []CommercetoolsTax
	TaxCategories []CommercetoolsTaxCategory `mapstructure:"tax_categories"`
	Stores        []CommercetoolsStore
	Zones         []CommercetoolsZone

	Components map[string]ComponentConfig `mapstructure:"-"`
}

func (s *SiteConfig) SetDefaults() {
	if defaults.CanUpdate(s.Frontend) {
		s.Frontend = &CommercetoolsFrontendSettings{
			CreateCredentials: true,
		}
		s.Frontend.SetDefaults()
	}
}

// ManagedStores returns all stores which are managed.
func (s *SiteConfig) ManagedStores() []CommercetoolsStore {
	managed := make([]CommercetoolsStore, 0)

	for _, store := range s.Stores {
		if store.Managed == nil || *store.Managed {
			managed = append(managed, store)
		}
	}
	return managed
}

func (c *SiteConfig) getComponentSiteConfig(name string) *ComponentConfig {
	compConfig, ok := c.Components[name]
	if !ok {
		compConfig = ComponentConfig{}
	}
	return &compConfig
}

type ComponentConfig struct {
	StoreVariables map[string]any `mapstructure:"store_variables"`
	StoreSecrets   map[string]any `mapstructure:"store_secrets"`
}

type CommercetoolsProjectSettings struct {
	Languages  []string `mapstructure:"languages"`
	Currencies []string `mapstructure:"currencies"`
	Countries  []string `mapstructure:"countries"`

	MessagesEnabled bool `mapstructure:"messages_enabled"`
}

type CommercetoolsFrontendSettings struct {
	CreateCredentials bool     `mapstructure:"create_credentials"`
	PermissionScopes  []string `mapstructure:"permission_scopes"`
}

func (s *CommercetoolsFrontendSettings) SetDefaults() {
	if defaults.CanUpdate(s.PermissionScopes) {
		s.PermissionScopes = []string{
			"create_anonymous_token",
			"manage_my_profile",
			"manage_my_orders",
			"manage_my_shopping_lists",
			"manage_my_payments",
			"view_products",
			"view_project_settings",
		}
	}
}

type CommercetoolsStore struct {
	Key                  string
	Name                 map[string]string
	Languages            []string
	DistributionChannels []string `mapstructure:"distribution_channels"`
	SupplyChannels       []string `mapstructure:"supply_channels"`

	// We use a pointer to a boolean. Otherwise the default value is false which
	// in turn is alwys set to true by the defaults module
	Managed *bool `mapstructure:"managed" default:"true"`

	// def __post_init__(self):
	//     if self.managed and not self.name:
	//         raise ValidationError("name is required")
}

type CommercetoolsChannel struct {
	Key         string
	Roles       []string
	Name        map[string]string
	Description map[string]string
}

type CommercetoolsTax struct {
	Country         string
	Amount          float64
	Name            string
	IncludedInPrice *bool `mapstructure:"included_in_price" default:"true"`
}

type CommercetoolsTaxCategory struct {
	Key   string
	Name  string
	Rates []CommercetoolsTax
}

type CommercetoolsZoneLocation struct {
	Country string
	State   string
}

type CommercetoolsZone struct {
	Name        string
	Description string
	Locations   []CommercetoolsZoneLocation
}
