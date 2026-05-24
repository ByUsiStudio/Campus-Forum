package models

import (
	"time"
)

type Comment struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	UserID     uint      `gorm:"index" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	ArticleID  uint      `gorm:"index" json:"article_id"`
	ParentID   *uint     `gorm:"index" json:"parent_id"` // 回复的评论ID，NULL表示顶级评论
	LikeCount  int       `gorm:"default:0" json:"like_count"`
	ReplyCount int       `gorm:"default:0" json:"reply_count"` // 回复数量
	Replies   []Comment `gorm:"-" json:"replies"`              // 回复列表（不存储在数据库）
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentLike 评论点赞
type CommentLike struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;uniqueIndex:idx_user_comment" json:"user_id"`
	CommentID uint      `gorm:"index;uniqueIndex:idx_user_comment" json:"comment_id"`
	CreatedAt time.Time `json:"created_at"`
}

// ViewHistory 浏览记录（用于IP校验）
type ViewHistory struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"index" json:"article_id"`
	IP        string    `gorm:"index;size:45" json:"ip"` // 支持IPv6
	UserID    *uint     `gorm:"index" json:"user_id"`    // 已登录用户的ID
	CreatedAt time.Time `json:"created_at"`
}
