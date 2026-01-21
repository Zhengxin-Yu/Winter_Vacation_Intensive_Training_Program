package models

import "time"

// User 对应 users 表（系统用户）。
// 用于登录鉴权与操作人员管理。
type User struct {
	ID           int64     `gorm:"column:id;primaryKey;autoIncrement"`                         // 用户ID（主键）
	Username     string    `gorm:"column:username;size:50;unique;not null"`                    // 用户名（唯一）
	PasswordHash string    `gorm:"column:password_hash;size:60;not null"`                      // 密码哈希
	Role         string    `gorm:"column:role;type:enum('staff','admin','guest');default:'staff';not null"` // 角色：staff/admin/guest
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime"`                           // 创建时间
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime"`                           // 更新时间
}

// TableName 指定数据库表名
func (User) TableName() string {
	return "users"
}
