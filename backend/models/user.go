package models

import (
    "time"
)

type User struct {
    ID          uint      `gorm:"primarykey" json:"id"`
    Username    string    `gorm:"uniqueIndex;size:50" json:"username"`
    QQNumber    string    `gorm:"uniqueIndex;size:20" json:"qq_number"`
    DisplayName string    `gorm:"size:50" json:"display_name"`
    Password    string    `gorm:"size:255" json:"-"`
    Avatar      string    `gorm:"size:500" json:"avatar"`
    Role        string    `gorm:"default:user" json:"role"` // admin, user
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}