package models

import (
	"encoding/json"
	"time"

	"gorm.io/gorm"
)

// LuggageHistory 对应 luggage_history 表（取件历史记录）
type LuggageHistory struct {
	ID            int64     `gorm:"column:id;primaryKey;autoIncrement"`    // 历史记录ID
	LuggageID     int64     `gorm:"column:luggage_id;not null"`            // 原行李ID
	GuestName     string    `gorm:"column:guest_name;size:100;not null"`   // 客人姓名
	ContactPhone  string    `gorm:"column:contact_phone;size:20"`          // 联系电话
	ContactEmail  string    `gorm:"column:contact_email;size:100"`         // 联系邮箱
	Description   string    `gorm:"column:description;type:text"`          // 行李描述
	Quantity      int       `gorm:"column:quantity;not null;default:1"`    // 行李数量
	SpecialNotes  string    `gorm:"column:special_notes;type:text"`        // 特殊备注
	PhotoURL      string    `gorm:"column:photo_url;size:255" json:"photo_url"` // 照片URL
	PhotoURLsRaw  string    `gorm:"column:photo_urls;type:text" json:"-"`       // 多图JSON（数据库字段）
	PhotoURLs     []string  `gorm:"-" json:"photo_urls,omitempty"`              // 多图数组（对外）
	HotelID       int64     `gorm:"column:hotel_id;not null"`              // 酒店ID
	StoreroomID   int64     `gorm:"column:storeroom_id;not null"`          // 寄存室ID
	RetrievalCode string    `gorm:"column:retrieval_code;size:8;not null"` // 取件码
	QRCodeURL     string    `gorm:"column:qr_code_url;size:255"`           // 二维码URL
	Status        string    `gorm:"column:status;size:20;not null"`        // 状态（retrieved）
	StoredBy      string    `gorm:"column:stored_by;size:50;not null"`     // 存放操作员用户名
	RetrievedBy   string    `gorm:"column:retrieved_by;size:50;not null"`  // 取件操作员用户名
	StoredAt      time.Time `gorm:"column:stored_at;not null"`             // 存放时间
	RetrievedAt   time.Time `gorm:"column:retrieved_at;not null"`          // 取件时间
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`      // 记录创建时间
}

// TableName 指定数据库表名
func (LuggageHistory) TableName() string {
	return "luggage_history"
}

// BeforeSave 在保存前把 PhotoURLs 写入 PhotoURLsRaw
func (item *LuggageHistory) BeforeSave(tx *gorm.DB) error {
	if item.PhotoURLs != nil {
		data, err := json.Marshal(item.PhotoURLs)
		if err != nil {
			return err
		}
		item.PhotoURLsRaw = string(data)
	}
	return nil
}

// AfterFind 在读取后把 PhotoURLsRaw 解析为 PhotoURLs
func (item *LuggageHistory) AfterFind(tx *gorm.DB) error {
	if item.PhotoURLsRaw == "" {
		return nil
	}
	var urls []string
	if err := json.Unmarshal([]byte(item.PhotoURLsRaw), &urls); err != nil {
		return err
	}
	item.PhotoURLs = urls
	return nil
}
