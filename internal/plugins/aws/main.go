package aws

import (
	"github.com/labd/mach-composer/internal/plugins/mcsdk"
)

func Serve() {
	p := NewAWSPlugin()
	mcsdk.ServePlugin(p)
}
