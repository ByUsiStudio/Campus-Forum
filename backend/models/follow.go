package models

import (
	"time"
)

// Follow 用户关注关系
type Follow struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	FollowerID  uint      `gorm:"index" json:"follower_id"` // 关注者ID
	Follower    User      `gorm:"foreignKey:FollowerID;-:constraint" json:"follower"`
	FollowingID uint      `gorm:"index" json:"following_id"` // 被关注者ID
	Following   User      `gorm:"foreignKey:FollowingID;-:constraint" json:"following"`
	CreatedAt   time.Time `json:"created_at"`
}
