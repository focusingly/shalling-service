package util

import (
	"fmt"
	"io"
	"iter"
)

// GetSerial 根据 seed 生成函数返回可迭代序列
func GetSerial[T any](seed func() T, optionLimit ...int) iter.Seq[T] {
	prevIndex := 0
	return func(yield func(T) bool) {
		switch {
		case len(optionLimit) == 0:
			// A infinity producer
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

type writeCounterWrapper struct {
	writer io.Writer
	count  int64
}

var _ io.Writer = (*writeCounterWrapper)(nil)

func (bc *writeCounterWrapper) Write(p []byte) (n int, err error) {
	n, err = bc.writer.Write(p)
	bc.count += int64(n)
	return
}

// NewByteWriteCountWrapper 对一个写入流进行包装, 返回一个额外的字节写入统计函数
func NewByteWriteCountWrapper(writer io.Writer) (wrapperWriter io.Writer, writeCounter func() int64) {
	wr := &writeCounterWrapper{
		writer: writer,
		count:  0,
	}

	return wr, func() int64 {
		return wr.count
	}
}
