package ptr

import (
	"reflect"
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

func IsNil(v any) bool {
	// 获取传入值的反射类型
	val := reflect.ValueOf(v)
	// 检查传入值的类型
	if !val.IsValid() {
		return true // 如果无效值，则为 nil
	}

	// 对于指针类型，直接判断指针是否为 nil
	if val.Kind() == reflect.Ptr {
		return val.IsNil() // 检查指针是否为 nil
	}

	// 对于接口类型，检查接口的值是否为 nil
	if val.Kind() == reflect.Interface {
		return val.IsNil() // 检查接口的值是否为 nil
	}

	// 对于其他类型，直接返回 false
	return false
}

// Optional returns the value `val` if it is non-nil, otherwise it returns the `fallback` value.
// This function works with any type `T` and handles special cases for pointers, interfaces,
// slices, maps, and channels by checking if they are nil.
//
// Parameters:
//   - val: The value to check.
//   - fallback: The value to return if `val` is nil.
//
// Returns:
//   - The value `val` if it is non-nil, otherwise the `fallback` value.
func Optional[T any](val T, fallback T) T {
	// 获取值的反射值
	v := reflect.ValueOf(val)

	// 检查特殊情况：接口或指针
	if !v.IsValid() {
		return fallback
	}

	// 根据不同类型判断是否为nil
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface, reflect.Slice, reflect.Map, reflect.Chan:
		if v.IsNil() {
			return fallback
		}
	}

	return val
}
