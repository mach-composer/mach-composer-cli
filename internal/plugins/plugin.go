package plugins

import (
	"fmt"
	"hash/crc32"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/elliotchance/pie/v2"
	"github.com/hashicorp/go-plugin"
	"github.com/mach-composer/mach-composer-plugin-sdk/protocol"
	"github.com/mach-composer/mach-composer-plugin-sdk/schema"
	"github.com/rs/zerolog/log"
)

func startPlugin(name string) (schema.MachComposerPlugin, error) {
	var pluginMap = map[string]plugin.Plugin{
		"MachComposerPlugin": &protocol.Plugin{
			Identifier: name,
		},
	}

	// Safety check to not run external plugins during test for now
	if strings.HasSuffix(os.Args[0], ".test") {
		panic(fmt.Sprintf("Not loading %s: invalid command: %s", name, os.Args[0]))
	}

	logger := NewHCLogAdapter(log.Logger)

	command := fmt.Sprintf("mach-composer-plugin-%s", name)
	args := []string{}

	if pie.Contains(LocalPluginNames, name) {
		command = os.Args[0]

		// If mach-composer is started from the $PATH the we need to resolve
		// the command
		if !strings.Contains(command, "/") {
			path, err := exec.LookPath(command)
			if err != nil {
				return nil, err
			}
			command = path
		}
		args = []string{"plugin", name}
	} else {
		path, err := exec.LookPath(command)
		if err != nil {
			return nil, err
		}

		command = path
	}

	pluginChecksum, err := getPluginChecksum(command)
	if err != nil {
		return nil, err
	}

	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: protocol.HandShakeConfig(),
		Plugins:         pluginMap,
		Cmd:             exec.Command(command, args...),
		Logger:          logger,
		SecureConfig: &plugin.SecureConfig{
			Hash:     crc32.NewIEEE(),
			Checksum: pluginChecksum,
		},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
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
