package variable

import "slices"

type MapVariable struct {
	baseVariable
	elements map[string]Variable
}

func NewMapVariable(elements map[string]Variable) *MapVariable {
	return &MapVariable{baseVariable: baseVariable{typ: Map}, elements: elements}
}

func (v *MapVariable) Elements() map[string]Variable {
	return v.elements
}

func (v *MapVariable) ReferencedComponents() []string {
	var references []string

	for _, element := range v.elements {
		references = append(references, element.ReferencedComponents()...)
	}

	return slices.Compact(references)
}

func (v *MapVariable) TransformValue(f TransformValueFunc) (any, error) {
	var data = make(map[string]any, len(v.elements))
	var err error

	for key, element := range v.elements {
		data[key], err = element.TransformValue(f)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}
