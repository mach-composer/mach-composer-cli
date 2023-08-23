package updater

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func TestMachConfigUpdater(t *testing.T) {
	data := []byte(utils.TrimIndent(`
        ---
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: 0.1.0
          endpoints:
            internal: internal
          integrations:
            - aws
            - commercetools
        -
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: "version"
          name: "foobar"
	`))

	updates := &UpdateSet{
		filename: "foo.yml",
		updates: []ChangeSet{
			{
				Component: &config.Component{
					Name:    "your-component",
					Version: "0.1.0",
				},
				LastVersion: "new-version",
			},
			{
				Component: &config.Component{
					Name:    "foobar",
					Version: "version",
				},
				LastVersion: "quoted",
			},
			{
				Component: &config.Component{
					Name:    "errror",
					Version: "not-in-file",
				},
				LastVersion: "quoted",
			},
		},
	}

	raw := &PartialRawConfig{}
	if err := yaml.Unmarshal(data, raw); err != nil {
		t.Fatal(err)
	}

	cfg := PartialConfig{
		ComponentsNode: &raw.Components,
	}

	output := machConfigUpdater(context.Background(), &cfg, data, updates)
	expected := []byte(utils.TrimIndent(`
        ---
        components:
        - name: your-component
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: "new-version"
          endpoints:
            internal: internal
          integrations:
            - aws
            - commercetools
        -
          source: "git::https://github.com/<username>/<your-component>.git//terraform"
          version: "quoted"
          name: "foobar"
	`))

	assert.EqualValues(t, string(expected), string(output))
}
