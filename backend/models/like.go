package models

import (
    "time"
)

type Like struct {
    ID        uint      `gorm:"primarykey" json:"id"`
    UserID    uint      `gorm:"index;uniqueIndex:idx_user_article" json:"user_id"`
    ArticleID uint      `gorm:"index;uniqueIndex:idx_user_article" json:"article_id"`
    CreatedAt time.Time `json:"created_at"`
}