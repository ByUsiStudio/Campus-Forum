package repository

import (
	"forum/models"

	"gorm.io/gorm"
)

// PersonalNotificationRepository 用户通知数据访问层
type PersonalNotificationRepository struct {
	*BaseRepository[models.PersonalNotification]
}

// NewPersonalNotificationRepository 创建用户通知 Repository
func NewPersonalNotificationRepository() *PersonalNotificationRepository {
	return &PersonalNotificationRepository{
		BaseRepository: NewBaseRepository[models.PersonalNotification](),
	}
}

// GetUserNotifications 获取用户的通知列表
func (r *PersonalNotificationRepository) GetUserNotifications(userID uint, page, pageSize int, unreadOnly bool) ([]models.PersonalNotification, int64, int64, error) {
	var notifications []models.PersonalNotification
	var total int64
	var unreadCount int64

	query := r.db.Model(&models.PersonalNotification{}).Where("user_id = ?", userID)
	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	query.Count(&total)
	r.db.Model(&models.PersonalNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Sender").
		Order("priority DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, unreadCount, err
}

// MarkAsRead 标记通知为已读
func (r *PersonalNotificationRepository) MarkAsRead(id uint) error {
	return r.db.Model(&models.PersonalNotification{}).Where("id = ?", id).Updates(map[string]interface{}{
		"is_read": true,
		"read_at": gorm.Expr("NOW()"),
	}).Error
}

// MarkAllAsRead 标记用户所有通知为已读
func (r *PersonalNotificationRepository) MarkAllAsRead(userID uint) error {
	return r.db.Model(&models.PersonalNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": gorm.Expr("NOW()"),
		}).Error
}

// DeleteByUser 删除用户的所有通知
func (r *PersonalNotificationRepository) DeleteByUser(userID uint) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.PersonalNotification{}).Error
}

// GetUnreadCount 获取用户未读通知数
func (r *PersonalNotificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.PersonalNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}
