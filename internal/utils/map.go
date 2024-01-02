package utils

import (
	"github.com/elliotchance/pie/v2"
)

func FilterMap(data map[string]any, skipKeys []string) map[string]any {
	var filteredData = map[string]any{}
	for key, datum := range data {
		if pie.Contains(skipKeys, key) {
			continue
		}
		filteredData[key] = datum
	}

	return filteredData
}
