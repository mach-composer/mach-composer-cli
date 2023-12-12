package utils

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
	"github.com/zclconf/go-cty/cty/json"
	"os"
	"os/exec"
)

// RunTerraform will execute a terraform command with the given arguments in the given directory.
func RunTerraform(ctx context.Context, catchOutputs bool, cwd string, args ...string) (string, error) {
	if _, err := os.Stat(cwd); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("the generated files are not found: %w", err)
		}
	}

	execPath, err := exec.LookPath("terraform")
	if err != nil {
		return "", err
	}

	return RunInteractive(ctx, catchOutputs, execPath, cwd, args...)
}

func GetTerraformOutputs(ctx context.Context, path string) (cty.Value, error) {
	var data json.SimpleJSONValue

	logger := log.Ctx(ctx).With().Str("path", path).Logger()

	output, err := RunTerraform(ctx, true, path, "output", "-json")
	if err != nil {
		logger.Error().Err(err).Msgf("failed to get terraform output: %s", err.Error())
		return cty.NilVal, err
	}

	logger.Debug().Str("output", output).Msgf("Fetched terraform output")

	if err = data.UnmarshalJSON([]byte(output)); err != nil {
		logger.Error().Err(err).Str("output", output).Msgf("failed to unmarshal terraform output: %s", err.Error())
		return cty.NilVal, err
	}

	return data.Value, nil
}

type SiteComponentOutput struct {
	Sensitive bool `cty:"sensitive"`
	Value     struct {
		Hash      *string    `cty:"hash"`
		Variables *cty.Value `cty:"variables"`
	} `cty:"value"`
	Type cty.Value `cty:"type"`
}

// ParseSiteComponentOutputByKey returns the output of a terraform command for the given key at the given path.
// If no output is found nil is returned.
func ParseSiteComponentOutputByKey(val cty.Value, key string) (*SiteComponentOutput, error) {
	if !val.Type().HasAttribute(key) {
		log.Debug().Msgf("no attribute found for key %s", key)
		return nil, nil
	}

	componentVal := val.GetAttr(key)

	var scOut SiteComponentOutput
	err := gocty.FromCtyValue(componentVal, &scOut)
	if err != nil {
		log.Err(err).Msgf("failed to convert terraform output to SiteComponentOutput: %s", err.Error())
		return nil, err
	}

	return &scOut, nil
}
