package models

import (
	"time"
)

// Notification 系统通知
type Notification struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Type      string    `gorm:"size:20;not null" json:"type"` // system, activity, update, warning
	Title     string    `gorm:"size:100;not null" json:"title"`
	Content   string    `gorm:"type:text;not null" json:"content"`
	Target    string    `gorm:"size:20;not null;default:'all'" json:"target"` // all, admin
	CreatedAt time.Time `json:"created_at"`
}

// UserNotification 用户通知阅读状态
type UserNotification struct {
	ID             uint       `gorm:"primarykey" json:"id"`
	UserID         uint       `gorm:"index;uniqueIndex:idx_user_notification" json:"user_id"`
	NotificationID uint       `gorm:"index;uniqueIndex:idx_user_notification" json:"notification_id"`
	IsRead         bool       `gorm:"default:false" json:"is_read"`
	ReadAt         *time.Time `json:"read_at"`
	CreatedAt      time.Time  `json:"created_at"`
}
