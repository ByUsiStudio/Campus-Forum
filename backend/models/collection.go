package models

import (
	"time"
)

// Collection 收藏夹模型
type Collection struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	UserID       uint      `gorm:"index" json:"user_id"`
	User         User      `gorm:"foreignKey:UserID" json:"user"`
	Name         string    `gorm:"size:100;not null" json:"name"`  // 收藏夹名称
	Description  string    `gorm:"size:500" json:"description"`    // 收藏夹描述
	IsPublic     bool      `gorm:"default:false" json:"is_public"` // 是否公开
	CoverImage   string    `gorm:"size:500" json:"cover_image"`    // 封面图片
	ArticleCount int       `gorm:"default:0" json:"article_count"` // 文章数量
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TableName 设置表名
func (Collection) TableName() string {
	return "collections"
}

// CollectionArticle 收藏夹文章关联模型
type CollectionArticle struct {
	ID           uint       `gorm:"primarykey" json:"id"`
	CollectionID uint       `gorm:"uniqueIndex:idx_collection_article;index" json:"collection_id"`
	Collection   Collection `gorm:"foreignKey:CollectionID" json:"collection"`
	ArticleID    uint       `gorm:"uniqueIndex:idx_collection_article;index" json:"article_id"`
	Article      Article    `gorm:"foreignKey:ArticleID" json:"article"`
	UserID       uint       `gorm:"index" json:"user_id"` // 添加收藏的用户ID
	Note         string     `gorm:"size:500" json:"note"` // 收藏笔记
	CreatedAt    time.Time  `json:"created_at"`
}

// TableName 设置表名
func (CollectionArticle) TableName() string {
	return "collection_articles"
}

// ArticleVersion 文章版本历史模型
type ArticleVersion struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	ArticleID   uint      `gorm:"index" json:"article_id"`
	Article     Article   `gorm:"foreignKey:ArticleID" json:"article"`
	UserID      uint      `gorm:"index" json:"user_id"`
	User        User      `gorm:"foreignKey:UserID" json:"user"`
	Title       string    `gorm:"size:200" json:"title"`         // 版本标题
	Content     string    `gorm:"type:text" json:"content"`      // 版本内容
	ContentHTML string    `gorm:"type:text" json:"content_html"` // 版本HTML内容
	Version     int       `gorm:"default:1" json:"version"`      // 版本号
	ChangeLog   string    `gorm:"size:500" json:"change_log"`    // 变更日志
	CreatedAt   time.Time `json:"created_at"`
}

// TableName 设置表名
func (ArticleVersion) TableName() string {
	return "article_versions"
}
