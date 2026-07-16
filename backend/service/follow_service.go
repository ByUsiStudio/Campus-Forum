package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type FollowService struct{}

var Follow = &FollowService{}

func (s *FollowService) SendFriendRequest(fromUserID, toUserID uint) error {
	if fromUserID == toUserID {
		return utils.NewError("不能添加自己为好友", 400)
	}

	var request models.FriendRequest
	result := database.DB.Where("sender_id = ? AND receiver_id = ?", fromUserID, toUserID).First(&request)
	if result.Error == nil {
		return utils.NewError("好友请求已发送", 400)
	}

	request = models.FriendRequest{
		SenderID:   fromUserID,
		ReceiverID: toUserID,
		Status:     0,
	}
	database.DB.Create(&request)

	return nil
}

func (s *FollowService) AcceptFriendRequest(userID, requestID uint) error {
	var request models.FriendRequest
	if result := database.DB.First(&request, requestID); result.Error != nil {
		return utils.NewError("请求不存在", 404)
	}

	if request.ReceiverID != userID {
		return utils.NewError("无权处理该请求", 403)
	}

	request.Status = 1
	database.DB.Save(&request)

	friend1 := models.Friend{
		UserID:   request.SenderID,
		FriendID: request.ReceiverID,
	}
	friend2 := models.Friend{
		UserID:   request.ReceiverID,
		FriendID: request.SenderID,
	}

	database.DB.Create(&friend1)
	database.DB.Create(&friend2)

	return nil
}

func (s *FollowService) RejectFriendRequest(userID, requestID uint) error {
	var request models.FriendRequest
	if result := database.DB.First(&request, requestID); result.Error != nil {
		return utils.NewError("请求不存在", 404)
	}

	if request.ReceiverID != userID {
		return utils.NewError("无权处理该请求", 403)
	}

	request.Status = 2
	database.DB.Save(&request)

	return nil
}

func (s *FollowService) DeleteFriend(userID, friendID uint) error {
	var friend models.Friend
	result := database.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID, friendID, friendID, userID).First(&friend)
	if result.Error != nil {
		return utils.NewError("不是好友关系", 400)
	}

	database.DB.Delete(&friend)

	return nil
}

func (s *FollowService) GetFriendList(userID uint) ([]models.Friend, error) {
	var friends []models.Friend
	err := database.DB.Where("user_id = ?", userID).Preload("Friend").Find(&friends).Error
	return friends, err
}

func (s *FollowService) GetFriendRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := database.DB.Where("receiver_id = ? AND status = ?", userID, 0).Preload("Sender").Find(&requests).Error
	return requests, err
}

func (s *FollowService) GetSentFriendRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := database.DB.Where("sender_id = ? AND status = ?", userID, 0).Preload("Receiver").Find(&requests).Error
	return requests, err
}

func (s *FollowService) UpdateFriendDisplayName(userID, friendID uint, displayName string) error {
	var friend models.Friend
	result := database.DB.Where("user_id = ? AND friend_id = ?", userID, friendID).First(&friend)
	if result.Error != nil {
		return utils.NewError("不是好友关系", 400)
	}

	friend.DisplayName = displayName
	database.DB.Save(&friend)

	return nil
}

func (s *FollowService) CheckFriendStatus(userID1, userID2 uint) (string, error) {
	var friend models.Friend
	result := database.DB.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)", userID1, userID2, userID2, userID1).First(&friend)
	if result.Error == nil {
		return "friend", nil
	}

	var request models.FriendRequest
	result = database.DB.Where("sender_id = ? AND receiver_id = ?", userID1, userID2).First(&request)
	if result.Error == nil {
		return "pending_sent", nil
	}

	result = database.DB.Where("sender_id = ? AND receiver_id = ?", userID2, userID1).First(&request)
	if result.Error == nil {
		return "pending_received", nil
	}

	return "none", nil
}

func (s *FollowService) GetMutualFriends(userID1, userID2 uint) ([]models.User, error) {
	var mutualFriends []models.User
	err := database.DB.Table("friends f1").
		Joins("JOIN friends f2 ON f1.friend_id = f2.friend_id").
		Joins("JOIN users u ON f1.friend_id = u.id").
		Where("f1.user_id = ? AND f2.user_id = ?", userID1, userID2).
		Select("u.*").
		Scan(&mutualFriends).Error

	return mutualFriends, err
}

func (s *FollowService) FollowUser(followerID, followedID uint) error {
	if followerID == followedID {
		return utils.NewError("不能关注自己", 400)
	}

	var follow models.UserFollow
	result := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow)
	if result.Error == nil {
		return utils.NewError("已关注该用户", 400)
	}

	follow = models.UserFollow{
		FollowerID: followerID,
		FollowedID: followedID,
	}
	database.DB.Create(&follow)

	return nil
}

func (s *FollowService) UnfollowUser(followerID, followedID uint) error {
	result := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).Delete(&models.UserFollow{})
	if result.RowsAffected == 0 {
		return utils.NewError("未关注该用户", 400)
	}

	return nil
}

func (s *FollowService) GetFollowers(userID uint, page, pageSize int) ([]models.User, int, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{}).
		Joins("JOIN user_follows ON user_follows.follower_id = users.id").
		Where("user_follows.followed_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("user_follows.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return users, totalPages, err
}

func (s *FollowService) GetFollowing(userID uint, page, pageSize int) ([]models.User, int, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{}).
		Joins("JOIN user_follows ON user_follows.followed_id = users.id").
		Where("user_follows.follower_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("user_follows.created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&users).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return users, totalPages, err
}

func (s *FollowService) CheckFollowing(followerID, followedID uint) (bool, error) {
	var follow models.UserFollow
	result := database.DB.Where("follower_id = ? AND followed_id = ?", followerID, followedID).First(&follow)
	if result.Error == nil {
		return true, nil
	}
	return false, nil
}
