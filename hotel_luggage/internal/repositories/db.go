package repositories

import (
	"log"

	"hotel_luggage/configs"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// DB 是全局数据库连接对象，方便在各层复用。
// 业务代码中通过 repositories.DB 直接使用即可。
var DB *gorm.DB

// InitDB 初始化数据库连接：
// 1. 读取配置
// 2. 通过 GORM 连接 MySQL
// 3. 失败则直接退出程序
func InitDB() *gorm.DB {
	cfg := configs.LoadDBConfig()
	db, err := gorm.Open(mysql.Open(cfg.DSN), &gorm.Config{})
	if err != nil {
		log.Fatalf("database connect failed: %v", err)
	}
	DB = db
	return db
}
