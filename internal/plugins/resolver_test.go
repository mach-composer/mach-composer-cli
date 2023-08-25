package plugins

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
