package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// CreateLuggageHistory 写入取件历史记录
func CreateLuggageHistory(record *models.LuggageHistory) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(record).Error
}

// ListHistoryByGuest 按客人姓名/手机号查询取件历史
func ListHistoryByGuest(guestName, contactPhone string) ([]models.LuggageHistory, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageHistory
	query := DB.Model(&models.LuggageHistory{})
	if guestName != "" {
		query = query.Where("guest_name = ?", guestName)
	}
	if contactPhone != "" {
		query = query.Where("contact_phone = ?", contactPhone)
	}
	err := query.Order("retrieved_at DESC").Find(&items).Error
	return items, err
}

// ListHistoryByHotel 按酒店查询取件历史（可选客人姓名/手机号）
func ListHistoryByHotel(hotelID int64, guestName, contactPhone string) ([]models.LuggageHistory, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageHistory
	query := DB.Model(&models.LuggageHistory{}).Where("hotel_id = ?", hotelID)
	if guestName != "" {
		query = query.Where("guest_name = ?", guestName)
	}
	if contactPhone != "" {
		query = query.Where("contact_phone = ?", contactPhone)
	}
	err := query.Order("retrieved_at DESC").Find(&items).Error
	return items, err
}
