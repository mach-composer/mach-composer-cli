package utils

import (
	"github.com/spf13/afero"
)

var (
	FS  = afero.NewOsFs()
	AFS = &afero.Afero{Fs: FS}
)
