package utils

func UniqueSlice[T comparable](input []T) []T {
	u := make([]T, 0, len(input))
	m := make(map[T]bool)
	for _, val := range input {
		if _, ok := m[val]; !ok {
			m[val] = true
			u = append(u, val)
		}
	}
	return u
}
