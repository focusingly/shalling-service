package ptr

import (
	"unsafe"
)

func ToPtr[T any](val T) *T {

	return &val
}

func String2Bytes(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

func Bytes2String(b []byte) string {
	return unsafe.String(unsafe.SliceData(b), len(b))
}
