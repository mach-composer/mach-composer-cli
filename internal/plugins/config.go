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
	executableName := fmt.Sprintf("mach-composer-plugin-%s_v%s", pc.name(), pc.Version)

	if runtime.GOOS == "windows" {
		executableName += ".exe"
	}

	return executableName
}

func (pc PluginConfig) path() string {
	path := path.Join(
		xdg.ConfigHome, "mach-composer", "plugins", pc.Source, pc.Version,
		fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH),
		pc.name())

	if runtime.GOOS == "windows" {
		path += ".exe"
	}

	return path
}

func NewDefaultPlugin(name string) PluginConfig {
	return PluginConfig{Source: name, Version: "builtin"}
}
