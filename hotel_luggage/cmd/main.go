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

	// 初始化 Gin 路由
	r := router.SetupRouter()

	// 启动服务，默认监听 8080 端口
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
