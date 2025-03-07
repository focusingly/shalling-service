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

func Compress[T any, E ~[]T](src E, eqFn func(current, receive T) bool) E {
	ret := []T{}

	for inputIndex := range src {
		if inputIndex < 1 {
			ret = append(ret, src[inputIndex])
		} else {
			exists := false
			for currentIndex := range ret {
				current, cmp := ret[currentIndex], src[inputIndex]
				if eqFn(current, cmp) {
					exists = true
					break
				}
			}
			if !exists {
				ret = append(ret, src[inputIndex])
			}
		}
	}

	return ret
}
