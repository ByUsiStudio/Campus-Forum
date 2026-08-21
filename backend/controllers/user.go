package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user models.User
	if result := database.DB.First(&user, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	var userTitles []models.UserTitle
	database.DB.Where("user_id = ?", uint(id)).Preload("Title").Find(&userTitles)

	var titles []models.Title
	for _, ut := range userTitles {
		if ut.Title.IsActive {
			titles = append(titles, ut.Title)
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"qq_number":    user.QQNumber,
		"display_name": user.DisplayName,
		"avatar":       user.Avatar,
		"role":         user.Role,
		"signature":    user.Signature,
		"status":       user.Status,
		"titles":       titles,
		"created_at":   user.CreatedAt,
	})
}

func GetUserArticles(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var articles []models.Article
	var total int64

	database.DB.Model(&models.Article{}).Where("user_id = ? AND status = ?", uint(id), "published").Count(&total)
	database.DB.Where("user_id = ? AND status = ?", uint(id), "published").
		Preload("User").
		Preload("Category").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&articles)

	var currentUserID uint
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	for i := range articles {
		var likeCount int64
		database.DB.Model(&models.Like{}).Where("article_id = ?", articles[i].ID).Count(&likeCount)
		articles[i].LikeCount = int(likeCount)

		var commentCount int64
		database.DB.Model(&models.Comment{}).Where("article_id = ?", articles[i].ID).Count(&commentCount)
		articles[i].CommentCount = int(commentCount)
	}

	maskAnonymousUsers(&articles, currentUserID)

	c.JSON(http.StatusOK, gin.H{
		"articles":    articles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func GetUserFollowing(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user models.User
	if result := database.DB.First(&user, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取用户添加的好友（UserID为该用户的好友）
	var friends []models.Friend
	database.DB.Preload("Friend").Where("user_id = ? AND status = 1", uint(id)).Find(&friends)

	friendUsers := make([]models.User, 0)
	for _, friend := range friends {
		friendUsers = append(friendUsers, friend.Friend)
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendUsers})
}

func GetUserFollowers(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的用户ID"})
		return
	}

	var user models.User
	if result := database.DB.First(&user, uint(id)); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在"})
		return
	}

	// 获取添加该用户为好友的人
	var friends []models.Friend
	database.DB.Preload("User").Where("friend_id = ? AND status = 1", uint(id)).Find(&friends)

	friendUsers := make([]models.User, 0)
	for _, friend := range friends {
		friendUsers = append(friendUsers, friend.User)
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendUsers})
}
