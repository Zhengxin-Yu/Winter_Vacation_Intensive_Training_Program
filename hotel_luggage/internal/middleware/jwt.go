package middleware

import (
	"net/http"
	"strings"

	"hotel_luggage/utils"

	"github.com/gin-gonic/gin"
)

// JWTAuth JWT 认证中间件
// 功能：
// 1. 从 HTTP Header 中提取 Authorization: Bearer <token>
// 2. 解析并验证 JWT token 的有效性（签名、过期时间等）
// 3. 将解析后的用户信息（username, role）存入 gin.Context
// 4. 如果 token 无效或缺失，返回 401 Unauthorized 并中止请求
//
// 使用方式：
//   auth := api.Group("/")
//   auth.Use(middleware.JWTAuth())
//
// 后续 handler 可通过以下方式获取用户信息：
//   username, _ := c.Get("username")
//   role, _ := c.Get("role")
func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. 获取 Authorization header
		auth := c.GetHeader("Authorization")
		
		// 2. 验证 header 格式是否为 "Bearer <token>"
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "missing or invalid authorization header",
			})
			c.Abort() // 中止请求，不再执行后续 handler
			return
		}

		// 3. 提取 token 字符串（去掉 "Bearer " 前缀）
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		
		// 4. 解析并验证 token
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			// token 无效（签名错误、已过期、格式错误等）
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "invalid token",
				"error":   err.Error(),
			})
			c.Abort()
			return
		}

		// 5. 将用户信息存入 Context，供后续 handler 使用
		c.Set("username", claims.Username) // 用户名
		c.Set("role", claims.Role)         // 角色（staff/admin）
		
		// 6. 继续执行后续 handler
		c.Next()
	}
}

// AdminOnly 管理员权限验证中间件
// 功能：限制只有 admin 角色的用户才能访问
// 注意：必须在 JWTAuth() 之后使用，因为需要依赖 JWTAuth 设置的 "role"
//
// 使用方式：
//   admin := api.Group("/admin")
//   admin.Use(middleware.JWTAuth(), middleware.AdminOnly())
func AdminOnly() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 从 Context 中获取角色信息（由 JWTAuth 中间件设置）
		role, _ := c.Get("role")
		
		// 验证角色是否为 admin
		if role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{
				"message": "admin only",
			})
			c.Abort() // 非管理员，中止请求
			return
		}
		
		// 管理员，继续执行
		c.Next()
	}
}
