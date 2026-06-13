package models

import (
	"time"
)

// SignInRecord 签到记录模型
type SignInRecord struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	SignInAt      time.Time `gorm:"index" json:"sign_in_at"`       // 签到时间
	SignInDate    string    `gorm:"index;size:10" json:"sign_in_date"` // 签到日期（YYYY-MM-DD格式，方便查询）
	ContinuousDay int       `gorm:"default:1" json:"continuous_day"` // 当次签到后的连续天数
	RewardPoints  int       `gorm:"default:0" json:"reward_points"` // 当次签到获得的积分
	IPAddress     string    `gorm:"size:50" json:"ip_address"`     // 签到IP
	CreatedAt     time.Time `json:"created_at"`
}

// TableName 设置表名
func (SignInRecord) TableName() string {
	return "sign_in_records"
}

// SignInConfig 签到配置模型
type SignInConfig struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	DailyPoints    int       `gorm:"default:1" json:"daily_points"`    // 每日签到基础积分
	WeeklyBonus    int       `gorm:"default:7" json:"weekly_bonus"`   // 连续7天签到奖励
	MonthlyBonus   int       `gorm:"default:30" json:"monthly_bonus"` // 连续30天签到奖励
	YearlyBonus    int       `gorm:"default:100" json:"yearly_bonus"` // 连续365天签到奖励
	Enabled        bool      `gorm:"default:true" json:"enabled"`     // 是否启用签到
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 设置表名
func (SignInConfig) TableName() string {
	return "sign_in_config"
}
