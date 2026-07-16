package controllers

import (
	"forum/models"
	"forum/repository"
	"forum/service"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	personalNotificationRepo = repository.NewPersonalNotificationRepository()
	userRepo                 = repository.NewUserRepository()
)

func SendUserNotification(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		UserID      uint   `json:"user_id" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		RelatedType string `json:"related_type"`
		RelatedID   uint   `json:"related_id"`
		Link        string `json:"link"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if _, err := userRepo.GetByID(req.UserID); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "目标用户不存在"})
		return
	}

	priority := "normal"
	if req.Priority != "" {
		priority = req.Priority
	}

	notification := models.PersonalNotification{
		SenderID:    userID,
		UserID:      req.UserID,
		Type:        req.Type,
		Title:       req.Title,
		Content:     req.Content,
		RelatedType: req.RelatedType,
		RelatedID:   req.RelatedID,
		Link:        req.Link,
		Priority:    priority,
		IsRead:      false,
	}

	if err := personalNotificationRepo.Create(&notification); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "通知发送成功", "notification": notification})
}

func SendBatchNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	var req struct {
		UserIDs     []uint `json:"user_ids" binding:"required"`
		Type        string `json:"type" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		RelatedType string `json:"related_type"`
		RelatedID   uint   `json:"related_id"`
		Link        string `json:"link"`
		Priority    string `json:"priority"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}

	if len(req.UserIDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户列表不能为空"})
		return
	}

	priority := "normal"
	if req.Priority != "" {
		priority = req.Priority
	}

	var notifications []models.PersonalNotification
	now := time.Now()

	for _, targetUserID := range req.UserIDs {
		notification := models.PersonalNotification{
			SenderID:    userID,
			UserID:      targetUserID,
			Type:        req.Type,
			Title:       req.Title,
			Content:     req.Content,
			RelatedType: req.RelatedType,
			RelatedID:   req.RelatedID,
			Link:        req.Link,
			Priority:    priority,
			IsRead:      false,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		notifications = append(notifications, notification)
	}

	for _, notification := range notifications {
		if err := personalNotificationRepo.Create(&notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建通知失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "批量通知发送成功", "sent_count": len(notifications)})
}

func GetUserNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	notifications, total, unreadCount, err := personalNotificationRepo.GetUserNotifications(userID, page, pageSize, unreadOnly)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"unread_count":  unreadCount,
		"page":          page,
		"page_size":     pageSize,
		"total_pages":   (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func GetNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	notification, err := personalNotificationRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notification": notification})
}

func MarkNotificationAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	notification, err := personalNotificationRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	if err := personalNotificationRepo.MarkAsRead(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}

func MarkAllNotificationsAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := personalNotificationRepo.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "全部标记已读"})
}

func DeletePersonalNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	notification, err := personalNotificationRepo.GetByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}
	if notification.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	if err := personalNotificationRepo.Delete(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func ClearAllNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := personalNotificationRepo.DeleteByUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "清空成功"})
}

func AdminGetUserNotifications(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	notifications, totalPages, err := service.UserNotification.AdminGetUserNotifications(uint(userID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifications, "total_pages": totalPages})
}

func GetUserUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")

	count, err := service.UserNotification.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"unread_count": count})
}
