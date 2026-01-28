package utils

import (
	"crypto/rand"
	"math/big"
)

// GenerateCode 生成指定长度的随机取件码（纯数字）
// 使用场景：行李寄存时生成唯一的取件码，客人凭此码取回行李
//
// 参数：
//   - length: 取件码长度（位数），如果 <= 0 则默认为 6 位
//
// 返回：
//   - string: 生成的取件码（纯数字字符串，如 "123456"）
//   - error: 生成失败时返回错误
//
// 特点：
// 1. 使用 crypto/rand 包的加密安全随机数生成器，保证随机性
// 2. 只包含数字 0-9，便于客人输入
// 3. 需要配合数据库唯一性约束，确保不重复
//
// 使用示例：
//   code, err := GenerateCode(6)  // 生成 6 位数字，如 "482917"
//   code, err := GenerateCode(8)  // 生成 8 位数字，如 "12345678"
//
// 注意：
// - 6 位数字共有 10^6 = 1,000,000 种组合
// - 8 位数字共有 10^8 = 100,000,000 种组合
// - 建议配合数据库 UNIQUE 约束防止冲突
func GenerateCode(length int) (string, error) {
	// 字符集：仅包含数字 0-9
	const charset = "0123456789"
	
	// 参数校验：如果长度无效，使用默认值 6
	if length <= 0 {
		length = 6
	}

	// 创建字节数组存储生成的取件码
	code := make([]byte, length)
	
	// 逐位生成随机数字
	for i := 0; i < length; i++ {
		// crypto/rand.Int() 生成加密安全的随机数
		// 范围：[0, len(charset))，即 [0, 10)
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			// 随机数生成失败（极少发生）
			return "", err
		}
		
		// 从字符集中选择对应位置的字符
		code[i] = charset[n.Int64()]
	}
	
	// 返回生成的取件码字符串
	return string(code), nil
}
