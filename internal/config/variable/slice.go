package variable

import "slices"

type SliceVariable struct {
	baseVariable
	Elements []Variable
}

func NewSliceVariable(elements []Variable) *SliceVariable {
	return &SliceVariable{baseVariable: baseVariable{typ: Slice}, Elements: elements}
}

func (v *SliceVariable) ReferencedComponents() []string {
	var references []string

	for _, element := range v.Elements {
		references = append(references, element.ReferencedComponents()...)
	}

	return slices.Compact(references)
}

func (v *SliceVariable) TransformValue(f TransformValueFunc) (any, error) {
	var data = make([]any, 0, len(v.Elements))

	for _, element := range v.Elements {
		dat, err := element.TransformValue(f)
		if err != nil {
			return nil, err
		}
		data = append(data, dat)
	}

	return data, nil
}
