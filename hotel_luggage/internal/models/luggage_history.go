package models

import "time"

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
	StoreroomID   int64     `gorm:"column:storeroom_id;not null"`          // 寄存室ID
	RetrievalCode string    `gorm:"column:retrieval_code;size:8;not null"` // 取件码
	QRCodeURL     string    `gorm:"column:qr_code_url;size:255"`           // 二维码URL
	Status        string    `gorm:"column:status;size:20;not null"`        // 状态（retrieved）
	StoredBy      int64     `gorm:"column:stored_by;not null"`             // 存放操作员ID
	RetrievedBy   int64     `gorm:"column:retrieved_by;not null"`          // 取件操作员ID
	StoredAt      time.Time `gorm:"column:stored_at;not null"`             // 存放时间
	RetrievedAt   time.Time `gorm:"column:retrieved_at;not null"`          // 取件时间
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime"`      // 记录创建时间
}

// TableName 指定数据库表名
func (LuggageHistory) TableName() string {
	return "luggage_history"
}
