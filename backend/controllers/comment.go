package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID := c.Param("id")

	var input struct {
		Content     string `json:"content" binding:"required"`
		ParentID    *uint  `json:"parent_id"`
		IsAnonymous bool   `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// XSS过滤：清理评论内容
	safeContent := utils.SanitizeHTML(input.Content)

	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	comment := models.Comment{
		Content:     safeContent,
		UserID:      userID,
		ArticleID:   article.ID,
		ParentID:    input.ParentID,
		IsAnonymous: input.IsAnonymous,
	}

	if result := database.DB.Create(&comment); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建评论失败"})
		return
	}

	if input.ParentID != nil {
		database.DB.Model(&models.Comment{}).Where("id = ?", *input.ParentID).UpdateColumn("reply_count", database.DB.Raw("reply_count + 1"))
	}

	database.DB.Preload("User").First(&comment, comment.ID)

	var commentResp models.Comment
	database.DB.Preload("User").First(&commentResp, comment.ID)

	if commentResp.IsAnonymous {
		commentResp.User = models.User{
			ID:       0,
			Username: "匿名用户",
			Avatar:   "",
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论成功",
		"comment": commentResp,
	})
}

func DeleteComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID := c.Param("id")
	role := c.GetString("role")

	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	if comment.UserID != userID && role != "admin" && role != "system" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除"})
		return
	}

	if comment.ParentID != nil {
		database.DB.Model(&models.Comment{}).Where("id = ?", *comment.ParentID).UpdateColumn("reply_count", database.DB.Raw("reply_count - 1"))
	}

	database.DB.Delete(&comment)

	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func LikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID := c.Param("id")

	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	var existingLike models.CommentLike
	if result := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingLike); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已点赞"})
		return
	}

	like := models.CommentLike{
		UserID:    userID,
		CommentID: comment.ID,
	}

	if result := database.DB.Create(&like); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	database.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount+1)

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

func UnlikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentID := c.Param("id")

	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	var existingLike models.CommentLike
	if result := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingLike); result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "未点赞"})
		return
	}

	database.DB.Delete(&existingLike)

	if comment.LikeCount > 0 {
		database.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount-1)
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}
