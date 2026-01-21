package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// GetStoreroomByID 按ID查询寄存室
func GetStoreroomByID(id int64) (models.LuggageStoreroom, error) {
	var room models.LuggageStoreroom
	if DB == nil {
		return room, errors.New("db not initialized")
	}
	err := DB.Where("id = ?", id).First(&room).Error
	return room, err
}

// ListStorerooms 查询寄存室列表
func ListStorerooms() ([]models.LuggageStoreroom, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var rooms []models.LuggageStoreroom
	err := DB.Order("id ASC").Find(&rooms).Error
	return rooms, err
}

// CreateStoreroom 创建寄存室
func CreateStoreroom(room *models.LuggageStoreroom) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(room).Error
}

// DeleteStoreroom 删除寄存室
func DeleteStoreroom(id int64) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Delete(&models.LuggageStoreroom{}, id).Error
}

// UpdateStoreroomStatus 更新寄存室启用状态
func UpdateStoreroomStatus(id int64, isActive bool) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Model(&models.LuggageStoreroom{}).
		Where("id = ?", id).
		Update("is_active", isActive).Error
}
