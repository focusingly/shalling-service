package arr

func MapSlice[T any, E any](source []T, mapFunc func(index int, current T) E) (out []E) {
	for index, val := range source {
		out = append(out, mapFunc(index, val))
	}

	return
}
