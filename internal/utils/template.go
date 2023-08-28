package utils

import (
	"bytes"
	"text/template"
)

func RenderGoTemplate(t string, data any) (string, error) {
	tpl, err := template.New("template").Parse(t)
	if err != nil {
		return "", err
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, data); err != nil {
		return "", err
	}
	return content.String(), nil
}
