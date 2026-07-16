package repository

import (
	"forum/models"
)

type CategoryRepository struct {
	*BaseRepository[models.Category]
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{
		BaseRepository: NewBaseRepository[models.Category](),
	}
}

func (r *CategoryRepository) GetCategoriesOrdered() ([]models.Category, error) {
	var categories []models.Category
	err := r.db.Order("sort_order ASC, created_at ASC").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) GetCategoryWithArticles(id uint) (*models.Category, error) {
	var category models.Category
	err := r.db.Preload("Articles").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}