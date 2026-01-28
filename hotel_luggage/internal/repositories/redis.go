package repositories

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 全局 Redis 客户端（可为空）
// 说明：
// - 用于缓存热点数据（如按取件码查询的结果）
// - 如果 Redis 连接失败，RedisClient 为 nil，系统会自动降级到直接查询数据库
// - 使用前需要判断 RedisClient != nil
//
// 设计理念：Redis 是可选的性能优化组件，不影响核心功能
//
// 使用示例：
//   if RedisClient != nil {
//       val, err := RedisClient.Get(ctx, "key").Result()
//   }
var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接（失败则自动降级）
// 功能：
// 1. 从环境变量读取 Redis 配置
// 2. 创建 Redis 客户端
// 3. 测试连接（Ping）
// 4. 连接失败：打印警告日志，RedisClient 设为 nil（降级）
// 5. 连接成功：设置全局 RedisClient 对象
//
// 环境变量配置：
//   REDIS_ADDR     - Redis 地址（默认：127.0.0.1:6379）
//   REDIS_PASSWORD - Redis 密码（默认：空，无密码）
//   REDIS_DB       - Redis 数据库编号（默认：0）
//
// 设置示例（Windows）：
//   set REDIS_ADDR=127.0.0.1:6379
//   set REDIS_PASSWORD=
//   set REDIS_DB=0
//
// 降级策略：
//   Redis 连接失败不会导致程序退出，而是打印警告日志并继续运行
//   业务代码需要判断 RedisClient != nil，失败时直接查询数据库
//
// 性能优化：
//   - 设置 2 秒连接超时，避免启动时长时间等待
//   - 连接成功后 Redis 可缓存热点数据，减少数据库查询压力
func InitRedis() {
	// 1. 读取 Redis 地址（默认 localhost:6379）
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	
	// 2. 读取 Redis 密码（默认无密码）
	password := os.Getenv("REDIS_PASSWORD")
	
	// 3. 读取 Redis 数据库编号（默认 DB 0）
	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if v, err := strconv.Atoi(dbStr); err == nil {
			db = v
		}
	}

	// 4. 创建 Redis 客户端
	client := redis.NewClient(&redis.Options{
		Addr:     addr,     // Redis 服务器地址
		Password: password, // 密码（如果有）
		DB:       db,       // 数据库编号
	})

	// 5. 测试连接（Ping，2秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	
	if err := client.Ping(ctx).Err(); err != nil {
		// 连接失败：打印警告，降级为不使用 Redis
		log.Printf("⚠️  Redis 连接失败，系统将降级使用数据库: %v", err)
		RedisClient = nil
		return
	}

	// 6. 连接成功：设置全局客户端
	RedisClient = client
	log.Println("✅ Redis 初始化成功")
}
