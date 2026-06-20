package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetUserLevel 获取用户等级信息
func GetUserLevel(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var level models.UserLevel
	if err := database.DB.Where("user_id = ?", userID).First(&level).Error; err != nil {
		// 自动创建初始等级
		level = models.UserLevel{
			UserID:     userID,
			Level:      1,
			Experience: 0,
			NextLevel:  100,
			Title:      "新手",
		}
		database.DB.Create(&level)
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    level,
	})
}

// GetUserExperienceRecords 获取用户经验记录
func GetUserExperienceRecords(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	var records []models.ExperienceRecord
	var total int64
	database.DB.Model(&models.ExperienceRecord{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Offset(offset).Limit(limit).
		Find(&records)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"records": records,
			"total":   total,
			"page":    page,
			"limit":   limit,
		},
	})
}

// GetUserAchievements 获取用户成就列表
func GetUserAchievements(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var achievements []models.UserAchievement
	database.DB.Where("user_id = ?", userID).
		Preload("Achievement").
		Order("unlocked_at desc").
		Find(&achievements)

	c.JSON(200, gin.H{
		"success": true,
		"data":    achievements,
	})
}

// GetAllAchievements 获取所有成就定义
func GetAllAchievements(c *gin.Context) {
	var achievements []models.Achievement
	database.DB.Find(&achievements)

	c.JSON(200, gin.H{
		"success": true,
		"data":    achievements,
	})
}

// GetLevelConfig 获取等级配置
func GetLevelConfig(c *gin.Context) {
	var configs []models.LevelConfig
	database.DB.Order("level asc").Find(&configs)

	c.JSON(200, gin.H{
		"success": true,
		"data":    configs,
	})
}

// CreateLevelConfig 创建等级配置（管理员）
func CreateLevelConfig(c *gin.Context) {
	var config models.LevelConfig
	if err := c.ShouldBindJSON(&config); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	if err := database.DB.Create(&config).Error; err != nil {
		utils.SendErrorResponse(c, 500, "创建失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    config,
	})
}

// UpdateLevelConfig 更新等级配置（管理员）
func UpdateLevelConfig(c *gin.Context) {
	configID := c.Param("id")

	var config models.LevelConfig
	if err := database.DB.First(&config, configID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "配置不存在")
		return
	}

	if err := c.ShouldBindJSON(&config); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	database.DB.Save(&config)

	c.JSON(200, gin.H{
		"success": true,
		"data":    config,
	})
}

// CreateAchievement 创建成就（管理员）
func CreateAchievement(c *gin.Context) {
	var achievement models.Achievement
	if err := c.ShouldBindJSON(&achievement); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	if err := database.DB.Create(&achievement).Error; err != nil {
		utils.SendErrorResponse(c, 500, "创建失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    achievement,
	})
}

// UpdateAchievement 更新成就（管理员）
func UpdateAchievement(c *gin.Context) {
	achievementID := c.Param("id")

	var achievement models.Achievement
	if err := database.DB.First(&achievement, achievementID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "成就不存在")
		return
	}

	if err := c.ShouldBindJSON(&achievement); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	database.DB.Save(&achievement)

	c.JSON(200, gin.H{
		"success": true,
		"data":    achievement,
	})
}

// DeleteAchievement 删除成就（管理员）
func DeleteAchievement(c *gin.Context) {
	achievementID := c.Param("id")

	// 删除用户成就关联
	database.DB.Where("achievement_id = ?", achievementID).Delete(&models.UserAchievement{})

	// 删除成就
	if err := database.DB.Delete(&models.Achievement{}, achievementID).Error; err != nil {
		utils.SendErrorResponse(c, 500, "删除失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "删除成功",
	})
}
