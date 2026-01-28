package repositories

import (
	"log"

	"hotel_luggage/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 全局数据库连接对象（GORM）
// 说明：
// - 在应用启动时通过 InitDB() 初始化
// - 所有 repository 层通过 repositories.DB 访问数据库
// - 使用全局变量方便在整个应用中共享连接池
//
// 使用示例：
//   var user models.User
//   repositories.DB.Where("username = ?", "user001").First(&user)
var DB *gorm.DB

// InitDB 初始化数据库连接
// 功能：
// 1. 从环境变量读取数据库配置（DSN）
// 2. 使用 GORM 连接 MySQL 数据库
// 3. 连接失败则直接退出程序（Fatalf）
// 4. 连接成功后设置全局 DB 对象
//
// 环境变量配置：
//   DB_DSN - 数据库连接字符串（Data Source Name）
//   示例：root:password@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local
//
// 调用时机：
//   在 main() 函数启动时调用，必须在所有数据库操作之前完成
//
// 错误处理：
//   数据库连接失败会直接 Fatal 退出程序（因为没有数据库无法提供服务）
//
// 返回：
//   *gorm.DB: GORM 数据库对象（同时也会设置到全局变量 DB）
func InitDB() *gorm.DB {
	// 1. 加载数据库配置（从环境变量读取 DSN）
	cfg := configs.LoadDBConfig()
	
	// 2. 使用 GORM 连接 MySQL
	// mysql.Open(cfg.DSN)：创建 MySQL 驱动
	// &gorm.Config{}：GORM 配置（使用默认配置）
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	
	// 3. 连接失败：打印错误并退出程序
	if err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	
	// 4. 设置全局数据库对象
	DB = db
	
	// 5. 打印成功日志（可选）
	log.Println("✅ 数据库连接成功")
	
	return db
}
