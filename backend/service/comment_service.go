package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
	"time"

	"gorm.io/gorm"
)

type CommentService struct{}

var Comment = &CommentService{}

func (s *CommentService) CreateComment(userID, articleID uint, content string, parentID uint, isReply bool) (*models.Comment, error) {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return nil, utils.NewError("文章不存在", 404)
	}

	content = utils.SanitizeHTML(content)

	var parentIDPtr *uint
	if parentID > 0 {
		parentIDPtr = &parentID
	}

	comment := models.Comment{
		UserID:    userID,
		ArticleID: articleID,
		Content:   content,
		ParentID:  parentIDPtr,
		LikeCount: 0,
	}

	if result := database.DB.Create(&comment); result.Error != nil {
		return nil, utils.NewError("创建评论失败", 500)
	}

	database.DB.Model(&article).UpdateColumn("comment_count", gorm.Expr("comment_count + 1"))

	return &comment, nil
}

func (s *CommentService) GetComments(articleID uint, page, pageSize int) ([]models.Comment, int, error) {
	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{}).Where("article_id = ?", articleID)
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return comments, totalPages, err
}

func (s *CommentService) UpdateComment(userID, commentID uint, content string) error {
	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		return utils.NewError("评论不存在", 404)
	}

	if comment.UserID != userID {
		return utils.NewError("无权修改该评论", 403)
	}

	comment.Content = utils.SanitizeHTML(content)
	comment.UpdatedAt = time.Now()
	database.DB.Save(&comment)
	return nil
}

func (s *CommentService) DeleteComment(userID, commentID uint) error {
	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		return utils.NewError("评论不存在", 404)
	}

	if comment.UserID != userID {
		return utils.NewError("无权删除该评论", 403)
	}

	database.DB.Delete(&comment)

	database.DB.Model(&models.Article{}).Where("id = ?", comment.ArticleID).UpdateColumn("comment_count", gorm.Expr("comment_count - 1"))

	return nil
}

func (s *CommentService) LikeComment(userID, commentID uint) error {
	var like models.CommentLike
	result := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&like)
	if result.Error == nil {
		return utils.NewError("已点赞该评论", 400)
	}

	like = models.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}
	database.DB.Create(&like)

	database.DB.Model(&models.Comment{}).Where("id = ?", commentID).UpdateColumn("like_count", gorm.Expr("like_count + 1"))

	return nil
}

func (s *CommentService) UnlikeComment(userID, commentID uint) error {
	result := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{})
	if result.RowsAffected == 0 {
		return utils.NewError("未点赞该评论", 400)
	}

	database.DB.Model(&models.Comment{}).Where("id = ?", commentID).UpdateColumn("like_count", gorm.Expr("like_count - 1"))

	return nil
}

func (s *CommentService) GetCommentsByArticle(articleID uint, page, pageSize int) ([]models.Comment, int, error) {
	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{}).Where("article_id = ?", articleID)
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return comments, totalPages, err
}
