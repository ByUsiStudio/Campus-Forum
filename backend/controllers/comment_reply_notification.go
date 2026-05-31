package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetCommentReplyNotifications 获取评论回复通知列表
func GetCommentReplyNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	var notifications []models.CommentReplyNotification
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&notifications)

	// 获取未读数量
	var unreadCount int64
	database.DB.Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	type NotificationWithInfo struct {
		models.CommentReplyNotification
		ReplyUser    models.User `json:"reply_user"`
		ArticleTitle string      `json:"article_title"`
		ReplyContent string      `json:"reply_content"`
	}

	var response []NotificationWithInfo
	for _, notification := range notifications {
		// 获取回复用户信息
		var reply models.Comment
		database.DB.Preload("User").First(&reply, notification.ReplyID)

		// 获取文章标题
		var article models.Article
		database.DB.First(&article, notification.ArticleID)

		response = append(response, NotificationWithInfo{
			CommentReplyNotification: notification,
			ReplyUser:                reply.User,
			ArticleTitle:             article.Title,
			ReplyContent:             reply.Content,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": response,
		"unread_count":  unreadCount,
	})
}

// MarkCommentReplyNotificationRead 标记评论回复通知为已读
func MarkCommentReplyNotificationRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")

	result := database.DB.Model(&models.CommentReplyNotification{}).
		Where("id = ? AND user_id = ?", notificationID, userID).
		Update("is_read", true)

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在或无权访问"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记已读成功"})
}

// MarkAllCommentReplyNotificationsRead 标记所有评论回复通知为已读
func MarkAllCommentReplyNotificationsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	database.DB.Model(&models.CommentReplyNotification{}).
		Where("user_id = ?", userID).
		Update("is_read", true)

	c.JSON(http.StatusOK, gin.H{"message": "全部标记已读成功"})
}
