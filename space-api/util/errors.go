package util

import "runtime/debug"

type (
	// 通用业务错误标识
	BizErr struct {
		Msg    string
		Reason any
	}

	// 限流错误标识
	LimitErr struct {
		BizErr
	}

	// 认证错误标识
	AuthErr struct {
		BizErr
	}

	// 表示严重的错误级别, 但这并不表示 panic
	FatalErr struct {
		BizErr
		StackTrace []byte // 堆栈记录信息
	}

	OutboundErr struct {
		BizErr
	}
)

// Error implements error.
func (b *BizErr) Error() string {
	return b.Msg
}

func CreateBizErr(msg string, reason any) *BizErr {
	return &BizErr{
		Msg:    msg,
		Reason: reason,
	}
}
func CreateAuthErr(msg string, reason any) *AuthErr {
	return &AuthErr{
		BizErr: BizErr{
			Msg:    msg,
			Reason: reason,
		},
	}
}
func CreateLimitErr(msg string, reason any) *LimitErr {
	return &LimitErr{
		BizErr: BizErr{
			Msg:    msg,
			Reason: reason,
		},
	}
}

func CreateFatalErr(msg string, reason any) *FatalErr {
	return &FatalErr{
		BizErr: BizErr{
			Msg:    msg,
			Reason: reason,
		},
		StackTrace: debug.Stack(),
	}
}

// ensure all biz err implements default error interface
var _ error = (*BizErr)(nil)
var _ error = (*LimitErr)(nil)
var _ error = (*AuthErr)(nil)
