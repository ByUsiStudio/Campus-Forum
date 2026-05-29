package models

import (
	"time"
)

type User struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	Username    string     `gorm:"uniqueIndex;size:50" json:"username"`
	QQNumber    string     `gorm:"uniqueIndex;size:20" json:"qq_number"`
	Email       string     `gorm:"uniqueIndex;size:100" json:"email"`
	DisplayName string     `gorm:"size:50" json:"display_name"`
	Password    string     `gorm:"size:255" json:"-"`
	Avatar      string     `gorm:"size:500" json:"avatar"`
	Role        string     `gorm:"default:user" json:"role"`     // admin, user, system
	Signature   string     `gorm:"size:200" json:"signature"`    // 个性化签名
	Status      string     `gorm:"default:normal" json:"status"` // normal, banned
	BanReason   string     `gorm:"size:500" json:"ban_reason"`   // 封禁原因
	BanTime     *time.Time `json:"ban_time"`                     // 封禁时间
	ResetToken  string     `gorm:"size:100" json:"-"`            // 密码重置Token
	ResetExpiry *time.Time `json:"-"`                            // Token过期时间
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
