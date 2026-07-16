package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type FavoriteService struct{}

var Favorite = &FavoriteService{}

func (s *FavoriteService) AddFavorite(userID, articleID uint) error {
	var favorite models.Favorite
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite)
	if result.Error == nil {
		return utils.NewError("已收藏该文章", 400)
	}

	favorite = models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}
	database.DB.Create(&favorite)

	database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("favorite_count", utils.Increment("favorite_count"))

	return nil
}

func (s *FavoriteService) RemoveFavorite(userID, articleID uint) error {
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Favorite{})
	if result.RowsAffected == 0 {
		return utils.NewError("未收藏该文章", 400)
	}

	database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("favorite_count", utils.Decrement("favorite_count"))

	return nil
}

func (s *FavoriteService) GetFavorites(userID uint, page, pageSize int) ([]models.Favorite, int, error) {
	var favorites []models.Favorite
	var total int64

	query := database.DB.Model(&models.Favorite{}).Where("user_id = ?", userID)
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

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return favorites, totalPages, err
}

func (s *FavoriteService) CheckFavorite(userID, articleID uint) (bool, error) {
	var count int64
	err := database.DB.Model(&models.Favorite{}).Where("user_id = ? AND article_id = ?", userID, articleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
