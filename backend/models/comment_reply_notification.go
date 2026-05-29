package models

import (
	"time"
)

// CommentReplyNotification 评论回复通知
type CommentReplyNotification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`         // 被回复的用户ID
	CommentID uint      `gorm:"index" json:"comment_id"`      // 回复的评论ID
	ReplyID   uint      `gorm:"index" json:"reply_id"`        // 回复内容的ID
	ArticleID uint      `gorm:"index" json:"article_id"`      // 文章ID
	IsRead    bool      `gorm:"default:false" json:"is_read"` // 是否已读
	CreatedAt time.Time `json:"created_at"`
}
