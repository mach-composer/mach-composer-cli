package generator

import (
	"errors"
	"fmt"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/labd/mach-composer/config"
	"github.com/sirupsen/logrus"
)

func registerFilters() {
	pongo2.RegisterFilter("string", filterString)
	pongo2.RegisterFilter("short_prefix", filterShortPrefix)
	pongo2.RegisterFilter("slugify", filterSlugify)
	pongo2.RegisterFilter("remove", filterRemove)
	pongo2.RegisterFilter("tf", filterTFValue)
	pongo2.RegisterFilter("tfvalue", filterTFValue)
	pongo2.RegisterFilter("azure_frontend_endpoint_name", AzureFrontendEndpointName)
	pongo2.RegisterFilter("service_plan_resource_name", AzureServicePlanResourceName)
	pongo2.RegisterFilter("get", FilterGetValueByKey)
	pongo2.RegisterFilter("render_commercetools_scopes", filterCommercetoolsScopes)
	pongo2.RegisterFilter("component_endpoint_name", filterComponentEndpointName)
}

func FilterGetValueByKey(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	val := in.Interface()

	if val == nil {
		return pongo2.AsValue(nil), nil
	}

	switch v := val.(type) {
	case map[string]any:
		key := param.String()
		retVal := v[key]
		return pongo2.AsValue(retVal), nil
	}

	return nil, &pongo2.Error{
		Sender:    "filter:get",
		OrigError: fmt.Errorf("invalid value for get: %v", val),
	}
}

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(Slugify(in.String())), nil
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

// Take a component and a site-endpoint, and return a Terraform reference to
// an output. The endpoint might have a different name in the component itself
// based on the mappings
func filterComponentEndpointName(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	component := in.Interface().(config.SiteComponent)
	endpoint := param.Interface().(config.Endpoint)
	for componentKey, epKey := range component.Definition.Endpoints {
		if epKey == endpoint.Key {
			return pongo2.AsSafeValue(componentKey), nil
		}

	}
	return nil, &pongo2.Error{
		Sender:    "filter:render_commercetools_scopes",
		OrigError: fmt.Errorf("endpoint %s not found on %s", endpoint.Key, component.Name),
	}
}

func filterTFValue(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	if in.IsString() {
		val, err := ParseTemplateVariable(in.String())
		if err != nil {
			logrus.Fatal(err.Error())
		}
		res := pongo2.AsSafeValue(fmt.Sprintf(`"%s"`, val))
		return res, nil
	}
	if in.IsInteger() {
		res := pongo2.AsSafeValue(fmt.Sprintf("%d", in.Integer()))
		return res, nil
	}
	if in.IsFloat() {
		res := pongo2.AsSafeValue(fmt.Sprintf("%f", in.Float()))
		return res, nil
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

	case map[string]string:
		return formatMap(data)

	case map[string]any:
		return formatMap(data)

	case map[any]any:
		// Should not be necessary if the formatMap is fixed
		items := make(map[string]any, 0)
		for k, v := range data {
			items[fmt.Sprint(k)] = v
		}
		return formatMap(items)

	default:
		return pongo2.AsValue(data), nil
	}

}

func formatMap[K comparable, V any](data map[K]V) (*pongo2.Value, *pongo2.Error) {
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
