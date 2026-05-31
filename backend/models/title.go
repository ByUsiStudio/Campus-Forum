package models

import (
	"time"
)

type Title struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;not null" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	Color       string    `gorm:"size:20;default:'#6750A4'" json:"color"`
	Icon        string    `gorm:"size:50" json:"icon"`
	IsActive    bool      `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time `json:"created_at"`
}

type UserTitle struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index;uniqueIndex:idx_user_title" json:"user_id"`
	TitleID   uint      `gorm:"index;uniqueIndex:idx_user_title" json:"title_id"`
	Title     Title     `gorm:"foreignKey:TitleID"`
	Reason    string    `gorm:"size:200" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}
