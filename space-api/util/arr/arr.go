package arr

func MapSlice[T any, E any](source []T, mapFunc func(index int, current T) E) (out []E) {
	for index, val := range source {
		out = append(out, mapFunc(index, val))
	}

	return
}

func FilterSlice[T any](source []T, filterFunc func(current T, index int) bool) (out []T) {
	for index, val := range source {
		if filterFunc(val, index) {
			out = append(out, val)
		}
	}

	return
}
