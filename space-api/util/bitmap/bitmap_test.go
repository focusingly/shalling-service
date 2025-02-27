package bitmap_test

import (
	"fmt"
	"runtime"
	"space-api/util/bitmap"
	"testing"

	"github.com/go-playground/assert/v2"
)

func job(t *testing.T, bs *bitmap.AtomicBitmap) {
	bs.SetMask(0)
	bs.SetMask(1)
	bs.SetMask(4)
	bs.SetMask(63)

	assert.Equal(t, bs.HasMask(0), true)
	assert.Equal(t, bs.HasMask(1), true)
	assert.Equal(t, bs.HasMask(4), true)
	assert.Equal(t, bs.HasMask(63), true)
	assert.Equal(t, bs.HasMask(6), false)

	bs.ClearMask(0)
	bs.ClearMask(1)
	assert.Equal(t, bs.HasMask(0), false)
	assert.Equal(t, bs.HasMask(1), false)
	assert.Equal(t, bs.HasMask(4), true)
	assert.Equal(t, bs.HasMask(63), true)
	assert.Equal(t, bs.HasMask(6), false)

	bs.ClearAll()
	assert.Equal(t, bs.HasMask(0), false)
	assert.Equal(t, bs.HasMask(1), false)
	assert.Equal(t, bs.HasMask(4), false)
	assert.Equal(t, bs.HasMask(63), false)
	assert.Equal(t, bs.HasMask(6), false)
}

func TestBitmap(t *testing.T) {
	bs := bitmap.NewBitmap(64)
	t.Parallel()

	for i := range runtime.NumCPU() {
		t.Run(fmt.Sprintf("test-%d", i), func(t *testing.T) {
			job(t, bs)
		})
	}

}

func BenchmarkBitmap(b *testing.B) {
	bs := bitmap.NewBitmap(64)

	for b.Loop() {
		bs.SetMask(0)
		bs.SetMask(1)
		bs.SetMask(4)
		bs.SetMask(63)

		assert.Equal(b, bs.HasMask(0), true)
		assert.Equal(b, bs.HasMask(1), true)
		assert.Equal(b, bs.HasMask(4), true)
		assert.Equal(b, bs.HasMask(63), true)
		assert.Equal(b, bs.HasMask(6), false)

		bs.ClearMask(0)
		bs.ClearMask(1)
		assert.Equal(b, bs.HasMask(0), false)
		assert.Equal(b, bs.HasMask(1), false)
		assert.Equal(b, bs.HasMask(4), true)
		assert.Equal(b, bs.HasMask(63), true)
		assert.Equal(b, bs.HasMask(6), false)

		bs.ClearAll()
		assert.Equal(b, bs.HasMask(0), false)
		assert.Equal(b, bs.HasMask(1), false)
		assert.Equal(b, bs.HasMask(4), false)
		assert.Equal(b, bs.HasMask(63), false)
		assert.Equal(b, bs.HasMask(6), false)
	}
}
