package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type UserNotificationService struct{}

var UserNotification = &UserNotificationService{}

func (s *UserNotificationService) SendNotification(userID uint, title, content string) error {
	notification := models.PersonalNotification{
		UserID:  userID,
		Title:   title,
		Content: content,
		IsRead:  false,
	}

	if result := database.DB.Create(&notification); result.Error != nil {
		return utils.NewError("发送通知失败", 500)
	}
	return nil
}

func (s *UserNotificationService) SendBatchNotifications(userIDs []uint, title, content string) error {
	for _, userID := range userIDs {
		notification := models.PersonalNotification{
			UserID:  userID,
			Title:   title,
			Content: content,
			IsRead:  false,
		}
		database.DB.Create(&notification)
	}
	return nil
}

func (s *UserNotificationService) GetUserNotifications(userID uint, page, pageSize int) ([]models.PersonalNotification, int, error) {
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
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return notifications, totalPages, err
}

func (s *UserNotificationService) GetNotification(userID, notificationID uint) (*models.PersonalNotification, error) {
	var notification models.PersonalNotification
	err := database.DB.Where("user_id = ? AND id = ?", userID, notificationID).First(&notification).Error
	if err != nil {
		return nil, utils.NewError("通知不存在", 404)
	}
	return &notification, nil
}

func (s *UserNotificationService) MarkNotificationAsRead(userID, notificationID uint) error {
	var notification models.PersonalNotification
	result := database.DB.Where("user_id = ? AND id = ?", userID, notificationID).First(&notification)
	if result.Error != nil {
		return utils.NewError("通知不存在", 404)
	}

	notification.IsRead = true
	database.DB.Save(&notification)
	return nil
}

func (s *UserNotificationService) MarkAllNotificationsAsRead(userID uint) error {
	database.DB.Model(&models.PersonalNotification{}).Where("user_id = ?", userID).Update("is_read", true)
	return nil
}

func (s *UserNotificationService) DeletePersonalNotification(userID, notificationID uint) error {
	var notification models.PersonalNotification
	result := database.DB.Where("user_id = ? AND id = ?", userID, notificationID).First(&notification)
	if result.Error != nil {
		return utils.NewError("通知不存在", 404)
	}

	database.DB.Delete(&notification)
	return nil
}

func (s *UserNotificationService) ClearAllNotifications(userID uint) error {
	database.DB.Where("user_id = ?", userID).Delete(&models.PersonalNotification{})
	return nil
}

func (s *UserNotificationService) AdminGetUserNotifications(userID uint, page, pageSize int) ([]models.PersonalNotification, int, error) {
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
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&notifications).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return notifications, totalPages, err
}
