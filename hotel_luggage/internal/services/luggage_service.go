package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"hotel_luggage/internal/models"
	"hotel_luggage/internal/repositories"
	"hotel_luggage/utils"

	"gorm.io/gorm"
)

// CreateLuggageRequest 创建行李寄存的业务输入
type CreateLuggageRequest struct {
	GuestName     string
	ContactPhone  string
	ContactEmail  string
	Description   string
	Quantity      int
	SpecialNotes  string
	PhotoURL      string
	PhotoURLs     []string
	RetrievalCode string
	StoreroomID   int64
	StaffName     string
	QRCodeURL     string
}

// CreateLuggage 生成寄存记录并自动生成取件码
func CreateLuggage(req CreateLuggageRequest) (models.LuggageItem, error) {
	if req.GuestName == "" {
		return models.LuggageItem{}, errors.New("guest name is empty")
	}
	if req.StoreroomID <= 0 {
		return models.LuggageItem{}, errors.New("invalid storeroom id")
	}
	if req.StaffName == "" {
		return models.LuggageItem{}, errors.New("staff_name is empty")
	}
	if req.Quantity <= 0 {
		req.Quantity = 1
	}
	if len(req.PhotoURLs) == 0 && req.PhotoURL != "" {
		req.PhotoURLs = []string{req.PhotoURL}
	}
	if len(req.PhotoURLs) > 0 && req.PhotoURL == "" {
		req.PhotoURL = req.PhotoURLs[0]
	}

	staff, err := repositories.GetUserByUsername(req.StaffName)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageItem{}, errors.New("staff not found")
		}
		return models.LuggageItem{}, err
	}
	if staff.Role != "staff" {
		return models.LuggageItem{}, errors.New("staff_name is not staff")
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
	if room.HotelID <= 0 {
		return models.LuggageItem{}, errors.New("storeroom hotel_id is missing")
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

	// 生成或使用取件码（6 位数字，最多尝试 5 次）
	code := req.RetrievalCode
	if code == "" {
		for i := 0; i < 5; i++ {
			c, err := utils.GenerateCode(6)
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
	}

	item := models.LuggageItem{
		GuestName:     req.GuestName,
		ContactPhone:  req.ContactPhone,
		ContactEmail:  req.ContactEmail,
		Description:   req.Description,
		Quantity:      req.Quantity,
		SpecialNotes:  req.SpecialNotes,
		PhotoURL:      req.PhotoURL,
		PhotoURLs:     req.PhotoURLs,
		HotelID:       room.HotelID,
		StoreroomID:   req.StoreroomID,
		RetrievalCode: code,
		QRCodeURL:     req.QRCodeURL,
		Status:        "stored",
		StoredBy:      req.StaffName,
	}

	// 如果未传入二维码URL，则默认指向二维码展示接口
	if item.QRCodeURL == "" {
		item.QRCodeURL = fmt.Sprintf("/qr/%s", code)
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
func FindLuggageByCode(code string) ([]models.LuggageItem, error) {
	if code == "" {
		return nil, errors.New("code is empty")
	}
	if items, ok, err := repositories.GetLuggageByCodeCache(code); err == nil && ok {
		return items, nil
	}
	items, err := repositories.FindLuggageByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("luggage not found")
		}
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("luggage not found")
	}
	_ = repositories.SetLuggageByCodeCache(code, items)
	return items, nil
}

// RetrieveLuggage 取件：根据取件码更新状态与取件人/时间
func RetrieveLuggage(code string, retrievedByUsername string) ([]models.LuggageItem, error) {
	if code == "" {
		return nil, errors.New("code is empty")
	}
	if retrievedByUsername == "" {
		return nil, errors.New("retrieved_by is empty")
	}

	user, err := repositories.GetUserByUsername(retrievedByUsername)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("user not found")
		}
		return nil, err
	}
	if user.Role != "staff" {
		return nil, errors.New("retrieved_by is not staff")
	}

	items, err := repositories.FindLuggageByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("luggage not found")
		}
		return nil, err
	}
	if len(items) == 0 {
		return nil, errors.New("luggage not found")
	}
	storedItems := make([]models.LuggageItem, 0, len(items))
	for _, item := range items {
		if item.Status == "stored" {
			storedItems = append(storedItems, item)
		}
	}
	if len(storedItems) == 0 {
		return nil, errors.New("luggage is not in stored status")
	}

	for _, item := range storedItems {
		if err := repositories.UpdateLuggageRetrieved(item.ID, user.Username); err != nil {
			return nil, err
		}
		// 写入取件历史
		history := models.LuggageHistory{
			LuggageID:     item.ID,
			GuestName:     item.GuestName,
			ContactPhone:  item.ContactPhone,
			ContactEmail:  item.ContactEmail,
			Description:   item.Description,
			Quantity:      item.Quantity,
			SpecialNotes:  item.SpecialNotes,
			PhotoURL:      item.PhotoURL,
			PhotoURLs:     item.PhotoURLs,
			HotelID:       item.HotelID,
			StoreroomID:   item.StoreroomID,
			RetrievalCode: item.RetrievalCode,
			QRCodeURL:     item.QRCodeURL,
			Status:        "retrieved",
			StoredBy:      item.StoredBy,
			RetrievedBy:   user.Username,
			StoredAt:      item.StoredAt,
			RetrievedAt:   time.Now(),
		}
		if err := repositories.CreateLuggageHistory(&history); err != nil {
			return nil, err
		}
		// 已取件的行李从数据库中删除（历史已保留）
		if err := repositories.DeleteLuggageByID(item.ID); err != nil {
			return nil, err
		}
	}
	_ = repositories.DeleteLuggageByCodeCache(code)

	return storedItems, nil
}

// ListLuggageByUser 获取用户寄存单列表
func ListLuggageByUser(username string, status string) ([]models.LuggageItem, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	return repositories.ListLuggageByUser(username, status)
}

// ListLuggageByGuest 按客人姓名/手机号查询寄存单列表
func ListLuggageByGuest(guestName, contactPhone, status string) ([]models.LuggageItem, error) {
	if guestName == "" && contactPhone == "" {
		return nil, errors.New("guest_name and contact_phone cannot both be empty")
	}
	return repositories.ListLuggageByGuest(guestName, contactPhone, status)
}

// ListGuestNames 获取所有寄存客人姓名（去重）
func ListGuestNames() ([]string, error) {
	return repositories.ListGuestNames()
}

// ListLuggageByStoreroom 按寄存室查询寄存单列表
func ListLuggageByStoreroom(storeroomID int64, status string) ([]models.LuggageItem, error) {
	if storeroomID <= 0 {
		return nil, errors.New("invalid storeroom id")
	}
	return repositories.ListLuggageByStoreroom(storeroomID, status)
}

// ListLuggageByHotelAndStatus 按酒店与状态查询寄存单列表
func ListLuggageByHotelAndStatus(hotelID int64, status string) ([]models.LuggageItem, error) {
	if hotelID <= 0 {
		return nil, errors.New("invalid hotel id")
	}
	return repositories.ListLuggageByHotelAndStatus(hotelID, status)
}

// ListGuestNamesByHotelAndStatus 查询某酒店下指定状态的客人姓名（去重）
func ListGuestNamesByHotelAndStatus(hotelID int64, status string) ([]string, error) {
	if hotelID <= 0 {
		return nil, errors.New("invalid hotel id")
	}
	return repositories.ListGuestNamesByHotelAndStatus(hotelID, status)
}

// ListLuggageUpdatesByHotel 查询某酒店寄存单修改记录
func ListLuggageUpdatesByHotel(hotelID int64) ([]models.LuggageUpdate, error) {
	if hotelID <= 0 {
		return nil, errors.New("invalid hotel id")
	}
	return repositories.ListUpdatesByHotel(hotelID)
}

// ListStoredLuggageByGuestName 获取某客人正在寄存的行李列表
func ListStoredLuggageByGuestName(hotelID int64, guestName string) ([]models.LuggageItem, error) {
	if hotelID <= 0 {
		return nil, errors.New("invalid hotel id")
	}
	if guestName == "" {
		return nil, errors.New("guest_name is empty")
	}
	return repositories.ListLuggageByHotelGuestAndStatus(hotelID, guestName, "stored")
}

// GetLuggageDetail 获取寄存单详情
func GetLuggageDetail(id int64) (models.LuggageItem, error) {
	if id <= 0 {
		return models.LuggageItem{}, errors.New("invalid luggage id")
	}
	item, err := repositories.GetLuggageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageItem{}, errors.New("luggage not found")
		}
		return models.LuggageItem{}, err
	}
	return item, nil
}

// GetLuggageDetailByCode 按取件码查询寄存单详情
func GetLuggageDetailByCode(code string) (models.LuggageItem, error) {
	if code == "" {
		return models.LuggageItem{}, errors.New("code is empty")
	}
	items, err := repositories.FindLuggageByCode(code)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return models.LuggageItem{}, errors.New("luggage not found")
		}
		return models.LuggageItem{}, err
	}
	if len(items) == 0 {
		return models.LuggageItem{}, errors.New("luggage not found")
	}
	return items[0], nil
}

// ListLuggageDetailByPhone 按客人手机号查询寄存单详情列表
func ListLuggageDetailByPhone(contactPhone string) ([]models.LuggageItem, error) {
	if contactPhone == "" {
		return nil, errors.New("contact_phone is empty")
	}
	return repositories.ListLuggageByGuest("", contactPhone, "")
}

// ListPickupCodesByUser 获取用户取件码列表
func ListPickupCodesByUser(username string, status string) ([]models.LuggageItem, error) {
	if username == "" {
		return nil, errors.New("username is empty")
	}
	return repositories.ListPickupCodesByUser(username, status)
}

// ListPickupCodesByPhone 按手机号查询取件码列表
func ListPickupCodesByPhone(contactPhone, status string) ([]models.LuggageItem, error) {
	if contactPhone == "" {
		return nil, errors.New("contact_phone is empty")
	}
	return repositories.ListPickupCodesByPhone(contactPhone, status)
}

// UpdateLuggageInfoRequest 修改寄存信息输入
type UpdateLuggageInfoRequest struct {
	GuestName    *string
	ContactPhone *string
	ContactEmail *string
	Description  *string
	Quantity     *int
	SpecialNotes *string
	PhotoURL     *string
	PhotoURLs    *[]string
	StoreroomID  *int64 // 新增：支持修改寄存室（迁移）
	UpdatedBy    string
}

// UpdateLuggageInfo 修改寄存信息（包含寄存室迁移）
func UpdateLuggageInfo(id int64, req UpdateLuggageInfoRequest) error {
	if id <= 0 {
		return errors.New("invalid luggage id")
	}

	item, err := repositories.GetLuggageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("luggage not found")
		}
		return err
	}

	// 如果要修改寄存室，需要验证目标寄存室
	if req.StoreroomID != nil && *req.StoreroomID != item.StoreroomID {
		// 验证目标寄存室是否存在
		targetRoom, err := repositories.GetStoreroomByID(*req.StoreroomID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("target storeroom not found")
			}
			return err
		}

		// 验证寄存室是否属于同一酒店
		if targetRoom.HotelID != item.HotelID {
			return errors.New("storeroom hotel mismatch")
		}

		// 验证寄存室是否可用
		if !targetRoom.IsActive {
			return errors.New("target storeroom is inactive")
		}

		// 验证寄存室是否已满
		count, err := repositories.CountStoredByStoreroom(*req.StoreroomID)
		if err != nil {
			return err
		}
		if count >= int64(targetRoom.Capacity) {
			return errors.New("target storeroom is full")
		}
	}

	updates := map[string]interface{}{}
	if req.GuestName != nil {
		updates["guest_name"] = *req.GuestName
	}
	if req.ContactPhone != nil {
		updates["contact_phone"] = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		updates["contact_email"] = *req.ContactEmail
	}
	if req.Description != nil {
		updates["description"] = *req.Description
	}
	if req.Quantity != nil {
		if *req.Quantity <= 0 {
			return errors.New("quantity must be greater than 0")
		}
		updates["quantity"] = *req.Quantity
	}
	if req.SpecialNotes != nil {
		updates["special_notes"] = *req.SpecialNotes
	}
	if req.PhotoURL != nil {
		updates["photo_url"] = *req.PhotoURL
	}
	if req.PhotoURLs != nil {
		data, err := json.Marshal(*req.PhotoURLs)
		if err != nil {
			return err
		}
		updates["photo_urls"] = string(data)
		if len(*req.PhotoURLs) > 0 {
			updates["photo_url"] = (*req.PhotoURLs)[0]
		} else {
			updates["photo_url"] = ""
		}
	}
	if req.StoreroomID != nil {
		updates["storeroom_id"] = *req.StoreroomID
	}

	if err := repositories.UpdateLuggageInfo(id, updates); err != nil {
		return err
	}

	// 记录修改历史
	updated := item
	if req.GuestName != nil {
		updated.GuestName = *req.GuestName
	}
	if req.ContactPhone != nil {
		updated.ContactPhone = *req.ContactPhone
	}
	if req.ContactEmail != nil {
		updated.ContactEmail = *req.ContactEmail
	}
	if req.Description != nil {
		updated.Description = *req.Description
	}
	if req.Quantity != nil {
		updated.Quantity = *req.Quantity
	}
	if req.SpecialNotes != nil {
		updated.SpecialNotes = *req.SpecialNotes
	}
	if req.PhotoURL != nil {
		updated.PhotoURL = *req.PhotoURL
	}
	if req.PhotoURLs != nil {
		updated.PhotoURLs = *req.PhotoURLs
		if len(*req.PhotoURLs) > 0 {
			updated.PhotoURL = (*req.PhotoURLs)[0]
		} else {
			updated.PhotoURL = ""
		}
	}
	if req.StoreroomID != nil {
		updated.StoreroomID = *req.StoreroomID
	}

	oldData, _ := json.Marshal(item)
	newData, _ := json.Marshal(updated)
	if req.UpdatedBy != "" {
		_ = repositories.CreateLuggageUpdate(&models.LuggageUpdate{
			HotelID:   item.HotelID,
			LuggageID: item.ID,
			UpdatedBy: req.UpdatedBy,
			OldData:   string(oldData),
			NewData:   string(newData),
		})
	}

	return nil
}

// UpdateLuggageCode 修改取件码
func UpdateLuggageCode(id int64, code string) error {
	if id <= 0 {
		return errors.New("invalid luggage id")
	}
	if code == "" {
		return errors.New("code is empty")
	}

	item, err := repositories.GetLuggageByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("luggage not found")
		}
		return err
	}

	return repositories.UpdateLuggageCode(item.ID, code)
}

// BindLuggageToUser 绑定行李到用户（按行李ID）
func BindLuggageToUser(luggageID int64, username string) error {
	if luggageID <= 0 || username == "" {
		return errors.New("invalid luggage_id or user_name")
	}

	item, err := repositories.GetLuggageByID(luggageID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("luggage not found")
		}
		return err
	}
	if item.Status != "stored" {
		return errors.New("luggage is not in stored status")
	}

	return repositories.BindLuggageToUser(item.ID, username)
}

// ListHistoryByGuest 按客人姓名/手机号查询取件历史
func ListHistoryByGuest(guestName, contactPhone string) ([]models.LuggageHistory, error) {
	if guestName == "" && contactPhone == "" {
		return nil, errors.New("guest_name and contact_phone cannot both be empty")
	}
	return repositories.ListHistoryByGuest(guestName, contactPhone)
}

// ListHistoryByHotel 按酒店查询取件历史
func ListHistoryByHotel(hotelID int64, guestName, contactPhone string) ([]models.LuggageHistory, error) {
	if hotelID <= 0 {
		return nil, errors.New("invalid hotel id")
	}
	return repositories.ListHistoryByHotel(hotelID, guestName, contactPhone)
}
