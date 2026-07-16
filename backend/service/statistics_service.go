package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
)

type StatisticsService struct{}

var Statistics = &StatisticsService{}

func (s *StatisticsService) GetUserStatistics(userID uint) (*models.UserStatistics, error) {
	var stats models.UserStatistics
	result := database.DB.Where("user_id = ?", userID).First(&stats)
	if result.Error != nil {
		return nil, utils.NewError("用户统计不存在", 404)
	}
	return &stats, nil
}

func (s *StatisticsService) GetDailyStatistics() ([]models.DailyStatistics, error) {
	var stats []models.DailyStatistics
	err := database.DB.Order("date DESC").Limit(30).Find(&stats).Error
	return stats, err
}

func (s *StatisticsService) GetSystemOverview() (*models.SystemOverview, error) {
	var overview models.SystemOverview
	result := database.DB.Order("created_at DESC").First(&overview)
	if result.Error != nil {
		return nil, utils.NewError("系统概览不存在", 404)
	}
	return &overview, nil
}

func (s *StatisticsService) GetUserActivity(userID uint) ([]models.UserActivity, error) {
	var activities []models.UserActivity
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Limit(20).Find(&activities).Error
	return activities, err
}

func (s *StatisticsService) GetStatisticsDashboard() (*models.SystemOverview, error) {
	var overview models.SystemOverview
	result := database.DB.Order("created_at DESC").First(&overview)
	if result.Error != nil {
		return nil, utils.NewError("系统概览不存在", 404)
	}
	return &overview, nil
}

func (s *StatisticsService) GetArticleStatistics(articleID uint) (*models.ArticleStatistics, error) {
	var stats models.ArticleStatistics
	result := database.DB.Where("article_id = ?", articleID).First(&stats)
	if result.Error != nil {
		return nil, utils.NewError("文章统计不存在", 404)
	}
	return &stats, nil
}
