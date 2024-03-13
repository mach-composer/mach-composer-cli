package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func HashSiteComponentConfig(sc config.SiteComponentConfig) (string, error) {
	var err error
	var tfHash string

	// If the source is local, we include the hash of the directory.
	if sc.Definition.Source.IsType(config.SourceTypeLocal) {
		tfHash, err = utils.ComputeDirHash(string(sc.Definition.Source))
		if err != nil {
			return "", err
		}
	}

	return utils.ComputeHash(struct {
		Name       string `json:"name"`
		Definition struct {
			Name    string        `json:"name"`
			Version string        `json:"version"`
			Source  config.Source `json:"source"`
			Branch  string        `json:"branch"`
		} `json:"definition"`
		Variables variable.VariablesMap `json:"variables"`
		Secrets   variable.VariablesMap `json:"secrets"`
		DependsOn []string              `json:"depends_on"`
		Terraform string                `json:"terraform"`
	}{
		Name: sc.Name,
		Definition: struct {
			Name    string        `json:"name"`
			Version string        `json:"version"`
			Source  config.Source `json:"source"`
			Branch  string        `json:"branch"`
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
