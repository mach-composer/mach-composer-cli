package config

import (
	"github.com/creasty/defaults"
)

type CommercetoolsSettings struct {
	ProjectKey      string                        `yaml:"project_key"`
	ClientID        string                        `yaml:"client_id"`
	ClientSecret    string                        `yaml:"client_secret"`
	Scopes          string                        `yaml:"scopes"`
	TokenURL        string                        `yaml:"token_url" default:"https://auth.europe-west1.gcp.commercetools.com"`
	APIURL          string                        `yaml:"api_url" default:"https://api.europe-west1.gcp.commercetools.com"`
	ProjectSettings *CommercetoolsProjectSettings `yaml:"project_settings"`

	Frontend *CommercetoolsFrontendSettings `yaml:"frontend"`

	Channels      []CommercetoolsChannel
	Taxes         []CommercetoolsTax
	TaxCategories []CommercetoolsTaxCategory `yaml:"tax_categories"`
	Stores        []CommercetoolsStore
	Zones         []CommercetoolsZone
}

func (s *CommercetoolsSettings) SetDefaults() {
	if defaults.CanUpdate(s.Frontend) {
		s.Frontend = &CommercetoolsFrontendSettings{
			CreateCredentials: true,
		}
		s.Frontend.SetDefaults()
	}
}

// ManagedStores returns all stores which are managed.
func (s *CommercetoolsSettings) ManagedStores() []CommercetoolsStore {
	managed := make([]CommercetoolsStore, 0)

	for _, store := range s.Stores {
		if store.Managed == nil || *store.Managed {
			managed = append(managed, store)
		}
	}
	return managed
}

type CommercetoolsProjectSettings struct {
	Languages  []string `yaml:"languages"`
	Currencies []string `yaml:"currencies"`
	Countries  []string `yaml:"countries"`

	MessagesEnabled *bool `yaml:"messages_enabled" default:"true"`
}

type CommercetoolsFrontendSettings struct {
	CreateCredentials bool     `yaml:"create_credentials" default:"false"`
	PermissionScopes  []string `yaml:"permission_scopes"`
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
	DistributionChannels []string `yaml:"distribution_channels"`
	SupplyChannels       []string `yaml:"supply_channels"`

	// We use a pointer to a boolean. Otherwise the default value is false which
	// in turn is alwys set to true by the defaults module
	Managed *bool `yaml:"managed" default:"true"`

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
	IncludedInPrice *bool `yaml:"included_in_price" default:"true"`
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
