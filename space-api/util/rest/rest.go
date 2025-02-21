package rest

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
	Code      BizCode `json:"code" yaml:"code" xml:"code" toml:"code"`
	Timestamp int64   `json:"timestamp" yaml:"timestamp" xml:"timestamp" toml:"timestamp"`
	Msg       string  `json:"msg" yaml:"msg" xml:"msg" toml:"msg"`
	Data      T       `json:"data" yaml:"data" xml:"data" toml:"data"`
}

func RestWithSuccess[T any](data T) *RestResult[T] {
	return &RestResult[T]{
		Code:      Ok,
		Timestamp: time.Now().UnixMilli(),
		Msg:       "ok",
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
