package models

import (
	"time"
)

type Announcement struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Content     string    `gorm:"type:text" json:"content"`
	ContentHTML string    `gorm:"type:text" json:"content_html"`
	UpdatedAt   time.Time `json:"updated_at"`
}
