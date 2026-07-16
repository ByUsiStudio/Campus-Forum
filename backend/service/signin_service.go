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

	var user models.User
	database.DB.First(&user, userID)

	var lastRecord models.SignInRecord
	database.DB.Where("user_id = ? AND created_at < ?", userID, startOfDay).Order("created_at DESC").First(&lastRecord)

	continuousDay := 1
	if lastRecord.ID > 0 {
		lastYear, lastMonth, lastDay := lastRecord.CreatedAt.Date()
		yesterday := startOfDay.Add(-24 * time.Hour)
		yYear, yMonth, yDay := yesterday.Date()
		if lastYear == yYear && lastMonth == yMonth && lastDay == yDay {
			continuousDay = lastRecord.ContinuousDay + 1
		}
	}

	if continuousDay > user.SignInDays {
		user.SignInDays = continuousDay
	}

	if user.SignInDays > user.MaxContinuousDays {
		user.MaxContinuousDays = user.SignInDays
	}

	rewardPoints := 1
	if user.SignInDays >= 7 {
		rewardPoints = 3
	} else if user.SignInDays >= 3 {
		rewardPoints = 2
	}
	user.TotalCoins += rewardPoints
	user.TotalSignIns++

	database.DB.Save(&user)

	record = models.SignInRecord{
		UserID:        userID,
		SignInDate:    startOfDay.Format("2006-01-02"),
		ContinuousDay: continuousDay,
		RewardPoints:  rewardPoints,
	}
	database.DB.Create(&record)

	Article.updateUserExperience(userID, "signin", rewardPoints*10)

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
		lastYear, lastMonth, lastDay := lastRecord.CreatedAt.Date()
		yesterday := startOfDay.Add(-24 * time.Hour)
		yYear, yMonth, yDay := yesterday.Date()
		if lastYear == yYear && lastMonth == yMonth && lastDay == yDay {
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
			Enabled:      true,
			DailyPoints:  1,
			WeeklyBonus:  7,
			MonthlyBonus: 30,
			YearlyBonus:  100,
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
