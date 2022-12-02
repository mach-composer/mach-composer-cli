package amplience

import (
	"github.com/mach-composer/mach-composer-plugin-sdk/plugin"
)

func Serve() {
	p := NewAmpliencePlugin()
	plugin.ServePlugin(p)
}
