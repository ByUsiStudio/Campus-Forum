package repository

import (
	"forum/models"
)

type LikeRepository struct {
	*BaseRepository[models.Like]
}

func NewLikeRepository() *LikeRepository {
	return &LikeRepository{
		BaseRepository: NewBaseRepository[models.Like](),
	}
}

func (r *LikeRepository) CheckLike(userID, articleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *LikeRepository) GetByUserAndArticle(userID, articleID uint) (*models.Like, error) {
	var like models.Like
	err := r.db.Where("user_id = ? AND article_id = ?", userID, articleID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}

func (r *LikeRepository) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Like{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}