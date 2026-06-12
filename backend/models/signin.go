package models

import (
	"time"
)

// SignInRecord 签到记录模型
type SignInRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	SignInAt  time.Time `gorm:"index" json:"sign_in_at"` // 签到日期（年月日）
	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置表名
func (SignInRecord) TableName() string {
	return "sign_in_records"
}
