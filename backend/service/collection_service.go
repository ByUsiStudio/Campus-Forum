package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
)

type CollectionService struct{}

var Collection = &CollectionService{}

func (s *CollectionService) GetCollections(userID uint) ([]models.Collection, error) {
	var collections []models.Collection
	err := database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&collections).Error
	return collections, err
}

func (s *CollectionService) GetCollection(collectionID uint) (*models.Collection, error) {
	var collection models.Collection
	err := database.DB.Preload("Articles").First(&collection, collectionID).Error
	if err != nil {
		return nil, utils.NewError("收藏夹不存在", 404)
	}
	return &collection, nil
}

func (s *CollectionService) CreateCollection(userID uint, name, description string) error {
	collection := models.Collection{
		UserID:      userID,
		Name:        name,
		Description: description,
	}

	if result := database.DB.Create(&collection); result.Error != nil {
		return utils.NewError("创建收藏夹失败", 500)
	}
	return nil
}

func (s *CollectionService) UpdateCollection(collectionID uint, name, description string) error {
	var collection models.Collection
	if result := database.DB.First(&collection, collectionID); result.Error != nil {
		return utils.NewError("收藏夹不存在", 404)
	}

	collection.Name = name
	collection.Description = description
	database.DB.Save(&collection)
	return nil
}

func (s *CollectionService) DeleteCollection(collectionID uint) error {
	var collection models.Collection
	if result := database.DB.First(&collection, collectionID); result.Error != nil {
		return utils.NewError("收藏夹不存在", 404)
	}

	database.DB.Delete(&collection)
	return nil
}

func (s *CollectionService) AddArticleToCollection(collectionID, articleID uint) error {
	var collectionArticle models.CollectionArticle
	result := database.DB.Where("collection_id = ? AND article_id = ?", collectionID, articleID).First(&collectionArticle)
	if result.Error == nil {
		return utils.NewError("文章已在收藏夹中", 400)
	}

	collectionArticle = models.CollectionArticle{
		CollectionID: collectionID,
		ArticleID:    articleID,
	}
	database.DB.Create(&collectionArticle)

	return nil
}

func (s *CollectionService) RemoveArticleFromCollection(collectionID, articleID uint) error {
	result := database.DB.Where("collection_id = ? AND article_id = ?", collectionID, articleID).Delete(&models.CollectionArticle{})
	if result.RowsAffected == 0 {
		return utils.NewError("文章不在收藏夹中", 400)
	}

	return nil
}

func (s *CollectionService) GetArticleVersions(articleID uint) ([]models.ArticleVersion, error) {
	var versions []models.ArticleVersion
	err := database.DB.Where("article_id = ?", articleID).Order("version DESC").Find(&versions).Error
	return versions, err
}

func (s *CollectionService) GetArticleVersion(articleID uint, version int) (*models.ArticleVersion, error) {
	var articleVersion models.ArticleVersion
	err := database.DB.Where("article_id = ? AND version = ?", articleID, version).First(&articleVersion).Error
	if err != nil {
		return nil, utils.NewError("版本不存在", 404)
	}
	return &articleVersion, nil
}

func (s *CollectionService) RestoreArticleVersion(articleID uint, version int) error {
	var articleVersion models.ArticleVersion
	result := database.DB.Where("article_id = ? AND version = ?", articleID, version).First(&articleVersion)
	if result.Error != nil {
		return utils.NewError("版本不存在", 404)
	}

	var article models.Article
	database.DB.First(&article, articleID)
	article.Content = articleVersion.Content
	article.ContentHTML = articleVersion.ContentHTML
	database.DB.Save(&article)

	return nil
}
