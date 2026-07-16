package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		ArticleID uint   `json:"article_id" binding:"required"`
		Content   string `json:"content" binding:"required"`
		ParentID  uint   `json:"parent_id"`
		IsReply   bool   `json:"is_reply"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Comment.CreateComment(userID, input.ArticleID, input.Content, input.ParentID, input.IsReply)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "评论成功"})
}

func GetComments(c *gin.Context) {
	articleID, _ := strconv.Atoi(c.Param("article_id"))
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, totalPages, err := service.Comment.GetComments(uint(articleID), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments, "total_pages": totalPages})
}

func UpdateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Comment.UpdateComment(userID, uint(id), input.Content)
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

func DeleteComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Comment.DeleteComment(userID, uint(id))
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

func LikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Comment.LikeComment(userID, uint(id))
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

func UnlikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Comment.UnlikeComment(userID, uint(id))
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
