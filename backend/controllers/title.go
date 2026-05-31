package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetAllTitles(c *gin.Context) {
	var titles []models.Title
	database.DB.Order("created_at DESC").Find(&titles)

	c.JSON(http.StatusOK, gin.H{"titles": titles})
}

func CreateTitle(c *gin.Context) {
	var input struct {
		Name        string `json:"name" binding:"required"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供头衔名称"})
		return
	}

	title := models.Title{
		Name:        input.Name,
		Description: input.Description,
		Color:       input.Color,
		Icon:        input.Icon,
		IsActive:    true,
	}

	database.DB.Create(&title)

	c.JSON(http.StatusOK, gin.H{
		"message": "头衔创建成功",
		"title":   title,
	})
}

func UpdateTitle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的头衔ID"})
		return
	}

	var title models.Title
	if result := database.DB.First(&title, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "头衔不存在"})
		return
	}

	var input struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		Color       string `json:"color"`
		Icon        string `json:"icon"`
		IsActive    bool   `json:"is_active"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Name != "" {
		title.Name = input.Name
	}
	if input.Description != "" {
		title.Description = input.Description
	}
	if input.Color != "" {
		title.Color = input.Color
	}
	if input.Icon != "" {
		title.Icon = input.Icon
	}
	title.IsActive = input.IsActive

	database.DB.Save(&title)

	c.JSON(http.StatusOK, gin.H{
		"message": "头衔更新成功",
		"title":   title,
	})
}

func DeleteTitle(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的头衔ID"})
		return
	}

	database.DB.Delete(&models.Title{}, uint(id))
	database.DB.Where("title_id = ?", uint(id)).Delete(&models.UserTitle{})

	c.JSON(http.StatusOK, gin.H{"message": "头衔删除成功"})
}

func GrantTitle(c *gin.Context) {
	var input struct {
		UserID  uint   `json:"user_id" binding:"required"`
		TitleID uint   `json:"title_id" binding:"required"`
		Reason  string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户ID和头衔ID"})
		return
	}

	var user models.User
	if result := database.DB.First(&user, input.UserID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var title models.Title
	if result := database.DB.First(&title, input.TitleID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "头衔不存在"})
		return
	}

	var existingUserTitle models.UserTitle
	if result := database.DB.Where("user_id = ? AND title_id = ?", input.UserID, input.TitleID).First(&existingUserTitle); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户已拥有该头衔"})
		return
	}

	userTitle := models.UserTitle{
		UserID:  input.UserID,
		TitleID: input.TitleID,
		Reason:  input.Reason,
	}

	database.DB.Create(&userTitle)

	c.JSON(http.StatusOK, gin.H{
		"message":    "头衔授予成功",
		"user_title": userTitle,
	})
}

func RevokeTitle(c *gin.Context) {
	var input struct {
		UserID  uint `json:"user_id" binding:"required"`
		TitleID uint `json:"title_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请提供用户ID和头衔ID"})
		return
	}

	result := database.DB.Where("user_id = ? AND title_id = ?", input.UserID, input.TitleID).Delete(&models.UserTitle{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "用户未拥有该头衔"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "头衔撤销成功"})
}

func GetUserTitles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var userTitles []models.UserTitle
	database.DB.Where("user_id = ?", uint(id)).Preload("Title").Find(&userTitles)

	c.JSON(http.StatusOK, gin.H{"titles": userTitles})
}
