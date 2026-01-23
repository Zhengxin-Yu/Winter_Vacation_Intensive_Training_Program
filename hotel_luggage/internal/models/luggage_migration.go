package models

import "time"

// LuggageUpdate 对应“行李寄存信息修改表”，记录寄存单修改前后信息
type LuggageUpdate struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`     // 记录ID
	HotelID   int64     `gorm:"column:hotel_id;not null"`               // 酒店ID
	LuggageID int64     `gorm:"column:luggage_id;not null"`             // 行李ID
	UpdatedBy string    `gorm:"column:updated_by;size:50;not null"`      // 操作员用户名
	OldData   string    `gorm:"column:old_data;type:text;not null"`      // 修改前快照（JSON）
	NewData   string    `gorm:"column:new_data;type:text;not null"`      // 修改后快照（JSON）
	UpdatedAt time.Time `gorm:"column:updated_at;autoCreateTime"`        // 修改时间
}

// TableName 指定数据库表名
func (LuggageUpdate) TableName() string {
	return "行李寄存信息修改表"
}
