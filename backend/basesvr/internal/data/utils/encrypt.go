package utils

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"io"
)

// GetPasswordSalt 生成一个随机的密码盐值。
// 密码盐值是一串随机数据，通常用于密码哈希过程中，以增加攻击者破解密码的难度。
// 该函数没有输入参数。
// 返回值:
//   - string: 生成的密码盐值的十六进制表示。
//   - error: 错误信息，如果生成随机盐值过程中出现错误，则返回该错误。
func GetPasswordSalt() (string, error) {
	// 创建一个长度为16的字节切片，用于存储生成的随机盐值。
	b := make([]byte, 16)

	// 使用rand.Reader生成随机字节，填充切片b。
	// rand.Reader是Go标准库提供的一个全局随机数生成器。
	// io.ReadFull确保切片b被完全填满随机字节，否则返回错误。
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		// 如果生成随机字节过程中出现错误，返回空字符串和错误信息。
		return "", err
	}

	// 将生成的随机盐值（字节切片b）编码为十六进制字符串，并返回。
	// hex.EncodeToString将字节切片转换为十六进制表示的字符串。
	return hex.EncodeToString(b), nil
}

// GenerateMd5WithSalt 对给定的密码和盐值进行组合并生成MD5哈希值。
// 这个函数的目的是增加密码存储的安全性，通过使用盐值来防止彩虹表攻击。
// 参数:
//
//	password - 需要哈希处理的原始密码字符串。
//	salt - 用于增加密码复杂度的盐值字符串，通常是一个随机字符串。
//
// 返回值:
//
//	返回经过MD5哈希处理后的密码字符串，该字符串是固定长度的32位十六进制字符。
func GenerateMd5WithSalt(password, salt string) string {
	// 将盐值附加到密码末尾，以增加密码的复杂性。
	password += salt

	// 使用md5.Sum函数计算附加了盐值的密码的MD5哈希值。
	hash := md5.Sum([]byte(password))

	// 将得到的哈希值转换为十六进制字符串表示，并返回。
	// 这里使用hex.EncodeToString函数是因为MD5哈希结果是字节切片，而我们通常需要字符串形式的哈希值。
	return hex.EncodeToString(hash[:])
}
