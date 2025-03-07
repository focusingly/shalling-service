package arr_test

import (
	"space-api/util/arr"
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestArrFunc(t *testing.T) {
	type k struct {
		val int
	}
	arr1 := []k{
		{1},
		{1},
		{2},
		{3},
		{2},
		{0},
		{0},
	}

	compact1 := arr.Compress(arr1, func(v1, v2 k) bool {
		return v1.val == v2.val
	})

	assert.Equal(t, compact1, []k{
		{1},
		{2},
		{3},
		{0},
	})
}
