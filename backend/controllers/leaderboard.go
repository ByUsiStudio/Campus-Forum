package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetLeaderboard 获取排行榜
func GetLeaderboard(c *gin.Context) {
	leaderboardType := c.DefaultQuery("type", "experience")
	period := c.DefaultQuery("period", "all_time")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))

	var leaderboard []models.Leaderboard
	query := database.DB.Where("type = ? AND period = ?", leaderboardType, period).
		Order("rank asc").
		Limit(limit).
		Preload("User")

	query.Find(&leaderboard)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"leaderboard": leaderboard,
			"total":       len(leaderboard),
			"type":        leaderboardType,
			"period":      period,
		},
	})
}

// GetUserRank 获取用户排名
func GetUserRank(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	leaderboardType := c.DefaultQuery("type", "experience")
	period := c.DefaultQuery("period", "all_time")

	var entry models.Leaderboard
	if err := database.DB.Where("user_id = ? AND type = ? AND period = ?", userID, leaderboardType, period).
		First(&entry).Error; err != nil {
		c.JSON(200, gin.H{
			"success": true,
			"data": gin.H{
				"rank":  0,
				"score": 0,
			},
		})
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    entry,
	})
}

// GetUserBadges 获取用户徽章
func GetUserBadges(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var badges []models.UserBadge
	database.DB.Where("user_id = ?", userID).
		Order("earned_at desc").
		Find(&badges)

	c.JSON(200, gin.H{
		"success": true,
		"data":    badges,
	})
}

// UpdateBadgeDisplay 更新徽章显示状态
func UpdateBadgeDisplay(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	badgeID := c.Param("id")

	var badge models.UserBadge
	if err := database.DB.Where("id = ? AND user_id = ?", badgeID, userID).First(&badge).Error; err != nil {
		utils.SendErrorResponse(c, 404, "徽章不存在")
		return
	}

	var req struct {
		IsDisplayed bool `json:"is_displayed"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	badge.IsDisplayed = req.IsDisplayed
	database.DB.Save(&badge)

	c.JSON(200, gin.H{
		"success": true,
		"data":    badge,
	})
}

// GrantBadge 授予徽章（管理员）
func GrantBadge(c *gin.Context) {
	var req struct {
		UserID      uint   `json:"user_id"`
		BadgeType   string `json:"badge_type"`
		BadgeName   string `json:"badge_name"`
		BadgeIcon   string `json:"badge_icon"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	badge := models.UserBadge{
		UserID:      req.UserID,
		BadgeType:   req.BadgeType,
		BadgeName:   req.BadgeName,
		BadgeIcon:   req.BadgeIcon,
		Description: req.Description,
		IsDisplayed: true,
		EarnedAt:    time.Now(),
	}

	if err := database.DB.Create(&badge).Error; err != nil {
		utils.SendErrorResponse(c, 500, "授予失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    badge,
	})
}

// RevokeBadge 撤销徽章（管理员）
func RevokeBadge(c *gin.Context) {
	badgeID := c.Param("id")

	if err := database.DB.Delete(&models.UserBadge{}, badgeID).Error; err != nil {
		utils.SendErrorResponse(c, 500, "撤销失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "撤销成功",
	})
}
