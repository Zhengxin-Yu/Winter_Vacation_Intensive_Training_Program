package configs

import "os"

// DBConfig 用于保存数据库连接相关配置。
// 目前只使用 DSN（Data Source Name）字符串。
type DBConfig struct {
	DSN string
}

// LoadDBConfig 从环境变量读取数据库配置。
// 优先读取 DB_DSN，如果为空则使用默认本地配置（仅用于开发环境）。
func LoadDBConfig() DBConfig {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		// 默认本地连接字符串，请根据本机 MySQL 用户名/密码调整
		// 格式：用户名:密码@tcp(地址:端口)/数据库名?参数
		dsn = "root:root@tcp(127.0.0.1:3306)/hotel_luggage?charset=utf8mb4&parseTime=True&loc=Local"
	}
	return DBConfig{DSN: dsn}
}
