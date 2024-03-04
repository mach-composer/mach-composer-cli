package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
	"github.com/rs/zerolog/log"
)

func hashSiteComponentConfig(sc config.SiteComponentConfig) (string, error) {
	var err error
	var tfHash string

	if !sc.Definition.IsGitSource() {
		tfHash, err = utils.ComputeDirHash(sc.Definition.Source)
		if err != nil {
			return "", err
		}
	} else {
		log.Debug().Msgf("Site component %s is a git source, so could not generate hash for terraform code", sc.Name)
	}

	return utils.ComputeHash(struct {
		Name       string `json:"name"`
		Definition struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Source  string `json:"source"`
			Branch  string `json:"branch"`
		} `json:"definition"`
		Variables variable.VariablesMap `json:"variables"`
		Secrets   variable.VariablesMap `json:"secrets"`
		DependsOn []string              `json:"depends_on"`
		Terraform string                `json:"terraform"`
	}{
		Name: sc.Name,
		Definition: struct {
			Name    string `json:"name"`
			Version string `json:"version"`
			Source  string `json:"source"`
			Branch  string `json:"branch"`
		}{
			Name:    sc.Definition.Name,
			Version: sc.Definition.Version,
			Source:  sc.Definition.Source,
			Branch:  sc.Definition.Branch,
		},
		Variables: sc.Variables,
		Secrets:   sc.Secrets,
		DependsOn: sc.DependsOn,
		Terraform: tfHash,
	})
}
