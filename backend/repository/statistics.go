package repository

import (
	"forum/models"
)

type UserStatisticsRepository struct {
	*BaseRepository[models.UserStatistics]
}

func NewUserStatisticsRepository() *UserStatisticsRepository {
	return &UserStatisticsRepository{
		BaseRepository: NewBaseRepository[models.UserStatistics](),
	}
}

func (r *UserStatisticsRepository) GetUserStatistics(userID uint) (*models.UserStatistics, error) {
	var stats models.UserStatistics
	err := r.db.Where("user_id = ?", userID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, err
}

type DailyStatisticsRepository struct {
	*BaseRepository[models.DailyStatistics]
}

func NewDailyStatisticsRepository() *DailyStatisticsRepository {
	return &DailyStatisticsRepository{
		BaseRepository: NewBaseRepository[models.DailyStatistics](),
	}
}

type ArticleStatisticsRepository struct {
	*BaseRepository[models.ArticleStatistics]
}

func NewArticleStatisticsRepository() *ArticleStatisticsRepository {
	return &ArticleStatisticsRepository{
		BaseRepository: NewBaseRepository[models.ArticleStatistics](),
	}
}

func (r *ArticleStatisticsRepository) GetArticleStatistics(articleID uint) (*models.ArticleStatistics, error) {
	var stats models.ArticleStatistics
	err := r.db.Where("article_id = ?", articleID).First(&stats).Error
	if err != nil {
		return nil, err
	}
	return &stats, err
}

type SystemOverviewRepository struct {
	*BaseRepository[models.SystemOverview]
}

func NewSystemOverviewRepository() *SystemOverviewRepository {
	return &SystemOverviewRepository{
		BaseRepository: NewBaseRepository[models.SystemOverview](),
	}
}

func (r *SystemOverviewRepository) GetLatestOverview() (*models.SystemOverview, error) {
	var overview models.SystemOverview
	err := r.db.Order("created_at DESC").First(&overview).Error
	if err != nil {
		return nil, err
	}
	return &overview, nil
}
