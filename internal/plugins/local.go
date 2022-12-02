package plugins

import (
	aws "github.com/mach-composer/mach-composer-plugin-aws/plugin"
	azure "github.com/mach-composer/mach-composer-plugin-azure/plugin"

	"github.com/labd/mach-composer/internal/plugins/amplience"
	"github.com/labd/mach-composer/internal/plugins/commercetools"
	"github.com/labd/mach-composer/internal/plugins/contentful"
	"github.com/labd/mach-composer/internal/plugins/sentry"
)

var LocalPluginNames = []string{"amplience", "aws", "azure", "contentful", "commercetools", "sentry"}

func GetLocalPlugins() map[string]func() {
	return map[string]func(){
		"amplience":     amplience.Serve,
		"aws":           aws.Serve,
		"azure":         azure.Serve,
		"contentful":    contentful.Serve,
		"commercetools": commercetools.Serve,
		"sentry":        sentry.Serve,
	}
}
