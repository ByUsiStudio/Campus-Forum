package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type LeaderboardService struct{}

var Leaderboard = &LeaderboardService{}

func (s *LeaderboardService) GetLeaderboard(page, pageSize int) ([]models.Leaderboard, int, error) {
	var rankings []models.Leaderboard
	var total int64

	query := database.DB.Model(&models.Leaderboard{})
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("rank ASC").
		Offset(offset).Limit(pageSize).
		Find(&rankings).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return rankings, totalPages, err
}

func (s *LeaderboardService) GetWeeklyLeaderboard() ([]models.Leaderboard, error) {
	var rankings []models.Leaderboard
	err := database.DB.Preload("User").
		Where("period = ?", "weekly").
		Order("rank ASC").
		Limit(10).
		Find(&rankings).Error
	return rankings, err
}

func (s *LeaderboardService) GetMonthlyLeaderboard() ([]models.Leaderboard, error) {
	var rankings []models.Leaderboard
	err := database.DB.Preload("User").
		Where("period = ?", "monthly").
		Order("rank ASC").
		Limit(10).
		Find(&rankings).Error
	return rankings, err
}

func (s *LeaderboardService) GetUserRank(userID uint) (*models.Leaderboard, error) {
	var ranking models.Leaderboard
	err := database.DB.Where("user_id = ?", userID).Preload("User").First(&ranking).Error
	if err != nil {
		return nil, utils.NewError("用户排名不存在", 404)
	}
	return &ranking, nil
}

func (s *LeaderboardService) GetUserBadges(userID uint) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	err := database.DB.Where("user_id = ?", userID).Order("granted_at DESC").Find(&badges).Error
	return badges, err
}

func (s *LeaderboardService) UpdateBadgeDisplay(userID, badgeID uint, display bool) error {
	var badge models.UserBadge
	result := database.DB.Where("user_id = ? AND badge_id = ?", userID, badgeID).First(&badge)
	if result.Error != nil {
		return utils.NewError("徽章不存在", 404)
	}

	badge.IsDisplayed = display
	database.DB.Save(&badge)
	return nil
}

func (s *LeaderboardService) GrantBadge(userID uint, badgeType, badgeName, badgeIcon string) error {
	var userBadge models.UserBadge
	result := database.DB.Where("user_id = ? AND badge_type = ?", userID, badgeType).First(&userBadge)
	if result.Error == nil {
		return utils.NewError("用户已拥有该徽章", 400)
	}

	userBadge = models.UserBadge{
		UserID:      userID,
		BadgeType:   badgeType,
		BadgeName:   badgeName,
		BadgeIcon:   badgeIcon,
		IsDisplayed: true,
	}
	database.DB.Create(&userBadge)

	return nil
}

func (s *LeaderboardService) RevokeBadge(userID uint, badgeType string) error {
	result := database.DB.Where("user_id = ? AND badge_type = ?", userID, badgeType).Delete(&models.UserBadge{})
	if result.RowsAffected == 0 {
		return utils.NewError("用户未拥有该徽章", 400)
	}

	return nil
}
