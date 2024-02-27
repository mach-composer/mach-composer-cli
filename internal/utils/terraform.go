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

type MissingHashError struct {
	message string
}

func NewMissingHashError(message string) *MissingHashError {
	return &MissingHashError{message: message}
}

func (m *MissingHashError) Error() string {
	return m.message
}

// RunTerraform will execute a terraform command with the given arguments in the given directory.
func RunTerraform(ctx context.Context, cwd string, catchOutputs bool, args ...string) (string, error) {
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

	output, err := RunTerraform(ctx, path, true, "output", "-json")
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

type HashOutput struct {
	Sensitive bool    `cty:"sensitive"`
	Type      string  `cty:"type"`
	Value     *string `cty:"value"`
}

// ParseHashOutput returns the hash output by the given key.
func ParseHashOutput(val cty.Value) (string, error) {
	if !val.Type().HasAttribute("hash") {
		return "", NewMissingHashError("no attribute with key hash found in terraform output")
	}

	componentVal := val.GetAttr("hash")

	var hashOutput HashOutput
	err := gocty.FromCtyValue(componentVal, &hashOutput)
	if err != nil {
		return "", err
	}

	if hashOutput.Value == nil {
		return "", NewMissingHashError("no value set for hash in terraform output")
	}

	return *hashOutput.Value, nil
}
