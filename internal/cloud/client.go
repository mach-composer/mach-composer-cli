package cloud

import (
	"net/http"

	"github.com/mach-composer/mcc-sdk-go/mccsdk"
)

type ClientConfig struct {
	URL        string
	HTTPClient *http.Client
}

func NewClient(cfg *ClientConfig) *mccsdk.APIClient {
	configuration := mccsdk.NewConfiguration()
	configuration.Servers = mccsdk.ServerConfigurations{
		{
			URL:         cfg.URL,
			Description: "Mach Composer Cloud Server",
		},
	}
	configuration.HTTPClient = cfg.HTTPClient

	return mccsdk.NewAPIClient(configuration)
}
