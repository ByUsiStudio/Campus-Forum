package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
	"time"
)

type AdminService struct{}

var Admin = &AdminService{}

func (s *AdminService) GetUsers(page, pageSize int, keyword string) ([]models.User, int, error) {
	var users []models.User
	var total int64

	query := database.DB.Model(&models.User{})
	if keyword != "" {
		query = query.Where("username LIKE ? OR display_name LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
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
		Find(&users).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return users, totalPages, err
}

func (s *AdminService) GetUser(userID uint) (*models.User, error) {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return nil, utils.NewError("用户不存在", 404)
	}
	return &user, nil
}

func (s *AdminService) UpdateUser(userID uint, displayName, avatar, signature string) error {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	if displayName != "" {
		user.DisplayName = displayName
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	if signature != "" {
		user.Signature = utils.SanitizeHTML(signature)
	}

	database.DB.Save(&user)
	return nil
}

func (s *AdminService) UpdateUserRole(userID uint, role string) error {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	user.Role = role
	database.DB.Save(&user)
	return nil
}

func (s *AdminService) BanUser(userID uint, reason string) error {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	user.Status = "banned"
	user.BanReason = reason
	database.DB.Save(&user)
	return nil
}

func (s *AdminService) UnbanUser(userID uint) error {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	user.Status = "normal"
	user.BanReason = ""
	database.DB.Save(&user)
	return nil
}

func (s *AdminService) DeleteUser(userID uint) error {
	var user models.User
	if result := database.DB.First(&user, userID); result.Error != nil {
		return utils.NewError("用户不存在", 404)
	}

	database.DB.Delete(&user)
	return nil
}

func (s *AdminService) GetAllArticles(page, pageSize int) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{})
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *AdminService) UpdateArticleStatus(articleID uint, status string) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	article.Status = status
	database.DB.Save(&article)
	return nil
}

func (s *AdminService) GetAllComments(page, pageSize int) ([]models.Comment, int, error) {
	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{})
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Article").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return comments, totalPages, err
}

func (s *AdminService) DeleteCommentAdmin(commentID uint) error {
	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		return utils.NewError("评论不存在", 404)
	}

	database.DB.Delete(&comment)
	return nil
}

func (s *AdminService) GetArticlesAdmin(page, pageSize int, status string) ([]models.Article, int, error) {
	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{})
	if status != "" {
		query = query.Where("status = ?", status)
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
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articles).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articles, totalPages, err
}

func (s *AdminService) DeleteArticleAdmin(articleID uint) error {
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		return utils.NewError("文章不存在", 404)
	}

	database.DB.Delete(&article)
	return nil
}

func (s *AdminService) GetCommentsAdmin(page, pageSize int) ([]models.Comment, int, error) {
	var comments []models.Comment
	var total int64

	query := database.DB.Model(&models.Comment{})
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").Preload("Article").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&comments).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return comments, totalPages, err
}

func (s *AdminService) GetLogs(page, pageSize int) ([]models.SystemLog, int, error) {
	var logs []models.SystemLog
	var total int64

	query := database.DB.Model(&models.SystemLog{})
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&logs).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return logs, totalPages, err
}

func (s *AdminService) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	var userCount int64
	database.DB.Model(&models.User{}).Count(&userCount)
	stats["total_users"] = userCount

	var articleCount int64
	database.DB.Model(&models.Article{}).Count(&articleCount)
	stats["total_articles"] = articleCount

	var commentCount int64
	database.DB.Model(&models.Comment{}).Count(&commentCount)
	stats["total_comments"] = commentCount

	var categoryCount int64
	database.DB.Model(&models.Category{}).Count(&categoryCount)
	stats["total_categories"] = categoryCount

	var today = time.Now().Format("2006-01-02")
	var todayUsers int64
	database.DB.Model(&models.User{}).Where("DATE(created_at) = ?", today).Count(&todayUsers)
	stats["today_new_users"] = todayUsers

	var todayArticles int64
	database.DB.Model(&models.Article{}).Where("DATE(created_at) = ?", today).Count(&todayArticles)
	stats["today_new_articles"] = todayArticles

	var onlineUsers int64
	database.DB.Model(&models.User{}).Where("online_status = ?", "online").Count(&onlineUsers)
	stats["online_users"] = onlineUsers

	return stats, nil
}
