package golangutils

import (
	"strings"
)

// FirstUpper 字符串首字母大写
func FirstUpper(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// FirstLower 字符串首字母小写
func FirstLower(s string) string {
	if s == "" {
		return ""
	}
	return strings.ToLower(s[:1]) + s[1:]
}

func ReplaceSpecialChars(name string) string {
	result := strings.ReplaceAll(name, " ", "")
	result = strings.ReplaceAll(result, "\uE000", "")
	result = strings.ReplaceAll(result, "-", "_")
	result = strings.ReplaceAll(result, "·", "_")
	result = strings.ReplaceAll(result, "—", "_")

	return result
}
