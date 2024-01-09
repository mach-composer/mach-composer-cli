package plugins

import (
	"encoding/json"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"hash/crc32"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"regexp"
	"runtime"
	"strings"

	"github.com/hashicorp/go-getter"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/rs/zerolog/log"
)

var registryEndpoint = "https://registry.mach.cloud"

var rePluginName = regexp.MustCompile("(?i)^([a-z0-9_-]+)/([a-z0-9_-]+)$")

type pluginExecutable struct {
	Checksum []byte
	Path     string
	Args     []string // Backwards compatible
}

type registryResponse struct {
	URL string `json:"url"`
}

func (p pluginExecutable) command() *exec.Cmd {
	return exec.Command(p.Path, p.Args...)
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

func resolvePlugin(pluginCfg PluginConfig) (*pluginExecutable, error) {

	if err := validateSource(pluginCfg.Source); err != nil {
		return nil, err
	}

	// Download the plugin if we don't have it yet
	if _, err := os.Stat(pluginCfg.path()); err != nil {
		log.Debug().Msgf("Plugin %s %s not found, trying to download", pluginCfg.name(), pluginCfg.Version)
		if err := downloadPlugin(pluginCfg); err != nil {
			return nil, fmt.Errorf("failed to download plugin %s: %w", pluginCfg.name(), err)
		}
	}

	pluginChecksum, err := getPluginChecksum(pluginCfg.path())
	if err != nil {
		return nil, err
	}

	result := &pluginExecutable{
		Path:     pluginCfg.path(),
		Checksum: pluginChecksum,
	}
	return result, nil
}

// downloadPlugin queries the registry for the download url, extracts it to a
// temporary directory and the moves it to the mach-composer plugin directory.
func downloadPlugin(pluginCfg PluginConfig) error {
	log.Info().Msgf("Downloading %s (%s)...", pluginCfg.Source, pluginCfg.Version)

	info, err := queryPluginRegistry(pluginCfg)
	if err != nil {
		return err
	}

	tempDir, err := os.MkdirTemp("", "mcp")
	if err != nil {
		return fmt.Errorf("failed to create temporary dir: %w", err)
	}
	defer os.RemoveAll(tempDir)

	client := getter.Client{
		DisableSymlinks: true,
		Src:             info.URL,
		Dst:             tempDir,
		Mode:            getter.ClientModeDir,
	}
	if err := client.Get(); err != nil {
		return err
	}

	filename := path.Join(tempDir, pluginCfg.executableName())

	if _, err := os.Stat(filename); err != nil {
		return err
	}

	target := pluginCfg.path()
	if err := os.MkdirAll(path.Dir(target), 0700); err != nil {
		return err
	}

	if err := os.Rename(filename, pluginCfg.path()); err != nil {
		return err
	}

	return nil
}

func queryPluginRegistry(pluginCfg PluginConfig) (*registryResponse, error) {
	params := url.Values{
		"version": {pluginCfg.Version},
		"os":      {runtime.GOOS},
		"arch":    {runtime.GOARCH},
	}

	u, err := url.Parse(registryEndpoint)
	if err != nil {
		return nil, fmt.Errorf("internal error: %w", err)
	}

	u.Path = fmt.Sprintf("/v1/plugins/%s", strings.ToLower(pluginCfg.Source))
	u.RawQuery = params.Encode()

	client := retryablehttp.NewClient()
	client.Logger = NewHCLogAdapter(log.Logger)

	r, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}

	defer r.Body.Close()

	if r.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("plugin does not exist")
	}

	info := &registryResponse{}
	if err := json.NewDecoder(r.Body).Decode(info); err != nil {
		return nil, err
	}
	return info, nil
}

func validateSource(name string) error {
	if !rePluginName.MatchString(name) {
		return fmt.Errorf("invalid plugin name: %s", name)
	}
	return nil
}
