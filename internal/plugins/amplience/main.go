package amplience

import (
	"github.com/labd/mach-composer/internal/plugins/mcsdk"
)

func Serve() {
	p := NewAmpliencePlugin()
	mcsdk.ServePlugin(p)
}
