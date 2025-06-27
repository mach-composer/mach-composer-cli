package config

type ProviderConfigs []ProviderConfig

func (pc *ProviderConfigs) Names() ([]string, error) {
	var names []string
	for _, config := range *pc {
		names = append(names, config.Name)
	}
	return names, nil
}

type ProviderConfig struct {
	Name          string         `mapstructure:"name"`
	Source        string         `mapstructure:"source"`
	Version       string         `mapstructure:"version"`
	Configuration map[string]any `mapstructure:"configuration,omitempty"`
}
