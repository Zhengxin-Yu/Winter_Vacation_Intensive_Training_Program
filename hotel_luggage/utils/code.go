package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateCode 生成指定长度的取件码（由大写字母和数字组成）
func GenerateCode(length int) (string, error) {
	const charset = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789" // 去除易混淆字符
	if length <= 0 {
		length = 8
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
