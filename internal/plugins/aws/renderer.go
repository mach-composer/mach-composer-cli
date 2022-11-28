package aws

import (
	"embed"

	"github.com/elliotchance/pie/v2"
	"github.com/flosch/pongo2/v5"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

//go:embed templates/*
var templates embed.FS

func renderResources(site, env string, cfg *SiteConfig, endpoints []EndpointConfig) (string, error) {
	templateSet := pongo2.NewSet("", &shared.EmbedLoader{Content: templates})
	template := pongo2.Must(templateSet.FromFile("resources/main.tf"))

	// check if there is an endpoint with a CDN enabled.
	// TODO: Is this used ?
	enableCDN := false
	dnsZones := []string{}
	for _, endpoint := range endpoints {
		if endpoint.EnableCDN {
			enableCDN = true
		}
		dnsZones = append(dnsZones, endpoint.Zone)
	}

	return template.Execute(pongo2.Context{
		"aws":       cfg,
		"siteName":  site,
		"envName":   env,
		"endpoints": endpoints,
		"enableCDN": enableCDN,
		"dnsZones":  pie.Unique(dnsZones),
	})
}
