package dotenv

import (
	"os"
	"strconv"
)

type envValue struct {
	v string
}

// String 转换为字符串
func (t *envValue) String(val ...string) string {
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
func (t *envValue) Bool(val ...bool) bool {
	d := false
	if len(val) >= 1 {
		d = val[0]
	}

	switch t.v {
	case "":
		return d
	case "0", "false":
		return false
	default:
		return true
	}
}

// Int64 转换为整型
func (t *envValue) Int64(val ...int64) int64 {
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
func (t *envValue) Float64(val ...float64) float64 {

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

// Getenv 获取环境变量
func Getenv(key string) *envValue {
	return &envValue{os.Getenv(key)}
}
