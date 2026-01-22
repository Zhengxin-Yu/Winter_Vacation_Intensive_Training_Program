package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// ListHotels 查询酒店列表
func ListHotels() ([]models.Hotel, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var hotels []models.Hotel
	err := DB.Order("id ASC").Find(&hotels).Error
	return hotels, err
}

// CreateHotel 创建酒店
func CreateHotel(hotel *models.Hotel) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(hotel).Error
}

// GetHotelByID 查询酒店
func GetHotelByID(id int64) (models.Hotel, error) {
	var hotel models.Hotel
	if DB == nil {
		return hotel, errors.New("db not initialized")
	}
	err := DB.Where("id = ?", id).First(&hotel).Error
	return hotel, err
}

// UpdateHotel 更新酒店信息
func UpdateHotel(id int64, updates map[string]interface{}) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Model(&models.Hotel{}).
		Where("id = ?", id).
		Updates(updates).Error
}

// DeleteHotel 删除酒店
func DeleteHotel(id int64) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Delete(&models.Hotel{}, id).Error
}
