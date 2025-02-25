package utils

import (
	"bytes"
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

	logger := log.Ctx(ctx).With().
		Str(CommandFieldName, execPath).
		Strs(ArgsFieldName, args).
		Str(CwdFieldName, cwd).
		Logger()

	logger.Debug().Msgf("Running: %s", cwd)
	defer logger.Debug().Msgf("Finished running: %s", cwd)

	w := logger.Hook(StdHook{Logger: logger})

	return RunInteractive(ctx, execPath, cwd, w, args...)
}

func GetTerraformOutputs(ctx context.Context, path string) (cty.Value, error) {
	var data json.SimpleJSONValue

	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return cty.NilVal, fmt.Errorf("the generated files are not found: %w", err)
		}
	}

	execPath, err := exec.LookPath("terraform")
	if err != nil {
		return cty.NilVal, err
	}

	w := new(bytes.Buffer)

	_, err = RunInteractive(ctx, execPath, path, w, "output", "-json")
	if err != nil {
		log.Error().Err(err).Msgf("failed to get terraform output: %s", err.Error())
		return cty.NilVal, err
	}
	output := w.String()

	log.Debug().Str("output", output).Msgf("Fetched terraform output")

	if err = data.UnmarshalJSON([]byte(output)); err != nil {
		log.Error().Err(err).Str("output", output).Msgf("failed to unmarshal terraform output: %s", err.Error())
		return cty.NilVal, err
	}

	return data.Value, nil
}
