package state

import "github.com/mach-composer/mach-composer-cli/internal/utils"

type LocalState struct {
	Path string `mapstructure:"path"`
}

type LocalRenderer struct {
	state *LocalState
	key   string
}

func (lr *LocalRenderer) Key() string {
	return lr.key
}

func (lr *LocalRenderer) Backend() (string, error) {
	templateContext := struct {
		State *LocalState
		Key   string
	}{
		State: lr.state,
		Key:   lr.key,
	}

	tpl := `
	backend "remote" {
		{{ if .State.Path }}
		path = "{{ .State.Path }}"
		{{ end }}
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}

func (lr *LocalRenderer) RemoteState() (string, error) {
	templateContext := struct {
		State *LocalState
		Key   string
	}{
		State: lr.state,
		Key:   lr.key,
	}

	//TODO: this needs fixing

	tpl := `
	data "terraform_remote_state" "{{ .Key }}" {
	  backend = "local"
	
	  config = {
		{{ if .State.Path }}
		path = "{{ .State.Path }}"
		{{ end }}
	  }
	}
	`
	return utils.RenderGoTemplate(tpl, templateContext)
}
