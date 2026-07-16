package repository

import (
	"forum/models"
)

type TitleRepository struct {
	*BaseRepository[models.Title]
}

func NewTitleRepository() *TitleRepository {
	return &TitleRepository{
		BaseRepository: NewBaseRepository[models.Title](),
	}
}

type UserTitleRepository struct {
	*BaseRepository[models.UserTitle]
}

func NewUserTitleRepository() *UserTitleRepository {
	return &UserTitleRepository{
		BaseRepository: NewBaseRepository[models.UserTitle](),
	}
}

func (r *UserTitleRepository) GetUserTitles(userID uint) ([]models.UserTitle, error) {
	var userTitles []models.UserTitle
	err := r.db.Where("user_id = ?", userID).Preload("Title").Order("granted_at DESC").Find(&userTitles).Error
	return userTitles, err
}

func (r *UserTitleRepository) CheckUserTitle(userID, titleID uint) (bool, error) {
	var count int64
	err := r.db.Model(&models.UserTitle{}).Where("user_id = ? AND title_id = ?", userID, titleID).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
