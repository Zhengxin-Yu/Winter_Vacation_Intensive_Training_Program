package services

import (
	"errors"
	"fmt"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"
	"hotel_luggage/utils"

	"gorm.io/gorm"
)

// CreateLuggageRequest 创建行李寄存的业务输入
type CreateLuggageRequest struct {
	GuestName    string
	ContactPhone string
	ContactEmail string
	Description  string
	Quantity     int
	SpecialNotes string
	StoreroomID  int64
	StoredBy     int64
	QRCodeURL    string
}

// CreateLuggage 生成寄存记录并自动生成取件码
func CreateLuggage(req CreateLuggageRequest) (models.LuggageItem, error) {
	if req.GuestName == "" {
		return models.LuggageItem{}, errors.New("guest name is empty")
	}
	if req.StoreroomID <= 0 {
		return models.LuggageItem{}, errors.New("invalid storeroom id")
	}
	if req.StoredBy <= 0 {
		return models.LuggageItem{}, errors.New("invalid stored_by")
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}

	// 校验寄存室是否存在且启用
	room, err := repositories.GetStoreroomByID(req.StoreroomID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageItem{}, errors.New("storeroom not found")
		}
		return models.LuggageItem{}, err
	}
	if !room.IsActive {
		return models.LuggageItem{}, errors.New("storeroom is inactive")
	}

	// 容量校验（当 capacity > 0 才判断）
	if room.Capacity > 0 {
		count, err := repositories.CountStoredByStoreroom(req.StoreroomID)
		if err != nil {
			return models.LuggageItem{}, err
		}
		if int(count) >= room.Capacity {
			return models.LuggageItem{}, errors.New("storeroom is full")
		}
	}

	// 生成唯一取件码（最多尝试 5 次）
	var code string
	for i := 0; i < 5; i++ {
		c, err := utils.GenerateCode(8)
		if err != nil {
			return models.LuggageItem{}, err
		}
		exists, err := repositories.RetrievalCodeExists(c)
		if err != nil {
			return models.LuggageItem{}, err
		}
		if !exists {
			code = c
			break
		}
	}
	if code == "" {
		return models.LuggageItem{}, fmt.Errorf("failed to generate unique retrieval code")
	}

	item := models.LuggageItem{
		GuestName:     req.GuestName,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
		Description:   req.Description,
		Quantity:      req.Quantity,
		SpecialNotes:  req.SpecialNotes,
		StoreroomID:   req.StoreroomID,
		RetrievalCode: code,
		QRCodeURL:     req.QRCodeURL,
		Status:        "stored",
		StoredBy:      req.StoredBy,
	}

	if err := repositories.CreateLuggage(&item); err != nil {
		return models.LuggageItem{}, err
	}
	return item, nil
}

// FindLuggageByUserInfo 按客人姓名/电话查询寄存记录
func FindLuggageByUserInfo(guestName, contactPhone string) ([]models.LuggageItem, error) {
	if guestName == "" && contactPhone == "" {
		return nil, errors.New("guest_name and contact_phone cannot both be empty")
	}
	return repositories.FindLuggageByUserInfo(guestName, contactPhone)
}

// FindLuggageByCode 按取件码查询寄存记录
func FindLuggageByCode(code string) (models.LuggageItem, error) {
	if code == "" {
		return models.LuggageItem{}, errors.New("code is empty")
	}
	item, err := repositories.FindLuggageByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageItem{}, errors.New("luggage not found")
		}
		return models.LuggageItem{}, err
	}
	return item, nil
}

// RetrieveLuggage 取件：根据取件码更新状态与取件人/时间
func RetrieveLuggage(code string, retrievedBy int64) (models.LuggageItem, error) {
	if code == "" {
		return models.LuggageItem{}, errors.New("code is empty")
	}
	if retrievedBy <= 0 {
		return models.LuggageItem{}, errors.New("invalid retrieved_by")
	}

	item, err := repositories.FindLuggageByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageItem{}, errors.New("luggage not found")
		}
		return models.LuggageItem{}, err
	}
	if item.Status != "stored" {
		return models.LuggageItem{}, errors.New("luggage is not in stored status")
	}

	if err := repositories.UpdateLuggageRetrieved(item.ID, retrievedBy); err != nil {
		return models.LuggageItem{}, err
	}

	// 返回最新状态（简单更新本地对象）
	item.Status = "retrieved"
	item.RetrievedBy = &retrievedBy
	return item, nil
}
