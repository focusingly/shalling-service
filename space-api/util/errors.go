package util

type (
	BizErr struct {
		Msg    string
		Reason any
	}

	LimitErr struct {
		BizErr
	}

	VerifyErr struct {
		BizErr
	}

	InnerErr struct {
		BizErr
		StackTrace []byte // 堆栈记录信息
	}
)

// Error implements error.
func (b *BizErr) Error() string {
	return b.Msg
}

// ensure all biz err implements default error interface
var _ error = (*BizErr)(nil)
var _ error = (*LimitErr)(nil)
var _ error = (*VerifyErr)(nil)
