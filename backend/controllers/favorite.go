package controllers

import (
	"forum/database"
	"forum/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// AddFavorite 添加收藏
func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleIDStr := c.Param("id")
	articleID := uint(parseUint(articleIDStr))

	// 检查文章是否存在
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查是否已收藏
	var existingFavorite models.Favorite
	if result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&existingFavorite); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经收藏过了"})
		return
	}

	// 创建收藏
	favorite := models.Favorite{
		UserID:    userID,
		ArticleID: articleID,
	}

	database.DB.Create(&favorite)
	database.DB.Model(&article).UpdateColumn("favorite_count", gorm.Expr("favorite_count + 1"))

	c.JSON(http.StatusOK, gin.H{
		"message":        "收藏成功",
		"favorite_count": article.FavoriteCount + 1,
	})
}

// RemoveFavorite 取消收藏
func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleIDStr := c.Param("id")
	articleID := uint(parseUint(articleIDStr))

	// 检查是否存在收藏
	var favorite models.Favorite
	if result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未收藏该文章"})
		return
	}

	// 删除收藏
	database.DB.Delete(&favorite)

	// 更新文章收藏数
	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error == nil {
		if article.FavoriteCount > 0 {
			database.DB.Model(&article).UpdateColumn("favorite_count", gorm.Expr("favorite_count - 1"))
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消收藏成功"})
}

// GetFavorites 获取用户收藏列表
func GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")

	var favorites []models.Favorite
	database.DB.Where("user_id = ?", userID).Order("created_at DESC").Find(&favorites)

	// 获取收藏的文章列表
	var articles []models.Article
	for _, favorite := range favorites {
		var article models.Article
		database.DB.Preload("User").First(&article, favorite.ArticleID)
		articles = append(articles, article)
	}

	c.JSON(http.StatusOK, gin.H{
		"articles": articles,
		"count":    len(articles),
	})
}

// CheckFavorite 检查文章是否已收藏
func CheckFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleIDStr := c.Param("id")
	articleID := uint(parseUint(articleIDStr))

	var favorite models.Favorite
	result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&favorite)

	c.JSON(http.StatusOK, gin.H{
		"favorited": result.Error == nil,
	})
}
