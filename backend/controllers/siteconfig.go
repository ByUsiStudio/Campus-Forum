package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSiteConfig(c *gin.Context) {
	var config models.SiteConfig
	result := database.DB.First(&config)
	if result.Error != nil {
		// 如果不存在，创建一个默认配置
		config = models.SiteConfig{
			SiteTitle: "校园论坛 - 分享与交流",
		}
		database.DB.Create(&config)
	}
	c.JSON(http.StatusOK, config)
}

func UpdateSiteConfig(c *gin.Context) {
	var input struct {
		SiteTitle string `json:"site_title"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var config models.SiteConfig
	database.DB.First(&config)

	if input.SiteTitle != "" {
		config.SiteTitle = input.SiteTitle
	}

	database.DB.Save(&config)
	c.JSON(http.StatusOK, gin.H{"message": "更新成功", "config": config})
}
