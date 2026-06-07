package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// 更新用户在线状态
func UpdateUserStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		Status string `json:"status"` // online, offline
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	now := time.Now()
	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	user.OnlineStatus = req.Status
	user.LastActiveAt = &now

	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新状态失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":        "状态更新成功",
		"online_status":  user.OnlineStatus,
		"last_active_at": user.LastActiveAt,
	})
}

// 获取所有用户状态
func GetAllUserStatuses(c *gin.Context) {
	var users []models.User
	database.DB.Select("id, username, display_name, avatar, online_status, last_active_at, created_at").
		Order("online_status DESC, last_active_at DESC").
		Find(&users)

	// 统计在线用户数
	onlineCount := 0
	for _, user := range users {
		if user.OnlineStatus == "online" {
			onlineCount++
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"users":        users,
		"online_count": onlineCount,
		"total_count":  len(users),
	})
}

// 获取单个用户状态
func GetUserStatus(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := database.DB.Select("id, username, display_name, avatar, online_status, last_active_at").
		First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})
}

// 获取在线用户列表
func GetOnlineUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	offset := (page - 1) * pageSize

	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{}).Where("online_status = ?", "online")
	query.Count(&total)
	query.Select("id, username, display_name, avatar, online_status, last_active_at, created_at").
		Order("last_active_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&users)

	c.JSON(http.StatusOK, gin.H{
		"users":       users,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// 批量更新用户状态（清理过期用户）
func CleanupUserStatuses(c *gin.Context) {
	// 将超过30分钟未活跃的用户标记为离线
	threshold := time.Now().Add(-30 * time.Minute)
	result := database.DB.Model(&models.User{}).
		Where("online_status = ? AND last_active_at < ?", "online", threshold).
		Update("online_status", "offline")

	c.JSON(http.StatusOK, gin.H{
		"message":       "清理完成",
		"updated_count": result.RowsAffected,
	})
}
