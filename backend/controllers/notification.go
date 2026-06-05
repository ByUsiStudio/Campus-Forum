package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// CreateNotification 创建通知（管理员）
func CreateNotification(c *gin.Context) {
	var input struct {
		Type    string `json:"type" binding:"required"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
		Target  string `json:"target"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Target == "" {
		input.Target = "all"
	}

	if input.Type != "system" && input.Type != "activity" && input.Type != "update" && input.Type != "warning" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的通知类型"})
		return
	}

	if input.Target != "all" && input.Target != "admin" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的目标群体"})
		return
	}

	// XSS过滤：清理通知标题和内容
	safeTitle := utils.SanitizeHTML(input.Title)
	safeContent := utils.SanitizeHTML(input.Content)

	notification := models.Notification{
		Type:    input.Type,
		Title:   safeTitle,
		Content: safeContent,
		Target:  input.Target,
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      "通知发送成功",
		"notification": notification,
	})
}

// GetNotifications 获取通知列表（用户）
func GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var notifications []models.Notification
	query := database.DB.Model(&models.Notification{})

	if role != "admin" {
		query = query.Where("target = ? OR target = ?", "all", "admin")
	}

	query.Order("created_at DESC").Find(&notifications)

	var userNotifications []models.UserNotification
	database.DB.Where("user_id = ?", userID).Find(&userNotifications)

	readMap := make(map[uint]bool)
	for _, un := range userNotifications {
		readMap[un.NotificationID] = un.IsRead
	}

	type NotificationResponse struct {
		models.Notification
		IsRead bool `json:"is_read"`
	}

	var response []NotificationResponse
	for _, n := range notifications {
		response = append(response, NotificationResponse{
			Notification: n,
			IsRead:       readMap[n.ID],
		})
	}

	c.JSON(http.StatusOK, gin.H{"notifications": response})
}

// GetAdminNotifications 获取所有通知（管理员）
func GetAdminNotifications(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var notifications []models.Notification
	var total int64

	database.DB.Model(&models.Notification{}).Count(&total)
	database.DB.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&notifications)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
	})
}

// MarkNotificationRead 标记通知为已读
func MarkNotificationRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")

	var notification models.Notification
	if result := database.DB.First(&notification, notificationID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	var userNotification models.UserNotification
	result := database.DB.Where("user_id = ? AND notification_id = ?", userID, notificationID).First(&userNotification)

	now := time.Now()
	if result.Error != nil {
		userNotification = models.UserNotification{
			UserID:         userID,
			NotificationID: notification.ID,
			IsRead:         true,
			ReadAt:         &now,
		}
		database.DB.Create(&userNotification)
	} else {
		userNotification.IsRead = true
		userNotification.ReadAt = &now
		database.DB.Save(&userNotification)
	}

	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

// MarkAllNotificationsRead 标记所有通知为已读
func MarkAllNotificationsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	var notifications []models.Notification
	database.DB.Find(&notifications)

	now := time.Now()
	for _, n := range notifications {
		var userNotification models.UserNotification
		result := database.DB.Where("user_id = ? AND notification_id = ?", userID, n.ID).First(&userNotification)

		if result.Error != nil {
			userNotification = models.UserNotification{
				UserID:         userID,
				NotificationID: n.ID,
				IsRead:         true,
				ReadAt:         &now,
			}
			database.DB.Create(&userNotification)
		} else if !userNotification.IsRead {
			userNotification.IsRead = true
			userNotification.ReadAt = &now
			database.DB.Save(&userNotification)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "已全部标记为已读"})
}

// DeleteNotification 删除通知（管理员）
func DeleteNotification(c *gin.Context) {
	notificationID := c.Param("id")

	result := database.DB.Delete(&models.Notification{}, notificationID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	database.DB.Where("notification_id = ?", notificationID).Delete(&models.UserNotification{})

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// GetUnreadCount 获取未读通知数量
func GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")

	var total int64
	query := database.DB.Model(&models.Notification{})
	if role != "admin" {
		query = query.Where("target = ? OR target = ?", "all", "admin")
	}
	query.Count(&total)

	var readCount int64
	database.DB.Model(&models.UserNotification{}).Where("user_id = ? AND is_read = ?", userID, true).Count(&readCount)

	unreadCount := total - readCount
	if unreadCount < 0 {
		unreadCount = 0
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": unreadCount})
}
