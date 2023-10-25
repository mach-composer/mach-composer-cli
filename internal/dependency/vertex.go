package dependency

import "regexp"

type Vertices []Node

func (v Vertices) FilterOne(pattern string) Node {
	for _, vertex := range v {
		matched, _ := regexp.MatchString(pattern, vertex.Path())
		if matched {
			return vertex
		}
	}
	return nil
}
