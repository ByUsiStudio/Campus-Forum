package controllers

import (
	"forum/service"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	users, totalPages, err := service.Admin.GetUsers(page, pageSize, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users, "total_pages": totalPages})
}

func GetUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	user, err := service.Admin.GetUser(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func UpdateUserRole(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Role string `json:"role" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Admin.UpdateUserRole(uint(id), input.Role)
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

func BanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Reason string `json:"reason" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := service.Admin.BanUser(uint(id), input.Reason)
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "封禁成功"})
}

func UnbanUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Admin.UnbanUser(uint(id))
	if err != nil {
		if appErr, ok := utils.IsAppError(err); ok {
			c.JSON(appErr.Code, gin.H{"error": appErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "解封成功"})
}

func DeleteUser(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Admin.DeleteUser(uint(id))
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

func GetArticlesAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	status := c.Query("status")

	articles, totalPages, err := service.Admin.GetArticlesAdmin(page, pageSize, status)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"articles": articles, "total_pages": totalPages})
}

func DeleteArticleAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Admin.DeleteArticleAdmin(uint(id))
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

func GetCommentsAdmin(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	comments, totalPages, err := service.Admin.GetCommentsAdmin(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments, "total_pages": totalPages})
}

func DeleteCommentAdmin(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	err := service.Admin.DeleteCommentAdmin(uint(id))
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

func GetLogs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	logs, totalPages, err := service.Admin.GetLogs(page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"logs": logs, "total_pages": totalPages})
}

func GetStats(c *gin.Context) {
	stats, err := service.Admin.GetStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, stats)
}
