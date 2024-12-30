package constants

// AccountPasswordPattern 是一个常量，定义了账户密码的正则表达式模式。
// 该模式要求密码至少包含8个字符，并且可以包含大小写字母、数字和特殊字符。
const (
	AccountPasswordPattern = "^[A-Za-z\\d\\S]{8,}"
)
