package repository

import (
	"forum/models"
	"time"
)

type SignInRecordRepository struct {
	*BaseRepository[models.SignInRecord]
}

func NewSignInRecordRepository() *SignInRecordRepository {
	return &SignInRecordRepository{
		BaseRepository: NewBaseRepository[models.SignInRecord](),
	}
}

func (r *SignInRecordRepository) GetTodaySignIn(userID uint) (*models.SignInRecord, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var record models.SignInRecord
	err := r.db.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfDay, endOfDay).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

func (r *SignInRecordRepository) GetSignInHistory(userID uint, page, pageSize int) ([]models.SignInRecord, int64, error) {
	var records []models.SignInRecord
	var total int64

	query := r.db.Model(&models.SignInRecord{}).Where("user_id = ?", userID)
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

func (r *SignInRecordRepository) GetSignInRankings(limit int) ([]models.SignInRecord, error) {
	var records []models.SignInRecord
	err := r.db.Select("user_id, COUNT(*) as count").
		Group("user_id").
		Order("count DESC").
		Limit(limit).
		Scan(&records).Error
	return records, err
}

type SignInConfigRepository struct {
	*BaseRepository[models.SignInConfig]
}

func NewSignInConfigRepository() *SignInConfigRepository {
	return &SignInConfigRepository{
		BaseRepository: NewBaseRepository[models.SignInConfig](),
	}
}

func (r *SignInConfigRepository) GetActiveConfig() (*models.SignInConfig, error) {
	var config models.SignInConfig
	err := r.db.Order("created_at DESC").First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
