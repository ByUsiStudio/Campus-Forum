package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type TopicService struct{}

var Topic = &TopicService{}

func (s *TopicService) CreateTopic(name, description string) error {
	topic := models.Topic{
		Name:        name,
		Description: description,
	}

	if result := database.DB.Create(&topic); result.Error != nil {
		return utils.NewError("创建话题失败", 500)
	}
	return nil
}

func (s *TopicService) UpdateTopic(topicID uint, name, description string) error {
	var topic models.Topic
	if result := database.DB.First(&topic, topicID); result.Error != nil {
		return utils.NewError("话题不存在", 404)
	}

	topic.Name = name
	topic.Description = description
	database.DB.Save(&topic)
	return nil
}

func (s *TopicService) DeleteTopic(topicID uint) error {
	var topic models.Topic
	if result := database.DB.First(&topic, topicID); result.Error != nil {
		return utils.NewError("话题不存在", 404)
	}

	database.DB.Delete(&topic)
	return nil
}

func (s *TopicService) FollowTopic(userID, topicID uint) error {
	var follow models.TopicFollow
	result := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).First(&follow)
	if result.Error == nil {
		return utils.NewError("已关注该话题", 400)
	}

	follow = models.TopicFollow{
		UserID:  userID,
		TopicID: topicID,
	}
	database.DB.Create(&follow)

	return nil
}

func (s *TopicService) UnfollowTopic(userID, topicID uint) error {
	result := database.DB.Where("user_id = ? AND topic_id = ?", userID, topicID).Delete(&models.TopicFollow{})
	if result.RowsAffected == 0 {
		return utils.NewError("未关注该话题", 400)
	}

	return nil
}

func (s *TopicService) GetFollowedTopics(userID uint) ([]models.TopicFollow, error) {
	var follows []models.TopicFollow
	err := database.DB.Where("user_id = ?", userID).Preload("Topic").Find(&follows).Error
	return follows, err
}

func (s *TopicService) AddTopicToArticle(articleID, topicID uint) error {
	var articleTopic models.ArticleTopic
	result := database.DB.Where("article_id = ? AND topic_id = ?", articleID, topicID).First(&articleTopic)
	if result.Error == nil {
		return utils.NewError("文章已关联该话题", 400)
	}

	articleTopic = models.ArticleTopic{
		ArticleID: articleID,
		TopicID:   topicID,
	}
	database.DB.Create(&articleTopic)

	database.DB.Model(&models.Topic{}).Where("id = ?", topicID).UpdateColumn("article_count", utils.Increment("article_count"))

	return nil
}

func (s *TopicService) RemoveTopicFromArticle(articleID, topicID uint) error {
	result := database.DB.Where("article_id = ? AND topic_id = ?", articleID, topicID).Delete(&models.ArticleTopic{})
	if result.RowsAffected == 0 {
		return utils.NewError("文章未关联该话题", 400)
	}

	database.DB.Model(&models.Topic{}).Where("id = ?", topicID).UpdateColumn("article_count", utils.Decrement("article_count"))

	return nil
}

func (s *TopicService) GetTopics() ([]models.Topic, error) {
	var topics []models.Topic
	err := database.DB.Order("created_at DESC").Find(&topics).Error
	return topics, err
}

func (s *TopicService) GetTopic(id uint) (*models.Topic, error) {
	var topic models.Topic
	err := database.DB.First(&topic, id).Error
	if err != nil {
		return nil, utils.NewError("话题不存在", 404)
	}
	return &topic, nil
}

func (s *TopicService) GetHotTopics(limit int) ([]models.Topic, error) {
	var topics []models.Topic
	err := database.DB.Order("article_count DESC").Limit(limit).Find(&topics).Error
	return topics, err
}

func (s *TopicService) GetArticlesByTopic(topicID uint, page, pageSize int) ([]models.ArticleTopic, int, error) {
	var articleTopics []models.ArticleTopic
	var total int64

	query := database.DB.Model(&models.ArticleTopic{}).Where("topic_id = ?", topicID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("Article").Preload("Article.User").Preload("Article.Category").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&articleTopics).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return articleTopics, totalPages, err
}
