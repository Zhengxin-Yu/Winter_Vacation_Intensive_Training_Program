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
	luggage.GET("/by_code", handlers.QueryLuggageByCode)
	luggage.GET("/list/by_guest_name", handlers.ListStoredLuggageByGuestName)

	luggage.GET("/storerooms", handlers.ListStorerooms)
	luggage.GET("/storerooms/:id/orders", handlers.ListLuggageByStoreroom)
	luggage.POST("/storerooms", handlers.CreateStoreroom)
	luggage.PUT("/storerooms/:id", handlers.UpdateStoreroomStatus)

	luggage.GET("/logs/stored", handlers.ListStoredLogs)
	luggage.GET("/logs/updated", handlers.ListUpdatedLogs)
	luggage.GET("/logs/retrieved", handlers.ListRetrievedLogs)

	luggage.PUT("/:id", handlers.UpdateLuggageInfo)
	luggage.POST("/:id/checkout", handlers.CheckoutLuggageByCode)
	luggage.GET("/:id/checkout", handlers.GetCheckoutInfoByCode)

	return r
}
