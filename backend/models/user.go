package models

import (
	"time"
)

type User struct {
	ID                uint       `gorm:"primarykey" json:"id"`
	Username          string     `gorm:"uniqueIndex;size:50" json:"username"`
	QQNumber          string     `gorm:"uniqueIndex;size:20" json:"qq_number"`
	DisplayName       string     `gorm:"size:50" json:"display_name"`
	Password          string     `gorm:"size:255" json:"-"`
	Avatar            string     `gorm:"size:500" json:"avatar"`
	Role              string     `gorm:"default:user" json:"role"`             // admin, user, system
	Signature         string     `gorm:"size:200" json:"signature"`            // 个性化签名
	Status            string     `gorm:"default:normal" json:"status"`         // normal, banned
	BanReason         string     `gorm:"size:500" json:"ban_reason"`           // 封禁原因
	BanTime           *time.Time `json:"ban_time"`                             // 封禁时间
	OnlineStatus      string     `gorm:"default:offline" json:"online_status"` // online, offline
	LastActiveAt      *time.Time `json:"last_active_at"`                       // 最后活跃时间
	ResetToken        string     `gorm:"size:100" json:"-"`                    // 密码重置验证码
	ResetIdentifier   string     `gorm:"size:100" json:"-"`                    // 密码重置标识token
	ResetExpiry       *time.Time `json:"-"`                                    // Token过期时间
	SignInDays        int        `gorm:"default:0" json:"sign_in_days"`        // 当前连续签到天数
	TotalSignIns      int        `gorm:"default:0" json:"total_sign_ins"`      // 总签到次数
	MaxContinuousDays int        `gorm:"default:0" json:"max_continuous_days"` // 最长连续签到天数
	TotalCoins        int        `gorm:"default:0" json:"total_coins"`         // 累计币
	CreatedAt         time.Time  `json:"created_at"`
	UpdatedAt         time.Time  `json:"updated_at"`
}
