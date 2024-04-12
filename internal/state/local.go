package state

import "github.com/mach-composer/mach-composer-cli/internal/utils"

type LocalState struct {
	Path string `mapstructure:"path"`
}

type LocalRenderer struct {
	BaseRenderer
	state *LocalState
}

func (lr *LocalRenderer) Backend() (string, error) {
	templateContext := struct {
		State      *LocalState
		Identifier string
	}{
		State:      lr.state,
		Identifier: lr.identifier,
	}

	tpl := `
	backend "local" {
		{{ if .State.Path }}
		path = "{{ .State.Path }}/{{ .Identifier }}.tfstate"
		{{ end }}
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (lr *LocalRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State      *LocalState
		Identifier string
		Key        string
	}{
		State:      lr.state,
		Identifier: lr.identifier,
		Key:        lr.key,
	}

	tpl := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "local"
	
	  config = {
		{{ if .State.Path }}
		path = "{{ .State.Path }}/{{ .Identifier }}.tfstate"
		{{ end }}
	  }
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}
