package plugins

import (
	"context"
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"strings"

	"github.com/hashicorp/go-plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/protocol"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/rs/zerolog/log"

	"github.com/labd/mach-composer/internal/utils"
)

type Plugin struct {
	schema.MachComposerPlugin
	Name string

	client    *plugin.Client
	isRunning bool
	config    PluginConfig
}

func (p *Plugin) start(ctx context.Context) error {
	if p.isRunning {
		return nil
	}

	var pluginMap = map[string]plugin.Plugin{
		"MachComposerPlugin": &protocol.Plugin{
			Identifier: p.Name,
		},
	}

	// Safety check to not run external plugins during test for now
	if strings.HasSuffix(os.Args[0], ".test") {
		panic(fmt.Sprintf("Not loading %s: invalid command: %s", p.Name, os.Args[0]))
	}

	logger := NewHCLogAdapter(log.Logger)

	executable, err := resolvePlugin(p.config)
	if err != nil {
		return err
	}

	p.client = plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: protocol.HandShakeConfig(),
		Plugins:         pluginMap,
		Cmd:             executable.command(),
		Logger:          logger,
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
		log.Fatal().Err(err).Msg("failed to connect to plugin")
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("MachComposerPlugin")
	if err != nil {
		log.Fatal().Err(err).Msg("failed to init the plugin")
	}

	plugin, ok := raw.(schema.MachComposerPlugin)
	if !ok {
		return fmt.Errorf("incompatible plugin resolved")
	}
	p.MachComposerPlugin = plugin
	p.isRunning = true

	return nil
}

func getPluginChecksum(filePath string) ([]byte, error) {
	h := crc32.NewIEEE()
	file, err := utils.AFS.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to get checksum of file: %w", err)
	}
	defer file.Close()

	_, err = io.Copy(h, file)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
