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
