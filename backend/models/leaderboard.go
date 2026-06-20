package models

import (
	"time"
)

// Leaderboard 排行榜模型
type Leaderboard struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Type      string    `gorm:"size:50;index" json:"type"`   // 排行榜类型：experience, articles, likes, comments, sign_in等
	Period    string    `gorm:"size:20;index" json:"period"` // 统计周期：daily, weekly, monthly, all_time
	UserID    uint      `gorm:"index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	Score     float64   `gorm:"default:0" json:"score"`    // 分数
	Rank      int       `gorm:"default:0" json:"rank"`     // 排名
	Date      string    `gorm:"size:10;index" json:"date"` // 统计日期 YYYY-MM-DD
	Week      int       `gorm:"default:0" json:"week"`     // 周数
	Month     string    `gorm:"size:7;index" json:"month"` // 月份 YYYY-MM
	Year      int       `gorm:"default:0" json:"year"`     // 年份
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Leaderboard) TableName() string {
	return "leaderboards"
}

// UserBadge 用户徽章模型
type UserBadge struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	UserID      uint      `gorm:"index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	BadgeType   string    `gorm:"size:50;index" json:"badge_type"`  // 徽章类型：top_author, active_user, contributor等
	BadgeName   string    `gorm:"size:100" json:"badge_name"`       // 徽章名称
	BadgeIcon   string    `gorm:"size:200" json:"badge_icon"`       // 徽章图标
	Description string    `gorm:"size:500" json:"description"`      // 徽章描述
	IsDisplayed bool      `gorm:"default:true" json:"is_displayed"` // 是否展示
	EarnedAt    time.Time `json:"earned_at"`                        // 获得时间
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 设置表名
func (UserBadge) TableName() string {
	return "user_badges"
}

// UserActivity 用户活跃度模型
type UserActivity struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	Date         string    `gorm:"uniqueIndex:idx_user_date;size:10" json:"date"` // 日期 YYYY-MM-DD
	LoginCount   int       `gorm:"default:0" json:"login_count"`                  // 登录次数
	PostCount    int       `gorm:"default:0" json:"post_count"`                   // 发帖数
	CommentCount int       `gorm:"default:0" json:"comment_count"`                // 评论数
	LikeCount    int       `gorm:"default:0" json:"like_count"`                   // 点赞数
	ViewCount    int       `gorm:"default:0" json:"view_count"`                   // 浏览数
	ActiveScore  float64   `gorm:"default:0" json:"active_score"`                 // 活跃度分数
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 设置表名
func (UserActivity) TableName() string {
	return "user_activities"
}
