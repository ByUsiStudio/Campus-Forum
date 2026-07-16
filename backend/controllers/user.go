package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserProfile(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := service.Auth.GetProfile(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":           user.ID,
		"username":     user.Username,
		"display_name": user.DisplayName,
		"avatar":       user.Avatar,
		"signature":    user.Signature,
		"role":         user.Role,
		"status":       user.Status,
		"created_at":   user.CreatedAt,
		"level":        user.Level,
		"exp":          user.Exp,
		"coins":        user.Coins,
	})
}

func GetUserArticles(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	articles, totalPages, err := service.Article.GetArticlesByUser(uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func GetUserComments(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, totalPages, err := service.Comment.GetCommentsByUser(uint(id), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments, "total_pages": totalPages})
}
