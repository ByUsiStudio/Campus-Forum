package controllers

import (
	"forum/database"
	"forum/models"
	"forum/repository"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

var systemLogRepo = repository.NewSystemLogRepository()

// LogAction 记录操作日志的辅助函数
func LogAction(c *gin.Context, userID uint, action, module, details string) {
	if userID == 0 {
		userID = c.GetUint("user_id")
	}
	if userID == 0 {
		return // 未登录用户不记录
	}

	ip := c.ClientIP()
	if ip == "" {
		ip = c.RemoteIP()
	}

	log := models.SystemLog{
		UserID:   userID,
		Action:   action,
		Module:   module,
		Details:  details,
		IP:       ip,
	}

	if err := database.DB.Create(&log).Error; err != nil {
		// 日志记录失败不影响主流程
	}
}

// GetSystemLogs 获取系统日志（管理员）
func GetSystemLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	module := c.Query("module")
	userID, _ := strconv.Atoi(c.Query("user_id"))

	var logs []models.SystemLog
	var total int64

	query := database.DB.Model(&models.SystemLog{})

	if module != "" {
		query = query.Where("module = ?", module)
	}
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}

	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&logs).Error

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":         logs,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetMyLogs 获取当前用户的操作日志
func GetMyLogs(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, total, err := systemLogRepo.GetLogsByUser(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取日志失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"logs":         logs,
		"total":        total,
		"page":         page,
		"page_size":    pageSize,
		"total_pages":  (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetLogModules 获取所有日志模块
func GetLogModules(c *gin.Context) {
	var modules []string
	database.DB.Model(&models.SystemLog{}).
		Distinct().
		Pluck("module", &modules)

	c.JSON(http.StatusOK, gin.H{"modules": modules})
}

// DeleteOldLogs 删除旧日志
func DeleteOldLogs(c *gin.Context) {
	days, _ := strconv.Atoi(c.DefaultQuery("days", "90"))
	if days < 7 {
		days = 7
	}

	if err := systemLogRepo.DeleteOldLogs(days); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除日志失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "旧日志删除成功"})
}
