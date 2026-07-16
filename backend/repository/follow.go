package repository

import (
	"forum/models"
)

type FriendRepository struct {
	*BaseRepository[models.Friend]
}

func NewFriendRepository() *FriendRepository {
	return &FriendRepository{
		BaseRepository: NewBaseRepository[models.Friend](),
	}
}

func (r *FriendRepository) GetFriendsByUser(userID uint) ([]models.Friend, error) {
	var friends []models.Friend
	err := r.db.Where("user_id = ?", userID).Preload("FriendUser").Find(&friends).Error
	return friends, err
}

func (r *FriendRepository) CheckFriendship(userID1, userID2 uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Friend{}).
		Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
			userID1, userID2, userID2, userID1).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FriendRepository) GetFriendship(userID1, userID2 uint) (*models.Friend, error) {
	var friend models.Friend
	err := r.db.Where("(user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)",
		userID1, userID2, userID2, userID1).First(&friend).Error
	if err != nil {
		return nil, err
	}
	return &friend, nil
}

type FriendRequestRepository struct {
	*BaseRepository[models.FriendRequest]
}

func NewFriendRequestRepository() *FriendRequestRepository {
	return &FriendRequestRepository{
		BaseRepository: NewBaseRepository[models.FriendRequest](),
	}
}

func (r *FriendRequestRepository) GetReceivedRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := r.db.Where("to_user_id = ? AND status = ?", userID, "pending").Preload("FromUser").Find(&requests).Error
	return requests, err
}

func (r *FriendRequestRepository) GetSentRequests(userID uint) ([]models.FriendRequest, error) {
	var requests []models.FriendRequest
	err := r.db.Where("from_user_id = ? AND status = ?", userID, "pending").Preload("ToUser").Find(&requests).Error
	return requests, err
}

func (r *FriendRequestRepository) GetRequest(fromUserID, toUserID uint) (*models.FriendRequest, error) {
	var request models.FriendRequest
	err := r.db.Where("from_user_id = ? AND to_user_id = ?", fromUserID, toUserID).First(&request).Error
	if err != nil {
		return nil, err
	}
	return &request, nil
}
