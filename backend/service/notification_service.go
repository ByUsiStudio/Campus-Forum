package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type NotificationService struct{}

var Notification = &NotificationService{}

func (s *NotificationService) GetNotifications(userID uint, page, pageSize int) ([]models.PersonalNotification, int, error) {
	var notifications []models.PersonalNotification
	var total int64

	query := database.DB.Model(&models.PersonalNotification{}).Where("user_id = ?", userID)
	query.Count(&total)

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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return notifications, totalPages, err
}

func (s *NotificationService) GetUnreadCount(userID uint) (int64, error) {
	var count int64
	err := database.DB.Model(&models.PersonalNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count).Error
	return count, err
}

func (s *NotificationService) MarkNotificationRead(userID, notificationID uint) error {
	var notification models.PersonalNotification
	if result := database.DB.First(&notification, notificationID); result.Error != nil {
		return utils.NewError("通知不存在", 404)
	}

	if notification.UserID != userID {
		return utils.NewError("无权操作该通知", 403)
	}

	notification.IsRead = true
	database.DB.Save(&notification)
	return nil
}

func (s *NotificationService) MarkAllNotificationsRead(userID uint) error {
	database.DB.Model(&models.PersonalNotification{}).Where("user_id = ?", userID).Update("is_read", true)
	return nil
}

func (s *NotificationService) CreateNotification(senderID, userID uint, notificationType, title, content, relatedType string, relatedID uint) error {
	notification := models.PersonalNotification{
		SenderID:    senderID,
		UserID:      userID,
		Type:        notificationType,
		Title:       title,
		Content:     content,
		RelatedType: relatedType,
		RelatedID:   relatedID,
		IsRead:      false,
	}

	database.DB.Create(&notification)
	return nil
}

func (s *NotificationService) DeleteNotification(userID, notificationID uint) error {
	var notification models.PersonalNotification
	if result := database.DB.First(&notification, notificationID); result.Error != nil {
		return utils.NewError("通知不存在", 404)
	}

	if notification.UserID != userID {
		return utils.NewError("无权删除该通知", 403)
	}

	database.DB.Delete(&notification)
	return nil
}

func (s *NotificationService) GetAdminNotifications(page, pageSize int) ([]models.Notification, int, error) {
	var notifications []models.Notification
	var total int64

	query := database.DB.Model(&models.Notification{}).Where("target = ?", "admin")
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return notifications, totalPages, err
}
