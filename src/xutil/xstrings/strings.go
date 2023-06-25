package xstrings

import "strconv"

func IsNumeric(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

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
