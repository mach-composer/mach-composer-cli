package commercetools

import (
	"embed"

	"github.com/flosch/pongo2/v5"
	"github.com/mach-composer/mach-composer-plugin-helpers/helpers"
)

//go:embed templates/*
var templates embed.FS

func renderResources(cfg *SiteConfig) (string, error) {
	templateSet := pongo2.NewSet("", &helpers.EmbedLoader{Content: templates})
	template := pongo2.Must(templateSet.FromFile("main.tf"))

	return template.Execute(pongo2.Context{
		"commercetools": cfg,
	})
}
