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
var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接（失败则降级为不使用 Redis）
func InitRedis() {
	addr := os.Getenv("REDIS_ADDR")
	if addr == "" {
		addr = "127.0.0.1:6379"
	}
	password := os.Getenv("REDIS_PASSWORD")
	db := 0
	if dbStr := os.Getenv("REDIS_DB"); dbStr != "" {
		if v, err := strconv.Atoi(dbStr); err == nil {
			db = v
		}
	}

	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		log.Printf("redis disabled: %v", err)
		RedisClient = nil
		return
	}

	RedisClient = client
}
