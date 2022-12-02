package plugins

import (
	aws "github.com/mach-composer/mach-composer-plugin-aws/plugin"
	azure "github.com/mach-composer/mach-composer-plugin-azure/plugin"

	"github.com/labd/mach-composer/internal/plugins/amplience"
	"github.com/labd/mach-composer/internal/plugins/commercetools"
	"github.com/labd/mach-composer/internal/plugins/contentful"
	"github.com/labd/mach-composer/internal/plugins/sentry"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
)

var localPlugins map[string]schema.MachComposerPlugin

func init() {
	localPlugins = map[string]schema.MachComposerPlugin{
		// "amplience": amplience.NewAmpliencePlugin(),
		// "aws":           aws.NewAWSPlugin(),
		// "azure":         azure.NewAzurePlugin(),
		// "contentful":    contentful.NewContentfulPlugin(),
		// "commercetools": commercetools.NewCommercetoolsPlugin(),
		// "sentry":        sentry.NewSentryPlugin(),
	}
}

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
