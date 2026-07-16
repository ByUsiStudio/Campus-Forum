package repository

import (
	"forum/models"
)

type TopicRepository struct {
	*BaseRepository[models.Topic]
}

func NewTopicRepository() *TopicRepository {
	return &TopicRepository{
		BaseRepository: NewBaseRepository[models.Topic](),
	}
}

func (r *TopicRepository) GetTopicsOrdered() ([]models.Topic, error) {
	var topics []models.Topic
	err := r.db.Order("created_at DESC").Find(&topics).Error
	return topics, err
}

func (r *TopicRepository) GetHotTopics(limit int) ([]models.Topic, error) {
	var topics []models.Topic
	err := r.db.Order("article_count DESC").Limit(limit).Find(&topics).Error
	return topics, err
}

func (r *TopicRepository) GetTopicWithArticles(id uint) (*models.Topic, error) {
	var topic models.Topic
	err := r.db.Preload("Articles").First(&topic, id).Error
	if err != nil {
		return nil, err
	}
	return &topic, nil
}

type TopicFollowRepository struct {
	*BaseRepository[models.TopicFollow]
}

func NewTopicFollowRepository() *TopicFollowRepository {
	return &TopicFollowRepository{
		BaseRepository: NewBaseRepository[models.TopicFollow](),
	}
}

func (r *TopicFollowRepository) CheckFollow(userID, topicID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.TopicFollow{}).Where("user_id = ? AND topic_id = ?", userID, topicID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *TopicFollowRepository) GetFollowedTopics(userID uint) ([]models.TopicFollow, error) {
	var follows []models.TopicFollow
	err := r.db.Where("user_id = ?", userID).Preload("Topic").Find(&follows).Error
	return follows, err
}

type ArticleTopicRepository struct {
	*BaseRepository[models.ArticleTopic]
}

func NewArticleTopicRepository() *ArticleTopicRepository {
	return &ArticleTopicRepository{
		BaseRepository: NewBaseRepository[models.ArticleTopic](),
	}
}

func (r *ArticleTopicRepository) GetTopicsByArticle(articleID uint) ([]models.ArticleTopic, error) {
	var articleTopics []models.ArticleTopic
	err := r.db.Where("article_id = ?", articleID).Preload("Topic").Find(&articleTopics).Error
	return articleTopics, err
}

func (r *ArticleTopicRepository) GetArticlesByTopic(topicID uint, page, pageSize int) ([]models.ArticleTopic, int64, error) {
	var articleTopics []models.ArticleTopic
	var total int64

	query := r.db.Model(&models.ArticleTopic{}).Where("topic_id = ?", topicID)
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

	return articleTopics, total, err
}