package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
)

type FollowService struct{}

var Follow = &FollowService{}

func (s *FollowService) SendFriendRequest(fromUserID, toUserID uint) error {
	if fromUserID == toUserID {
		return utils.NewError("不能添加自己为好友", 400)
	}

	var request models.FriendRequest
	result := database.DB.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).First(&request)
	if result.Error == nil {
		return utils.NewError("好友请求已发送", 400)
	}

	request = models.FriendRequest{
		FromUserID: fromUserID,
		ToUserID:   toUserID,
		Status:     "pending",
	}
	database.DB.Create(&request)

	return nil
}

func (s *FollowService) AcceptFriendRequest(userID, requestID uint) error {
	var request models.FriendRequest
	if result := database.DB.First(&request, requestID); result.Error != nil {
		return utils.NewError("请求不存在", 404)
	}

	if request.ToUserID != userID {
		return utils.NewError("无权处理该请求", 403)
	}

	request.Status = "accepted"
	database.DB.Save(&request)

	friend1 := models.Friend{
		UserID:   request.FromUserID,
		FriendID: request.ToUserID,
	}
	friend2 := models.Friend{
		UserID:   request.ToUserID,
		FriendID: request.FromUserID,
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

	if request.ToUserID != userID {
		return utils.NewError("无权处理该请求", 403)
	}

	request.Status = "rejected"
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
	err := database.DB.Where("user_id = ?", userID).Preload("FriendUser").Find(&friends).Error
	return friends, err
}

func (s *FollowService) GetFriendRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := database.DB.Where("to_user_id = ? AND status = ?", userID, "pending").Preload("FromUser").Find(&requests).Error
	return requests, err
}

func (s *FollowService) GetSentFriendRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := database.DB.Where("from_user_id = ? AND status = ?", userID, "pending").Preload("ToUser").Find(&requests).Error
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
	result = database.DB.Where("from_user_id = ? AND to_user_id = ?", userID1, userID2).First(&request)
	if result.Error == nil {
		return "pending_sent", nil
	}

	result = database.DB.Where("from_user_id = ? AND to_user_id = ?", userID2, userID1).First(&request)
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
