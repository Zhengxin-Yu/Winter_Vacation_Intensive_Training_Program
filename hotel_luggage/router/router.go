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
	// 查看寄存接口：按用户信息查询
	r.GET("/storage/search", handlers.QueryLuggageByUserInfo)
	// 查看寄存接口：按取件码查询
	r.GET("/storage/by-code", handlers.QueryLuggageByCode)
	// 取件接口：通过取件码完成取件
	r.POST("/storage/retrieve", handlers.RetrieveLuggage)

	// 寄存室管理接口
	r.GET("/storerooms", handlers.ListStorerooms)
	r.POST("/storerooms", handlers.CreateStoreroom)
	r.DELETE("/storerooms/:id", handlers.DeleteStoreroom)
	r.PUT("/storerooms/:id/status", handlers.UpdateStoreroomStatus)

	// 行李迁移接口
	r.POST("/storerooms/migrate", handlers.MigrateLuggage)

	return r
}
