package terraform

import (
	"os"
	"path/filepath"
)

const PlanFile = "terraform.plan"

func hasTerraformPlan(path string) (string, error) {
	filename := filepath.Join(path, PlanFile)
	if _, err := os.Stat(filename); err == nil {
		return filename, nil
	}
	return "", nil
}
