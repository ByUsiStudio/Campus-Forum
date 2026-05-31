package models

import (
	"time"
)

type Comment struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Content     string    `gorm:"type:text;not null" json:"content"`
	UserID      uint      `gorm:"index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	ArticleID   uint      `gorm:"index" json:"article_id"`
	ParentID    *uint     `gorm:"index" json:"parent_id"`
	LikeCount   int       `gorm:"default:0" json:"like_count"`
	ReplyCount  int       `gorm:"default:0" json:"reply_count"`
	Replies     []Comment `gorm:"-" json:"replies"`
	IsAnonymous bool      `gorm:"default:false" json:"is_anonymous"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type CommentLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;uniqueIndex:idx_user_comment" json:"user_id"`
	CommentID uint      `gorm:"index;uniqueIndex:idx_user_comment" json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
}

type ViewHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"index" json:"article_id"`
	IP        string    `gorm:"index;size:45" json:"ip"`
	UserID    *uint     `gorm:"index" json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}
