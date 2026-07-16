package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"

	"gorm.io/gorm"
)

type UserService struct{}

var User = &UserService{}

func (s *UserService) GetUserByID(userID uint) (*models.User, error) {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return nil, utils.NewError("用户不存在", 404)
	}
	return &user, nil
}

func (s *UserService) GetUserArticles(userID uint, page, pageSize int) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).Where("user_id = ? AND status = ?", userID, "published")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *UserService) UpdateUserStatus(userID uint, status, description string) error {
	var userStatus models.UserActivity
	result := database.DB.Where("user_id = ?", userID).First(&userStatus)
	if result.Error != nil {
		userStatus = models.UserActivity{UserID: userID}
	}

	userStatus.Status = status
	userStatus.Description = description
	database.DB.Save(&userStatus)

	return nil
}

func (s *UserService) GetUserStatus(userID uint) (*models.UserActivity, error) {
	var userStatus models.UserActivity
	result := database.DB.Where("user_id = ?", userID).First(&userStatus)
	if result.Error != nil {
		return nil, utils.NewError("用户状态不存在", 404)
	}
	return &userStatus, nil
}

func (s *UserService) UpdateOnlineStatus(userID uint, status string) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("online_status", status).Error
}

func (s *UserService) UpdateLastActive(userID uint) error {
	return database.DB.Model(&models.User{}).Where("id = ?", userID).Update("last_active_at", gorm.Expr("NOW()")).Error
}

func (s *UserService) GetAllUserStatuses(page, pageSize int) ([]models.UserActivity, int, error) {
	var statuses []models.UserActivity
	var total int64

	query := database.DB.Model(&models.UserActivity{})
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&statuses).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return statuses, totalPages, err
}

func (s *UserService) GetOnlineUsers() ([]models.User, error) {
	var users []models.User
	err := database.DB.Where("online_status = ?", "online").Find(&users).Error
	return users, err
}

func (s *UserService) CleanupUserStatuses() error {
	return database.DB.Where("updated_at < NOW() - INTERVAL 1 HOUR").Delete(&models.UserActivity{}).Error
}
