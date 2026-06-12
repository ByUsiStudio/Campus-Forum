package controllers

import (
	"fmt"
	"forum/database"
	"forum/models"
	"forum/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateArticle(c *gin.Context) {
	userID := c.GetUint("user_id")

	var input struct {
		Title       string `json:"title" binding:"required"`
		Content     string `json:"content" binding:"required"`
		CategoryID  uint   `json:"category_id" binding:"required"`
		VoiceURL    string `json:"voice_url"`
		IsAnonymous bool   `json:"is_anonymous"`
		Status      string `json:"status"` // draft, published
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 默认为发布状态
	status := "published"
	if input.Status == "draft" {
		status = "draft"
	}

	contentHTML := utils.MarkdownToHTML(input.Content)

	article := models.Article{
		Title:       input.Title,
		Content:     input.Content,
		ContentHTML: contentHTML,
		UserID:      userID,
		CategoryID:  input.CategoryID,
		Status:      status,
		VoiceURL:    input.VoiceURL,
		IsAnonymous: input.IsAnonymous,
	}

	if result := database.DB.Create(&article); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "创建文章失败"})
		return
	}

	database.DB.Preload("User").Preload("Category").First(&article, article.ID)

	// 只有发布的文章才发送通知给好友
	if status == "published" && !input.IsAnonymous {
		var friends []models.Friend
		database.DB.Where("friend_id = ? AND status = 1", userID).Find(&friends)

		for _, friend := range friends {
			// 使用 PersonalNotification 发送通知
			personalNotif := models.PersonalNotification{
				UserID:  friend.UserID,
				Type:    "new_article",
				Title:   "好友发布新文章",
				Content: fmt.Sprintf(`{"article_id": %d, "sender_id": %d}`, article.ID, userID),
			}
			database.DB.Create(&personalNotif)
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
	query.Order("is_pinned DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

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

	database.DB.Model(&article).UpdateColumn("view_count", gorm.Expr("view_count + 1"))

	var likeCount int64
	database.DB.Model(&models.Like{}).Where("article_id = ?", article.ID).Count(&likeCount)
	article.LikeCount = int(likeCount)

	var commentCount int64
	database.DB.Model(&models.Comment{}).Where("article_id = ?", article.ID).Count(&commentCount)
	article.CommentCount = int(commentCount)

	var comments []models.Comment
	var total int64
	database.DB.Model(&models.Comment{}).Where("article_id = ?", article.ID).Count(&total)
	database.DB.Where("article_id = ? AND parent_id IS NULL", article.ID).Preload("User").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&comments)

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

		// 递归获取嵌套回复
		comments[i].Replies = getNestedReplies(comments[i].ID, currentUserID, 0)
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
	if result := database.DB.Where("status != ?", "deleted").First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在或已删除"})
		return
	}

	// 权限检查：管理员可以删除任何文章，普通用户只能删除自己的文章
	if role != "admin" && role != "system" && article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "您没有权限删除此文章"})
		return
	}

	// 软删除：将状态设置为 deleted
	article.Status = "deleted"
	if err := database.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "删除失败"})
		return
	}

	// 记录操作日志
	logOperation(c, ActionDelete, ModuleArticle,
		"删除文章 - 文章ID: "+strconv.FormatUint(uint64(article.ID), 10)+", 标题: "+article.Title)

	c.JSON(http.StatusOK, gin.H{"message": "文章已删除"})
}

// RestoreArticle 恢复已删除的文章
func RestoreArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if article.Status != "deleted" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "文章未被删除，无需恢复"})
		return
	}

	article.Status = "published"
	if err := database.DB.Save(&article).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "恢复失败"})
		return
	}

	// 记录操作日志
	logOperation(c, ActionUpdate, ModuleArticle,
		"恢复文章 - 文章ID: "+strconv.FormatUint(uint64(article.ID), 10)+", 标题: "+article.Title)

	c.JSON(http.StatusOK, gin.H{"message": "文章已恢复"})
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

	database.DB.Model(&article).UpdateColumn("like_count", gorm.Expr("like_count + 1"))

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

	database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("like_count", gorm.Expr("like_count - 1"))

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

	database.DB.Model(&models.Article{}).Where("user_id = ? AND status != ?", userID, "deleted").Count(&total)
	database.DB.Where("user_id = ? AND status != ?", userID, "deleted").Preload("User").Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

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

func SearchArticles(c *gin.Context) {
	keyword := c.Query("keyword")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if keyword == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "搜索关键词不能为空"})
		return
	}

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize

	var articles []models.Article
	var total int64

	query := database.DB.Model(&models.Article{}).
		Where("status = ? AND (title LIKE ? OR content LIKE ?)", "published", "%"+keyword+"%", "%"+keyword+"%")

	query.Count(&total)
	query.Preload("User").Preload("Category").Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

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

func ShareArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	database.DB.Model(&article).UpdateColumn("share_count", gorm.Expr("share_count + 1"))

	c.JSON(http.StatusOK, gin.H{"message": "分享成功", "share_count": article.ShareCount + 1})
}

// 置顶文章（管理员）
func PinArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	now := time.Now()
	database.DB.Model(&article).Updates(map[string]interface{}{
		"is_pinned": true,
		"pinned_at": &now,
	})

	c.JSON(http.StatusOK, gin.H{"message": "置顶成功", "article": article})
}

// 取消置顶文章（管理员）
func UnpinArticle(c *gin.Context) {
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	database.DB.Model(&article).Updates(map[string]interface{}{
		"is_pinned": false,
		"pinned_at": nil,
	})

	c.JSON(http.StatusOK, gin.H{"message": "取消置顶成功", "article": article})
}

// 获取用户草稿列表
func GetDraftArticles(c *gin.Context) {
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

	database.DB.Model(&models.Article{}).Where("user_id = ? AND status = ?", userID, "draft").Count(&total)
	database.DB.Where("user_id = ? AND status = ?", userID, "draft").Preload("User").Preload("Category").Order("updated_at DESC").Offset(offset).Limit(pageSize).Find(&articles)

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

// 发布草稿
func PublishDraft(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")

	var article models.Article
	if result := database.DB.First(&article, id); result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "文章不存在"})
		return
	}

	if article.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权操作"})
		return
	}

	if article.Status != "draft" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "不是草稿状态"})
		return
	}

	// 更新为发布状态
	database.DB.Model(&article).Update("status", "published")

	// 发送通知给好友
	if !article.IsAnonymous {
		var friends []models.Friend
		database.DB.Where("friend_id = ? AND status = 1", userID).Find(&friends)

		for _, friend := range friends {
			// 使用 PersonalNotification 发送通知
			personalNotif := models.PersonalNotification{
				UserID:  friend.UserID,
				Type:    "new_article",
				Title:   "好友发布新文章",
				Content: fmt.Sprintf(`{"article_id": %d, "sender_id": %d}`, article.ID, userID),
			}
			database.DB.Create(&personalNotif)
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "发布成功", "article": article})
}

// getNestedReplies 递归获取嵌套回复
// maxDepth 限制最大嵌套深度，防止无限递归
func getNestedReplies(parentID uint, currentUserID uint, depth int) []models.Comment {
	// 限制最大嵌套深度为5层
	if depth >= 5 {
		return nil
	}

	var replies []models.Comment
	database.DB.Where("parent_id = ?", parentID).Preload("User").Order("created_at ASC").Find(&replies)

	for i := range replies {
		isReplyOwner := replies[i].UserID == currentUserID
		if replies[i].IsAnonymous && !isReplyOwner {
			replies[i].User = models.User{
				ID:          0,
				Username:    "anonymous",
				DisplayName: "匿名用户",
				Avatar:      "",
			}
		}

		// 递归获取子回复
		replies[i].Replies = getNestedReplies(replies[i].ID, currentUserID, depth+1)
	}

	return replies
}
