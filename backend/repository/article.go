package repository

import (
	"forum/database"
	"forum/models"

	"gorm.io/gorm"
)

// ArticleRepository 文章数据访问层
type ArticleRepository struct {
	*BaseRepository[models.Article]
}

// NewArticleRepository 创建文章 Repository
func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{
		BaseRepository: NewBaseRepository[models.Article](),
	}
}

// GetArticleWithUser 获取文章及其作者信息
func (r *ArticleRepository) GetArticleWithUser(id uint) (*models.Article, error) {
	var article models.Article
	err := r.db.Preload("User").Preload("Category").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

// GetArticlesByUser 获取用户的文章
func (r *ArticleRepository) GetArticlesByUser(userID uint, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	return articles, total, err
}

// GetPublishedArticles 获取已发布的文章
func (r *ArticleRepository) GetPublishedArticles(page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("status = ?", "published")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Category").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	return articles, total, err
}

// GetDrafts 获取用户草稿
func (r *ArticleRepository) GetDrafts(userID uint, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("user_id = ? AND status = ?", userID, "draft")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").
		Order("updated_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	return articles, total, err
}

// SearchArticles 搜索文章
func (r *ArticleRepository) SearchArticles(keyword string, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("status = ?", "published")
	if keyword != "" {
		query = query.Where("(title LIKE ? OR content LIKE ?)", "%"+keyword+"%", "%"+keyword+"%")
	}

	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Category").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	return articles, total, err
}

// GetPinnedArticles 获取置顶文章
func (r *ArticleRepository) GetPinnedArticles() ([]models.Article, error) {
	var articles []models.Article
	err := r.db.Where("is_pinned = ? AND status = ?", true, "published").
		Preload("User").Preload("Category").
		Order("pinned_at DESC").
		Find(&articles).Error
	return articles, err
}

// SetPin 置顶/取消置顶文章
func (r *ArticleRepository) SetPin(id uint, isPinned bool) error {
	updates := map[string]interface{}{
		"is_pinned": isPinned,
	}
	if isPinned {
		updates["pinned_at"] = database.DB.Raw("NOW()")
	} else {
		updates["pinned_at"] = nil
	}
	return r.db.Model(&models.Article{}).Where("id = ?", id).Updates(updates).Error
}

// GetByCategory 按分类获取文章
func (r *ArticleRepository) GetByCategory(categoryID uint, page, pageSize int) ([]models.Article, int64, error) {
	var articles []models.Article
	var total int64

	query := r.db.Model(&models.Article{}).Where("category_id = ? AND status = ?", categoryID, "published")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Category").
		Order("is_pinned DESC, created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	return articles, total, err
}

// IncrementViewCount 增加浏览量
func (r *ArticleRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&models.Article{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1")).Error
}
