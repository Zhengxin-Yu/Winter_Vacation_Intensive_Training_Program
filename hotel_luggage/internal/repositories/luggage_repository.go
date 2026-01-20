package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// CreateLuggage 创建行李寄存记录
func CreateLuggage(item *models.LuggageItem) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(item).Error
}

// RetrievalCodeExists 判断取件码是否已存在
func RetrievalCodeExists(code string) (bool, error) {
	if DB == nil {
		return false, errors.New("db not initialized")
	}
	var count int64
	err := DB.Model(&models.LuggageItem{}).Where("retrieval_code = ?", code).Count(&count).Error
	return count > 0, err
}

// CountStoredByStoreroom 统计某寄存室内“已存放”的行李数量
func CountStoredByStoreroom(storeroomID int64) (int64, error) {
	if DB == nil {
		return 0, errors.New("db not initialized")
	}
	var count int64
	err := DB.Model(&models.LuggageItem{}).
		Where("storeroom_id = ? AND status = ?", storeroomID, "stored").
		Count(&count).Error
	return count, err
}
