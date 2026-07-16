package repository

import (
	"forum/models"

	"gorm.io/gorm"
)

type CommentRepository struct {
	*BaseRepository[models.Comment]
}

func NewCommentRepository() *CommentRepository {
	return &CommentRepository{
		BaseRepository: NewBaseRepository[models.Comment](),
	}
}

func (r *CommentRepository) GetCommentsByArticle(articleID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Where("article_id = ?", articleID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("created_at ASC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepository) GetCommentWithUser(id uint) (*models.Comment, error) {
	var comment models.Comment
	err := r.db.Preload("User").First(&comment, id).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *CommentRepository) GetCommentsByUser(userID uint, page, pageSize int) ([]models.Comment, int64, error) {
	var comments []models.Comment
	var total int64

	query := r.db.Model(&models.Comment{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Article").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	return comments, total, err
}

func (r *CommentRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + 1")).Error
}

func (r *CommentRepository) DecrementLikeCount(id uint) error {
	return r.db.Model(&models.Comment{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count - 1")).Error
}

func (r *CommentRepository) CountByArticle(articleID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Comment{}).Where("article_id = ?", articleID).Count(&count).Error
	return count, err
}

type CommentLikeRepository struct {
	*BaseRepository[models.CommentLike]
}

func NewCommentLikeRepository() *CommentLikeRepository {
	return &CommentLikeRepository{
		BaseRepository: NewBaseRepository[models.CommentLike](),
	}
}

func (r *CommentLikeRepository) CheckLike(userID, commentID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.CommentLike{}).Where("user_id = ? AND comment_id = ?", userID, commentID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *CommentLikeRepository) GetByUserAndComment(userID, commentID uint) (*models.CommentLike, error) {
	var like models.CommentLike
	err := r.db.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like).Error
	if err != nil {
		return nil, err
	}
	return &like, nil
}
