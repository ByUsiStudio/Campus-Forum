package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 定义常用操作类型和模块
const (
	ActionCreate  = "create"
	ActionUpdate  = "update"
	ActionDelete  = "delete"
	ActionLogin   = "login"
	ActionLogout  = "logout"
	ActionSend    = "send"
	ActionGrant   = "grant"
	ActionRevoke  = "revoke"
	ActionBan     = "ban"
	ActionUnban   = "unban"
	ActionPin     = "pin"
	ActionUnpin   = "unpin"
	ActionApprove = "approve"
	ActionReject  = "reject"

	ModuleUser         = "user"
	ModuleArticle      = "article"
	ModuleComment      = "comment"
	ModuleNotification = "notification"
	ModulePermission   = "permission"
	ModuleReport       = "report"
	ModuleAnnouncement = "announcement"
	ModuleCategory     = "category"
	ModuleTitle        = "title"
	ModuleSiteConfig   = "site_config"
	ModuleDeletion     = "deletion"
)

// logOperation 记录操作日志的通用函数
func logOperation(c *gin.Context, action, module, details string) {
	LogAction(c, 0, action, module, details)
}

// getLogDetails 生成详细的日志描述
func getLogDetails(name string, id interface{}) string {
	switch v := id.(type) {
	case uint:
		return name + " ID: " + strconv.FormatUint(uint64(v), 10)
	case int:
		return name + " ID: " + strconv.Itoa(v)
	case string:
		return name + ": " + v
	default:
		return name
	}
}

// SendUserNotification 发送单独通知给用户
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

	// 验证目标用户是否存在
	var targetUser models.User
	if err := database.DB.First(&targetUser, req.UserID).Error; err != nil {
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

	if err := database.DB.Create(&notification).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建通知失败"})
		return
	}

	logOperation(c, ActionSend, ModuleNotification,
		"向用户发送通知 - 用户ID: "+strconv.FormatUint(uint64(req.UserID), 10)+", 标题: "+req.Title)

	c.JSON(http.StatusOK, gin.H{
		"message":      "通知发送成功",
		"notification": notification,
	})
}

// SendBatchNotifications 批量发送通知
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

	// 批量创建
	if err := database.DB.Create(&notifications).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "批量创建通知失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":       "批量通知发送成功",
		"sent_count":    len(notifications),
		"notifications": notifications,
	})
}

// GetUserNotifications 获取用户通知列表
func GetUserNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	unreadOnly := c.Query("unread_only") == "true"

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.PersonalNotification{}).Where("user_id = ?", userID)

	if unreadOnly {
		query = query.Where("is_read = ?", false)
	}

	var notifications []models.PersonalNotification
	var total int64
	var unreadCount int64

	query.Count(&total)
	database.DB.Model(&models.PersonalNotification{}).Where("user_id = ? AND is_read = ?", userID, false).Count(&unreadCount)

	query.Preload("Sender").
		Order("priority DESC, created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"unread_count":  unreadCount,
		"page":          page,
		"page_size":     pageSize,
		"total_pages":   (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetNotification 获取单个通知详情
func GetNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")

	var notification models.PersonalNotification
	if err := database.DB.Preload("Sender").Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	// 标记为已读
	if !notification.IsRead {
		now := time.Now()
		database.DB.Model(&notification).Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
		notification.IsRead = true
		notification.ReadAt = &now
	}

	c.JSON(http.StatusOK, gin.H{"notification": notification})
}

// MarkNotificationAsRead 标记通知为已读
func MarkNotificationAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")

	var notification models.PersonalNotification
	if err := database.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	if !notification.IsRead {
		now := time.Now()
		database.DB.Model(&notification).Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})
	}

	c.JSON(http.StatusOK, gin.H{"message": "标记成功"})
}

// MarkAllNotificationsAsRead 标记所有通知为已读
func MarkAllNotificationsAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")

	now := time.Now()
	database.DB.Model(&models.PersonalNotification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]interface{}{
			"is_read": true,
			"read_at": &now,
		})

	c.JSON(http.StatusOK, gin.H{"message": "全部标记已读"})
}

// DeletePersonalNotification 删除个人通知
func DeletePersonalNotification(c *gin.Context) {
	userID := c.GetUint("user_id")
	notificationID := c.Param("id")

	var notification models.PersonalNotification
	if err := database.DB.Where("id = ? AND user_id = ?", notificationID, userID).First(&notification).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "通知不存在"})
		return
	}

	database.DB.Delete(&notification)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// ClearAllNotifications 清空所有通知
func ClearAllNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")

	database.DB.Where("user_id = ?", userID).Delete(&models.PersonalNotification{})
	c.JSON(http.StatusOK, gin.H{"message": "清空成功"})
}

// AdminGetUserNotifications 管理员获取用户通知（任意用户）
func AdminGetUserNotifications(c *gin.Context) {
	targetUserID, _ := strconv.Atoi(c.Param("user_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var notifications []models.PersonalNotification
	var total int64

	query := database.DB.Model(&models.PersonalNotification{}).Where("user_id = ?", targetUserID)
	query.Count(&total)

	query.Preload("Sender").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&notifications)

	c.JSON(http.StatusOK, gin.H{
		"notifications": notifications,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
		"total_pages":   (total + int64(pageSize) - 1) / int64(pageSize),
	})
}
