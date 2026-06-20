package models

import (
	"time"
)

// UserStatistics 用户统计数据模型
type UserStatistics struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	UserID           uint      `gorm:"uniqueIndex" json:"user_id"`
	User             User      `gorm:"foreignKey:UserID" json:"user"`
	TotalArticles    int       `gorm:"default:0" json:"total_articles"`    // 总文章数
	TotalComments    int       `gorm:"default:0" json:"total_comments"`    // 总评论数
	TotalLikes       int       `gorm:"default:0" json:"total_likes"`       // 获得的总点赞数
	TotalViews       int       `gorm:"default:0" json:"total_views"`       // 获得的总浏览数
	TotalShares      int       `gorm:"default:0" json:"total_shares"`      // 获得的总分享数
	TotalFavorites   int       `gorm:"default:0" json:"total_favorites"`  // 获得的总收藏数
	TotalFollowers   int       `gorm:"default:0" json:"total_followers"`   // 总粉丝数
	TotalFollowing   int       `gorm:"default:0" json:"total_following"`  // 总关注数
	ActiveDays       int       `gorm:"default:0" json:"active_days"`      // 活跃天数
	LastActiveDate   *time.Time `json:"last_active_date"`                  // 最后活跃日期
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName 设置表名
func (UserStatistics) TableName() string {
	return "user_statistics"
}

// DailyStatistics 每日统计数据模型
type DailyStatistics struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	Date           string    `gorm:"uniqueIndex;size:10" json:"date"`    // 日期 YYYY-MM-DD
	NewUsers       int       `gorm:"default:0" json:"new_users"`         // 新增用户数
	ActiveUsers    int       `gorm:"default:0" json:"active_users"`     // 活跃用户数
	NewArticles    int       `gorm:"default:0" json:"new_articles"`     // 新增文章数
	NewComments    int       `gorm:"default:0" json:"new_comments"`     // 新增评论数
	TotalViews     int       `gorm:"default:0" json:"total_views"`     // 总浏览量
	TotalLikes     int       `gorm:"default:0" json:"total_likes"`     // 总点赞数
	TotalShares    int       `gorm:"default:0" json:"total_shares"`     // 总分享数
	TotalSignIns   int       `gorm:"default:0" json:"total_sign_ins"`   // 总签到数
	PeakOnline     int       `gorm:"default:0" json:"peak_online"`      // 峰值在线人数
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 设置表名
func (DailyStatistics) TableName() string {
	return "daily_statistics"
}

// ArticleStatistics 文章统计数据模型
type ArticleStatistics struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	ArticleID     uint      `gorm:"uniqueIndex" json:"article_id"`
	Article       Article   `gorm:"foreignKey:ArticleID" json:"article"`
	DailyViews    int       `gorm:"default:0" json:"daily_views"`     // 当日浏览量
	DailyLikes     int       `gorm:"default:0" json:"daily_likes"`     // 当日点赞数
	DailyComments  int       `gorm:"default:0" json:"daily_comments"` // 当日评论数
	DailyShares    int       `gorm:"default:0" json:"daily_shares"`     // 当日分享数
	WeeklyViews    int       `gorm:"default:0" json:"weekly_views"`     // 本周浏览量
	MonthlyViews   int       `gorm:"default:0" json:"monthly_views"`   // 本月浏览量
	LastResetAt    *time.Time `json:"last_reset_at"`                    // 最后重置时间
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// TableName 设置表名
func (ArticleStatistics) TableName() string {
	return "article_statistics"
}

// SystemOverview 系统概览数据模型
type SystemOverview struct {
	ID               uint      `gorm:"primarykey" json:"id"`
	TotalUsers       int       `gorm:"default:0" json:"total_users"`       // 总用户数
	TotalArticles    int       `gorm:"default:0" json:"total_articles"`   // 总文章数
	TotalComments    int       `gorm:"default:0" json:"total_comments"`   // 总评论数
	TotalCategories  int       `gorm:"default:0" json:"total_categories"` // 总分类数
	OnlineUsers      int       `gorm:"default:0" json:"online_users"`     // 当前在线用户数
	TodayActiveUsers int       `gorm:"default:0" json:"today_active_users"` // 今日活跃用户数
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName 设置表名
func (SystemOverview) TableName() string {
	return "system_overview"
}