package utils

import (
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// HashAndSalt 加密密码
func HashAndSalt(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

func ValidatePWD(hashedPwd, plainPwd string) bool {
	byteHash := []byte(hashedPwd)
	err := bcrypt.CompareHashAndPassword(byteHash, []byte(plainPwd))
	return err == nil
}

func GenerateUUID() string {
	// 使用标准库中的 uuid 包生成 UUID
	id := uuid.New()
	return id.String()
}
