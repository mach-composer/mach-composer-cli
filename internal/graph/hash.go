package graph

import (
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

func HashSiteComponent(sc *SiteComponent) (string, error) {
	var err error
	var tfHash string
	var variablesHash string

	// If the source is local, we include the hash of the directory.
	if sc.SiteComponentConfig.Definition.Source.IsType(config.SourceTypeLocal) {
		tfHash, err = utils.ComputeDirHash(string(sc.SiteComponentConfig.Definition.Source))
		if err != nil {
			return "", err
		}
	}

	// If a variables file has been set, we include the hash in the component hash.
	// This is necessary because the variables file can be changed without
	// changing the config itself, which would not update the hash. It is very costly though,
	// as it will rerun all the components in the project
	if sc.ProjectConfig.MachComposer.VariablesFile != "" {
		variablesHash, err = utils.ComputeFileHash(sc.ProjectConfig.MachComposer.VariablesFile)
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
		Variables     variable.VariablesMap `json:"variables"`
		Secrets       variable.VariablesMap `json:"secrets"`
		DependsOn     []string              `json:"depends_on"`
		Terraform     string                `json:"terraform"`
		VariablesFile string                `json:"variables_file"`
	}{
		Name: sc.SiteComponentConfig.Name,
		Definition: struct {
			Name    string        `json:"name"`
			Version string        `json:"version"`
			Source  config.Source `json:"source"`
			Branch  string        `json:"branch"`
		}{
			Name:    sc.SiteComponentConfig.Definition.Name,
			Version: sc.SiteComponentConfig.Definition.Version,
			Source:  sc.SiteComponentConfig.Definition.Source,
			Branch:  sc.SiteComponentConfig.Definition.Branch,
		},
		Variables:     sc.SiteComponentConfig.Variables,
		Secrets:       sc.SiteComponentConfig.Secrets,
		DependsOn:     sc.SiteComponentConfig.DependsOnKeys,
		Terraform:     tfHash,
		VariablesFile: variablesHash,
	})
}
