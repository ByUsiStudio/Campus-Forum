package controllers

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetTopics 获取话题列表
func GetTopics(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	isHot := c.Query("is_hot")

	offset := (page - 1) * limit

	var topics []models.Topic
	query := database.DB.Order("article_count desc")

	if isHot == "true" {
		query = query.Where("is_hot = ?", true)
	}

	var total int64
	database.DB.Model(&models.Topic{}).Count(&total)
	query.Offset(offset).Limit(limit).Find(&topics)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"topics": topics,
			"total":  total,
			"page":   page,
			"limit":  limit,
		},
	})
}

// GetTopic 获取话题详情
func GetTopic(c *gin.Context) {
	topicID := c.Param("id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	offset := (page - 1) * limit

	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "话题不存在")
		return
	}

	var articles []models.ArticleTopic
	var total int64
	database.DB.Model(&models.ArticleTopic{}).Where("topic_id = ?", topicID).Count(&total)
	database.DB.Where("topic_id = ?", topicID).
		Preload("Article").
		Preload("Article.User").
		Offset(offset).Limit(limit).
		Find(&articles)

	c.JSON(200, gin.H{
		"success": true,
		"data": gin.H{
			"topic":    topic,
			"articles": articles,
			"total":    total,
			"page":     page,
			"limit":    limit,
		},
	})
}

// GetHotTopics 获取热门话题
func GetHotTopics(c *gin.Context) {
	var topics []models.Topic
	database.DB.Where("is_hot = ?", true).
		Order("article_count desc").
		Limit(10).
		Find(&topics)

	c.JSON(200, gin.H{
		"success": true,
		"data":    topics,
	})
}

// CreateTopic 创建话题（管理员）
func CreateTopic(c *gin.Context) {
	var topic models.Topic
	if err := c.ShouldBindJSON(&topic); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	if err := database.DB.Create(&topic).Error; err != nil {
		utils.SendErrorResponse(c, 500, "创建失败")
		return
	}

	c.JSON(200, gin.H{
		"success": true,
		"data":    topic,
	})
}

// UpdateTopic 更新话题（管理员）
func UpdateTopic(c *gin.Context) {
	topicID := c.Param("id")

	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "话题不存在")
		return
	}

	if err := c.ShouldBindJSON(&topic); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	database.DB.Save(&topic)

	c.JSON(200, gin.H{
		"success": true,
		"data":    topic,
	})
}

// DeleteTopic 删除话题（管理员）
func DeleteTopic(c *gin.Context) {
	topicID := c.Param("id")

	// 删除所有关联
	database.DB.Where("topic_id = ?", topicID).Delete(&models.ArticleTopic{})
	database.DB.Where("topic_id = ?", topicID).Delete(&models.TopicFollow{})

	// 删除话题
	database.DB.Delete(&models.Topic{}, topicID)

	c.JSON(200, gin.H{
		"success": true,
		"message": "删除成功",
	})
}

// FollowTopic 关注话题
func FollowTopic(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	topicID := c.Param("id")

	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "话题不存在")
		return
	}

	// 检查是否已关注
	var existing models.TopicFollow
	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&existing).Error; err == nil {
		utils.SendErrorResponse(c, 400, "已关注该话题")
		return
	}

	follow := models.TopicFollow{
		UserID:  userID,
		TopicID: topic.ID,
	}
	if err := database.DB.Create(&follow).Error; err != nil {
		utils.SendErrorResponse(c, 500, "关注失败")
		return
	}

	// 更新关注数
	database.DB.Model(&topic).UpdateColumn("follow_count", topic.FollowCount+1)

	c.JSON(200, gin.H{
		"success": true,
		"message": "关注成功",
	})
}

// UnfollowTopic 取消关注话题
func UnfollowTopic(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	topicID := c.Param("id")

	if err := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).Delete(&models.TopicFollow{}).Error; err != nil {
		utils.SendErrorResponse(c, 500, "取消关注失败")
		return
	}

	// 更新关注数
	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err == nil {
		database.DB.Model(&topic).UpdateColumn("follow_count", topic.FollowCount-1)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "取消关注成功",
	})
}

// GetFollowedTopics 获取用户关注的话题
func GetFollowedTopics(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)

	var follows []models.TopicFollow
	database.DB.Where("user_id = ?", userID).
		Preload("Topic").
		Find(&follows)

	c.JSON(200, gin.H{
		"success": true,
		"data":    follows,
	})
}

// AddTopicToArticle 为文章添加话题
func AddTopicToArticle(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	articleID := c.Param("article_id")

	var req struct {
		TopicID uint `json:"topic_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.SendErrorResponse(c, 400, "参数错误")
		return
	}

	// 验证文章所有权
	var article models.Article
	if err := database.DB.Where("id = ? AND user_id = ?", articleID, userID).First(&article).Error; err != nil {
		utils.SendErrorResponse(c, 404, "文章不存在或无权限")
		return
	}

	// 验证话题
	var topic models.Topic
	if err := database.DB.First(&topic, req.TopicID).Error; err != nil {
		utils.SendErrorResponse(c, 404, "话题不存在")
		return
	}

	// 检查是否已添加
	var existing models.ArticleTopic
	if err := database.DB.Where("article_id = ? AND topic_id = ?", articleID, req.TopicID).First(&existing).Error; err == nil {
		utils.SendErrorResponse(c, 400, "文章已添加该话题")
		return
	}

	topicArticle := models.ArticleTopic{
		ArticleID: article.ID,
		TopicID:   topic.ID,
	}
	if err := database.DB.Create(&topicArticle).Error; err != nil {
		utils.SendErrorResponse(c, 500, "添加失败")
		return
	}

	// 更新话题文章数
	database.DB.Model(&topic).UpdateColumn("article_count", topic.ArticleCount+1)

	c.JSON(200, gin.H{
		"success": true,
		"message": "添加成功",
	})
}

// RemoveTopicFromArticle 从文章移除话题
func RemoveTopicFromArticle(c *gin.Context) {
	userID := c.MustGet("user_id").(uint)
	articleID := c.Param("article_id")
	topicID := c.Param("topic_id")

	// 验证文章所有权
	var article models.Article
	if err := database.DB.Where("id = ? AND user_id = ?", articleID, userID).First(&article).Error; err != nil {
		utils.SendErrorResponse(c, 404, "文章不存在或无权限")
		return
	}

	if err := database.DB.Where("article_id = ? AND topic_id = ?", articleID, topicID).Delete(&models.ArticleTopic{}).Error; err != nil {
		utils.SendErrorResponse(c, 500, "移除失败")
		return
	}

	// 更新话题文章数
	var topic models.Topic
	if err := database.DB.First(&topic, topicID).Error; err == nil {
		database.DB.Model(&topic).UpdateColumn("article_count", topic.ArticleCount-1)
	}

	c.JSON(200, gin.H{
		"success": true,
		"message": "移除成功",
	})
}
