package controllers

import (
    "forum/database"
    "forum/models"
    "forum/utils"
    "net/http"
    "strconv"
    _"strings"

    "github.com/gin-gonic/gin"
)

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
    
    c.JSON(http.StatusOK, gin.H{
        "articles": articles,
        "total":    total,
        "page":     page,
        "page_size": pageSize,
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
    
    // 增加浏览量
    database.DB.Model(&article).UpdateColumn("view_count", article.ViewCount+1)
    
    // 获取评论
    var comments []models.Comment
    database.DB.Preload("User").Where("article_id = ?", article.ID).Order("created_at DESC").Find(&comments)
    
    // 获取当前用户是否点赞
    var liked bool
    if userID, exists := c.Get("user_id"); exists {
        var like models.Like
        if database.DB.Where("user_id = ? AND article_id = ?", userID, article.ID).First(&like).Error == nil {
            liked = true
        }
    }
    
    c.JSON(http.StatusOK, gin.H{
        "article":  article,
        "comments": comments,
        "liked":    liked,
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
    
    // 如果是管理员，直接删除
    if role == "admin" {
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
    
    result := database.DB.Where("user_id = ? AND article_id = ?", userID, articleID).Delete(&models.Like{})
    if result.RowsAffected > 0 {
        database.DB.Model(&models.Article{}).Where("id = ?", articleID).UpdateColumn("like_count", database.DB.Raw("like_count - 1"))
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
        Content string `json:"content" binding:"required"`
    }
    
    if err := c.ShouldBindJSON(&input); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    
    comment := models.Comment{
        Content:   input.Content,
        UserID:    userID,
        ArticleID: uint(parseUint(articleID)),
    }
    
    database.DB.Create(&comment)
    database.DB.Preload("User").First(&comment, comment.ID)
    
    c.JSON(http.StatusOK, gin.H{
        "message": "评论成功",
        "comment": comment,
    })
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
    
    database.DB.Delete(&comment)
    c.JSON(http.StatusOK, gin.H{"message": "删除成功"})
}

func parseUint(s string) uint {
    i, _ := strconv.ParseUint(s, 10, 64)
    return uint(i)
}