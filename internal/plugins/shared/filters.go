package shared

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/flosch/pongo2/v5"
)

func init() {
	registerFilters()
}

func registerFilters() {
	MustRegisterFilter("render_commercetools_scopes", filterCommercetoolsScopes)
	MustRegisterFilter("ctvalue", FilterTFValue)
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
		res := pongo2.AsSafeValue(in.String())
		return res, nil

	case in.IsInteger():
		res := pongo2.AsSafeValue(fmt.Sprintf("%d", in.Integer()))
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
			sl = append(sl, fmt.Sprintf(`"%s"`, in.Index(i).String()))
		}

		result := pongo2.AsSafeValue(fmt.Sprintf("[%s]", strings.Join(sl, ", ")))
		return result, nil

	default:
		return nil, nil
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
