package mcsdk

import (
	"encoding/gob"

	"github.com/hashicorp/go-plugin"
	"github.com/sirupsen/logrus"
)

// handshakeConfigs are used to just do a basic handshake between
// a plugin and host. If the handshake fails, a user friendly error is shown.
// This prevents users from executing bad plugins or executing a plugin
// directory. It is a UX feature, not a security feature.
var handshakeConfig = plugin.HandshakeConfig{
	ProtocolVersion:  1,
	MagicCookieKey:   "MACH_COMPOSER",
	MagicCookieValue: "plugin",
}

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}

func ServePlugin(p MachComposerPlugin) {
	logger := &Logger{
		logger: logrus.New(),
	}

	if val, ok := p.(*Adapter); ok {
		val.SetLogger(logger)
	}

	// pluginMap is the map of plugins we can dispense.
	var pluginMap = map[string]plugin.Plugin{
		"MachComposerPlugin": &Plugin{Impl: p},
	}

	plugin.Serve(&plugin.ServeConfig{
		HandshakeConfig: handshakeConfig,
		Plugins:         pluginMap,
		Logger:          logger,
	})
}
