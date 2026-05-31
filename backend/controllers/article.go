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
		if article.IsAnonymous {
			resp.User.DisplayName = "匿名用户"
			resp.User.Avatar = ""
		} else {
			resp.User.DisplayName = article.User.DisplayName
			resp.User.Avatar = article.User.Avatar
		}
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

	var comments []models.Comment
	database.DB.Preload("User").Where("article_id = ? AND parent_id IS NULL", article.ID).Order("created_at DESC").Find(&comments)

	for i := range comments {
		var replies []models.Comment
		database.DB.Preload("User").Where("parent_id = ?", comments[i].ID).Order("created_at ASC").Find(&replies)
		comments[i].Replies = replies
	}

	articleTitles := getUserTitles(article.UserID)

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
		if comment.IsAnonymous {
			resp.User.DisplayName = "匿名用户"
			resp.User.Avatar = ""
			resp.User.Titles = []models.Title{}
		} else {
			resp.User.DisplayName = comment.User.DisplayName
			resp.User.Avatar = comment.User.Avatar
			resp.User.Titles = commentTitles
		}

		for _, reply := range comment.Replies {
			replyTitles := getUserTitles(reply.UserID)
			replyResp := CommentResponse{
				Comment: reply,
			}
			replyResp.User.ID = reply.User.ID
			replyResp.User.Username = reply.User.Username
			if reply.IsAnonymous {
				replyResp.User.DisplayName = "匿名用户"
				replyResp.User.Avatar = ""
				replyResp.User.Titles = []models.Title{}
			} else {
				replyResp.User.DisplayName = reply.User.DisplayName
				replyResp.User.Avatar = reply.User.Avatar
				replyResp.User.Titles = replyTitles
			}
			resp.Replies = append(resp.Replies, replyResp)
		}
		commentResponses = append(commentResponses, resp)
	}

	var liked bool
	if userID, exists := c.Get("user_id"); exists {
		var like models.Like
		if database.DB.Where("user_id = ? AND article_id = ?", userID, article.ID).First(&like).Error == nil {
			liked = true
		}
	}

	var commentLikedMap = make(map[uint]bool)
	if userID, exists := c.Get("user_id"); exists {
		for _, comment := range commentResponses {
			var commentLike models.CommentLike
			if database.DB.Where("user_id = ? AND comment_id = ?", userID, comment.ID).First(&commentLike).Error == nil {
				commentLikedMap[comment.ID] = true
			}
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
	if article.IsAnonymous {
		articleResp.User.DisplayName = "匿名用户"
		articleResp.User.Avatar = ""
	} else {
		articleResp.User.DisplayName = article.User.DisplayName
		articleResp.User.Avatar = article.User.Avatar
	}
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

	if article.UserID != userID && c.GetString("role") != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "无权限修改"})
		return
	}

	var input struct {
		Title       string `json:"title"`
		Content     string `json:"content"`
		CategoryID  uint   `json:"category_id"`
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
	if input.CategoryID != 0 {
		article.CategoryID = input.CategoryID
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
