package util

import "time"

type BizCode int

const (
	Cached BizCode = (iota+1)*1000 + 1
	Ok
	Limit
	NotAllowed
	Error
)

type RestResult[T any] struct {
	Code      BizCode `json:"code"`
	Timestamp int64   `json:"timestamp"`
	Msg       string  `json:"msg"`
	Data      T       `json:"data"`
}

func RestWithSuccess[T any](data T) *RestResult[T] {
	return &RestResult[T]{
		Code:      Ok,
		Timestamp: time.Now().UnixMilli(),
		Msg:       "Ok",
		Data:      data,
	}
}

func RestWithError(msg string) *RestResult[any] {
	return &RestResult[any]{
		Code:      NotAllowed,
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
	}
}

func RestWithUnknown(msg string) *RestResult[any] {
	return &RestResult[any]{
		Code:      Error,
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
	}
}

func RestWithNotAllowed(msg string) *RestResult[any] {
	return &RestResult[any]{
		Code:      Error,
		Timestamp: time.Now().UnixMilli(),
		Msg:       msg,
	}
}

// TernaryExpr go 三元表达式的替换
func TernaryExpr[T any](condition bool, val1, val2 T) T {
	if condition {
		return val1
	}

	return val2
}
