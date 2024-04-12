package state

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type GcpState struct {
	Bucket string `mapstructure:"bucket"`
	Prefix string `mapstructure:"prefix"`
}

func (a *GcpState) Identifier(identifier string) string {
	if a.Prefix == "" {
		return identifier
	}
	return fmt.Sprintf("%s/%s", a.Prefix, identifier)
}

type GcpRenderer struct {
	BaseRenderer
	state *GcpState
}

func (gr *GcpRenderer) Backend() (string, error) {
	templateContext := struct {
		Bucket string
		Prefix string
	}{
		Bucket: gr.state.Bucket,
		Prefix: gr.state.Identifier(gr.identifier),
	}

	tpl := `
	backend "gcs" {
	  bucket  = "{{ .Bucket }}"
	  prefix = "{{ .Prefix }}"
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (gr *GcpRenderer) RemoteState() (string, error) {
	templateContext := struct {
		Identifier string
		Bucket     string
		Prefix     string
		Key        string
	}{
		Identifier: gr.identifier,
		Bucket:     gr.state.Bucket,
		Prefix:     gr.state.Identifier(gr.identifier),
		Key:        gr.key,
	}

	template := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "gcp"
	
	  config = {
		  bucket  = "{{ .Bucket }}"
		  prefix = "{{ .Prefix }}"
	  }
	}
	`
	return utils.RenderGoTemplate(template, templateContext)
}
