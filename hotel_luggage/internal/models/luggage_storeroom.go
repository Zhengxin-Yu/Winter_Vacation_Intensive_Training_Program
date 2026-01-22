package models

import "time"

// LuggageStoreroom 对应 luggage_storerooms 表（寄存室）。
// 记录寄存室名称、容量、启用状态等信息。
type LuggageStoreroom struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`     // 寄存室ID（主键）
	HotelID   int64     `gorm:"column:hotel_id;not null"`               // 所属酒店ID
	Name      string    `gorm:"column:name;size:100;not null"`          // 寄存室名称
	Location  string    `gorm:"column:location;size:255"`               // 位置描述
	Capacity  int       `gorm:"column:capacity;not null;default:0"`     // 容量
	IsActive  bool      `gorm:"column:is_active;not null;default:true"` // 是否启用
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`       // 创建时间
}

// TableName 指定数据库表名
func (LuggageStoreroom) TableName() string {
	return "luggage_storerooms"
}
