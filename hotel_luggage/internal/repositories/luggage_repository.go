package repositories

import (
	"errors"

	"hotel_luggage/internal/models"

	"gorm.io/gorm"
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

// FindLuggageByUserInfo 按客人姓名/电话查询寄存记录
func FindLuggageByUserInfo(guestName, contactPhone string) ([]models.LuggageItem, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageItem
	query := DB.Model(&models.LuggageItem{})
	if guestName != "" {
		query = query.Where("guest_name = ?", guestName)
	}
	if contactPhone != "" {
		query = query.Where("contact_phone = ?", contactPhone)
	}
	err := query.Order("stored_at DESC").Find(&items).Error
	return items, err
}

// FindLuggageByCode 按取件码查询寄存记录
func FindLuggageByCode(code string) (models.LuggageItem, error) {
	if DB == nil {
		return models.LuggageItem{}, errors.New("db not initialized")
	}
	var item models.LuggageItem
	err := DB.Where("retrieval_code = ?", code).First(&item).Error
	return item, err
}

// GetLuggageByID 按ID查询行李记录
func GetLuggageByID(id int64) (models.LuggageItem, error) {
	if DB == nil {
		return models.LuggageItem{}, errors.New("db not initialized")
	}
	var item models.LuggageItem
	err := DB.Where("id = ?", id).First(&item).Error
	return item, err
}

// UpdateLuggageStoreroom 更新行李所在寄存室
func UpdateLuggageStoreroom(id int64, toStoreroomID int64) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Model(&models.LuggageItem{}).
		Where("id = ?", id).
		Update("storeroom_id", toStoreroomID).Error
}

// UpdateLuggageRetrieved 更新行李为已取件状态
func UpdateLuggageRetrieved(id int64, retrievedBy int64) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Model(&models.LuggageItem{}).
		Where("id = ?", id).
		Updates(map[string]interface{}{
			"status":       "retrieved",
			"retrieved_by": retrievedBy,
			"retrieved_at": gorm.Expr("NOW()"),
		}).Error
}

// ListLuggageByUser 查询某用户创建的寄存单列表
// status 可选：stored/retrieved/migrated
func ListLuggageByUser(userID int64, status string) ([]models.LuggageItem, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageItem
	query := DB.Model(&models.LuggageItem{}).Where("stored_by = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("stored_at DESC").Find(&items).Error
	return items, err
}

// ListLuggageByGuest 按客人姓名/手机号查询寄存单列表
func ListLuggageByGuest(guestName, contactPhone, status string) ([]models.LuggageItem, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageItem
	query := DB.Model(&models.LuggageItem{})
	if guestName != "" {
		query = query.Where("guest_name = ?", guestName)
	}
	if contactPhone != "" {
		query = query.Where("contact_phone = ?", contactPhone)
	}
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("stored_at DESC").Find(&items).Error
	return items, err
}

// ListPickupCodesByUser 按用户查询取件码列表（从行李表中提取）
func ListPickupCodesByUser(userID int64, status string) ([]models.LuggageItem, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageItem
	query := DB.Model(&models.LuggageItem{}).Where("stored_by = ?", userID)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("stored_at DESC").Find(&items).Error
	return items, err
}

// ListPickupCodesByPhone 按手机号查询取件码列表
func ListPickupCodesByPhone(contactPhone, status string) ([]models.LuggageItem, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageItem
	query := DB.Model(&models.LuggageItem{}).Where("contact_phone = ?", contactPhone)
	if status != "" {
		query = query.Where("status = ?", status)
	}
	err := query.Order("stored_at DESC").Find(&items).Error
	return items, err
}
