package repository

import (
	"forum/models"
)

type DeletionRequestRepository struct {
	*BaseRepository[models.DeletionRequest]
}

func NewDeletionRequestRepository() *DeletionRequestRepository {
	return &DeletionRequestRepository{
		BaseRepository: NewBaseRepository[models.DeletionRequest](),
	}
}

func (r *DeletionRequestRepository) GetPendingRequests(page, pageSize int) ([]models.DeletionRequest, int64, error) {
	var requests []models.DeletionRequest
	var total int64

	query := r.db.Model(&models.DeletionRequest{}).Where("status = ?", "pending")
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

	return requests, total, err
}

type AnnouncementRepository struct {
	*BaseRepository[models.Announcement]
}

func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{
		BaseRepository: NewBaseRepository[models.Announcement](),
	}
}

func (r *AnnouncementRepository) GetActiveAnnouncement() (*models.Announcement, error) {
	var announcement models.Announcement
	err := r.db.Order("created_at DESC").First(&announcement).Error
	if err != nil {
		return nil, err
	}
	return &announcement, nil
}

type SiteConfigRepository struct {
	*BaseRepository[models.SiteConfig]
}

func NewSiteConfigRepository() *SiteConfigRepository {
	return &SiteConfigRepository{
		BaseRepository: NewBaseRepository[models.SiteConfig](),
	}
}

func (r *SiteConfigRepository) GetActiveConfig() (*models.SiteConfig, error) {
	var config models.SiteConfig
	err := r.db.Order("created_at DESC").First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type SidebarConfigRepository struct {
	*BaseRepository[models.SidebarConfig]
}

func NewSidebarConfigRepository() *SidebarConfigRepository {
	return &SidebarConfigRepository{
		BaseRepository: NewBaseRepository[models.SidebarConfig](),
	}
}

func (r *SidebarConfigRepository) GetActiveConfig() (*models.SidebarConfig, error) {
	var config models.SidebarConfig
	err := r.db.Order("created_at DESC").First(&config).Error
	if err != nil {
		return nil, err
	}
	return &config, nil
}
