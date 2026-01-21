package models

import "time"

// LuggageItem 对应 luggage_items 表（行李寄存记录）。
// 包含客人信息、行李信息、取件码、状态等核心字段。
type LuggageItem struct {
	ID            int64      `gorm:"column:id;primaryKey;autoIncrement"`                                                 // 行李ID（主键）
	GuestName     string     `gorm:"column:guest_name;size:100;not null"`                                                // 客人姓名
	ContactPhone  string     `gorm:"column:contact_phone;size:20"`                                                       // 联系电话
	ContactEmail  string     `gorm:"column:contact_email;size:100"`                                                      // 联系邮箱
	Description   string     `gorm:"column:description;type:text"`                                                       // 行李描述
	Quantity      int        `gorm:"column:quantity;not null;default:1"`                                                 // 行李数量
	SpecialNotes  string     `gorm:"column:special_notes;type:text"`                                                     // 特殊备注
	StoreroomID   int64      `gorm:"column:storeroom_id;not null"`                                                       // 寄存室ID（外键）
	RetrievalCode string     `gorm:"column:retrieval_code;size:8;unique;not null"`                                       // 取回码
	QRCodeURL     string     `gorm:"column:qr_code_url;size:255"`                                                        // 二维码URL
	Status        string     `gorm:"column:status;type:enum('stored','retrieved','migrated');default:'stored';not null"` // 行李状态
	StoredBy      string     `gorm:"column:stored_by;size:50;not null"`                                                  // 存放操作员用户名
	RetrievedBy   *string    `gorm:"column:retrieved_by;size:50"`                                                        // 取回操作员用户名（可为空）
	RetrievedAt   *time.Time `gorm:"column:retrieved_at"`                                                                // 取回时间（可为空）
	StoredAt      time.Time  `gorm:"column:stored_at;autoCreateTime"`                                                    // 存放时间
	UpdatedAt     time.Time  `gorm:"column:updated_at;autoUpdateTime"`                                                   // 更新时间
}

// TableName 指定数据库表名
func (LuggageItem) TableName() string {
	return "luggage_items"
}
