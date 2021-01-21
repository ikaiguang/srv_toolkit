package tkcrypt

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordEncrypt 加密密码
func PasswordEncrypt(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// PasswordCompare 比较密码
func PasswordCompare(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
