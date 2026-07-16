package repository

import (
	"forum/models"
)

type NotificationRepository struct {
	*BaseRepository[models.Notification]
}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{
		BaseRepository: NewBaseRepository[models.Notification](),
	}
}

func (r *NotificationRepository) GetUserNotifications(userID uint, page, pageSize int) ([]models.Notification, int64, error) {
	var notifications []models.Notification
	var total int64

	query := r.db.Model(&models.Notification{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications).Error

	return notifications, total, err
}

func (r *NotificationRepository) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Notification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (r *NotificationRepository) MarkAllRead(userID uint) error {
	return r.db.Model(&models.Notification{}).Where("user_id = ?", userID).Update("is_read", true).Error
}

type UserNotificationRepository struct {
	*BaseRepository[models.UserNotification]
}

func NewUserNotificationRepository() *UserNotificationRepository {
	return &UserNotificationRepository{
		BaseRepository: NewBaseRepository[models.UserNotification](),
	}
}
