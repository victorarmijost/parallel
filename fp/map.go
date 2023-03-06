package fp

// Applies a map function over an array
func Map[v any, w any](input []v, mf func(v) w) []w {
	result := make([]w, len(input))
	for _, i := range input {
		result = append(result, mf(i))
	}

	return result
}
