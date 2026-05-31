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

	var articles []models.Article
	database.DB.Where("user_id = ? AND status = ?", uint(id), "published").
		Preload("User").
		Order("created_at DESC").
		Find(&articles)

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
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

	var follows []models.Follow
	database.DB.Where("follower_id = ?", uint(id)).Find(&follows)

	followingIDs := make([]uint, 0)
	for _, follow := range follows {
		followingIDs = append(followingIDs, follow.FollowingID)
	}

	var users []models.User
	if len(followingIDs) > 0 {
		database.DB.Where("id IN ?", followingIDs).Find(&users)
	}

	c.JSON(http.StatusOK, gin.H{"following": users})
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

	var follows []models.Follow
	database.DB.Where("following_id = ?", uint(id)).Find(&follows)

	followerIDs := make([]uint, 0)
	for _, follow := range follows {
		followerIDs = append(followerIDs, follow.FollowerID)
	}

	var users []models.User
	if len(followerIDs) > 0 {
		database.DB.Where("id IN ?", followerIDs).Find(&users)
	}

	c.JSON(http.StatusOK, gin.H{"followers": users})
}
