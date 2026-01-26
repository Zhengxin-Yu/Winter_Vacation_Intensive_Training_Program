package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateCode 生成指定长度的取件码（仅数字）
func GenerateCode(length int) (string, error) {
	const charset = "0123456789"
	if length <= 0 {
		length = 6
	}

	code := make([]byte, length)
	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		code[i] = charset[n.Int64()]
	}
	return string(code), nil
}
