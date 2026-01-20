package router

import (
	"hotel_luggage/internal/handlers"

	"github.com/gin-gonic/gin"
)

// SetupRouter 负责初始化路由配置。
// 这里是路由的统一入口，后续所有接口都应在此注册。
func SetupRouter() *gin.Engine {
	// gin.Default() 自带 Logger 和 Recovery 中间件
	r := gin.Default()

	// 健康检查接口：用于确认服务是否能正常响应
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// 登录接口：账号密码验证
	r.POST("/login", handlers.Login)

	// 创建用户接口：生成 bcrypt 密码哈希并入库
	r.POST("/users", handlers.CreateUser)

	// 行李寄存接口：生成寄存记录与取件码
	r.POST("/storage", handlers.CreateLuggage)

	return r
}
