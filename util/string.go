package util

import (
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// GenerateFromPassword 비밀번호 생성
func GenerateFromPassword(password string) string {
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(hash)
}

// IsEmpty ...
func IsEmpty(str string) bool {
	if len(strings.TrimSpace(str)) == 0 {
		return true
	}
	return false
}
