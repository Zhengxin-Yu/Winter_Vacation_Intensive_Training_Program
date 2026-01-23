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

	// 静态文件
	r.Static("/uploads", "./uploads")

	// 健康检查接口：用于确认服务是否能正常响应
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	// 首页接口：功能入口
	r.GET("/home", handlers.Home)
	// 二维码展示接口（公开）
	r.GET("/qr/:code", handlers.GetQRCode)

	// /api 分组
	api := r.Group("/api")

	// public 组（无需认证）
	api.POST("/login", handlers.Login)

	// auth 组（需要登录）
	auth := api.Group("/")
	auth.Use(middleware.JWTAuth())

	// 行李模块（统一前缀 /api/luggage）
	luggage := auth.Group("/luggage")
	luggage.POST("", handlers.CreateLuggage)
	luggage.GET("/search", handlers.QueryLuggageByUserInfo)
	luggage.GET("/by_code", handlers.QueryLuggageByCode)
	luggage.GET("/by_phone", handlers.QueryLuggageByPhone)
	luggage.GET("/list", handlers.ListLuggageByUser)
	luggage.GET("/list/by_guest", handlers.ListLuggageByGuest)
	luggage.GET("/detail/by_code", handlers.GetLuggageDetailByCode)
	luggage.GET("/detail/by_phone", handlers.ListLuggageDetailByPhone)
	luggage.GET("/pickup_codes", handlers.ListPickupCodesByUser)
	luggage.GET("/pickup_codes/by_phone", handlers.ListPickupCodesByPhone)
	luggage.GET("/history", handlers.ListHistoryByGuest)
	luggage.GET("/users", handlers.ListUsersByHotel)

	luggage.GET("/storerooms", handlers.ListStorerooms)
	luggage.GET("/storerooms/:id/orders", handlers.ListLuggageByStoreroom)
	luggage.POST("/storerooms", handlers.CreateStoreroom)
	luggage.PUT("/storerooms/:id", handlers.UpdateStoreroomStatus)

	luggage.GET("/logs/stored", handlers.ListStoredLogs)
	luggage.GET("/logs/updated", handlers.ListUpdatedLogs)
	luggage.GET("/logs/retrieved", handlers.ListRetrievedLogs)

	luggage.PUT("/:id", handlers.UpdateLuggageInfo)
	luggage.PUT("/:id/code", handlers.UpdateLuggageCode)
	luggage.POST("/bind", handlers.BindLuggage)
	luggage.POST("/:id/checkout", handlers.CheckoutLuggageByCode)
	luggage.POST("/:id/transfer", handlers.TransferLuggageByID)
	luggage.GET("/:id/transfers", handlers.ListTransfersByLuggageID)
	luggage.GET("/:id", handlers.GetLuggageDetail)

	auth.POST("/upload", handlers.Upload)

	// admin 组（需要管理员权限）
	admin := auth.Group("/admin")
	admin.Use(middleware.AdminOnly())

	// 员工管理
	admin.POST("/employees", handlers.CreateEmployee)
	admin.GET("/employees", handlers.ListEmployees)
	admin.DELETE("/employees/:id", handlers.DeleteEmployee)

	// 管理员管理
	admin.POST("/admins", handlers.CreateAdmin)
	admin.GET("/admins", handlers.ListAdmins)
	admin.DELETE("/admins/:id", handlers.DeleteAdmin)

	// 酒店管理
	admin.POST("/hotels", handlers.CreateHotel)
	admin.GET("/hotels", handlers.ListHotels)
	admin.PUT("/hotels/:id", handlers.UpdateHotel)
	admin.DELETE("/hotels/:id", handlers.DeleteHotel)

	// 寄存室管理（旧路径保留）
	admin.POST("/storages", handlers.CreateStoreroom)
	admin.GET("/storages", handlers.ListStorerooms)
	admin.PUT("/storages/:id", handlers.UpdateStoreroomStatus)
	admin.DELETE("/storages/:id", handlers.DeleteStoreroom)

	return r
}
