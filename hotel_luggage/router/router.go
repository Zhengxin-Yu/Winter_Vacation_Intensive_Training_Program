package router

import (
	"hotel_luggage/internal/handlers"
	"hotel_luggage/internal/middleware"

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
	// 首页接口：功能入口
	r.GET("/home", handlers.Home)

	// 登录接口：账号密码验证
	r.POST("/login", handlers.Login)
	// 创建用户接口：生成 bcrypt 密码哈希并入库
	r.POST("/users", handlers.CreateUser)
	// 二维码展示接口（公开）
	r.GET("/qr/:code", handlers.GetQRCode)

	// 需要鉴权的接口
	auth := r.Group("/")
	auth.Use(middleware.JWTAuth())

	// 行李寄存接口：生成寄存记录与取件码
	auth.POST("/storage", handlers.CreateLuggage)
	// 查看寄存接口：按用户信息查询
	auth.GET("/storage/search", handlers.QueryLuggageByUserInfo)
	// 查看寄存接口：按取件码查询
	auth.GET("/storage/by-code", handlers.QueryLuggageByCode)
	// 取件接口：通过取件码完成取件
	auth.POST("/storage/retrieve", handlers.RetrieveLuggage)
	// 寄存单列表：按用户查询
	auth.GET("/storage/list", handlers.ListLuggageByUser)
	// 寄存单列表：按客人姓名/手机号查询
	auth.GET("/storage/list/by-guest", handlers.ListLuggageByGuest)
	// 寄存单详情
	auth.GET("/storage/detail", handlers.GetLuggageDetail)
	// 寄存单详情：按取件码查询
	auth.GET("/storage/detail/by-code", handlers.GetLuggageDetailByCode)
	// 寄存单详情：按手机号查询
	auth.GET("/storage/detail/by-phone", handlers.ListLuggageDetailByPhone)
	// 查看取件码页面：取件码列表
	auth.GET("/pickup-codes", handlers.ListPickupCodesByUser)
	// 查看取件码页面：按手机号查询
	auth.GET("/pickup-codes/by-phone", handlers.ListPickupCodesByPhone)
	// 修改寄存信息
	auth.PUT("/storage/:id", handlers.UpdateLuggageInfo)
	// 修改取件码
	auth.PUT("/storage/:id/code", handlers.UpdateLuggageCode)
	// 行李绑定到用户
	auth.POST("/storage/bind", handlers.BindLuggage)
	// 取件历史查询
	auth.GET("/storage/history/by-guest", handlers.ListHistoryByGuest)

	// 酒店管理接口
	auth.GET("/hotels", handlers.ListHotels)
	auth.POST("/hotels", handlers.CreateHotel)
	auth.PUT("/hotels/:id", handlers.UpdateHotel)
	auth.DELETE("/hotels/:id", handlers.DeleteHotel)

	// 寄存室管理接口
	auth.GET("/storerooms", handlers.ListStorerooms)
	auth.POST("/storerooms", handlers.CreateStoreroom)
	auth.DELETE("/storerooms/:id", handlers.DeleteStoreroom)
	auth.PUT("/storerooms/:id/status", handlers.UpdateStoreroomStatus)

	// 行李迁移接口
	auth.POST("/storerooms/migrate", handlers.MigrateLuggage)

	return r
}
