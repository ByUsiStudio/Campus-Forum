package models

import (
	"time"
)

// Friend 用户好友关系（双向）
type Friend struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	UserID        uint      `gorm:"index" json:"user_id"`         // 用户ID
	User          User      `gorm:"foreignKey:UserID;-:constraint" json:"user"`
	FriendID      uint      `gorm:"index" json:"friend_id"`       // 好友ID
	Friend        User      `gorm:"foreignKey:FriendID;-:constraint" json:"friend"`
	DisplayName   string    `gorm:"size:100" json:"display_name"` // 好友备注名
	Status        int       `gorm:"default:1" json:"status"`      // 状态: 1-正常, 0-已删除
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// FriendRequest 好友请求
type FriendRequest struct {
	ID            uint      `gorm:"primarykey" json:"id"`
	SenderID      uint      `gorm:"index" json:"sender_id"`     // 请求发送者ID
	Sender        User      `gorm:"foreignKey:SenderID;-:constraint" json:"sender"`
	ReceiverID    uint      `gorm:"index" json:"receiver_id"`   // 请求接收者ID
	Receiver      User      `gorm:"foreignKey:ReceiverID;-:constraint" json:"receiver"`
	Message       string    `gorm:"size:500" json:"message"`    // 请求消息
	Status        int       `gorm:"default:0" json:"status"`     // 状态: 0-待处理, 1-已同意, 2-已拒绝
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}