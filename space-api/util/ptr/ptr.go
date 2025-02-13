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

// IsNil 判断一个变量是否为 nil, 而不被动态类型存在干扰
func IsNil(v interface{}) bool {
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
