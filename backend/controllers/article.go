package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"net/http"
	"strconv"
	_ "strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func getUserTitles(userID uint) []models.Title {
	var userTitles []models.UserTitle
	database.DB.Where("user_id = ?", userID).Preload("Title").Find(&userTitles)

	var titles []models.Title
	for _, ut := range userTitles {
		if ut.Title.IsActive {
			titles = append(titles, ut.Title)
		}
	}
	return titles
}

func CreateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		Title      string `json:"title" binding:"required"`
		Content    string `json:"content" binding:"required"`
		CategoryID uint   `json:"category_id" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 将Markdown转换为HTML
	contentHTML := utils.MarkdownToHTML(input.Content)

	article := models.Article{
		Title:       input.Title,
		Content:     input.Content,
		ContentHTML: contentHTML,
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Status:      "published",
	}

	if result := database.DB.Create(&article); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	// 预加载关联数据
	database.DB.Preload("User").Preload("Category").First(&article, article.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "创建成功",
		"article": article,
	})
}

func GetArticles(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	categoryID := c.Query("category_id")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	query := database.DB.Model(&models.Article{}).Where("status = ?", "published").Preload("User").Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	var articles []models.Article
	var total int64

	query.Count(&total)
	query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

	type ArticleResponse struct {
		models.Article
		User struct {
			ID          uint           `json:"id"`
			Username    string         `json:"username"`
			DisplayName string         `json:"display_name"`
			Avatar      string         `json:"avatar"`
			Titles      []models.Title `json:"titles"`
		} `json:"user"`
	}

	var response []ArticleResponse
	for _, article := range articles {
		titles := getUserTitles(article.UserID)
		resp := ArticleResponse{
			Article: article,
		}
		resp.User.ID = article.User.ID
		resp.User.Username = article.User.Username
		resp.User.DisplayName = article.User.DisplayName
		resp.User.Avatar = article.User.Avatar
		resp.User.Titles = titles
		response = append(response, resp)
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":    response,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article

	if result := database.DB.Preload("User").Preload("Category").Where("id = ? AND status = ?", id, "published").First(&article); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 浏览量校验：仅登录用户可增加浏览量，同一用户1分钟内只允许增加一浏览量
	clientIP := c.ClientIP()
	articleIDUint := uint(parseUint(id))

	canIncreaseView := false
	if userID, exists := c.Get("user_id"); exists {
		uid := userID.(uint)
		oneMinuteAgo := time.Now().Add(-time.Minute)

		var existingView models.ViewHistory
		if result := database.DB.Where("article_id = ? AND user_id = ? AND created_at > ?", articleIDUint, uid, oneMinuteAgo).First(&existingView); result.Error != nil {
			canIncreaseView = true
		}
	}

	if canIncreaseView {
		database.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)

		userIDPtr := new(uint)
		if userID, exists := c.Get("user_id"); exists {
			*userIDPtr = userID.(uint)
		}

		viewHistory := models.ViewHistory{
			ArticleID: articleIDUint,
			IP:        clientIP,
			UserID:    userIDPtr,
		}
		database.DB.Create(&viewHistory)

		article.ViewCount++
	}

	// 获取评论（包含回复）
	var comments []models.Comment
	database.DB.Preload("User").Where("article_id = ? AND parent_id IS NULL", article.ID).Order("created_at DESC").Find(&comments)

	// 获取每个评论的回复
	for i := range comments {
		var replies []models.Comment
		database.DB.Preload("User").Where("parent_id = ?", comments[i].ID).Order("created_at ASC").Find(&replies)
		comments[i].Replies = replies
	}

	// 获取文章作者头衔
	articleTitles := getUserTitles(article.UserID)

	// 获取评论用户头衔
	type CommentResponse struct {
		models.Comment
		User struct {
			ID          uint           `json:"id"`
			Username    string         `json:"username"`
			DisplayName string         `json:"display_name"`
			Avatar      string         `json:"avatar"`
			Titles      []models.Title `json:"titles"`
		} `json:"user"`
		Replies []CommentResponse `json:"replies"`
	}

	var commentResponses []CommentResponse
	for _, comment := range comments {
		commentTitles := getUserTitles(comment.UserID)
		resp := CommentResponse{
			Comment: comment,
		}
		resp.User.ID = comment.User.ID
		resp.User.Username = comment.User.Username
		resp.User.DisplayName = comment.User.DisplayName
		resp.User.Avatar = comment.User.Avatar
		resp.User.Titles = commentTitles

		// 处理回复
		for _, reply := range comment.Replies {
			replyTitles := getUserTitles(reply.UserID)
			replyResp := CommentResponse{
				Comment: reply,
			}
			replyResp.User.ID = reply.User.ID
			replyResp.User.Username = reply.User.Username
			replyResp.User.DisplayName = reply.User.DisplayName
			replyResp.User.Avatar = reply.User.Avatar
			replyResp.User.Titles = replyTitles
			resp.Replies = append(resp.Replies, replyResp)
		}
		commentResponses = append(commentResponses, resp)
	}

	// 获取当前用户是否点赞
	var liked bool
	if userID, exists := c.Get("user_id"); exists {
		var like models.Like
		if database.DB.Where("user_id = ? AND article_id = ?", userID, article.ID).First(&like).Error == nil {
			liked = true
		}
	}

	// 获取当前用户对每条评论的点赞状态
	var commentLikedMap = make(map[uint]bool)
	if userID, exists := c.Get("user_id"); exists {
		for _, comment := range commentResponses {
			var commentLike models.CommentLike
			if database.DB.Where("user_id = ? AND comment_id = ?", userID, comment.ID).First(&commentLike).Error == nil {
				commentLikedMap[comment.ID] = true
			}
			// 检查回复的点赞状态
			for _, reply := range comment.Replies {
				if database.DB.Where("user_id = ? AND comment_id = ?", userID, reply.ID).First(&commentLike).Error == nil {
					commentLikedMap[reply.ID] = true
				}
			}
		}
	}

	type ArticleResponse struct {
		models.Article
		User struct {
			ID          uint           `json:"id"`
			Username    string         `json:"username"`
			DisplayName string         `json:"display_name"`
			Avatar      string         `json:"avatar"`
			Titles      []models.Title `json:"titles"`
		} `json:"user"`
	}

	articleResp := ArticleResponse{
		Article: article,
	}
	articleResp.User.ID = article.User.ID
	articleResp.User.Username = article.User.Username
	articleResp.User.DisplayName = article.User.DisplayName
	articleResp.User.Avatar = article.User.Avatar
	articleResp.User.Titles = articleTitles

	c.JSON(http.StatusOK, gin.H{
		"article":       articleResp,
		"comments":      commentResponses,
		"liked":         liked,
		"comment_liked": commentLikedMap,
	})
}

func UpdateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 检查权限
	if article.UserID != userID && c.GetString("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改"})
		return
	}

	var input struct {
		Title      string `json:"title"`
		Content    string `json:"content"`
		CategoryID uint   `json:"category_id"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if input.Title != "" {
		article.Title = input.Title
	}
	if input.Content != "" {
		article.Content = input.Content
		article.ContentHTML = utils.MarkdownToHTML(input.Content)
	}
	if input.CategoryID != 0 {
		article.CategoryID = input.CategoryID
	}

	database.DB.Save(&article)
	database.DB.Preload("User").Preload("Category").First(&article, article.ID)

	c.JSON(http.StatusOK, gin.H{
		"message": "更新成功",
		"article": article,
	})
}

func DeleteArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	// 如果是管理员或系统管理员，直接删除
	if role == "admin" || role == "system" {
		database.DB.Delete(&article)
		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
		return
	}

	// 普通用户需要提交删除申请
	var input struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	deletionRequest := models.DeletionRequest{
		ArticleID: article.ID,
		UserID:    userID,
		Reason:    input.Reason,
		Status:    "pending",
	}

	database.DB.Create(&deletionRequest)

	c.JSON(http.StatusOK, gin.H{"message": "删除申请已提交，等待管理员审核"})
}

func LikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID := c.Param("id")

	// 检查是否已点赞
	var existingLike models.Like
	if result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&existingLike); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经点过赞了"})
		return
	}

	like := models.Like{
		UserID:    userID,
		ArticleID: uint(parseUint(articleID)),
	}

	database.DB.Create(&like)
	database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("like_count", database.DB.Raw("like_count + 1"))

	// 获取最新点赞数
	var article models.Article
	database.DB.First(&article, articleID)

	c.JSON(http.StatusOK, gin.H{
		"message":    "点赞成功",
		"like_count": article.LikeCount,
	})
}

func UnlikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID := c.Param("id")

	result := database.DB.Where("user_id = ? AND article_id = ?", userID, uint(parseUint(articleID))).Delete(&models.Like{})
	if result.RowsAffected > 0 {
		database.DB.Model(&models.Article{}).Where("id = ?", uint(parseUint(articleID))).UpdateColumn("like_count", database.DB.Raw("like_count - 1"))
	}

	var article models.Article
	database.DB.First(&article, articleID)

	c.JSON(http.StatusOK, gin.H{
		"message":    "取消点赞",
		"like_count": article.LikeCount,
	})
}

func CreateComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID := c.Param("id")

	var input struct {
		Content  string `json:"content" binding:"required"`
		ParentID *uint  `json:"parent_id"` // 回复的评论ID
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 如果是回复，检查父评论是否存在
	var parentCommentUserID uint
	if input.ParentID != nil {
		var parentComment models.Comment
		if result := database.DB.First(&parentComment, *input.ParentID); result.Error != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "父评论不存在"})
			return
		}
		// 记录父评论用户ID用于发送通知
		parentCommentUserID = parentComment.UserID
		// 增加父评论的回复数量
		database.DB.Model(&parentComment).UpdateColumn("reply_count", parentComment.ReplyCount+1)
	}

	comment := models.Comment{
		Content:   input.Content,
		UserID:    userID,
		ArticleID: uint(parseUint(articleID)),
		ParentID:  input.ParentID,
	}

	database.DB.Create(&comment)
	database.DB.Preload("User").First(&comment, comment.ID)

	// 如果是回复且回复的不是自己，发送通知
	if input.ParentID != nil && parentCommentUserID != userID {
		replyNotification := models.CommentReplyNotification{
			UserID:    parentCommentUserID,
			CommentID: *input.ParentID,
			ReplyID:   comment.ID,
			ArticleID: uint(parseUint(articleID)),
		}
		database.DB.Create(&replyNotification)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "评论成功",
		"comment": comment,
	})
}

// LikeComment 点赞评论
func LikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentIDStr := c.Param("id")
	commentID := uint(parseUint(commentIDStr))

	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	var existingLike models.CommentLike
	if result := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).First(&existingLike); result.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "已经点过赞了"})
		return
	}

	like := models.CommentLike{
		UserID:    userID,
		CommentID: commentID,
	}

	database.DB.Create(&like)
	database.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount+1)

	c.JSON(http.StatusOK, gin.H{
		"message":    "点赞成功",
		"like_count": comment.LikeCount + 1,
	})
}

// UnlikeComment 取消评论点赞
func UnlikeComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	commentIDStr := c.Param("id")
	commentID := uint(parseUint(commentIDStr))

	result := database.DB.Where("user_id = ? AND comment_id = ?", userID, commentID).Delete(&models.CommentLike{})
	if result.RowsAffected > 0 {
		var comment models.Comment
		database.DB.First(&comment, commentID)
		database.DB.Model(&comment).UpdateColumn("like_count", comment.LikeCount-1)

		c.JSON(http.StatusOK, gin.H{
			"message":    "取消点赞",
			"like_count": comment.LikeCount - 1,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "未点赞"})
}

func DeleteComment(c *gin.Context) {
	userID := c.GetUint("user_id")
	role := c.GetString("role")
	commentID := c.Param("id")

	var comment models.Comment
	if result := database.DB.First(&comment, commentID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "评论不存在"})
		return
	}

	// 检查权限
	if comment.UserID != userID && role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限删除"})
		return
	}

	// 如果是回复评论，减少父评论的回复数量
	if comment.ParentID != nil {
		database.DB.Model(&models.Comment{}).Where("id = ?", *comment.ParentID).UpdateColumn("reply_count", gorm.Expr("reply_count - 1"))
	}

	database.DB.Delete(&comment)
	c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}
