package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCollections 获取用户收藏夹列表
func GetCollections(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var collections []models.Collection
	database.DB.Where("user_id = ?", userID).
		Order("created_at desc").
		Find(&collections)

	c.JSON(200, gin.H{
		"success": true,
		"data":    collections,
	})
}

// GetCollection 获取单个收藏夹详情
func GetCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	collectionID := c.Param("id")

	var collection models.Collection
	if err := database.DB.Where("id = ? AND user_id = ?", collectionID, userID).First(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 404, "收藏夹不存在")
		return
	}

	var articles []models.CollectionArticle
	database.DB.Where("collection_id = ?", collectionID).
		Preload("Article").
		Preload("Article.User").
		Order("created_at desc").
		Find(&articles)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"collection": collection,
			"articles":   articles,
		},
	})
}

// CreateCollection 创建收藏夹
func CreateCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var collection models.Collection
	if err := c.ShouldBindJSON(&collection); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	collection.UserID = userID

	if err := database.DB.Create(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 500, "创建失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    collection,
	})
}

// UpdateCollection 更新收藏夹
func UpdateCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	collectionID := c.Param("id")

	var collection models.Collection
	if err := database.DB.Where("id = ? AND user_id = ?", collectionID, userID).First(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 404, "收藏夹不存在")
		return
	}

	if err := c.ShouldBindJSON(&collection); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	database.DB.Save(&collection)

	c.JSON(200, gin.H{
		"success": true,
		"data":    collection,
	})
}

// DeleteCollection 删除收藏夹
func DeleteCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	collectionID := c.Param("id")

	var collection models.Collection
	if err := database.DB.Where("id = ? AND user_id = ?", collectionID, userID).First(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 404, "收藏夹不存在")
		return
	}

	database.DB.Where("collection_id = ?", collectionID).Delete(&models.CollectionArticle{})
	database.DB.Delete(&collection)

	c.JSON(200, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// AddArticleToCollection 添加文章到收藏夹
func AddArticleToCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	collectionID := c.Param("id")

	var req struct {
		ArticleID uint   `json:"article_id"`
		Note      string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	var collection models.Collection
	if err := database.DB.Where("id = ? AND user_id = ?", collectionID, userID).First(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 404, "收藏夹不存在")
		return
	}

	var article models.Article
	if err := database.DB.First(&article, req.ArticleID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "文章不存在")
		return
	}

	var existing models.CollectionArticle
	if err := database.DB.Where("collection_id = ? AND article_id = ?", collectionID, req.ArticleID).First(&existing).Error; err == nil {
		utils.SendErrorResponse(c, 400, "文章已在收藏夹中")
		return
	}

	collectionArticle := models.CollectionArticle{
		CollectionID: collection.ID,
		ArticleID:    req.ArticleID,
		Note:         req.Note,
	}
	if err := database.DB.Create(&collectionArticle).Error; err != nil {
		utils.SendErrorResponse(c, 500, "添加失败")
		return
	}

	database.DB.Model(&collection).UpdateColumn("article_count", gorm.Expr("article_count + 1"))

	c.JSON(200, gin.H{
		"success": true,
		"message": "添加成功",
	})
}

// RemoveArticleFromCollection 从收藏夹移除文章
func RemoveArticleFromCollection(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	collectionID := c.Param("id")
	articleID := c.Param("article_id")

	var collection models.Collection
	if err := database.DB.Where("id = ? AND user_id = ?", collectionID, userID).First(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 404, "收藏夹不存在")
		return
	}

	if err := database.DB.Where("collection_id = ? AND article_id = ?", collectionID, articleID).Delete(&models.CollectionArticle{}).Error; err != nil {
		utils.SendErrorResponse(c, 500, "移除失败")
		return
	}

	database.DB.Model(&collection).UpdateColumn("article_count", gorm.Expr("article_count - 1"))

	c.JSON(200, gin.H{
		"success": true,
		"message": "移除成功",
	})
}

// GetCollectionArticles 获取收藏夹文章列表
func GetCollectionArticles(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	collectionID := c.Param("id")

	var collection models.Collection
	if err := database.DB.Where("id = ? AND user_id = ?", collectionID, userID).First(&collection).Error; err != nil {
		utils.SendErrorResponse(c, 404, "收藏夹不存在")
		return
	}

	var articles []models.CollectionArticle
	database.DB.Where("collection_id = ?", collectionID).
		Preload("Article").
		Preload("Article.User").
		Order("created_at desc").
		Find(&articles)

	c.JSON(200, gin.H{
		"success": true,
		"data":    articles,
	})
}

// GetArticleVersions 获取文章版本历史
func GetArticleVersions(c *gin.Context) {
	articleID := c.Param("id")

	var versions []models.ArticleVersion
	database.DB.Where("article_id = ?", articleID).
		Order("version desc").
		Find(&versions)

	c.JSON(200, gin.H{
		"success": true,
		"data":    versions,
	})
}

// GetArticleVersion 获取指定版本内容
func GetArticleVersion(c *gin.Context) {
	articleID := c.Param("id")
	versionNum := c.Param("version")

	var version models.ArticleVersion
	if err := database.DB.Where("article_id = ? AND version = ?", articleID, versionNum).First(&version).Error; err != nil {
		utils.SendErrorResponse(c, 404, "版本不存在")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    version,
	})
}

// RestoreArticleVersion 恢复文章版本
func RestoreArticleVersion(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	articleID := c.Param("id")
	versionNum := c.Param("version")

	var article models.Article
	if err := database.DB.Where("id = ? AND user_id = ?", articleID, userID).First(&article).Error; err != nil {
		utils.SendErrorResponse(c, 404, "文章不存在或无权限")
		return
	}

	var version models.ArticleVersion
	if err := database.DB.Where("article_id = ? AND version = ?", articleID, versionNum).First(&version).Error; err != nil {
		utils.SendErrorResponse(c, 404, "版本不存在")
		return
	}

	article.Title = version.Title
	article.Content = version.Content
	article.ContentHTML = version.ContentHTML
	database.DB.Save(&article)

	c.JSON(200, gin.H{
		"success": true,
		"message": "恢复成功",
	})
}
