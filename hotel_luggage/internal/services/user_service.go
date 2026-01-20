package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser 创建用户：
// 1. 校验参数
// 2. 密码哈希
// 3. 写入数据库
func CreateUser(username, password, role string) (models.User, error) {
	if username == "" || password == "" {
		return models.User{}, errors.New("username or password is empty")
	}
	if role == "" {
		role = "staff"
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
	}
	if err := repositories.CreateUser(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}
