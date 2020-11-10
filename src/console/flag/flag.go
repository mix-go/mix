package flag

import (
	"fmt"
	"strconv"
)

// 参数值
type flagValue struct {
	v     string
	exist bool
}

// String 转换为字符串
func (t *flagValue) String(val ...string) string {
	d := ""
	if len(val) >= 1 {
		d = val[0]
	}

	if t.v == "" {
		return d
	}

	return t.v
}

// Bool 转换为布尔
func (t *flagValue) Bool(val ...bool) bool {
	d := false
	if len(val) >= 1 {
		d = val[0]
	}

	if !t.exist {
		return d
	}

	switch t.v {
	case "false":
		return false
	default:
		return true
	}
}

// Int64 转换为整型
func (t *flagValue) Int64(val ...int64) int64 {
	d := int64(0)
	if len(val) >= 1 {
		d = val[0]
	}

	if t.v == "" {
		return d
	}

	v, _ := strconv.ParseInt(t.v, 10, 64)
	return v
}

// Float64 转换为浮点
func (t *flagValue) Float64(val ...float64) float64 {

	d := float64(0)
	if len(val) >= 1 {
		d = val[0]
	}

	if t.v == "" {
		return d
	}

	v, _ := strconv.ParseFloat(t.v, 64)
	return v
}

// Match 匹配参数名称
func Match(names ...string) *flagValue {
	for _, name := range names {
		v, exist := value(name)
		if exist {
			return &flagValue{v, exist}
		}
	}
	return &flagValue{}
}

// 获取指定参数的值
func value(name string) (string, bool) {
	key := ""
	if len(name) == 1 {
		key = fmt.Sprintf("-%s", name)
	} else {
		key = fmt.Sprintf("--%s", name)
	}
	if v, ok := Options()[key]; ok {
		return v, true
	}
	return "", false
}
