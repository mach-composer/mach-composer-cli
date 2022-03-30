package generator

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/labd/mach-composer-go/config"
)

func registerFilters() {
	pongo2.RegisterFilter("string", filterString)
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("tf", filterTFValue)
	pongo2.RegisterFilter("tfvalue", filterTFValue)
	pongo2.RegisterFilter("azure_region_short", filterAzureRegionShort)
	pongo2.RegisterFilter("azure_region_long", filterAzureRegionLong)
	pongo2.RegisterFilter("azure_frontend_endpoint_name", AzureFrontendEndpointName)
	pongo2.RegisterFilter("service_plan_resource_name", AzureServicePlanResourceName)
	pongo2.RegisterFilter("get", FilterGetValueByKey)
	pongo2.RegisterFilter("render_commercetools_scopes", filterCommercetoolsScopes)
	pongo2.RegisterFilter("component_endpoint_name", filterComponentEndpointName)
	pongo2.RegisterFilter("has_cloud_integration", filterHasCloudIntegration)
}

func FilterGetValueByKey(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.Interface()

	if val == nil {
		return pongo2.AsValue(nil), nil
	}

	switch v := val.(type) {
	case map[string]interface{}:
		key := param.String()
		retval := v[key]
		return pongo2.AsValue(retval), nil
	}

	return nil, &pongo2.Error{
		Sender:    "filter:get",
		OrigError: fmt.Errorf("invalid value for get: %v", val),
	}
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(Slugify(in.String())), nil
}

func filterReplace(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	items := strings.Split(param.String(), ",")

	if len(items) != 2 {
		return nil, &pongo2.Error{
			Sender:    "filter:replace",
			OrigError: errors.New("replace needs two arguments"),
		}
	}

	output := strings.Replace(in.String(), items[0], items[1], -1)
	return pongo2.AsValue(output), nil
}

func filterString(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsBool() {
		if in.Bool() {
			return pongo2.AsValue("true"), nil
		} else {
			return pongo2.AsValue("false"), nil
		}
	}
	return in, nil
}

// Take an component and an site-endpoint, and return a Terraform reference to
// an output. The endpoint might have a different name in the component itself
// based on the mappings
func filterComponentEndpointName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	component := in.Interface().(config.SiteComponent)
	endpoint := param.Interface().(config.Endpoint)
	for component_key, ep_key := range component.Definition.Endpoints {
		if ep_key == endpoint.Key {
			return pongo2.AsSafeValue(component_key), nil
		}

	}
	return nil, &pongo2.Error{
		Sender:    "filter:render_commercetools_scopes",
		OrigError: fmt.Errorf("endpoint %s not found on %s", endpoint.Key, component.Name),
	}
}

func filterTFValue(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsString() {
		val := pongo2.AsSafeValue(fmt.Sprintf(`"%s"`, in.String()))
		return val, nil
	}
	if in.IsInteger() {
		val := pongo2.AsSafeValue(fmt.Sprintf("%d", in.Integer()))
		return val, nil
	}
	if in.IsFloat() {
		val := pongo2.AsSafeValue(fmt.Sprintf("%f", in.Float()))
		return val, nil
	}
	if in.IsBool() {
		if in.IsTrue() {
			return pongo2.AsValue("true"), nil
		} else {
			return pongo2.AsValue("false"), nil
		}
	}
	if in.CanSlice() {
		sl := make([]string, 0, in.Len())
		for i := 0; i < in.Len(); i++ {
			sl = append(sl, fmt.Sprintf(`"%s"`, in.Index(i).String()))
		}

		result := pongo2.AsSafeValue(fmt.Sprintf("[%s]", strings.Join(sl, ", ")))
		return result, nil
	}

	switch data := in.Interface().(type) {

	case map[interface{}]interface{}:
		{
			items := make([]string, 0)
			for k, v := range data {
				formatted, err := filterTFValue(pongo2.AsSafeValue(v), nil)
				if err != nil {
					continue
				}

				items = append(items, fmt.Sprintf("\t\t%v = %s,", k, formatted))
			}

			raw := fmt.Sprintf("{\n%s\n\t}", strings.Join(items, "\n"))
			return pongo2.AsSafeValue(raw), nil
		}

	case map[string]interface{}:
		{
			items := make([]string, 0)
			for k, v := range data {
				formatted, err := filterTFValue(pongo2.AsSafeValue(v), nil)
				if err != nil {
					continue
				}

				items = append(items, fmt.Sprintf("  %v = %s", k, formatted))
			}

			raw := fmt.Sprintf("{\n%s\n\t}", strings.Join(items, "\n"))
			return pongo2.AsSafeValue(raw), nil
		}
	default:
		return pongo2.AsValue(data), nil
	}

}

func filterCommercetoolsScopes(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if !in.CanSlice() {
		return nil, &pongo2.Error{
			Sender:    "filter:render_commercetools_scopes",
			OrigError: errors.New("input is not sliceable"),
		}
	}

	projectKey := param.String()
	sl := make([]string, in.Len())
	for i := 0; i < in.Len(); i++ {
		sl = append(sl, fmt.Sprintf(`"%s:%s",`, in.Index(i).String(), projectKey))
	}

	result := pongo2.AsSafeValue(fmt.Sprintf("[\n  %s\n]", strings.Join(sl, "")))
	return result, nil
}

func filterHasCloudIntegration(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.Interface()

	switch v := val.(type) {
	case config.SiteComponent:
		{
			if v.Definition == nil {
				log.Fatalf("Component %s was not resolved properly (missing definition)", v.Name)
			}
			for _, i := range v.Definition.Integrations {
				if i == "aws" || i == "azure" {
					return pongo2.AsValue(true), nil
				}
			}
			return pongo2.AsValue(false), nil
		}
	}
	return in, nil

}
