package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"text/template"
)

func RenderGoTemplate(t string, data any) (string, error) {
	funcMap := template.FuncMap{
		"tfValue": func(v any) string {
			switch val := v.(type) {
			case string:
				return `"` + val + `"`
			case int, int64, float64, bool:
				return fmt.Sprintf("%v", val)
			default:
				b, _ := json.Marshal(val)
				return string(b)
			}
		},
		"isMap": func(v any) bool {
			_, ok := v.(map[string]any)
			return ok
		},
	}

	tpl, err := template.New("template").Funcs(funcMap).Parse(t)
	if err != nil {
		return "", err
	}

	var content bytes.Buffer
	if err := tpl.Execute(&content, data); err != nil {
		return "", err
	}
	return content.String(), nil
}
