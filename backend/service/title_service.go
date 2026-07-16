package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
)

type TitleService struct{}

var Title = &TitleService{}

func (s *TitleService) GetAllTitles() ([]models.Title, error) {
	var titles []models.Title
	err := database.DB.Find(&titles).Error
	return titles, err
}

func (s *TitleService) CreateTitle(name, description string, color string) error {
	title := models.Title{
		Name:        name,
		Description: description,
		Color:       color,
	}

	if result := database.DB.Create(&title); result.Error != nil {
		return utils.NewError("创建头衔失败", 500)
	}
	return nil
}

func (s *TitleService) UpdateTitle(titleID uint, name, description string, color string) error {
	var title models.Title
	if result := database.DB.First(&title, titleID); result.Error != nil {
		return utils.NewError("头衔不存在", 404)
	}

	title.Name = name
	title.Description = description
	title.Color = color
	database.DB.Save(&title)
	return nil
}

func (s *TitleService) DeleteTitle(titleID uint) error {
	var title models.Title
	if result := database.DB.First(&title, titleID); result.Error != nil {
		return utils.NewError("头衔不存在", 404)
	}

	database.DB.Delete(&title)
	return nil
}

func (s *TitleService) GrantTitle(userID, titleID uint) error {
	var userTitle models.UserTitle
	result := database.DB.Where("user_id = ? AND title_id = ?", userID, titleID).First(&userTitle)
	if result.Error == nil {
		return utils.NewError("用户已拥有该头衔", 400)
	}

	userTitle = models.UserTitle{
		UserID:  userID,
		TitleID: titleID,
	}
	database.DB.Create(&userTitle)

	return nil
}

func (s *TitleService) RevokeTitle(userID, titleID uint) error {
	result := database.DB.Where("user_id = ? AND title_id = ?", userID, titleID).Delete(&models.UserTitle{})
	if result.RowsAffected == 0 {
		return utils.NewError("用户未拥有该头衔", 400)
	}

	return nil
}

func (s *TitleService) GetUserTitles(userID uint) ([]models.UserTitle, error) {
	var userTitles []models.UserTitle
	err := database.DB.Where("user_id = ?", userID).Preload("Title").Order("granted_at DESC").Find(&userTitles).Error
	return userTitles, err
}
