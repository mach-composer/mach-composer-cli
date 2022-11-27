package plugins

import (
	"github.com/labd/mach-composer/internal/plugins/amplience"
	"github.com/labd/mach-composer/internal/plugins/aws"
	"github.com/labd/mach-composer/internal/plugins/azure"
	"github.com/labd/mach-composer/internal/plugins/commercetools"
	"github.com/labd/mach-composer/internal/plugins/contentful"
	"github.com/labd/mach-composer/internal/plugins/sentry"
)

var localPlugins map[string]MachComposerPlugin

func init() {
	localPlugins = map[string]MachComposerPlugin{
		"amplience":     amplience.NewAmpliencePlugin(),
		"aws":           aws.NewAWSPlugin(),
		"azure":         azure.NewAzurePlugin(),
		"contentful":    contentful.NewContentfulPlugin(),
		"commercetools": commercetools.NewCommercetoolsPlugin(),
		"sentry":        sentry.NewSentryPlugin(),
	}
}
