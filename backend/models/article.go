package models

import (
	"time"
)

type Article struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	Title         string    `gorm:"size:200;not null" json:"title"`
	Content       string    `gorm:"type:text" json:"content"`
	ContentHTML   string    `gorm:"type:text" json:"content_html"`
	UserID        uint      `gorm:"index" json:"user_id"`
	User          User      `gorm:"foreignKey:UserID" json:"user"`
	CategoryID    uint      `gorm:"index" json:"category_id"`
	Category      Category  `gorm:"foreignKey:CategoryID" json:"category"`
	LikeCount     int       `gorm:"default:0" json:"like_count"`
	CommentCount  int       `gorm:"default:0" json:"comment_count"`
	FavoriteCount int       `gorm:"default:0" json:"favorite_count"`
	ViewCount     int       `gorm:"default:0" json:"view_count"`
	ShareCount    int       `gorm:"default:0" json:"share_count"`
	CoinCount     int       `gorm:"default:0" json:"coin_count"`       // 投币数
	Status        string    `gorm:"default:published" json:"status"`   // published, draft, deleted, pending
	IsPinned      bool      `gorm:"default:false" json:"is_pinned"`    // 是否置顶
	PinnedAt      *time.Time `json:"pinned_at"`                         // 置顶时间
	VoiceURL      string    `gorm:"size:500" json:"voice_url"`         // 语音文件URL
	IsAnonymous   bool      `gorm:"default:false" json:"is_anonymous"` // 是否匿名发布
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
