package service

import (
	"forum/database"
	"forum/models"
	"forum/utils"
	"math"
)

type LevelService struct{}

var Level = &LevelService{}

func (s *LevelService) GetUserLevel(userID uint) (*models.UserLevel, error) {
	var userLevel models.UserLevel
	result := database.DB.Where("user_id = ?", userID).First(&userLevel)
	if result.Error != nil {
		return nil, utils.NewError("用户等级不存在", 404)
	}
	return &userLevel, nil
}

func (s *LevelService) GetUserExperienceRecords(userID uint, page, pageSize int) ([]models.ExperienceRecord, int, error) {
	var records []models.ExperienceRecord
	var total int64

	query := database.DB.Model(&models.ExperienceRecord{}).Where("user_id = ?", userID)
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&records).Error

	totalPages := int(math.Ceil(float64(total) / float64(pageSize)))

	return records, totalPages, err
}

func (s *LevelService) GetUserAchievements(userID uint) ([]models.UserAchievement, error) {
	var achievements []models.UserAchievement
	err := database.DB.Where("user_id = ?", userID).Preload("Achievement").Order("achieved_at DESC").Find(&achievements).Error
	return achievements, err
}

func (s *LevelService) GetAllAchievements() ([]models.Achievement, error) {
	var achievements []models.Achievement
	err := database.DB.Find(&achievements).Error
	return achievements, err
}

func (s *LevelService) GetLevelConfig() ([]models.LevelConfig, error) {
	var configs []models.LevelConfig
	err := database.DB.Order("level ASC").Find(&configs).Error
	if err != nil || len(configs) == 0 {
		configs = []models.LevelConfig{
			{Level: 1, RequiredExp: 0, Name: "新手"},
			{Level: 2, RequiredExp: 100, Name: "初级会员"},
			{Level: 3, RequiredExp: 300, Name: "中级会员"},
			{Level: 4, RequiredExp: 500, Name: "高级会员"},
			{Level: 5, RequiredExp: 1000, Name: "资深会员"},
			{Level: 6, RequiredExp: 2000, Name: "精英会员"},
			{Level: 7, RequiredExp: 5000, Name: "专家"},
			{Level: 8, RequiredExp: 10000, Name: "大师"},
		}
		database.DB.Create(&configs)
	}
	return configs, nil
}

func (s *LevelService) CreateLevelConfig(config models.LevelConfig) error {
	if result := database.DB.Create(&config); result.Error != nil {
		return utils.NewError("创建等级配置失败", 500)
	}
	return nil
}

func (s *LevelService) UpdateLevelConfig(configID uint, level, requiredExp int, name string) error {
	var config models.LevelConfig
	if result := database.DB.First(&config, configID); result.Error != nil {
		return utils.NewError("等级配置不存在", 404)
	}

	config.Level = level
	config.RequiredExp = requiredExp
	config.Name = name
	database.DB.Save(&config)
	return nil
}

func (s *LevelService) CreateAchievement(achievement models.Achievement) error {
	if result := database.DB.Create(&achievement); result.Error != nil {
		return utils.NewError("创建成就失败", 500)
	}
	return nil
}

func (s *LevelService) UpdateAchievement(achievementID uint, name, description string, icon string, requirement int) error {
	var achievement models.Achievement
	if result := database.DB.First(&achievement, achievementID); result.Error != nil {
		return utils.NewError("成就不存在", 404)
	}

	achievement.Name = name
	achievement.Description = description
	achievement.Icon = icon
	achievement.Requirement = requirement
	database.DB.Save(&achievement)
	return nil
}

func (s *LevelService) DeleteAchievement(achievementID uint) error {
	var achievement models.Achievement
	if result := database.DB.First(&achievement, achievementID); result.Error != nil {
		return utils.NewError("成就不存在", 404)
	}

	database.DB.Delete(&achievement)
	return nil
}
