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
	clientIP := c.ClientIP()

	// 获取今天的日期
	today := time.Now().Format("2006-01-02")

	// 检查今天是否已经签到（使用字符串日期比较，提高性能）
	var existingRecord models.SignInRecord
	result := database.DB.Where("user_id = ? AND sign_in_date = ?", userID, today).First(&existingRecord)
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

	// 获取签到配置
	var config models.SignInConfig
	if err := database.DB.FirstOrCreate(&config, models.SignInConfig{ID: 1}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取签到配置失败"})
		return
	}

	if !config.Enabled {
		c.JSON(http.StatusForbidden, gin.H{"error": "签到功能已关闭"})
		return
	}

	// 计算连续签到天数
	newSignInDays := 1
	rewardPoints := config.DailyPoints

	// 获取昨天的日期
	yesterday := time.Now().Add(-24 * time.Hour).Format("2006-01-02")

	// 检查昨天是否签到
	var yesterdayRecord models.SignInRecord
	yesterdayResult := database.DB.Where("user_id = ? AND sign_in_date = ?", userID, yesterday).First(&yesterdayRecord)

	if yesterdayResult.Error == nil {
		newSignInDays = user.SignInDays + 1
	}

	// 检查连续签到奖励
	var bonusPoints int
	if newSignInDays == config.WeeklyBonus {
		bonusPoints = 5 // 连续7天额外奖励5积分
	} else if newSignInDays == config.MonthlyBonus {
		bonusPoints = 15 // 连续30天额外奖励15积分
	} else if newSignInDays == config.YearlyBonus {
		bonusPoints = 50 // 连续365天额外奖励50积分
	}

	totalReward := rewardPoints + bonusPoints

	// 创建签到记录
	signInRecord := models.SignInRecord{
		UserID:        userID,
		SignInAt:      time.Now(),
		SignInDate:    today,
		ContinuousDay: newSignInDays,
		RewardPoints:  totalReward,
		IPAddress:     clientIP,
	}

	if err := database.DB.Create(&signInRecord).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "签到失败"})
		return
	}

	// 更新用户签到统计
	maxContinuousDays := user.MaxContinuousDays
	if newSignInDays > maxContinuousDays {
		maxContinuousDays = newSignInDays
	}

	database.DB.Model(&user).Updates(map[string]interface{}{
		"sign_in_days":        newSignInDays,
		"total_sign_ins":      user.TotalSignIns + 1,
		"max_continuous_days": maxContinuousDays,
		"total_points":        user.TotalPoints + totalReward,
	})

	c.JSON(http.StatusOK, gin.H{
		"message":             "签到成功",
		"sign_in_days":        newSignInDays,
		"total_sign_ins":      user.TotalSignIns + 1,
		"reward_points":       totalReward,
		"bonus_points":        bonusPoints,
		"max_continuous_days": maxContinuousDays,
		"total_points":        user.TotalPoints + totalReward,
	})
}

// GetSignInStatus 获取用户签到状态
func GetSignInStatus(c *gin.Context) {
	userID := c.GetUint("user_id")

	// 获取今天的日期
	today := time.Now().Format("2006-01-02")

	// 检查今天是否已经签到
	var existingRecord models.SignInRecord
	result := database.DB.Where("user_id = ? AND sign_in_date = ?", userID, today).First(&existingRecord)
	hasSignedIn := result.Error == nil

	// 获取用户信息
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取用户信息失败"})
		return
	}

	// 获取签到配置
	var config models.SignInConfig
	database.DB.FirstOrCreate(&config, models.SignInConfig{ID: 1})

	// 计算本月签到天数
	monthStart := time.Now().Truncate(24*time.Hour).AddDate(0, 0, -time.Now().Day()+1)
	var monthSignInCount int64
	database.DB.Model(&models.SignInRecord{}).Where("user_id = ? AND sign_in_at >= ?", userID, monthStart).Count(&monthSignInCount)

	// 计算本周签到天数
	weekday := int(time.Now().Weekday())
	if weekday == 0 {
		weekday = 7
	}
	weekStart := time.Now().AddDate(0, 0, -(weekday - 1)).Truncate(24 * time.Hour)
	var weekSignInCount int64
	database.DB.Model(&models.SignInRecord{}).Where("user_id = ? AND sign_in_at >= ?", userID, weekStart).Count(&weekSignInCount)

	c.JSON(http.StatusOK, gin.H{
		"has_signed_in":       hasSignedIn,
		"sign_in_days":        user.SignInDays,
		"total_sign_ins":      user.TotalSignIns,
		"max_continuous_days": user.MaxContinuousDays,
		"total_points":        user.TotalPoints,
		"month_sign_in_count": monthSignInCount,
		"week_sign_in_count":  weekSignInCount,
		"config": gin.H{
			"daily_points":  config.DailyPoints,
			"weekly_bonus":  config.WeeklyBonus,
			"monthly_bonus": config.MonthlyBonus,
			"yearly_bonus":  config.YearlyBonus,
		},
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
		"records":     records,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetSignInConfig 获取签到配置（管理员）
func GetSignInConfig(c *gin.Context) {
	var config models.SignInConfig
	if err := database.DB.FirstOrCreate(&config, models.SignInConfig{ID: 1}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取签到配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"config": config,
	})
}

// UpdateSignInConfig 更新签到配置（管理员）
func UpdateSignInConfig(c *gin.Context) {
	var config models.SignInConfig
	if err := database.DB.FirstOrCreate(&config, models.SignInConfig{ID: 1}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "获取签到配置失败"})
		return
	}

	var input struct {
		DailyPoints  int   `json:"daily_points"`
		WeeklyBonus  int   `json:"weekly_bonus"`
		MonthlyBonus int   `json:"monthly_bonus"`
		YearlyBonus  int   `json:"yearly_bonus"`
		Enabled      *bool `json:"enabled"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误"})
		return
	}

	updates := map[string]interface{}{}
	if input.DailyPoints > 0 {
		updates["daily_points"] = input.DailyPoints
	}
	if input.WeeklyBonus > 0 {
		updates["weekly_bonus"] = input.WeeklyBonus
	}
	if input.MonthlyBonus > 0 {
		updates["monthly_bonus"] = input.MonthlyBonus
	}
	if input.YearlyBonus > 0 {
		updates["yearly_bonus"] = input.YearlyBonus
	}
	if input.Enabled != nil {
		updates["enabled"] = *input.Enabled
	}

	if err := database.DB.Model(&config).Updates(updates).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "更新配置失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"config":  config,
	})
}

// GetSignInRankings 获取签到排行榜
func GetSignInRankings(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	// 获取连续签到天数排行榜
	var continuousRankings []models.User
	database.DB.Order("sign_in_days DESC, total_sign_ins DESC").Limit(limit).Find(&continuousRankings)

	// 获取累计积分排行榜
	var pointsRankings []models.User
	database.DB.Order("total_points DESC, total_sign_ins DESC").Limit(limit).Find(&pointsRankings)

	c.JSON(http.StatusOK, gin.H{
		"continuous_rankings": continuousRankings,
		"points_rankings":     pointsRankings,
	})
}
