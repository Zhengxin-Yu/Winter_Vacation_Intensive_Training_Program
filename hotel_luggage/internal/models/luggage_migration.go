package models

import "time"

// LuggageMigration 对应 luggage_migrations 表（行李迁移记录）。
// 记录每次行李在寄存室之间迁移的详细信息。
type LuggageMigration struct {
	ID              int64     `gorm:"column:id;primaryKey;autoIncrement"` // 迁移记录ID（主键）
	HotelID         int64     `gorm:"column:hotel_id;not null"`           // 酒店ID
	LuggageID       int64     `gorm:"column:luggage_id;not null"`         // 行李ID（外键）
	FromStoreroomID int64     `gorm:"column:from_storeroom_id;not null"`  // 原寄存室ID（外键）
	ToStoreroomID   int64     `gorm:"column:to_storeroom_id;not null"`    // 目标寄存室ID（外键）
	MigratedBy      int64     `gorm:"column:migrated_by;not null"`        // 迁移操作员ID（外键）
	MigratedAt      time.Time `gorm:"column:migrated_at;autoCreateTime"`  // 迁移时间
}

// TableName 指定数据库表名
func (LuggageMigration) TableName() string {
	return "luggage_migrations"
}
