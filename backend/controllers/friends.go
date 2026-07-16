package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SendFriendRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.Atoi(c.Param("id"))

	err := service.Follow.SendFriendRequest(userID, uint(targetID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友请求已发送"})
}

func AcceptFriendRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	requestID, _ := strconv.Atoi(c.Param("id"))

	err := service.Follow.AcceptFriendRequest(userID, uint(requestID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友请求已接受"})
}

func RejectFriendRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	requestID, _ := strconv.Atoi(c.Param("id"))

	err := service.Follow.RejectFriendRequest(userID, uint(requestID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友请求已拒绝"})
}

func DeleteFriend(c *gin.Context) {
	userID := c.GetUint("user_id")
	friendID, _ := strconv.Atoi(c.Param("id"))

	err := service.Follow.DeleteFriend(userID, uint(friendID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "已删除好友"})
}

func GetFriendList(c *gin.Context) {
	userID := c.GetUint("user_id")

	friends, err := service.Follow.GetFriendList(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

func GetFriendRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	requests, err := service.Follow.GetFriendRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

func GetSentFriendRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	requests, err := service.Follow.GetSentFriendRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sent_requests": requests})
}

func UpdateFriendDisplayName(c *gin.Context) {
	userID := c.GetUint("user_id")
	friendID, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		DisplayName string `json:"display_name"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Follow.UpdateFriendDisplayName(userID, uint(friendID), input.DisplayName)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友昵称已更新"})
}

func CheckFriendStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.Atoi(c.Param("id"))

	status, err := service.Follow.CheckFriendStatus(userID, uint(targetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status})
}

func GetMutualFriends(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID, _ := strconv.Atoi(c.Param("id"))

	friends, err := service.Follow.GetMutualFriends(userID, uint(targetID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"mutual_friends": friends, "count": len(friends)})
}
