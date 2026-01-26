package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"hotel_luggage/internal/models"

	"github.com/redis/go-redis/v9"
)

const luggageByCodeTTL = time.Minute

func luggageByCodeKey(code string) string {
	return "luggage:code:" + code
}

// GetLuggageByCodeCache 从缓存中读取行李信息
func GetLuggageByCodeCache(code string) ([]models.LuggageItem, bool, error) {
	if RedisClient == nil {
		return nil, false, nil
	}
	if code == "" {
		return nil, false, errors.New("code is empty")
	}
	val, err := RedisClient.Get(context.Background(), luggageByCodeKey(code)).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, false, nil
		}
		return nil, false, err
	}
	var items []models.LuggageItem
	if err := json.Unmarshal([]byte(val), &items); err != nil {
		return nil, false, err
	}
	return items, true, nil
}

// SetLuggageByCodeCache 写入行李信息到缓存
func SetLuggageByCodeCache(code string, items []models.LuggageItem) error {
	if RedisClient == nil {
		return nil
	}
	if code == "" {
		return errors.New("code is empty")
	}
	data, err := json.Marshal(items)
	if err != nil {
		return err
	}
	return RedisClient.Set(context.Background(), luggageByCodeKey(code), data, luggageByCodeTTL).Err()
}

// DeleteLuggageByCodeCache 删除行李缓存
func DeleteLuggageByCodeCache(code string) error {
	if RedisClient == nil {
		return nil
	}
	if code == "" {
		return errors.New("code is empty")
	}
	return RedisClient.Del(context.Background(), luggageByCodeKey(code)).Err()
}

