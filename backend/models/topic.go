package models

import (
	"time"
)

// Topic 话题标签模型
type Topic struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	Name         string    `gorm:"uniqueIndex;size:50;not null" json:"name"` // 话题名称
	DisplayName string    `gorm:"size:100" json:"display_name"`             // 显示名称
	Description  string    `gorm:"size:500" json:"description"`              // 话题描述
	Icon         string    `gorm:"size:200" json:"icon"`                     // 话题图标
	CoverImage   string    `gorm:"size:500" json:"cover_image"`             // 封面图片
	ArticleCount int       `gorm:"default:0" json:"article_count"`          // 文章数量
	FollowCount  int       `gorm:"default:0" json:"follow_count"`           // 关注数量
	IsHot        bool      `gorm:"default:false;index" json:"is_hot"`       // 是否热门
	IsOfficial   bool      `gorm:"default:false" json:"is_official"`        // 是否官方
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Topic) TableName() string {
	return "topics"
}

// ArticleTopic 文章话题关联模型
type ArticleTopic struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	ArticleID uint      `gorm:"uniqueIndex:idx_article_topic;index" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID" json:"article"`
	TopicID   uint      `gorm:"uniqueIndex:idx_article_topic;index" json:"topic_id"`
	Topic     Topic     `gorm:"foreignKey:TopicID" json:"topic"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置表名
func (ArticleTopic) TableName() string {
	return "article_topics"
}

// TopicFollow 话题关注模型
type TopicFollow struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_topic;index" json:"user_id"`
	User      User      `gorm:"foreignKey:UserID" json:"user"`
	TopicID   uint      `gorm:"uniqueIndex:idx_user_topic;index" json:"topic_id"`
	Topic     Topic     `gorm:"foreignKey:TopicID" json:"topic"`
	CreatedAt time.Time `json:"created_at"`
}

// TableName 设置表名
func (TopicFollow) TableName() string {
	return "topic_follows"
}

// HotTopic 热门话题模型
type HotTopic struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	TopicID      uint      `gorm:"index" json:"topic_id"`
	Topic        Topic     `gorm:"foreignKey:TopicID" json:"topic"`
	Score        float64   `gorm:"default:0" json:"score"`              // 热度分数
	Rank         int       `gorm:"default:0" json:"rank"`               // 排名
	ViewCount    int       `gorm:"default:0" json:"view_count"`         // 浏览量
	DiscussCount int       `gorm:"default:0" json:"discuss_count"`      // 讨论数
	Period       string    `gorm:"size:20;index" json:"period"`          // 统计周期：daily, weekly, monthly
	Date         string    `gorm:"size:10;index" json:"date"`            // 统计日期
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 设置表名
func (HotTopic) TableName() string {
	return "hot_topics"
}