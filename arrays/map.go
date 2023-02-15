package arrays

// Applies a map function over an array
func Map[v any, w any](inputs []v, mf func(v) w) []w {
	result := make([]w, len(inputs))
	for i := range inputs {
		result[i] = mf(inputs[i])
	}

	return result
}
