package models

import (
	"time"
)

type UserFollow struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	FollowerID uint      `gorm:"index;uniqueIndex:idx_follow" json:"follower_id"`
	Follower   User      `gorm:"foreignKey:FollowerID" json:"follower"`
	FollowedID uint      `gorm:"index;uniqueIndex:idx_follow" json:"followed_id"`
	Followed   User      `gorm:"foreignKey:FollowedID" json:"followed"`
	CreatedAt  time.Time `json:"created_at"`
}

func (UserFollow) TableName() string {
	return "user_follows"
}
