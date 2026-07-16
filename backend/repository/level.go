package repository

import (
	"forum/models"
)

type UserLevelRepository struct {
	*BaseRepository[models.UserLevel]
}

func NewUserLevelRepository() *UserLevelRepository {
	return &UserLevelRepository{
		BaseRepository: NewBaseRepository[models.UserLevel](),
	}
}

func (r *UserLevelRepository) GetUserLevel(userID uint) (*models.UserLevel, error) {
	var userLevel models.UserLevel
	err := r.db.Where("user_id = ?", userID).First(&userLevel).Error
	if err != nil {
		return nil, err
	}
	return &userLevel, nil
}

type ExperienceRecordRepository struct {
	*BaseRepository[models.ExperienceRecord]
}

func NewExperienceRecordRepository() *ExperienceRecordRepository {
	return &ExperienceRecordRepository{
		BaseRepository: NewBaseRepository[models.ExperienceRecord](),
	}
}

func (r *ExperienceRecordRepository) GetUserExperienceRecords(userID uint, page, pageSize int) ([]models.ExperienceRecord, int64, error) {
	var records []models.ExperienceRecord
	var total int64

	query := r.db.Model(&models.ExperienceRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error

	return records, total, err
}

type LevelConfigRepository struct {
	*BaseRepository[models.LevelConfig]
}

func NewLevelConfigRepository() *LevelConfigRepository {
	return &LevelConfigRepository{
		BaseRepository: NewBaseRepository[models.LevelConfig](),
	}
}

func (r *LevelConfigRepository) GetLevelConfigsOrdered() ([]models.LevelConfig, error) {
	var configs []models.LevelConfig
	err := r.db.Order("level ASC").Find(&configs).Error
	return configs, err
}

type AchievementRepository struct {
	*BaseRepository[models.Achievement]
}

func NewAchievementRepository() *AchievementRepository {
	return &AchievementRepository{
		BaseRepository: NewBaseRepository[models.Achievement](),
	}
}

type UserAchievementRepository struct {
	*BaseRepository[models.UserAchievement]
}

func NewUserAchievementRepository() *UserAchievementRepository {
	return &UserAchievementRepository{
		BaseRepository: NewBaseRepository[models.UserAchievement](),
	}
}

func (r *UserAchievementRepository) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	var achievements []models.UserAchievement
	err := r.db.Where("user_id = ?", userID).Preload("Achievement").Order("achieved_at DESC").Find(&achievements).Error
	return achievements, err
}