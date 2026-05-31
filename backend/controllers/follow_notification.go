package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetFollowNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	var notifications []models.FollowNotification
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications)

	type FollowNotificationResponse struct {
		models.FollowNotification
		Sender  models.User    `json:"sender"`
		Article models.Article `json:"article"`
	}

	var responses []FollowNotificationResponse
	for _, n := range notifications {
		var sender models.User
		var article models.Article
		database.DB.First(&sender, n.SenderID)
		database.DB.Preload("User").Preload("Category").First(&article, n.ArticleID)

		responses = append(responses, FollowNotificationResponse{
			FollowNotification: n,
			Sender:             sender,
			Article:            article,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": responses,
	})
}

func MarkFollowNotificationRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")

	var notification models.FollowNotification
	if result := database.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	notification.IsRead = true
	database.DB.Save(&notification)

	c.JSON(http.StatusOK, gin.H{"message": "已标记为已读"})
}

func MarkAllFollowNotificationsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	database.DB.Model(&models.FollowNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"message": "已全部标记为已读"})
}

func GetFollowNotificationUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	var count int64
	database.DB.Model(&models.FollowNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&count)

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
