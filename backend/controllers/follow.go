package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func FollowUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.Atoi(c.Param("target_id"))

	err := service.Follow.FollowUser(userID, uint(targetID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "关注成功"})
}

func UnfollowUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.Atoi(c.Param("target_id"))

	err := service.Follow.UnfollowUser(userID, uint(targetID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消关注成功"})
}

func GetFollowers(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	followers, totalPages, err := service.Follow.GetFollowers(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"followers": followers, "total_pages": totalPages})
}

func GetFollowing(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	following, totalPages, err := service.Follow.GetFollowing(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"following": following, "total_pages": totalPages})
}

func CheckFollowing(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.Atoi(c.Param("target_id"))

	isFollowing, err := service.Follow.CheckFollowing(userID, uint(targetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_following": isFollowing})
}
