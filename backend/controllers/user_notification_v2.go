package controllers

import (
	"forum/models"
	"forum/repository"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

var (
	personalNotificationRepo = repository.NewPersonalNotificationRepository()
	userRepo                 = repository.NewUserRepository()
)

// SendUserNotification 发送单独通知给用户 (使用 Repository 版本)
func SendUserNotificationV2(c *gin.Context) {
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

	// 验证目标用户是否存在
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

	c.JSON(http.StatusOK, gin.H{
		"message":      "通知发送成功",
		"notification": notification,
	})
}

// SendBatchNotificationsV2 批量发送通知 (使用 Repository 版本)
func SendBatchNotificationsV2(c *gin.Context) {
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

	// 批量创建
	for _, notification := range notifications {
		if err := personalNotificationRepo.Create(&notification); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建通知失败"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "批量通知发送成功",
		"sent_count":    len(notifications),
		"notifications": notifications,
	})
}

// GetUserNotificationsV2 获取用户通知列表 (使用 Repository 版本)
func GetUserNotificationsV2(c *gin.Context) {
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

// MarkNotificationAsReadV2 标记通知为已读 (使用 Repository 版本)
func MarkNotificationAsReadV2(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	// 验证通知是否属于该用户
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

// MarkAllNotificationsAsReadV2 标记所有通知为已读 (使用 Repository 版本)
func MarkAllNotificationsAsReadV2(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := personalNotificationRepo.MarkAllAsRead(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "标记失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "全部标记已读"})
}

// DeletePersonalNotificationV2 删除个人通知 (使用 Repository 版本)
func DeletePersonalNotificationV2(c *gin.Context) {
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

// ClearAllNotificationsV2 清空所有通知 (使用 Repository 版本)
func ClearAllNotificationsV2(c *gin.Context) {
	userID := c.GetUint("user_id")

	if err := personalNotificationRepo.DeleteByUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "清空失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "清空成功"})
}
