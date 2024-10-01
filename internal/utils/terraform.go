package utils

import (
	"context"
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/json"
	"os"
	"os/exec"
)

// RunTerraform will execute a terraform command with the given arguments in the given directory.
func RunTerraform(ctx context.Context, cwd string, args ...string) (string, error) {
	if _, err := os.Stat(cwd); err != nil {
		if os.IsNotExist(err) {
			return "", fmt.Errorf("the generated files are not found: %w", err)
		}
	}

	execPath, err := exec.LookPath("terraform")
	if err != nil {
		return "", err
	}

	return RunInteractive(ctx, execPath, cwd, args...)
}

func GetTerraformOutputs(ctx context.Context, path string) (cty.Value, error) {
	var data json.SimpleJSONValue

	output, err := RunTerraform(ctx, path, "output", "-json")
	if err != nil {
		log.Error().Err(err).Msgf("failed to get terraform output: %s", err.Error())
		return cty.NilVal, err
	}

	log.Debug().Str("output", output).Msgf("Fetched terraform output")

	if err = data.UnmarshalJSON([]byte(output)); err != nil {
		log.Error().Err(err).Str("output", output).Msgf("failed to unmarshal terraform output: %s", err.Error())
		return cty.NilVal, err
	}

	return data.Value, nil
}
