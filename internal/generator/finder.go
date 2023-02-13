package generator

import (
	"reflect"
	"strings"
)

func findDependsOn(s map[string]any) []string {
	result := []string{}
	for _, v := range s {

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
			items := findDependsOn(v.(map[string]any))
			result = append(result, items...)
		}
	}
	return result
}
