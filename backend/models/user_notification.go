package models

import (
	"time"
)

// PersonalNotification 用户个人通知（单独通知）
type PersonalNotification struct {
	ID          uint       `gorm:"primarykey" json:"id"`
	SenderID    uint       `gorm:"index" json:"sender_id"` // 发送者ID
	Sender      User       `gorm:"foreignKey:SenderID" json:"sender"`
	UserID      uint       `gorm:"index" json:"user_id"` // 接收者ID
	User        User       `gorm:"foreignKey:UserID" json:"user"`
	Type        string     `gorm:"size:30" json:"type"` // system, warning, promotion, reminder
	Title       string     `gorm:"size:200" json:"title"`
	Content     string     `gorm:"type:text" json:"content"`
	RelatedType string     `gorm:"size:30" json:"related_type"` // article, comment, user, etc.
	RelatedID   uint       `gorm:"index" json:"related_id"`     // 关联内容的ID
	Link        string     `gorm:"size:500" json:"link"`        // 可跳转的链接
	IsRead      bool       `gorm:"default:false" json:"is_read"`
	ReadAt      *time.Time `json:"read_at"`
	Priority    string     `gorm:"size:10;default:'normal'" json:"priority"` // high, normal, low
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// PermissionGroup 权限组
type PermissionGroup struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Name        string    `gorm:"size:50;uniqueIndex" json:"name"`
	Description string    `gorm:"size:200" json:"description"`
	Level       int       `gorm:"default:1" json:"level"`          // 权限级别，数字越大权限越高
	IsDefault   bool      `gorm:"default:false" json:"is_default"` // 是否为默认组
	IsActive    bool      `gorm:"default:true" json:"is_active"`   // 是否启用
	Permissions string    `gorm:"type:text" json:"permissions"`    // JSON格式的权限列表
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserPermissionGroup 用户权限组关联
type UserPermissionGroup struct {
	ID                uint            `gorm:"primarykey" json:"id"`
	UserID            uint            `gorm:"index" json:"user_id"`
	User              User            `gorm:"foreignKey:UserID" json:"user"`
	PermissionGroupID uint            `gorm:"index" json:"permission_group_id"`
	PermissionGroup   PermissionGroup `gorm:"foreignKey:PermissionGroupID" json:"permission_group"`
	ExpiresAt         *time.Time      `json:"expires_at"`              // 过期时间，nil表示永久
	GrantedBy         uint            `gorm:"index" json:"granted_by"` // 授权人ID
	CreatedAt         time.Time       `json:"created_at"`
}

// SystemLog 系统操作日志
type SystemLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	User       User      `gorm:"foreignKey:UserID" json:"user"`
	Action     string    `gorm:"size:50" json:"action"`       // 操作类型
	Module     string    `gorm:"size:50;index" json:"module"` // 模块
	TargetType string    `gorm:"size:50" json:"target_type"`  // 目标类型
	TargetID   uint      `gorm:"index" json:"target_id"`      // 目标ID
	Details    string    `gorm:"type:text" json:"details"`    // 操作详情
	IP         string    `gorm:"size:50" json:"ip"`           // IP地址
	UserAgent  string    `gorm:"size:500" json:"user_agent"`  // 用户代理
	CreatedAt  time.Time `json:"created_at"`
}
