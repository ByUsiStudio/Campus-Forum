package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
)

type CategoryService struct{}

var Category = &CategoryService{}

func (s *CategoryService) CreateCategory(name, description string, icon string) (*models.Category, error) {
	category := models.Category{
		Name:        name,
		Description: description,
		SortOrder:   0,
	}

	if result := database.DB.Create(&category); result.Error != nil {
		return nil, utils.NewError("创建分区失败", 500)
	}
	return &category, nil
}

func (s *CategoryService) UpdateCategory(categoryID uint, name, description string, icon string) error {
	var category models.Category
	if result := database.DB.First(&category, categoryID); result.Error != nil {
		return utils.NewError("分区不存在", 404)
	}

	category.Name = name
	category.Description = description
	database.DB.Save(&category)
	return nil
}

func (s *CategoryService) DeleteCategory(categoryID uint) error {
	var category models.Category
	if result := database.DB.First(&category, categoryID); result.Error != nil {
		return utils.NewError("分区不存在", 404)
	}

	database.DB.Delete(&category)
	return nil
}

func (s *CategoryService) GetCategories() ([]models.Category, error) {
	var categories []models.Category
	err := database.DB.Order("sort_order ASC, created_at ASC").Find(&categories).Error
	return categories, err
}

func (s *CategoryService) GetCategory(categoryID uint) (*models.Category, error) {
	var category models.Category
	if result := database.DB.First(&category, categoryID); result.Error != nil {
		return nil, utils.NewError("分区不存在", 404)
	}
	return &category, nil
}
