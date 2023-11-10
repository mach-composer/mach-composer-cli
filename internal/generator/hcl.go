package generator

import (
	"encoding/json"
	"fmt"
	"github.com/mach-composer/mach-composer-cli/internal/config"
	"github.com/mach-composer/mach-composer-cli/internal/config/variable"
	"github.com/mach-composer/mach-composer-cli/internal/state"
	"regexp"
	"strings"

	"github.com/hashicorp/hcl/v2/hclwrite"
	ctyjson "github.com/zclconf/go-cty/cty/json"
)

var regexVars = regexp.MustCompilePOSIX(`"\$\$\{([^}]+)}"`)

func serializeToHCL(attributeName string, data variable.VariablesMap, deploymentType config.DeploymentType,
	repository *state.Repository) (string, error) {
	var transformFunc variable.TransformValueFunc
	switch deploymentType {
	case config.DeploymentSite:
		transformFunc = variable.ModuleTransformFunc()
		break
	case config.DeploymentSiteComponent:

		transformFunc = variable.RemoteStateTransformFunc(repository)
		break
	default:
		return "", fmt.Errorf("invalid deployment type: %s", deploymentType)
	}

	val, err := data.Transform(transformFunc)
	if err != nil {
		return "", err
	}

	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return "", err
	}

	var ctyJsonVal ctyjson.SimpleJSONValue
	if err := ctyJsonVal.UnmarshalJSON(jsonBytes); err != nil {
		return "", err
	}

	f := hclwrite.NewEmptyFile()
	rootBody := f.Body()
	rootBody.SetAttributeValue(attributeName, ctyJsonVal.Value)

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
