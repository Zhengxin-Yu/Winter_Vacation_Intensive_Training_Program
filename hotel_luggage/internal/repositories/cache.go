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
func GetLuggageByCodeCache(code string) (models.LuggageItem, bool, error) {
	if RedisClient == nil {
		return models.LuggageItem{}, false, nil
	}
	if code == "" {
		return models.LuggageItem{}, false, errors.New("code is empty")
	}
	val, err := RedisClient.Get(context.Background(), luggageByCodeKey(code)).Result()
	if err != nil {
		if err == redis.Nil {
			return models.LuggageItem{}, false, nil
		}
		return models.LuggageItem{}, false, err
	}
	var item models.LuggageItem
	if err := json.Unmarshal([]byte(val), &item); err != nil {
		return models.LuggageItem{}, false, err
	}
	return item, true, nil
}

// SetLuggageByCodeCache 写入行李信息到缓存
func SetLuggageByCodeCache(code string, item models.LuggageItem) error {
	if RedisClient == nil {
		return nil
	}
	if code == "" {
		return errors.New("code is empty")
	}
	data, err := json.Marshal(item)
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

