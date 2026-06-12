package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// SignIn 用户签到
func SignIn(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取今天的日期（年月日）
	today := time.Now().Format("2006-01-02")

	// 检查今天是否已经签到
	var existingRecord models.SignInRecord
	result := database.DB.Where("user_id = ? AND DATE(sign_in_at) = ?", userID, today).First(&existingRecord)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "今天已经签到过了"})
		return
	}

	// 获取用户信息
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 获取昨天的日期
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	// 检查昨天是否签到（判断连续签到）
	var yesterdayRecord models.SignInRecord
	yesterdayResult := database.DB.Where("user_id = ? AND DATE(sign_in_at) = ?", userID, yesterday).First(&yesterdayRecord)

	// 计算连续签到天数
	newSignInDays := 1
	if yesterdayResult.Error == nil {
		newSignInDays = user.SignInDays + 1
	}

	// 创建签到记录
	signInRecord := models.SignInRecord{
		UserID:   userID,
		SignInAt: time.Now(),
	}

	if err := database.DB.Create(&signInRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签到失败"})
		return
	}

	// 更新用户签到统计
	database.DB.Model(&user).Updates(map[string]interface{}{
		"sign_in_days":   newSignInDays,
		"total_sign_ins": user.TotalSignIns + 1,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":        "签到成功",
		"sign_in_days":   newSignInDays,
		"total_sign_ins": user.TotalSignIns + 1,
	})
}

// GetSignInStatus 获取用户签到状态
func GetSignInStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取今天的日期
	today := time.Now().Format("2006-01-02")

	// 检查今天是否已经签到
	var existingRecord models.SignInRecord
	result := database.DB.Where("user_id = ? AND DATE(sign_in_at) = ?", userID, today).First(&existingRecord)
	hasSignedIn := result.Error == nil

	// 获取用户信息
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"has_signed_in":  hasSignedIn,
		"sign_in_days":   user.SignInDays,
		"total_sign_ins": user.TotalSignIns,
	})
}

// GetSignInHistory 获取用户签到历史
func GetSignInHistory(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "30"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 30
	}

	offset := (page - 1) * pageSize

	var records []models.SignInRecord
	var total int64

	database.DB.Model(&models.SignInRecord{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).Order("sign_in_at DESC").Offset(offset).Limit(pageSize).Find(&records)

	c.JSON(http.StatusOK, gin.H{
		"records":    records,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}