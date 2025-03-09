package query

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Int64Array 是一个自定义类型，用于将 int64 数组序列化为字符串数组
type Int64Array []int64

// MarshalJSON 实现 json.Marshaler 接口
func (a Int64Array) MarshalJSON() ([]byte, error) {
	if len(a) == 0 {
		return []byte("[]"), nil
	}

	var builder strings.Builder
	builder.WriteString("[")

	for i, v := range a {
		if i > 0 {
			builder.WriteString(",")
		}
		builder.WriteString(fmt.Sprintf("\"%d\"", v))
	}

	builder.WriteString("]")
	return []byte(builder.String()), nil
}

// UnmarshalJSON 实现 json.Unmarshaler 接口
func (a *Int64Array) UnmarshalJSON(data []byte) error {
	// 解析 JSON 数组
	var stringArray []string
	if err := json.Unmarshal(data, &stringArray); err != nil {
		// 尝试解析为 int64 数组，处理客户端可能发送的整数数组
		var int64Array []int64
		if err := json.Unmarshal(data, &int64Array); err != nil {
			return err
		}
		*a = int64Array
		return nil
	}

	// 将字符串数组转换为 int64 数组
	result := make([]int64, len(stringArray))
	for i, s := range stringArray {
		n, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return err
		}
		result[i] = n
	}

	*a = result
	return nil
}
