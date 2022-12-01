package commercetools

import (
	"github.com/labd/mach-composer/internal/plugins/mcsdk"
)

func Serve() {
	p := NewCommercetoolsPlugin()
	mcsdk.ServePlugin(p)
}
