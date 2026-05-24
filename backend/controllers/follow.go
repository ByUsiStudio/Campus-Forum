package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// FollowUser 关注用户
func FollowUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	followingID := c.GetUint("id")

	if userID == followingID {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不能关注自己"})
		return
	}

	var follow models.Follow
	result := database.DB.Where("follower_id = ? AND following_id = ?", userID, followingID).First(&follow)
	if result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已关注该用户"})
		return
	}

	follow = models.Follow{
		FollowerID:  userID,
		FollowingID: followingID,
	}

	if err := database.DB.Create(&follow).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "关注失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "关注成功"})
}

// UnfollowUser 取消关注用户
func UnfollowUser(c *gin.Context) {
	userID := c.GetUint("user_id")
	followingID := c.GetUint("id")

	result := database.DB.Where("follower_id = ? AND following_id = ?", userID, followingID).Delete(&models.Follow{})
	if result.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "未关注该用户"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消关注成功"})
}

// GetFollowingList 获取关注列表
func GetFollowingList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var follows []models.Follow
	database.DB.Where("follower_id = ?", userID).Find(&follows)

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

// GetFollowerList 获取粉丝列表
func GetFollowerList(c *gin.Context) {
	userID := c.GetUint("user_id")

	var follows []models.Follow
	database.DB.Where("following_id = ?", userID).Find(&follows)

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

// CheckFollowStatus 检查关注状态
func CheckFollowStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
	targetID := c.GetUint("id")

	var follow models.Follow
	result := database.DB.Where("follower_id = ? AND following_id = ?", userID, targetID).First(&follow)
	isFollowing := result.Error == nil

	// 检查是否被对方关注（互相关注）
	var reverseFollow models.Follow
	result = database.DB.Where("follower_id = ? AND following_id = ?", targetID, userID).First(&reverseFollow)
	isFollowed := result.Error == nil

	// 获取目标用户信息
	var targetUser models.User
	database.DB.First(&targetUser, targetID)

	c.JSON(http.StatusOK, gin.H{
		"is_following":   isFollowing,
		"is_followed":    isFollowed,
		"mutual":         isFollowing && isFollowed,
		"following_user": targetUser,
	})
}

// GetMutualFriends 获取互相关注的用户
func GetMutualFriends(c *gin.Context) {
	userID := c.GetUint("user_id")

	var following []models.Follow
	database.DB.Where("follower_id = ?", userID).Find(&following)

	followingIDs := make(map[uint]bool)
	for _, f := range following {
		followingIDs[f.FollowingID] = true
	}

	var followers []models.Follow
	database.DB.Where("following_id = ?", userID).Find(&followers)

	var mutual []models.User
	for _, f := range followers {
		if followingIDs[f.FollowerID] {
			var user models.User
			database.DB.First(&user, f.FollowerID)
			mutual = append(mutual, user)
		}
	}

	c.JSON(http.StatusOK, gin.H{"mutual": mutual})
}
