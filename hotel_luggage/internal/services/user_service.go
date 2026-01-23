package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// CreateUser 创建用户：
// 1. 校验参数
// 2. 密码哈希
// 3. 写入数据库
func CreateUser(username, password, role string, hotelID *int64) (models.User, error) {
	if username == "" || password == "" {
		return models.User{}, errors.New("username or password is empty")
	}
	if role == "" {
		role = "staff"
	}
	if role != "staff" && role != "admin" {
		return models.User{}, errors.New("invalid role")
	}

	// staff/admin 必须关联酒店
	if hotelID == nil || *hotelID <= 0 {
		return models.User{}, errors.New("hotel_id is required for staff/admin")
	}

	if _, err := repositories.GetHotelByID(*hotelID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("hotel not found")
		}
		return models.User{}, err
	}

	// 用户名唯一校验
	if _, err := repositories.GetUserByUsername(username); err == nil {
		return models.User{}, errors.New("username already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return models.User{}, err
	}

	// 生成 bcrypt 哈希
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return models.User{}, err
	}

	user := models.User{
		Username:     username,
		PasswordHash: string(hash),
		Role:         role,
		HotelID:      hotelID,
	}
	if err := repositories.CreateUser(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

// ListUsersByRole 查询用户列表
func ListUsersByRole(role string) ([]models.User, error) {
	if role != "staff" && role != "admin" {
		return nil, errors.New("invalid role")
	}
	return repositories.ListUsersByRole(role)
}

// DeleteUserByRole 删除指定角色的用户
func DeleteUserByRole(id int64, role string) error {
	if id <= 0 {
		return errors.New("invalid id")
	}
	user, err := repositories.GetUserByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user not found")
		}
		return err
	}
	if user.Role != role {
		return errors.New("role mismatch")
	}
	return repositories.DeleteUserByID(id)
}
