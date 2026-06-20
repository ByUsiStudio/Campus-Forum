package models

import (
	"time"
)

// Achievement 成就定义模型
type Achievement struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Code        string    `gorm:"uniqueIndex;size:50" json:"code"`  // 成就代码
	Name        string    `gorm:"size:100" json:"name"`             // 成就名称
	Description string    `gorm:"size:500" json:"description"`      // 成就描述
	Icon        string    `gorm:"size:200" json:"icon"`             // 成就图标
	Category    string    `gorm:"size:50;index" json:"category"`   // 成就分类：post, comment, social, special等
	Condition   string    `gorm:"type:text" json:"condition"`      // 成就条件（JSON格式）
	Reward      int       `gorm:"default:0" json:"reward"`          // 成就奖励经验值
	Rarity      string    `gorm:"size:20;default:'common'" json:"rarity"` // 稀有度：common, rare, epic, legendary
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Achievement) TableName() string {
	return "achievements"
}

// UserAchievement 用户获得的成就模型
type UserAchievement struct {
	ID            uint        `gorm:"primarykey" json:"id"`
	UserID        uint        `gorm:"uniqueIndex:idx_user_achievement;index" json:"user_id"`
	User          User        `gorm:"foreignKey:UserID" json:"user"`
	AchievementID uint        `gorm:"uniqueIndex:idx_user_achievement;index" json:"achievement_id"`
	Achievement   Achievement `gorm:"foreignKey:AchievementID" json:"achievement"`
	Progress      int         `gorm:"default:0" json:"progress"`           // 进度
	UnlockedAt    *time.Time  `json:"unlocked_at"`                        // 解锁时间
	CreatedAt     time.Time   `json:"created_at"`
}

// TableName 设置表名
func (UserAchievement) TableName() string {
	return "user_achievements"
}