package models

import (
	"time"
)

type FollowNotification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	SenderID  uint      `gorm:"index" json:"sender_id"`
	ArticleID uint      `gorm:"index" json:"article_id"`
	Type      string    `gorm:"size:20;default:'new_article'" json:"type"`
	IsRead    bool      `gorm:"default:false" json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}
