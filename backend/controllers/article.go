package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		CategoryID  uint   `json:"category_id" binding:"required"`
		IsAnonymous bool   `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	article, err := service.Article.CreateArticle(userID, input.Title, input.Content, input.CategoryID, input.IsAnonymous)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "创建成功", "article": article})
}

func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID, _ := strconv.Atoi(c.Query("category_id"))

	articles, totalPages, err := service.Article.GetArticles(page, pageSize, uint(categoryID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func GetArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	article, err := service.Article.GetArticle(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, article)
}

func UpdateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Article.UpdateArticle(userID, uint(id), input.Title, input.Content, input.CategoryID)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "更新成功"})
}

func DeleteArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.DeleteArticle(userID, uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func RestoreArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.RestoreArticle(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "恢复成功"})
}

func LikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.LikeArticle(userID, uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

func UnlikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.UnlikeArticle(userID, uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}

func GetMyArticles(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, totalPages, err := service.Article.GetMyArticles(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func SearchArticles(c *gin.Context) {
	keyword := c.Query("q")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	articles, totalPages, err := service.Article.SearchArticles(keyword, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func ShareArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.ShareArticle(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "分享成功"})
}

func PinArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.PinArticle(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "置顶成功"})
}

func UnpinArticle(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.UnpinArticle(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消置顶成功"})
}

func GetDraftArticles(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	articles, totalPages, err := service.Article.GetDraftArticles(userID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func PublishDraft(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.PublishDraft(userID, uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功"})
}

func CoinArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Article.CoinArticle(userID, uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "投币成功"})
}
