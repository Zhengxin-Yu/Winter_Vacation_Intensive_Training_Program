package models

import "time"

// Hotel 对应 hotels 表（酒店信息）
type Hotel struct {
	ID        int64     `gorm:"column:id;primaryKey;autoIncrement"`     // 酒店ID
	Name      string    `gorm:"column:name;size:100;not null"`          // 酒店名称
	Address   string    `gorm:"column:address;size:255"`                // 地址
	Phone     string    `gorm:"column:phone;size:50"`                   // 联系电话
	IsActive  bool      `gorm:"column:is_active;not null;default:true"` // 是否启用
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`       // 创建时间
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`       // 更新时间
}

// TableName 指定数据库表名
func (Hotel) TableName() string {
	return "hotels"
}
