package sentry

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
)

func Serve() {
	p := NewSentryPlugin()
	plugin.ServePlugin(p)
}
