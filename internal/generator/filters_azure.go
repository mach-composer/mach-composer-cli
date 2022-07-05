package generator

import (
	"fmt"

	"github.com/flosch/pongo2/v5"
	"github.com/labd/mach-composer/internal/config"
)

func AzureFrontendEndpointName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.Interface().(config.Endpoint)

	if val.Azure != nil && val.Azure.InternalName != "" {
		return filterTFValue(pongo2.AsSafeValue(val.Azure.InternalName), nil)
	}
	return filterTFValue(pongo2.AsSafeValue(val.Key), nil)

}

// AzureServicePlanResourceName Retrieve the resource name for a Azure app service plan.
// The reason to make this conditional is because of backwards compatability;
// existing environments already have a `functionapp` resource. We want to keep that intact.
func AzureServicePlanResourceName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.String()

	if val == "default" {
		return pongo2.AsSafeValue("functionapps"), nil
	}
	return pongo2.AsSafeValue(fmt.Sprintf("functionapps_%s", val)), nil
}
