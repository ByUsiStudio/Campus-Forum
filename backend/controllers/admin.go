package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetStatistics 获取统计数据
func GetStatistics(c *gin.Context) {
	var userCount int64
	var articleCount int64
	var commentCount int64
	var viewCount int64

	database.DB.Model(&models.User{}).Count(&userCount)
	database.DB.Model(&models.Article{}).Where("status = ?", "published").Count(&articleCount)
	database.DB.Model(&models.Comment{}).Count(&commentCount)

	// 获取总浏览量
	var articles []models.Article
	database.DB.Select("view_count").Find(&articles)
	for _, article := range articles {
		viewCount += int64(article.ViewCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"user_count":    userCount,
		"article_count": articleCount,
		"comment_count": commentCount,
		"view_count":    viewCount,
	})
}

// GetAllUsers 获取所有用户（管理员）
func GetAllUsers(c *gin.Context) {
	var users []models.User
	database.DB.Order("created_at DESC").Find(&users)

	type UserResponse struct {
		ID          uint           `json:"id"`
		Username    string         `json:"username"`
		DisplayName string         `json:"display_name"`
		Avatar      string         `json:"avatar"`
		QQNumber    string         `json:"qq_number"`
		Role        string         `json:"role"`
		Status      string         `json:"status"`
		CreatedAt   time.Time      `json:"created_at"`
		Titles      []models.Title `json:"titles"`
	}

	var response []UserResponse
	for _, user := range users {
		var userTitles []models.UserTitle
		database.DB.Where("user_id = ?", user.ID).Preload("Title").Find(&userTitles)

		var titles []models.Title
		for _, ut := range userTitles {
			if ut.Title.IsActive {
				titles = append(titles, ut.Title)
			}
		}

		resp := UserResponse{
			ID:          user.ID,
			Username:    user.Username,
			DisplayName: user.DisplayName,
			Avatar:      user.Avatar,
			QQNumber:    user.QQNumber,
			Role:        user.Role,
			Status:      user.Status,
			CreatedAt:   user.CreatedAt,
			Titles:      titles,
		}
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, gin.H{"users": response})
}

// UpdateUserRole 更新用户角色
func UpdateUserRole(c *gin.Context) {
	userID := c.Param("id")

	var input struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Role != "admin" && input.Role != "user" && input.Role != "system" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的角色"})
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Update("role", input.Role)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

// DeleteUser 删除用户
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	// 不允许删除自己
	if c.GetUint("user_id") == uint(parseUint(userID)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能删除自己"})
		return
	}

	result := database.DB.Delete(&models.User{}, userID)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// BanUser 封禁用户
func BanUser(c *gin.Context) {
	userID := c.Param("id")

	// 不允许封禁自己
	if c.GetUint("user_id") == uint(parseUint(userID)) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能封禁自己"})
		return
	}

	var input struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"status":     "banned",
		"ban_reason": input.Reason,
	})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "封禁成功"})
}

// UnbanUser 解封用户
func UnbanUser(c *gin.Context) {
	userID := c.Param("id")

	result := database.DB.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"status":     "normal",
		"ban_reason": "",
	})

	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "解封成功"})
}

// GetAllArticles 获取所有文章（管理员）
func GetAllArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Article{}).Preload("User").Preload("Category")

	if status != "" {
		query = query.Where("status = ?", status)
	}

	var articles []models.Article
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"articles":    articles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetAllComments 获取所有评论（管理员）
func GetAllComments(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var comments []models.Comment
	var total int64

	database.DB.Model(&models.Comment{}).Preload("User").Preload("Article").Count(&total)
	database.DB.Preload("User").Preload("Article").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	c.JSON(http.StatusOK, gin.H{
		"comments":    comments,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// DeleteCommentAdmin 管理员删除评论
func DeleteCommentAdmin(c *gin.Context) {
	commentID := c.Param("id")

	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 如果是回复评论，减少父评论的回复数量
	if comment.ParentID != nil {
		database.DB.Model(&models.Comment{}).Where("id = ?", *comment.ParentID).UpdateColumn("reply_count", gorm.Expr("reply_count - 1"))
	}

	result := database.DB.Delete(&comment)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

// UpdateArticleStatus 更新文章状态
func UpdateArticleStatus(c *gin.Context) {
	articleID := c.Param("id")

	var input struct {
		Status string `json:"status" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Status != "published" && input.Status != "pending" && input.Status != "deleted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的状态"})
		return
	}

	result := database.DB.Model(&models.Article{}).Where("id = ?", articleID).Update("status", input.Status)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}
