package shared

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/flosch/pongo2/v5"
)

func init() {
	registerFilters()
}

func registerFilters() {
	MustRegisterFilter("render_commercetools_scopes", filterCommercetoolsScopes)
	MustRegisterFilter("tf", FilterTFValue)
	MustRegisterFilter("string", filterString)
	MustRegisterFilter("slugify", filterSlugify)
}

// mustRegisterFilter behaves like pongo2.RegisterFilter, but panics on an error.
func MustRegisterFilter(name string, filterFunc pongo2.FilterFunction) {
	if err := pongo2.RegisterFilter(name, filterFunc); err != nil {
		panic(fmt.Errorf("pongo2.RegisterFilter(%q): %v", name, err))
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

func FilterTFValue(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	switch {
	case in.IsString():
		return pongo2.AsSafeValue(formatString(in.String())), nil

	case in.IsInteger():
		res := pongo2.AsSafeValue(strconv.Itoa(in.Integer()))
		return res, nil

	case in.IsFloat():
		buf := strconv.FormatFloat(in.Float(), 'f', -1, 64)
		res := pongo2.AsSafeValue(buf)
		return res, nil

	case in.IsBool():
		if in.IsTrue() {
			return pongo2.AsValue("true"), nil
		}
		return pongo2.AsValue("false"), nil

	case in.CanSlice():
		sl := make([]string, 0, in.Len())
		for i := 0; i < in.Len(); i++ {
			sl = append(sl, formatString(in.Index(i).String()))
		}

		result := pongo2.AsSafeValue(fmt.Sprintf("[%s]", strings.Join(sl, ", ")))
		return result, nil

	default:
		return formatAny(in.Interface())
	}
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

func filterSlugify(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	return pongo2.AsValue(Slugify(in.String())), nil
}

var reTerraformRef = regexp.MustCompile(`^\$\{[^\s|\}]+\}$`)

// Format the string for terraform usage. If the value is a terraform reference
// then we don't quote the value and the remove ${} syntax.
// For example:
//
//	"${module.foo.bar}"     => module.foo.bar
//	"foobar"                => "foobar"
//	"foo ${module.foo.bar}" => "foo ${module.foo.bar}"
func formatString(val string) string {
	if reTerraformRef.MatchString(val) {
		return val[2 : len(val)-1]
	}
	return fmt.Sprintf(`"%s"`, EscapeChars(val))
}

func formatAny(val any) (*pongo2.Value, *pongo2.Error) {
	switch data := val.(type) {
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

	case string:
		return pongo2.AsSafeValue(formatString(data)), nil

	default:
		return pongo2.AsSafeValue(data), nil
	}
}

func formatMap[K comparable, V any](data map[K]V) (*pongo2.Value, *pongo2.Error) {
	items := make([]string, 0)
	for k, v := range data {
		formatted, err := formatAny(v)
		if err != nil {
			return nil, err
		}
		items = append(items, fmt.Sprintf("\t\t%v = %s,", k, formatted))
	}

	raw := fmt.Sprintf("{\n%s\n\t}", strings.Join(items, "\n"))
	return pongo2.AsSafeValue(raw), nil
}
