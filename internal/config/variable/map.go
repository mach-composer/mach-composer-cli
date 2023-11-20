package variable

import "slices"

type MapVariable struct {
	baseVariable
	Elements map[string]Variable
}

func NewMapVariable(elements map[string]Variable) *MapVariable {
	return &MapVariable{baseVariable: baseVariable{typ: Map}, Elements: elements}
}

func (v *MapVariable) ReferencedComponents() []string {
	var references []string

	for _, element := range v.Elements {
		references = append(references, element.ReferencedComponents()...)
	}

	return slices.Compact(references)
}

func (v *MapVariable) TransformValue(f TransformValueFunc) (any, error) {
	var data = make(map[string]any, len(v.Elements))
	var err error

	for key, element := range v.Elements {
		data[key], err = element.TransformValue(f)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
