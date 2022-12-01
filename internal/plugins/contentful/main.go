package contentful

import (
	"github.com/labd/mach-composer/internal/plugins/mcsdk"
)

func Serve() {
	p := NewContentfulPlugin()
	mcsdk.ServePlugin(p)
}
