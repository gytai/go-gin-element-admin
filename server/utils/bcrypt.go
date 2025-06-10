package utils

import (
	"crypto/md5"
	"fmt"
)

// 固定的密码盐
const PASSWORD_SALT = "go-gin-element-admin-2025"

// BcryptHash 使用 MD5 + 固定盐对密码进行加密
func BcryptHash(password string) string {
	// 用户密码 + 固定盐
	saltedPassword := password + PASSWORD_SALT
	// MD5加密并转为小写
	hash := md5.Sum([]byte(saltedPassword))
	return fmt.Sprintf("%x", hash)
}

// BcryptCheck 校验密码
func BcryptCheck(password, hash string) bool {
	// 使用相同的方式加密输入的密码
	inputHash := BcryptHash(password)
	// 比较哈希值
	return inputHash == hash
}
