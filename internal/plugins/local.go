package plugins

import (
	amplience "github.com/mach-composer/mach-composer-plugin-amplience/plugin"
	aws "github.com/mach-composer/mach-composer-plugin-aws/plugin"
	azure "github.com/mach-composer/mach-composer-plugin-azure/plugin"
	commercetools "github.com/mach-composer/mach-composer-plugin-commercetools/plugin"
	contentful "github.com/mach-composer/mach-composer-plugin-contentful/plugin"
	sentry "github.com/mach-composer/mach-composer-plugin-sentry/plugin"
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
