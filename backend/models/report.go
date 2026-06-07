package models

import (
	"time"
)

type Report struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	ReporterID  uint       `gorm:"index" json:"reporter_id"`
	Reporter    User       `gorm:"foreignKey:ReporterID" json:"reporter"`
	TargetType  string     `gorm:"size:50;index" json:"target_type"` // article, comment, user
	TargetID    uint       `gorm:"index" json:"target_id"`
	Reason      string     `gorm:"type:text" json:"reason"`
	Description string     `gorm:"type:text" json:"description"`
	Status      string     `gorm:"default:pending" json:"status"` // pending, resolved, rejected
	HandlerID   *uint      `gorm:"index" json:"handler_id"`
	Handler     *User      `gorm:"foreignKey:HandlerID" json:"handler"`
	HandleNote  string     `gorm:"type:text" json:"handle_note"`
	HandledAt   *time.Time `json:"handled_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}
