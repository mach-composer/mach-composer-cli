package shared

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
)

var regexVars = regexp.MustCompilePOSIX(`"\$\$\{([^\}]+)\}"`)

func SerializeToHCL(attributeName string, data any) string {
	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	rootBody.SetAttributeValue(attributeName, transformToCTY(data))
	return fixVariableReference(string(f.Bytes()))
}

func fixVariableReference(data string) string {
	matches := regexVars.FindAllStringSubmatch(data, -1)
	for _, match := range matches {
		data = strings.Replace(data, match[0], match[1], 1)
	}

	return data
}

func transformToCTY(source any) cty.Value {
	switch v := source.(type) {
	case string:
		return cty.StringVal(v)
	case int:
		return cty.NumberIntVal(int64(v))
	case float32:
		return cty.NumberFloatVal(float64(v))
	case float64:
		return cty.NumberFloatVal(v)
	case bool:
		return cty.BoolVal(v)
	}

	val := reflect.ValueOf(source)
	if val.Kind() == reflect.Map {
		result := map[string]cty.Value{}
		for _, e := range val.MapKeys() {
			v := val.MapIndex(e)
			k := fmt.Sprintf("%v", e.Interface())
			result[k] = transformToCTY(v.Interface())
		}
		return cty.ObjectVal(result)
	}

	if val.Kind() == reflect.Slice {
		result := []cty.Value{}
		for i := 0; i < val.Len(); i++ {
			result = append(result, transformToCTY(val.Index(i).Interface()))
		}
		return cty.ListVal(result)
	}

	return cty.NilVal
}
