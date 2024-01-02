package plugins

import (
	"context"
	"fmt"
	"github.com/hashicorp/go-plugin"
	protocolv1 "github.com/mach-composer/mach-composer-plugin-sdk/protocol"
	schemav1 "github.com/mach-composer/mach-composer-plugin-sdk/schema"
	protocolv2 "github.com/mach-composer/mach-composer-plugin-sdk/v2/protocol"
	schemav2 "github.com/mach-composer/mach-composer-plugin-sdk/v2/schema"
	"github.com/rs/zerolog/log"
	"hash/crc32"
	"os"
	"strings"
)

type PluginHandler struct {
	schemav2.MachComposerPlugin
	Name string

	client    *plugin.Client
	isRunning bool
	Config    PluginConfig
}

func (p *PluginHandler) Close() {
	p.client.Kill()
	p.isRunning = false
	p.client = nil
}

func (p *PluginHandler) Start(ctx context.Context) error {
	if p.isRunning {
		return nil
	}

	var versionedPlugins = map[int]plugin.PluginSet{
		1: {"MachComposerPlugin": &protocolv1.Plugin{Identifier: p.Name}},
		2: {"MachComposerPlugin": &protocolv2.Plugin{Identifier: p.Name}},
	}

	// Safety check to not run external handlers during test for now
	if strings.HasSuffix(os.Args[0], ".test") {
		panic(fmt.Sprintf("Not loading %s: invalid command: %s", p.Name, os.Args[0]))
	}

	logger := NewHCLogAdapter(log.Logger)

	executable, err := resolvePlugin(p.Config)
	if err != nil {
		return err
	}

	p.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: plugin.HandshakeConfig{
			ProtocolVersion:  1,
			MagicCookieKey:   protocolv2.HandShakeConfig().MagicCookieKey,
			MagicCookieValue: protocolv2.HandShakeConfig().MagicCookieValue,
		},
		VersionedPlugins: versionedPlugins,
		Cmd:              executable.command(),
		Logger:           logger,
		SecureConfig: &plugin.SecureConfig{
			Hash:     crc32.NewIEEE(),
			Checksum: executable.Checksum,
		},
	})

	go func() {
		<-ctx.Done()
		log.Debug().Msg("Received cancellation")
		p.client.Kill()
	}()

	// Connect via RPC
	rpcClient, err := p.client.Client()
	if err != nil {
		return fmt.Errorf("failed to connect to plugin: %s", err.Error())
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("MachComposerPlugin")
	if err != nil {
		return fmt.Errorf("failed to dispense plugin: %s", err.Error())
	}

	protocolVersion := p.client.NegotiatedVersion()
	switch protocolVersion {
	case 1:
		machComposerPlugin, ok := raw.(schemav1.MachComposerPlugin)
		if !ok {
			return fmt.Errorf("incompatible machComposerPlugin resolved")
		}
		p.MachComposerPlugin = NewPluginV1Adapter(machComposerPlugin)
		log.Debug().Msgf("Loading plugin with protocol version %d", protocolVersion)

	case 2:
		machComposerPlugin, ok := raw.(schemav2.MachComposerPlugin)
		if !ok {
			return fmt.Errorf("incompatible machComposerPlugin resolved")
		}
		p.MachComposerPlugin = machComposerPlugin
		log.Debug().Msgf("Loading plugin with protocol version %d", protocolVersion)
	default:
		return fmt.Errorf("incompatible protocol version %d. Try upgrading your mach-composer binary to the latest version", protocolVersion)
	}
	p.isRunning = true

	return nil
}
