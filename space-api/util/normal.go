package util

import (
	"fmt"
	"io"
	"iter"
)

// GetSerial 生成一个序列生成器函数，该函数根据提供的种子函数和可选的限制参数生成序列。
//
// Parameters:
//   - seed: 一个无参函数，返回类型为 T 的值，用于生成序列中的每个元素。
//   - optionLimit: 可选参数，指定生成序列的最大长度。如果未提供，则生成无限序列。
//
// Returns:
//   - iter.Seq[T]: 一个序列生成器函数，接受一个 yield 函数作为参数，该函数用于处理生成的每个元素。
//
// Usages:
//   - 如果未提供 optionLimit 参数，则生成一个无限序列，直到 yield 函数返回 false。
//   - 如果提供了一个正整数的 optionLimit 参数，则生成一个长度不超过该限制的序列。
//   - 如果提供的 optionLimit 参数为非正整数或超过一个参数，则会引发 panic。
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

// TernaryExpr 三元表达式的替换
//
// Parameters:
//
//	boolVal 条件值
//	matched 条件匹配成功的返回值产生器
//	fallback 条件匹配失败的返回值产生器
//
// Notes:
//   - 对于判断空指针是要尤其注意, 如下示例的用法是错误的(这会导致严重的空指针 panic, 因为无论是否为空, 都在传参时执行了指针取值操作),
//     ```
//     var i *int;
//     TernaryExpr(i !=nil, *i, 0)
//     ```
//   - 请考虑使用 `TernaryExprWithProducer` 替换
//
// Returns:
//   - T
func TernaryExpr[T any](boolVal bool, success, fallback T) T {
	if boolVal {
		return success
	}

	return fallback
}

func TernaryExprWithProducer[T any](boolVal bool, success, fallback func() T) T {
	if boolVal {
		return success()
	}

	return fallback()
}

// GetWithFallback attempts to execute the provided producer function and returns its result.
// If the producer function returns an error, the fallback value is returned instead.
//
// Type Parameters:
//
//	T - the type of the value produced by the producer function and the fallback value.
//
// Parameters:
//
//	producer - a function that produces a value of type T and may return an error.
//	fallback - a value of type T to be returned if the producer function fails.
//
// Returns:
//
//	The value produced by the producer function if it succeeds, otherwise the fallback value.
func GetWithFallback[T any](producer func() (T, error), fallback T) T {
	if p, e := producer(); e != nil {
		return fallback
	} else {
		return p
	}
}
