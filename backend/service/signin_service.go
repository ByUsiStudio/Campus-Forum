package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
	"time"
)

type SignInService struct{}

var SignIn = &SignInService{}

func (s *SignInService) SignIn(userID uint) error {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var record models.SignInRecord
	result := database.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfDay, endOfDay).First(&record)
	if result.Error == nil {
		return utils.NewError("今日已签到", 400)
	}

	record = models.SignInRecord{
		UserID: userID,
	}
	database.DB.Create(&record)

	var user models.User
	database.DB.First(&user, userID)
	user.TotalSignIns++

	var lastRecord models.SignInRecord
	database.DB.Where("user_id = ? AND created_at < ?", userID, startOfDay).Order("created_at DESC").First(&lastRecord)
	if lastRecord.ID > 0 {
		lastSignInDate := lastRecord.CreatedAt.Date()
		yesterday := startOfDay.Add(-24 * time.Hour).Date()
		if lastSignInDate == yesterday {
			user.SignInDays++
		} else {
			user.SignInDays = 1
		}
	} else {
		user.SignInDays = 1
	}

	if user.SignInDays > user.MaxContinuousDays {
		user.MaxContinuousDays = user.SignInDays
	}

	coins := 1
	if user.SignInDays >= 7 {
		coins = 3
	} else if user.SignInDays >= 3 {
		coins = 2
	}
	user.TotalCoins += coins

	database.DB.Save(&user)

	s.updateUserExperience(userID, "signin", coins*10)

	return nil
}

func (s *SignInService) GetSignInStatus(userID uint) (*models.SignInRecord, int, error) {
	now := time.Now()
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	endOfDay := startOfDay.Add(24 * time.Hour)

	var record models.SignInRecord
	result := database.DB.Where("user_id = ? AND created_at >= ? AND created_at < ?", userID, startOfDay, endOfDay).First(&record)
	if result.Error == nil {
		return &record, 0, nil
	}

	var user models.User
	database.DB.First(&user, userID)

	var lastRecord models.SignInRecord
	database.DB.Where("user_id = ? AND created_at < ?", userID, startOfDay).Order("created_at DESC").First(&lastRecord)
	if lastRecord.ID > 0 {
		lastSignInDate := lastRecord.CreatedAt.Date()
		yesterday := startOfDay.Add(-24 * time.Hour).Date()
		if lastSignInDate == yesterday {
			return nil, user.SignInDays + 1, nil
		}
	}

	return nil, 1, nil
}

func (s *SignInService) GetSignInHistory(userID uint, page, pageSize int) ([]models.SignInRecord, int, error) {
	var records []models.SignInRecord
	var total int64

	query := database.DB.Model(&models.SignInRecord{}).Where("user_id = ?", userID)
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return records, totalPages, err
}

func (s *SignInService) GetSignInRankings(limit int) ([]models.SignInRecord, error) {
	var records []models.SignInRecord
	err := database.DB.Select("user_id, COUNT(*) as count").
		Group("user_id").
		Order("count DESC").
		Limit(limit).
		Scan(&records).Error
	return records, err
}

func (s *SignInService) GetSignInConfig() (*models.SignInConfig, error) {
	var config models.SignInConfig
	err := database.DB.Order("created_at DESC").First(&config).Error
	if err != nil {
		config = models.SignInConfig{
			Enabled:              true,
			DailyCoins:           1,
			Continuous3DaysBonus: 2,
			Continuous7DaysBonus: 3,
		}
		database.DB.Create(&config)
		return &config, nil
	}
	return &config, nil
}

func (s *SignInService) UpdateSignInConfig(config models.SignInConfig) error {
	database.DB.Create(&config)
	return nil
}
