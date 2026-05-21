package models

import (
    "time"
)

type DeletionRequest struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    ArticleID   uint      `gorm:"index" json:"article_id"`
    Article     Article   `gorm:"foreignKey:ArticleID" json:"article"`
    UserID      uint      `gorm:"index" json:"user_id"`
    User        User      `gorm:"foreignKey:UserID" json:"user"`
    Reason      string    `gorm:"type:text" json:"reason"`
    Status      string    `gorm:"default:pending" json:"status"` // pending, approved, rejected
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}