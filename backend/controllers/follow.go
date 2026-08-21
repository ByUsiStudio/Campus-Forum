package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// SendFriendRequest 发送好友请求
func SendFriendRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FriendID uint   `json:"friend_id"`
		Message  string `json:"message"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	if userID == req.FriendID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能添加自己为好友"})
		return
	}

	// 检查是否已存在好友关系
	var existingFriend models.Friend
	result := database.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, req.FriendID, req.FriendID, userID).First(&existingFriend)
	if result.Error == nil && existingFriend.Status == 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已添加该用户为好友"})
		return
	}

	// 检查是否已有待处理的请求
	var existingRequest models.FriendRequest
	result = database.DB.Where("(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)", userID, req.FriendID, req.FriendID, userID).First(&existingRequest)
	if result.Error == nil && existingRequest.Status == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已存在待处理的好友请求"})
		return
	}

	friendRequest := models.FriendRequest{
		SenderID:   userID,
		ReceiverID: req.FriendID,
		Message:    req.Message,
		Status:     0,
	}

	if err := database.DB.Create(&friendRequest).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "发送好友请求失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "好友请求已发送"})
}

// AcceptFriendRequest 同意好友请求
func AcceptFriendRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	requestID := c.GetUint("id")

	var friendRequest models.FriendRequest
	result := database.DB.First(&friendRequest, requestID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "好友请求不存在"})
		return
	}

	if friendRequest.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权处理此请求"})
		return
	}

	if friendRequest.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求状态不正确"})
		return
	}

	// 更新请求状态
	friendRequest.Status = 1
	database.DB.Save(&friendRequest)

	// 创建双向好友关系
	friend1 := models.Friend{
		UserID:      userID,
		FriendID:    friendRequest.SenderID,
		DisplayName: "",
		Status:      1,
	}
	friend2 := models.Friend{
		UserID:      friendRequest.SenderID,
		FriendID:    userID,
		DisplayName: "",
		Status:      1,
	}

	database.DB.Create(&friend1)
	database.DB.Create(&friend2)

	c.JSON(http.StatusOK, gin.H{"message": "已同意好友请求"})
}

// RejectFriendRequest 拒绝好友请求
func RejectFriendRequest(c *gin.Context) {
	userID := c.GetUint("user_id")
	requestID := c.GetUint("id")

	var friendRequest models.FriendRequest
	result := database.DB.First(&friendRequest, requestID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "好友请求不存在"})
		return
	}

	if friendRequest.ReceiverID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权处理此请求"})
		return
	}

	if friendRequest.Status != 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求状态不正确"})
		return
	}

	friendRequest.Status = 2
	database.DB.Save(&friendRequest)

	c.JSON(http.StatusOK, gin.H{"message": "已拒绝好友请求"})
}

// DeleteFriend 删除好友
func DeleteFriend(c *gin.Context) {
	userID := c.GetUint("user_id")
	friendID := c.GetUint("id")

	// 删除双向好友关系
	database.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).Delete(&models.Friend{})
	database.DB.Where("user_id = ? AND friend_id = ?", friendID, userID).Delete(&models.Friend{})

	c.JSON(http.StatusOK, gin.H{"message": "已删除好友"})
}

// GetFriendList 获取好友列表
func GetFriendList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var friends []models.Friend
	database.DB.Where("user_id = ? AND status = 1", userID).Preload("Friend").Find(&friends)

	friendUsers := make([]models.User, 0)
	for _, friend := range friends {
		friendUsers = append(friendUsers, friend.Friend)
	}

	c.JSON(http.StatusOK, gin.H{"friends": friendUsers})
}

// GetFriendRequests 获取好友请求列表
func GetFriendRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	var requests []models.FriendRequest
	database.DB.Where("receiver_id = ? AND status = 0", userID).Preload("Sender").Find(&requests)

	c.JSON(http.StatusOK, gin.H{"requests": requests})
}

// GetSentFriendRequests 获取已发送的好友请求列表
func GetSentFriendRequests(c *gin.Context) {
	userID := c.GetUint("user_id")

	var requests []models.FriendRequest
	database.DB.Where("sender_id = ? AND status = 0", userID).Preload("Receiver").Find(&requests)

	c.JSON(http.StatusOK, gin.H{"sent_requests": requests})
}

// UpdateFriendDisplayName 更新好友备注名
func UpdateFriendDisplayName(c *gin.Context) {
	userID := c.GetUint("user_id")
	friendID := c.GetUint("id")

	var req struct {
		DisplayName string `json:"display_name"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求参数错误"})
		return
	}

	result := database.DB.Model(&models.Friend{}).Where("user_id = ? AND friend_id = ?", userID, friendID).Update("display_name", req.DisplayName)
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "好友关系不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "备注名已更新"})
}

// CheckFriendStatus 检查好友状态
func CheckFriendStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID := c.GetUint("id")

	var friend models.Friend
	result := database.DB.Where("user_id = ? AND friend_id = ? AND status = 1", userID, targetID).First(&friend)
	isFriend := result.Error == nil

	// 获取目标用户信息
	var targetUser models.User
	database.DB.First(&targetUser, targetID)

	c.JSON(http.StatusOK, gin.H{
		"is_friend": isFriend,
		"user":      targetUser,
	})
}

// GetMutualFriends 获取共同好友
func GetMutualFriends(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID := c.GetUint("id")

	// 获取当前用户的好友列表
	var myFriends []models.Friend
	database.DB.Where("user_id = ? AND status = 1", userID).Find(&myFriends)

	myFriendIDs := make(map[uint]bool)
	for _, f := range myFriends {
		myFriendIDs[f.FriendID] = true
	}

	// 获取目标用户的好友列表
	var targetFriends []models.Friend
	database.DB.Where("user_id = ? AND status = 1", targetID).Find(&targetFriends)

	var mutual []models.User
	for _, f := range targetFriends {
		if myFriendIDs[f.FriendID] {
			var user models.User
			database.DB.First(&user, f.FriendID)
			mutual = append(mutual, user)
		}
	}

	c.JSON(http.StatusOK, gin.H{"mutual_friends": mutual})
}
