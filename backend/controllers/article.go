package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func maskAnonymousUser(article *models.Article, isOwner bool) {
	if article.IsAnonymous && !isOwner {
		article.User = models.User{
			ID:          0,
			Username:    "anonymous",
			DisplayName: "匿名用户",
			Avatar:      "",
		}
	}
}

func maskAnonymousUsers(articles *[]models.Article, currentUserID uint) {
	for i := range *articles {
		isOwner := (*articles)[i].UserID == currentUserID
		maskAnonymousUser(&(*articles)[i], isOwner)
	}
}

func CreateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		CategoryID  uint   `json:"category_id" binding:"required"`
		VoiceURL    string `json:"voice_url"`
		IsAnonymous bool   `json:"is_anonymous"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	contentHTML := utils.MarkdownToHTML(input.Content)

	article := models.Article{
		Title:       input.Title,
		Content:     input.Content,
		ContentHTML: contentHTML,
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Status:      "published",
		VoiceURL:    input.VoiceURL,
		IsAnonymous: input.IsAnonymous,
	}

	if result := database.DB.Create(&article); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	database.DB.Preload("User").Preload("Category").First(&article, article.ID)

	if !input.IsAnonymous {
		var follows []models.Follow
		database.DB.Where("following_id = ?", userID).Find(&follows)

		for _, follow := range follows {
			notification := models.FollowNotification{
				UserID:    follow.FollowerID,
				SenderID:  userID,
				ArticleID: article.ID,
				Type:      "new_article",
			}
			database.DB.Create(&notification)
		}
	}

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

	var currentUserID uint
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
	}

	for i := range articles {
		var likeCount int64
		database.DB.Model(&models.Like{}).Where("article_id = ?", articles[i].ID).Count(&likeCount)
		articles[i].LikeCount = int(likeCount)

		var commentCount int64
		database.DB.Model(&models.Comment{}).Where("article_id = ?", articles[i].ID).Count(&commentCount)
		articles[i].CommentCount = int(commentCount)
	}

	maskAnonymousUsers(&articles, currentUserID)

	c.JSON(http.StatusOK, gin.H{
		"articles":    articles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func GetArticle(c *gin.Context) {
	id := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var article models.Article
	if result := database.DB.Preload("User").Preload("Category").First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	database.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)

	var likeCount int64
	database.DB.Model(&models.Like{}).Where("article_id = ?", article.ID).Count(&likeCount)
	article.LikeCount = int(likeCount)

	var commentCount int64
	database.DB.Model(&models.Comment{}).Where("article_id = ?", article.ID).Count(&commentCount)
	article.CommentCount = int(commentCount)

	var comments []models.Comment
	var total int64
	database.DB.Model(&models.Comment{}).Where("article_id = ?", article.ID).Count(&total)
	database.DB.Where("article_id = ? AND parent_id = 0", article.ID).Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

	var currentUserID uint
	isOwner := false
	if userID, exists := c.Get("user_id"); exists {
		currentUserID = userID.(uint)
		isOwner = article.UserID == currentUserID
	}

	for i := range comments {
		var replyCount int64
		database.DB.Model(&models.Comment{}).Where("parent_id = ?", comments[i].ID).Count(&replyCount)
		comments[i].ReplyCount = int(replyCount)

		var replies []models.Comment
		database.DB.Where("parent_id = ?", comments[i].ID).Preload("User").Order("created_at ASC").Find(&replies)

		for j := range replies {
			isReplyOwner := replies[j].UserID == currentUserID
			if replies[j].IsAnonymous && !isReplyOwner {
				replies[j].User = models.User{
					ID:          0,
					Username:    "anonymous",
					DisplayName: "匿名用户",
					Avatar:      "",
				}
			}
		}
		comments[i].Replies = replies
	}

	var liked = false
	if userID, exists := c.Get("user_id"); exists {
		var like models.Like
		if result := database.DB.Where("user_id = ? AND article_id = ?", userID, id).First(&like); result.Error == nil {
			liked = true
		}
	}

	commentLiked := make(map[uint]bool)
	if userID, exists := c.Get("user_id"); exists {
		var likes []models.CommentLike
		database.DB.Where("user_id = ? AND comment_id IN (SELECT id FROM comments WHERE article_id = ?)", userID, article.ID).Find(&likes)
		for _, like := range likes {
			commentLiked[like.CommentID] = true
		}
	}

	maskAnonymousUser(&article, isOwner)

	for i := range comments {
		isCommentOwner := comments[i].UserID == currentUserID
		if comments[i].IsAnonymous && !isCommentOwner {
			comments[i].User = models.User{
				ID:          0,
				Username:    "anonymous",
				DisplayName: "匿名用户",
				Avatar:      "",
			}
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"article":       article,
		"comments":      comments,
		"total":         total,
		"page":          page,
		"page_size":     pageSize,
		"total_pages":   (total + int64(pageSize) - 1) / int64(pageSize),
		"liked":         liked,
		"comment_liked": commentLiked,
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

	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权修改此文章"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		CategoryID  *uint  `json:"category_id"`
		VoiceURL    string `json:"voice_url"`
		IsAnonymous *bool  `json:"is_anonymous"`
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
	if input.CategoryID != nil {
		article.CategoryID = *input.CategoryID
	}
	if input.VoiceURL != "" {
		article.VoiceURL = input.VoiceURL
	}
	if input.IsAnonymous != nil {
		article.IsAnonymous = *input.IsAnonymous
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

	if role == "admin" || role == "system" {
		database.DB.Delete(&article)
		c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
		return
	}

	var input struct {
		Reason string `json:"reason"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		input.Reason = "用户请求删除"
	}

	deletionRequest := models.DeletionRequest{
		ArticleID: article.ID,
		UserID:    userID,
		Reason:    input.Reason,
		Status:    "pending",
	}
	database.DB.Create(&deletionRequest)

	c.JSON(http.StatusOK, gin.H{"message": "删除申请已提交，等待审核"})
}

func LikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var article models.Article
	if result := database.DB.First(&article, articleID); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	var existingLike models.Like
	if result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&existingLike); result.Error == nil {
		c.JSON(http.StatusOK, gin.H{"message": "已经点赞过了"})
		return
	}

	like := models.Like{
		UserID:    uint(userID),
		ArticleID: uint(articleID),
	}

	if result := database.DB.Create(&like); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "点赞成功"})
}

func UnlikeArticle(c *gin.Context) {
	userID := c.GetUint("user_id")
	articleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的文章ID"})
		return
	}

	var like models.Like
	if result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).First(&like); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "未点赞此文章"})
		return
	}

	if result := database.DB.Delete(&like); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "取消点赞失败"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "取消点赞成功"})
}

func GetMyArticles(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var articles []models.Article
	var total int64

	database.DB.Model(&models.Article{}).Where("user_id = ?", userID).Count(&total)
	database.DB.Where("user_id = ?", userID).Preload("User").Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

	for i := range articles {
		var likeCount int64
		database.DB.Model(&models.Like{}).Where("article_id = ?", articles[i].ID).Count(&likeCount)
		articles[i].LikeCount = int(likeCount)

		var commentCount int64
		database.DB.Model(&models.Comment{}).Where("article_id = ?", articles[i].ID).Count(&commentCount)
		articles[i].CommentCount = int(commentCount)
	}

	c.JSON(http.StatusOK, gin.H{
		"articles":    articles,
		"total":       total,
		"page":        page,
		"page_size":   pageSize,
		"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

func parseUint(s string) uint {
	i, _ := strconv.ParseUint(s, 10, 64)
	return uint(i)
}
