package xutil

import (
	"strconv"
	"strings"
)

func SubString(s string, start int, length int) string {
	runes := []rune(s)
	end := start + length

	if start < 0 {
		start = 0
	}

	if end > len(runes) {
		end = len(runes)
	}

	return string(runes[start:end])
}

func Capitalize(s string) string {
	if len(s) > 0 {
		first := strings.ToUpper(s[:1]) // 将首字母转化为大写
		return first + s[1:]
	}
	return s
}

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}
