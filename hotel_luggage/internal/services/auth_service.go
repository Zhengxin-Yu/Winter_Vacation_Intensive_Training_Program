package services

import (
	"errors"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Login 校验用户名和密码，成功返回用户信息
func Login(username, password string) (models.User, error) {
	if username == "" || password == "" {
		return models.User{}, errors.New("username or password is empty")
	}

	user, err := repositories.GetUserByUsername(username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.User{}, errors.New("invalid username or password")
		}
		return models.User{}, err
	}

	// bcrypt 对比密码哈希
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return models.User{}, errors.New("invalid username or password")
	}

	return user, nil
}
