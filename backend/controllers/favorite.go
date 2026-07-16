package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func AddFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, _ := strconv.Atoi(c.Param("article_id"))

	err := service.Article.AddFavorite(userID, uint(articleID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "收藏成功"})
}

func RemoveFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, _ := strconv.Atoi(c.Param("article_id"))

	err := service.Article.RemoveFavorite(userID, uint(articleID))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消收藏成功"})
}

func GetFavorites(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	articles, totalPages, err := service.Article.GetFavorites(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func CheckFavorite(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, _ := strconv.Atoi(c.Param("article_id"))

	isFavorite, err := service.Article.CheckFavorite(userID, uint(articleID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"is_favorite": isFavorite})
}
