package arr

func MapSlice[T any, E any](source []T, mapFunc func(index int, current T) E) (out []E) {
	out = []E{}
	for index, val := range source {
		out = append(out, mapFunc(index, val))
	}

	return
}

func FilterSlice[T any](source []T, filterFunc func(current T, index int) bool) (out []T) {
	out = []T{}
	for index, val := range source {
		if filterFunc(val, index) {
			out = append(out, val)
		}
	}

	return
}

func Compress[T any, E ~[]T](arr E, containFunc func(T, E) bool) E {
	newPack := []T{}
	for _, val := range arr {
		if !containFunc(val, newPack) {
			newPack = append(newPack, val)
		}
	}

	return newPack
}
