package azure

import (
	"errors"
	"strings"

	"github.com/flosch/pongo2/v5"

	"github.com/labd/mach-composer/internal/plugins/shared"
)

func init() {
	shared.MustRegisterFilter("azure_frontend_endpoint_name", AzureFrontendEndpointName)
	shared.MustRegisterFilter("service_plan_resource_name", AzureServicePlanResourceName)
	shared.MustRegisterFilter("component_endpoint_name", filterComponentEndpointName)
	shared.MustRegisterFilter("short_prefix", filterShortPrefix)
	shared.MustRegisterFilter("remove", filterRemove)
}

func AzureFrontendEndpointName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// val, ok := in.Interface().(config.Endpoint)
	// if !ok {
	// 	return nil, &pongo2.Error{
	// 		Sender:    "filter:azure_frontend_endpoint_name",
	// 		OrigError: fmt.Errorf("filter expected argument of type Endpoint"),
	// 	}
	// }

	// if val.Azure != nil && val.Azure.InternalName != "" {
	// 	return shared.FilterTFValue(pongo2.AsSafeValue(val.Azure.InternalName), nil)
	// }
	// return shared.FilterTFValue(pongo2.AsSafeValue(val.Key), nil)
	return pongo2.AsSafeValue("TODO"), nil
}

// AzureServicePlanResourceName Retrieve the resource name for a Azure app service plan.
// The reason to make this conditional is because of backwards compatability;
// existing environments already have a `functionapp` resource. We want to keep that intact.
func AzureServicePlanResourceName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := azureServicePlanResourceName(in.String())
	return pongo2.AsSafeValue(val), nil
}

// Take a component and a site-endpoint, and return a Terraform reference to
// an output. The endpoint might have a different name in the component itself
// based on the mappings
func filterComponentEndpointName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	// component, ok := in.Interface().(config.SiteComponent)
	// if !ok {
	// 	return nil, &pongo2.Error{
	// 		Sender:    "filter:component_endpoint_name",
	// 		OrigError: fmt.Errorf("filter only works on site component"),
	// 	}
	// }
	// endpoint, ok := param.Interface().(config.Endpoint)
	// if !ok {
	// 	return nil, &pongo2.Error{
	// 		Sender:    "filter:component_endpoint_name",
	// 		OrigError: fmt.Errorf("filter expected argument of type Endpoint"),
	// 	}
	// }

	// for componentKey, epKey := range component.Definition.Endpoints {
	// 	if epKey == endpoint.Key {
	// 		return pongo2.AsSafeValue(componentKey), nil
	// 	}
	// }
	// return nil, &pongo2.Error{
	// 	Sender:    "filter:component_endpoint_name",
	// 	OrigError: fmt.Errorf("endpoint %s not found on %s", endpoint.Key, component.Name),
	// }
	return nil, nil
}

// Specific function created to be backwards compatible with Python version
// It replaces env names with 1 letter codes.
// TODO: Research why/if this is still needed
func filterShortPrefix(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return nil, &pongo2.Error{
			Sender:    "filter:short_string",
			OrigError: errors.New("filter only applicable on strings"),
		}
	}

	val := in.String()
	val = strings.Replace(val, "dev", "d", -1)
	val = strings.Replace(val, "tst", "t", -1)
	val = strings.Replace(val, "acc", "a", -1)
	val = strings.Replace(val, "prd", "p", -1)
	return pongo2.AsValue(val), nil
}

func filterRemove(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.IsString() {
		return nil, &pongo2.Error{
			Sender:    "filter:remove",
			OrigError: errors.New("filter only applicable on strings"),
		}
	}
	if !param.IsString() {
		return nil, &pongo2.Error{
			Sender:    "filter:remove",
			OrigError: errors.New("filter requires a param"),
		}
	}

	output := strings.Replace(in.String(), param.String(), "", -1)
	return pongo2.AsValue(output), nil
}
