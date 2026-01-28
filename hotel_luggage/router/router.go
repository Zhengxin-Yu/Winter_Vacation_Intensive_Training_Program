package router

import (
	"hotel_luggage/internal/handlers"
	"hotel_luggage/internal/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRouter 初始化并配置所有路由
// 功能：
// 1. 创建 Gin 引擎（包含日志和错误恢复中间件）
// 2. 配置全局中间件（CORS、文件上传限制）
// 3. 注册公开接口（无需认证）
// 4. 注册受保护接口（需要 JWT 认证）
// 5. 返回配置完成的路由引擎
//
// 路由架构：
// - 公开接口：/api/login（登录）
// - 受保护接口：/api/luggage/... （需要 JWT token）
// - 静态文件：/uploads/... （行李照片）
// - 健康检查：/ping
//
// 返回：
//   *gin.Engine: 配置完成的路由引擎（可直接调用 Run() 启动服务）
func SetupRouter() *gin.Engine {
	// ========================================
	// 1. 创建 Gin 引擎
	// ========================================
	// gin.Default() 自动包含：
	// - Logger 中间件：记录请求日志
	// - Recovery 中间件：捕获 panic，避免服务崩溃
	r := gin.Default()
	
	// 设置文件上传大小限制：5MB（5 << 20 = 5 * 1024 * 1024）
	r.MaxMultipartMemory = 5 << 20

	// ========================================
	// 2. 配置 CORS 跨域中间件
	// ========================================
	// 用途：允许前端（可能在不同域名/端口）调用后端接口
	// 原理：浏览器发送跨域请求前会先发送 OPTIONS 预检请求
	r.Use(func(c *gin.Context) {
		// 允许所有来源跨域访问（生产环境建议指定具体域名）
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		
		// 允许携带凭证（如 Cookie、Authorization header）
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		
		// 允许的请求头（包括 Authorization，用于传递 JWT token）
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		
		// 允许的 HTTP 方法
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		// 处理 OPTIONS 预检请求（浏览器自动发送）
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204) // 返回 204 No Content
			return
		}

		// 继续处理实际请求
		c.Next()
	})

	// ========================================
	// 3. 静态文件服务
	// ========================================
	// 访问路径：http://host:port/uploads/2026/01/xxx.jpg
	// 映射到：./uploads/2026/01/xxx.jpg
	r.Static("/uploads", "./uploads")

	// ========================================
	// 4. 健康检查接口（无需认证）
	// ========================================
	// 用途：监控系统、负载均衡器可通过此接口判断服务是否正常
	// 访问：GET http://host:port/ping
	// 响应：{"message": "pong"}
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	// ========================================
	// 5. API 路由分组
	// ========================================
	// 统一前缀：/api
	api := r.Group("/api")

	// ========================================
	// 5.1 公开接口（无需认证）
	// ========================================
	api.POST("/login", handlers.Login) // 用户登录（返回 JWT token）

	// ========================================
	// 5.2 受保护接口（需要 JWT 认证）
	// ========================================
	// 所有 /api/* 路径都需要通过 JWT 认证
	auth := api.Group("/")
	auth.Use(middleware.JWTAuth()) // 应用 JWT 认证中间件

	// ========================================
	// 5.3 行李管理模块（/api/luggage）
	// ========================================
	luggage := auth.Group("/luggage")
	
	// --- 行李寄存与查询 ---
	luggage.POST("", handlers.CreateLuggage)                         // 创建行李寄存记录
	luggage.GET("/by_code", handlers.QueryLuggageByCode)            // 按取件码查询行李
	luggage.GET("/list/by_guest_name", handlers.ListStoredLuggageByGuestName) // 按客人姓名查询寄存中的行李

	// --- 寄存室管理 ---
	luggage.GET("/storerooms", handlers.ListStorerooms)             // 获取当前酒店所有寄存室
	luggage.GET("/storerooms/:id/orders", handlers.ListLuggageByStoreroom) // 获取指定寄存室的所有行李
	luggage.POST("/storerooms", handlers.CreateStoreroom)           // 创建新寄存室
	luggage.PUT("/storerooms/:id", handlers.UpdateStoreroomStatus) // 更新寄存室状态（启用/停用）

	// --- 日志查询 ---
	luggage.GET("/logs/stored", handlers.ListStoredLogs)            // 获取寄存记录（status=stored）
	luggage.GET("/logs/updated", handlers.ListUpdatedLogs)          // 获取修改记录（含寄存室迁移）
	luggage.GET("/logs/retrieved", handlers.ListRetrievedLogs)      // 获取取件记录（status=retrieved）

	// --- 行李操作 ---
	luggage.PUT("/:id", handlers.UpdateLuggageInfo)                  // 修改寄存信息（支持寄存室迁移，自动记录历史）
	luggage.POST("/:id/checkout", handlers.CheckoutLuggageByCode)   // 确认取件（更新状态、取件人、取件时间）
	luggage.GET("/:id/checkout", handlers.GetCheckoutInfoByCode)    // 获取取件信息（客人姓名、联系方式等）

	// ========================================
	// 5.4 文件上传（/api/upload）
	// ========================================
	// 用途：上传行李照片
	// 存储策略：优先 MinIO，失败则降级到本地 ./uploads 目录
	auth.POST("/upload", handlers.Upload)

	// ========================================
	// 6. 返回配置完成的路由引擎
	// ========================================
	return r
}
