package models

import (
	"time"
)

// Follow 关注关系
type Follow struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	FollowerID  uint      `gorm:"index;uniqueIndex:idx_follower_following" json:"follower_id"`  // 关注者ID
	FollowingID uint      `gorm:"index;uniqueIndex:idx_follower_following" json:"following_id"` // 被关注者ID
	CreatedAt   time.Time `json:"created_at"`
}

// ChatMessage 聊天消息
type ChatMessage struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	SenderID   uint      `gorm:"index" json:"sender_id"`
	ReceiverID uint      `gorm:"index" json:"receiver_id"`
	Content    string    `gorm:"type:text" json:"content"`
	IsRead     bool      `gorm:"default:false" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
}

// ChatSession 聊天会话（用于记录用户之间的聊天状态）
type ChatSession struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	User1ID       uint      `gorm:"index" json:"user1_id"`
	User2ID       uint      `gorm:"index" json:"user2_id"`
	LastMessageID *uint     `gorm:"index" json:"last_message_id"`
	UnreadCount   int       `gorm:"default:0" json:"unread_count"`
	IsBlocked     bool      `gorm:"default:false" json:"is_blocked"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
