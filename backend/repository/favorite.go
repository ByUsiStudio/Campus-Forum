package repository

import (
	"forum/models"
)

type FavoriteRepository struct {
	*BaseRepository[models.Favorite]
}

func NewFavoriteRepository() *FavoriteRepository {
	return &FavoriteRepository{
		BaseRepository: NewBaseRepository[models.Favorite](),
	}
}

func (r *FavoriteRepository) CheckFavorite(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *FavoriteRepository) GetByUserAndArticle(userID, articleID uint) (*models.Favorite, error) {
	var favorite models.Favorite
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite).Error
	if err != nil {
		return nil, err
	}
	return &favorite, nil
}

func (r *FavoriteRepository) GetFavoritesByUser(userID uint, page, pageSize int) ([]models.Favorite, int64, error) {
	var favorites []models.Favorite
	var total int64

	query := r.db.Model(&models.Favorite{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Article").Preload("Article.User").Preload("Article.Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&favorites).Error

	return favorites, total, err
}

func (r *FavoriteRepository) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Favorite{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}
