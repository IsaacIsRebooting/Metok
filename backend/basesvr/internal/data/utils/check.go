package utils

import "regexp"

// IsValidWithRegex 检查字符串是否与正则表达式模式匹配。
// 参数:
//
//	pattern: 正则表达式模式
//	str: 需要检查的字符串
//
// 返回值:
//
//	如果字符串与模式匹配，返回 true；否则返回 false
func IsValidWithRegex(pattern, str string) bool {
	// 编译正则表达式模式
	regex := regexp.MustCompile(pattern)

	// 使用编译后的正则表达式检查字符串是否匹配
	matched := regex.MatchString(str)

	return matched
}
