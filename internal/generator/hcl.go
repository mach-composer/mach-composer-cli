package generator

import (
	"encoding/json"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	"github.com/zclconf/go-cty/cty"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var regexVars = regexp.MustCompilePOSIX(`"\$\$\{([^\}]+)\}"`)
var regexTFVar = regexp.MustCompilePOSIX(`\$\{([^\}]+)\}`)

func serializeToHCL(attributeName string, data any) (string, error) {
	val, err := asCTY(data)
	if err != nil {
		return "", err
	}

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	rootBody.SetAttributeValue(attributeName, val)

	result := fixVariableReference(string(f.Bytes()))
	return result, nil
}

func fixVariableReference(data string) string {
	matches := regexVars.FindAllStringSubmatch(data, -1)
	for _, match := range matches {
		replacement := match[1]

		// Unescape quotes. Required for secret references, e.g.:
		// 	data.sops_external.variables.data[\"my-key\"]
		// should become:
		// 	data.sops_external.variables.data["my-key"]
		replacement = strings.ReplaceAll(replacement, `\"`, `"`)
		data = strings.Replace(data, match[0], replacement, 1)
	}

	return data
}

func asCTY(source any) (cty.Value, error) {
	jsonBytes, err := json.Marshal(source)
	if err != nil {
		return cty.NilVal, err
	}
	var ctyJsonVal ctyjson.SimpleJSONValue
	if err := ctyJsonVal.UnmarshalJSON(jsonBytes); err != nil {
		return cty.NilVal, err
	}

	return ctyJsonVal.Value, nil
}

func findVariables(input string) []string {
	matches := regexTFVar.FindAllStringSubmatch(input, -1)
	result := []string{}
	for _, match := range matches {
		item := match[1]
		result = append(result, item)
	}
	return result
}
