// Package crypto 提供加密相关的工具函数。
package crypto

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword 对密码进行加密。
// 参数 password 为明文密码。
// 返回加密后的密码哈希值和可能的错误。
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPassword 验证密码是否匹配。
// 参数 password 为明文密码，hash 为加密后的密码哈希值。
// 返回 true 表示密码匹配，false 表示不匹配。
func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
