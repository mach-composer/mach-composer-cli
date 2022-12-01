package mcsdk

import (
	"net/rpc"

	"github.com/hashicorp/go-plugin"
)

// This is the implementation of plugin.Plugin so we can serve/consume this
//
// This has two methods: Server must return an RPC server for this plugin
// type. We construct a GreeterRPCServer for this.
//
// Client must return an implementation of our interface that communicates
// over an RPC client. We return GreeterRPC for this.
//
// Ignore MuxBroker. That is used to create more multiplexed streams on our
// plugin connection and is a more advanced use case.
type Plugin struct {
	// Impl Injection
	Impl       MachComposerPlugin
	Identifier string
}

func (p *Plugin) Server(*plugin.MuxBroker) (any, error) {
	return &PluginRPCServer{Impl: p.Impl}, nil
}

func (p *Plugin) Client(b *plugin.MuxBroker, c *rpc.Client) (any, error) {
	return &PluginRPC{client: c, identifier: p.Identifier}, nil
}
