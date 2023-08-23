package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/mach-composer/mach-composer-cli/internal/config"
)

func TestFileLocations(t *testing.T) {
	cfg := &config.MachConfig{
		Sites: []config.SiteConfig{
			{
				Identifier: "my-site",
			},
		},
	}

	options := &GenerateOptions{
		OutputPath: "deployments/something",
	}

	actual := FileLocations(cfg, options)
	expected := map[string]string{
		"my-site": "deployments/something/my-site",
	}

	assert.EqualValues(t, expected, actual)
}
