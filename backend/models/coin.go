package models

import (
	"time"
)

type CoinRecord struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	ArticleID uint      `gorm:"index" json:"article_id"`
	CoinCount int       `gorm:"default:1" json:"coin_count"` // 投币数量
	CreatedAt time.Time `json:"created_at"`
}

func (CoinRecord) TableName() string {
	return "coin_records"
}