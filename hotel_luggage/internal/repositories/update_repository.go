package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// CreateLuggageUpdate 写入寄存单修改记录
func CreateLuggageUpdate(record *models.LuggageUpdate) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(record).Error
}

// ListUpdatesByHotel 按酒店查询寄存单修改记录
func ListUpdatesByHotel(hotelID int64) ([]models.LuggageUpdate, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageUpdate
	err := DB.Where("hotel_id = ?", hotelID).Order("updated_at DESC").Find(&items).Error
	return items, err
}
