package plugins

import (
	"fmt"
	"path"
	"path/filepath"
	"runtime"

	"github.com/adrg/xdg"
)

type PluginConfig struct {
	Source  string
	Version string
}

func (pc PluginConfig) name() string {
	return filepath.Base(pc.Source)
}

func (pc PluginConfig) executableName() string {
	return fmt.Sprintf("mach-composer-plugin-%s_v%s", pc.name(), pc.Version)
}

func (pc PluginConfig) path() string {
	return path.Join(
		xdg.ConfigHome, "mach-composer", "plugins", pc.Source, pc.Version,
		fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH),
		pc.name())
}

func NewDefaultPlugin(name string) PluginConfig {
	return PluginConfig{Source: name, Version: "builtin"}
}
