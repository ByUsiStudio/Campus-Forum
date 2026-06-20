package models

import (
	"time"
)

// UserLevel 用户等级模型
type UserLevel struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"uniqueIndex" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Level       int       `gorm:"default:1" json:"level"`           // 当前等级
	Experience  int       `gorm:"default:0" json:"experience"`      // 当前经验值
	NextLevel   int       `gorm:"default:100" json:"next_level"`    // 升级所需经验
	Title       string    `gorm:"size:50" json:"title"`             // 等级称号
	Badge       string    `gorm:"size:200" json:"badge"`            // 等级徽章图标
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 设置表名
func (UserLevel) TableName() string {
	return "user_levels"
}

// ExperienceRecord 经验记录模型
type ExperienceRecord struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Type        string    `gorm:"size:50;index" json:"type"`        // 获得经验的类型：login, post, comment, like, share等
	Amount      int       `gorm:"default:0" json:"amount"`          // 获得的经验值
	Description string    `gorm:"size:200" json:"description"`       // 描述
	RelatedID   *uint     `json:"related_id"`                       // 关联ID（文章ID、评论ID等）
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 设置表名
func (ExperienceRecord) TableName() string {
	return "experience_records"
}

// LevelConfig 等级配置模型
type LevelConfig struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Level       int       `gorm:"uniqueIndex" json:"level"`         // 等级
	MinExp      int       `gorm:"default:0" json:"min_exp"`         // 所需最小经验
	Title       string    `gorm:"size:50" json:"title"`             // 等级称号
	Badge       string    `gorm:"size:200" json:"badge"`            // 等级徽章
	Privileges  string    `gorm:"type:text" json:"privileges"`      // 等级特权（JSON格式）
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// TableName 设置表名
func (LevelConfig) TableName() string {
	return "level_configs"
}