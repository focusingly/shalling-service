package util_test

import (
	"fmt"
	"slices"
	"testing"
)

func TestSort(*testing.T) {
	s := []int{1, 3, 1, 6, 2}
	slices.SortFunc(s, func(a, b int) int {
		return b - a
	})

	fmt.Println(s)
}
