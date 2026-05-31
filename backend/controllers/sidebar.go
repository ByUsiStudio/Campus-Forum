package controllers

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetSidebarConfig(c *gin.Context) {
	var config models.SidebarConfig
	result := database.DB.First(&config)

	if result.Error != nil {
		// 返回默认配置
		defaultItems := []models.SidebarItem{
			{Title: "首页", Link: "/", Icon: "🏠"},
			{Title: "表白墙", Link: "/category/1", Icon: "❤️"},
		}
		c.JSON(http.StatusOK, gin.H{
			"items": defaultItems,
		})
		return
	}

	var items []models.SidebarItem
	json.Unmarshal([]byte(config.Items), &items)

	c.JSON(http.StatusOK, gin.H{
		"items": items,
	})
}

func UpdateSidebarConfig(c *gin.Context) {
	var input struct {
		Items []models.SidebarItem `json:"items" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemsJSON, err := json.Marshal(input.Items)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "序列化失败"})
		return
	}

	var config models.SidebarConfig
	result := database.DB.First(&config)
	if result.Error != nil {
		config.Items = string(itemsJSON)
		database.DB.Create(&config)
	} else {
		config.Items = string(itemsJSON)
		database.DB.Save(&config)
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
