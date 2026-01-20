package repositories

import (
	"errors"

	"hotel_luggage/internal/models"
)

// GetUserByUsername 按用户名查询用户，找不到则返回 gorm.ErrRecordNotFound
func GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if DB == nil {
		return user, errors.New("db not initialized")
	}
	err := DB.Where("username = ?", username).First(&user).Error
	return user, err
}

// CreateUser 创建用户记录（写入数据库）
func CreateUser(user *models.User) error {
	if DB == nil {
		return errors.New("db not initialized")
	}
	return DB.Create(user).Error
}
