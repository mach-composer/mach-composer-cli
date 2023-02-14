package generator

import (
	"reflect"
	"strings"

	"github.com/elliotchance/pie/v2"
)

func findDependsOn(s map[string]any) []string {
	result := []string{}
	for _, v := range s {
		if v == nil {
			continue
		}

		if value, ok := v.(string); ok {
			vars := findVariables(value)

			for _, variable := range vars {
				if strings.HasPrefix(variable, "module.") {
					items := strings.SplitN(variable, ".", 3)
					result = append(result, strings.Join(items[0:2], "."))
				}
			}
		}

		if reflect.TypeOf(v).Kind() == reflect.Map {
			if data, ok := v.(map[string]any); ok {
				items := findDependsOn(data)
				result = append(result, items...)
			}
		}
	}
	return pie.Unique(result)
}
