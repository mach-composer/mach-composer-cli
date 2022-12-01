package sentry

import (
	"github.com/labd/mach-composer/internal/plugins/mcsdk"
)

func Serve() {
	p := NewSentryPlugin()
	mcsdk.ServePlugin(p)
}
