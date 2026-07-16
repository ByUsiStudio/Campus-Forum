package repository

import (
	"forum/models"
)

type LeaderboardRepository struct {
	*BaseRepository[models.Leaderboard]
}

func NewLeaderboardRepository() *LeaderboardRepository {
	return &LeaderboardRepository{
		BaseRepository: NewBaseRepository[models.Leaderboard](),
	}
}

func (r *LeaderboardRepository) GetRankings(page, pageSize int) ([]models.Leaderboard, int64, error) {
	var rankings []models.Leaderboard
	var total int64

	query := r.db.Model(&models.Leaderboard{})
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

	return rankings, total, err
}

func (r *LeaderboardRepository) GetUserRank(userID uint) (*models.Leaderboard, error) {
	var ranking models.Leaderboard
	err := r.db.Where("user_id = ?", userID).Preload("User").First(&ranking).Error
	if err != nil {
		return nil, err
	}
	return &ranking, nil
}

type UserBadgeRepository struct {
	*BaseRepository[models.UserBadge]
}

func NewUserBadgeRepository() *UserBadgeRepository {
	return &UserBadgeRepository{
		BaseRepository: NewBaseRepository[models.UserBadge](),
	}
}

func (r *UserBadgeRepository) GetUserBadges(userID uint) ([]models.UserBadge, error) {
	var badges []models.UserBadge
	err := r.db.Where("user_id = ?", userID).Order("granted_at DESC").Find(&badges).Error
	return badges, err
}

func (r *UserBadgeRepository) CheckUserBadge(userID, badgeID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserBadge{}).Where("user_id = ? AND badge_id = ?", userID, badgeID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}