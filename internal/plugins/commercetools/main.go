package commercetools

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
)

func Serve() {
	p := NewCommercetoolsPlugin()
	plugin.ServePlugin(p)
}
