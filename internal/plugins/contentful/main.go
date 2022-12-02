package contentful

import "github.com/mach-composer/mach-composer-plugin-sdk/plugin"

func Serve() {
	p := NewContentfulPlugin()
	plugin.ServePlugin(p)
}
