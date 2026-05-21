package models

import (
    "time"
)

type Comment struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    Content   string    `gorm:"type:text;not null" json:"content"`
    UserID    uint      `gorm:"index" json:"user_id"`
    User      User      `gorm:"foreignKey:UserID" json:"user"`
    ArticleID uint      `gorm:"index" json:"article_id"`
    CreatedAt time.Time `json:"created_at"`
}