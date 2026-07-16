package service

import (
	"encoding/json"
	"forum/database"
	"forum/models"
	"forum/utils"
)

type AdminConfigService struct{}

var AdminConfig = &AdminConfigService{}

func (s *AdminConfigService) GetDeletionRequests(page, pageSize int) ([]models.DeletionRequest, int, error) {
	var requests []models.DeletionRequest
	var total int64

	query := database.DB.Model(&models.DeletionRequest{}).Where("status = ?", "pending")
	query.Count(&total)

	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	offset := (page - 1) * pageSize
	err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).Limit(pageSize).
		Find(&requests).Error

	totalPages := int(total / int64(pageSize))
	if total%int64(pageSize) > 0 {
		totalPages++
	}

	return requests, totalPages, err
}

func (s *AdminConfigService) ApproveDeletion(requestID uint) error {
	var request models.DeletionRequest
	if result := database.DB.First(&request, requestID); result.Error != nil {
		return utils.NewError("请求不存在", 404)
	}

	request.Status = "approved"
	database.DB.Save(&request)

	var user models.User
	database.DB.First(&user, request.UserID)
	user.Status = "deleted"
	database.DB.Save(&user)

	return nil
}

func (s *AdminConfigService) RejectDeletion(requestID uint) error {
	var request models.DeletionRequest
	if result := database.DB.First(&request, requestID); result.Error != nil {
		return utils.NewError("请求不存在", 404)
	}

	request.Status = "rejected"
	database.DB.Save(&request)

	return nil
}

func (s *AdminConfigService) GetAnnouncement() (*models.Announcement, error) {
	var announcement models.Announcement
	err := database.DB.Order("created_at DESC").First(&announcement).Error
	if err != nil {
		return nil, utils.NewError("公告不存在", 404)
	}
	return &announcement, nil
}

func (s *AdminConfigService) UpdateAnnouncement(content string) error {
	announcement := models.Announcement{
		Content: content,
	}
	database.DB.Create(&announcement)
	return nil
}

func (s *AdminConfigService) GetSiteConfig() (*models.SiteConfig, error) {
	var config models.SiteConfig
	err := database.DB.Order("created_at DESC").First(&config).Error
	if err != nil {
		config = models.SiteConfig{
			SiteTitle:       "校园论坛",
			SiteDescription: "校园论坛",
			SiteLogo:        "",
		}
		database.DB.Create(&config)
		return &config, nil
	}
	return &config, nil
}

func (s *AdminConfigService) UpdateSiteConfig(config models.SiteConfig) error {
	database.DB.Create(&config)
	return nil
}

func (s *AdminConfigService) GetSidebarConfig() (*models.SidebarConfig, error) {
	var config models.SidebarConfig
	err := database.DB.Order("created_at DESC").First(&config).Error
	if err != nil {
		return nil, utils.NewError("侧边栏配置不存在", 404)
	}
	return &config, nil
}

func (s *AdminConfigService) UpdateSidebarConfig(items []models.SidebarItem) error {
	itemsJSON, _ := json.Marshal(items)
	config := models.SidebarConfig{
		Items: string(itemsJSON),
	}
	database.DB.Create(&config)
	return nil
}
