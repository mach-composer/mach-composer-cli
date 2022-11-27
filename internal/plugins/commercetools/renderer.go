package commercetools

import (
	"embed"

	"github.com/flosch/pongo2/v5"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

//go:embed templates/*
var templates embed.FS

func Render(cfg *SiteConfig) (string, error) {
	templateSet := pongo2.NewSet("", &shared.EmbedLoader{Content: templates})
	template := pongo2.Must(templateSet.FromFile("main.tf"))

	return template.Execute(pongo2.Context{
		"commercetools": cfg,
	})
}
