package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"gorm.io/gorm"
)

// ListHotels 获取酒店列表
func ListHotels() ([]models.Hotel, error) {
	return repositories.ListHotels()
}

// CreateHotel 创建酒店
func CreateHotel(name, address, phone string, isActive bool) (models.Hotel, error) {
	if name == "" {
		return models.Hotel{}, errors.New("name is empty")
	}
	hotel := models.Hotel{
		Name:     name,
		Address:  address,
		Phone:    phone,
		IsActive: isActive,
	}
	if err := repositories.CreateHotel(&hotel); err != nil {
		return models.Hotel{}, err
	}
	return hotel, nil
}

// UpdateHotel 更新酒店信息
func UpdateHotel(id int64, name, address, phone *string, isActive *bool) error {
	if id <= 0 {
		return errors.New("invalid hotel id")
	}

	_, err := repositories.GetHotelByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("hotel not found")
		}
		return err
	}

	updates := map[string]interface{}{}
	if name != nil {
		updates["name"] = *name
	}
	if address != nil {
		updates["address"] = *address
	}
	if phone != nil {
		updates["phone"] = *phone
	}
	if isActive != nil {
		updates["is_active"] = *isActive
	}
	if len(updates) == 0 {
		return errors.New("no fields to update")
	}

	return repositories.UpdateHotel(id, updates)
}

// DeleteHotel 删除酒店
func DeleteHotel(id int64) error {
	if id <= 0 {
		return errors.New("invalid hotel id")
	}
	_, err := repositories.GetHotelByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("hotel not found")
		}
		return err
	}
	return repositories.DeleteHotel(id)
}
