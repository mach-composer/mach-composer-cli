package generator

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flosch/pongo2/v5"
	"github.com/sirupsen/logrus"
)

func registerFilters() {
	mustRegisterFilter("interpolate", filterInterpolate)
}

// mustRegisterFilter behaves like pongo2.RegisterFilter, but panics on an error.
func mustRegisterFilter(name string, filterFunc pongo2.FilterFunction) {
	if err := pongo2.RegisterFilter(name, filterFunc); err != nil {
		panic(fmt.Errorf("pongo2.RegisterFilter(%q): %v", name, err))
	}
}

func filterInterpolate(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	switch {
	case in.IsString():
		val, err := ParseTemplateVariable(in.String())
		if err != nil {
			logrus.Fatal(err.Error())
		}
		res := pongo2.AsSafeValue(EscapeChars(val))
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
}

func formatMap[K comparable, V any](data map[K]V) (*pongo2.Value, *pongo2.Error) {
	items := make([]string, 0)
	for k, v := range data {
		formatted, err := filterInterpolate(pongo2.AsSafeValue(v), nil)
		if err != nil {
			continue
		}
		items = append(items, fmt.Sprintf("\t\t%v = %s,", k, formatted))
	}

	raw := fmt.Sprintf("{\n%s\n\t}", strings.Join(items, "\n"))
	return pongo2.AsSafeValue(raw), nil
}
