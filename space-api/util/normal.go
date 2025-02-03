package util

import (
	"fmt"
	"iter"
)

// GetSerial 根据 seed 生成函数返回可迭代序列
func GetSerial[T any](seed func() T, optionLimit ...int) iter.Seq[T] {
	prevIndex := 0
	return func(yield func(T) bool) {
		switch {
		case len(optionLimit) == 0:
			for yield(seed()) {
			}
		case len(optionLimit) == 1:
			limit := optionLimit[0]
			if limit > 0 {
				for prevIndex < limit {
					prevIndex++
					if !yield(seed()) {
						break
					}
				}
			} else {
				panic(fmt.Sprintf("require a positive num, but got: %d", optionLimit[0]))
			}
		default:
			panic(fmt.Sprintf("require one arg, but got %d", len(optionLimit)))
		}
	}
}
