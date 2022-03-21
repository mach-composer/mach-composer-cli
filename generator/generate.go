package generator

import (
	"embed"

	"github.com/flosch/pongo2/v5"
	"github.com/labd/mach-composer-go/config"
)

//go:embed templates/*
var templates embed.FS

func Render(cfg *config.Root, site *config.Site) (string, error) {
	templateSet := pongo2.NewSet("", &EmbedLoader{Content: templates})

	registerFilters()

	// TODO: these filter differ from the current mach-composer python version
	// due to jinja2 vs pongo2 limitations. Needs more research
	pongo2.RegisterFilter("replace", filterReplace)

	var tpl = pongo2.Must(templateSet.FromFile("site.tf"))

	out, err := tpl.Execute(pongo2.Context{"global": cfg.Global, "site": site})

	if err != nil {
		panic(err)
	}
	return out, nil
}
