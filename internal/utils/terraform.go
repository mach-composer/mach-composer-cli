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

type HashOutput struct {
	Sensitive bool    `cty:"sensitive"`
	Type      string  `cty:"type"`
	Value     *string `cty:"value"`
}

// ParseHashOutput returns the hash output by the given key.
func ParseHashOutput(val cty.Value) (string, error) {
	if !val.Type().HasAttribute("hash") {
		return "", fmt.Errorf("no attribute with key hash found in terraform output")
	}

	componentVal := val.GetAttr("hash")

	var hashOutput HashOutput
	err := gocty.FromCtyValue(componentVal, &hashOutput)
	if err != nil {
		log.Err(err).Msgf("failed to convert terraform output to HashOutput: %s", err.Error())
		return "", err
	}

	if hashOutput.Value == nil {
		return "", fmt.Errorf("no value set for hash in terraform output")
	}

	return *hashOutput.Value, nil
}
