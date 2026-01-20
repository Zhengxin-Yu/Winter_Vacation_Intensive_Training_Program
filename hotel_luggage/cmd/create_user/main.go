package main

import (
	"flag"
	"fmt"
	"log"

	"hotel_luggage/internal/repositories"
	"hotel_luggage/internal/services"
)

// 命令行工具：创建用户并生成 bcrypt 密码哈希
// 用法示例：
// go run ./cmd/create_user -u admin -p 123456 -r admin
func main() {
	// 读取命令行参数
	username := flag.String("u", "", "用户名")
	password := flag.String("p", "", "密码（明文）")
	role := flag.String("r", "staff", "角色（staff/admin）")
	flag.Parse()

	if *username == "" || *password == "" {
		log.Fatal("参数缺失：必须提供 -u 和 -p")
	}

	// 初始化数据库连接
	repositories.InitDB()

	// 创建用户（自动生成 bcrypt 哈希）
	user, err := services.CreateUser(*username, *password, *role)
	if err != nil {
		log.Fatalf("创建用户失败: %v", err)
	}

	fmt.Printf("创建成功：id=%d, username=%s, role=%s\n", user.ID, user.Username, user.Role)
}
