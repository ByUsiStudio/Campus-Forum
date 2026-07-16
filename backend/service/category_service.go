package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
)

type CategoryService struct{}

var Category = &CategoryService{}

func (s *CategoryService) CreateCategory(name, description string, sortOrder int) error {
	category := models.Category{
		Name:        name,
		Description: description,
		SortOrder:   sortOrder,
	}

	if result := database.DB.Create(&category); result.Error != nil {
		return utils.NewError("创建分区失败", 500)
	}
	return nil
}

func (s *CategoryService) UpdateCategory(categoryID uint, name, description string, sortOrder int) error {
	var category models.Category
	if result := database.DB.First(&category, categoryID); result.Error != nil {
		return utils.NewError("分区不存在", 404)
	}

	category.Name = name
	category.Description = description
	category.SortOrder = sortOrder
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