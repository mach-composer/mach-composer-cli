package plugins

import (
	"fmt"
	"hash/crc32"
	"io"
	"log"
	"os"
	"os/exec"

	"github.com/hashicorp/go-plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/helpers"
	"github.com/mach-composer/mach-composer-plugin-sdk/protocol"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
)

func StartPlugin(name string) (schema.MachComposerPlugin, error) {
	var pluginMap = map[string]plugin.Plugin{
		"MachComposerPlugin": &protocol.Plugin{
			Identifier: name,
		},
	}

	// Create the checksum based on the current executable. This should be
	// retrieved from a 'catalog' and cached.
	pluginChecksum, err := getPluginChecksum(os.Args[0])
	if err != nil {
		return nil, err
	}

	// We're a host! Start by launching the plugin process.
	logger := helpers.NewLogger(nil)

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: protocol.HandShakeConfig(),
		Plugins:         pluginMap,
		Cmd:             exec.Command(os.Args[0], "plugin", name),
		Logger:          logger,
		SecureConfig: &plugin.SecureConfig{
			Hash:     crc32.NewIEEE(),
			Checksum: pluginChecksum,
		},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		log.Fatal(err)
	}

	// Request the plugin
	raw, err := rpcClient.Dispense("MachComposerPlugin")
	if err != nil {
		log.Fatal(err)
	}
	plugin, ok := raw.(schema.MachComposerPlugin)
	if !ok {
		return nil, fmt.Errorf("invalid plugin resolved")
	}
	return plugin, nil
}

func getPluginChecksum(filePath string) ([]byte, error) {
	h := crc32.NewIEEE()
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(h, file)
	if err != nil {
		return nil, err
	}

	return h.Sum(nil), nil
}
