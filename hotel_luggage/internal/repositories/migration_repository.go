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
