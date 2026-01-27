package main

import (
	"log"

	"hotel_luggage/internal/repositories"
	"hotel_luggage/router"
)

// main 是程序入口：
// 1. 先初始化数据库连接（GORM）
// 2. 再初始化路由
// 3. 启动 HTTP 服务
func main() {
	// 初始化数据库连接（失败会直接退出）
	repositories.InitDB()
	// 初始化 Redis（失败则自动降级）
	repositories.InitRedis()
	// 初始化 MinIO（失败则自动降级到本地存储）
	repositories.InitMinIO()

	// 初始化 Gin 路由
	r := router.SetupRouter()

	// 启动服务，监听指定 IP 与端口
	if err := r.Run("10.154.101.161:8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
