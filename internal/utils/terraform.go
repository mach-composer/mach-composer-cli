package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
)

type TerraformOutput map[string]SiteComponentOutput

func (t *TerraformOutput) GetSiteComponentOutput(key string) (*SiteComponentOutput, bool) {
	if t == nil {
		return nil, false
	}

	if output, ok := (*t)[key]; ok {
		return &output, true
	}

	return nil, false
}

type SiteComponentOutput struct {
	Value struct {
		Hash      string      `json:"hash"`
		Variables interface{} `json:"variables"`
	} `json:"value"`
}

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

func GetTerraformOutput(ctx context.Context, path string) (*TerraformOutput, error) {
	var data TerraformOutput

	output, err := RunTerraform(ctx, path, "output", "-json")
	if err != nil {
		return nil, err
	}

	if err = json.Unmarshal([]byte(output), &data); err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, nil
	}

	return &data, nil
}
