package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"gorm.io/gorm"
)

// MigrateLuggageRequest 行李迁移的业务输入
type MigrateLuggageRequest struct {
	LuggageID     int64
	ToStoreroomID int64
	MigratedBy    int64
}

// MigrateLuggage 行李迁移：更新寄存室并记录迁移日志
func MigrateLuggage(req MigrateLuggageRequest) (models.LuggageMigration, error) {
	if req.LuggageID <= 0 || req.ToStoreroomID <= 0 || req.MigratedBy <= 0 {
		return models.LuggageMigration{}, errors.New("invalid request")
	}

	// 1) 查询行李记录
	item, err := repositories.GetLuggageByID(req.LuggageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageMigration{}, errors.New("luggage not found")
		}
		return models.LuggageMigration{}, err
	}
	if item.Status != "stored" {
		return models.LuggageMigration{}, errors.New("luggage is not in stored status")
	}

	// 2) 校验目标寄存室
	room, err := repositories.GetStoreroomByID(req.ToStoreroomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageMigration{}, errors.New("target storeroom not found")
		}
		return models.LuggageMigration{}, err
	}
	if item.HotelID != room.HotelID {
		return models.LuggageMigration{}, errors.New("storeroom hotel mismatch")
	}
	if !room.IsActive {
		return models.LuggageMigration{}, errors.New("target storeroom is inactive")
	}
	if room.Capacity > 0 {
		count, err := repositories.CountStoredByStoreroom(req.ToStoreroomID)
		if err != nil {
			return models.LuggageMigration{}, err
		}
		if int(count) >= room.Capacity {
			return models.LuggageMigration{}, errors.New("target storeroom is full")
		}
	}

	// 3) 更新行李寄存室
	if err := repositories.UpdateLuggageStoreroom(item.ID, req.ToStoreroomID); err != nil {
		return models.LuggageMigration{}, err
	}

	// 4) 写入迁移日志
	log := models.LuggageMigration{
		HotelID:         room.HotelID,
		LuggageID:       item.ID,
		FromStoreroomID: item.StoreroomID,
		ToStoreroomID:   req.ToStoreroomID,
		MigratedBy:      req.MigratedBy,
	}
	if err := repositories.CreateMigrationLog(&log); err != nil {
		return models.LuggageMigration{}, err
	}

	return log, nil
}
