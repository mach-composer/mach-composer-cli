package state

import (
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/utils"
)

type GcpState struct {
	Bucket string `mapstructure:"bucket"`
	Prefix string `mapstructure:"prefix"`
}

func (a *GcpState) Key(key string) string {
	if a.Prefix == "" {
		return key
	}
	return fmt.Sprintf("%s/%s", a.Prefix, key)
}

type GcpRenderer struct {
	key   string
	state *GcpState
}

func (gr *GcpRenderer) Key() string {
	return gr.key
}

func (gr *GcpRenderer) Backend() (string, error) {
	templateContext := struct {
		Bucket string
		Prefix string
	}{
		Bucket: gr.state.Bucket,
		Prefix: gr.state.Key(gr.key),
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
		Bucket string
		Prefix string
	}{
		Bucket: gr.state.Bucket,
		Prefix: gr.state.Key(gr.key),
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
