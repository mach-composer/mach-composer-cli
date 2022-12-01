package azure

import (
	"github.com/labd/mach-composer/internal/plugins/mcsdk"
)

func Serve() {
	p := NewAzurePlugin()
	mcsdk.ServePlugin(p)
}
