package utils

import (
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// jwtSecret JWT 签名密钥，从环境变量 JWT_SECRET 读取
// 如果未设置环境变量，使用默认值 "change-me"（生产环境必须修改！）
var jwtSecret = []byte(getJWTSecret())

// Claims JWT 自定义载荷（Payload）
// 包含业务自定义字段（username, role）和标准字段（过期时间、签发时间等）
//
// JWT 结构说明：
// - Header：算法类型（HS256）
// - Payload：自定义数据 + 标准声明
// - Signature：签名（防止篡改）
type Claims struct {
	Username string `json:"username"` // 用户名
	Role     string `json:"role"`     // 角色（staff/admin）
	jwt.RegisteredClaims              // 标准字段：过期时间、签发时间等
}

// GenerateToken 生成 JWT token
// 参数：
//   - username: 用户名
//   - role: 用户角色（staff/admin）
// 返回：
//   - string: JWT token 字符串（用于 Authorization: Bearer <token>）
//   - error: 生成失败时返回错误
//
// Token 有效期：24小时
// 算法：HS256（HMAC-SHA256）
//
// 使用示例：
//   token, err := GenerateToken("user001", "staff")
//   // 返回示例：eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
func GenerateToken(username, role string) (string, error) {
	// 设置过期时间为当前时间 + 24小时
	expire := time.Now().Add(24 * time.Hour)
	
	// 构造载荷（Payload）
	claims := Claims{
		Username: username,
		Role:     role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expire),    // 过期时间
			IssuedAt:  jwt.NewNumericDate(time.Now()), // 签发时间
		},
	}
	
	// 使用 HS256 算法创建 token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	// 使用密钥签名并返回完整的 token 字符串
	return token.SignedString(jwtSecret)
}

// ParseToken 解析并验证 JWT token
// 参数：
//   - tokenStr: JWT token 字符串（通常从 Authorization: Bearer <token> 中提取）
// 返回：
//   - *Claims: 解析后的载荷数据（包含 username, role 等）
//   - error: 验证失败时返回错误
//
// 验证内容：
// 1. 签名是否正确（防止篡改）
// 2. Token 是否过期
// 3. Token 格式是否正确
//
// 使用示例：
//   claims, err := ParseToken("eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...")
//   if err != nil {
//       // token 无效
//   }
//   username := claims.Username
func ParseToken(tokenStr string) (*Claims, error) {
	// 解析 token，并验证签名
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 返回签名密钥用于验证
		return jwtSecret, nil
	})
	
	// 解析失败（格式错误、签名错误、已过期等）
	if err != nil {
		return nil, err
	}
	
	// 类型断言，获取自定义载荷
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	
	return claims, nil
}

// getJWTSecret 从环境变量获取 JWT 签名密钥
// 环境变量：JWT_SECRET
// 默认值：change-me（仅用于开发环境，生产环境必须设置强密钥！）
//
// 设置方式（Windows）：
//   set JWT_SECRET=your-strong-secret-key-here
//
// 设置方式（Linux/Mac）：
//   export JWT_SECRET=your-strong-secret-key-here
func getJWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		// ⚠️ 生产环境必须设置 JWT_SECRET 环境变量！
		secret = "change-me"
	}
	return secret
}
