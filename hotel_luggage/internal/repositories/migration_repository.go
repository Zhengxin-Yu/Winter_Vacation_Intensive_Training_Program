package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// CreateMigrationLog 创建行李迁移日志
func CreateMigrationLog(log *models.LuggageMigration) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(log).Error
}

// ListMigrationsByLuggageID 查询迁移历史
func ListMigrationsByLuggageID(luggageID int64) ([]models.LuggageMigration, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageMigration
	err := DB.Where("luggage_id = ?", luggageID).Order("migrated_at DESC").Find(&items).Error
	return items, err
}

// ListMigrationsByHotel 查询某酒店的迁移日志
func ListMigrationsByHotel(hotelID int64) ([]models.LuggageMigration, error) {
	if DB == nil {
		return nil, errors.New("db not initialized")
	}
	var items []models.LuggageMigration
	err := DB.Where("hotel_id = ?", hotelID).Order("migrated_at DESC").Find(&items).Error
	return items, err
}
