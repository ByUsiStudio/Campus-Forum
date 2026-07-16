package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
	"time"

	"gorm.io/gorm"
)

type ArticleService struct{}

var Article = &ArticleService{}

func (s *ArticleService) CreateArticle(userID uint, title, content string, categoryID uint, isAnonymous bool) (*models.Article, error) {
	contentHTML := utils.RenderMarkdown(content)

	article := models.Article{
		Title:         title,
		Content:       content,
		ContentHTML:   contentHTML,
		UserID:        userID,
		CategoryID:    categoryID,
		Status:        "published",
		IsAnonymous:   isAnonymous,
		IsPinned:      false,
		LikeCount:     0,
		CommentCount:  0,
		FavoriteCount: 0,
		ViewCount:     0,
		ShareCount:    0,
		CoinCount:     0,
	}

	if result := database.DB.Create(&article); result.Error != nil {
		return nil, utils.NewError("创建文章失败", 500)
	}

	s.updateUserExperience(userID, "create_article", 10)
	s.updateUserStatistics(userID, "article_count", 1)

	return &article, nil
}

func (s *ArticleService) GetArticles(page, pageSize int, categoryID uint) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).Where("status = ?", "published")
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *ArticleService) GetArticle(id uint) (*models.Article, error) {
	var article models.Article
	err := database.DB.Preload("User").Preload("Category").First(&article, id).Error
	if err != nil {
		return nil, utils.NewError("文章不存在", 404)
	}

	database.DB.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	return &article, nil
}

func (s *ArticleService) UpdateArticle(userID, articleID uint, title, content string, categoryID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	if article.UserID != userID {
		return utils.NewError("无权修改该文章", 403)
	}

	article.Title = title
	article.Content = content
	article.ContentHTML = utils.RenderMarkdown(content)
	article.CategoryID = categoryID
	article.UpdatedAt = time.Now()

	database.DB.Save(&article)
	return nil
}

func (s *ArticleService) DeleteArticle(userID, articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	if article.UserID != userID {
		return utils.NewError("无权删除该文章", 403)
	}

	article.Status = "deleted"
	database.DB.Save(&article)
	return nil
}

func (s *ArticleService) RestoreArticle(articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	article.Status = "published"
	database.DB.Save(&article)
	return nil
}

func (s *ArticleService) LikeArticle(userID, articleID uint) error {
	var like models.Like
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&like)
	if result.Error == nil {
		return utils.NewError("已点赞该文章", 400)
	}

	like = models.Like{
		UserID:    userID,
		ArticleID: articleID,
	}
	database.DB.Create(&like)

	database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("like_count", gorm.Expr("like_count + 1"))

	s.updateUserExperience(userID, "like_article", 1)

	return nil
}

func (s *ArticleService) UnlikeArticle(userID, articleID uint) error {
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Like{})
	if result.RowsAffected == 0 {
		return utils.NewError("未点赞该文章", 400)
	}

	database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("like_count", gorm.Expr("like_count - 1"))

	return nil
}

func (s *ArticleService) GetMyArticles(userID uint, page, pageSize int) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).Where("user_id = ? AND status != ?", userID, "deleted")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *ArticleService) SearchArticles(keyword string, page, pageSize int) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).Where("status = ?", "published")
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *ArticleService) ShareArticle(articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	database.DB.Model(&article).UpdateColumn("share_count", gorm.Expr("share_count + 1"))
	return nil
}

func (s *ArticleService) PinArticle(articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	article.IsPinned = true
	now := time.Now()
	article.PinnedAt = &now
	database.DB.Save(&article)
	return nil
}

func (s *ArticleService) UnpinArticle(articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	article.IsPinned = false
	article.PinnedAt = nil
	database.DB.Save(&article)
	return nil
}

func (s *ArticleService) GetDraftArticles(userID uint, page, pageSize int) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).Where("user_id = ? AND status = ?", userID, "draft")
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *ArticleService) PublishDraft(userID, articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	if article.UserID != userID {
		return utils.NewError("无权发布该文章", 403)
	}

	article.Status = "published"
	database.DB.Save(&article)

	s.updateUserExperience(userID, "publish_article", 5)

	return nil
}

func (s *ArticleService) CoinArticle(userID, articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	if user.TotalCoins < 1 {
		return utils.NewError("币不足", 400)
	}

	database.DB.Model(&user).UpdateColumn("total_coins", gorm.Expr("total_coins - 1"))
	database.DB.Model(&article).UpdateColumn("coin_count", gorm.Expr("coin_count + 1"))

	coinRecord := models.CoinRecord{
		UserID:    userID,
		ArticleID: articleID,
		Amount:    1,
		Type:      "coin",
	}
	database.DB.Create(&coinRecord)

	s.updateUserExperience(userID, "coin_article", 5)

	return nil
}

func (s *ArticleService) updateUserExperience(userID uint, action string, amount int) {
	var userLevel models.UserLevel
	result := database.DB.Where("user_id = ?", userID).First(&userLevel)
	if result.Error != nil {
		userLevel = models.UserLevel{
			UserID:     userID,
			Level:      1,
			Experience: 0,
		}
		database.DB.Create(&userLevel)
	}

	userLevel.Experience += amount
	for {
		var levelConfig models.LevelConfig
		result := database.DB.Where("level = ?", userLevel.Level+1).First(&levelConfig)
		if result.Error != nil || userLevel.Experience < levelConfig.RequiredExp {
			break
		}
		userLevel.Level++
	}

	database.DB.Save(&userLevel)

	expRecord := models.ExperienceRecord{
		UserID: userID,
		Action: action,
		Amount: amount,
		Level:  userLevel.Level,
	}
	database.DB.Create(&expRecord)
}

func (s *ArticleService) updateUserStatistics(userID uint, field string, increment int) {
	var stats models.UserStatistics
	result := database.DB.Where("user_id = ?", userID).First(&stats)
	if result.Error != nil {
		stats = models.UserStatistics{UserID: userID}
		database.DB.Create(&stats)
	}

	switch field {
	case "article_count":
		stats.ArticleCount += increment
	case "comment_count":
		stats.CommentCount += increment
	case "like_count":
		stats.LikeCount += increment
	case "favorite_count":
		stats.FavoriteCount += increment
	}

	database.DB.Save(&stats)
}
