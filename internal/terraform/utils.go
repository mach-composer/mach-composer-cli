package terraform

import (
	"github.com/rs/zerolog/log"
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

func terraformIsInitialized(path string) bool {
	tfLockFile := filepath.Join(path, ".terraform.lock.hcl")
	if _, err := os.Stat(tfLockFile); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Fatal().Err(err)
	}
	return true
}
