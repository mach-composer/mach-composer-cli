package plugins

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func TestResolveBuiltinPlugin(t *testing.T) {
	oldArgs := os.Args
	defer func() { os.Args = oldArgs }()

	os.Args = []string{"./mach-composer"}

	// Create dummy executable
	fh, err := utils.AFS.Create("mach-composer")
	require.NoError(t, err)
	if _, err := fh.WriteString("executable"); err != nil {
		require.NoError(t, err)
	}
	fh.Close()

	cfg := NewDefaultPlugin("aws")
	plugin, err := resolvePlugin(cfg)
	require.NoError(t, err)

	assert.Equal(t, "./mach-composer", plugin.Path)
	assert.Equal(t, []string{"plugin", "aws"}, plugin.Args)
}

func TestQueryPluginRegistry(t *testing.T) {
	// Set up test data
	pluginCfg := PluginConfig{
		Version: "1.0.0",
		Source:  "example-plugin",
	}

	expectedResponse := &registryResponse{
		URL: "https://example.org/download.zip",
	}

	// Create a mock HTTP server to handle the request
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		os := r.URL.Query().Get("os")
		arch := r.URL.Query().Get("arch")
		version := r.URL.Query().Get("version")

		assert.Equal(t, runtime.GOOS, os)
		assert.Equal(t, runtime.GOARCH, arch)
		assert.Equal(t, "1.0.0", version)

		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(expectedResponse)
		require.NoError(t, err)
	}))
	defer ts.Close()

	// Modify the URL's host to point to the mock server
	endpoint := registryEndpoint
	defer func() { registryEndpoint = endpoint }()
	registryEndpoint = ts.URL

	// Call the function being tested
	res, err := queryPluginRegistry(pluginCfg)
	require.NoError(t, err)
	assert.Equal(t, expectedResponse, res)
}
