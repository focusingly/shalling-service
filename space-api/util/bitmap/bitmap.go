package bitmap

import (
	"fmt"
	"sync/atomic"
)

type AtomicBitmap struct {
	data []uint64
}

func NewBitmap(size uint) *AtomicBitmap {
	if size == 0 {
		panic(fmt.Errorf("want a positive num, but got %d", size))
	}

	numWords := (size + 63) / 64
	return &AtomicBitmap{
		data: make([]uint64, numWords),
	}
}

func (b *AtomicBitmap) SetMask(index int) {
	wordIndex := index / 64
	bitIndex := index % 64
	atomic.OrUint64(&b.data[wordIndex], 1<<bitIndex)
}

func (b *AtomicBitmap) ClearMask(index int) {
	wordIndex := index / 64
	bitIndex := index % 64
	atomic.AndUint64(&b.data[wordIndex], ^(1 << bitIndex))
}

func (b *AtomicBitmap) HasMask(index int) bool {
	wordIndex := index / 64
	bitIndex := index % 64
	return (atomic.LoadUint64(&b.data[wordIndex]))&(1<<bitIndex) != 0
}

func (b *AtomicBitmap) ClearAll() {
	for i := range b.data {
		atomic.StoreUint64(&b.data[i], 0)
	}
}
