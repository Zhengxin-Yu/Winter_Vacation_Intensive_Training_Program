package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"gorm.io/gorm"
)

// CreateStoreroomRequest 创建寄存室的业务输入
type CreateStoreroomRequest struct {
	HotelID  int64
	Name     string
	Location string
	Capacity int
	IsActive bool
}

// ListStorerooms 获取寄存室列表（按酒店）
func ListStorerooms(hotelID int64) ([]models.LuggageStoreroom, error) {
	if hotelID <= 0 {
		return nil, errors.New("invalid hotel id")
	}
	return repositories.ListStorerooms(hotelID)
}

// CreateStoreroom 创建寄存室
func CreateStoreroom(req CreateStoreroomRequest) (models.LuggageStoreroom, error) {
	if req.HotelID <= 0 {
		return models.LuggageStoreroom{}, errors.New("invalid hotel id")
	}
	if req.Name == "" {
		return models.LuggageStoreroom{}, errors.New("name is empty")
	}
	if req.Capacity < 0 {
		return models.LuggageStoreroom{}, errors.New("capacity cannot be negative")
	}

	// 校验酒店是否存在
	if _, err := repositories.GetHotelByID(req.HotelID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageStoreroom{}, errors.New("hotel not found")
		}
		return models.LuggageStoreroom{}, err
	}

	room := models.LuggageStoreroom{
		HotelID:  req.HotelID,
		Name:     req.Name,
		Location: req.Location,
		Capacity: req.Capacity,
		IsActive: req.IsActive,
	}
	if err := repositories.CreateStoreroom(&room); err != nil {
		return models.LuggageStoreroom{}, err
	}
	return room, nil
}

// DeleteStoreroom 删除寄存室（有行李则禁止删除）
func DeleteStoreroom(id int64) error {
	if id <= 0 {
		return errors.New("invalid storeroom id")
	}

	// 判断是否存在
	_, err := repositories.GetStoreroomByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("storeroom not found")
		}
		return err
	}

	// 如果寄存室内还有行李（stored），禁止删除
	count, err := repositories.CountStoredByStoreroom(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("storeroom has luggage, cannot delete")
	}

	return repositories.DeleteStoreroom(id)
}

// UpdateStoreroomStatus 更新寄存室状态（启用/停用）
func UpdateStoreroomStatus(id int64, isActive bool) error {
	if id <= 0 {
		return errors.New("invalid storeroom id")
	}

	_, err := repositories.GetStoreroomByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("storeroom not found")
		}
		return err
	}

	return repositories.UpdateStoreroomStatus(id, isActive)
}
