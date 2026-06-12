package models

import (
	"time"
)

// Conversation 会话
type Conversation struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Type         string    `gorm:"size:20;default:'private'" json:"type"` // private, group
	Participants string    `gorm:"size:500" json:"participants"`          // 参与者ID列表，逗号分隔
	LastMsgID    uint      `gorm:"default:0" json:"last_msg_id"`
	LastMsgTime  time.Time `json:"last_msg_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// Message 消息
type Message struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	ConversationID uint      `gorm:"index" json:"conversation_id"`
	SenderID       uint      `gorm:"index" json:"sender_id"`
	Content        string    `gorm:"type:text" json:"content"`
	Type           string    `gorm:"size:20;default:'text'" json:"type"` // text, image, file
	Status         int       `gorm:"default:1" json:"status"`            // 1-正常, 0-已删除
	CreatedAt      time.Time `json:"created_at"`
}

// UserConversation 用户会话关联（用于记录未读数等）
type UserConversation struct {
	ID             uint      `gorm:"primarykey" json:"id"`
	UserID         uint      `gorm:"index" json:"user_id"`
	ConversationID uint      `gorm:"index" json:"conversation_id"`
	UnreadCount    int       `gorm:"default:0" json:"unread_count"`
	LastReadTime   time.Time `json:"last_read_time"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}